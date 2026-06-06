package api

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAPIError_Hint(t *testing.T) {
	cases := []struct {
		name   string
		err    *APIError
		expect string
	}{
		{"plan", &APIError{StatusCode: 402}, "plan doesn't include"},
		{"forbidden", &APIError{StatusCode: 403}, "lacks permission"},
		{"rate", &APIError{StatusCode: 429}, "rate limit"},
		{"auth", &APIError{StatusCode: 401}, "authentication failed"},
		{"notfound", &APIError{StatusCode: 404}, "not found"},
		{"server", &APIError{StatusCode: 500}, "transient"},
		{"stamp-known", &APIError{StatusCode: 400, Message: "rejected AEP1005"}, "certificate has expired"},
		{"stamp-unknown", &APIError{StatusCode: 400, Body: `{"code":"AEP9999"}`}, "not idempotent"},
		{"none", &APIError{StatusCode: 400, Message: "bad field"}, ""},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := tc.err.Hint()
			if tc.expect == "" {
				assert.Empty(t, got)
				return
			}
			assert.Contains(t, got, tc.expect)
		})
	}
}

func TestAPIError_ErrorIncludesHint(t *testing.T) {
	e := &APIError{StatusCode: 402, Message: "Plan suspended"}
	s := e.Error()
	assert.True(t, strings.Contains(s, "HTTP 402"))
	assert.True(t, strings.Contains(s, "→"))
	assert.Contains(t, s, "plan doesn't include")
}

func TestParseAPIError_Shapes(t *testing.T) {
	// bare object
	e := parseAPIError(422, []byte(`{"message":"Validation failed","code":"1001"}`))
	assert.Equal(t, "Validation failed", e.Message)
	assert.Equal(t, "1001", e.Code)
	// wrapped
	e = parseAPIError(400, []byte(`{"error":{"message":"Bad"}}`))
	assert.Equal(t, "Bad", e.Message)
	// unparseable → still carries status + body
	e = parseAPIError(500, []byte(`oops`))
	assert.Equal(t, 500, e.StatusCode)
	assert.Equal(t, "oops", e.Body)
}
