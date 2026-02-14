package torrentclaw

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"
)

func TestNewClientDefaults(t *testing.T) {
	c := NewClient()
	if c.baseURL != defaultBaseURL {
		t.Errorf("baseURL = %q, want %q", c.baseURL, defaultBaseURL)
	}
	if c.userAgent != defaultUserAgent {
		t.Errorf("userAgent = %q, want %q", c.userAgent, defaultUserAgent)
	}
	if c.apiKey != "" {
		t.Errorf("apiKey = %q, want empty", c.apiKey)
	}
	if c.maxRetries != defaultMaxRetries {
		t.Errorf("maxRetries = %d, want %d", c.maxRetries, defaultMaxRetries)
	}
}

func TestNewClientOptions(t *testing.T) {
	hc := &http.Client{Timeout: 5 * time.Second}
	c := NewClient(
		WithBaseURL("https://custom.example.com"),
		WithAPIKey("test-key-123"),
		WithUserAgent("my-app/1.0"),
		WithHTTPClient(hc),
		WithRetry(5, 2*time.Second, 60*time.Second),
	)
	if c.baseURL != "https://custom.example.com" {
		t.Errorf("baseURL = %q", c.baseURL)
	}
	if c.apiKey != "test-key-123" {
		t.Errorf("apiKey = %q", c.apiKey)
	}
	if c.userAgent != "my-app/1.0" {
		t.Errorf("userAgent = %q", c.userAgent)
	}
	if c.httpClient != hc {
		t.Error("httpClient not set")
	}
	if c.maxRetries != 5 {
		t.Errorf("maxRetries = %d, want 5", c.maxRetries)
	}
}

func TestWithTimeout(t *testing.T) {
	c := NewClient(WithTimeout(30 * time.Second))
	if c.httpClient.Timeout != 30*time.Second {
		t.Errorf("timeout = %v, want 30s", c.httpClient.Timeout)
	}
}

func TestSetHeaders(t *testing.T) {
	tests := []struct {
		name   string
		apiKey string
		want   map[string]string
	}{
		{
			name:   "without API key",
			apiKey: "",
			want: map[string]string{
				headerUserAgent:    defaultUserAgent,
				headerSearchSource: searchSource,
			},
		},
		{
			name:   "with API key",
			apiKey: "my-key",
			want: map[string]string{
				headerUserAgent:    defaultUserAgent,
				headerSearchSource: searchSource,
				headerAPIKey:       "my-key",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := NewClient(WithAPIKey(tt.apiKey))
			req, _ := http.NewRequest(http.MethodGet, "https://example.com", nil)
			c.setHeaders(req)
			for key, want := range tt.want {
				got := req.Header.Get(key)
				if got != want {
					t.Errorf("header %q = %q, want %q", key, got, want)
				}
			}
			if tt.apiKey == "" && req.Header.Get(headerAPIKey) != "" {
				t.Error("API key header should not be set when apiKey is empty")
			}
		})
	}
}

func TestDoJSON_Success(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got := r.Header.Get(headerUserAgent); got != defaultUserAgent {
			t.Errorf("User-Agent = %q, want %q", got, defaultUserAgent)
		}
		if got := r.Header.Get(headerSearchSource); got != searchSource {
			t.Errorf("X-Search-Source = %q, want %q", got, searchSource)
		}
		if got := r.Header.Get("Accept"); got != "application/json" {
			t.Errorf("Accept = %q, want application/json", got)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]int{"total": 42})
	}))
	defer srv.Close()

	c := NewClient(WithBaseURL(srv.URL))
	var dst struct{ Total int }
	err := c.doJSON(context.Background(), "/test", nil, &dst)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if dst.Total != 42 {
		t.Errorf("Total = %d, want 42", dst.Total)
	}
}

func TestDoJSON_APIError(t *testing.T) {
	tests := []struct {
		name       string
		statusCode int
		body       string
		wantBody   bool
	}{
		{name: "400 bad request", statusCode: 400, body: "invalid query", wantBody: true},
		{name: "404 not found", statusCode: 404, body: "not found", wantBody: true},
		{name: "500 server error", statusCode: 500, body: "internal details", wantBody: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tt.statusCode)
				w.Write([]byte(tt.body))
			}))
			defer srv.Close()

			c := NewClient(WithBaseURL(srv.URL), WithRetry(0, 0, 0))
			var dst struct{}
			err := c.doJSON(context.Background(), "/test", nil, &dst)
			if err == nil {
				t.Fatal("expected error")
			}
			apiErr, ok := err.(*APIError)
			if !ok {
				t.Fatalf("expected *APIError, got %T", err)
			}
			if apiErr.StatusCode != tt.statusCode {
				t.Errorf("StatusCode = %d, want %d", apiErr.StatusCode, tt.statusCode)
			}
			if tt.wantBody && apiErr.Body == "" {
				t.Error("expected Body to be non-empty for 4xx")
			}
			if !tt.wantBody && apiErr.Body != "" {
				t.Errorf("expected empty Body for 5xx, got %q", apiErr.Body)
			}
		})
	}
}

func TestDoJSON_RetryOn429(t *testing.T) {
	attempts := 0
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		attempts++
		if attempts < 3 {
			w.WriteHeader(http.StatusTooManyRequests)
			w.Write([]byte("rate limited"))
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
	}))
	defer srv.Close()

	c := NewClient(
		WithBaseURL(srv.URL),
		WithRetry(3, 1*time.Millisecond, 10*time.Millisecond),
	)
	var dst struct{ Status string }
	err := c.doJSON(context.Background(), "/test", nil, &dst)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if dst.Status != "ok" {
		t.Errorf("Status = %q, want ok", dst.Status)
	}
	if attempts != 3 {
		t.Errorf("attempts = %d, want 3", attempts)
	}
}

func TestDoJSON_RetryExhausted(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusTooManyRequests)
		w.Write([]byte("rate limited"))
	}))
	defer srv.Close()

	c := NewClient(
		WithBaseURL(srv.URL),
		WithRetry(2, 1*time.Millisecond, 10*time.Millisecond),
	)
	var dst struct{}
	err := c.doJSON(context.Background(), "/test", nil, &dst)
	if err == nil {
		t.Fatal("expected error after retries exhausted")
	}
	apiErr, ok := err.(*APIError)
	if !ok {
		t.Fatalf("expected *APIError, got %T", err)
	}
	if apiErr.StatusCode != 429 {
		t.Errorf("StatusCode = %d, want 429", apiErr.StatusCode)
	}
}

func TestDoJSON_ContextCanceled(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(1 * time.Second)
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	c := NewClient(WithBaseURL(srv.URL), WithRetry(0, 0, 0))
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	var dst struct{}
	err := c.doJSON(ctx, "/test", nil, &dst)
	if err == nil {
		t.Fatal("expected error for canceled context")
	}
}

func TestBackoffDuration(t *testing.T) {
	c := NewClient(WithRetry(5, 1*time.Second, 30*time.Second))
	tests := []struct {
		attempt int
		want    time.Duration
	}{
		{0, 1 * time.Second},
		{1, 2 * time.Second},
		{2, 4 * time.Second},
		{3, 8 * time.Second},
		{4, 16 * time.Second},
		{5, 30 * time.Second}, // capped
	}
	for _, tt := range tests {
		got := c.backoffDuration(tt.attempt)
		if got != tt.want {
			t.Errorf("backoffDuration(%d) = %v, want %v", tt.attempt, got, tt.want)
		}
	}
}

func TestAPIKeyHeader(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		got := r.Header.Get(headerAPIKey)
		if got != "secret-key" {
			t.Errorf("X-API-Key = %q, want %q", got, "secret-key")
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]bool{"ok": true})
	}))
	defer srv.Close()

	c := NewClient(WithBaseURL(srv.URL), WithAPIKey("secret-key"))
	var dst struct{ OK bool }
	err := c.doJSON(context.Background(), "/test", nil, &dst)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestDoJSON_InvalidJSON(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte("not json"))
	}))
	defer srv.Close()

	c := NewClient(WithBaseURL(srv.URL), WithRetry(0, 0, 0))
	var dst struct{}
	err := c.doJSON(context.Background(), "/test", nil, &dst)
	if err == nil {
		t.Fatal("expected error for invalid JSON")
	}
	if !strings.Contains(err.Error(), "decode") {
		t.Errorf("error = %q, want to contain 'decode'", err.Error())
	}
}

func TestDoJSON_ContextCanceledDuringBackoff(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusTooManyRequests)
		w.Write([]byte("rate limited"))
	}))
	defer srv.Close()

	ctx, cancel := context.WithCancel(context.Background())
	c := NewClient(
		WithBaseURL(srv.URL),
		WithRetry(5, 1*time.Second, 10*time.Second), // long backoff
	)

	go func() {
		time.Sleep(50 * time.Millisecond)
		cancel()
	}()

	var dst struct{}
	err := c.doJSON(ctx, "/test", nil, &dst)
	if err == nil {
		t.Fatal("expected error")
	}
	if err != context.Canceled {
		t.Errorf("err = %v, want context.Canceled", err)
	}
}

func TestDoJSON_RetryOn502(t *testing.T) {
	attempts := 0
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		attempts++
		if attempts < 2 {
			w.WriteHeader(http.StatusBadGateway)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"ok": "true"})
	}))
	defer srv.Close()

	c := NewClient(
		WithBaseURL(srv.URL),
		WithRetry(2, 1*time.Millisecond, 10*time.Millisecond),
	)
	var dst struct{ Ok string }
	err := c.doJSON(context.Background(), "/test", nil, &dst)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if attempts != 2 {
		t.Errorf("attempts = %d, want 2", attempts)
	}
}

func TestDoJSON_RetryOn503(t *testing.T) {
	attempts := 0
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		attempts++
		if attempts < 2 {
			w.WriteHeader(http.StatusServiceUnavailable)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"ok": "true"})
	}))
	defer srv.Close()

	c := NewClient(
		WithBaseURL(srv.URL),
		WithRetry(2, 1*time.Millisecond, 10*time.Millisecond),
	)
	var dst struct{ Ok string }
	err := c.doJSON(context.Background(), "/test", nil, &dst)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if attempts != 2 {
		t.Errorf("attempts = %d, want 2", attempts)
	}
}

func TestDoJSON_NoRetryOn401(t *testing.T) {
	attempts := 0
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		attempts++
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("unauthorized"))
	}))
	defer srv.Close()

	c := NewClient(
		WithBaseURL(srv.URL),
		WithRetry(3, 1*time.Millisecond, 10*time.Millisecond),
	)
	var dst struct{}
	err := c.doJSON(context.Background(), "/test", nil, &dst)
	if err == nil {
		t.Fatal("expected error")
	}
	if attempts != 1 {
		t.Errorf("attempts = %d, want 1 (no retries for 401)", attempts)
	}
}

func TestDoJSON_WithQueryParams(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got := r.URL.Query().Get("foo"); got != "bar" {
			t.Errorf("foo = %q, want bar", got)
		}
		if got := r.URL.Query().Get("num"); got != "42" {
			t.Errorf("num = %q, want 42", got)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]bool{"ok": true})
	}))
	defer srv.Close()

	c := NewClient(WithBaseURL(srv.URL), WithRetry(0, 0, 0))
	q := url.Values{}
	q.Set("foo", "bar")
	q.Set("num", "42")
	var dst struct{ Ok bool }
	err := c.doJSON(context.Background(), "/test", q, &dst)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestDoJSON_InvalidBaseURL(t *testing.T) {
	c := NewClient(WithBaseURL("://invalid"), WithRetry(0, 0, 0))
	var dst struct{}
	err := c.doJSON(context.Background(), "/test", nil, &dst)
	if err == nil {
		t.Fatal("expected error for invalid base URL")
	}
	if !strings.Contains(err.Error(), "invalid base URL") {
		t.Errorf("error = %q, want to contain 'invalid base URL'", err.Error())
	}
}

func TestDoRaw_InvalidBaseURL(t *testing.T) {
	c := NewClient(WithBaseURL("://invalid"), WithRetry(0, 0, 0))
	_, err := c.doRaw(context.Background(), "/test")
	if err == nil {
		t.Fatal("expected error for invalid base URL")
	}
	if !strings.Contains(err.Error(), "invalid base URL") {
		t.Errorf("error = %q, want to contain 'invalid base URL'", err.Error())
	}
}

func TestAddIntParam(t *testing.T) {
	q := url.Values{}
	addIntParam(q, "zero", 0)
	if q.Has("zero") {
		t.Error("zero value should not be added")
	}
	addIntParam(q, "val", 42)
	if q.Get("val") != "42" {
		t.Errorf("val = %q, want 42", q.Get("val"))
	}
}

func TestAddFloatParam(t *testing.T) {
	q := url.Values{}
	addFloatParam(q, "zero", 0)
	if q.Has("zero") {
		t.Error("zero value should not be added")
	}
	addFloatParam(q, "rating", 7.5)
	if q.Get("rating") != "7.5" {
		t.Errorf("rating = %q, want 7.5", q.Get("rating"))
	}
}

func TestAddStringParam(t *testing.T) {
	q := url.Values{}
	addStringParam(q, "empty", "")
	if q.Has("empty") {
		t.Error("empty value should not be added")
	}
	addStringParam(q, "key", "value")
	if q.Get("key") != "value" {
		t.Errorf("key = %q, want value", q.Get("key"))
	}
}

func TestVersion(t *testing.T) {
	if Version == "" {
		t.Error("Version should not be empty")
	}
}

func TestDoRaw_ContextCanceledDuringBackoff(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadGateway)
	}))
	defer srv.Close()

	ctx, cancel := context.WithCancel(context.Background())
	c := NewClient(
		WithBaseURL(srv.URL),
		WithRetry(5, 1*time.Second, 10*time.Second),
	)

	go func() {
		time.Sleep(50 * time.Millisecond)
		cancel()
	}()

	_, err := c.doRaw(ctx, "/test")
	if err == nil {
		t.Fatal("expected error")
	}
	if err != context.Canceled {
		t.Errorf("err = %v, want context.Canceled", err)
	}
}

func TestWithTimeout_NilHTTPClient(t *testing.T) {
	// Simulate the edge case where httpClient is nil when WithTimeout is applied.
	c := &Client{}
	opt := WithTimeout(10 * time.Second)
	opt(c)
	if c.httpClient == nil {
		t.Fatal("httpClient should be initialized")
	}
	if c.httpClient.Timeout != 10*time.Second {
		t.Errorf("timeout = %v, want 10s", c.httpClient.Timeout)
	}
}
