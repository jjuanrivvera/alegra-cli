package api

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// APIError represents an error returned by the Alegra API.
//
// Alegra returns errors in a few shapes depending on the endpoint. The most
// common is a JSON object with a top-level "message" and optional "code" and
// "details". We capture the raw body so callers can inspect anything we did
// not model explicitly.
type APIError struct {
	StatusCode int    `json:"-"`
	Message    string `json:"message"`
	Code       string `json:"code,omitempty"`
	// Details holds field-level validation errors keyed by field name when the
	// API provides them.
	Details map[string]any `json:"details,omitempty"`
	// Body is the raw, unparsed response body, useful for debugging.
	Body string `json:"-"`
}

func (e *APIError) Error() string {
	if e.Message == "" {
		return fmt.Sprintf("alegra: HTTP %d (%s)", e.StatusCode, http.StatusText(e.StatusCode))
	}
	if e.Code != "" {
		return fmt.Sprintf("alegra: HTTP %d: %s (code %s)", e.StatusCode, e.Message, e.Code)
	}
	return fmt.Sprintf("alegra: HTTP %d: %s", e.StatusCode, e.Message)
}

// IsNotFound reports whether the error is a 404.
func (e *APIError) IsNotFound() bool { return e.StatusCode == http.StatusNotFound }

// IsUnauthorized reports whether the error is a 401/403 (bad credentials/scope).
func (e *APIError) IsUnauthorized() bool {
	return e.StatusCode == http.StatusUnauthorized || e.StatusCode == http.StatusForbidden
}

// IsRateLimited reports whether the error is a 429.
func (e *APIError) IsRateLimited() bool { return e.StatusCode == http.StatusTooManyRequests }

// parseAPIError builds an APIError from an HTTP response body.
func parseAPIError(statusCode int, body []byte) *APIError {
	apiErr := &APIError{StatusCode: statusCode, Body: string(body)}

	// Alegra sometimes wraps the error object, sometimes returns it bare.
	// Try the bare shape first, then a {"error": {...}} wrapper.
	if err := json.Unmarshal(body, apiErr); err == nil && apiErr.Message != "" {
		return apiErr
	}

	var wrapper struct {
		Error APIError `json:"error"`
	}
	if err := json.Unmarshal(body, &wrapper); err == nil && wrapper.Error.Message != "" {
		wrapper.Error.StatusCode = statusCode
		wrapper.Error.Body = string(body)
		return &wrapper.Error
	}

	return apiErr
}

// AsAPIError extracts an *APIError from an error chain, if present.
func AsAPIError(err error) (*APIError, bool) {
	var apiErr *APIError
	for err != nil {
		if e, ok := err.(*APIError); ok {
			return e, true
		}
		type unwrapper interface{ Unwrap() error }
		u, ok := err.(unwrapper)
		if !ok {
			break
		}
		err = u.Unwrap()
	}
	return apiErr, false
}
