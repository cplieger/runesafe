package runesafe

import "log/slog"

// Untrusted marks a string as untrusted upstream text at the moment it
// enters the program — an API response field, an upstream error message, a
// file name, a title — so every human-facing sink it later reaches applies
// the package's rune policy automatically. The type makes provenance
// portable: the trust decision is recorded once, in the DTO struct that
// decodes the upstream payload, instead of re-derived at every emit site.
//
// Ingestion is free. A string-kinded named type without an UnmarshalText
// method decodes natively (encoding/json assigns into it directly), so
// tagging a decode-struct field preserves the raw bytes exactly as
// received. Emission is where the policy fires, through the standard
// interfaces:
//
//   - slog: LogValue implements slog.LogValuer, so a value passed as a
//     bare attr ("title", v) resolves to its Sanitize'd form in every
//     handler, through groups, before encoding.
//   - fmt and errors: String implements fmt.Stringer, so %s, %v, %q — and
//     fmt.Errorf("upstream said %s", v) — render sanitized text. An error
//     built this way is safe at construction, the one boundary that covers
//     error values (slog handlers stringify errors inside the encoder,
//     after any attribute rewriting).
//   - encoders: MarshalText implements encoding.TextMarshaler, so
//     encoding/json and any TextMarshaler-aware encoder emit the
//     Sanitize'd form, however deeply the value nests in a document.
//
// Raw returns the exact bytes for the paths that must not be transformed:
// matching, dedupe keys, byte caps, context-aware escapers. A plain
// string(v) conversion yields the same bytes but silently drops the tag;
// prefer Raw so intentional unwrapping stays greppable.
//
// Two rules keep the type honest:
//
//   - Machine-read persistence stores Raw. MarshalText fires inside every
//     json.Marshal, so a tagged field written to a state file and read
//     back would round-trip sanitized-not-raw. Structs persisted for the
//     program's own consumption keep plain string fields, populated via
//     Raw at construction; the tagged form is for human-facing documents
//     only.
//   - Untrusted does not replace construction-time sanitization for text
//     that must be safe unconditionally through every future sink (a
//     captured error body embedded in a returned value), and context-aware
//     sinks (Markdown cells, URLs, HTML) still need their own escaping on
//     top of the rune policy, composed over Raw.
//
// LogValue, String, and MarshalText apply Sanitize (the keepCRLF=true
// preset): correct for JSON sinks, and safe for slog's TextHandler, whose
// quoting escapes a kept CR or LF. For a hand-built single-line sink whose
// encoder escapes nothing, use SingleLine explicitly.
type Untrusted string

// LogValue implements slog.LogValuer: a tagged attr value resolves to its
// Sanitize'd form in every handler before encoding.
func (u Untrusted) LogValue() slog.Value {
	return slog.StringValue(Sanitize(string(u)))
}

// String implements fmt.Stringer: %s, %v, %q, and fmt.Errorf render the
// Sanitize'd form, so an error wrapping the value is safe at construction.
func (u Untrusted) String() string {
	return Sanitize(string(u))
}

// MarshalText implements encoding.TextMarshaler: encoding/json and any
// TextMarshaler-aware encoder emit the Sanitize'd form at any nesting
// depth. Decoding is deliberately untouched (no UnmarshalText), so raw
// bytes survive ingestion; see the type comment for the machine-read
// persistence rule this asymmetry imposes.
func (u Untrusted) MarshalText() ([]byte, error) {
	return []byte(Sanitize(string(u))), nil
}

// SingleLine returns the SanitizeSingleLine'd form, for hand-built
// single-line sinks whose encoder does not escape CR/LF.
func (u Untrusted) SingleLine() string {
	return SanitizeSingleLine(string(u))
}

// Raw returns the exact bytes as received, for matching, dedupe keys, byte
// caps, and context-aware escapers. Prefer Raw over a string conversion so
// intentional unwrapping stays greppable.
func (u Untrusted) Raw() string {
	return string(u)
}
