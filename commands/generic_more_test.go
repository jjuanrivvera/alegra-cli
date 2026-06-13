package commands

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseMapping(t *testing.T) {
	got, err := parseMapping([]string{"Name=name,NIT=identification.number", "Email=email"})
	require.NoError(t, err)
	assert.Equal(t, map[string]string{
		"Name":  "name",
		"NIT":   "identification.number",
		"Email": "email",
	}, got)

	_, err = parseMapping([]string{"bogus"})
	assert.Error(t, err)
}

func TestInferValue(t *testing.T) {
	assert.Equal(t, "", inferValue(""))
	assert.Equal(t, true, inferValue("true"))
	assert.Equal(t, false, inferValue("false"))
	assert.Nil(t, inferValue("null"))
	// Clean integers stay integers (int64), not float64 — so large values keep
	// full precision (H2).
	assert.Equal(t, int64(42), inferValue("42"))
	assert.Equal(t, int64(-7), inferValue("-7"))
	assert.Equal(t, int64(123456789012345678), inferValue("123456789012345678"))
	// Genuine decimals coerce to float64.
	assert.Equal(t, 1.5, inferValue("1.5"))
	assert.Equal(t, "plain", inferValue("plain"))
	assert.Equal(t, "quoted", inferValue(`"quoted"`))
	assert.Equal(t, []any{1.0, 2.0}, inferValue(`[1,2]`))
	assert.Equal(t, map[string]any{"id": 1.0}, inferValue(`{"id":1}`))

	// H2: identifiers that look numeric but must NOT be silently rewritten stay
	// as strings — leading zeros, a leading "+", underscores, scientific forms
	// that change the digits, trailing-zero decimals, and NaN/Inf.
	for _, s := range []string{
		"007123456",               // leading zeros (NIT/account number)
		"+44",                     // leading plus
		"1_000",                   // underscore separator
		"1.5e2",                   // scientific form (would become 150)
		"1.50",                    // trailing zero would be dropped
		"99999999999999999999999", // beyond int64/float64 exact range
		"NaN", "Inf", "-Inf", "Infinity",
	} {
		assert.Equal(t, s, inferValue(s), "inferValue(%q) must stay a string", s)
	}
}

func TestBodyFlags_Build(t *testing.T) {
	// --set fields, values JSON-inferred.
	raw, err := bodyFlags{sets: []string{"name=Acme", `type=["client"]`}}.build()
	require.NoError(t, err)
	var m map[string]any
	require.NoError(t, json.Unmarshal(raw, &m))
	assert.Equal(t, "Acme", m["name"])
	assert.Equal(t, []any{"client"}, m["type"])

	// --data JSON.
	raw, err = bodyFlags{data: `{"a":1}`}.build()
	require.NoError(t, err)
	assert.JSONEq(t, `{"a":1}`, string(raw))

	// --file.
	f := filepath.Join(t.TempDir(), "body.json")
	require.NoError(t, os.WriteFile(f, []byte(`{"from":"file"}`), 0o600))
	raw, err = bodyFlags{file: f}.build()
	require.NoError(t, err)
	assert.JSONEq(t, `{"from":"file"}`, string(raw))

	// Nothing provided → error.
	_, err = bodyFlags{}.build()
	assert.Error(t, err)

	// Invalid --set → error.
	_, err = bodyFlags{sets: []string{"noequals"}}.build()
	assert.Error(t, err)
}

func TestBodyFlags_BuildOptional_Empty(t *testing.T) {
	raw, err := bodyFlags{}.buildOptional()
	require.NoError(t, err)
	assert.Nil(t, raw)
}

func TestConfirm(t *testing.T) {
	for _, in := range []string{"y\n", "yes\n", "Y\n"} {
		cmd := &cobra.Command{}
		cmd.SetIn(bytes.NewBufferString(in))
		cmd.SetOut(&bytes.Buffer{})
		assert.True(t, confirm(cmd, "proceed"), "input %q", in)
	}
	for _, in := range []string{"n\n", "\n", "whatever\n"} {
		cmd := &cobra.Command{}
		cmd.SetIn(bytes.NewBufferString(in))
		cmd.SetOut(&bytes.Buffer{})
		assert.False(t, confirm(cmd, "proceed"), "input %q", in)
	}
}

// Every registered resource builds its list command at init, so by the time
// tests run any filter dropped over an empty definition or a flag collision is
// already recorded. An entry here means a resource silently lost a filter —
// rename the flag in the resource definition.
func TestNoListFiltersAreSilentlyDropped(t *testing.T) {
	assert.Empty(t, droppedListFilters)
}
