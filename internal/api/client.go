// Package api implements a typed client for the Alegra accounting API
// (https://developer.alegra.com). It provides a generic Resource[T] abstraction
// so each accounting resource (contacts, invoices, items, ...) is a thin typed
// wrapper over a shared request/pagination/retry/rate-limit core.
package api

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

// DefaultBaseURL is the production Alegra API v1 root.
const DefaultBaseURL = "https://api.alegra.com/api/v1"

// MaxPageLimit is the maximum page size Alegra accepts for list endpoints.
const MaxPageLimit = 30

// errDryRun is returned internally when a request is short-circuited in
// dry-run mode. Callers treat it as a successful no-op.
var errDryRun = errors.New("dry-run: request not executed")

// HTTPDoer is the subset of *http.Client we depend on, to ease testing.
type HTTPDoer interface {
	Do(req *http.Request) (*http.Response, error)
}

// Client is a configured Alegra API client.
type Client struct {
	baseURL     string
	email       string
	token       string
	bearerToken string // optional: OAuth bearer (marketplace apps)

	httpClient  HTTPDoer
	rateLimiter *RateLimiter
	retryPolicy *RetryPolicy
	logger      *slog.Logger
	userAgent   string

	// rps is the configured rate-limiter rate; the limiter is built in New()
	// after all options are applied so it captures the final logger regardless
	// of option order (WithLogger may come after WithRequestsPerSecond).
	rps float64

	defaultLimit int

	dryRun       bool
	showToken    bool
	dryRunWriter io.Writer
}

// Option configures a Client.
type Option func(*Client)

// WithBaseURL overrides the API root (useful for tests or sandboxes).
func WithBaseURL(u string) Option {
	return func(c *Client) {
		if u != "" {
			c.baseURL = strings.TrimRight(u, "/")
		}
	}
}

// WithBasicAuth sets email + API token credentials (the standard Alegra auth).
func WithBasicAuth(email, token string) Option {
	return func(c *Client) { c.email, c.token = email, token }
}

// WithBearerToken sets an OAuth bearer token, used by marketplace apps.
func WithBearerToken(token string) Option {
	return func(c *Client) { c.bearerToken = token }
}

// WithHTTPClient injects a custom HTTP doer (tests use this).
func WithHTTPClient(h HTTPDoer) Option {
	return func(c *Client) {
		if h != nil {
			c.httpClient = h
		}
	}
}

// WithLogger sets the structured logger.
func WithLogger(l *slog.Logger) Option {
	return func(c *Client) {
		if l != nil {
			c.logger = l
		}
	}
}

// WithRequestsPerSecond configures the adaptive rate limiter.
func WithRequestsPerSecond(rps float64) Option {
	return func(c *Client) {
		if rps > 0 {
			c.rps = rps
		}
	}
}

// WithUserAgent overrides the User-Agent header.
func WithUserAgent(ua string) Option {
	return func(c *Client) {
		if ua != "" {
			c.userAgent = ua
		}
	}
}

// WithDryRun, when enabled, prints the equivalent curl command for each request
// to w (default os.Stdout) instead of executing it. If showToken is false the
// Authorization header is redacted.
func WithDryRun(enabled, showToken bool, w io.Writer) Option {
	return func(c *Client) {
		c.dryRun = enabled
		c.showToken = showToken
		if w != nil {
			c.dryRunWriter = w
		}
	}
}

// WithDefaultLimit sets the default page size for list operations.
func WithDefaultLimit(n int) Option {
	return func(c *Client) {
		if n > 0 && n <= MaxPageLimit {
			c.defaultLimit = n
		}
	}
}

// New builds a Client from options.
func New(opts ...Option) *Client {
	c := &Client{
		baseURL:      DefaultBaseURL,
		logger:       slog.Default(),
		userAgent:    "alegra-cli",
		defaultLimit: MaxPageLimit,
		dryRunWriter: io.Discard,
	}
	for _, o := range opts {
		o(c)
	}
	if c.httpClient == nil {
		c.httpClient = &http.Client{Timeout: 60 * time.Second}
	}
	if c.rateLimiter == nil {
		// NewRateLimiter applies its own 5.0 default when rps <= 0.
		c.rateLimiter = NewRateLimiter(c.rps, c.logger)
	}
	if c.retryPolicy == nil {
		c.retryPolicy = DefaultRetryPolicy(c.logger)
	}
	if c.dryRunWriter == nil {
		c.dryRunWriter = io.Discard
	}
	return c
}

// BaseURL returns the configured API root.
func (c *Client) BaseURL() string { return c.baseURL }

// RateLimit returns the most recent server-reported quota (limit, remaining,
// reset seconds); values are -1 until a request has been made.
func (c *Client) RateLimit() (limit, remaining, reset int) {
	return c.rateLimiter.Snapshot()
}

// DefaultLimit returns the default page size.
func (c *Client) DefaultLimit() int { return c.defaultLimit }

// buildURL joins the base URL, a resource path, and query parameters.
func (c *Client) buildURL(path string, query url.Values) string {
	path = strings.TrimLeft(path, "/")
	u := c.baseURL + "/" + path
	if len(query) > 0 {
		u += "?" + query.Encode()
	}
	return u
}

// newRequest constructs an authenticated *http.Request.
func (c *Client) newRequest(ctx context.Context, method, rawURL string, body []byte) (*http.Request, error) {
	var rdr io.Reader
	if body != nil {
		rdr = bytes.NewReader(body)
	}
	req, err := http.NewRequestWithContext(ctx, method, rawURL, rdr)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", c.userAgent)
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	switch {
	case c.bearerToken != "":
		req.Header.Set("Authorization", "Bearer "+c.bearerToken)
	case c.email != "" || c.token != "":
		req.SetBasicAuth(c.email, c.token)
	}
	return req, nil
}

// do executes a request with rate limiting and retries, returning the response
// body bytes on success. The caller owns interpreting the body.
func (c *Client) do(ctx context.Context, method, path string, query url.Values, body []byte) ([]byte, error) {
	rawURL := c.buildURL(path, query)

	if c.dryRun {
		c.printCurl(method, rawURL, body)
		return nil, errDryRun
	}

	for attempt := 0; attempt <= c.retryPolicy.MaxRetries; attempt++ {
		if err := c.rateLimiter.Wait(ctx); err != nil {
			return nil, err
		}

		req, err := c.newRequest(ctx, method, rawURL, body)
		if err != nil {
			return nil, err
		}

		c.logger.Debug("alegra request", "method", method, "url", rawURL, "attempt", attempt)
		resp, err := c.httpClient.Do(req)

		// Adapt to the server-reported quota on every response.
		if resp != nil {
			c.rateLimiter.Observe(
				headerInt(resp, "X-Rate-Limit-Limit"),
				headerInt(resp, "X-Rate-Limit-Remaining"),
				headerInt(resp, "X-Rate-Limit-Reset"),
			)
		}

		if c.retryPolicy.shouldRetry(method, resp, err) && attempt < c.retryPolicy.MaxRetries {
			if resp != nil && resp.StatusCode == http.StatusTooManyRequests {
				c.rateLimiter.Throttle()
			}
			wait := c.retryPolicy.backoff(attempt, resp)
			if resp != nil {
				_ = resp.Body.Close()
			}
			c.logger.Debug("alegra retry", "wait", wait.String(), "attempt", attempt)
			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			case <-time.After(wait):
			}
			continue
		}

		if err != nil {
			// http.Client.Do can return a non-nil response together with a
			// non-nil error (e.g. a redirect-policy failure), so close the body
			// before abandoning the request to avoid leaking the connection.
			if resp != nil {
				_ = resp.Body.Close()
			}
			return nil, fmt.Errorf("alegra: request failed: %w", err)
		}

		respBody, readErr := io.ReadAll(resp.Body)
		_ = resp.Body.Close()
		if readErr != nil {
			return nil, fmt.Errorf("alegra: reading response: %w", readErr)
		}

		if resp.StatusCode >= 400 {
			return nil, parseAPIError(resp.StatusCode, respBody)
		}

		c.rateLimiter.Restore()
		return respBody, nil
	}

	// Unreachable with a sane policy (MaxRetries >= 0): the final attempt always
	// returns inside the loop. Kept as a guard against a negative MaxRetries.
	return nil, errors.New("alegra: exhausted retries")
}

// doJSON sends an optional JSON body and decodes a JSON response into out.
// In dry-run mode it prints the request and returns nil without touching out.
func (c *Client) doJSON(ctx context.Context, method, path string, query url.Values, body, out any) error {
	var raw []byte
	if body != nil {
		switch b := body.(type) {
		case []byte:
			raw = b
		case json.RawMessage:
			raw = b
		default:
			marshaled, err := json.Marshal(body)
			if err != nil {
				return fmt.Errorf("alegra: marshaling request body: %w", err)
			}
			raw = marshaled
		}
	}

	respBody, err := c.do(ctx, method, path, query, raw)
	if err != nil {
		if errors.Is(err, errDryRun) {
			return nil
		}
		return err
	}
	if out == nil || len(respBody) == 0 {
		return nil
	}
	if err := json.Unmarshal(respBody, out); err != nil {
		return fmt.Errorf("alegra: decoding response: %w", err)
	}
	return nil
}

// printCurl writes a copy-pasteable curl command for the request.
func (c *Client) printCurl(method, rawURL string, body []byte) {
	var b strings.Builder
	b.WriteString("curl -X ")
	b.WriteString(method)
	b.WriteString(" \\\n  ")
	b.WriteString(shellQuote(rawURL))
	b.WriteString(" \\\n  -H 'Accept: application/json'")

	// --show-token reveals the auth scheme but never the live secret: the Basic
	// branch already emits a placeholder, and printing a real bearer token here
	// would leak a long-lived credential into terminal scrollback, CI logs, and
	// shell history. Both schemes are redacted to keep that boundary symmetric.
	auth := "<redacted>"
	if c.showToken {
		switch {
		case c.bearerToken != "":
			auth = "Bearer <token>"
		default:
			auth = "Basic <base64(email:token)>"
		}
	}
	b.WriteString(" \\\n  -H " + shellQuote("Authorization: "+auth))

	if body != nil {
		b.WriteString(" \\\n  -H 'Content-Type: application/json'")
		b.WriteString(" \\\n  -d " + shellQuote(string(body)))
	}
	b.WriteString("\n")
	_, _ = io.WriteString(c.dryRunWriter, b.String())
}

func shellQuote(s string) string {
	return "'" + strings.ReplaceAll(s, "'", `'\''`) + "'"
}

// headerInt parses an integer response header, returning -1 when absent/invalid.
func headerInt(resp *http.Response, name string) int {
	v := resp.Header.Get(name)
	if v == "" {
		return -1
	}
	n, err := strconv.Atoi(strings.TrimSpace(v))
	if err != nil {
		return -1
	}
	return n
}
