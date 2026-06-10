package commands

import (
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zalando/go-keyring"

	"github.com/jjuanrivvera/alegra-cli/internal/catalog"
	"github.com/jjuanrivvera/alegra-cli/internal/config"
)

const satTestSQL = `INSERT INTO cfdi_40_productos_servicios VALUES('52141506','Refrigeradores','','','','2022-01-01','',1,'Neveras, Frigoríficos');
INSERT INTO cfdi_40_productos_servicios VALUES('81111504','Programación de software de aplicaciones','','','','2022-01-01','',1,'');
`

// satTestSource points the SAT sync at a local fixture server for the test's
// lifetime and returns the temp config home it set up.
func satTestSource(t *testing.T) {
	t.Helper()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.URL.Path, "version.txt") {
			_, _ = w.Write([]byte("10.7.test"))
			return
		}
		_, _ = w.Write([]byte(satTestSQL))
	}))
	t.Cleanup(srv.Close)
	origData, origVersion := catalog.SATDataURL, catalog.SATVersionURL
	catalog.SATDataURL, catalog.SATVersionURL = srv.URL+"/data.sql", srv.URL+"/version.txt"
	t.Cleanup(func() { catalog.SATDataURL, catalog.SATVersionURL = origData, origVersion })
}

func TestCatalogSAT_SyncAndSearch(t *testing.T) {
	satTestSource(t)
	t.Setenv(config.EnvConfig, filepath.Join(t.TempDir(), "config.yaml"))
	t.Setenv(config.EnvProfile, "")

	// Searching before syncing must point at the fix.
	_, err := runRoot(t, "catalog", "product-keys", "refrigerador")
	require.ErrorContains(t, err, "sync-sat")

	out, err := runRoot(t, "catalog", "sync-sat")
	require.NoError(t, err)
	assert.Contains(t, out, "Synced 2 product keys")
	assert.Contains(t, out, "10.7.test")

	// Accent- and case-insensitive search across name and similar names.
	out, err = runRoot(t, "catalog", "product-keys", "refrigerador", "-o", "json")
	require.NoError(t, err)
	assert.Contains(t, out, "52141506")

	out, err = runRoot(t, "catalog", "product-keys", "neveras")
	require.NoError(t, err)
	assert.Contains(t, out, "52141506")

	out, err = runRoot(t, "catalog", "product-keys", "programacion")
	require.NoError(t, err)
	assert.Contains(t, out, "81111504")

	_, err = runRoot(t, "catalog", "product-keys", "zzz-no-match")
	assert.ErrorContains(t, err, "no product keys match")
}

// mxServer is an API stub for a Mexican account (init + doctor flows).
func mxServer(t *testing.T) *httptest.Server {
	t.Helper()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.URL.Path {
		case "/users/self":
			_, _ = w.Write([]byte(`{"name":"Tester","role":"admin"}`))
		case "/company":
			_, _ = w.Write([]byte(`{"name":"Tacos SA","applicationVersion":"mexico"}`))
		case "/number-templates":
			_, _ = w.Write([]byte(`[]`))
		default:
			_, _ = w.Write([]byte(`{}`))
		}
	}))
	t.Cleanup(srv.Close)
	return srv
}

func setupMXEnv(t *testing.T) {
	t.Helper()
	keyring.MockInit()
	srv := mxServer(t)
	t.Setenv(config.EnvBaseURL, srv.URL)
	t.Setenv(config.EnvEmail, "")
	t.Setenv(config.EnvToken, "")
	t.Setenv(config.EnvProfile, "")
	t.Setenv(config.EnvConfig, filepath.Join(t.TempDir(), "config.yaml"))
}

// runRootStdin is runRoot with a scripted stdin (init's prompts read from it).
func runRootStdin(t *testing.T, stdin string, args ...string) (string, error) {
	t.Helper()
	rootCmd.SetIn(strings.NewReader(stdin))
	t.Cleanup(func() { rootCmd.SetIn(nil) })
	return runRoot(t, args...)
}

func TestInit_MexicoOffersSATCatalog(t *testing.T) {
	satTestSource(t)

	t.Run("accept downloads the catalog", func(t *testing.T) {
		setupMXEnv(t)
		out, err := runRootStdin(t, "y\n", "init", "--email", "mx@x.com", "--token", "tok")
		require.NoError(t, err)
		assert.Contains(t, out, "Mexican account detected")
		assert.Contains(t, out, "ok (2 keys, catalog 10.7.test)")

		out, err = runRoot(t, "catalog", "product-keys", "refrigerador")
		require.NoError(t, err)
		assert.Contains(t, out, "52141506")
	})

	t.Run("decline leaves a hint and succeeds", func(t *testing.T) {
		setupMXEnv(t)
		out, err := runRootStdin(t, "n\n", "init", "--email", "mx@x.com", "--token", "tok")
		require.NoError(t, err)
		assert.Contains(t, out, "sync-sat` anytime")
		assert.Contains(t, out, "Next steps")
	})

	t.Run("EOF (non-interactive) skips non-fatally", func(t *testing.T) {
		setupMXEnv(t)
		out, err := runRootStdin(t, "", "init", "--email", "mx@x.com", "--token", "tok")
		require.NoError(t, err)
		assert.Contains(t, out, "sync-sat` anytime")
	})

	t.Run("download failure never fails init", func(t *testing.T) {
		setupMXEnv(t)
		origData := catalog.SATDataURL
		catalog.SATDataURL = "http://127.0.0.1:1/unreachable"
		t.Cleanup(func() { catalog.SATDataURL = origData })
		out, err := runRootStdin(t, "y\n", "init", "--email", "mx@x.com", "--token", "tok")
		require.NoError(t, err)
		assert.Contains(t, out, "failed")
		assert.Contains(t, out, "sync-sat` anytime")
	})
}

func TestInit_NonMexicoDoesNotOffer(t *testing.T) {
	satTestSource(t)
	keyring.MockInit()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if r.URL.Path == "/company" {
			_, _ = w.Write([]byte(`{"applicationVersion":"colombia"}`))
			return
		}
		_, _ = w.Write([]byte(`{"name":"Tester"}`))
	}))
	t.Cleanup(srv.Close)
	t.Setenv(config.EnvBaseURL, srv.URL)
	t.Setenv(config.EnvEmail, "")
	t.Setenv(config.EnvToken, "")
	t.Setenv(config.EnvProfile, "")
	t.Setenv(config.EnvConfig, filepath.Join(t.TempDir(), "config.yaml"))

	out, err := runRoot(t, "init", "--email", "co@x.com", "--token", "tok")
	require.NoError(t, err)
	assert.NotContains(t, out, "SAT")
}

func TestDoctor_MexicoReportsSATCatalog(t *testing.T) {
	satTestSource(t)
	setupMXEnv(t)
	t.Setenv(config.EnvEmail, "mx@x.com")
	t.Setenv(config.EnvToken, "tok")

	out, err := runRoot(t, "doctor")
	require.NoError(t, err)
	assert.Contains(t, out, "not synced — run `alegra catalog sync-sat`")

	_, err = runRoot(t, "catalog", "sync-sat")
	require.NoError(t, err)

	out, err = runRoot(t, "doctor")
	require.NoError(t, err)
	assert.Contains(t, out, "2 product keys (version 10.7.test")
}
