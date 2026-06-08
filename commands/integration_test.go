package commands

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/jjuanrivvera/alegra-cli/internal/config"
)

// runRoot executes the real command tree with args against a clean global state,
// returning combined stdout+stderr. It resets the once-cached resolution and the
// persistent flag globals so each call re-reads the test environment.
func runRoot(t *testing.T, args ...string) (string, error) {
	t.Helper()
	cfgOnce = sync.Once{}
	cfg, resolved, cfgErr = nil, nil, nil
	flagProfile, flagOutput, flagBaseURL, flagColumns = "", "", "", nil
	flagDryRun, flagShowToken, flagVerbose, flagRPS = false, false, false, 0

	var buf bytes.Buffer
	rootCmd.SetOut(&buf)
	rootCmd.SetErr(&buf)
	rootCmd.SetArgs(args)
	err := rootCmd.Execute()
	return buf.String(), err
}

// integrationServer returns a handler covering the endpoints the command tree
// hits across CRUD + doctor.
func integrationServer(t *testing.T) *httptest.Server {
	t.Helper()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		write := func(s string) { _, _ = w.Write([]byte(s)) }
		switch {
		case r.Method == http.MethodGet && r.URL.Path == "/contacts" && r.URL.Query().Get("metadata") == "true":
			write(`{"metadata":{"total":42},"data":[{"id":"1","name":"Acme"}]}`)
		case r.Method == http.MethodGet && r.URL.Path == "/contacts":
			write(`[{"id":"1","name":"Acme"}]`)
		case r.Method == http.MethodGet && r.URL.Path == "/contacts/1":
			write(`{"id":"1","name":"Acme"}`)
		case r.Method == http.MethodPost && r.URL.Path == "/contacts":
			write(`{"id":"2","name":"Beta"}`)
		case r.Method == http.MethodPut && r.URL.Path == "/contacts/1":
			write(`{"id":"1","name":"NewName"}`)
		case r.Method == http.MethodDelete && r.URL.Path == "/contacts/2":
			write(`{}`)
		case r.Method == http.MethodPost && r.URL.Path == "/invoices/5/void":
			write(`{"id":"5","status":"void"}`)
		case r.Method == http.MethodGet && r.URL.Path == "/users/self":
			write(`{"name":"Tester","role":"admin"}`)
		case r.Method == http.MethodGet && r.URL.Path == "/company":
			write(`{"name":"Acme","applicationVersion":"colombia","regime":"Común"}`)
		case r.Method == http.MethodGet && r.URL.Path == "/number-templates":
			write(`[{"id":"7","name":"FE"}]`)
		case r.Method == http.MethodGet && r.URL.Path == "/reconciliations":
			write(`[]`)
		case r.Method == http.MethodGet && r.URL.Path == "/reports/sales-by-client":
			write(`{"data":[{"clientName":"Acme","total":1000}]}`)
		case r.Method == http.MethodGet && r.URL.Path == "/reports/sales-by-client-totals":
			write(`{"total":1000}`)
		default:
			w.WriteHeader(http.StatusNotFound)
			write(`{"message":"not found"}`)
		}
	}))
	t.Cleanup(srv.Close)
	return srv
}

func TestIntegration_CommandTree(t *testing.T) {
	srv := integrationServer(t)
	t.Setenv(config.EnvBaseURL, srv.URL)
	t.Setenv(config.EnvEmail, "e@x.com")
	t.Setenv(config.EnvToken, "tok")
	t.Setenv(config.EnvProfile, "")
	t.Setenv(config.EnvConfig, filepath.Join(t.TempDir(), "config.yaml"))

	t.Run("version", func(t *testing.T) {
		out, err := runRoot(t, "version")
		require.NoError(t, err)
		assert.Contains(t, out, "alegra")
	})

	t.Run("list", func(t *testing.T) {
		out, err := runRoot(t, "contacts", "list", "-o", "json")
		require.NoError(t, err)
		assert.Contains(t, out, "Acme")
	})

	t.Run("get", func(t *testing.T) {
		out, err := runRoot(t, "contacts", "get", "1", "-o", "json")
		require.NoError(t, err)
		assert.Contains(t, out, "Acme")
	})

	t.Run("get not found errors", func(t *testing.T) {
		_, err := runRoot(t, "contacts", "get", "999", "-o", "json")
		assert.Error(t, err)
	})

	t.Run("doctor", func(t *testing.T) {
		out, err := runRoot(t, "doctor")
		require.NoError(t, err)
		assert.Contains(t, out, "colombia")
	})

	// create uses --set (append-style flag) — run once so it can't accumulate.
	t.Run("create", func(t *testing.T) {
		out, err := runRoot(t, "contacts", "create", "--set", "name=Beta", "-o", "json")
		require.NoError(t, err)
		assert.Contains(t, out, "Beta")
	})

	t.Run("update", func(t *testing.T) {
		out, err := runRoot(t, "contacts", "update", "1", "--set", "name=NewName", "-o", "json")
		require.NoError(t, err)
		assert.Contains(t, out, "NewName")
	})

	t.Run("action void", func(t *testing.T) {
		out, err := runRoot(t, "invoices", "void", "5", "-o", "json")
		require.NoError(t, err)
		assert.Contains(t, out, "void")
	})

	t.Run("auth status", func(t *testing.T) {
		out, err := runRoot(t, "auth", "status")
		require.NoError(t, err)
		assert.Contains(t, out, "authenticated")
	})

	t.Run("config view", func(t *testing.T) {
		_, err := runRoot(t, "config", "view")
		require.NoError(t, err)
	})

	t.Run("delete", func(t *testing.T) {
		_, err := runRoot(t, "contacts", "delete", "2", "-y")
		require.NoError(t, err)
	})

	t.Run("report (data envelope)", func(t *testing.T) {
		out, err := runRoot(t, "reports", "sales-by-client", "--from", "2026-01-01", "--to", "2026-03-31", "-o", "json")
		require.NoError(t, err)
		assert.Contains(t, out, "Acme")
	})

	t.Run("report (raw)", func(t *testing.T) {
		out, err := runRoot(t, "reports", "sales-by-client-totals", "--from", "2026-01-01", "--to", "2026-03-31", "-o", "json")
		require.NoError(t, err)
		assert.Contains(t, out, "1000")
	})

	// --count uses a bool flag on the list command; run it last so the toggled
	// state can't leak into a plain list.
	t.Run("list count", func(t *testing.T) {
		out, err := runRoot(t, "contacts", "list", "--count")
		require.NoError(t, err)
		assert.Contains(t, out, "42")
	})
}
