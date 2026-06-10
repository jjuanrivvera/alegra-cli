package api

import (
	"bytes"
	"encoding/json"
	"regexp"
	"strconv"
	"testing"
)

// The flexible JSON types parse values Alegra serializes inconsistently (string
// vs number vs null). Beyond never panicking, the fuzzers assert value-level
// properties: weak "doesn't crash" checks let a real bug through once (ids
// above 2^53 were silently rounded through float64), so each decoder pins the
// strongest property that holds for its input class. Without `-fuzz` they run
// the seed corpus as regular tests.

// jsonIntLiteral matches a bare JSON integer token (no exponent, no fraction,
// no leading zeros — shapes json.Unmarshal would reject anyway). "-0" is
// excluded: normalizing it to "0" is intended, not precision loss.
var jsonIntLiteral = regexp.MustCompile(`^(-?[1-9][0-9]*|0)$`)

func FuzzID(f *testing.F) {
	for _, s := range []string{`"12"`, `12`, `12.0`, `12.5`, `null`, `""`, `"abc"`, `true`, `[]`, `9007199254740993`} {
		f.Add([]byte(s))
	}
	f.Fuzz(func(t *testing.T, data []byte) {
		var id ID
		if json.Unmarshal(data, &id) != nil {
			return
		}
		// An integer literal must survive decoding verbatim: any difference
		// means precision was lost on the way through.
		if tok := bytes.TrimSpace(data); jsonIntLiteral.Match(tok) {
			if id.String() != string(tok) {
				t.Fatalf("integer literal %s decoded to %q", tok, id)
			}
		}
		if _, err := json.Marshal(id); err != nil {
			t.Fatalf("re-marshal failed for %q: %v", data, err)
		}
	})
}

func FuzzInt(f *testing.F) {
	for _, s := range []string{`30`, `"30"`, `""`, `null`, `12.9`, `"abc"`, `{}`, `9007199254740993`, `"9007199254740993"`} {
		f.Add([]byte(s))
	}
	f.Fuzz(func(t *testing.T, data []byte) {
		var n Int
		if json.Unmarshal(data, &n) != nil {
			return
		}
		// Bare or quoted integer literals within int64 range must decode to
		// the exact value, not a float64 approximation.
		tok := bytes.TrimSpace(data)
		if tok[0] == '"' {
			var s string
			if json.Unmarshal(tok, &s) == nil && jsonIntLiteral.MatchString(s) {
				tok = []byte(s)
			}
		}
		if jsonIntLiteral.Match(tok) {
			if want, err := strconv.ParseInt(string(tok), 10, 64); err == nil && int64(n) != want {
				t.Fatalf("integer literal %s decoded to %d, want %d", tok, int64(n), want)
			}
		}
		if _, err := json.Marshal(n); err != nil {
			t.Fatalf("re-marshal failed for %q: %v", data, err)
		}
	})
}

func FuzzMoney(f *testing.F) {
	for _, s := range []string{`1000`, `"1000.50"`, `""`, `null`, `"abc"`, `[1]`, `-42.5`} {
		f.Add([]byte(s))
	}
	f.Fuzz(func(t *testing.T, data []byte) {
		var m Money
		if json.Unmarshal(data, &m) != nil {
			return
		}
		// A bare JSON number must decode to exactly what ParseFloat yields —
		// Money is a float64, so this is the strongest property available.
		tok := bytes.TrimSpace(data)
		if len(tok) > 0 && tok[0] != '"' && string(tok) != "null" {
			if want, err := strconv.ParseFloat(string(tok), 64); err == nil && float64(m) != want {
				t.Fatalf("number %s decoded to %v, want %v", tok, float64(m), want)
			}
		}
		if _, err := json.Marshal(m); err != nil {
			t.Fatalf("re-marshal failed for %q: %v", data, err)
		}
	})
}

func FuzzStringOrSlice(f *testing.F) {
	for _, s := range []string{`"client"`, `["client","provider"]`, `null`, `123`, `[1,2]`} {
		f.Add([]byte(s))
	}
	f.Fuzz(func(t *testing.T, data []byte) {
		var s StringOrSlice
		if json.Unmarshal(data, &s) != nil {
			return
		}
		// A single JSON string must decode to exactly that one element.
		tok := bytes.TrimSpace(data)
		if len(tok) > 0 && tok[0] == '"' {
			var want string
			if json.Unmarshal(tok, &want) == nil {
				if len(s) != 1 || s[0] != want {
					t.Fatalf("string %s decoded to %v", tok, []string(s))
				}
			}
		}
	})
}
