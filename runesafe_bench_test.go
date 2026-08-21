package runesafe_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/cplieger/runesafe/v2"
)

// This file exists because runesafe sits on a per-log-line hot path: a sanitizer
// runs once per emitted attribute, so its cost is multiplied by an app's entire
// log volume, and the package's central cost property was nowhere verified until
// this landed.
//
// That property is that CLEAN text is sanitized without allocating. It is
// strings.Map's property, not this package's — Map builds no output buffer until
// the first rune actually changes, and returns the ORIGINAL string when none
// does — and Sanitize and SanitizeSingleLine are one strings.Map call each, so
// they inherit it. It matters because almost no log attribute contains a C0
// control or a bidi override, so the clean path is the overwhelmingly common one:
// routing every untrusted attribute through the sanitizer costs nothing at all
// for honest values. Measured here, on go1.27.0: zero allocations for a clean
// input from 64 bytes to 256 KiB, under both policies, ASCII and multi-byte
// alike.
//
// It is worth GATING rather than trusting because every way to lose it is
// ordinary: replacing strings.Map with a strings.Builder loop, copying the input
// before mapping it, normalising or re-decoding it first. None of those changes
// what the sanitizers RETURN, so every functional test in this package stays
// green while the emit path goes from zero allocations per attribute to one — the
// exact class of regression a test suite that only checks output cannot see.
//
// Two kinds of check here, doing different jobs:
//
//   - The Test functions GATE a measured claim at merge time. testing.AllocsPerRun
//     is exact for an all-or-nothing allocator, so a fast path asserts == 0 rather
//     than a threshold. For a path that legitimately allocates, the valuable
//     property is that its cost is PROPORTIONATE to the input and cannot be
//     amplified: those assert that the allocation count is CONSTANT across inputs
//     three orders of magnitude apart, which is what an upstream-controlled length
//     would otherwise break.
//   - The Benchmark functions feed the weekly benchmark tracker a trend series.
//     They are size-parameterised so an accidental super-linear scan — a rescan
//     per unsafe rune, a per-rune string concatenation — shows up as a jump
//     between sizes rather than a uniform slowdown that reads as runner noise.
//
// Every number quoted in a comment here was measured before the assertion around
// it was written, and every assertion was seen to fail with the property broken.

// sinkString and sinkBool consume benchmark and AllocsPerRun results so the
// compiler cannot delete the call being measured. They are shared by the other
// benchmark files in this package; slog's Value has its own sink beside the
// Untrusted benchmarks, because only that file needs the import.
var (
	sinkString string
	sinkBool   bool
)

// Fixture sizes. typicalBytes is the shape this library actually meets — an
// upstream title or error message headed for one log attribute — and the scaling
// set spans it either way so a size series can expose a super-linear scan.
const (
	typicalBytes = 4096
	smallBytes   = 256
	largeBytes   = 65536
)

// repeatTo builds a deterministic fixture of exactly n bytes by repeating unit
// and padding the remainder with spaces. The padding rune is deliberately a
// space: it is one byte, valid UTF-8, and safe under both policies, so it can
// neither split a multi-byte rune at the tail nor change which class the fixture
// belongs to. Fixtures are built in a benchmark's setup, never inside its timed
// loop, and depend on nothing outside this file — the weekly tracker runs
// benchmarks with -run='^$', so no test function has executed first.
func repeatTo(unit string, n int) string {
	if n <= 0 {
		return ""
	}
	var b strings.Builder
	b.Grow(n)
	for b.Len()+len(unit) <= n {
		b.WriteString(unit)
	}
	for b.Len() < n {
		b.WriteByte(' ')
	}
	return b.String()
}

// cleanASCII is the fast path: a release-name-shaped attribute carrying nothing
// the policy touches, which is what the overwhelming majority of real log
// attributes look like.
func cleanASCII(n int) string {
	return repeatTo("Release.Name.S01E02.1080p.WEB-DL.x264-GROUP ", n)
}

// cleanUTF8 is the fast path again, over multi-byte runes. It is a separate class
// because strings.Map's no-change check runs per RUNE while the scan advances by
// bytes, so a regression that mishandles multi-byte input can be invisible on
// ASCII.
func cleanUTF8(n int) string {
	return repeatTo("Ясность — 日本語のタイトル — émission ", n)
}

// unsafeTail is the worst case for the fast path's bail-out: strings.Map scans
// every rune before it finds the one that changes, then still copies the whole
// prefix. A value shaped like this pays both the full scan and the full copy, so
// it bounds the class from above — and it is the realistic hostile shape, since
// an attacker appending one escape introducer to an otherwise honest title
// produces exactly it.
func unsafeTail(n int) string {
	if n <= 1 {
		return "\x1b"
	}
	return cleanASCII(n-1) + "\x1b"
}

// unsafeDense is text the sanitizer must rewrite throughout: a C0 control, a C1
// control, a bidi override and a line separator per unit, which is what a
// Trojan-Source-shaped title carries. Every rune the policy replaces here is
// multi-byte except the C0, so the rewritten form is SHORTER than the input,
// which is why this class never grows the output buffer.
func unsafeDense(n int) string {
	return repeatTo("a\x1bb\u009cc\u202ed\u2028e ", n)
}

// invalidUTF8 is the one class whose sanitized form GROWS: strings.Map turns each
// invalid byte into the three-byte U+FFFD, so a dense fixture expands roughly 3x
// and the output buffer has to be regrown. It is the shape a byte-truncated
// upstream response actually has.
func invalidUTF8(n int) string {
	return repeatTo("head\xff\xfe\x80tail ", n)
}

// TestSanitizeCleanTextIsAllocationFree gates the package's central cost claim:
// a value carrying no unsafe rune is returned as-is, not copied.
//
// AllocsPerRun averages over runs and these functions allocate either always or
// never, so any non-zero average means a clean value started being copied. Both
// policies are checked at both ends of the size range, because a copy introduced
// by a length check or a pre-pass would be cheap enough at 64 bytes to look like
// noise in a chart and is only unmistakable as an exact zero.
//
// The table holds values that are clean under BOTH policies, so it can drive both
// loops. CR/LF gets its own case below because it is where the two policies' costs
// diverge: newline-bearing text is clean under the multi-line policy (a JSON
// encoder escapes it) and must stay allocation-free, while the same value is a
// rewrite under the strict one.
func TestSanitizeCleanTextIsAllocationFree(t *testing.T) {
	const runs = 20
	tests := []struct {
		name string
		desc string
		in   string
	}{
		{"ascii_small", "clean ASCII, 64 bytes", cleanASCII(64)},
		{"ascii_large", "clean ASCII, 256 KiB", cleanASCII(262144)},
		{"utf8_multibyte", "clean multi-byte UTF-8, 4 KiB", cleanUTF8(typicalBytes)},
		{"literal_replacement_char", "a literal U+FFFD, already a safe rune", "title \uFFFD tail"},
		{"empty", "the empty string", ""},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if got := testing.AllocsPerRun(runs, func() {
				sinkString = runesafe.Sanitize(tc.in)
			}); got != 0 {
				t.Errorf("Sanitize(%s) allocated %v times per run, want 0: strings.Map "+
					"returns the original string when the mapping changes no rune, so "+
					"sanitizing an honest value must cost nothing", tc.desc, got)
			}
		})
	}
	// The strict policy is gated rather than assumed: it is a second strings.Map
	// call, and a copy could be added to one preset and not the other.
	for _, tc := range tests {
		t.Run("single_line_"+tc.name, func(t *testing.T) {
			if got := testing.AllocsPerRun(runs, func() {
				sinkString = runesafe.SanitizeSingleLine(tc.in)
			}); got != 0 {
				t.Errorf("SanitizeSingleLine(%s) allocated %v times per run, want 0: the "+
					"strict preset is one strings.Map call too, so it shares the "+
					"no-change fast path", tc.desc, got)
			}
		})
	}
	t.Run("crlf_kept_by_multiline_policy", func(t *testing.T) {
		const in = "line one\r\nline two\r\n"
		if got := testing.AllocsPerRun(runs, func() {
			sinkString = runesafe.Sanitize(in)
		}); got != 0 {
			t.Errorf("Sanitize(%q) allocated %v times per run, want 0: CR and LF are safe "+
				"under the multi-line policy, so a newline-bearing value is a clean value "+
				"here and must not be copied", in, got)
		}
	})
}

// TestSanitizeIsIdempotentForFree pins the cost half of a claim the README makes
// in prose: "Sanitizing is idempotent, so double-sanitizing at two layers is
// harmless." Idempotence alone only makes the second pass harmless to the VALUE.
// What makes it harmless to the caller is that a sanitized string is by
// construction a clean one, so the second pass takes the no-change fast path and
// allocates nothing — which is what lets a library sanitize at construction and
// an app sanitize again at its emit boundary without paying twice. Without this,
// the advice to double-sanitize would be advice to double-allocate.
func TestSanitizeIsIdempotentForFree(t *testing.T) {
	const runs = 20
	for _, tc := range []struct {
		name string
		desc string
		in   string
	}{
		{"rewritten_unsafe", "rewritten dense unsafe text, 4 KiB", unsafeDense(typicalBytes)},
		{"repaired_invalid_utf8", "repaired invalid UTF-8, 4 KiB", invalidUTF8(typicalBytes)},
	} {
		t.Run(tc.name, func(t *testing.T) {
			once := runesafe.Sanitize(tc.in)
			if got := testing.AllocsPerRun(runs, func() {
				sinkString = runesafe.Sanitize(once)
			}); got != 0 {
				t.Errorf("Sanitize(Sanitize(%s)) allocated %v times per run on the second "+
					"pass, want 0: the first pass leaves no rune the policy touches, so "+
					"double-sanitizing at two layers must cost nothing", tc.desc, got)
			}
		})
	}
}

// TestSanitizeRewriteCostDoesNotScaleWithInput covers the path that legitimately
// allocates. Zero is not the property to want here — a rewritten value has to be
// built somewhere — so what matters is that the cost is PROPORTIONATE to the
// input and cannot be amplified: an upstream that sends 1024x more hostile text
// must not make sanitizing cost 1024x more ALLOCATIONS, only 1024x more bytes
// inside one of them.
//
// Measured on go1.27.0: exactly one allocation for a rewritten value, from 256
// bytes to 256 KiB, however many unsafe runes it carries. strings.Map grows its
// buffer once to len(s)+utf8.UTFMax, and every rune this policy replaces becomes
// a single-byte space, so the rewritten form can only be shorter — one buffer is
// always enough. A per-rune concatenation or a regrowth per replacement would
// break the constancy immediately, which is what this asserts.
func TestSanitizeRewriteCostDoesNotScaleWithInput(t *testing.T) {
	const runs = 10
	sizes := []int{256, 4096, 65536, 262144}
	counts := make([]float64, 0, len(sizes))
	for _, n := range sizes {
		in := unsafeDense(n)
		counts = append(counts, testing.AllocsPerRun(runs, func() {
			sinkString = runesafe.Sanitize(in)
		}))
	}
	for i, got := range counts {
		if got != 1 {
			t.Errorf("Sanitize(dense unsafe text, %d bytes) allocated %v times per run, "+
				"want 1: the rewritten form is never longer than the input, so one "+
				"buffer must suffice", sizes[i], got)
		}
		if got != counts[0] {
			t.Errorf("Sanitize(dense unsafe text, %d bytes) allocated %v times per run "+
				"but %v at %d bytes, want the same count: sanitizing cost must not be "+
				"amplifiable by sending more hostile text", sizes[i], got, counts[0], sizes[0])
		}
	}
}

// TestSanitizeInvalidUTF8CostIsBounded covers the one input class whose sanitized
// form GROWS. Each invalid byte becomes the three-byte U+FFFD, so the buffer
// strings.Map sized at len(s)+utf8.UTFMax is genuinely too small and has to be
// regrown — this path cannot be a single allocation, and asserting one would be
// asserting a bug.
//
// The assertion is a small upper BOUND rather than constancy, because constancy is
// not what this measures. Over 512 bytes to 1 MiB on go1.27.0 the count runs 2, 2,
// 2, 2, 2, 4, 3 — it neither stays fixed nor rises monotonically, because it is
// strings.Builder's amortized growth schedule crossing allocator size classes, not
// anything this package decides. Pinning those steps would pin an implementation
// detail of a dependency, which is the wrong thing to gate.
//
// What is worth gating is the difference between amortized growth and per-byte
// allocation. Amortized growth is logarithmic in the output, so a handful of
// allocations covers any input; a rewrite that allocated per invalid byte, or
// concatenated per rune, would report five figures at the largest fixture. The
// bound is set with headroom above the observed maximum so a future toolchain's
// growth schedule cannot redden it, and it still separates the two behaviours by
// four orders of magnitude.
func TestSanitizeInvalidUTF8CostIsBounded(t *testing.T) {
	const (
		runs     = 5
		maxAlloc = 8
	)
	sizes := []int{1024, 32768, 1048576}
	counts := make([]float64, 0, len(sizes))
	for _, n := range sizes {
		in := invalidUTF8(n)
		got := testing.AllocsPerRun(runs, func() {
			sinkString = runesafe.Sanitize(in)
		})
		counts = append(counts, got)
		if got > maxAlloc {
			t.Errorf("Sanitize(invalid UTF-8, %d bytes) allocated %v times per run, want at "+
				"most %d: repairing invalid bytes may regrow the output buffer a few times, "+
				"but the cost must be the buffer's amortized growth and not one allocation "+
				"per invalid byte", n, got, maxAlloc)
		}
	}
	t.Logf("repairing invalid UTF-8 costs %v allocations at %v bytes: amortized growth, "+
		"not per-byte", counts, sizes)
}

// TestCapBytesNeverAllocates pins CapBytes as a pure slicing operation. It
// returns a prefix of its argument, so it must never copy — which is what makes
// it free to compose after a sanitizer, and free to call speculatively on a value
// that turns out to be within the cap.
//
// The boundary cases are the point. A cut that lands mid-rune walks backwards to
// the nearest rune start, and a plausible rewrite of that walk — re-encoding the
// prefix, or assembling the result through a Builder — would allocate while
// returning the same bytes. Both cap directions and the multi-byte backoff are
// therefore checked as exact zeros, not against a threshold.
func TestCapBytesNeverAllocates(t *testing.T) {
	const runs = 50
	ascii := cleanASCII(smallBytes)
	multi := cleanUTF8(smallBytes)
	for _, tc := range []struct {
		name string
		desc string
		in   string
		n    int
	}{
		{"ascii_under_cap", "ASCII, cap below the length", ascii, smallBytes / 2},
		{"ascii_at_cap", "ASCII, cap exactly the length", ascii, smallBytes},
		{"ascii_over_cap", "ASCII, cap above the length", ascii, smallBytes + 7},
		{"multibyte_backoff", "multi-byte, cut needing a rune-start backoff", multi, smallBytes / 2},
		{"multibyte_over_cap", "multi-byte, cap above the length", multi, smallBytes + 7},
		{"non_positive_cap", "ASCII, non-positive cap", ascii, 0},
	} {
		t.Run(tc.name, func(t *testing.T) {
			if got := testing.AllocsPerRun(runs, func() {
				sinkString = runesafe.CapBytes(tc.in, tc.n)
			}); got != 0 {
				t.Errorf("CapBytes(%s, %d) allocated %v times per run, want 0: the result "+
					"is a prefix of the input, so capping must never copy", tc.desc, tc.n, got)
			}
		})
	}
}

// TestSanitizeSingleLineBoundedWithinCapIsAllocationFree gates the cost half of
// this preset's documented contract: "a result within the cap is returned
// untouched — no marker, byte-identical to the unbounded form". Byte-identical is
// a value claim; what makes the preset safe to put on every log attribute is that
// for a value both clean and within the cap it is cost-identical to doing nothing
// as well. That is the common case at a log site, and it is what an added marker
// check or a defensive copy would silently take away.
func TestSanitizeSingleLineBoundedWithinCapIsAllocationFree(t *testing.T) {
	const runs = 50
	in := cleanASCII(smallBytes)
	for _, tc := range []struct {
		name string
		desc string
		in   string
		n    int
	}{
		{"under_cap", "clean value under the cap", in, smallBytes * 2},
		{"at_cap", "clean value exactly at the cap", in, smallBytes},
		{"empty_negative_cap", "empty input under a negative cap", "", -1},
	} {
		t.Run(tc.name, func(t *testing.T) {
			if got := testing.AllocsPerRun(runs, func() {
				sinkString = runesafe.SanitizeSingleLineBounded(tc.in, tc.n)
			}); got != 0 {
				t.Errorf("SanitizeSingleLineBounded(%s, %d) allocated %v times per run, "+
					"want 0: a clean within-cap value is returned as-is, so the bounded "+
					"preset must cost nothing on the common log-attribute case",
					tc.desc, tc.n, got)
			}
		})
	}
}

// BenchmarkSanitize is the package's primary trend series: the cost of one
// sanitizer call over the input classes a real emit boundary meets.
//
// clean_ascii carries the full size sweep because it is the path almost every real
// call takes, and it is the one where two different regressions are both visible:
// a fixed per-call overhead shows up proportionally largest at 256 bytes, while a
// super-linear scan shows up as a jump between 4 KiB and 64 KiB. unsafe_dense
// carries the same two upper sizes so the rewrite path's slope can be read against
// the fast path's directly. The remaining three classes are measured at the one
// representative size: their value is the RATIO to clean_ascii there — what the
// fast path's bail-out, multi-byte decoding, and U+FFFD growth each cost — which
// one point already gives.
//
// SetBytes is set here because these calls genuinely scan every byte of their
// input, so bytes per second is a real unit. The fleet's reducer drops MB/s from
// the published series deliberately — throughput is the one go-test metric where
// smaller is worse, so charting it on a smaller-is-better series would invert
// every verdict — so it serves the local reader only.
func BenchmarkSanitize(b *testing.B) {
	for _, class := range []struct {
		name  string
		build func(int) string
		sizes []int
	}{
		{"clean_ascii", cleanASCII, []int{smallBytes, typicalBytes, largeBytes}},
		{"unsafe_dense", unsafeDense, []int{typicalBytes, largeBytes}},
		{"unsafe_tail", unsafeTail, []int{typicalBytes}},
		{"clean_utf8", cleanUTF8, []int{typicalBytes}},
		{"invalid_utf8", invalidUTF8, []int{typicalBytes}},
	} {
		b.Run(class.name, func(b *testing.B) {
			for _, n := range class.sizes {
				in := class.build(n)
				b.Run(fmt.Sprintf("bytes_%d", n), func(b *testing.B) {
					b.ReportAllocs()
					b.SetBytes(int64(len(in)))
					for b.Loop() {
						sinkString = runesafe.Sanitize(in)
					}
				})
			}
		})
	}
}

// BenchmarkCapBytes measures the rune-boundary cap around its cut. One iteration
// performs four caps, not one, and both of those choices are forced by the
// measurement rather than preferred: a single cap costs one to two nanoseconds,
// which is at the timer's resolution — measured as 1.284 ns then 1.997 ns for the
// same call on consecutive runs, with a separate ASCII case landing either side of
// it — so a one-call iteration publishes a series whose movement is noise. Four
// calls per iteration lifts it clear of the floor and lets one series carry the
// boundary: a cut needing a multi-byte backoff, a cut on a rune boundary already,
// the cap exactly at the length, and a cap past it.
//
// Read it as a step function. CapBytes compares two lengths and steps back at most
// three bytes, so its cost is independent of the input's size by construction; the
// regression this can catch is a change that makes the cap proportional to the
// input — re-decoding the prefix to find the boundary, or assembling the result —
// which would move the series by orders of magnitude rather than by percent.
//
// No SetBytes for the same reason: the call does not scan its input, so a
// bytes-per-second figure derived from the fixture's length is meaningless. The
// first draft of this file printed 3.5 TB/s.
func BenchmarkCapBytes(b *testing.B) {
	ascii := cleanASCII(typicalBytes)
	multi := cleanUTF8(typicalBytes)
	b.ReportAllocs()
	for b.Loop() {
		sinkString = runesafe.CapBytes(multi, typicalBytes/2)
		sinkString = runesafe.CapBytes(ascii, typicalBytes/2)
		sinkString = runesafe.CapBytes(ascii, typicalBytes)
		sinkString = runesafe.CapBytes(ascii, typicalBytes+7)
	}
}

// BenchmarkSanitizeSingleLineBounded measures the log-attribute preset on the side
// of its cap that does the work: sanitize, cap on a rune boundary, concatenate the
// marker. It is the composition every consumer of this library hand-rolled before
// the preset existed, so it is the series that says what adopting it costs.
//
// Only the truncating case is a series. The within-cap case is the same
// strings.Map call plus a length comparison, so BenchmarkSanitize/clean_ascii at
// the same size already measures it, and the property that matters there — that a
// clean within-cap value is returned without allocating rather than copied — is an
// exact zero a benchmark cannot assert and
// TestSanitizeSingleLineBoundedWithinCapIsAllocationFree does.
//
// No SetBytes: the result is a fixed 200-odd bytes whatever the input size, so
// bytes per second would describe the fixture rather than the work.
func BenchmarkSanitizeSingleLineBounded(b *testing.B) {
	in := cleanASCII(typicalBytes)
	b.ReportAllocs()
	for b.Loop() {
		sinkString = runesafe.SanitizeSingleLineBounded(in, 200)
	}
}
