package catalog

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAvailable(t *testing.T) {
	got := Available()
	assert.Contains(t, got, "colombia")
	assert.Contains(t, got, "mexico")
	assert.GreaterOrEqual(t, len(got), 4)
}

func TestNormalize(t *testing.T) {
	cases := map[string]string{
		"CO":         "colombia",
		"colombia":   "colombia",
		"México":     "mexico",
		"mexico":     "mexico",
		"costaRica":  "costarica",
		"Costa Rica": "costarica",
		"pe":         "peru",
		"":           "",
	}
	for in, want := range cases {
		assert.Equalf(t, want, Normalize(in), "Normalize(%q)", in)
	}
}

func TestLoadAndCategory(t *testing.T) {
	cat := mustLoad(t, "colombia")
	assert.Equal(t, "Colombia", cat.Label)
	assert.NotEmpty(t, cat.Categories)

	// Friendly alias resolves to the slugified Spanish key.
	units, ok := cat.Category("units")
	require.True(t, ok, "units alias should resolve")
	assert.Equal(t, "unidades-de-medida", units.Key)
	assert.NotEmpty(t, units.Entries)

	ids, ok := cat.Category("identification-types")
	require.True(t, ok)
	assert.NotEmpty(t, ids.Entries)

	_, ok = cat.Category("does-not-exist")
	assert.False(t, ok)

	assert.Contains(t, cat.CategoryKeys(), "unidades-de-medida")
}

func TestLoadErrors(t *testing.T) {
	_, err := Load("")
	assert.Error(t, err, "empty country should error")

	_, err = Load("narnia")
	assert.Error(t, err, "unknown country should error")
}

func mustLoad(t *testing.T, country string) *Catalog {
	t.Helper()
	c, err := Load(country)
	require.NoError(t, err)
	return c
}
