package api

import (
	"context"
	"errors"
	"log/slog"
	"math"
	"net/http"
	"time"
)

// RetryPolicy controls automatic retries for transient failures.
type RetryPolicy struct {
	MaxRetries     int
	InitialBackoff time.Duration
	MaxBackoff     time.Duration
	Logger         *slog.Logger
}

// DefaultRetryPolicy returns a sensible default: 3 retries with exponential
// backoff capped at 8s.
func DefaultRetryPolicy(logger *slog.Logger) *RetryPolicy {
	if logger == nil {
		logger = slog.Default()
	}
	return &RetryPolicy{
		MaxRetries:     3,
		InitialBackoff: 500 * time.Millisecond,
		MaxBackoff:     8 * time.Second,
		Logger:         logger,
	}
}

// shouldRetry reports whether a response/error pair warrants a retry.
//
// We retry on rate-limiting (429), server errors (5xx), and transient network
// errors. We never retry on context cancellation or 4xx client errors (other
// than 429).
func (p *RetryPolicy) shouldRetry(resp *http.Response, err error) bool {
	if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
		return false
	}
	if err != nil {
		return true // network-level error
	}
	if resp == nil {
		return false
	}
	if resp.StatusCode == http.StatusTooManyRequests {
		return true
	}
	return resp.StatusCode >= 500
}

// backoff computes the wait duration before the given attempt (0-indexed),
// honoring a Retry-After header when present.
func (p *RetryPolicy) backoff(attempt int, resp *http.Response) time.Duration {
	if resp != nil {
		if ra := resp.Header.Get("Retry-After"); ra != "" {
			if secs, err := time.ParseDuration(ra + "s"); err == nil {
				return secs
			}
		}
		// On a 429, Alegra tells us exactly how long until the window resets.
		if resp.StatusCode == http.StatusTooManyRequests {
			if rs := resp.Header.Get("X-Rate-Limit-Reset"); rs != "" {
				if secs, err := time.ParseDuration(rs + "s"); err == nil && secs > 0 {
					if secs > p.MaxBackoff {
						return p.MaxBackoff
					}
					return secs
				}
			}
		}
	}
	d := time.Duration(float64(p.InitialBackoff) * math.Pow(2, float64(attempt)))
	return min(d, p.MaxBackoff)
}
