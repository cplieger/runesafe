package runesafe_test

import (
	"strings"
	"testing"
	"unicode"
	"unicode/utf8"

	"github.com/cplieger/runesafe"
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
			name     string
			fn       func(string) string
			keepCRLF bool
		}{
			{"Sanitize", runesafe.Sanitize, true},
			{"SanitizeSingleLine", runesafe.SanitizeSingleLine, false},
		}
		for _, p := range presets {
			out := p.fn(in)
			if !utf8.ValidString(out) {
				t.Errorf("%s(%q) = %q, not valid UTF-8", p.name, in, out)
			}
			for _, r := range out {
				if runesafe.IsUnsafe(r, p.keepCRLF) {
					t.Errorf("%s(%q) = %q still carries unsafe rune %U", p.name, in, out, r)
				}
			}
			if got, want := utf8.RuneCountInString(out), utf8.RuneCountInString(in); got != want {
				t.Errorf("%s(%q) changed rune count: %d, want %d", p.name, in, got, want)
			}
			var b strings.Builder
			for _, r := range in {
				if runesafe.IsUnsafe(r, p.keepCRLF) {
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
// alone, and within the valid rune range IsBidiControl agrees with the
// standard library's unicode.Bidi_Control table.
func FuzzIsUnsafePolicyConsistency(f *testing.F) {
	for _, r := range []rune{
		0, '\n', '\r', 0x1b, 0x1f, ' ', '~', 0x7f, 0x80, 0x9b, 0x9f, 0xa0,
		'\u061c', '\u200e', '\u2027', '\u2028', '\u202a', '\u202e', '\u2066',
		'\u2069', 'a', '葬', unicode.MaxRune, -1,
	} {
		f.Add(r)
	}
	f.Fuzz(func(t *testing.T, r rune) {
		keep, strict := runesafe.IsUnsafe(r, true), runesafe.IsUnsafe(r, false)
		if runesafe.IsBidiControl(r) && (!keep || !strict) {
			t.Errorf("IsBidiControl(%U) is true but IsUnsafe = (keepCRLF %v, strict %v), want unsafe under both", r, keep, strict)
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
