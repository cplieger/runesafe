package runesafe

import (
	"strings"
	"unicode"
	"unicode/utf8"
)

// IsBidiControl reports whether r is one of Unicode's Bidi_Control format
// characters: the singleton marks U+061C (ALM) and U+200E/U+200F (LRM/RLM),
// the override/embedding range U+202A-U+202E (LRE/RLE/PDF/LRO/RLO), and the
// isolate range U+2066-U+2069 (LRI/RLI/FSI/PDI). Any of them in untrusted
// text can visually reorder rendered output (Trojan-Source-style report and
// link spoofing), so every output sanitizer treats the full set as unsafe.
// The set matches unicode.Bidi_Control exactly, without the table lookup.
func IsBidiControl(r rune) bool {
	return r == '\u061c' || r == '\u200e' || r == '\u200f' ||
		(r >= '\u202a' && r <= '\u202e') ||
		(r >= '\u2066' && r <= '\u2069')
}

// IsUnsafe reports whether r is unsafe in untrusted text bound for a log,
// JSON, or rendered-output sink: a C0 control, DEL, a C1 control
// (U+0080-U+009F, single-rune terminal-escape introducers), a Unicode bidi
// control (IsBidiControl), or the U+2028/U+2029 line separators. keepCRLF
// selects the CR/LF policy: true treats CR and LF as safe, for sinks whose
// encoder escapes them (JSON); false treats them as unsafe like every other
// C0 control, for single-line sinks where a raw newline forges a new record.
// Sanitize and SanitizeSingleLine apply the two policies to whole strings.
func IsUnsafe(r rune, keepCRLF bool) bool {
	switch {
	case r < 0x20:
		return !keepCRLF || (r != '\n' && r != '\r')
	case r == 0x7f:
		return true
	case r >= 0x80 && r <= 0x9f:
		return true
	case IsBidiControl(r) || r == '\u2028' || r == '\u2029':
		return true
	}
	return false
}

// IsUnsafeNonASCII reports whether r is an unsafe rune above the ASCII
// range: a C1 control (U+0080-U+009F), a Unicode bidi control
// (IsBidiControl), or the U+2028/U+2029 line separators. It is the IsUnsafe
// policy minus C0, DEL, and the CR/LF axis (all ASCII), for composed
// escapers whose sink already covers ASCII: a URL percent-encoder escapes
// C0, DEL, and whitespace itself, but url.Parse accepts these non-ASCII
// runes raw, and a terminal or Markdown viewer must never receive them.
// The CR/LF policy switch is moot above ASCII, so there is no keepCRLF
// parameter: IsUnsafeNonASCII(r) equals IsUnsafe(r, keepCRLF) && r >
// unicode.MaxASCII under either policy.
func IsUnsafeNonASCII(r rune) bool {
	return r > unicode.MaxASCII && IsUnsafe(r, true)
}

// Sanitize makes an untrusted string safe for slog/JSON sinks by replacing
// each unsafe rune with a space: C0 controls (except CR/LF, which JSON
// encoders escape), DEL, C1 controls (U+0080-U+009F, single-rune
// terminal-escape introducers emitted raw by encoding/json and slog's
// JSONHandler), Unicode bidi controls, and the U+2028/U+2029 line
// separators. Invalid UTF-8 bytes become U+FFFD (a strings.Map property),
// so the result is always valid UTF-8. Apply it to every untrusted
// attribute at the emit boundary — one policy shared by all of an app's
// sinks, so they cannot drift. For a single-line sink where CR/LF must also
// go, use SanitizeSingleLine.
func Sanitize(s string) string {
	return sanitize(s, true)
}

// SanitizeSingleLine makes an untrusted string safe for a single-line sink —
// a plain-text log line, a one-line error message, a rendered table cell —
// by replacing each unsafe rune with a space under the strict keepCRLF=false
// policy: everything Sanitize replaces, plus CR and LF, whose raw appearance
// in a single-line sink forges a record boundary. Invalid UTF-8 bytes become
// U+FFFD, so the result is always valid UTF-8 and carries no line break.
func SanitizeSingleLine(s string) string {
	return sanitize(s, false)
}

// sanitize applies the IsUnsafe policy to every rune of s, replacing each
// unsafe rune with a space via strings.Map (which also converts each invalid
// UTF-8 byte to U+FFFD, a safe rune under both policies).
func sanitize(s string, keepCRLF bool) string {
	return strings.Map(func(r rune) rune {
		if IsUnsafe(r, keepCRLF) {
			return ' '
		}
		return r
	}, s)
}

// CapBytes truncates s to at most n bytes without splitting a multi-byte
// rune at the cut: the cut backs off to the nearest rune start, so the
// result never ends in a partial rune. It exists because sanitizing can grow
// a string (each invalid UTF-8 byte becomes the three-byte U+FFFD), so a
// pre-sanitize byte cap does not survive sanitization, and a naive re-cap
// can split a rune — leaving a partial-rune tail whose raw 0x80-0x9F bytes
// a non-UTF-8 terminal reads as C1 escape introducers, re-minting the very
// class of byte the sanitizer removed. For valid UTF-8 input (any Sanitize or
// SanitizeSingleLine output) the backoff discards at most three bytes below
// n and the result is a valid-UTF-8 prefix of s. A non-positive n returns
// the empty string.
func CapBytes(s string, n int) string {
	if n <= 0 {
		return ""
	}
	if len(s) <= n {
		return s
	}
	cut := n
	for cut > 0 && !utf8.RuneStart(s[cut]) {
		cut--
	}
	return s[:cut]
}
