package api

import (
	"net/http"
	"testing"
	"time"
)

func TestBackoffClampsRetryAfter(t *testing.T) {
	p := DefaultRetryPolicy(nil) // MaxBackoff = 8s

	resp := &http.Response{
		StatusCode: http.StatusServiceUnavailable,
		Header:     http.Header{"Retry-After": {"999999999"}}, // absurd: ~31 years
	}
	if got := p.backoff(0, resp); got > p.MaxBackoff {
		t.Errorf("Retry-After not clamped: got %v, want <= %v", got, p.MaxBackoff)
	}

	// A sane Retry-After under the cap is honored as-is.
	resp.Header.Set("Retry-After", "2")
	if got := p.backoff(0, resp); got != 2*time.Second {
		t.Errorf("Retry-After=2s: got %v, want 2s", got)
	}
}

func TestBackoffClampsRateLimitReset(t *testing.T) {
	p := DefaultRetryPolicy(nil)
	resp := &http.Response{
		StatusCode: http.StatusTooManyRequests,
		Header:     http.Header{"X-Rate-Limit-Reset": {"999999999"}},
	}
	if got := p.backoff(0, resp); got > p.MaxBackoff {
		t.Errorf("X-Rate-Limit-Reset not clamped: got %v, want <= %v", got, p.MaxBackoff)
	}
}
