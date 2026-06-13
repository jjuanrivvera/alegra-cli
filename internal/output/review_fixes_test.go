package output

import (
	"bytes"
	"strings"
	"testing"
)

// TestRenderYAML_NumbersAreNotQuoted pins finding M3: YAML output must emit
// numbers as numbers, not quoted strings, matching the JSON renderer.
func TestRenderYAML_NumbersAreNotQuoted(t *testing.T) {
	var buf bytes.Buffer
	data := map[string]any{"total": 1234.5, "count": 7, "name": "Acme"}
	if err := Render(&buf, data, FormatYAML, nil); err != nil {
		t.Fatal(err)
	}
	out := buf.String()
	if !strings.Contains(out, "count: 7") {
		t.Errorf("integer should be unquoted (count: 7); got:\n%s", out)
	}
	if strings.Contains(out, `"7"`) || strings.Contains(out, `'7'`) {
		t.Errorf("integer must not be quoted; got:\n%s", out)
	}
	if !strings.Contains(out, "total: 1234.5") {
		t.Errorf("float should be unquoted (total: 1234.5); got:\n%s", out)
	}
}

// TestRenderCSV_SingleObject pins the theme-4 fix: a single record (e.g. `get`)
// must render as one CSV row, not nothing.
func TestRenderCSV_SingleObject(t *testing.T) {
	var buf bytes.Buffer
	data := map[string]any{"id": "5", "name": "Acme"}
	if err := Render(&buf, data, FormatCSV, []string{"id", "name"}); err != nil {
		t.Fatal(err)
	}
	out := buf.String()
	if !strings.Contains(out, "id,name") || !strings.Contains(out, "5,Acme") {
		t.Errorf("single object should render as one CSV row; got:\n%q", out)
	}
}

// TestCSVField_HardenedCases covers the L14 additions: a formula trigger hiding
// behind leading whitespace/newline, and "-Infinity" (which ParseFloat accepts).
func TestCSVField_HardenedCases(t *testing.T) {
	cases := []struct{ in, want string }{
		{" =1+1", "' =1+1"},           // leading space before a formula
		{"\n=cmd()", "'\n=cmd()"},     // formula behind a newline
		{"  @SUM(A1)", "'  @SUM(A1)"}, // formula behind spaces
		{"-Infinity", "'-Infinity"},   // not a finite number
		{"-Inf", "'-Inf"},             // not a finite number
		{"-42.5", "-42.5"},            // a real negative number stays raw
	}
	for _, c := range cases {
		if got := csvField(c.in); got != c.want {
			t.Errorf("csvField(%q) = %q, want %q", c.in, got, c.want)
		}
	}
}
