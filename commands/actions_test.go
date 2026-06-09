package commands

import (
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/jjuanrivvera/alegra-cli/internal/config"
)

// actionServer returns an empty array for any GET (list) and an empty object for
// any write, which is enough to drive every resource's command through happily.
func actionServer(t *testing.T) *httptest.Server {
	t.Helper()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if r.Method == http.MethodGet {
			_, _ = w.Write([]byte(`[]`))
			return
		}
		_, _ = w.Write([]byte(`{}`))
	}))
	t.Cleanup(srv.Close)
	return srv
}

func setActionEnv(t *testing.T, srv *httptest.Server) {
	t.Setenv(config.EnvBaseURL, srv.URL)
	t.Setenv(config.EnvEmail, "e@x.com")
	t.Setenv(config.EnvToken, "tok")
	t.Setenv(config.EnvProfile, "")
	t.Setenv(config.EnvConfig, filepath.Join(t.TempDir(), "config.yaml"))
}

// TestIntegration_AllResourceLists runs `<resource> list` for every registered
// resource, exercising each resource's `New` accessor closure (otherwise only
// hit for the few resources touched by other tests).
func TestIntegration_AllResourceLists(t *testing.T) {
	srv := actionServer(t)
	setActionEnv(t, srv)

	for _, c := range rootCmd.Commands() {
		hasList := false
		for _, sub := range c.Commands() {
			if sub.Name() == "list" {
				hasList = true
				break
			}
		}
		if !hasList {
			continue
		}
		name := c.Name()
		t.Run(name, func(t *testing.T) {
			_, err := runRoot(t, name, "list", "-o", "json")
			require.NoError(t, err, "%s list", name)
		})
	}
}

// TestIntegration_ResourceActions drives the custom record/collection actions so
// the generic action builders' branches (incl. the empty-response "OK" path) are
// covered.
func TestIntegration_ResourceActions(t *testing.T) {
	srv := actionServer(t)
	setActionEnv(t, srv)

	cases := [][]string{
		{"invoices", "open", "1"},
		{"invoices", "email", "1", "-d", `{"emails":["a@b.com"]}`},
		{"invoices", "preview", "-d", "{}"},
		{"credit-notes", "email", "1", "-d", "{}"},
		{"estimates", "email", "1", "-d", "{}"},
		{"remissions", "void", "1"},
		{"remissions", "open", "1"},
		{"transportation-receipts", "void", "1"},
		{"transportation-receipts", "open", "1"},
		{"transportation-receipts", "email", "1", "-d", "{}"},
		{"transportation-receipts", "preview", "-d", "{}"},
		{"payments", "void", "1"},
		{"payments", "open", "1"},
		{"payments", "stamp", "1"},
		{"purchase-orders", "void", "1"},
		{"purchase-orders", "email", "1", "-d", "{}"},
		{"purchase-orders", "comments", "1", "-d", "{}"},
		{"bank-accounts", "transfer", "1", "-d", "{}"},
		{"bills", "close", "1", "-d", "{}"},
		{"bills", "comments", "1", "-d", "{}"},
		{"bills", "advances", "1", "-d", `{"advances":[]}`},
		{"bills", "attach", "1", "-d", `{"file":"Zm9v","name":"x.pdf"}`},
		{"bills", "perceptions", "1", "-d", `{"perceptions":[]}`},
		{"bills", "retentions", "1", "-d", `{"retentions":[]}`},
		{"bills", "comment-update", "1", "10", "-d", `{"comment":"edited"}`},
		{"bills", "comment-delete", "1", "10", "--yes"},
		{"bills", "attachment-delete", "5", "--yes"},
		{"bills", "import-by-cufe", "-d", "{}"},
	}
	for _, args := range cases {
		_, err := runRoot(t, args...)
		require.NoError(t, err, "args %v", args)
	}
}

// TestItemStock covers `items stock`: the per-warehouse breakdown and the
// not-inventariable (service item) branch. It needs a server that returns an
// item object for GET /items/{id} (the shared actionServer returns [] for GET).
func TestItemStock(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.URL.Path {
		case "/items/5":
			_, _ = w.Write([]byte(`{"id":"5","name":"Widget","inventory":{"unit":"unit","availableQuantity":12,"warehouses":[{"id":"1","name":"Main","availableQuantity":10},{"id":"2","name":"Backup","availableQuantity":2}]}}`))
		case "/items/9":
			_, _ = w.Write([]byte(`{"id":"9","name":"Consulting"}`))
		default:
			_, _ = w.Write([]byte(`{}`))
		}
	}))
	t.Cleanup(srv.Close)
	setActionEnv(t, srv)

	out, err := runRoot(t, "items", "stock", "5", "-o", "json")
	require.NoError(t, err)
	require.Contains(t, out, "Main")
	require.Contains(t, out, "Backup")

	out, err = runRoot(t, "items", "stock", "9")
	require.NoError(t, err)
	require.Contains(t, out, "not inventariable")
}

// TestIntegration_MiscPaths covers a few remaining branches: report dry-run,
// verbose logging, and count.
func TestIntegration_MiscPaths(t *testing.T) {
	srv := actionServer(t)
	setActionEnv(t, srv)

	for _, args := range [][]string{
		{"reports", "sales-by-client", "--from", "2026-01-01", "--to", "2026-03-31", "--dry-run"},
		{"reports", "income-statement", "--from", "2026-01-01", "--to", "2026-03-31", "--dry-run"},
		{"contacts", "list", "-v"},
		{"contacts", "get", "1", "--dry-run"},
	} {
		_, err := runRoot(t, args...)
		require.NoError(t, err, "args %v", args)
	}
}
