package api

import (
	"encoding/json"
	"testing"
)

// The flexible JSON types parse values Alegra serializes inconsistently (string
// vs number vs null). These fuzz tests assert the decoders never panic and that
// successfully-decoded values re-marshal cleanly. Without `-fuzz` they run the
// seed corpus as regular tests.

func FuzzID(f *testing.F) {
	for _, s := range []string{`"12"`, `12`, `12.0`, `12.5`, `null`, `""`, `"abc"`, `true`, `[]`} {
		f.Add([]byte(s))
	}
	f.Fuzz(func(t *testing.T, data []byte) {
		var id ID
		if json.Unmarshal(data, &id) != nil {
			return
		}
		if _, err := json.Marshal(id); err != nil {
			t.Fatalf("re-marshal failed for %q: %v", data, err)
		}
		_ = id.String()
	})
}

func FuzzInt(f *testing.F) {
	for _, s := range []string{`30`, `"30"`, `""`, `null`, `12.9`, `"abc"`, `{}`} {
		f.Add([]byte(s))
	}
	f.Fuzz(func(t *testing.T, data []byte) {
		var n Int
		if json.Unmarshal(data, &n) != nil {
			return
		}
		if _, err := json.Marshal(n); err != nil {
			t.Fatalf("re-marshal failed for %q: %v", data, err)
		}
	})
}

func FuzzMoney(f *testing.F) {
	for _, s := range []string{`1000`, `"1000.50"`, `""`, `null`, `"abc"`, `[1]`} {
		f.Add([]byte(s))
	}
	f.Fuzz(func(t *testing.T, data []byte) {
		var m Money
		if json.Unmarshal(data, &m) != nil {
			return
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
		_ = json.Unmarshal(data, &s) // must not panic regardless of input
	})
}
