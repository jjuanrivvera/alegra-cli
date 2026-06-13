package api

import (
	"context"
	"errors"
	"net/http"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// fastRetry returns a retry policy with near-zero backoff for tests.
func fastRetry(maxRetries int) *RetryPolicy {
	return &RetryPolicy{MaxRetries: maxRetries, InitialBackoff: time.Millisecond, MaxBackoff: time.Millisecond}
}

func TestClient_RetriesThenSucceeds(t *testing.T) {
	var calls int32
	c := newTestClient(t, func(w http.ResponseWriter, _ *http.Request) {
		if atomic.AddInt32(&calls, 1) == 1 {
			w.WriteHeader(http.StatusInternalServerError) // first attempt fails
			return
		}
		_, _ = w.Write([]byte(`{"id":"1"}`))
	})
	c.retryPolicy = fastRetry(3)

	var out map[string]any
	require.NoError(t, c.GetInto(context.Background(), "contacts/1", nil, &out))
	assert.Equal(t, "1", out["id"])
	assert.Equal(t, int32(2), atomic.LoadInt32(&calls), "should have retried once")
}

func TestClient_GivesUpAfterMaxRetries(t *testing.T) {
	var calls int32
	c := newTestClient(t, func(w http.ResponseWriter, _ *http.Request) {
		atomic.AddInt32(&calls, 1)
		w.WriteHeader(http.StatusServiceUnavailable) // always 503
	})
	c.retryPolicy = fastRetry(2)

	var out map[string]any
	err := c.GetInto(context.Background(), "contacts", nil, &out)
	require.Error(t, err)
	apiErr, ok := AsAPIError(err)
	require.True(t, ok)
	assert.Equal(t, http.StatusServiceUnavailable, apiErr.StatusCode)
	assert.Equal(t, int32(3), atomic.LoadInt32(&calls), "1 initial + 2 retries")
}

func TestShouldRetry(t *testing.T) {
	p := &RetryPolicy{}

	// 429 is always retryable regardless of method: the server rejected the
	// request before processing it, so no state changed.
	assert.True(t, p.shouldRetry(http.MethodGet, &http.Response{StatusCode: 429}, nil))
	assert.True(t, p.shouldRetry(http.MethodPost, &http.Response{StatusCode: 429}, nil))

	// Idempotent methods retry on 5xx and transient network errors.
	assert.True(t, p.shouldRetry(http.MethodGet, &http.Response{StatusCode: 503}, nil))
	assert.True(t, p.shouldRetry(http.MethodGet, nil, errors.New("dial tcp: connection refused")))
	assert.True(t, p.shouldRetry(http.MethodPut, &http.Response{StatusCode: 500}, nil))
	assert.True(t, p.shouldRetry(http.MethodDelete, nil, errors.New("connection reset")))

	// Non-idempotent methods (POST/PATCH) must NOT retry on 5xx or network
	// errors: the request may already have taken effect (e.g. a stamped
	// invoice), and resending would duplicate it. This is finding H3.
	assert.False(t, p.shouldRetry(http.MethodPost, &http.Response{StatusCode: 503}, nil))
	assert.False(t, p.shouldRetry(http.MethodPost, nil, errors.New("dial tcp: connection refused")))
	assert.False(t, p.shouldRetry(http.MethodPatch, &http.Response{StatusCode: 500}, nil))

	// Never retry success, 4xx, a nil response, or context cancellation.
	assert.False(t, p.shouldRetry(http.MethodGet, &http.Response{StatusCode: 404}, nil))
	assert.False(t, p.shouldRetry(http.MethodGet, &http.Response{StatusCode: 200}, nil))
	assert.False(t, p.shouldRetry(http.MethodGet, nil, nil))
	assert.False(t, p.shouldRetry(http.MethodGet, nil, context.Canceled))
	assert.False(t, p.shouldRetry(http.MethodPost, nil, context.DeadlineExceeded))
}
