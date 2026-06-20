package commands

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync/atomic"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zalando/go-keyring"

	"github.com/jjuanrivvera/alegra-cli/internal/config"
)

// These tests exercise the failure paths the integration suite never hit: the
// fake servers previously only returned 200/404, so partial batch failures,
// API errors mid-operation, and persistence failures were all untested.

// failureTestEnv points the CLI at srv with env-only credentials and an
// isolated config dir, returning that dir.
func failureTestEnv(t *testing.T, srv *httptest.Server) string {
	t.Helper()
	dir := t.TempDir()
	t.Setenv(config.EnvBaseURL, srv.URL)
	t.Setenv(config.EnvEmail, "e@x.com")
	t.Setenv(config.EnvToken, "tok")
	t.Setenv(config.EnvProfile, "")
	t.Setenv(config.EnvConfig, filepath.Join(dir, "config.yaml"))
	return dir
}

func TestImport_PartialFailureFailsLoudly(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		body, _ := io.ReadAll(r.Body)
		var m map[string]any
		_ = json.Unmarshal(body, &m)
		if m["name"] == "Bad" {
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte(`{"message":"name is invalid","code":400}`))
			return
		}
		_, _ = w.Write([]byte(`{"id":"9","name":"ok"}`))
	}))
	t.Cleanup(srv.Close)
	dir := failureTestEnv(t, srv)

	csvFile := filepath.Join(dir, "rows.csv")
	require.NoError(t, os.WriteFile(csvFile, []byte("name\nGood\nBad\nAlsoGood\n"), 0o600))

	out, err := runRoot(t, "contacts", "import", "-f", csvFile)
	// A partial import must fail the command (exit code), report which row
	// broke, and still account for the rows that were created.
	require.Error(t, err, "partial failure must not look like success to pipelines")
	assert.Contains(t, err.Error(), "1 row(s) failed")
	assert.Contains(t, out, "[row 2] FAILED")
	assert.Contains(t, out, "Imported 2, failed 1")
}

func TestEmit_BatchAPIErrorFailsAndKeepsCacheClean(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if strings.HasSuffix(r.URL.Path, "/invoices/stamp") {
			w.WriteHeader(http.StatusBadRequest) // 4xx: not retried, fails fast
			_, _ = w.Write([]byte(`{"message":"EPR001: certificado vencido"}`))
			return
		}
		_, _ = w.Write([]byte(`{}`))
	}))
	t.Cleanup(srv.Close)
	dir := failureTestEnv(t, srv)

	out, err := runRoot(t, "invoices", "emit", "7", "--force")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "1 batch(es) failed")
	assert.Contains(t, out, "FAILED")

	// A failed batch must never be recorded as emitted.
	data, rerr := os.ReadFile(filepath.Join(dir, "emitted.json"))
	if rerr == nil {
		assert.NotContains(t, string(data), ":7", "failed invoice must not be marked emitted")
	}
}

func TestEmit_CacheSaveFailureStopsEmission(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("read-only directory permissions are not enforced on Windows")
	}
	if os.Geteuid() == 0 {
		t.Skip("root bypasses directory write permissions, so the cache save cannot be forced to fail")
	}
	var stamps int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if strings.HasSuffix(r.URL.Path, "/invoices/stamp") {
			atomic.AddInt32(&stamps, 1)
			_, _ = w.Write([]byte(`{"stamped":true}`))
			return
		}
		_, _ = w.Write([]byte(`{}`))
	}))
	t.Cleanup(srv.Close)
	dir := failureTestEnv(t, srv)

	// Make the config dir unwritable so the post-batch cache save fails.
	require.NoError(t, os.Chmod(dir, 0o500))
	t.Cleanup(func() { _ = os.Chmod(dir, 0o700) })

	// 12 ids → 2 batches. The save failure after batch 1 must stop the run
	// before batch 2 is stamped, and the error must name the affected ids.
	ids := []string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "10", "11", "12"}
	out, err := runRoot(t, append([]string{"invoices", "emit", "--force"}, ids...)...)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "could not record them in the idempotency cache")
	assert.Contains(t, err.Error(), "1") // affected ids are listed
	assert.Contains(t, out, "stamped")
	assert.Equal(t, int32(1), atomic.LoadInt32(&stamps), "must stop stamping once the guard cannot record progress")
}

func TestDelete_AbortsWithoutConfirmation(t *testing.T) {
	var deletes int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if r.Method == http.MethodDelete {
			atomic.AddInt32(&deletes, 1)
		}
		_, _ = w.Write([]byte(`{}`))
	}))
	t.Cleanup(srv.Close)
	failureTestEnv(t, srv)

	// Piped/empty stdin without --yes: the prompt cannot be answered, so the
	// delete must abort — and the DELETE request must never be sent.
	rootCmd.SetIn(strings.NewReader(""))
	t.Cleanup(func() { rootCmd.SetIn(nil) })

	_, err := runRoot(t, "contacts", "delete", "1")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "aborted")
	assert.Zero(t, atomic.LoadInt32(&deletes), "no DELETE may be sent without confirmation")
}

func TestAuthLogin_NeverPersistsTokenToConfig(t *testing.T) {
	keyring.MockInit()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"Tester","applicationVersion":"colombia"}`))
	}))
	t.Cleanup(srv.Close)
	dir := failureTestEnv(t, srv)
	t.Setenv(config.EnvToken, "") // force login to use the flag token

	const secret = "super-secret-token"
	_, err := runRoot(t, "auth", "login", "--email", "a@x.com", "--token", secret, "--profile", "p1")
	require.NoError(t, err)

	// The token belongs in the keyring only; the YAML must never contain it.
	data, rerr := os.ReadFile(filepath.Join(dir, "config.yaml"))
	require.NoError(t, rerr)
	assert.NotContains(t, string(data), secret, "plaintext token leaked into config.yaml")
}
