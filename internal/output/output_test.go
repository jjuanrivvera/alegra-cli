package output

import (
	"bytes"
	"strings"
	"testing"
)

func TestCSVField(t *testing.T) {
	cases := []struct {
		in, want string
	}{
		{"", ""},
		{"plain text", "plain text"},
		{"Acme Corp", "Acme Corp"},
		{"-42.5", "-42.5"},        // real negative number: untouched
		{"-1", "-1"},              // real negative number: untouched
		{"=1+1", "'=1+1"},         // formula
		{"+1+1", "'+1+1"},         // formula
		{"@SUM(A1)", "'@SUM(A1)"}, // formula
		{"=cmd|'/c calc'!A1", "'=cmd|'/c calc'!A1"},
		{"-2+3+cmd|'/c calc'!A0", "'-2+3+cmd|'/c calc'!A0"}, // leading '-' but not numeric
		{"\t=1", "'\t=1"}, // leading tab
		{"\r=1", "'\r=1"}, // leading CR
	}
	for _, c := range cases {
		if got := csvField(c.in); got != c.want {
			t.Errorf("csvField(%q) = %q, want %q", c.in, got, c.want)
		}
	}
}

// TestRenderCSVNeutralizesFormulas guards the export path end-to-end: a contact
// name set by a third party must not become a live spreadsheet formula.
func TestRenderCSVNeutralizesFormulas(t *testing.T) {
	rows := []map[string]any{
		{"id": "1", "name": `=cmd|'/c calc.exe'!A1`},
		{"id": "2", "name": `@SUM(1+9)*cmd`},
		{"id": "3", "name": `-9+cmd|'/c calc'!A0`},
		{"id": "4", "name": "Legitimate Co"},
	}
	var buf bytes.Buffer
	if err := Render(&buf, rows, FormatCSV, []string{"id", "name"}); err != nil {
		t.Fatalf("render: %v", err)
	}
	out := buf.String()

	for _, bad := range []string{",=cmd", ",@SUM", ",-9+cmd"} {
		if strings.Contains(out, bad) {
			t.Errorf("formula not neutralized; CSV contains %q:\n%s", bad, out)
		}
	}
	if !strings.Contains(out, ",Legitimate Co") {
		t.Errorf("benign value should pass through untouched:\n%s", out)
	}
}
