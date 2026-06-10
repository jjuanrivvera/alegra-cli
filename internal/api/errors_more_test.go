package api

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAPIError_Predicates(t *testing.T) {
	assert.True(t, (&APIError{StatusCode: 404}).IsNotFound())
	assert.False(t, (&APIError{StatusCode: 500}).IsNotFound())

	assert.True(t, (&APIError{StatusCode: 401}).IsUnauthorized())
	assert.True(t, (&APIError{StatusCode: 403}).IsUnauthorized())
	assert.False(t, (&APIError{StatusCode: 404}).IsUnauthorized())

	assert.True(t, (&APIError{StatusCode: 429}).IsRateLimited())
	assert.False(t, (&APIError{StatusCode: 200}).IsRateLimited())
}

func TestAPIError_HintByStatus(t *testing.T) {
	cases := map[int]string{
		401: "authentication failed",
		402: "plan",
		403: "permission",
		404: "not found",
		429: "rate limit",
		500: "server error",
		503: "server error",
	}
	for code, want := range cases {
		got := (&APIError{StatusCode: code}).Hint()
		assert.Contains(t, got, want, "status %d", code)
	}
	// An unmapped status has no hint.
	assert.Empty(t, (&APIError{StatusCode: 418}).Hint())
}

func TestAPIError_HintStampCodes(t *testing.T) {
	// Known stamping code → mapped remedy.
	assert.Contains(t, (&APIError{StatusCode: 400, Message: "rejected EPR503"}).Hint(),
		"timed out")
	// Unknown stamping code → generic non-idempotent warning with the code.
	got := (&APIError{StatusCode: 400, Body: `{"code":"D1234"}`}).Hint()
	assert.Contains(t, got, "not idempotent")
	assert.Contains(t, got, "D1234")
}

func TestAsAPIError(t *testing.T) {
	wrapped := fmt.Errorf("context: %w", &APIError{StatusCode: 402, Message: "plan required"})
	ae, ok := AsAPIError(wrapped)
	assert.True(t, ok)
	assert.Equal(t, 402, ae.StatusCode)

	_, ok = AsAPIError(errors.New("plain error"))
	assert.False(t, ok)

	_, ok = AsAPIError(nil)
	assert.False(t, ok)
}

func TestAsAPIError_UnwrapsChainsAndJoins(t *testing.T) {
	apiErr := &APIError{StatusCode: 404, Message: "not found"}

	got, ok := AsAPIError(fmt.Errorf("wrapped: %w", apiErr))
	assert.True(t, ok)
	assert.Same(t, apiErr, got)

	// errors.Join produces a multi-error tree; only errors.As walks it.
	got, ok = AsAPIError(errors.Join(errors.New("other"), apiErr))
	assert.True(t, ok)
	assert.Same(t, apiErr, got)

	got, ok = AsAPIError(errors.New("plain"))
	assert.False(t, ok)
	assert.Nil(t, got)
}
