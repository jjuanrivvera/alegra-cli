package api

import (
	"context"
	"net/http"
	"strconv"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestParseRetryAfter_UnitSuffixNotMisread guards against the "+s" parsing hack:
// a Retry-After carrying a stray unit ("1m") must not be misread as 1ms.
func TestParseRetryAfter_UnitSuffixNotMisread(t *testing.T) {
	d, ok := parseRetryAfter("30")
	require.True(t, ok)
	assert.Equal(t, 30*time.Second, d)

	d, ok = parseRetryAfter("0.5") // fractional delta-seconds still works
	require.True(t, ok)
	assert.Equal(t, 500*time.Millisecond, d)

	// Values carrying a unit are not valid RFC 9110 delta-seconds; they must be
	// ignored (fall back to exponential), never misinterpreted.
	_, ok = parseRetryAfter("1m")
	assert.False(t, ok, `"1m" must not be parsed (the old code turned it into 1ms)`)
	_, ok = parseRetryAfter("30s")
	assert.False(t, ok)
	_, ok = parseRetryAfter("NaN")
	assert.False(t, ok)
}

// TestResource_Create_DoesNotRetryPOST pins finding H3: a non-idempotent POST
// that fails transiently must surface the error rather than silently retrying,
// because the server may already have processed it (e.g. stamped an invoice) and
// a resend would duplicate the fiscal document.
func TestResource_Create_DoesNotRetryPOST(t *testing.T) {
	var posts int32
	c := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		atomic.AddInt32(&posts, 1)
		w.WriteHeader(http.StatusServiceUnavailable) // transient 5xx
	})
	c.retryPolicy = fastRetry(3)

	res := NewResource[map[string]any](c, "invoices")
	_, err := res.Create(context.Background(), map[string]any{"x": 1})
	require.Error(t, err, "a failed POST must surface the error, not silently retry")
	assert.Equal(t, int32(1), atomic.LoadInt32(&posts),
		"POST must not be auto-retried (H3): a retried create/stamp could duplicate a record")
}

// TestResource_Delete_RetriesIdempotent confirms the other side of H3: DELETE is
// idempotent, so a transient failure is still retried (resilience preserved).
func TestResource_Delete_RetriesIdempotent(t *testing.T) {
	var calls int32
	c := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodDelete, r.Method)
		if atomic.AddInt32(&calls, 1) == 1 {
			w.WriteHeader(http.StatusServiceUnavailable)
			return
		}
		_, _ = w.Write([]byte(`{}`))
	})
	c.retryPolicy = fastRetry(3)

	res := NewResource[map[string]any](c, "contacts")
	require.NoError(t, res.Delete(context.Background(), "5"))
	assert.Equal(t, int32(2), atomic.LoadInt32(&calls), "DELETE is idempotent and should retry once")
}

// TestResource_ListAll_StopsOnEmptyTrailingPage exercises the exact-multiple
// boundary (T2): two full pages then an empty page. A `<=` instead of `<` in the
// stop condition would mis-handle this.
func TestResource_ListAll_StopsOnEmptyTrailingPage(t *testing.T) {
	c := newTestClient(t, pagedArrayHandler([]string{`{"id":"1"}`, `{"id":"2"}`, `{"id":"3"}`, `{"id":"4"}`}))
	res := NewResource[map[string]any](c, "contacts")
	all, err := res.ListAll(context.Background(), ListParams{Limit: 2}, 0)
	require.NoError(t, err)
	assert.Len(t, all, 4)
}

// TestResource_ListAll_StopsAtMaxPages pins the page-cap guard (T2): an endpoint
// that always returns a full page must not loop forever.
func TestResource_ListAll_StopsAtMaxPages(t *testing.T) {
	c := newTestClient(t, func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(`[{"id":"1"},{"id":"2"}]`)) // always a full page
	})
	res := NewResource[map[string]any](c, "contacts")
	all, err := res.ListAll(context.Background(), ListParams{Limit: 2}, 3)
	require.NoError(t, err)
	assert.Len(t, all, 6, "3 pages * 2 records, then stop at the cap")
}

// TestResource_Count_WalksWhenNoTotalReported pins finding M2: when the endpoint
// reports no total, Count must walk every page and report the real total — not
// the size of the 1-record metadata probe.
func TestResource_Count_WalksWhenNoTotalReported(t *testing.T) {
	c := newTestClient(t, pagedArrayHandler([]string{`{"id":"1"}`, `{"id":"2"}`, `{"id":"3"}`}))
	res := NewResource[map[string]any](c, "contacts")
	n, err := res.Count(context.Background(), ListParams{})
	require.NoError(t, err)
	assert.Equal(t, int64(3), n, "a no-total endpoint must count via a full walk, not the 1-record probe")
}

// TestResource_Count_DoesNotMutateCallerParams pins finding L8: Count must not
// leak metadata=true (or any change) back into the caller's ListParams.Extra.
func TestResource_Count_DoesNotMutateCallerParams(t *testing.T) {
	c := newTestClient(t, func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(`{"metadata":{"total":7}}`))
	})
	res := NewResource[map[string]any](c, "contacts")
	params := ListParams{Extra: map[string][]string{"status": {"open"}}}
	_, err := res.Count(context.Background(), params)
	require.NoError(t, err)
	assert.Empty(t, params.Extra.Get("metadata"), "Count must not mutate the caller's Extra map")
	assert.Equal(t, "open", params.Extra.Get("status"), "the caller's existing filters must be untouched")
}

// pagedArrayHandler returns a fake list endpoint that honors start/limit and
// emits a bare JSON array (no total field), so it can drive pagination tests.
func pagedArrayHandler(records []string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start, _ := strconv.Atoi(r.URL.Query().Get("start"))
		limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
		if start > len(records) {
			start = len(records)
		}
		end := len(records)
		if limit > 0 {
			end = min(start+limit, len(records))
		}
		if end < start {
			end = start
		}
		_, _ = w.Write([]byte("[" + strings.Join(records[start:end], ",") + "]"))
	}
}
