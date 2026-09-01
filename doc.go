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
//     CR/LF policy split: IsUnsafeMultiLine vs IsUnsafeSingleLine).
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
// The four classes are enumerated as CODE POINTS, never looked up in a Unicode
// table, so a Unicode upgrade cannot move the policy. Verified rune by rune over
// all 1,114,112 code points: every predicate and both sanitizers return identical
// answers on Go 1.26.7 (Unicode 15.0.0) and Go 1.27.0 (Unicode 17.0.0). The tables
// are consulted only by the tests, as an independent oracle: they fail if a future
// Unicode adds a member to a class enumerated here, and each class's size is pinned
// so mirroring such an addition has to be deliberate. What a Unicode upgrade does
// move is how a SINK renders this package's output — slog's TextHandler decides
// whether to quote an attribute via unicode.IsPrint — which is an encoder's choice
// about a graphic rune, not a change in what this package replaces.
//
// The Untrusted string type makes the policy travel with the value: tag an
// upstream field at its decode struct and every standard sink applies
// Sanitize automatically — slog via LogValuer, fmt and error construction
// via Stringer, encoders via TextMarshaler — while Raw preserves the exact
// bytes for matching, dedupe keys, and composed escapers. Machine-read
// persistence must store Raw (encoding fires the sanitizer), and
// construction-time sanitization remains the boundary for text that must
// be safe unconditionally through every future sink.
//
// IsUnsafeMultiLine and IsUnsafeSingleLine classify one rune per policy,
// IsUnsafeNonASCII exposes the above-ASCII subset (C1, bidi controls, the
// separators) for escapers whose sink already covers ASCII, and
// IsBidiControl exposes the Bidi_Control subset. Sanitize (CR/LF kept, for
// JSON-encoded sinks) and SanitizeSingleLine (CR/LF replaced too, for
// single-line sinks) apply the two policies to whole strings, replacing
// each unsafe rune with a space. CapBytes truncates on a rune boundary, so
// a byte cap applied after sanitizing cannot re-introduce a partial-rune C1
// tail, and SanitizeSingleLineBounded packages the log-bound composition —
// SanitizeSingleLine, then CapBytes on the sanitized form, then a "..."
// marker outside the cap — for upstream-controlled values headed into
// capped log attributes. SanitizeCapped and SanitizeSingleLineCapped are
// the general form of that composition, one per CR/LF policy, for the
// callers the preset cannot serve: the marker is caller-supplied and
// charged against the cap, so the returned text is a hard total bound (what
// a persisted record's write limit or a vendor payload budget needs), and
// the truncation fact comes back as a second result instead of having to be
// inferred from the marker. Neither caps BEFORE sanitizing: their bound is
// on what comes back, not on the work done to produce it. SanitizeBudgeted
// and SanitizeSingleLineBudgeted are that same composition with the order
// reversed, for an upstream value whose size the caller does not control,
// where walking all of it is itself the exposure; Budget is their aggregate
// form, spending one shared byte budget across several values and keeping ONE
// truncation fact for the whole aggregate. Keeping a value's tail behind a
// prefixed marker remains the caller's own rune-boundary walk.
//
// Apply the sanitizer at the emit boundary (the slog call site, just before JSON
// encoding) so comparisons and dedupe keys keep operating on the raw value, and
// use one policy per application so two sinks cannot drift.
//
// A caller that also removes a KNOWN secret value from the same text redacts
// on both sides of the sanitizer and caps last, in the order
// redact, sanitize, redact, cap. The before position earns its place because
// the sanitizer rewrites an unsafe rune INSIDE a value, so a 32-character key
// carrying one interior unsafe rune (a C0, C1 or bidi control, and under
// SanitizeSingleLine also a CR or LF) reaches the sink as a near-complete
// fragment that a full-value needle can no longer match. The after position
// earns its place because a space and U+FFFD are the only runes the sanitizer
// can PRODUCE, so a secret containing either can be constructed out of text
// that did not match before it ran; an invalid UTF-8 byte becoming U+FFFD is
// the cheapest such case. The cap goes last so a value straddling the cut is
// already gone rather than sliced into a surviving prefix, and note that
// SanitizeSingleLineBounded, the Capped pair, the Budgeted pair and Budget all
// CONTAIN a cap and so delete bytes (Budget.Write caps BEFORE sanitizing),
// which leaves a redaction placed after one of them matching against a needle
// that cap may have split: redact before such a preset and again on its
// output, or compose the primitives. The full consumer-side contract is
// documented on httpx.RedactSecretString.
//
// The package is one deliberately small policy. It is not an HTML/XSS
// sanitizer, does not normalize Unicode (NFC/NFKC), case-folds and case-maps
// nothing, and does not touch zero-width or confusable runes; context-aware
// sinks (Markdown cells, URLs, HTML) need their own escaping on top of this
// rune policy.
package runesafe
