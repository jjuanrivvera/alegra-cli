package api

import (
	"context"
	"net/http"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// These tests pin the client's behavior under hostile responses: malformed
// bodies, non-JSON error pages, cancellation mid-backoff, and the 429
// throttle/restore cycle. None of them were previously exercised — the fake
// servers only ever spoke well-formed JSON.

func TestDoJSON_MalformedJSONOn200(t *testing.T) {
	c := newTestClient(t, func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(`{"id": "trunca`)) // 200 OK, broken body
	})

	var out map[string]any
	err := c.GetInto(context.Background(), "contacts/1", nil, &out)
	require.Error(t, err, "a 200 with a broken body must not be reported as success")
	assert.NotContains(t, err.Error(), "panic")
}

func TestDo_HTMLErrorBodyProducesUsableError(t *testing.T) {
	c := newTestClient(t, func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.WriteHeader(http.StatusBadRequest) // 4xx: not retried
		_, _ = w.Write([]byte(`<html><body><h1>400 Bad Request</h1></body></html>`))
	})

	var out map[string]any
	err := c.GetInto(context.Background(), "contacts", nil, &out)
	require.Error(t, err)
	apiErr, ok := AsAPIError(err)
	require.True(t, ok, "non-JSON error bodies must still yield an APIError")
	assert.Equal(t, http.StatusBadRequest, apiErr.StatusCode)
	assert.Contains(t, apiErr.Body, "Bad Request", "raw body must be preserved for diagnosis")
}

func TestDo_ContextCancelDuringBackoffAborts(t *testing.T) {
	c := newTestClient(t, func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusServiceUnavailable) // always retryable
	})
	// A long backoff makes the wait the dominant phase; cancellation must cut
	// it short instead of sleeping out the full window.
	c.retryPolicy = &RetryPolicy{MaxRetries: 3, InitialBackoff: 30 * time.Second, MaxBackoff: 30 * time.Second}

	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	start := time.Now()
	var out map[string]any
	err := c.GetInto(ctx, "contacts", nil, &out)
	require.Error(t, err)
	assert.ErrorIs(t, err, context.DeadlineExceeded)
	assert.Less(t, time.Since(start), 5*time.Second, "cancellation must abort the backoff wait")
}

func TestDo_429ThrottlesThenRestores(t *testing.T) {
	var calls int32
	c := newTestClient(t, func(w http.ResponseWriter, _ *http.Request) {
		if atomic.AddInt32(&calls, 1) == 1 {
			w.WriteHeader(http.StatusTooManyRequests)
			return
		}
		_, _ = w.Write([]byte(`{"id":"1"}`))
	})
	c.retryPolicy = fastRetry(2)

	before := c.rateLimiter.current
	var out map[string]any
	require.NoError(t, c.GetInto(context.Background(), "contacts/1", nil, &out))

	// The 429 halves the rate; the subsequent success restores toward base
	// (1.5x) — so the net effect of halve-then-restore lands below the
	// starting rate but above the floor.
	after := c.rateLimiter.current
	assert.Less(t, after, before, "a 429 must leave the limiter slower than it started")
	assert.Greater(t, after, before/2, "the success after the 429 must have restored part of the rate")
}

func TestBackoff_RetryAfterHTTPDateFallsBack(t *testing.T) {
	p := &RetryPolicy{MaxRetries: 3, InitialBackoff: time.Second, MaxBackoff: 8 * time.Second}
	resp := &http.Response{
		StatusCode: http.StatusTooManyRequests,
		Header:     http.Header{"Retry-After": []string{"Wed, 21 Oct 2026 07:28:00 GMT"}},
	}
	// HTTP-date is a valid RFC 9110 Retry-After form the parser doesn't
	// understand; it must fall back to exponential backoff, not panic or
	// return zero/negative.
	d := p.backoff(0, resp)
	assert.GreaterOrEqual(t, d, time.Second)
	assert.LessOrEqual(t, d, 8*time.Second)
}

func TestDo_NoRateLimitHeadersKeepsSnapshotUnset(t *testing.T) {
	c := newTestClient(t, func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(`{"id":"1"}`)) // no X-Rate-Limit-* headers
	})

	var out map[string]any
	require.NoError(t, c.GetInto(context.Background(), "contacts/1", nil, &out))
	limit, remaining, reset := c.rateLimiter.Snapshot()
	assert.Equal(t, -1, limit)
	assert.Equal(t, -1, remaining)
	assert.Equal(t, -1, reset)
}
