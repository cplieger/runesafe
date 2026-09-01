package runesafe_test

import (
	"bytes"
	"encoding/json"
	jsonv2 "encoding/json/v2"
	"fmt"
	"log/slog"
	"strings"
	"testing"

	"github.com/cplieger/runesafe/v2"
)

// TestUntrustedSinkForms pins every emission form against the preset
// oracles — String, MarshalText, and LogValue must equal Sanitize;
// SingleLine must equal SanitizeSingleLine — while Raw round-trips the
// exact input bytes, including invalid UTF-8.
func TestUntrustedSinkForms(t *testing.T) {
	inputs := []string{
		"",
		"plain text",
		"葬送のフリーレン",
		"a\x1b[2Jb",
		"a\u009bb",
		"a\u202evil\u202cb",
		"a\u2028b\u2029c",
		"line1\nline2\rline3",
		"a\xffb",
	}
	for _, in := range inputs {
		t.Run(fmt.Sprintf("%q", in), func(t *testing.T) {
			u := runesafe.Untrusted(in)
			want := runesafe.Sanitize(in)
			if got := u.String(); got != want {
				t.Errorf("String() = %q, want Sanitize form %q", got, want)
			}
			text, err := u.MarshalText()
			if err != nil {
				t.Errorf("MarshalText() error: %v", err)
			}
			if got := string(text); got != want {
				t.Errorf("MarshalText() = %q, want Sanitize form %q", got, want)
			}
			v := u.LogValue()
			if v.Kind() != slog.KindString {
				t.Errorf("LogValue().Kind() = %v, want KindString", v.Kind())
			}
			if got := v.String(); got != want {
				t.Errorf("LogValue() = %q, want Sanitize form %q", got, want)
			}
			if got, wantStrict := u.SingleLine(), runesafe.SanitizeSingleLine(in); got != wantStrict {
				t.Errorf("SingleLine() = %q, want SanitizeSingleLine form %q", got, wantStrict)
			}
			if got := u.Raw(); got != in {
				t.Errorf("Raw() = %q, want exact input %q", got, in)
			}
		})
	}
}

// TestUntrustedJSONAsymmetry pins the decode-raw / encode-sanitized
// contract: a JSON document carrying a C1 introducer and a bidi override
// decodes into the tagged field byte-exact, and marshaling the same struct
// emits the sanitized form — including for a nested field, the class no
// sink-side rewriter can reach.
func TestUntrustedJSONAsymmetry(t *testing.T) {
	type inner struct {
		Title runesafe.Untrusted `json:"title"`
	}
	type doc struct {
		Title  runesafe.Untrusted `json:"title"`
		Nested inner              `json:"nested"`
	}
	raw := "Frieren\u009b\u202egpj.exe"
	blob := `{"title":"Frieren\u009b\u202egpj.exe","nested":{"title":"Frieren\u009b\u202egpj.exe"}}`
	var d doc
	if err := json.Unmarshal([]byte(blob), &d); err != nil {
		t.Fatalf("Unmarshal: %v", err)
	}
	if d.Title.Raw() != raw || d.Nested.Title.Raw() != raw {
		t.Fatalf("decode changed bytes: top %q nested %q, want %q", d.Title.Raw(), d.Nested.Title.Raw(), raw)
	}
	out, err := json.Marshal(d)
	if err != nil {
		t.Fatalf("Marshal: %v", err)
	}
	want := runesafe.Sanitize(raw)
	for _, r := range []rune{'\u009b', '\u202e'} {
		if bytes.ContainsRune(out, r) {
			t.Errorf("marshaled document carries raw %U: %s", r, out)
		}
	}
	var back struct {
		Title  string `json:"title"`
		Nested struct {
			Title string `json:"title"`
		} `json:"nested"`
	}
	if err := json.Unmarshal(out, &back); err != nil {
		t.Fatalf("re-decode: %v", err)
	}
	if back.Title != want || back.Nested.Title != want {
		t.Errorf("marshaled forms = %q / %q, want Sanitize form %q", back.Title, back.Nested.Title, want)
	}
}

// TestUntrustedSlogResolution proves the LogValuer fires in both built-in
// handlers, as a bare kv attr and inside a group: the encoded record
// carries the sanitized form and never the raw C1/bidi runes.
func TestUntrustedSlogResolution(t *testing.T) {
	u := runesafe.Untrusted("Frieren\u009b\u202egpj.exe")
	handlers := map[string]func(*bytes.Buffer) slog.Handler{
		"json": func(b *bytes.Buffer) slog.Handler { return slog.NewJSONHandler(b, nil) },
		"text": func(b *bytes.Buffer) slog.Handler { return slog.NewTextHandler(b, nil) },
	}
	for name, mk := range handlers {
		t.Run(name, func(t *testing.T) {
			var buf bytes.Buffer
			slog.New(mk(&buf)).Info("emit", "title", u, slog.Group("g", "title", u))
			got := buf.String()
			for _, r := range []rune{'\u009b', '\u202e'} {
				if strings.ContainsRune(got, r) {
					t.Errorf("%s record carries raw %U: %s", name, r, got)
				}
			}
			if !strings.Contains(got, "Frieren") {
				t.Errorf("%s record lost the safe text: %s", name, got)
			}
		})
	}
}

// TestUntrustedErrorConstruction pins the error-class coverage: an error
// built with fmt.Errorf("%s", v) carries the sanitized form at
// construction, before any sink sees it.
func TestUntrustedErrorConstruction(t *testing.T) {
	u := runesafe.Untrusted("bad request\u009b\nlevel=ERROR forged")
	err := fmt.Errorf("upstream said %s", u)
	msg := err.Error()
	if strings.ContainsRune(msg, '\u009b') {
		t.Errorf("error message carries raw C1: %q", msg)
	}
	if want := "upstream said " + runesafe.Sanitize(u.Raw()); msg != want {
		t.Errorf("Error() = %q, want %q", msg, want)
	}
}

// TestUntrustedComparableRaw pins the compute contract: equality and map
// keys operate on the raw bytes (two values differing only in an unsafe
// rune stay distinct), so matching and dedupe keep working on tagged
// fields without unwrapping.
func TestUntrustedComparableRaw(t *testing.T) {
	a, b := runesafe.Untrusted("x\u202ey"), runesafe.Untrusted("x y")
	if a == b {
		t.Error("raw-distinct values compare equal; equality must be on raw bytes")
	}
	if a.String() != b.String() {
		t.Errorf("sanitized forms differ: %q vs %q", a.String(), b.String())
	}
	m := map[runesafe.Untrusted]int{a: 1, b: 2}
	if len(m) != 2 {
		t.Errorf("map collapsed raw-distinct keys: %v", m)
	}
}

// TestUntrustedKeptCRLFCannotForgeRecord pins the doc-comment claim that
// the keepCRLF=true preset in LogValue is safe for BOTH built-in handlers:
// a kept CR/LF is escaped by the handler's own encoding (JSONHandler string
// escaping; TextHandler strconv quoting), so an upstream newline can never
// forge a second log record. An environmental-conformance pin, like
// TestClassifierUnicodeConformance: it guards the stdlib behavior the
// package's safety claim rests on across Go upgrades.
func TestUntrustedKeptCRLFCannotForgeRecord(t *testing.T) {
	u := runesafe.Untrusted("ok\r\nlevel=ERROR msg=forged")
	handlers := map[string]func(*bytes.Buffer) slog.Handler{
		"json": func(b *bytes.Buffer) slog.Handler { return slog.NewJSONHandler(b, nil) },
		"text": func(b *bytes.Buffer) slog.Handler { return slog.NewTextHandler(b, nil) },
	}
	for name, mk := range handlers {
		t.Run(name, func(t *testing.T) {
			var buf bytes.Buffer
			slog.New(mk(&buf)).Info("emit", "title", u)
			got := buf.String()
			if n := strings.Count(got, "\n"); n != 1 {
				t.Errorf("%s record carries %d newlines, want exactly 1 (the record terminator): %q", name, n, got)
			}
			if strings.Contains(got, "\r") {
				t.Errorf("%s record carries a raw CR: %q", name, got)
			}
			if !strings.HasSuffix(got, "\n") {
				t.Errorf("%s record does not end with the terminator: %q", name, got)
			}
		})
	}
}

// TestUntrustedMapKeySanitized pins the one emission path encoding/json used to
// bypass: a string-kinded map key. The v2-backed encoding/json routes the key
// through MarshalText, so a key is sanitized like every other sink. Measured
// across the bump on this test's own input: go1.26.7 emits {"k\u009b\u202e":1}
// and go1.27.0 emits {"k  ":1}, with no v1 option restoring the old bytes.
func TestUntrustedMapKeySanitized(t *testing.T) {
	key := runesafe.Untrusted("k\u009b\u202e")
	out, err := json.Marshal(map[runesafe.Untrusted]int{key: 1})
	if err != nil {
		t.Fatalf("Marshal: %v", err)
	}
	for _, r := range []rune{'\u009b', '\u202e'} {
		if bytes.ContainsRune(out, r) {
			t.Errorf("map key still emits raw %U: %s", r, out)
		}
	}
	if want := `{"` + key.String() + `":1}`; string(out) != want {
		t.Errorf("Marshal(map[Untrusted]int) = %s, want the MarshalText form %s", out, want)
	}
}

// TestUntrustedSurvivesJSONv2Strictness pins what the tag buys a consumer that
// moves its encoder to encoding/json/v2, whose defaults REJECT invalid UTF-8 in a
// JSON string instead of substituting U+FFFD the way v1 does. A plain string field
// carrying one invalid byte aborts the marshal with `jsontext: invalid UTF-8
// within "/v"` once the object is already partly written, while the tagged field
// emits cleanly and byte-identically to v1 — MarshalText hands the encoder the
// Sanitize form, always valid UTF-8 by construction, so the stricter boundary has
// nothing to refuse.
func TestUntrustedSurvivesJSONv2Strictness(t *testing.T) {
	const bad = "a\xffb\u009b\u202e"

	tagged := struct {
		V runesafe.Untrusted `json:"v"`
	}{V: runesafe.Untrusted(bad)}
	out, err := jsonv2.Marshal(tagged)
	if err != nil {
		t.Fatalf("json/v2 Marshal of a tagged field: %v", err)
	}
	if want := `{"v":"` + runesafe.Sanitize(bad) + `"}`; string(out) != want {
		t.Errorf("json/v2 Marshal = %s, want the Sanitize form %s", out, want)
	}
	if v1, err := json.Marshal(tagged); err != nil || string(v1) != string(out) {
		t.Errorf("v1 Marshal = %s (err %v), want v2's bytes %s", v1, err, out)
	}

	plain := struct {
		V string `json:"v"`
	}{V: bad}
	if _, err := jsonv2.Marshal(plain); err == nil {
		t.Error("json/v2 accepted invalid UTF-8 in an untagged string field; the tag's value here rests on it refusing")
	}
}
