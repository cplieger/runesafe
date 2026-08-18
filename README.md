# runesafe

[![Go Reference](https://pkg.go.dev/badge/github.com/cplieger/runesafe/v2.svg)](https://pkg.go.dev/github.com/cplieger/runesafe/v2)
[![Go version](https://img.shields.io/github/go-mod/go-version/cplieger/runesafe)](https://github.com/cplieger/runesafe/blob/main/go.mod)
[![Test coverage](https://img.shields.io/endpoint?url=https://raw.githubusercontent.com/cplieger/runesafe/badges/coverage.json)](https://github.com/cplieger/runesafe/actions/workflows/coverage.yml)
[![OpenSSF Best Practices](https://www.bestpractices.dev/projects/13657/badge)](https://www.bestpractices.dev/projects/13657)
[![OpenSSF Scorecard](https://api.scorecard.dev/projects/github.com/cplieger/runesafe/badge)](https://scorecard.dev/viewer/?uri=github.com/cplieger/runesafe)

> One rune-safety policy for untrusted upstream text headed to slog, JSON, or rendered output

A standalone Go library for a boundary every scraper-adjacent app meets: text the app did not author (API response fields, upstream error messages, file names, titles) that will be emitted into a log line, a JSON document, or a rendered report. Four rune classes survive that trip with their control semantics intact and let the upstream author forge or garble what the operator sees:

- **C0 controls (U+0000–U+001F) and DEL (U+007F)**: terminal escape sequences and log-record forgery. ESC introduces CSI/OSC sequences that can retitle the terminal, clear the screen, or write to the clipboard; a raw newline splits one record into two, fabricating a whole log line. CR/LF are optionally kept for sinks whose encoder escapes them (JSON).
- **C1 controls (U+0080–U+009F)**: single-rune escape introducers (CSI U+009B, OSC U+009D, …) with the same terminal powers as ESC sequences. `encoding/json` and slog's `JSONHandler` emit them **raw**, so escaping C0 alone does not close the terminal-injection hole.
- **Unicode Bidi_Control format characters** (U+061C, U+200E/U+200F, U+202A–U+202E, U+2066–U+2069): visually reorder rendered text, Trojan-Source-style, so a link or verdict reads differently than it compares. The set matches `unicode.Bidi_Control` exactly.
- **Line and paragraph separators U+2028/U+2029**: legal unescaped in JSON but line terminators to JavaScript and many viewers, so they split records like a raw newline.

runesafe classifies these runes and provides the shared sanitizers, so an app's slog emitter, JSON report writer, and renderer apply an identical policy instead of three drifting hand-rolled ones.

Standard library only, zero dependencies.

## Install

```sh
go get github.com/cplieger/runesafe/v2@latest
```

## Usage

### The shared sanitizer

Apply at the emit boundary (the slog call site, or just before JSON encoding) so comparisons and dedupe keys keep operating on the raw value:

```go
slog.Warn("better release available",
    "title", runesafe.Sanitize(upstream.Title),
    "group", runesafe.Sanitize(upstream.Group))
```

Each unsafe rune becomes a space; CR and LF pass through (slog's JSON encoder escapes them); invalid UTF-8 bytes become U+FFFD, so the result is always valid UTF-8. Sanitizing is idempotent, so double-sanitizing at two layers is harmless.

### Single-line sinks

For a sink where a raw newline forges a record boundary (a plain-text log line, a one-line error message, a rendered table cell), `SanitizeSingleLine` applies the strict policy (CR and LF become spaces too):

```go
msg := runesafe.SanitizeSingleLine(upstreamErr.Error())
```

Choose per sink: `Sanitize` when the sink's encoder provably escapes CR/LF (JSON), `SanitizeSingleLine` when it does not.

### Bounding sanitized text

Sanitizing can grow a string (each invalid UTF-8 byte becomes the three-byte U+FFFD), so a byte cap belongs **after** sanitizing. A naive `s[:n]` can then split a multi-byte rune, leaving a partial-rune tail whose raw 0x80–0x9F bytes a non-UTF-8 terminal reads as C1 escape introducers: the very class the sanitizer just removed. `CapBytes` cuts on a rune boundary:

```go
body := runesafe.CapBytes(runesafe.Sanitize(raw), maxBodyBytes)
```

For the common log-attribute case (single-line, capped, visibly marked), `SanitizeSingleLineBounded` packages the composition: `SanitizeSingleLine`, then `CapBytes` on the sanitized form, then `"..."` appended **outside** the cap. `n` budgets the retained body, so a truncated result is at most n+3 bytes; a within-cap result comes back byte-identical, with no marker. Truncated output always ends in the marker, but the converse does not hold (input may itself end in `...`); a caller that must know whether truncation occurred composes the primitives itself. A non-positive `n` yields `"..."` alone for non-empty input, and `""` stays `""`:

```go
slog.Warn("upstream rejected request",
    "reason", runesafe.SanitizeSingleLineBounded(upstreamErr.Error(), 200))
```

When the cap is a real budget — a record persisted under a write limit, a payload assembled under a vendor byte cap, a fixed-width column — `SanitizeCapped` (CR/LF kept, for JSON sinks) and `SanitizeSingleLineCapped` (strict) take the marker from the caller and charge it **against** `n`, so the returned text never exceeds the limit, and they return the truncation fact rather than leaving it to be inferred from the marker:

```go
attr, cut := runesafe.SanitizeSingleLineCapped(key, maxLoggedKeyBytes, "...")
slog.Warn("unknown upstream keys", "key", attr, "key_truncated", cut)
```

A cap too small to hold the marker drops the marker rather than emitting a fragment of it (`cut` still reports the elision); an empty marker caps silently. The marker is emitted verbatim, so build it from program text, not from untrusted input.

### Bounding the sanitizer's work

Everything above bounds what comes back. When the value's size is not the caller's to control — an upstream response field, an error message interpolating one, a file name — the bound has to cover the work as well: walking a multi-megabyte value with `strings.Map` inside a memory-limited process is a work-amplification denial of service ([CWE-400](https://cwe.mitre.org/data/definitions/400.html)) whatever the output cap is. `SanitizeBudgeted` and `SanitizeSingleLineBudgeted` are the `Capped` pair with the cap moved **ahead** of the sanitizer: the raw bytes are cut on a rune boundary first, only that chunk is sanitized, the growth sanitizing can cause is re-capped, and the caller's marker is still charged inside `n`.

```go
label, cut := runesafe.SanitizeSingleLineBudgeted(upstreamLabel, 64, "...(truncated)")
slog.Warn("skipped block", "type", label, "type_truncated", cut)
```

For a value carrying no unsafe rune the two orders are byte-identical, cut for cut, so moving a call site onto this one changes nothing an honest value emits. They diverge for exactly one input class: a value whose **raw** form exceeds `n` but whose **sanitized** form would have fitted, which takes enough multi-byte unsafe runes to collapse it under the cap (each becomes a single-byte space). This order cuts and marks such a value; the mark is honest either way, because bytes really were dropped before sanitizing.

`Budget` is the aggregate form, for an attribute assembled from several untrusted values that must share ONE byte budget and report ONE truncation fact. Joining first and capping afterwards allocates the whole untrusted aggregate before the bound applies, which is the same amplification one value at a time:

```go
b := runesafe.NewBudget(maxAttrBytes, "...")
for i, group := range upstream.Groups {
    if i > 0 && !b.Write(", ") {
        break
    }
    if !b.Write(group) { // false once anything has been dropped
        break
    }
}
attr, cut := b.Result()
slog.Warn("better release available", "groups", attr, "groups_truncated", cut)
```

Separators go through `Write` too, so a hostile value **count** cannot grow the attribute past the budget either. The marker is appended once for the whole aggregate, charged inside the budget, so `Result` never exceeds `max(n, 0)` bytes. `Write` reports whether the aggregate is still whole rather than whether the budget has room: a caller must attempt the write it cannot fit, because the refusal is what latches the fact.

### Rune classification

Two predicates expose the policy rune-by-rune, one per CR/LF policy, so a call
site names its sink instead of passing a boolean:

```go
runesafe.IsUnsafeMultiLine('\x1b')  // true:  ESC is always unsafe
runesafe.IsUnsafeMultiLine('\n')    // false: the sink's encoder escapes it
runesafe.IsUnsafeSingleLine('\n')   // true:  single-line sink; a newline forges a record
```

### A custom replacement policy

For a sink that needs a different replacement (strip instead of space), compose the predicate for your sink:

```go
// Remove (rather than blank) unsafe runes for a compact identifier.
id = strings.Map(func(r rune) rune {
    if runesafe.IsUnsafeSingleLine(r) {
        return -1
    }
    return r
}, id)
```

For context-aware escapers that percent-encode rather than replace (for example a Markdown link-URL escaper that must keep the URL usable), `IsUnsafeNonASCII` classifies the above-ASCII subset of the policy: C1 controls, the Bidi_Control set, and U+2028/U+2029. A URL encoder already covers ASCII controls and whitespace itself:

```go
// Percent-encode the policy runes url.Parse accepts but a viewer must never see raw.
for _, r := range u {
    if runesafe.IsUnsafeNonASCII(r) {
        for _, b := range []byte(string(r)) {
            fmt.Fprintf(&out, "%%%02X", b)
        }
        continue
    }
    out.WriteRune(r)
}
```

`IsBidiControl` exposes just the Bidi_Control subset for policies scoped to reordering alone.

### Tagging provenance: the `Untrusted` type

Per-call sanitizing re-derives the same fact (this text is untrusted) at every emit site, and a forgotten wrap is invisible. `Untrusted` records the fact once, where it is actually known: the struct that decodes the upstream payload.

```go
type Episode struct {
    Title runesafe.Untrusted `json:"title"`
}
```

Decoding is untouched: a string-kinded named type unmarshals natively, raw bytes in. Emission fires the policy through the standard interfaces:

- slog resolves the value sanitized (`slog.LogValuer`).
- `fmt` renders it sanitized (`fmt.Stringer`), so `fmt.Errorf("upstream said %s", v)` carries no escape introducers from construction on; that is the one boundary that covers error values. This form keeps CR/LF, so route such an error's text through `SingleLine()` if it is ever bound for a hand-built sink that escapes nothing.
- `encoding/json` emits it sanitized at any nesting depth (`encoding.TextMarshaler`). Map keys are the exception: `encoding/json` reads a string-kinded key's bytes directly without calling `MarshalText`, so key marshaled documents by `v.String()`, never by the tagged value.

Compute paths keep the exact bytes via `Raw()`:

```go
slog.Warn("better release available", "title", ep.Title) // sanitized automatically
if ep.Title.Raw() == onDisk.Title.Raw() { /* matching stays raw */ }
```

Two rules keep it honest. Structs persisted for the program's own re-reading store `Raw()` in plain `string` fields: `MarshalText` fires inside every `json.Marshal`, so a tagged field in a state file would round-trip sanitized, not raw. And a `string(v)` conversion silently drops the tag; use `Raw()` so intentional unwrapping stays greppable.

## Migrating from v1

| v1 | v2 |
| --- | --- |
| `IsUnsafe(r, true)` | `IsUnsafeMultiLine(r)` |
| `IsUnsafe(r, false)` | `IsUnsafeSingleLine(r)` |
| `github.com/cplieger/runesafe` | `github.com/cplieger/runesafe/v2` |

## API

| Symbol | Contract |
| --- | --- |
| `Sanitize(s string) string` | Replaces each unsafe rune (multi-line policy, `IsUnsafeMultiLine`) with a space. Valid UTF-8 out, rune count preserved, idempotent. |
| `SanitizeSingleLine(s string) string` | The strict preset (`IsUnsafeSingleLine`): everything `Sanitize` replaces, plus CR and LF. |
| `SanitizeSingleLineBounded(s string, n int) string` | `SanitizeSingleLine`, then a rune-boundary cap of the sanitized form at n bytes with `"..."` appended outside the cap (truncated result ≤ n+3 bytes; within-cap input byte-identical, no marker). Non-positive n yields `"..."` for non-empty input; `""` stays `""`. |
| `SanitizeCapped(s string, n int, marker string) (text string, cut bool)` | `Sanitize`, then a rune-boundary cap with the caller's marker counted **inside** the cap: the text is always ≤ max(n, 0) bytes. `cut` is true exactly when the sanitized form was shortened. A cap below `len(marker)` drops the marker; the marker is emitted verbatim, never sanitized. |
| `SanitizeSingleLineCapped(s string, n int, marker string) (text string, cut bool)` | The strict twin of `SanitizeCapped` (CR and LF replaced too), same cap, marker and `cut` contract. |
| `SanitizeBudgeted(s string, n int, marker string) (text string, cut bool)` | `SanitizeCapped` with the cap **ahead** of the sanitizer: `s` is cut on a rune boundary at n bytes first, only that chunk is sanitized, its growth is re-capped, and the marker is charged inside n. Bounds the sanitizer's WORK, for a value whose size the caller does not control. Identical to `SanitizeCapped` for any value carrying no unsafe rune. |
| `SanitizeSingleLineBudgeted(s string, n int, marker string) (text string, cut bool)` | The strict twin of `SanitizeBudgeted` (CR and LF replaced too), same pre-cap, marker and `cut` contract. |
| `NewBudget(n int, marker string) *Budget` / `NewSingleLineBudget(...)` | One shared n-byte budget for SEVERAL untrusted values, one per CR/LF policy. Each value is capped before it is sanitized; the marker is charged inside n. |
| `(*Budget).Write(raw string) bool` | Append `raw`'s sanitized, pre-capped prefix against the remaining budget. Reports whether the aggregate is still whole — deliberately not "the budget has room", so a caller loops until a write is actually refused. Separators go through it too. |
| `(*Budget).Result() (text string, cut bool)` | The aggregate and the ONE truncation fact latched across every write, marked once and never longer than max(n, 0) bytes. Reads without spending; calling it twice returns the same pair. |
| `CapBytes(s string, n int) string` | Truncates to at most n bytes on a rune boundary; never ends in a partial rune. Non-positive n returns "". |
| `IsUnsafeMultiLine(r rune) bool` | One rune under the multi-line policy: C0 except CR/LF, DEL, C1, Bidi_Control, U+2028/U+2029. |
| `IsUnsafeSingleLine(r rune) bool` | The strict per-rune policy: everything `IsUnsafeMultiLine` refuses, plus CR and LF. The v1 `IsUnsafe(r, keepCRLF)` split into these two names so a call site reads its policy — and so deleting the bool cannot silently pick one. |
| `IsUnsafeNonASCII(r rune) bool` | The above-ASCII subset: C1, Bidi_Control, U+2028/U+2029. For escapers whose sink already covers ASCII (URL percent-encoders). |
| `IsBidiControl(r rune) bool` | Exactly `unicode.Is(unicode.Bidi_Control, r)`, without the table lookup. |
| `Untrusted` (string type) | Provenance tag for upstream text: decodes raw, emits `Sanitize`'d through `slog.LogValuer`, `fmt.Stringer`, and `encoding.TextMarshaler`. |
| `Untrusted.Raw() string` | The exact bytes as received: matching, dedupe keys, compute-side caps, and composed escapers operate on this. An emit-bound byte cap belongs on the sanitized form: `CapBytes(v.String(), n)`. |
| `Untrusted.SingleLine() string` | The strict `SanitizeSingleLine` form, for hand-built single-line sinks. |

## Adoption guidance

- **One policy per app.** Route every untrusted attribute through `Sanitize` / `SanitizeSingleLine` (or one app-local wrapper around `IsUnsafe`) so two sinks cannot disagree about what is dangerous.
- **Construction-time sanitization still has a place.** For text that must be safe unconditionally through every future sink, like a captured error body, sanitize at construction instead of tagging.
- **Context-aware sinks still need their own escaping.** A Markdown table cell, a link URL, or an HTML page has injection vectors this rune policy does not address (pipes, brackets, angle brackets); apply the sink's escaper on top.

## Unsupported by Design

| Feature | Rationale |
| --- | --- |
| HTML/XSS sanitization | Different threat model and sink; use a real HTML sanitizer. This policy is for logs, JSON, and rendered reports. |
| Unicode normalization (NFC/NFKC) | Normalization changes text identity; a display-safety policy must not. Normalize separately if the app needs it. |
| Zero-width and confusable runes | Invisible-character and homoglyph spoofing is a rabbit hole with legitimate-text collateral (ZWJ in emoji, real Cyrillic). The policy targets runes with control semantics, where replacement is always correct. |
| Configurable replacement rune | The space is the policy; a different replacement is a two-line `strings.Map` over `IsUnsafe` (shown above). |
| Removing instead of replacing | Deletion changes rune offsets and can splice adjacent fragments into new tokens; 1:1 replacement preserves shape. Compose `IsUnsafe` for the rare sink that genuinely wants removal. |

## Disclaimer

This project is built with care and follows security best practices, but it is intended for personal / self-hosted use. No guarantees of fitness for production environments. Use at your own risk.

This project was built with AI-assisted tooling using [Claude](https://claude.com), [GPT](https://openai.com), and [Kiro](https://kiro.dev). The human maintainer defines architecture, supervises implementation, and makes all final decisions.

## License

Apache-2.0. See [LICENSE](LICENSE).
