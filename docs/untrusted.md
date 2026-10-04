# The Untrusted type

This page explains the `Untrusted` string type, for a developer deciding whether to tag decoded fields instead of sanitizing at each log call.

Sanitizing at each call site repeats one decision, that this text is untrusted, everywhere the text is emitted, and a forgotten call leaves no trace. `Untrusted` records the decision once, in the struct that decodes the upstream payload, where the decision is known.

```go
type Episode struct {
    Title runesafe.Untrusted `json:"title"`
}
```

## Decoding keeps the string value

`Untrusted` is a named string type with no `UnmarshalText` method, so `encoding/json` uses its normal string decoding and runesafe's sanitizer does not run. `Raw` returns that decoded string unchanged.

## Every standard sink emits it sanitized

- slog: `LogValue` implements `slog.LogValuer`, so a tagged value resolves to its `Sanitize` form in every handler, inside groups too.
- `fmt` and errors: `String` implements `fmt.Stringer`, so `%s`, `%v`, `%q` and `fmt.Errorf("upstream said %s", v)` render sanitized text. An error built this way carries no escape introducers from the moment it is created.
- Encoders: `MarshalText` implements `encoding.TextMarshaler`, so `encoding/json` emits the sanitized form at any nesting depth. A map key is sanitized too, because Go 1.27's `encoding/json` routes a string-kinded map key through `MarshalText`.

All three use `Sanitize`, which keeps CR and LF. That is correct for JSON, and slog's `TextHandler` quotes a value that holds them. For an error bound for a hand-built single-line sink that escapes nothing, build the message from `SingleLine` instead.

runesafe ships no slog handler or `ReplaceAttr` hook. slog handlers turn an error into text inside the encoder, after any attribute rewriting. So an error is covered only when it was built from sanitized text, and the type does that for you.

For text with nothing to replace, `LogValue`, `String`, `SingleLine` and `Raw` allocate nothing. `MarshalText` copies once, because it returns a `[]byte`.

## Raw for matching

Compute paths read the unsanitized value through `Raw`:

```go
slog.Warn("better release available", "title", ep.Title) // sanitized automatically
if ep.Title.Raw() == stored.Title { /* matching stays raw */ }
```

Use `Raw` for matching, dedupe keys and context-aware escapers. A `string(v)` conversion gives the same bytes but drops the tag without a trace, while `Raw` keeps intentional unwrapping easy to search for.

Matching on `Raw` compares bytes, so two values that differ only in an unsafe rune stay different. A cap that bounds the emitted form belongs on the sanitized text, `CapBytes(v.String(), n)`, because sanitizing can grow the raw bytes.

## Two rules for using it

- Store `Raw` in plain `string` fields in any struct your program saves and reads back. `MarshalText` runs inside every `json.Marshal`, so a tagged field in a state file would come back sanitized, not raw.
- Keep sanitizing at construction for text that must be safe in every future sink, such as a captured error body inside a returned value. A Markdown cell, a URL or an HTML page still needs its own escaping, built on `Raw`.
