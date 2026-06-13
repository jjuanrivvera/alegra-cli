package api

import (
	"context"
	"errors"
	"log/slog"
	"math"
	"net/http"
	"strconv"
	"strings"
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

// idempotentMethod reports whether repeating method is safe — i.e. running it
// twice leaves the same server state as running it once. GET/HEAD/PUT/DELETE/
// OPTIONS/TRACE qualify (RFC 9110 §9.2.2); POST and PATCH do NOT.
//
// This gates automatic retries: a retried POST /invoices/stamp could emit a
// duplicate electronic invoice at the tax authority (DIAN/SAT), which is
// legally irreversible. When we can't be sure a non-idempotent request didn't
// already take effect, we must surface the error instead of silently resending.
func idempotentMethod(method string) bool {
	switch strings.ToUpper(method) {
	case http.MethodGet, http.MethodHead, http.MethodPut,
		http.MethodDelete, http.MethodOptions, http.MethodTrace:
		return true
	default:
		return false
	}
}

// shouldRetry reports whether a response/error pair warrants a retry of a
// request using the given HTTP method.
//
// 429 (rate limited) is always retryable: the server rejected the request
// before processing it, so no state changed regardless of method. Transient
// network errors and 5xx are retried only for idempotent methods, because for a
// POST/PATCH the request may have already taken effect on the server even though
// we never saw the response. We never retry context cancellation or 4xx (other
// than 429).
func (p *RetryPolicy) shouldRetry(method string, resp *http.Response, err error) bool {
	if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
		return false
	}
	if resp != nil && resp.StatusCode == http.StatusTooManyRequests {
		return true
	}
	if !idempotentMethod(method) {
		return false
	}
	if err != nil {
		return true // transient network-level error on an idempotent request
	}
	if resp == nil {
		return false
	}
	return resp.StatusCode >= 500
}

// backoff computes the wait duration before the given attempt (0-indexed),
// honoring a Retry-After header when present.
func (p *RetryPolicy) backoff(attempt int, resp *http.Response) time.Duration {
	if resp != nil {
		if ra := resp.Header.Get("Retry-After"); ra != "" {
			// Clamp to MaxBackoff: a server (or a proxy/CDN/WAF in front of it)
			// can send an absurd Retry-After that would otherwise hang the client
			// for years.
			if d, ok := parseRetryAfter(ra); ok && d > 0 {
				return min(d, p.MaxBackoff)
			}
		}
		// On a 429, Alegra tells us exactly how long until the window resets.
		if resp.StatusCode == http.StatusTooManyRequests {
			if rs := resp.Header.Get("X-Rate-Limit-Reset"); rs != "" {
				if d, ok := parseRetryAfter(rs); ok && d > 0 {
					return min(d, p.MaxBackoff)
				}
			}
		}
	}
	// Exponential backoff. No jitter: the CLI issues requests sequentially from a
	// single process, so there is no fleet of clients to de-synchronize — jitter
	// would only add nondeterminism without reducing any real contention.
	d := time.Duration(float64(p.InitialBackoff) * math.Pow(2, float64(attempt)))
	return min(d, p.MaxBackoff)
}

// parseRetryAfter parses a Retry-After / X-Rate-Limit-Reset value in either RFC
// 9110 form: delta-seconds ("120") or an HTTP-date ("Wed, 21 Oct 2025 07:28:00
// GMT"). A past HTTP-date yields a non-positive duration, which callers treat as
// "no usable hint" and fall back to exponential backoff.
//
// The seconds form is parsed as a plain number, NOT via time.ParseDuration with
// an appended "s": that would misread a value carrying a stray unit — "1m" would
// become "1ms" (1 millisecond) and "30s" would fail outright — so a misbehaving
// proxy could trick the client into hammering the server.
func parseRetryAfter(v string) (time.Duration, bool) {
	v = strings.TrimSpace(v)
	if secs, err := strconv.ParseFloat(v, 64); err == nil && !math.IsNaN(secs) && !math.IsInf(secs, 0) {
		return time.Duration(secs * float64(time.Second)), true
	}
	if t, err := http.ParseTime(v); err == nil {
		return time.Until(t), true
	}
	return 0, false
}
