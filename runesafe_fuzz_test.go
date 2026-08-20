package runesafe_test

import (
	"encoding/json"
	"strings"
	"testing"
	"unicode"
	"unicode/utf8"

	"github.com/cplieger/runesafe/v2"
)

// fuzzSeeds is the adversarial corpus for the sanitizers: terminal escape
// sequences (C0 ESC and the single-rune C1 introducers), bidi overrides and
// isolates, log-forgery newlines, the JSON-legal line separators, invalid
// UTF-8, and plain multi-byte text that must pass through unchanged.
var fuzzSeeds = []string{
	"",
	"plain text",
	"葬送のフリーレン",
	"a\x1b[2Jb",
	"a\x1b]0;owned\x07b",
	"a\u009b2Jb",
	"a\u009d0;owned\u009cb",
	"line1\nline2\rline3",
	"a\u202evil\u202cb",
	"\u2066isolate\u2069",
	"a\u061c\u200e\u200fb",
	"a\u2028b\u2029c",
	"a\x00\x7f\u0080\u009fb",
	"a\xffb",
	"\xe8\x91",
	"\xed\xa0\x80",
	" \t\n\r",
	// One seed per Unicode class that MOVED between Go 1.26.7's Unicode 15 and
	// Go 1.27.0's Unicode 17. None of these runes is in any of the four unsafe
	// classes on either side, which is the point: the committed corpus is what
	// the weekly fuzz re-explores from and it held none of them, so the claim
	// that the policy is Unicode-version-independent went unprobed at exactly
	// the code points that changed. Each pairs its delta rune with an unsafe
	// one so the sanitizer has work to do at a multi-byte boundary.
	"\u202e\u0390\u1fd3\u03b0\u1fe3\ufb05\ufb06\u009b", // the three pairs SimpleFold now folds together
	"\u0130\u212aKi\u009b",                             // the only two non-ASCII runes that lower to ASCII
	"a\U0001171e\u202eb",                               // U+1171E: Mn -> Mc, the matrix's one category removal
	"\u0295\u1c89\u1c8a\u009b",                         // U+0295: Ll -> Lo, beside the cased pair added in its place
	"\U000105c0\U00010d40\U0001e5d0\u2028",             // assigned in 16/17, unassigned (Cn) in 15
	"\u019b\u0264\ua7d3\ua7d5\u202e",                   // gained an uppercase mapping in 17
}

// FuzzSanitizeSafeIdempotent drives both sanitizer presets with arbitrary
// strings and asserts the full contract of each: the output is valid UTF-8,
// carries no rune its own policy classifies unsafe (cross-function
// consistency with IsUnsafe), preserves the input's rune count (replacement
// is 1:1, never a drop or splice), equals an independent rune-by-rune walk
// (differential oracle for the strings.Map plumbing), and is a fixed point
// (sanitizing is idempotent, so double-sanitizing at two layers is safe).
// It also pins the composition law relating the presets: SanitizeSingleLine
// of a Sanitize output equals SanitizeSingleLine of the raw input, because
// the strict policy is a superset of the keep-CR/LF policy.
func FuzzSanitizeSafeIdempotent(f *testing.F) {
	for _, s := range fuzzSeeds {
		f.Add(s)
	}
	f.Fuzz(func(t *testing.T, in string) {
		presets := []struct {
			name   string
			fn     func(string) string
			unsafe func(rune) bool
		}{
			{"Sanitize", runesafe.Sanitize, runesafe.IsUnsafeMultiLine},
			{"SanitizeSingleLine", runesafe.SanitizeSingleLine, runesafe.IsUnsafeSingleLine},
		}
		for _, p := range presets {
			out := p.fn(in)
			if !utf8.ValidString(out) {
				t.Errorf("%s(%q) = %q, not valid UTF-8", p.name, in, out)
			}
			for _, r := range out {
				if p.unsafe(r) {
					t.Errorf("%s(%q) = %q still carries unsafe rune %U", p.name, in, out, r)
				}
			}
			if got, want := utf8.RuneCountInString(out), utf8.RuneCountInString(in); got != want {
				t.Errorf("%s(%q) changed rune count: %d, want %d", p.name, in, got, want)
			}
			var b strings.Builder
			for _, r := range in {
				if p.unsafe(r) {
					b.WriteRune(' ')
				} else {
					b.WriteRune(r)
				}
			}
			if want := b.String(); out != want {
				t.Errorf("%s(%q) = %q, rune-walk oracle says %q", p.name, in, out, want)
			}
			if again := p.fn(out); again != out {
				t.Errorf("%s not idempotent: %q -> %q -> %q", p.name, in, out, again)
			}
		}
		if got, want := runesafe.SanitizeSingleLine(runesafe.Sanitize(in)), runesafe.SanitizeSingleLine(in); got != want {
			t.Errorf("SanitizeSingleLine(Sanitize(%q)) = %q, want %q (composition law)", in, got, want)
		}
	})
}

// FuzzIsUnsafePolicyConsistency drives the rune classifiers with arbitrary
// int32 values (including negatives and beyond MaxRune, which string
// iteration can never produce but a direct caller can) and asserts the
// policy lattice: a bidi control is unsafe under both policies, keepCRLF
// only ever shrinks the unsafe set, the two policies diverge on CR and LF
// alone, IsUnsafeNonASCII is exactly the above-ASCII restriction of the
// policy (so it sits inside both presets and never flags ASCII), and within
// the valid rune range IsBidiControl agrees with the standard library's
// unicode.Bidi_Control table.
func FuzzIsUnsafePolicyConsistency(f *testing.F) {
	for _, r := range []rune{
		0, '\n', '\r', 0x1b, 0x1f, ' ', '~', 0x7f, 0x80, 0x9b, 0x9f, 0xa0,
		'\u061c', '\u200e', '\u2027', '\u2028', '\u2029', '\u202a', '\u202e',
		'\u2066', '\u2069', 'a', '葬', unicode.MaxRune, -1,
		// The Unicode 17 delta, rune by rune: the three now-folding pairs, the
		// two runes that lower to ASCII, the Mn -> Mc and Ll -> Lo category
		// moves, and two runes that were unassigned in Unicode 15. The
		// classifier must answer identically on both sides of that bump.
		'\u0390', '\u1fd3', '\u03b0', '\u1fe3', '\ufb05', '\ufb06',
		'\u0130', '\u212a', '\U0001171e', '\u0295', '\u1c89',
		'\U000105c0', '\U0001e5d0',
	} {
		f.Add(r)
	}
	f.Fuzz(func(t *testing.T, r rune) {
		keep, strict := runesafe.IsUnsafeMultiLine(r), runesafe.IsUnsafeSingleLine(r)
		if runesafe.IsBidiControl(r) && (!keep || !strict) {
			t.Errorf("IsBidiControl(%U) is true but IsUnsafe = (keepCRLF %v, strict %v), want unsafe under both", r, keep, strict)
		}
		if got, want := runesafe.IsUnsafeNonASCII(r), keep && r > unicode.MaxASCII; got != want {
			t.Errorf("IsUnsafeNonASCII(%U) = %v, want IsUnsafe(r, true) && r > MaxASCII = %v", r, got, want)
		}
		if runesafe.IsUnsafeNonASCII(r) && !strict {
			t.Errorf("IsUnsafeNonASCII(%U) is true but the strict policy says safe; the subset must sit inside both presets", r)
		}
		if keep && !strict {
			t.Errorf("IsUnsafe(%U) unsafe with keepCRLF=true but safe with false; keepCRLF must only shrink the unsafe set", r)
		}
		if strict && !keep && r != '\n' && r != '\r' {
			t.Errorf("IsUnsafe(%U) diverges between policies; only CR and LF may diverge", r)
		}
		if r >= 0 && r <= unicode.MaxRune {
			if got, want := runesafe.IsBidiControl(r), unicode.Is(unicode.Bidi_Control, r); got != want {
				t.Errorf("IsBidiControl(%U) = %v, unicode.Bidi_Control says %v", r, got, want)
			}
		}
	})
}

// FuzzCapBytes drives the rune-boundary cap with arbitrary strings and cap
// values and asserts its contract: the result is a prefix of the input, no
// longer than max(n, 0) bytes, idempotent under the same cap, and for valid
// UTF-8 input it stays valid UTF-8 (never ends in a partial rune) while the
// backoff discards fewer than utf8.UTFMax bytes below the cap.
func FuzzCapBytes(f *testing.F) {
	for _, s := range fuzzSeeds {
		f.Add(s, 3)
		f.Add(s, 0)
		f.Add(s, len(s))
	}
	f.Add("葬送のフリーレン", 7)
	f.Add("a\U0001f600b", 4)
	f.Add("\x80\x81\x82", 2)
	f.Add("abc", -5)
	f.Fuzz(func(t *testing.T, in string, n int) {
		out := runesafe.CapBytes(in, n)
		if !strings.HasPrefix(in, out) {
			t.Errorf("CapBytes(%q, %d) = %q, not a prefix of the input", in, n, out)
		}
		if n <= 0 && out != "" {
			t.Errorf("CapBytes(%q, %d) = %q, want empty for non-positive cap", in, n, out)
		}
		if n > 0 && len(out) > n {
			t.Errorf("CapBytes(%q, %d) = %q, longer than the cap (%d bytes)", in, n, out, len(out))
		}
		if again := runesafe.CapBytes(out, n); again != out {
			t.Errorf("CapBytes not idempotent: %q -> %q -> %q under cap %d", in, out, again, n)
		}
		if utf8.ValidString(in) {
			if !utf8.ValidString(out) {
				t.Errorf("CapBytes(%q, %d) = %q, valid input became invalid UTF-8", in, n, out)
			}
			if n > 0 && len(in) > n && n-len(out) >= utf8.UTFMax {
				t.Errorf("CapBytes(%q, %d) = %q discarded %d bytes below the cap, want < %d", in, n, out, n-len(out), utf8.UTFMax)
			}
		}
	})
}

// FuzzUntrustedContract drives the provenance type with arbitrary strings
// and asserts its full contract against the preset oracles: Raw round-trips
// the exact input bytes, String/MarshalText/LogValue all equal Sanitize,
// SingleLine equals SanitizeSingleLine, MarshalText never errors, a JSON
// round-trip of a tagged field yields the Sanitize form (decode-raw,
// encode-sanitized), and re-tagging an emitted form is a fixed point
// (idempotence carries over from the presets).
func FuzzUntrustedContract(f *testing.F) {
	for _, s := range fuzzSeeds {
		f.Add(s)
	}
	f.Fuzz(func(t *testing.T, in string) {
		u := runesafe.Untrusted(in)
		if u.Raw() != in {
			t.Errorf("Raw() = %q, want exact input %q", u.Raw(), in)
		}
		want := runesafe.Sanitize(in)
		if got := u.String(); got != want {
			t.Errorf("String() = %q, want %q", got, want)
		}
		text, err := u.MarshalText()
		if err != nil {
			t.Errorf("MarshalText() error: %v", err)
		}
		if string(text) != want {
			t.Errorf("MarshalText() = %q, want %q", string(text), want)
		}
		if got := u.LogValue().String(); got != want {
			t.Errorf("LogValue() = %q, want %q", got, want)
		}
		if got, wantStrict := u.SingleLine(), runesafe.SanitizeSingleLine(in); got != wantStrict {
			t.Errorf("SingleLine() = %q, want %q", got, wantStrict)
		}
		var wrapped struct {
			V runesafe.Untrusted `json:"v"`
		}
		wrapped.V = u
		blob, err := json.Marshal(wrapped)
		if err != nil {
			t.Fatalf("Marshal: %v", err)
		}
		var back struct {
			V string `json:"v"`
		}
		if err := json.Unmarshal(blob, &back); err != nil {
			t.Fatalf("Unmarshal: %v", err)
		}
		if back.V != want {
			t.Errorf("JSON round-trip = %q, want Sanitize form %q", back.V, want)
		}
		if again := runesafe.Untrusted(u.String()).String(); again != u.String() {
			t.Errorf("re-tagged String not a fixed point: %q -> %q", u.String(), again)
		}
	})
}

// FuzzSanitizeSingleLineBounded pins the preset's invariants for arbitrary
// input and cap: the result is valid UTF-8 carrying no unsafe rune under the
// strict policy, its length never exceeds max(n,0)+3 (the marker rides
// outside the cap), a within-cap result — and empty input under any cap,
// negative included — is byte-identical to the unbounded SanitizeSingleLine
// form, and an over-cap result is that form's rune-safe prefix plus the
// marker. It also holds the preset byte-identical to its pre-rebuild
// implementation (legacySanitizeSingleLineBounded), the contract every
// existing call site depends on.
func FuzzSanitizeSingleLineBounded(f *testing.F) {
	f.Add("hello", 5)
	f.Add("a\nb\x1bc", 3)
	f.Add("\xff\xfe\xfd", 4)
	f.Add("", 0)
	f.Add("", -1)
	f.Add("é\u202e\u2028x", 2)
	f.Fuzz(func(t *testing.T, s string, n int) {
		got := runesafe.SanitizeSingleLineBounded(s, n)
		if want := legacySanitizeSingleLineBounded(s, n); got != want {
			t.Fatalf("SanitizeSingleLineBounded(%q, %d) = %q, pre-rebuild form was %q", s, n, got, want)
		}
		if !utf8.ValidString(got) {
			t.Fatalf("invalid UTF-8 output: %q", got)
		}
		for _, r := range got {
			if runesafe.IsUnsafeSingleLine(r) {
				t.Fatalf("unsafe rune %U survived: %q", r, got)
			}
		}
		if bound := max(n, 0); len(got) > bound+3 {
			t.Fatalf("output %d bytes exceeds cap %d plus marker: %q", len(got), n, got)
		}
		full := runesafe.SanitizeSingleLine(s)
		if full == "" || len(full) <= n {
			if got != full {
				t.Fatalf("within-cap output %q differs from SanitizeSingleLine %q", got, full)
			}
		} else {
			if !strings.HasSuffix(got, "...") {
				t.Fatalf("over-cap output lacks the truncation marker: %q", got)
			}
			if prefix := strings.TrimSuffix(got, "..."); !strings.HasPrefix(full, prefix) {
				t.Fatalf("truncated body %q is not a prefix of the sanitized form %q", prefix, full)
			}
		}
	})
}

// FuzzSanitizeCapped drives the caller-marker primitive on both CR/LF
// policies with arbitrary input, cap, and marker, and pins the three
// guarantees its consumers rely on: the returned text NEVER exceeds
// max(n, 0) bytes however long the marker is (the hard total bound a
// persisted record or a vendor payload budget needs), cut is true exactly
// when the sanitized form did not fit (and the text is then that form's
// rune-safe prefix, plus the marker verbatim whenever the cap can hold it),
// and an uncut result is byte-identical to the matching unbounded preset.
// The marker is caller program text and deliberately not sanitized, so the
// rune-policy and UTF-8 assertions apply to the body the primitive itself
// produced.
func FuzzSanitizeCapped(f *testing.F) {
	for _, s := range fuzzSeeds {
		f.Add(s, 5, "...")
	}
	f.Add("hello world", 8, "...")
	f.Add("hello world", 2, "...")
	f.Add("hello", 3, "...")
	f.Add("hello world", 10, "...(truncated)")
	f.Add("葬送のフリーレン", 10, "...")
	f.Add("a\U0001f600b", 5, "*")
	f.Add("abc", 0, "")
	f.Add("abc", -1, "...")
	f.Add("", 4, "…")
	f.Fuzz(func(t *testing.T, s string, n int, marker string) {
		variants := []struct {
			name   string
			fn     func(string, int, string) (string, bool)
			full   string
			unsafe func(rune) bool
		}{
			{"SanitizeCapped", runesafe.SanitizeCapped, runesafe.Sanitize(s), runesafe.IsUnsafeMultiLine},
			{"SanitizeSingleLineCapped", runesafe.SanitizeSingleLineCapped, runesafe.SanitizeSingleLine(s), runesafe.IsUnsafeSingleLine},
		}
		for _, v := range variants {
			out, cut := v.fn(s, n, marker)
			if bound := max(n, 0); len(out) > bound {
				t.Fatalf("%s(%q, %d, %q) = %q: %d bytes exceeds the hard cap %d", v.name, s, n, marker, out, len(out), bound)
			}
			if want := v.full != "" && len(v.full) > n; cut != want {
				t.Fatalf("%s(%q, %d, %q) cut = %v, want %v (sanitized form is %d bytes)", v.name, s, n, marker, cut, want, len(v.full))
			}
			body := out
			if cut {
				if n >= len(marker) {
					if !strings.HasSuffix(out, marker) {
						t.Fatalf("%s(%q, %d, %q) = %q: cap holds the marker but the output does not end in it", v.name, s, n, marker, out)
					}
					body = out[:len(out)-len(marker)]
				}
				if !strings.HasPrefix(v.full, body) {
					t.Fatalf("%s(%q, %d, %q) = %q: body %q is not a prefix of the sanitized form %q", v.name, s, n, marker, out, body, v.full)
				}
			} else if out != v.full {
				t.Fatalf("%s(%q, %d, %q) = %q: an uncut result must equal the unbounded form %q", v.name, s, n, marker, out, v.full)
			}
			if !utf8.ValidString(body) {
				t.Fatalf("%s(%q, %d, %q) = %q: body %q is not valid UTF-8", v.name, s, n, marker, out, body)
			}
			for _, r := range body {
				if v.unsafe(r) {
					t.Fatalf("%s(%q, %d, %q) = %q: unsafe rune %U survived in the body", v.name, s, n, marker, out, r)
				}
			}
		}
	})
}

// FuzzBudget drives the work-bounding primitive on both CR/LF policies with
// arbitrary input, budget, and marker, and pins what separates it from the
// Capped pair. Three guarantees are shared with that pair and asserted the same
// way: the text never exceeds max(n, 0) bytes however long the marker is, the
// body the primitive produced is valid UTF-8 carrying no rune its own policy
// calls unsafe, and an uncut result is byte-identical to the matching unbounded
// preset. Two are its own. The WORK bound is structural rather than timed: the
// body must be a prefix of the sanitized form of a rune-boundary prefix of s at
// most max(n, 0) bytes long, so nothing in the result can have come from beyond
// the budget — which is only possible if the raw cap ran before the sanitizer.
// And the packaged form must equal the hand-rolled composition of CapBytes,
// Sanitize and CapBytes it exists to replace (composedBudgeted). The aggregate
// form is held to the same total bound across two writes, with the truncation
// fact monotone once latched.
func FuzzBudget(f *testing.F) {
	for _, s := range fuzzSeeds {
		f.Add(s, 5, "...")
	}
	f.Add("A"+strings.Repeat("\u202e", 64), 8, "...")
	f.Add(strings.Repeat("\u202e", 5), 10, "...")
	f.Add("hello world", 8, "...")
	f.Add("hello world", 2, "...")
	f.Add("hello world", 10, "...(truncated)")
	f.Add("葬送のフリーレン", 10, "...")
	f.Add("a\U0001f600b", 5, "*")
	f.Add("\xff\xff\xff", 8, "...")
	f.Add("abc", 0, "")
	f.Add("abc", -1, "...")
	f.Add("", 4, "…")
	f.Fuzz(func(t *testing.T, s string, n int, marker string) {
		variants := []struct {
			name      string
			fn        func(string, int, string) (string, bool)
			newBudget func(int, string) *runesafe.Budget
			sanitize  func(string) string
			unsafe    func(rune) bool
		}{
			{"SanitizeBudgeted", runesafe.SanitizeBudgeted, runesafe.NewBudget, runesafe.Sanitize, runesafe.IsUnsafeMultiLine},
			{"SanitizeSingleLineBudgeted", runesafe.SanitizeSingleLineBudgeted, runesafe.NewSingleLineBudget, runesafe.SanitizeSingleLine, runesafe.IsUnsafeSingleLine},
		}
		bound := max(n, 0)
		for _, v := range variants {
			out, cut := v.fn(s, n, marker)
			if len(out) > bound {
				t.Fatalf("%s(%q, %d, %q) = %q: %d bytes exceeds the hard cap %d", v.name, s, n, marker, out, len(out), bound)
			}
			wantText, wantCut := composedBudgeted(s, n, marker, v.sanitize)
			if out != wantText || cut != wantCut {
				t.Fatalf("%s(%q, %d, %q) = (%q, %v), hand-composed form is (%q, %v)", v.name, s, n, marker, out, cut, wantText, wantCut)
			}
			body := out
			if cut && n >= len(marker) {
				if !strings.HasSuffix(out, marker) {
					t.Fatalf("%s(%q, %d, %q) = %q: cap holds the marker but the output does not end in it", v.name, s, n, marker, out)
				}
				body = out[:len(out)-len(marker)]
			}
			if !cut && out != v.sanitize(s) {
				t.Fatalf("%s(%q, %d, %q) = %q: an uncut result must equal the unbounded form %q", v.name, s, n, marker, out, v.sanitize(s))
			}
			if walked := v.sanitize(runesafe.CapBytes(s, n)); !strings.HasPrefix(walked, body) {
				t.Fatalf("%s(%q, %d, %q) = %q: body %q is not a prefix of %q, the sanitized form of the budgeted raw prefix",
					v.name, s, n, marker, out, body, walked)
			}
			if !utf8.ValidString(body) {
				t.Fatalf("%s(%q, %d, %q) = %q: body %q is not valid UTF-8", v.name, s, n, marker, out, body)
			}
			for _, r := range body {
				if v.unsafe(r) {
					t.Fatalf("%s(%q, %d, %q) = %q: unsafe rune %U survived in the body", v.name, s, n, marker, out, r)
				}
			}

			b := v.newBudget(n, marker)
			b.Write(s)
			if one, oneCut := b.Result(); one != out || oneCut != cut {
				t.Fatalf("one Write into %s's budget = (%q, %v), want the single-value form (%q, %v)", v.name, one, oneCut, out, cut)
			}
			b.Write(s)
			two, twoCut := b.Result()
			if len(two) > bound {
				t.Fatalf("two writes of %q under budget %d = %q: %d bytes exceeds the hard cap %d", s, n, two, len(two), bound)
			}
			if cut && !twoCut {
				t.Fatalf("two writes of %q under budget %d = (%q, false): the truncation fact must stay latched", s, n, two)
			}
		}
	})
}
