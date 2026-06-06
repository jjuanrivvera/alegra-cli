package commands

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSetDotPath(t *testing.T) {
	m := map[string]any{}
	setDotPath(m, "name", "Acme")
	setDotPath(m, "identification.type", "NIT")
	setDotPath(m, "identification.number", "901123456")
	assert.Equal(t, "Acme", m["name"])
	id, ok := m["identification"].(map[string]any)
	require.True(t, ok)
	assert.Equal(t, "NIT", id["type"])
	assert.Equal(t, "901123456", id["number"])
}

func TestParseKeyVals(t *testing.T) {
	got, err := parseKeyVals([]string{"Name=name", "NIT=identification.number"})
	require.NoError(t, err)
	assert.Equal(t, map[string]string{"Name": "name", "NIT": "identification.number"}, got)

	_, err = parseKeyVals([]string{"bad"})
	assert.Error(t, err)
}
