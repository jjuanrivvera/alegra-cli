package api

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/time/rate"
)

func TestRateLimiter_ObserveThrottlesNearLimit(t *testing.T) {
	rl := NewRateLimiter(5.0, nil)

	// Plenty of budget → stay at baseline.
	rl.Observe(150, 120, 40)
	assert.Equal(t, rate.Limit(5.0), rl.limiter.Limit())

	// Down to the last fifth → spread remaining across the reset window.
	rl.Observe(150, 20, 40) // 20/40 = 0.5 rps
	assert.InDelta(t, 0.5, float64(rl.limiter.Limit()), 0.001)

	// Exhausted → floor, never zero.
	rl.Observe(150, 0, 30)
	assert.InDelta(t, 0.2, float64(rl.limiter.Limit()), 0.001)
}

func TestRateLimiter_ObserveNoHeaders(t *testing.T) {
	rl := NewRateLimiter(5.0, nil)
	rl.Observe(-1, -1, -1) // missing headers: no change
	assert.Equal(t, rate.Limit(5.0), rl.limiter.Limit())
}
