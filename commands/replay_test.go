package commands

import (
	"encoding/base64"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/jjuanrivvera/alegra-cli/internal/config"
)

// recordedRequest captures the shape of one request the CLI sent, so replay
// tests can assert wire behavior (auth, pagination params) and not just output.
type recordedRequest struct {
	Method string
	Path   string
	Query  string
	Auth   string
}

// replayServer serves recorded API fixtures from testdata/replay — the middle
// layer between hand-written httptest JSON and live smoke tests: realistic
// response bodies (string amounts, padded decimals, mixed id types) without
// needing an Alegra account. Fixture files are named METHOD_path[_startN].json.
func replayServer(t *testing.T) (*httptest.Server, *[]recordedRequest) {
	t.Helper()
	var mu sync.Mutex
	var reqs []recordedRequest
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		reqs = append(reqs, recordedRequest{
			Method: r.Method,
			Path:   r.URL.Path,
			Query:  r.URL.RawQuery,
			Auth:   r.Header.Get("Authorization"),
		})
		mu.Unlock()

		name := r.Method + strings.ReplaceAll(r.URL.Path, "/", "_")
		// Paginated endpoints keep one fixture per page, keyed by offset.
		if start := r.URL.Query().Get("start"); start != "" && start != "0" {
			if b, err := os.ReadFile(filepath.Join("testdata", "replay", name+"_start"+start+".json")); err == nil {
				w.Header().Set("Content-Type", "application/json")
				_, _ = w.Write(b)
				return
			}
		}
		for _, candidate := range []string{name + "_start0.json", name + ".json"} {
			if b, err := os.ReadFile(filepath.Join("testdata", "replay", candidate)); err == nil {
				w.Header().Set("Content-Type", "application/json")
				_, _ = w.Write(b)
				return
			}
		}
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(`{"message":"no replay fixture for ` + name + `"}`))
	}))
	t.Cleanup(srv.Close)
	return srv, &reqs
}

func TestReplay_RecordedFixtures(t *testing.T) {
	srv, reqs := replayServer(t)
	t.Setenv(config.EnvBaseURL, srv.URL)
	t.Setenv(config.EnvEmail, "e@x.com")
	t.Setenv(config.EnvToken, "tok")
	t.Setenv(config.EnvProfile, "")
	t.Setenv(config.EnvConfig, filepath.Join(t.TempDir(), "config.yaml"))

	t.Run("invoice amounts survive byte-for-byte", func(t *testing.T) {
		out, err := runRoot(t, "invoices", "list", "-o", "json")
		require.NoError(t, err)
		// The fixture carries amounts a float64 pipeline rewrites: zero-padded
		// cents ("19928.00" → 19928), long padding ("0.0000000000" → 0), and a
		// 17-significant-digit total (…56.78 → …56.8). All must come out intact
		// (Money emits them as bare JSON numbers, preserving the exact text).
		assert.Contains(t, out, `"total": 19928.00`)
		assert.Contains(t, out, `"totalPaid": 0.0000000000`)
		assert.Contains(t, out, "1234567890123456.78")
		assert.NotContains(t, out, "1234567890123456.8\n")
		assert.Contains(t, out, "Rivera Refrigeración S.A.S.")
	})

	t.Run("get renders nested recorded shapes", func(t *testing.T) {
		out, err := runRoot(t, "invoices", "get", "1", "-o", "json")
		require.NoError(t, err)
		assert.Contains(t, out, "FE341")
		assert.Contains(t, out, `"price": 16746.218487`)
		assert.Contains(t, out, "IVA 19%")
	})

	t.Run("table renders recorded list", func(t *testing.T) {
		out, err := runRoot(t, "invoices", "list")
		require.NoError(t, err)
		assert.Contains(t, out, "TOTAL")
		assert.Contains(t, out, "19928.00")
	})

	t.Run("--all walks offset pagination", func(t *testing.T) {
		before := len(*reqs)
		out, err := runRoot(t, "items", "list", "--all", "-o", "json")
		require.NoError(t, err)
		assert.Contains(t, out, "Repuesto 001")
		assert.Contains(t, out, "Repuesto 032")

		pages := (*reqs)[before:]
		require.Len(t, pages, 2, "30-item page must trigger exactly one follow-up fetch")
		assert.Contains(t, pages[0].Query, "start=0")
		assert.Contains(t, pages[0].Query, "limit=30")
		assert.Contains(t, pages[1].Query, "start=30")
	})

	t.Run("requests authenticate with basic email:token", func(t *testing.T) {
		want := "Basic " + base64.StdEncoding.EncodeToString([]byte("e@x.com:tok"))
		require.NotEmpty(t, *reqs)
		for _, r := range *reqs {
			assert.Equal(t, want, r.Auth, "%s %s", r.Method, r.Path)
		}
	})
}
