package runesafe_test

import (
	"log/slog"
	"testing"

	"github.com/cplieger/runesafe/v2"
)

// This file exists because Untrusted is the adoption shape the README pushes
// hardest — tag the field once at its decode struct and every sink sanitizes
// automatically — and that recommendation only holds if the automatic call is
// cheap. A tagged field resolves through LogValue on EVERY log line that carries
// it, in every handler, so its cost is multiplied by the app's log volume rather
// than by the number of call sites the author can see. If tagging were more
// expensive than a hand-written Sanitize call at the emit site, the type's whole
// argument would invert, and nothing measured it.
//
// It is not more expensive: LogValue is Sanitize plus slog.StringValue, and
// slog.StringValue stores the string without copying it, so for a clean value the
// entire resolution is allocation-free. Measured on go1.27.0 with a slog.Value
// sink: 0 allocations for clean text, 1 for text the policy has to rewrite.
//
// One asymmetry the README does not mention and this file pins: the type is
// allocation-free through slog and fmt but NOT through an encoder. MarshalText has
// to return a []byte, and a string cannot become one without a copy, so JSON
// emission of a clean tagged value costs exactly one allocation that the same
// value costs nothing to log. That is a property of encoding.TextMarshaler's
// signature rather than of this package, it is unavoidable without changing that
// interface, and it is asserted here so the number is a decision rather than a
// surprise.
//
// The sink is a slog.Value, deliberately. Consuming LogValue's result through an
// `any` variable measures 1 allocation for clean text — that is the interface
// boxing of the returned Value, which no real caller pays: slog calls LogValue and
// keeps the result in a slog.Value field. Measuring through the wrong sink here
// would have published a cost that does not exist.
var sinkValue slog.Value

// TestUntrustedSinksAreAllocationFreeForCleanText gates the cost claim under the
// type's adoption advice: tagging a field is free for the values an honest
// upstream sends, at every sink the tag fires through.
//
// All three read paths are checked because they are three separate methods that
// happen to share an implementation today; a future edit that gave one of them its
// own copy, a length check, or a defensive conversion would be invisible to every
// test that only compares output.
func TestUntrustedSinksAreAllocationFreeForCleanText(t *testing.T) {
	const runs = 50
	u := runesafe.Untrusted(cleanASCII(typicalBytes))
	for _, tc := range []struct {
		name string
		desc string
		call func()
	}{
		{"log_value", "LogValue, the slog.LogValuer path", func() { sinkValue = u.LogValue() }},
		{"string", "String, the fmt.Stringer and fmt.Errorf path", func() { sinkString = u.String() }},
		{"single_line", "SingleLine, the strict hand-built-sink path", func() { sinkString = u.SingleLine() }},
		{"raw", "Raw, the compute path that must not transform", func() { sinkString = u.Raw() }},
	} {
		t.Run(tc.name, func(t *testing.T) {
			if got := testing.AllocsPerRun(runs, tc.call); got != 0 {
				t.Errorf("Untrusted(clean ASCII, %d bytes).%s allocated %v times per run, "+
					"want 0: a tagged field resolves on every log line that carries it, so "+
					"tagging an honest value must cost nothing", typicalBytes, tc.desc, got)
			}
		})
	}
}

// TestUntrustedMarshalTextCopiesOnce pins the one sink the allocation-free property
// does not reach, so that the asymmetry is recorded rather than rediscovered.
// encoding.TextMarshaler returns []byte, and converting the sanitized string to
// one is a copy no fast path can remove: a clean tagged value therefore costs
// exactly one allocation to JSON-encode and none to log. Asserting the number
// keeps the copy single — a rewrite that assembled the bytes through a Builder, or
// sanitized into an intermediate, would double it.
func TestUntrustedMarshalTextCopiesOnce(t *testing.T) {
	const runs = 50
	u := runesafe.Untrusted(cleanASCII(typicalBytes))
	var buf []byte
	if got := testing.AllocsPerRun(runs, func() {
		buf, _ = u.MarshalText()
	}); got != 1 {
		t.Errorf("Untrusted(clean ASCII, %d bytes).MarshalText allocated %v times per run, "+
			"want 1: the []byte the TextMarshaler contract requires is one unavoidable "+
			"copy of the sanitized string, and only one", typicalBytes, got)
	}
	sinkString = string(buf)
}

// BenchmarkUntrustedLogValue is the per-log-line series for the tag: what one
// tagged attribute adds to an emitted record, before any handler formatting.
//
// It measures the LogValuer resolution alone rather than a whole slog call, on
// purpose — routing it through a real handler would put slog's own encoder in the
// series, so a stdlib change would read here as a runesafe regression.
//
// The fixture is clean because that is both the common case and the only one where
// the tag's cost is a claim rather than an inheritance: this series should sit on
// top of BenchmarkSanitize/clean_ascii at the same size, at zero allocations, and
// any daylight between them is wrapper overhead that should not exist. The
// rewrite case is deliberately not a second series — it is BenchmarkSanitize's
// unsafe_dense measurement plus the same wrapper, so it would publish a duplicate
// trend under a second name.
func BenchmarkUntrustedLogValue(b *testing.B) {
	u := runesafe.Untrusted(cleanASCII(typicalBytes))
	b.ReportAllocs()
	for b.Loop() {
		sinkValue = u.LogValue()
	}
}
