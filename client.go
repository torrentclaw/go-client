package torrentclaw

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

// Version is the library version, used in the default User-Agent header.
const Version = "0.2.0"

const (
	defaultBaseURL   = "https://torrentclaw.com"
	defaultTimeout   = 15 * time.Second
	defaultUserAgent = "torrentclaw-go-client/" + Version

	defaultMaxRetries    = 3
	defaultRetryBaseWait = 1 * time.Second
	defaultRetryMaxWait  = 30 * time.Second

	headerAPIKey         = "X-API-Key"
	headerAuthorization  = "Authorization"
	headerSearchSource   = "X-Search-Source"
	headerUserAgent      = "User-Agent"
	headerDebridProvider = "X-Debrid-Provider"
	headerDebridKey      = "X-Debrid-Key"

	searchSource = "go-client"
)

// Client is a TorrentClaw API client. Use [NewClient] to create one.
type Client struct {
	baseURL       string
	apiKey        string
	bearerToken   string
	userAgent     string
	httpClient    *http.Client
	maxRetries    int
	retryBaseWait time.Duration
	retryMaxWait  time.Duration
}

// Option configures a [Client].
type Option func(*Client)

// WithBaseURL sets a custom API base URL. The default is https://torrentclaw.com.
func WithBaseURL(u string) Option {
	return func(c *Client) { c.baseURL = u }
}

// WithAPIKey sets the API key sent as the X-API-Key header.
func WithAPIKey(key string) Option {
	return func(c *Client) { c.apiKey = key }
}

// WithBearerToken sets a bearer token sent as the Authorization header.
// If both WithBearerToken and WithAPIKey are used, the bearer token takes precedence.
func WithBearerToken(token string) Option {
	return func(c *Client) { c.bearerToken = token }
}

// WithUserAgent sets a custom User-Agent header.
func WithUserAgent(ua string) Option {
	return func(c *Client) { c.userAgent = ua }
}

// WithHTTPClient sets a custom *http.Client for all requests.
func WithHTTPClient(hc *http.Client) Option {
	return func(c *Client) { c.httpClient = hc }
}

// WithTimeout sets the HTTP client timeout. The default is 15 seconds.
// This option is safe to use regardless of option ordering; it sets the
// timeout on the client's internal HTTP client.
func WithTimeout(d time.Duration) Option {
	return func(c *Client) {
		if c.httpClient == nil {
			c.httpClient = &http.Client{}
		}
		c.httpClient.Timeout = d
	}
}

// WithRetry configures the retry policy for transient errors (429, 5xx).
// maxRetries is the maximum number of retries (0 disables retrying).
// baseWait is the initial wait duration before the first retry.
// maxWait caps the exponential backoff duration.
func WithRetry(maxRetries int, baseWait, maxWait time.Duration) Option {
	return func(c *Client) {
		c.maxRetries = maxRetries
		c.retryBaseWait = baseWait
		c.retryMaxWait = maxWait
	}
}

// NewClient creates a new TorrentClaw API client with the given options.
func NewClient(opts ...Option) *Client {
	c := &Client{
		baseURL:   defaultBaseURL,
		userAgent: defaultUserAgent,
		httpClient: &http.Client{
			Timeout: defaultTimeout,
		},
		maxRetries:    defaultMaxRetries,
		retryBaseWait: defaultRetryBaseWait,
		retryMaxWait:  defaultRetryMaxWait,
	}
	for _, opt := range opts {
		opt(c)
	}
	return c
}

// doJSON performs an HTTP GET request, retries on transient errors, and
// decodes the JSON response into dst.
func (c *Client) doJSON(ctx context.Context, path string, query url.Values, dst any) error {
	u, err := url.Parse(c.baseURL)
	if err != nil {
		return fmt.Errorf("torrentclaw: invalid base URL: %w", err)
	}
	u.Path = path
	if query != nil {
		u.RawQuery = query.Encode()
	}

	var lastErr error
	attempts := 1 + c.maxRetries
	for i := range attempts {
		if err := ctx.Err(); err != nil {
			return err
		}

		req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
		if err != nil {
			return fmt.Errorf("torrentclaw: failed to create request: %w", err)
		}
		c.setHeaders(req)
		req.Header.Set("Accept", "application/json")

		resp, err := c.httpClient.Do(req)
		if err != nil {
			return fmt.Errorf("torrentclaw: request failed: %w", err)
		}

		if resp.StatusCode == http.StatusOK {
			err := json.NewDecoder(resp.Body).Decode(dst)
			resp.Body.Close()
			if err != nil {
				return fmt.Errorf("torrentclaw: failed to decode response: %w", err)
			}
			return nil
		}

		body := readErrorBody(resp)
		resp.Body.Close()

		apiErr := newAPIError(resp.StatusCode, body)
		lastErr = apiErr

		if !apiErr.IsRetryable() || i == attempts-1 {
			return apiErr
		}

		wait := c.backoffDuration(i)
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(wait):
		}
	}
	return lastErr
}

// doRaw performs an HTTP GET request, retries on transient errors, and
// returns the raw response body bytes.
func (c *Client) doRaw(ctx context.Context, path string, query url.Values) ([]byte, error) {
	u, err := url.Parse(c.baseURL)
	if err != nil {
		return nil, fmt.Errorf("torrentclaw: invalid base URL: %w", err)
	}
	u.Path = path
	if query != nil {
		u.RawQuery = query.Encode()
	}

	var lastErr error
	attempts := 1 + c.maxRetries
	for i := range attempts {
		if err := ctx.Err(); err != nil {
			return nil, err
		}

		req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
		if err != nil {
			return nil, fmt.Errorf("torrentclaw: failed to create request: %w", err)
		}
		c.setHeaders(req)

		resp, err := c.httpClient.Do(req)
		if err != nil {
			return nil, fmt.Errorf("torrentclaw: request failed: %w", err)
		}

		if resp.StatusCode == http.StatusOK {
			data, err := io.ReadAll(resp.Body)
			resp.Body.Close()
			if err != nil {
				return nil, fmt.Errorf("torrentclaw: failed to read response body: %w", err)
			}
			return data, nil
		}

		body := readErrorBody(resp)
		resp.Body.Close()

		apiErr := newAPIError(resp.StatusCode, body)
		lastErr = apiErr

		if !apiErr.IsRetryable() || i == attempts-1 {
			return nil, apiErr
		}

		wait := c.backoffDuration(i)
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-time.After(wait):
		}
	}
	return nil, lastErr
}

// setHeaders applies common headers to an outgoing request.
func (c *Client) setHeaders(req *http.Request) {
	req.Header.Set(headerUserAgent, c.userAgent)
	req.Header.Set(headerSearchSource, searchSource)
	if c.bearerToken != "" {
		req.Header.Set(headerAuthorization, "Bearer "+c.bearerToken)
	} else if c.apiKey != "" {
		req.Header.Set(headerAPIKey, c.apiKey)
	}
}

// backoffDuration computes the wait time for retry attempt i using
// exponential backoff capped at retryMaxWait.
func (c *Client) backoffDuration(attempt int) time.Duration {
	wait := time.Duration(float64(c.retryBaseWait) * math.Pow(2, float64(attempt)))
	if wait > c.retryMaxWait {
		wait = c.retryMaxWait
	}
	return wait
}

// readErrorBody reads up to 512 bytes of the response body for error context.
// For server errors (5xx), an empty string is returned to avoid leaking internals.
func readErrorBody(resp *http.Response) string {
	if resp.StatusCode >= 500 {
		return ""
	}
	b, err := io.ReadAll(io.LimitReader(resp.Body, 512))
	if err != nil {
		return ""
	}
	return string(b)
}

// addIntParam adds an integer query parameter if the value is non-zero.
func addIntParam(q url.Values, key string, val int) {
	if val != 0 {
		q.Set(key, strconv.Itoa(val))
	}
}

// addFloatParam adds a float query parameter if the value is non-zero.
func addFloatParam(q url.Values, key string, val float64) {
	if val != 0 {
		q.Set(key, strconv.FormatFloat(val, 'f', -1, 64))
	}
}

// addStringParam adds a string query parameter if the value is non-empty.
func addStringParam(q url.Values, key, val string) {
	if val != "" {
		q.Set(key, val)
	}
}

// addBoolParam adds a boolean query parameter if the value is true.
func addBoolParam(q url.Values, key string, val bool) {
	if val {
		q.Set(key, "true")
	}
}

// doPost performs an HTTP POST request with a JSON body, retries on transient
// errors, and decodes the JSON response into dst. Extra headers (e.g. debrid
// provider credentials) are applied on top of the common headers.
func (c *Client) doPost(ctx context.Context, path string, body any, dst any, extraHeaders map[string]string) error {
	u, err := url.Parse(c.baseURL)
	if err != nil {
		return fmt.Errorf("torrentclaw: invalid base URL: %w", err)
	}
	u.Path = path

	payload, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("torrentclaw: failed to marshal request body: %w", err)
	}

	var lastErr error
	attempts := 1 + c.maxRetries
	for i := range attempts {
		if err := ctx.Err(); err != nil {
			return err
		}

		req, err := http.NewRequestWithContext(ctx, http.MethodPost, u.String(), bytes.NewReader(payload))
		if err != nil {
			return fmt.Errorf("torrentclaw: failed to create request: %w", err)
		}
		c.setHeaders(req)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Accept", "application/json")
		for k, v := range extraHeaders {
			req.Header.Set(k, v)
		}

		resp, err := c.httpClient.Do(req)
		if err != nil {
			return fmt.Errorf("torrentclaw: request failed: %w", err)
		}

		if resp.StatusCode == http.StatusOK {
			err := json.NewDecoder(resp.Body).Decode(dst)
			resp.Body.Close()
			if err != nil {
				return fmt.Errorf("torrentclaw: failed to decode response: %w", err)
			}
			return nil
		}

		errBody := readErrorBody(resp)
		resp.Body.Close()

		apiErr := newAPIError(resp.StatusCode, errBody)
		lastErr = apiErr

		if !apiErr.IsRetryable() || i == attempts-1 {
			return apiErr
		}

		wait := c.backoffDuration(i)
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(wait):
		}
	}
	return lastErr
}
