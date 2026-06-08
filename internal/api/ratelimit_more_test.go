package api

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRateLimiter_ThrottleHalvesToFloor(t *testing.T) {
	r := NewRateLimiter(10, nil)
	r.Throttle()
	assert.InDelta(t, 5.0, r.current, 0.001, "first throttle halves")
	for range 20 {
		r.Throttle()
	}
	assert.GreaterOrEqual(t, r.current, 0.5, "never drops below the floor")
	assert.InDelta(t, 0.5, r.current, 0.001)
}

func TestRateLimiter_RestoreTowardBase(t *testing.T) {
	r := NewRateLimiter(10, nil)
	r.Throttle() // 10 -> 5
	r.Restore()  // 5 * 1.5 = 7.5
	assert.InDelta(t, 7.5, r.current, 0.001)
	for range 10 {
		r.Restore()
	}
	assert.InDelta(t, 10.0, r.current, 0.001, "never exceeds the base rate")
}

func TestRateLimiter_Snapshot(t *testing.T) {
	r := NewRateLimiter(5, nil)
	l, rem, reset := r.Snapshot()
	assert.Equal(t, -1, l)
	assert.Equal(t, -1, rem)
	assert.Equal(t, -1, reset)

	r.Observe(150, 140, 60)
	l, rem, reset = r.Snapshot()
	assert.Equal(t, 150, l)
	assert.Equal(t, 140, rem)
	assert.Equal(t, 60, reset)
}
