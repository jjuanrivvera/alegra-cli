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
