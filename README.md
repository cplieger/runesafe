# runesafe

[![Go Reference](https://pkg.go.dev/badge/github.com/cplieger/runesafe.svg)](https://pkg.go.dev/github.com/cplieger/runesafe)
[![Go version](https://img.shields.io/github/go-mod/go-version/cplieger/runesafe)](https://github.com/cplieger/runesafe/blob/main/go.mod)
[![Test coverage](https://img.shields.io/endpoint?url=https://raw.githubusercontent.com/cplieger/runesafe/badges/coverage.json)](https://github.com/cplieger/runesafe/actions/workflows/coverage.yml)
[![Mutation](https://img.shields.io/endpoint?url=https://raw.githubusercontent.com/cplieger/runesafe/badges/mutation.json)](https://github.com/cplieger/runesafe/issues?q=label%3Agremlins-tracker)
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

`IsBidiControl` exposes the Bidi_Control subset for context-aware escapers that percent-encode rather than replace (for example a Markdown link-URL escaper that must keep the URL usable).

## API

| Symbol | Contract |
| --- | --- |
| `Sanitize(s string) string` | Replaces each unsafe rune (keepCRLF=true policy) with a space. Valid UTF-8 out, rune count preserved, idempotent. |
| `SanitizeSingleLine(s string) string` | The strict preset (keepCRLF=false): everything `Sanitize` replaces, plus CR and LF. |
| `CapBytes(s string, n int) string` | Truncates to at most n bytes on a rune boundary; never ends in a partial rune. Non-positive n returns "". |
| `IsUnsafe(r rune, keepCRLF bool) bool` | One rune under the policy: C0 (CR/LF exempt when keepCRLF), DEL, C1, Bidi_Control, U+2028/U+2029. |
| `IsBidiControl(r rune) bool` | Exactly `unicode.Is(unicode.Bidi_Control, r)`, without the table lookup. |

## Origin

Extracted from [seadex-scout](https://github.com/cplieger/seadex-scout)'s `internal/textsafe`, where one policy is shared by the audit report's renderers, the daemon's finding emitter, and the upstream-error-message sanitizer — three sinks that previously risked drifting on which runes they considered dangerous. The classification (C1 controls and bidi controls pass slog's JSON escaping raw) came out of an adversarial review of what actually reaches a terminal when JSON logs shipped to a log store are viewed. `CapBytes` hardening originated in [arrapi](https://github.com/cplieger/arrapi)'s captured-error-body path, where a post-sanitize re-cap had to avoid re-minting C1 tail bytes.

## Adoption guidance

- **Sanitize at the emit boundary, not at parse time.** Keep the raw value for matching, dedupe keys, and comparisons; sanitize the copy that leaves the process for human eyes.
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
