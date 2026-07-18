// Package runesafe classifies runes that are unsafe in untrusted text bound
// for logs, JSON, or rendered output, and provides shared sanitizers that
// neutralize them.
//
// Untrusted upstream text — API response fields, upstream error messages,
// file names, titles — eventually reaches a sink that renders it: a slog
// line read in a terminal, a JSON report opened in a viewer, a Markdown
// table. Four classes of rune survive the trip into those sinks with their
// control semantics intact and let the upstream author forge or garble what
// the operator sees:
//
//   - C0 controls (U+0000-U+001F) and DEL (U+007F): terminal escape
//     sequences (ESC introduces CSI/OSC sequences that can retitle the
//     terminal, clear the screen, or write to the clipboard) and log-record
//     forgery (a raw newline splits one record into two, letting upstream
//     text fabricate a whole log line). JSON encoders escape C0, so CR and
//     LF may be kept for sinks whose encoder provably escapes them (the
//     keepCRLF policy switch).
//   - C1 controls (U+0080-U+009F): single-rune escape introducers (CSI
//     U+009B, OSC U+009D, ...) with the same terminal powers as ESC
//     sequences. encoding/json and slog's JSONHandler emit them raw, so
//     escaping C0 alone does not close the terminal-injection hole.
//   - Unicode Bidi_Control format characters (U+061C, U+200E/U+200F,
//     U+202A-U+202E, U+2066-U+2069): visually reorder rendered text
//     (Trojan-Source-style spoofing — a link or verdict reads differently
//     than it compares).
//   - The line and paragraph separators U+2028/U+2029: legal unescaped in
//     JSON but line terminators to JavaScript and many viewers, so they
//     split records like a raw newline.
//
// IsUnsafe classifies one rune under an explicit CR/LF policy,
// IsUnsafeNonASCII exposes the above-ASCII subset (C1, bidi controls, the
// separators) for escapers whose sink already covers ASCII, and
// IsBidiControl exposes the Bidi_Control subset. Sanitize (CR/LF kept, for
// JSON-encoded sinks) and SanitizeSingleLine (CR/LF replaced too, for
// single-line sinks) apply the two policies to whole strings, replacing
// each unsafe rune with a space. CapBytes truncates on a rune boundary, so
// a byte cap applied after sanitizing cannot re-introduce a partial-rune C1
// tail. Apply the sanitizer at the emit boundary (the slog call site, just
// before JSON encoding) so comparisons and dedupe keys keep operating on
// the raw value, and use one policy per application so two sinks cannot
// drift.
//
// The package is one deliberately small policy. It is not an HTML/XSS
// sanitizer, does not normalize Unicode (NFC/NFKC), and does not touch
// zero-width or confusable runes; context-aware sinks (Markdown cells,
// URLs, HTML) need their own escaping on top of this rune policy.
package runesafe
