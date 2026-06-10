package commands

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/jjuanrivvera/alegra-cli/internal/api"
	"github.com/jjuanrivvera/alegra-cli/internal/config"
)

func TestStrOr(t *testing.T) {
	m := map[string]any{"name": "Acme", "nilval": nil}
	assert.Equal(t, "Acme", strOr(m, "name", "?"))
	assert.Equal(t, "?", strOr(m, "missing", "?"))
	assert.Equal(t, "?", strOr(m, "nilval", "?"))
}

func TestOrNone(t *testing.T) {
	assert.Equal(t, "(not set)", orNone(""))
	assert.Equal(t, "you@x.com", orNone("you@x.com"))
}

func TestIsBuiltinCommand(t *testing.T) {
	assert.True(t, isBuiltinCommand("doctor"))
	assert.True(t, isBuiltinCommand("config"))
	assert.False(t, isBuiltinCommand("totally-not-a-command"))
}

func TestReportFlagsValues(t *testing.T) {
	f := &reportFlags{from: "2026-01-01", to: "2026-03-31", start: 5, limit: 10}

	q := f.values(true)
	assert.Equal(t, "2026-01-01", q.Get("from"))
	assert.Equal(t, "2026-03-31", q.Get("to"))
	assert.Equal(t, "5", q.Get("start"))
	assert.Equal(t, "10", q.Get("limit"))

	q = f.values(false)
	assert.Equal(t, "2026-01-01", q.Get("from"))
	assert.Empty(t, q.Get("start"), "non-paginated reports omit pagination")
	assert.Empty(t, q.Get("limit"))
}

func TestInvoiceIDAndEmitKey(t *testing.T) {
	assert.Equal(t, "5", invoiceID(api.Invoice{ID: api.ID("5")}))
	assert.Equal(t, "prod:5", emitKey("prod", "5"))
}

func TestEmitCacheRoundTrip(t *testing.T) {
	dir := t.TempDir()
	t.Setenv(config.EnvConfig, filepath.Join(dir, "config.yaml"))

	// Missing cache is a fresh start, not an error.
	cache, err := loadEmitCache()
	require.NoError(t, err)
	assert.Empty(t, cache)

	require.NoError(t, saveEmitCache(map[string]bool{"prod:5": true}))
	assert.Equal(t, filepath.Join(dir, "emitted.json"), emitCachePath())

	cache, err = loadEmitCache()
	require.NoError(t, err)
	assert.True(t, cache["prod:5"])

	// The atomic write must not leave temp files behind.
	entries, rerr := os.ReadDir(dir)
	require.NoError(t, rerr)
	for _, e := range entries {
		assert.NotContains(t, e.Name(), "emitted-", "leftover temp file %s", e.Name())
	}
}

func TestLoadEmitCache_CorruptFileIsAnError(t *testing.T) {
	dir := t.TempDir()
	t.Setenv(config.EnvConfig, filepath.Join(dir, "config.yaml"))
	require.NoError(t, os.WriteFile(filepath.Join(dir, "emitted.json"), []byte("{not json"), 0o600))

	// A corrupt cache must never read as empty: that would drop the
	// idempotency guard and allow double emission.
	_, err := loadEmitCache()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "corrupt cache")
}

func TestRootCmd(t *testing.T) {
	assert.Same(t, rootCmd, RootCmd())
}

func TestFormatValidationError(t *testing.T) {
	err := formatValidationError("contacts", "colombia", []string{"name is required"})
	require.Error(t, err)
	msg := err.Error()
	assert.Contains(t, msg, "pre-flight validation failed")
	assert.Contains(t, msg, "contacts")
	assert.Contains(t, msg, "colombia")
	assert.Contains(t, msg, "name is required")

	// No country: still renders the resource context.
	assert.Contains(t, formatValidationError("invoices", "", []string{"x"}).Error(), "invoices")
}

func TestVerifyCredentials(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/users/self", r.URL.Path)
		_, _ = w.Write([]byte(`{"name":"Tester"}`))
	}))
	defer srv.Close()
	client := api.New(api.WithBaseURL(srv.URL), api.WithBasicAuth("e@x.com", "tok"), api.WithRequestsPerSecond(1000))
	assert.NoError(t, verifyCredentials(context.Background(), client))
}

func TestPromptSecretNonInteractive(t *testing.T) {
	cmd := &cobra.Command{}
	cmd.SetIn(bytes.NewBufferString("s3cr3t\n"))
	cmd.SetOut(&bytes.Buffer{})
	got, err := promptSecret(cmd, "token: ")
	require.NoError(t, err)
	assert.Equal(t, "s3cr3t", got)
}

func TestCurrentProfileName(t *testing.T) {
	dir := t.TempDir()
	t.Setenv(config.EnvConfig, filepath.Join(dir, "config.yaml"))
	t.Setenv(config.EnvProfile, "")
	// runRoot resets flag globals before each run, not after — a prior test's
	// --profile would otherwise leak into this direct helper call.
	prev := flagProfile
	flagProfile = ""
	t.Cleanup(func() { flagProfile = prev })

	cfg := config.New()
	cfg.DefaultProfile = "prod"
	cfg.SetProfile(&config.Profile{Name: "prod"})
	require.NoError(t, cfg.Save())

	assert.Equal(t, "prod", currentProfileName())
}
