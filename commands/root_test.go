package commands

import (
	"bytes"
	"strings"
	"testing"
)

func TestWarnInsecureBaseURL(t *testing.T) {
	cases := []struct {
		url      string
		wantWarn bool
	}{
		{"https://api.alegra.com/api/v1", false},
		{"http://localhost:8080", false},
		{"http://127.0.0.1:9000/api", false},
		{"http://[::1]:9000", false},
		{"http://api.alegra.com/api/v1", true},
		{"http://evil.example/api", true},
		{"http://192.168.1.50:3000", true},
	}
	for _, c := range cases {
		var buf bytes.Buffer
		warnInsecureBaseURL(&buf, c.url)
		gotWarn := strings.Contains(buf.String(), "not HTTPS")
		if gotWarn != c.wantWarn {
			t.Errorf("warnInsecureBaseURL(%q): warned=%v, want %v (output: %q)",
				c.url, gotWarn, c.wantWarn, buf.String())
		}
	}
}
