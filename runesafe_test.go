package runesafe_test

import (
	"testing"
	"unicode"
	"unicode/utf8"

	"github.com/cplieger/runesafe"
)

// TestSanitize pins the shared unsafe-rune policy for the slog/JSON sinks:
// every C0 control except CR/LF (which JSON encoders escape), DEL, C1
// controls, the Unicode Bidi_Control set, and the U+2028/U+2029 line
// separators become spaces, while plain ASCII and non-control Unicode pass
// through unchanged.
func TestSanitize(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want string
	}{
		{"C0 escape introducer", "a\x1b[2Jb", "a [2Jb"},
		{"C0 NUL", "a\x00b", "a b"},
		{"C0 BEL", "a\x07b", "a b"},
		{"tab", "a\tb", "a b"},
		{"LF preserved", "a\nb", "a\nb"},
		{"CR preserved", "a\rb", "a\rb"},
		{"DEL", "a\x7fb", "a b"},
		{"C1 CSI", "a\u009bb", "a b"},
		{"C1 OSC", "a\u009db", "a b"},
		{"C1 range start", "a\u0080b", "a b"},
		{"C1 range end", "a\u009fb", "a b"},
		{"bidi ALM", "a\u061cb", "a b"},
		{"bidi LRM", "a\u200eb", "a b"},
		{"bidi RLO", "a\u202eb", "a b"},
		{"bidi isolate FSI", "a\u2068b", "a b"},
		{"line separator", "a\u2028b", "a b"},
		{"paragraph separator", "a\u2029b", "a b"},
		{"adjacent unsafe runes", "a\x1b\u202e\u009bb", "a   b"},
		{"plain ASCII unchanged", "Frieren: Beyond Journey's End", "Frieren: Beyond Journey's End"},
		{"plain unicode unchanged", "葬送のフリーレン", "葬送のフリーレン"},
		{"boundary neighbors unchanged", " ~\u00a0\u2027\u206a", " ~\u00a0\u2027\u206a"},
		{"empty string", "", ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := runesafe.Sanitize(tt.in); got != tt.want {
				t.Errorf("Sanitize(%q) = %q, want %q", tt.in, got, tt.want)
			}
		})
	}
}

// TestSanitizeSingleLine pins the strict keepCRLF=false preset: CR and LF
// become spaces like every other C0 control, the other rune classes match
// Sanitize, and safe text passes through unchanged.
func TestSanitizeSingleLine(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want string
	}{
		{"LF replaced", "a\nb", "a b"},
		{"CR replaced", "a\rb", "a b"},
		{"CRLF replaced", "line1\r\nline2", "line1  line2"},
		{"C0 escape introducer", "a\x1b[2Jb", "a [2Jb"},
		{"C1 CSI", "a\u009bb", "a b"},
		{"bidi RLO", "a\u202eb", "a b"},
		{"line separator", "a\u2028b", "a b"},
		{"forged log record", "ok\nlevel=ERROR fake", "ok level=ERROR fake"},
		{"plain unicode unchanged", "葬送のフリーレン", "葬送のフリーレン"},
		{"empty string", "", ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := runesafe.SanitizeSingleLine(tt.in); got != tt.want {
				t.Errorf("SanitizeSingleLine(%q) = %q, want %q", tt.in, got, tt.want)
			}
		})
	}
}

// TestSanitizeInvalidUTF8 pins the strings.Map property the doc comment
// promises for both presets: each invalid UTF-8 byte decodes as U+FFFD
// (safe, kept as the replacement rune), so the output is always valid UTF-8
// and never carries raw invalid bytes into a JSON encoder.
func TestSanitizeInvalidUTF8(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want string
	}{
		{"lone continuation byte", "a\xffb", "a\ufffdb"},
		{"truncated multibyte", "a\xe8\x91b", "a\ufffd\ufffdb"},
		{"surrogate half encoding", "a\xed\xa0\x80b", "a\ufffd\ufffd\ufffdb"},
		{"invalid byte beside unsafe rune", "\xff\x1b", "\ufffd "},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for fn, got := range map[string]string{
				"Sanitize":           runesafe.Sanitize(tt.in),
				"SanitizeSingleLine": runesafe.SanitizeSingleLine(tt.in),
			} {
				if got != tt.want {
					t.Errorf("%s(%q) = %q, want %q", fn, tt.in, got, tt.want)
				}
				if !utf8.ValidString(got) {
					t.Errorf("%s(%q) = %q is not valid UTF-8", fn, tt.in, got)
				}
			}
		})
	}
}

// TestCapBytes pins the rune-boundary byte cap: the result is a prefix of
// the input no longer than n bytes that never ends in a partial rune, so a
// cap applied to sanitized output cannot re-introduce raw 0x80-0x9F tail
// bytes.
func TestCapBytes(t *testing.T) {
	tests := []struct {
		name string
		in   string
		n    int
		want string
	}{
		{"under cap", "abc", 10, "abc"},
		{"exact cap", "abc", 3, "abc"},
		{"ascii cut", "abcdef", 4, "abcd"},
		{"two-byte rune backoff", "aé", 2, "a"},
		{"three-byte rune backoff", "葬送", 4, "葬"},
		{"three-byte rune exact boundary", "葬送", 3, "葬"},
		{"four-byte rune backoff", "a\U0001f600", 3, "a"},
		{"replacement-rune growth capped", "\ufffd\ufffd", 4, "\ufffd"},
		{"zero cap", "abc", 0, ""},
		{"negative cap", "abc", -1, ""},
		{"empty input", "", 5, ""},
		{"all continuation bytes", "\x80\x81\x82", 2, ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := runesafe.CapBytes(tt.in, tt.n); got != tt.want {
				t.Errorf("CapBytes(%q, %d) = %q, want %q", tt.in, tt.n, got, tt.want)
			}
		})
	}
}

// TestIsUnsafeCRLFPolicy pins the keepCRLF switch: CR and LF are safe only
// when the sink's encoder escapes them (keepCRLF true); a single-line sink
// (keepCRLF false) treats them as unsafe like every other C0 control, and
// the other rune classes are unsafe under both policies.
func TestIsUnsafeCRLFPolicy(t *testing.T) {
	for _, r := range []rune{'\n', '\r'} {
		if runesafe.IsUnsafe(r, true) {
			t.Errorf("IsUnsafe(%U, keepCRLF=true) = true, want false", r)
		}
		if !runesafe.IsUnsafe(r, false) {
			t.Errorf("IsUnsafe(%U, keepCRLF=false) = false, want true", r)
		}
	}
	for _, r := range []rune{0x00, 0x1b, '\t', 0x7f, '\u009b', '\u202e', '\u2028'} {
		if !runesafe.IsUnsafe(r, true) || !runesafe.IsUnsafe(r, false) {
			t.Errorf("IsUnsafe(%U) = false under a keepCRLF policy, want true under both", r)
		}
	}
	for _, r := range []rune{'a', ' ', 'é', '葬'} {
		if runesafe.IsUnsafe(r, true) || runesafe.IsUnsafe(r, false) {
			t.Errorf("IsUnsafe(%U) = true, want false under both policies", r)
		}
	}
}

// TestIsBidiControl pins the complete Bidi_Control set (the U+061C/U+200E/
// U+200F singleton marks plus the override and isolate ranges) and the
// adjacent near-miss code points that must stay safe.
func TestIsBidiControl(t *testing.T) {
	unsafe := []rune{
		'\u061c', '\u200e', '\u200f',
		'\u202a', '\u202b', '\u202c', '\u202d', '\u202e',
		'\u2066', '\u2067', '\u2068', '\u2069',
	}
	for _, r := range unsafe {
		if !runesafe.IsBidiControl(r) {
			t.Errorf("IsBidiControl(%U) = false, want true", r)
		}
	}
	safe := []rune{'\u061b', '\u061d', '\u200d', '\u2029', '\u2065', '\u206a', 'a'}
	for _, r := range safe {
		if runesafe.IsBidiControl(r) {
			t.Errorf("IsBidiControl(%U) = true, want false", r)
		}
	}
}

// TestClassifierUnicodeConformance sweeps every valid rune and checks the
// classifiers against the standard library's Unicode tables as an
// independent oracle: IsBidiControl must equal unicode.Is(Bidi_Control, r)
// exactly (the doc comment's claim), IsUnsafe must equal the documented
// union — Cc controls (unicode.IsControl covers C0, DEL, and C1 exactly),
// Bidi_Control, and the U+2028/U+2029 separators — with the two policies
// diverging on CR/LF alone, and IsUnsafeNonASCII must equal both the
// above-ASCII restriction of that union and the literal C1|bidi|separator
// enumeration a percent-escaper composes. This is the drift guard against a
// future hand-edit of the hardcoded ranges.
func TestClassifierUnicodeConformance(t *testing.T) {
	for r := rune(0); r <= unicode.MaxRune; r++ {
		bidi := unicode.Is(unicode.Bidi_Control, r)
		if got := runesafe.IsBidiControl(r); got != bidi {
			t.Fatalf("IsBidiControl(%U) = %v, unicode.Bidi_Control says %v", r, got, bidi)
		}
		unsafeStrict := unicode.IsControl(r) || bidi || r == '\u2028' || r == '\u2029'
		if got := runesafe.IsUnsafe(r, false); got != unsafeStrict {
			t.Fatalf("IsUnsafe(%U, false) = %v, oracle says %v", r, got, unsafeStrict)
		}
		unsafeKeep := unsafeStrict && r != '\n' && r != '\r'
		if got := runesafe.IsUnsafe(r, true); got != unsafeKeep {
			t.Fatalf("IsUnsafe(%U, true) = %v, oracle says %v", r, got, unsafeKeep)
		}
		nonASCII := unsafeStrict && r > unicode.MaxASCII
		if got := runesafe.IsUnsafeNonASCII(r); got != nonASCII {
			t.Fatalf("IsUnsafeNonASCII(%U) = %v, oracle says %v", r, got, nonASCII)
		}
		enumerated := (r >= 0x80 && r <= 0x9f) || bidi || r == '\u2028' || r == '\u2029'
		if got := runesafe.IsUnsafeNonASCII(r); got != enumerated {
			t.Fatalf("IsUnsafeNonASCII(%U) = %v, C1|bidi|separator enumeration says %v", r, got, enumerated)
		}
	}
}

// TestIsUnsafeNonASCII pins the above-ASCII subset predicate: C1 controls,
// the Bidi_Control set, and the U+2028/U+2029 separators are flagged, while
// every ASCII rune — including the C0 controls, CR/LF, and DEL that IsUnsafe
// classifies — and safe non-ASCII neighbors stay false.
func TestIsUnsafeNonASCII(t *testing.T) {
	unsafe := []rune{
		'\u0080', '\u009b', '\u009d', '\u009f',
		'\u061c', '\u200e', '\u202e', '\u2066', '\u2069',
		'\u2028', '\u2029',
	}
	for _, r := range unsafe {
		if !runesafe.IsUnsafeNonASCII(r) {
			t.Errorf("IsUnsafeNonASCII(%U) = false, want true", r)
		}
	}
	safe := []rune{0x00, '\x1b', '\t', '\n', '\r', 0x7f, 'a', '~', '\u00a0', '\u2027', '\u206a', 'é', '葬'}
	for _, r := range safe {
		if runesafe.IsUnsafeNonASCII(r) {
			t.Errorf("IsUnsafeNonASCII(%U) = true, want false", r)
		}
	}
}

// TestSanitizeSingleLineBounded covers the log-bound preset: within-cap
// input is byte-identical to SanitizeSingleLine (no marker), an over-cap
// result truncates on a rune boundary with the "..." marker outside the cap,
// the cap measures the SANITIZED form (invalid bytes grow into U+FFFD), and
// the non-positive-cap and empty-input edges hold.
func TestSanitizeSingleLineBounded(t *testing.T) {
	tests := []struct {
		name string
		in   string
		n    int
		want string
	}{
		{"within cap untouched", "hello", 10, "hello"},
		{"exactly at cap untouched", "hello", 5, "hello"},
		{"over cap truncated with marker", "hello world", 5, "hello..."},
		{"sanitizes before capping", "a\nb\x1bc", 10, "a b c"},
		{"rune boundary respected", "aé", 2, "a..."},
		{"cap measures sanitized growth", "\xff\xff", 3, "\uFFFD..."},
		{"non-positive cap yields the marker alone", "abc", 0, "..."},
		{"empty input stays empty", "", 0, ""},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if got := runesafe.SanitizeSingleLineBounded(tc.in, tc.n); got != tc.want {
				t.Errorf("SanitizeSingleLineBounded(%q, %d) = %q, want %q", tc.in, tc.n, got, tc.want)
			}
		})
	}
}
