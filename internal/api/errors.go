package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"strings"
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
	var b strings.Builder
	b.WriteString("alegra: HTTP ")
	fmt.Fprintf(&b, "%d", e.StatusCode)
	if e.Message != "" {
		b.WriteString(": ")
		b.WriteString(e.Message)
	} else {
		fmt.Fprintf(&b, " (%s)", http.StatusText(e.StatusCode))
	}
	if e.Code != "" {
		fmt.Fprintf(&b, " (code %s)", e.Code)
	}
	if hint := e.Hint(); hint != "" {
		b.WriteString("\n  → ")
		b.WriteString(hint)
	}
	return b.String()
}

// stampCodeRe matches electronic-stamping provider error codes such as AEP1005,
// EPR503, or D1000 that Alegra surfaces when emission fails.
var stampCodeRe = regexp.MustCompile(`\b(AEP\d{3,5}|EPR\d{3}|D\d{4})\b`)

// stampHints maps the most common stamping error codes to a plain remedy.
var stampHints = map[string]string{
	"AEP1001": "the digital certificate is invalid — re-upload it in Alegra → Electronic invoicing settings",
	"AEP1003": "no digital certificate configured — upload one in Alegra settings",
	"AEP1005": "the digital certificate has expired — renew it in Alegra settings",
	"AEP6008": "this account isn't enabled to emit electronic invoices — complete DIAN habilitación",
	"AEP6009": "numbering rejected (duplicate/out of range) — use a valid resolution number",
	"AEP6011": "numbering rejected (duplicate/out of range) — use a valid resolution number",
	"AEP2011": "the document is already sent/processing — check its status before retrying",
	"EPR500":  "the tax authority (DIAN) communication failed — verify status, then retry",
	"EPR503":  "the tax authority (DIAN) timed out — verify status, then retry",
}

// Hint returns a short, actionable remediation string for the error, or "".
func (e *APIError) Hint() string {
	// Stamping codes can appear in the message or raw body.
	if code := stampCodeRe.FindString(e.Message + " " + e.Body); code != "" {
		if h, ok := stampHints[code]; ok {
			return fmt.Sprintf("%s (%s)", h, code)
		}
		return fmt.Sprintf("electronic stamping failed (%s) — check status before retrying; emission is not idempotent", code)
	}
	switch e.StatusCode {
	case http.StatusUnauthorized: // 401
		return "authentication failed — run `alegra auth login` or check ALEGRA_EMAIL/ALEGRA_TOKEN"
	case http.StatusPaymentRequired: // 402
		return "your Alegra plan doesn't include this action (or the account is suspended) — upgrade the plan"
	case http.StatusForbidden: // 403
		return "this user lacks permission for this action — check the user's role in Alegra"
	case http.StatusNotFound: // 404
		return "not found — verify the id with a `list` (or the account may be suspended)"
	case http.StatusTooManyRequests: // 429
		return "rate limit hit (150 req/min) — slow down; prefer --count over --all"
	case http.StatusInternalServerError, http.StatusServiceUnavailable: // 500/503
		return "Alegra had a server error — this is usually transient, try again shortly"
	}
	return ""
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
	if errors.As(err, &apiErr) {
		return apiErr, true
	}
	return nil, false
}
