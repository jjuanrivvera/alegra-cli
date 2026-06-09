package commands

import (
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/jjuanrivvera/alegra-cli/internal/config"
)

func TestScalar(t *testing.T) {
	assert.Equal(t, "hi", scalar("hi"))
	assert.Equal(t, "true", scalar(true))
	assert.Equal(t, "false", scalar(false))
	assert.Equal(t, "12", scalar(float64(12))) // whole numbers drop the decimal
	assert.Equal(t, "1.5", scalar(float64(1.5)))
	assert.Equal(t, "", scalar(nil))
	assert.Equal(t, "", scalar(map[string]any{"x": 1})) // nested values are not scalars
}

func TestIDCompletions(t *testing.T) {
	type rec struct {
		ID     string `json:"id"`
		Name   string `json:"name"`
		Number string `json:"number"`
	}
	items := []rec{
		{ID: "10", Name: "Acme"},
		{ID: "11", Number: "FV-2"}, // falls back to number when name is empty
		{ID: "20"},                 // no label
	}

	all := idCompletions(items, "")
	require.Equal(t, []string{"10\tAcme", "11\tFV-2", "20"}, all)

	// Prefix filtering keeps only ids that start with the typed text.
	assert.Equal(t, []string{"10\tAcme", "11\tFV-2"}, idCompletions(items, "1"))
	assert.Empty(t, idCompletions(items, "9"))
}

func TestCompletionLabel(t *testing.T) {
	assert.Equal(t, "Acme", completionLabel(map[string]any{"name": "Acme", "email": "a@x.com"}))
	assert.Equal(t, "a@x.com", completionLabel(map[string]any{"email": "a@x.com"}))
	assert.Equal(t, "", completionLabel(map[string]any{"unrelated": "x"}))
}

func TestFixedCompleter(t *testing.T) {
	f := fixedCompleter("table", "json", "yaml", "csv")
	got, dir := f(nil, nil, "j")
	assert.Equal(t, []string{"json"}, got)
	assert.Equal(t, cobra.ShellCompDirectiveNoFileComp, dir)

	got, _ = f(nil, nil, "")
	assert.Len(t, got, 4)
}

func TestCompleteColumns(t *testing.T) {
	cmd := &cobra.Command{}
	withColumns(cmd, []string{"id", "date", "dueDate", "status"})

	got, dir := completeColumns(cmd, nil, "d")
	assert.Equal(t, []string{"date", "dueDate"}, got)
	assert.Equal(t, cobra.ShellCompDirectiveNoSpace|cobra.ShellCompDirectiveNoFileComp, dir)

	// Comma-separated: prior segments are preserved, the last is completed.
	got, _ = completeColumns(cmd, nil, "id,s")
	assert.Equal(t, []string{"id,status"}, got)

	// A command without the annotation yields nothing rather than panicking.
	got, _ = completeColumns(&cobra.Command{}, nil, "")
	assert.Empty(t, got)
}

func TestFilterEnum(t *testing.T) {
	// Explicit values win.
	assert.Equal(t, []string{"a", "b"}, filterEnum(listFilter{Values: []string{"a", "b"}}))
	// A compact comma-separated usage is treated as an enum.
	assert.Equal(t, []string{"open", "closed", "void"}, filterEnum(listFilter{Usage: "open,closed,void"}))
	// Prose usage (contains spaces) is not an enum.
	assert.Nil(t, filterEnum(listFilter{Usage: "Filter by client ID"}))
	// Single-value usage is not an enum.
	assert.Nil(t, filterEnum(listFilter{Usage: "name"}))
}

func TestCompleteProfiles(t *testing.T) {
	// config.New resolves its path from EnvConfig, so Save writes where Load reads.
	t.Setenv(config.EnvConfig, filepath.Join(t.TempDir(), "config.yaml"))

	cfg := config.New()
	cfg.SetProfile(&config.Profile{Name: "prod", Email: "prod@x.com"})
	cfg.SetProfile(&config.Profile{Name: "sandbox", Email: "sbx@x.com"})
	require.NoError(t, cfg.Save())

	got, dir := completeProfiles(nil, nil, "")
	assert.Equal(t, cobra.ShellCompDirectiveNoFileComp, dir)
	require.Len(t, got, 2)
	assert.Equal(t, "prod\tprod@x.com", got[0]) // sorted, annotated with email
	assert.Equal(t, "sandbox\tsbx@x.com", got[1])

	got, _ = completeProfiles(nil, nil, "pro")
	assert.Equal(t, []string{"prod\tprod@x.com"}, got)
}

// TestCompletion_LiveIDs drives the real `__complete` path end-to-end so the
// ValidArgsFunction wiring (get/update/delete/actions) is exercised against a
// server, not just the projection helper.
func TestCompletion_LiveIDs(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if r.Method == http.MethodGet && r.URL.Path == "/invoices" {
			_, _ = w.Write([]byte(`[{"id":"100","status":"open"},{"id":"101","status":"void"}]`))
			return
		}
		w.WriteHeader(http.StatusNotFound)
	}))
	t.Cleanup(srv.Close)

	t.Setenv(config.EnvBaseURL, srv.URL)
	t.Setenv(config.EnvEmail, "e@x.com")
	t.Setenv(config.EnvToken, "tok")
	t.Setenv(config.EnvProfile, "")
	t.Setenv(config.EnvConfig, filepath.Join(t.TempDir(), "config.yaml"))

	// The label falls through labelFields to the first present field; for the
	// Invoice type (which has no name/number column) that is the status.
	t.Run("get completes live ids with labels", func(t *testing.T) {
		out, err := runRoot(t, "__complete", "invoices", "get", "")
		require.NoError(t, err)
		assert.Contains(t, out, "100\topen")
		assert.Contains(t, out, "101\tvoid")
	})

	t.Run("prefix narrows the live ids", func(t *testing.T) {
		out, err := runRoot(t, "__complete", "invoices", "delete", "100")
		require.NoError(t, err)
		assert.Contains(t, out, "100\topen")
		assert.NotContains(t, out, "101")
	})

	t.Run("action subcommand completes ids", func(t *testing.T) {
		out, err := runRoot(t, "__complete", "invoices", "void", "")
		require.NoError(t, err)
		assert.Contains(t, out, "100\topen")
	})

	t.Run("static enum filter completes", func(t *testing.T) {
		out, err := runRoot(t, "__complete", "invoices", "list", "--status", "")
		require.NoError(t, err)
		assert.Contains(t, out, "open")
		assert.Contains(t, out, "void")
	})
}
