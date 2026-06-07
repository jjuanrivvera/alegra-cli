package commands

import (
	"context"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/jjuanrivvera/alegra-cli/internal/api"
	"github.com/jjuanrivvera/alegra-cli/internal/config"
)

// detectCountry reads applicationVersion from /company and lowercases it.
func TestDetectCountry(t *testing.T) {
	client := api.New(
		api.WithBaseURL(newCompanyServer(t, `{"name":"Acme","applicationVersion":"costaRica"}`)),
		api.WithBasicAuth("test@example.com", "token"),
		api.WithRequestsPerSecond(1000),
	)
	assert.Equal(t, "costarica", detectCountry(context.Background(), client))
}

func TestDetectCountry_NoField(t *testing.T) {
	client := api.New(
		api.WithBaseURL(newCompanyServer(t, `{"name":"Acme"}`)),
		api.WithBasicAuth("test@example.com", "token"),
		api.WithRequestsPerSecond(1000),
	)
	assert.Equal(t, "", detectCountry(context.Background(), client))
}

func newCompanyServer(t *testing.T, body string) string {
	t.Helper()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(body))
	}))
	t.Cleanup(srv.Close)
	return srv.URL
}

// resolveCountry precedence: --country flag > detected per-profile country >
// global offline hint.
func TestResolveCountryPrecedence(t *testing.T) {
	t.Setenv(config.EnvConfig, filepath.Join(t.TempDir(), "config.yaml"))

	save := func(profileCountry, hint string) {
		cfg := config.New()
		cfg.DefaultProfile = "default"
		cfg.SetProfile(&config.Profile{Name: "default", Email: "x@y.z", Country: profileCountry})
		cfg.Settings.Country = hint
		require.NoError(t, cfg.Save())
	}

	save("mexico", "colombia")
	assert.Equal(t, "peru", resolveCountry("Peru"), "explicit flag wins (and is lowercased)")
	assert.Equal(t, "mexico", resolveCountry(""), "detected profile country wins over the hint")

	save("", "colombia")
	assert.Equal(t, "colombia", resolveCountry(""), "falls back to the offline hint when nothing detected")

	save("", "")
	assert.Equal(t, "", resolveCountry(""), "no country anywhere")
}
