package commands

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/jjuanrivvera/alegra-cli/internal/config"
)

func TestCatalog(t *testing.T) {
	// Isolate from any real config so --country is the only country source.
	t.Setenv(config.EnvProfile, "")
	t.Setenv(config.EnvConfig, filepath.Join(t.TempDir(), "config.yaml"))

	out, err := runRoot(t, "catalog", "units", "--country", "colombia")
	require.NoError(t, err)
	require.Contains(t, out, "Unidad")

	out, err = runRoot(t, "catalog", "--country", "mexico", "-o", "json")
	require.NoError(t, err)
	require.Contains(t, out, "unidades-de-medida")

	// `reference` alias + friendly category alias both resolve.
	out, err = runRoot(t, "reference", "identification-types", "--country", "colombia")
	require.NoError(t, err)
	require.Contains(t, out, "Cédula de ciudadanía")

	// Unknown category errors helpfully.
	_, err = runRoot(t, "catalog", "nope", "--country", "colombia")
	require.Error(t, err)

	// No country and no configured account -> clear error.
	_, err = runRoot(t, "catalog", "units")
	require.Error(t, err)
}
