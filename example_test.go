package runesafe_test

import (
	"encoding/json"
	"fmt"
	"strings"

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

// ExampleSanitizeSingleLineBounded bounds an upstream error message for a
// capped log attribute: sanitization happens before the cap is measured, a
// within-cap result comes back untouched, and an over-cap result is cut on
// a rune boundary with the "..." marker appended outside the cap (at most
// n+3 bytes).
func ExampleSanitizeSingleLineBounded() {
	fmt.Printf("%q\n", runesafe.SanitizeSingleLineBounded("bad\nrequest", 20))
	fmt.Printf("%q\n", runesafe.SanitizeSingleLineBounded("葬送のフリーレン", 7))
	// Output:
	// "bad request"
	// "葬送..."
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

// ExampleUntrusted tags an upstream field at the decode boundary: the raw
// bytes survive ingestion for matching (Raw), while fmt, errors, and JSON
// emission all render the sanitized form automatically.
func ExampleUntrusted() {
	var ep struct {
		Title runesafe.Untrusted `json:"title"`
	}
	_ = json.Unmarshal([]byte(`{"title":"Frieren\u202egpj.exe"}`), &ep)

	fmt.Println(len(ep.Title.Raw()))             // raw bytes preserved for compute
	fmt.Println(ep.Title)                        // fmt renders sanitized
	fmt.Println(fmt.Errorf("bad: %s", ep.Title)) // errors carry sanitized text
	out, _ := json.Marshal(ep)
	fmt.Println(string(out)) // encoders emit sanitized
	// Output:
	// 17
	// Frieren gpj.exe
	// bad: Frieren gpj.exe
	// {"title":"Frieren gpj.exe"}
}

// ExampleSanitizeCapped bounds a multi-line upstream body for a JSON sink:
// CR/LF survive the sanitizer, the caller's own marker is charged against
// the cap so the result never exceeds it, and cut carries the truncation
// fact for the caller to record.
func ExampleSanitizeCapped() {
	body, cut := runesafe.SanitizeCapped("line one\nline two\x1b[2J", 16, "...(truncated)")
	fmt.Printf("%q %v %d\n", body, cut, len(body))

	body, cut = runesafe.SanitizeCapped("short\nbody", 16, "...(truncated)")
	fmt.Printf("%q %v\n", body, cut)
	// Output:
	// "li...(truncated)" true 16
	// "short\nbody" false
}

// ExampleSanitizeSingleLineCapped bounds an upstream error for a capped log
// attribute and logs the truncation as a fact rather than inferring it from
// the marker.
func ExampleSanitizeSingleLineCapped() {
	reason, cut := runesafe.SanitizeSingleLineCapped("bad request\nlevel=ERROR forged", 14, "...")
	fmt.Printf("%q %v %d\n", reason, cut, len(reason))
	// Output: "bad request..." true 14
}

// ExampleSanitizeBudgeted bounds an upstream value whose size the caller does
// not control. The value is three megabytes of right-to-left overrides behind
// one letter, and the budget is eight bytes: the pre-cap form cuts the RAW
// input at the budget and sanitizes only that, so the two overrides that fit
// below byte eight are all it ever looked at. SanitizeCapped sanitizes the
// whole value first — every override collapsing to a single-byte space — and
// then fills the same budget from the megabyte of spaces that produces, which
// is the work the pre-cap exists to avoid.
func ExampleSanitizeBudgeted() {
	huge := "A" + strings.Repeat("\u202e", 1<<20)

	text, cut := runesafe.SanitizeBudgeted(huge, 8, "...")
	fmt.Printf("%q %v\n", text, cut)

	text, cut = runesafe.SanitizeCapped(huge, 8, "...")
	fmt.Printf("%q %v\n", text, cut)
	// Output:
	// "A  ..." true
	// "A    ..." true
}

// ExampleBudget spends one shared byte budget across several untrusted values
// and keeps ONE truncation fact for the whole aggregate: each value is capped
// before it is sanitized, the separators are charged too (so a hostile value
// COUNT cannot grow the attribute either), and the marker is charged inside the
// budget, so the joined attribute lands exactly on it with a single marker
// rather than one per value.
func ExampleBudget() {
	b := runesafe.NewBudget(24, "...")
	for i, group := range []string{"SubsPlease", "Erai-raws\u202e", "LostYears", "PMR"} {
		if i > 0 && !b.Write(", ") {
			break
		}
		if !b.Write(group) {
			break
		}
	}
	attr, cut := b.Result()
	fmt.Printf("%q %v %d\n", attr, cut, len(attr))
	// Output: "SubsPlease, Erai-raws..." true 24
}
