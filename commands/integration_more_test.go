package commands

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zalando/go-keyring"

	"github.com/jjuanrivvera/alegra-cli/internal/config"
)

// broadServer answers every endpoint the extended command-tree test drives.
func broadServer(t *testing.T) *httptest.Server {
	t.Helper()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		write := func(s string) { _, _ = w.Write([]byte(s)) }
		path := r.URL.Path
		switch {
		case r.Method == http.MethodGet && path == "/contacts" && r.URL.Query().Get("metadata") == "true":
			write(`{"metadata":{"total":7},"data":[{"id":"1","name":"Acme"}]}`)
		case r.Method == http.MethodGet && path == "/contacts":
			write(`[{"id":"1","name":"Acme"}]`)
		case r.Method == http.MethodPost && path == "/contacts":
			write(`{"id":"9","name":"Imported"}`)
		case r.Method == http.MethodGet && path == "/invoices":
			write(`[{"id":"7","numberTemplate":{"id":"1"}}]`)
		case r.Method == http.MethodPost && path == "/invoices/stamp":
			write(`{"ok":true}`)
		case r.Method == http.MethodGet && path == "/company":
			write(`{"name":"Acme","applicationVersion":"colombia"}`)
		case r.Method == http.MethodPut && path == "/company":
			write(`{"name":"NewCo"}`)
		case r.Method == http.MethodGet && path == "/users/self":
			write(`{"name":"Tester","role":"admin"}`)
		case r.Method == http.MethodGet && path == "/journals/entries/graph":
			write(`{"balance":100}`)
		case r.Method == http.MethodGet && path == "/categories/settings":
			write(`{"setting":"x"}`)
		case r.Method == http.MethodPut && path == "/categories/settings":
			write(`{"setting":"y"}`)
		case r.Method == http.MethodGet && path == "/reports/sales-by-seller":
			write(`{"data":[{"sellerName":"Sam","total":500}]}`)
		case r.Method == http.MethodGet && path == "/reports/income-statement":
			write(`{"total":1000}`)
		case r.Method == http.MethodGet && path == "/reports/account-statement":
			write(`{"data":[]}`)
		default:
			w.WriteHeader(http.StatusNotFound)
			write(`{"message":"not found"}`)
		}
	}))
	t.Cleanup(srv.Close)
	return srv
}

func TestIntegration_ExtendedCommands(t *testing.T) {
	srv := broadServer(t)
	t.Setenv(config.EnvBaseURL, srv.URL)
	t.Setenv(config.EnvEmail, "e@x.com")
	t.Setenv(config.EnvToken, "tok")
	t.Setenv(config.EnvProfile, "")
	t.Setenv(config.EnvConfig, filepath.Join(t.TempDir(), "config.yaml"))

	t.Run("list variants", func(t *testing.T) {
		for _, args := range [][]string{
			{"contacts", "list", "--all", "-o", "json"},
			{"contacts", "list", "--query", "acme", "-o", "json"},
			{"contacts", "list", "--param", "status=open", "-o", "json"},
			{"contacts", "list", "--since", "this-month", "--until", "today", "-o", "json"},
			{"contacts", "list", "--order-field", "name", "--order-direction", "ASC", "-o", "yaml"},
		} {
			out, err := runRoot(t, args...)
			require.NoError(t, err, "args %v", args)
			assert.Contains(t, out, "Acme", "args %v", args)
		}
	})

	t.Run("export csv to file", func(t *testing.T) {
		outFile := filepath.Join(t.TempDir(), "contacts.csv")
		_, err := runRoot(t, "contacts", "export", "--out", outFile)
		require.NoError(t, err)
		data, rerr := os.ReadFile(outFile)
		require.NoError(t, rerr)
		assert.Contains(t, string(data), "Acme")
	})

	t.Run("export json stdout", func(t *testing.T) {
		out, err := runRoot(t, "contacts", "export", "--format", "json")
		require.NoError(t, err)
		assert.Contains(t, out, "Acme")
	})

	t.Run("import csv", func(t *testing.T) {
		csvFile := filepath.Join(t.TempDir(), "rows.csv")
		require.NoError(t, os.WriteFile(csvFile, []byte("name\nBeta\nGamma\n"), 0o600))
		out, err := runRoot(t, "contacts", "import", "-f", csvFile)
		require.NoError(t, err)
		// Assert the exact tally, not just a stray "2": "Imported 0, failed 2"
		// would also contain "2" (T5).
		assert.Contains(t, out, "Imported 2, failed 0")
	})

	t.Run("emit explicit id", func(t *testing.T) {
		out, err := runRoot(t, "invoices", "emit", "7")
		require.NoError(t, err)
		assert.Contains(t, out, "stamped")
	})

	t.Run("emit all", func(t *testing.T) {
		out, err := runRoot(t, "invoices", "emit", "--all", "--force")
		require.NoError(t, err)
		assert.Contains(t, out, "Emitted")
	})

	t.Run("emit dry-run", func(t *testing.T) {
		out, err := runRoot(t, "invoices", "emit", "7", "--force", "--dry-run")
		require.NoError(t, err)
		assert.Contains(t, out, "would stamp")
	})

	t.Run("company get/update", func(t *testing.T) {
		out, err := runRoot(t, "company", "get", "-o", "json")
		require.NoError(t, err)
		assert.Contains(t, out, "Acme")

		out, err = runRoot(t, "company", "update", "--set", "name=NewCo", "-o", "json")
		require.NoError(t, err)
		assert.Contains(t, out, "NewCo")
	})

	t.Run("resource actions", func(t *testing.T) {
		out, err := runRoot(t, "users", "self", "-o", "json")
		require.NoError(t, err)
		assert.Contains(t, out, "Tester")

		_, err = runRoot(t, "journals", "balance", "-o", "json")
		require.NoError(t, err)

		_, err = runRoot(t, "categories", "settings", "-o", "json")
		require.NoError(t, err)

		_, err = runRoot(t, "categories", "set-settings", "--set", "foo=bar", "-o", "json")
		require.NoError(t, err)
	})

	t.Run("reports", func(t *testing.T) {
		out, err := runRoot(t, "reports", "sales-by-seller", "--from", "2026-01-01", "--to", "2026-03-31", "-o", "json")
		require.NoError(t, err)
		assert.Contains(t, out, "Sam")

		_, err = runRoot(t, "reports", "income-statement", "--from", "2026-01-01", "--to", "2026-03-31", "-o", "json")
		require.NoError(t, err)

		_, err = runRoot(t, "reports", "account-statement", "--from", "2026-01-01", "--to", "2026-03-31", "-o", "json")
		require.NoError(t, err)
	})

	t.Run("auth login and logout", func(t *testing.T) {
		keyring.MockInit()
		out, err := runRoot(t, "auth", "login", "--email", "new@x.com", "--token", "newtok", "--profile", "fresh")
		require.NoError(t, err)
		assert.Contains(t, out, "Logged in")

		_, err = runRoot(t, "auth", "logout", "--profile", "fresh")
		require.NoError(t, err)
	})
}

// A corrupt emitted.json must stop emission before any API call: proceeding
// with an empty guard could double-emit already-stamped invoices.
func TestEmit_AbortsOnCorruptCache(t *testing.T) {
	dir := t.TempDir()
	t.Setenv(config.EnvBaseURL, "http://127.0.0.1:1")
	t.Setenv(config.EnvEmail, "e@x.com")
	t.Setenv(config.EnvToken, "tok")
	t.Setenv(config.EnvProfile, "")
	t.Setenv(config.EnvConfig, filepath.Join(dir, "config.yaml"))
	require.NoError(t, os.WriteFile(filepath.Join(dir, "emitted.json"), []byte("{corrupt"), 0o600))

	_, err := runRoot(t, "invoices", "emit", "7")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "idempotency cache")
}
