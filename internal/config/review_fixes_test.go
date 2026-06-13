package config

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestSave_FilePermissions pins T6: the config file (which may hold a token) is
// written 0600 and its directory 0700, so a regression that widened them — and
// exposed a stored credential — would fail CI.
func TestSave_FilePermissions(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("unix file modes are not meaningful on Windows")
	}
	path := filepath.Join(t.TempDir(), "cfgdir", "config.yaml")
	t.Setenv(EnvConfig, path)

	c := New() // DefaultPath() resolves to EnvConfig
	c.SetProfile(&Profile{Name: "default", Email: "e@x.com", Token: "shh"})
	require.NoError(t, c.Save())

	info, err := os.Stat(path)
	require.NoError(t, err)
	assert.Equal(t, os.FileMode(0o600), info.Mode().Perm(), "config file must be 0600")

	dirInfo, err := os.Stat(filepath.Dir(path))
	require.NoError(t, err)
	assert.Equal(t, os.FileMode(0o700), dirInfo.Mode().Perm(), "config dir must be 0700")
}
