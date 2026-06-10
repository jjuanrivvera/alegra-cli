package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewDefaults(t *testing.T) {
	c := New()
	require.NotNil(t, c.Settings)
	assert.Equal(t, "table", c.Settings.DefaultOutputFormat)
	assert.Equal(t, 5.0, c.Settings.RequestsPerSecond)
	assert.Equal(t, "info", c.Settings.LogLevel)
	assert.NotNil(t, c.Profiles)
}

func TestDefaultPathHonorsEnv(t *testing.T) {
	t.Setenv(EnvConfig, "/tmp/custom/alegra.yaml")
	assert.Equal(t, "/tmp/custom/alegra.yaml", DefaultPath())
}

func TestSaveLoadRoundTrip(t *testing.T) {
	path := filepath.Join(t.TempDir(), "config.yaml")
	t.Setenv(EnvConfig, path)

	c := New()
	c.DefaultProfile = "prod"
	c.SetProfile(&Profile{Name: "prod", Email: "you@biz.com", BaseURL: "https://api.alegra.com/api/v1", Country: "colombia"})
	c.Aliases = map[string]string{"unpaid": "invoices list --status open"}
	require.NoError(t, c.Save())
	assert.Equal(t, path, c.Path())

	got, err := Load()
	require.NoError(t, err)
	assert.Equal(t, "prod", got.DefaultProfile)
	assert.Equal(t, "you@biz.com", got.Profile("prod").Email)
	assert.Equal(t, "colombia", got.Profile("prod").Country)
	assert.Equal(t, "invoices list --status open", got.Aliases["unpaid"])
}

func TestLoadMissingReturnsDefault(t *testing.T) {
	t.Setenv(EnvConfig, filepath.Join(t.TempDir(), "does-not-exist.yaml"))
	c, err := Load()
	require.NoError(t, err)
	assert.NotNil(t, c.Settings)
	assert.Empty(t, c.Profiles)
}

func TestProfileCreatesInMemory(t *testing.T) {
	c := New()
	p := c.Profile("ghost") // not stored
	assert.Equal(t, "ghost", p.Name)
	assert.Empty(t, p.Email)
}

func TestActiveProfileNamePrecedence(t *testing.T) {
	c := New()
	c.DefaultProfile = "configured"

	t.Setenv(EnvProfile, "")
	assert.Equal(t, "flag", c.ActiveProfileName("flag"), "explicit override wins")
	assert.Equal(t, "configured", c.ActiveProfileName(""), "falls back to default profile")

	t.Setenv(EnvProfile, "fromenv")
	assert.Equal(t, "fromenv", c.ActiveProfileName(""), "env beats configured default")
	assert.Equal(t, "flag", c.ActiveProfileName("flag"), "override still beats env")

	empty := New()
	t.Setenv(EnvProfile, "")
	assert.Equal(t, "default", empty.ActiveProfileName(""), "ultimate fallback")
}

func TestResolveDefaultsAndEnvOverrides(t *testing.T) {
	c := New()
	c.SetProfile(&Profile{Name: "p", Email: "file@x.com", Token: "filetok"})

	// Clear all relevant env so we see file/defaults.
	for _, k := range []string{EnvBaseURL, EnvEmail, EnvToken, EnvBearer, EnvOutput, EnvLogLevel} {
		t.Setenv(k, "")
	}
	r := c.Resolve("p")
	assert.Equal(t, "https://api.alegra.com/api/v1", r.BaseURL)
	assert.Equal(t, "file@x.com", r.Email)
	assert.Equal(t, "filetok", r.Token)
	assert.Equal(t, "table", r.OutputFormat)
	assert.Equal(t, 5.0, r.RequestsPerSecond)

	// Env overrides take precedence over the profile/defaults.
	t.Setenv(EnvBaseURL, "https://sandbox.example/api/v1")
	t.Setenv(EnvEmail, "env@x.com")
	t.Setenv(EnvToken, "envtok")
	t.Setenv(EnvOutput, "json")
	r2 := c.Resolve("p")
	assert.Equal(t, "https://sandbox.example/api/v1", r2.BaseURL)
	assert.Equal(t, "env@x.com", r2.Email)
	assert.Equal(t, "envtok", r2.Token)
	assert.Equal(t, "json", r2.OutputFormat)
}

func TestResolveClampsRequestsPerSecond(t *testing.T) {
	c := New()
	c.Settings.RequestsPerSecond = 0 // invalid → fallback
	r := c.Resolve("default")
	assert.Equal(t, 5.0, r.RequestsPerSecond)
}

func TestFirstNonEmpty(t *testing.T) {
	assert.Equal(t, "b", firstNonEmpty("", "b", "c"))
	assert.Equal(t, "a", firstNonEmpty("a", "b"))
	assert.Equal(t, "", firstNonEmpty("", ""))
}

func TestParseRPS(t *testing.T) {
	assert.Equal(t, 5.0, ParseRPS("", 5.0))
	assert.Equal(t, 12.5, ParseRPS("12.5", 5.0))
	assert.Equal(t, 5.0, ParseRPS("nope", 5.0), "invalid falls back")
	assert.Equal(t, 5.0, ParseRPS("-3", 5.0), "non-positive falls back")
}

func TestLoadCorruptYAMLIsAnError(t *testing.T) {
	path := filepath.Join(t.TempDir(), "config.yaml")
	t.Setenv(EnvConfig, path)
	require.NoError(t, os.WriteFile(path, []byte("profiles: [not: a: map\n\t"), 0o600))

	// A corrupt config must surface a parse error, never be silently
	// replaced by defaults (that would drop every profile).
	_, err := Load()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "parsing config")
}

func TestSaveIsAtomicAndLeavesNoTempFiles(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.yaml")
	t.Setenv(EnvConfig, path)

	c := New()
	c.SetProfile(&Profile{Name: "p", Email: "a@b.c"})
	require.NoError(t, c.Save())
	require.NoError(t, c.Save()) // overwrite path: rename over existing file

	entries, err := os.ReadDir(dir)
	require.NoError(t, err)
	require.Len(t, entries, 1, "temp files must not survive a save")
	assert.Equal(t, "config.yaml", entries[0].Name())

	got, err := Load()
	require.NoError(t, err)
	assert.Equal(t, "a@b.c", got.Profile("p").Email)
}
