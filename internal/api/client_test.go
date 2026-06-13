package api

import (
	"bytes"
	"context"
	"net/http"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNew_Defaults(t *testing.T) {
	c := New()
	assert.Equal(t, DefaultBaseURL, c.BaseURL())
	assert.Equal(t, MaxPageLimit, c.DefaultLimit())
	limit, remaining, reset := c.RateLimit()
	assert.Equal(t, -1, limit)
	assert.Equal(t, -1, remaining)
	assert.Equal(t, -1, reset)
}

func TestOptions(t *testing.T) {
	c := New(
		WithBaseURL("https://x.example/api/v1/"), // trailing slash trimmed
		WithDefaultLimit(10),
		WithUserAgent("ua/1"),
		WithBearerToken("tok"),
	)
	assert.Equal(t, "https://x.example/api/v1", c.BaseURL())
	assert.Equal(t, 10, c.DefaultLimit())
	assert.Equal(t, "ua/1", c.userAgent)
	assert.Equal(t, "tok", c.bearerToken)

	// Invalid / ignored options keep the defaults.
	assert.Equal(t, MaxPageLimit, New(WithDefaultLimit(0)).DefaultLimit())
	assert.Equal(t, MaxPageLimit, New(WithDefaultLimit(99)).DefaultLimit())
	assert.Equal(t, DefaultBaseURL, New(WithBaseURL("")).BaseURL())
	assert.NotNil(t, New(WithHTTPClient(nil)).httpClient, "nil doer ignored")
	assert.NotNil(t, New(WithLogger(nil)).logger, "nil logger ignored")
}

func TestBuildURL(t *testing.T) {
	c := New(WithBaseURL("https://x.example/api/v1"))
	assert.Equal(t, "https://x.example/api/v1/contacts", c.buildURL("contacts", nil))
	assert.Equal(t, "https://x.example/api/v1/contacts", c.buildURL("/contacts", nil))

	q := url.Values{}
	q.Set("start", "0")
	assert.Equal(t, "https://x.example/api/v1/contacts?start=0", c.buildURL("contacts", q))
}

func TestHeaderInt(t *testing.T) {
	resp := &http.Response{Header: http.Header{}}
	resp.Header.Set("X-Rate-Limit-Limit", "150")
	resp.Header.Set("Bad", "abc")
	assert.Equal(t, 150, headerInt(resp, "X-Rate-Limit-Limit"))
	assert.Equal(t, -1, headerInt(resp, "Missing"))
	assert.Equal(t, -1, headerInt(resp, "Bad"))
}

func TestShellQuote(t *testing.T) {
	assert.Equal(t, "'abc'", shellQuote("abc"))
	assert.Equal(t, `'a'\''b'`, shellQuote("a'b"))
}

func TestDryRunPrintsCurl_Redacted(t *testing.T) {
	buf := &bytes.Buffer{}
	c := New(WithBaseURL("https://x.example/api/v1"), WithBasicAuth("e@x.com", "secret"), WithDryRun(true, false, buf))

	_, err := c.do(context.Background(), http.MethodGet, "contacts", nil, nil)
	require.ErrorIs(t, err, errDryRun)

	out := buf.String()
	assert.Contains(t, out, "curl -X GET")
	assert.Contains(t, out, "https://x.example/api/v1/contacts")
	assert.Contains(t, out, "<redacted>")
	assert.NotContains(t, out, "secret")
}

func TestDryRunPrintsCurl_ShowTokenAndBody(t *testing.T) {
	buf := &bytes.Buffer{}
	c := New(WithBearerToken("brr"), WithDryRun(true, true, buf))

	_, err := c.do(context.Background(), http.MethodPost, "invoices", nil, []byte(`{"a":1}`))
	require.ErrorIs(t, err, errDryRun)

	out := buf.String()
	assert.Contains(t, out, "curl -X POST")
	// --show-token reveals the auth scheme but never the live secret (L11): the
	// real bearer token must not leak into the curl preview.
	assert.Contains(t, out, "Authorization: Bearer <token>")
	assert.NotContains(t, out, "brr")
	assert.Contains(t, out, "Content-Type: application/json")
	assert.Contains(t, out, `-d '{"a":1}'`)
}
