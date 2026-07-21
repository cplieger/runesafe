# runesafe

[![Go Reference](https://pkg.go.dev/badge/github.com/cplieger/runesafe.svg)](https://pkg.go.dev/github.com/cplieger/runesafe)
[![Go version](https://img.shields.io/github/go-mod/go-version/cplieger/runesafe)](https://github.com/cplieger/runesafe/blob/main/go.mod)
[![Test coverage](https://img.shields.io/endpoint?url=https://raw.githubusercontent.com/cplieger/runesafe/badges/coverage.json)](https://github.com/cplieger/runesafe/actions/workflows/coverage.yml)
[![OpenSSF Best Practices](https://www.bestpractices.dev/projects/13657/badge)](https://www.bestpractices.dev/projects/13657)
[![OpenSSF Scorecard](https://api.scorecard.dev/projects/github.com/cplieger/runesafe/badge)](https://scorecard.dev/viewer/?uri=github.com/cplieger/runesafe)

> One rune-safety policy for untrusted upstream text headed to slog, JSON, or rendered output

A standalone Go library for a boundary every scraper-adjacent app meets: text the app did not author — API response fields, upstream error messages, file names, titles — that will be emitted into a log line, a JSON document, or a rendered report. Four rune classes survive that trip with their control semantics intact and let the upstream author forge or garble what the operator sees:

- **C0 controls (U+0000–U+001F) and DEL (U+007F)** — terminal escape sequences (ESC introduces CSI/OSC sequences that can retitle the terminal, clear the screen, or write to the clipboard) and log-record forgery (a raw newline splits one record into two, fabricating a whole log line). CR/LF are optionally kept for sinks whose encoder escapes them (JSON).
- **C1 controls (U+0080–U+009F)** — single-rune escape introducers (CSI U+009B, OSC U+009D, …) with the same terminal powers as ESC sequences. `encoding/json` and slog's `JSONHandler` emit them **raw**, so escaping C0 alone does not close the terminal-injection hole.
- **Unicode Bidi_Control format characters** (U+061C, U+200E/U+200F, U+202A–U+202E, U+2066–U+2069) — visually reorder rendered text, Trojan-Source-style: a link or verdict reads differently than it compares. The set matches `unicode.Bidi_Control` exactly (drift-guarded by an exhaustive sweep test).
- **Line and paragraph separators U+2028/U+2029** — legal unescaped in JSON but line terminators to JavaScript and many viewers, so they split records like a raw newline.

runesafe classifies these runes and provides the shared sanitizers, so an app's slog emitter, JSON report writer, and renderer apply an identical policy instead of three drifting hand-rolled ones.

Standard library only, zero dependencies.

## Install

```sh
go get github.com/cplieger/runesafe@latest
```

## Usage

### The shared sanitizer

Apply at the emit boundary — the slog call site, or just before JSON encoding — so comparisons and dedupe keys keep operating on the raw value:

```go
slog.Warn("better release available",
    "title", runesafe.Sanitize(upstream.Title),
    "group", runesafe.Sanitize(upstream.Group))
```

Each unsafe rune becomes a space; CR and LF pass through (slog's JSON encoder escapes them); invalid UTF-8 bytes become U+FFFD, so the result is always valid UTF-8. Sanitizing is idempotent — double-sanitizing at two layers is harmless.

### Single-line sinks

For a sink where a raw newline forges a record boundary — a plain-text log line, a one-line error message, a rendered table cell — `SanitizeSingleLine` applies the strict policy (CR and LF become spaces too):

```go
msg := runesafe.SanitizeSingleLine(upstreamErr.Error())
```

Choose per sink: `Sanitize` when the sink's encoder provably escapes CR/LF (JSON), `SanitizeSingleLine` when it does not.

### Bounding sanitized text

Sanitizing can grow a string (each invalid UTF-8 byte becomes the three-byte U+FFFD), so a byte cap belongs **after** sanitizing — and a naive `s[:n]` can split a multi-byte rune, leaving a partial-rune tail whose raw 0x80–0x9F bytes a non-UTF-8 terminal reads as C1 escape introducers: the very class the sanitizer just removed. `CapBytes` cuts on a rune boundary:

```go
body := runesafe.CapBytes(runesafe.Sanitize(raw), maxBodyBytes)
```

For the common log-attribute case — single-line, capped, visibly marked — `SanitizeSingleLineBounded` packages the composition: `SanitizeSingleLine`, then `CapBytes` on the sanitized form, then `"..."` appended **outside** the cap. `n` budgets the retained body, so a truncated result is at most n+3 bytes; a within-cap result comes back byte-identical, with no marker. Truncated output always ends in the marker, but the converse does not hold (input may itself end in `...`); a caller that must know whether truncation occurred composes the primitives itself. A non-positive `n` yields `"..."` alone for non-empty input, and `""` stays `""`:

```go
slog.Warn("upstream rejected request",
    "reason", runesafe.SanitizeSingleLineBounded(upstreamErr.Error(), 200))
```

### Rune classification

`IsUnsafe` exposes the policy rune-by-rune, with an explicit CR/LF switch:

```go
runesafe.IsUnsafe('\x1b', true)  // true  — ESC is always unsafe
runesafe.IsUnsafe('\n', true)    // false — the sink's encoder escapes it
runesafe.IsUnsafe('\n', false)   // true  — single-line sink: a newline forges a record
```

### A custom replacement policy

For a sink that needs a different replacement — strip instead of space — compose `IsUnsafe` yourself:

```go
// Remove (rather than blank) unsafe runes for a compact identifier.
id = strings.Map(func(r rune) rune {
    if runesafe.IsUnsafe(r, false) {
        return -1
    }
    return r
}, id)
```

For context-aware escapers that percent-encode rather than replace (for example a Markdown link-URL escaper that must keep the URL usable), `IsUnsafeNonASCII` classifies the above-ASCII subset of the policy — C1 controls, the Bidi_Control set, U+2028/U+2029 — since a URL encoder already covers ASCII controls and whitespace itself:

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

Per-call sanitizing re-derives the same fact — this text is untrusted — at every emit site, and a forgotten wrap is invisible. `Untrusted` records the fact once, where it is actually known: the struct that decodes the upstream payload.

```go
type Episode struct {
    Title runesafe.Untrusted `json:"title"`
}
```

Decoding is untouched (a string-kinded named type unmarshals natively, raw bytes in). Emission fires the policy through the standard interfaces — slog resolves the value sanitized (`slog.LogValuer`), `fmt` renders it sanitized (`fmt.Stringer`, so `fmt.Errorf("upstream said %s", v)` carries no escape introducers from construction on, the one boundary that covers error values; the form keeps CR/LF, so route such an error's text through `SingleLine()` if it is ever bound for a hand-built sink that escapes nothing), and `encoding/json` emits it sanitized at any nesting depth (`encoding.TextMarshaler`; map keys are the exception -- `encoding/json` reads a string-kinded key's bytes directly without calling `MarshalText`, so key marshaled documents by `v.String()`, never by the tagged value). Compute paths keep the exact bytes via `Raw()`:

```go
slog.Warn("better release available", "title", ep.Title) // sanitized automatically
if ep.Title.Raw() == onDisk.Title.Raw() { /* matching stays raw */ }
```

Two rules keep it honest. Structs persisted for the program's own re-reading store `Raw()` in plain `string` fields — `MarshalText` fires inside every `json.Marshal`, so a tagged field in a state file would round-trip sanitized, not raw. And a `string(v)` conversion silently drops the tag; use `Raw()` so intentional unwrapping stays greppable.

## API

| Symbol | Contract |
| --- | --- |
| `Sanitize(s string) string` | Replaces each unsafe rune (keepCRLF=true policy) with a space. Valid UTF-8 out, rune count preserved, idempotent. |
| `SanitizeSingleLine(s string) string` | The strict preset (keepCRLF=false): everything `Sanitize` replaces, plus CR and LF. |
| `SanitizeSingleLineBounded(s string, n int) string` | `SanitizeSingleLine`, then a rune-boundary cap of the sanitized form at n bytes with `"..."` appended outside the cap (truncated result ≤ n+3 bytes; within-cap input byte-identical, no marker). Non-positive n yields `"..."` for non-empty input; `""` stays `""`. |
| `CapBytes(s string, n int) string` | Truncates to at most n bytes on a rune boundary; never ends in a partial rune. Non-positive n returns "". |
| `IsUnsafe(r rune, keepCRLF bool) bool` | One rune under the policy: C0 (CR/LF exempt when keepCRLF), DEL, C1, Bidi_Control, U+2028/U+2029. |
| `IsUnsafeNonASCII(r rune) bool` | The above-ASCII subset: C1, Bidi_Control, U+2028/U+2029. For escapers whose sink already covers ASCII (URL percent-encoders). |
| `IsBidiControl(r rune) bool` | Exactly `unicode.Is(unicode.Bidi_Control, r)`, without the table lookup. |
| `Untrusted` (string type) | Provenance tag for upstream text: decodes raw, emits `Sanitize`'d through `slog.LogValuer`, `fmt.Stringer`, and `encoding.TextMarshaler`. |
| `Untrusted.Raw() string` | The exact bytes as received — matching, dedupe keys, compute-side caps, and composed escapers operate on this. An emit-bound byte cap belongs on the sanitized form: `CapBytes(v.String(), n)`. |
| `Untrusted.SingleLine() string` | The strict `SanitizeSingleLine` form, for hand-built single-line sinks. |

## Origin

Extracted from [seadex-scout](https://github.com/cplieger/seadex-scout)'s `internal/textsafe`, where one policy is shared by the audit report's renderers, the daemon's finding emitter, and the upstream-error-message sanitizer — three sinks that previously risked drifting on which runes they considered dangerous. The classification (C1 controls and bidi controls pass slog's JSON escaping raw) came out of an adversarial review of what actually reaches a terminal when JSON logs shipped to a log store are viewed. `CapBytes` hardening originated in [arrapi](https://github.com/cplieger/arrapi)'s captured-error-body path, where a post-sanitize re-cap had to avoid re-minting C1 tail bytes.

## Adoption guidance

- **Sanitize at the emit boundary, not at parse time.** Keep the raw value for matching, dedupe keys, and comparisons; sanitize the copy that leaves the process for human eyes. Tagging a field `Untrusted` implements exactly this split with no per-site calls: raw in, sanitized out at every standard sink.
- **Machine-read persistence stores `Raw()`.** A tagged field inside a struct written to a state file round-trips sanitized (`MarshalText` fires in `json.Marshal`); state structs keep plain `string` fields populated from `Raw()`. Construction-time sanitization also remains the right boundary for text that must be safe unconditionally through every future sink, like a captured error body.
- **One policy per app.** Route every untrusted attribute through `Sanitize` / `SanitizeSingleLine` (or one app-local wrapper around `IsUnsafe`) so two sinks cannot disagree about what is dangerous.
- **Context-aware sinks still need their own escaping.** A Markdown table cell, a link URL, or an HTML page has injection vectors this rune policy does not address (pipes, brackets, angle brackets); apply the sink's escaper on top.
- **Choose the preset per sink, not globally.** JSON sinks keep CR/LF (the encoder escapes them); single-line sinks must not.
- **Cap after sanitizing, on a rune boundary.** A byte cap applied before sanitizing can be outgrown by U+FFFD replacement; a cap applied with a raw slice can split a rune. `CapBytes` after `Sanitize` is the safe order.

## Unsupported by design

| Feature | Rationale |
| --- | --- |
| HTML/XSS sanitization | Different threat model and sink; use a real HTML sanitizer. This policy is for logs, JSON, and rendered reports. |
| Unicode normalization (NFC/NFKC) | Normalization changes text identity; a display-safety policy must not. Normalize separately if the app needs it. |
| Zero-width and confusable runes | Invisible-character and homoglyph spoofing is a rabbit hole with legitimate-text collateral (ZWJ in emoji, real Cyrillic). The policy targets runes with control semantics, where replacement is always correct. |
| Configurable replacement rune | The space is the policy; a different replacement is a two-line `strings.Map` over `IsUnsafe` (shown above). |
| Removing instead of replacing | Deletion changes rune offsets and can splice adjacent fragments into new tokens; 1:1 replacement preserves shape. Compose `IsUnsafe` for the rare sink that genuinely wants removal. |

## Disclaimer

This project is built with care and follows security best practices, but it is intended for personal / self-hosted use. No guarantees of fitness for production environments. Use at your own risk.

This project was built with AI-assisted tooling using [Claude Opus](https://www.anthropic.com/claude) and [Kiro](https://kiro.dev). The human maintainer defines architecture, supervises implementation, and makes all final decisions.

## License

GPL-3.0 — see [LICENSE](LICENSE).
