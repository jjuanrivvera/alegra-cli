package output

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFormatValid(t *testing.T) {
	for _, f := range []Format{FormatTable, FormatJSON, FormatYAML, FormatCSV} {
		assert.True(t, f.Valid(), "%s should be valid", f)
	}
	assert.False(t, Format("xml").Valid())
}

func TestRenderJSON(t *testing.T) {
	var buf bytes.Buffer
	require.NoError(t, Render(&buf, map[string]any{"name": "Acme"}, FormatJSON, nil))
	assert.Contains(t, buf.String(), `"name": "Acme"`)
}

func TestRenderYAML(t *testing.T) {
	var buf bytes.Buffer
	require.NoError(t, Render(&buf, map[string]any{"name": "Acme"}, FormatYAML, nil))
	assert.Contains(t, buf.String(), "name: Acme")
}

func TestRenderCSV(t *testing.T) {
	var buf bytes.Buffer
	rows := []map[string]any{{"id": "1", "name": "Acme"}, {"id": "2", "name": "Beta"}}
	require.NoError(t, Render(&buf, rows, FormatCSV, []string{"id", "name"}))
	out := buf.String()
	assert.Contains(t, out, "id,name")
	assert.Contains(t, out, "1,Acme")
	assert.Contains(t, out, "2,Beta")
}

func TestRenderTable_MultipleRows(t *testing.T) {
	var buf bytes.Buffer
	rows := []map[string]any{
		{"id": "1", "name": "Acme"},
		{"id": "2", "name": "Beta"},
	}
	require.NoError(t, Render(&buf, rows, FormatTable, []string{"id", "name"}))
	out := buf.String()
	assert.Contains(t, out, "ID")
	assert.Contains(t, out, "NAME")
	assert.Contains(t, out, "Acme")
	assert.Contains(t, out, "Beta")
}

func TestRenderTable_TruncatesLongValues(t *testing.T) {
	var buf bytes.Buffer
	long := strings.Repeat("x", 100)
	require.NoError(t, Render(&buf, []map[string]any{{"a": "1", "name": long}, {"a": "2", "name": "y"}}, FormatTable, []string{"a", "name"}))
	assert.Contains(t, buf.String(), "…", "long cell should be truncated with an ellipsis")
}

func TestRenderTable_SingleObjectKeyValue(t *testing.T) {
	var buf bytes.Buffer
	require.NoError(t, Render(&buf, map[string]any{"name": "Acme", "id": "1"}, FormatTable, nil))
	out := buf.String()
	assert.Contains(t, out, "NAME")
	assert.Contains(t, out, "Acme")
}

func TestRenderTable_Empty(t *testing.T) {
	var buf bytes.Buffer
	require.NoError(t, Render(&buf, []any{}, FormatTable, nil))
	assert.Contains(t, buf.String(), "(no results)")
}

func TestRenderUnsupportedFormat(t *testing.T) {
	var buf bytes.Buffer
	err := Render(&buf, map[string]any{"a": 1}, Format("xml"), nil)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "unsupported output format")
}

func TestResolveColumns(t *testing.T) {
	rows := []map[string]any{{"id": "1", "name": "A", "nested": map[string]any{"x": 1}}}
	// Explicit columns pass through unchanged.
	assert.Equal(t, []string{"name", "id"}, resolveColumns(rows, []string{"name", "id"}))

	// Derived columns: scalars only, preferred order first; nested map excluded.
	got := resolveColumns(rows, nil)
	assert.Equal(t, []string{"id", "name"}, got)
	assert.NotContains(t, got, "nested")
}

func TestResolveColumns_CapsAtTen(t *testing.T) {
	row := map[string]any{}
	for _, k := range []string{"c01", "c02", "c03", "c04", "c05", "c06", "c07", "c08", "c09", "c10", "c11", "c12"} {
		row[k] = "v"
	}
	got := resolveColumns([]map[string]any{row}, nil)
	assert.Len(t, got, 10)
}

func TestScalarStringAndValueString(t *testing.T) {
	assert.Equal(t, "", scalarString(nil))
	assert.Equal(t, "true", scalarString(true))
	assert.Equal(t, "abc", scalarString("abc"))
	assert.Equal(t, "3.5", scalarString(3.5))
	assert.Equal(t, "42", scalarString(json.Number("42")))

	// Nested structures compact to JSON via valueString.
	assert.Equal(t, `{"x":1}`, valueString(map[string]any{"x": float64(1)}))
	assert.Equal(t, "plain", valueString("plain"))
}
