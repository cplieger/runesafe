// Package runesafe classifies runes that are unsafe in untrusted text bound
// for logs, JSON or rendered output, and provides shared sanitizers that
// neutralize them.
//
// Untrusted upstream text eventually reaches a slog line read in a terminal, a
// JSON report or a Markdown table. Four classes of rune keep their control
// semantics on that trip and let the upstream author forge or garble what the
// operator sees:
//
//   - C0 controls and DEL: terminal escape sequences, and a raw newline that
//     splits one log record into two.
//   - C1 controls: single-rune escape introducers that encoding/json and slog's
//     JSONHandler emit raw.
//   - Unicode Bidi_Control characters, which visually reorder rendered text.
//   - U+2028 and U+2029, line terminators to JavaScript and many viewers.
//
// The classes are written as code points, never looked up in a Unicode table,
// so a Unicode upgrade cannot move the policy.
//
// [IsUnsafeMultiLine] and [IsUnsafeSingleLine] classify one rune under the two
// CR/LF policies. [Sanitize] keeps CR and LF for JSON-encoded sinks, and
// [SanitizeSingleLine] replaces them too, each unsafe rune becoming a space.
// [CapBytes], [SanitizeSingleLineBounded], the Capped and Budgeted pairs and
// [Budget] bound the result. The [Untrusted] string type applies Sanitize
// automatically in every standard sink while [Untrusted.Raw] keeps the exact
// bytes.
//
// Sanitize at the emit boundary, so comparisons and dedupe keys keep the raw
// value. A caller that also removes a known secret redacts, sanitizes, redacts
// again and caps last.
//
// It is not an HTML sanitizer and does not normalize Unicode. docs/policy.md,
// docs/capping.md, docs/untrusted.md and docs/non-goals.md hold the full
// contract.
package runesafe
