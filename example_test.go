package runesafe_test

import (
	"fmt"

	"github.com/cplieger/runesafe"
)

// ExampleSanitize sanitizes an upstream-controlled title before it becomes
// a slog attribute: the terminal-escape introducer and the bidi override
// become spaces, ordinary text passes through.
func ExampleSanitize() {
	title := "Frieren\x1b[2J \u202egpj.exe"
	fmt.Printf("%q\n", runesafe.Sanitize(title))
	// Output: "Frieren [2J  gpj.exe"
}

// ExampleSanitizeSingleLine flattens an upstream error message for a
// one-line sink: the newline that would forge a second log record becomes a
// space along with the escape introducer.
func ExampleSanitizeSingleLine() {
	msg := "bad request\nlevel=ERROR forged\x1b[2J"
	fmt.Printf("%q\n", runesafe.SanitizeSingleLine(msg))
	// Output: "bad request level=ERROR forged [2J"
}

// ExampleCapBytes bounds a sanitized string without splitting a multi-byte
// rune: the cut backs off to a rune boundary instead of leaving a partial
// rune's raw tail bytes.
func ExampleCapBytes() {
	s := "葬送のフリーレン" // three bytes per rune
	fmt.Printf("%q\n", runesafe.CapBytes(s, 7))
	// Output: "葬送"
}

// ExampleIsUnsafe shows the CR/LF policy switch: a newline is safe for a
// JSON sink whose encoder escapes it, and unsafe for a single-line sink.
func ExampleIsUnsafe() {
	fmt.Println(runesafe.IsUnsafe('\n', true), runesafe.IsUnsafe('\n', false))
	fmt.Println(runesafe.IsUnsafe('\u009b', true), runesafe.IsUnsafe('\u009b', false))
	// Output:
	// false true
	// true true
}

// ExampleIsBidiControl classifies a right-to-left override against a plain
// letter.
func ExampleIsBidiControl() {
	fmt.Println(runesafe.IsBidiControl('\u202e'), runesafe.IsBidiControl('a'))
	// Output: true false
}

// ExampleIsUnsafeNonASCII classifies the runes a URL percent-escaper must
// encode even though url.Parse accepts them raw: the C1 escape introducer
// and the bidi override are flagged, while an ASCII control (the sink's own
// escaping covers it) and a plain letter are not.
func ExampleIsUnsafeNonASCII() {
	fmt.Println(runesafe.IsUnsafeNonASCII('\u009b'), runesafe.IsUnsafeNonASCII('\u202e'))
	fmt.Println(runesafe.IsUnsafeNonASCII('\x1b'), runesafe.IsUnsafeNonASCII('é'))
	// Output:
	// true true
	// false false
}
