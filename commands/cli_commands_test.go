package commands

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/jjuanrivvera/alegra-cli/internal/config"
)

func TestConfigCommands(t *testing.T) {
	t.Setenv(config.EnvConfig, filepath.Join(t.TempDir(), "config.yaml"))
	t.Setenv(config.EnvProfile, "")

	_, err := runRoot(t, "config", "path")
	require.NoError(t, err)

	out, err := runRoot(t, "config", "set-profile", "--name", "prod", "--email", "you@biz.com", "--base-url", "https://api.alegra.com/api/v1")
	require.NoError(t, err)
	assert.Contains(t, out, "prod")

	out, err = runRoot(t, "config", "use", "prod")
	require.NoError(t, err)
	assert.Contains(t, out, "prod")

	out, err = runRoot(t, "config", "list-profiles")
	require.NoError(t, err)
	assert.Contains(t, out, "prod")

	out, err = runRoot(t, "config", "view")
	require.NoError(t, err)
	assert.Contains(t, out, "you@biz.com")
}

func TestAliasCommands(t *testing.T) {
	t.Setenv(config.EnvConfig, filepath.Join(t.TempDir(), "config.yaml"))
	t.Setenv(config.EnvProfile, "")

	out, err := runRoot(t, "alias", "set", "unpaid", "invoices list --status open")
	require.NoError(t, err)
	assert.Contains(t, out, "unpaid")

	out, err = runRoot(t, "alias", "list")
	require.NoError(t, err)
	assert.Contains(t, out, "unpaid")

	out, err = runRoot(t, "alias", "remove", "unpaid")
	require.NoError(t, err)
	assert.Contains(t, out, "unpaid")
}

func TestSkillsCommands(t *testing.T) {
	dir := t.TempDir()

	out, err := runRoot(t, "skills", "install", "--dir", dir)
	require.NoError(t, err)
	assert.Contains(t, out, "Installed")
	_, statErr := os.Stat(filepath.Join(dir, "alegra-cli", "SKILL.md"))
	require.NoError(t, statErr)

	out, err = runRoot(t, "skills", "path", "--agent", "cursor")
	require.NoError(t, err)
	assert.Contains(t, out, "alegra-cli")

	out, err = runRoot(t, "skills", "print")
	require.NoError(t, err)
	assert.Contains(t, out, "name: alegra-cli")

	// dry-run lists files without writing.
	out, err = runRoot(t, "skills", "install", "--dir", filepath.Join(dir, "dry"), "--dry-run")
	require.NoError(t, err)
	assert.Contains(t, out, "Would install")
}
