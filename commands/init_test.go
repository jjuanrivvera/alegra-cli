package commands

import (
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zalando/go-keyring"

	"github.com/jjuanrivvera/alegra-cli/internal/config"
)

func TestIntegration_Init(t *testing.T) {
	keyring.MockInit()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.URL.Path {
		case "/users/self":
			_, _ = w.Write([]byte(`{"name":"Tester"}`))
		case "/company":
			_, _ = w.Write([]byte(`{"applicationVersion":"colombia"}`))
		default:
			_, _ = w.Write([]byte(`{}`))
		}
	}))
	t.Cleanup(srv.Close)
	t.Setenv(config.EnvBaseURL, srv.URL)
	t.Setenv(config.EnvEmail, "")
	t.Setenv(config.EnvToken, "")
	t.Setenv(config.EnvProfile, "")
	t.Setenv(config.EnvConfig, filepath.Join(t.TempDir(), "config.yaml"))

	out, err := runRoot(t, "init", "--email", "new@x.com", "--token", "tok", "--profile", "onboard")
	require.NoError(t, err)
	assert.Contains(t, out, "Setup complete")
	assert.Contains(t, out, "colombia")
	assert.Contains(t, out, "Next steps")
}

func TestIntegration_VersionJSON(t *testing.T) {
	out, err := runRoot(t, "version", "--json")
	require.NoError(t, err)
	assert.Contains(t, out, `"version"`)
	assert.Contains(t, out, `"goVersion"`)
}
