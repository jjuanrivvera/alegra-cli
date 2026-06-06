package api

import (
	"context"
	"log/slog"
	"sync"

	"golang.org/x/time/rate"
)

// RateLimiter throttles outbound requests. Alegra enforces per-plan request
// quotas; staying well under them avoids 429s. The limiter is adaptive: when
// the server signals throttling (429), the effective rate is reduced.
type RateLimiter struct {
	limiter *rate.Limiter
	mu      sync.Mutex
	base    float64
	current float64
	logger  *slog.Logger

	// Last quota reported by the server (-1 = unknown).
	lastLimit, lastRemaining, lastReset int
}

// Snapshot returns the most recent server-reported quota (limit, remaining,
// reset seconds). Values are -1 until a response has been observed.
func (r *RateLimiter) Snapshot() (limit, remaining, reset int) {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.lastLimit, r.lastRemaining, r.lastReset
}

// NewRateLimiter creates a limiter allowing requestsPerSecond, with a small
// burst so short command bursts are not serialized unnecessarily.
func NewRateLimiter(requestsPerSecond float64, logger *slog.Logger) *RateLimiter {
	if requestsPerSecond <= 0 {
		requestsPerSecond = 5.0
	}
	if logger == nil {
		logger = slog.Default()
	}
	burst := max(int(requestsPerSecond), 1)
	return &RateLimiter{
		limiter:       rate.NewLimiter(rate.Limit(requestsPerSecond), burst),
		base:          requestsPerSecond,
		current:       requestsPerSecond,
		logger:        logger,
		lastLimit:     -1,
		lastRemaining: -1,
		lastReset:     -1,
	}
}

// Wait blocks until a request may proceed or ctx is cancelled.
func (r *RateLimiter) Wait(ctx context.Context) error {
	return r.limiter.Wait(ctx)
}

// Throttle halves the current rate (down to a floor) after a 429, so we back
// off proactively rather than hammering the API.
func (r *RateLimiter) Throttle() {
	r.mu.Lock()
	defer r.mu.Unlock()
	next := r.current / 2
	if next < 0.5 {
		next = 0.5
	}
	if next != r.current {
		r.current = next
		r.limiter.SetLimit(rate.Limit(next))
		r.logger.Debug("rate limiter throttled", "requests_per_second", next)
	}
}

// Observe adapts the limiter to Alegra's reported quota from the
// X-Rate-Limit-* response headers. While plenty of budget remains it runs at the
// configured baseline; in the final fifth of the quota it spreads the remaining
// calls across the reset window so we glide to the limit instead of slamming it.
func (r *RateLimiter) Observe(limit, remaining, resetSeconds int) {
	if limit <= 0 {
		return
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	r.lastLimit, r.lastRemaining, r.lastReset = limit, remaining, resetSeconds
	var next float64
	if float64(remaining) > 0.2*float64(limit) {
		next = r.base
	} else {
		secs := max(resetSeconds, 1)
		next = min(max(float64(remaining)/float64(secs), 0.2), r.base)
	}
	if next != r.current {
		r.current = next
		r.limiter.SetLimit(rate.Limit(next))
		r.logger.Debug("rate limiter adjusted from headers",
			"remaining", remaining, "limit", limit, "reset_s", resetSeconds, "rps", next)
	}
}

// Restore gradually returns the rate toward the configured baseline after a
// successful request.
func (r *RateLimiter) Restore() {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.current >= r.base {
		return
	}
	next := r.current * 1.5
	if next > r.base {
		next = r.base
	}
	r.current = next
	r.limiter.SetLimit(rate.Limit(next))
}
