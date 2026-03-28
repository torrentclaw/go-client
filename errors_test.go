package torrentclaw

import (
	"errors"
	"testing"
)

func TestAPIError_Error(t *testing.T) {
	tests := []struct {
		name string
		err  *APIError
		want string
	}{
		{
			name: "with body",
			err:  &APIError{StatusCode: 400, Message: "Bad request", Body: "invalid query"},
			want: "torrentclaw: Bad request (HTTP 400): invalid query",
		},
		{
			name: "without body",
			err:  &APIError{StatusCode: 500, Message: "Server error", Body: ""},
			want: "torrentclaw: Server error (HTTP 500)",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.err.Error()
			if got != tt.want {
				t.Errorf("Error() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestAPIError_IsRetryable(t *testing.T) {
	tests := []struct {
		code int
		want bool
	}{
		{400, false},
		{401, false},
		{403, false},
		{404, false},
		{429, true},
		{500, true},
		{502, true},
		{503, true},
		{504, false},
	}
	for _, tt := range tests {
		err := &APIError{StatusCode: tt.code}
		if got := err.IsRetryable(); got != tt.want {
			t.Errorf("IsRetryable() for %d = %v, want %v", tt.code, got, tt.want)
		}
	}
}

func TestAPIError_IsRateLimited(t *testing.T) {
	tests := []struct {
		code int
		want bool
	}{
		{200, false},
		{400, false},
		{429, true},
		{500, false},
	}
	for _, tt := range tests {
		err := &APIError{StatusCode: tt.code}
		if got := err.IsRateLimited(); got != tt.want {
			t.Errorf("IsRateLimited() for %d = %v, want %v", tt.code, got, tt.want)
		}
	}
}

func TestAPIError_IsNotFound(t *testing.T) {
	tests := []struct {
		code int
		want bool
	}{
		{200, false},
		{400, false},
		{404, true},
		{500, false},
	}
	for _, tt := range tests {
		err := &APIError{StatusCode: tt.code}
		if got := err.IsNotFound(); got != tt.want {
			t.Errorf("IsNotFound() for %d = %v, want %v", tt.code, got, tt.want)
		}
	}
}

func TestStatusMessage(t *testing.T) {
	tests := []struct {
		code int
		want string
	}{
		{400, "Bad request — check that all parameters are valid"},
		{401, "Authentication required — provide a valid API key"},
		{403, "Forbidden — insufficient permissions for this endpoint"},
		{404, "Not found — the requested resource does not exist"},
		{429, "Rate limit exceeded — wait before retrying"},
		{500, "TorrentClaw server error"},
		{502, "TorrentClaw is temporarily unavailable"},
		{503, "TorrentClaw is under maintenance"},
		{418, "API request failed with status 418"},
	}
	for _, tt := range tests {
		got := statusMessage(tt.code)
		if got != tt.want {
			t.Errorf("statusMessage(%d) = %q, want %q", tt.code, got, tt.want)
		}
	}
}

func TestNewAPIError(t *testing.T) {
	err := newAPIError(404, "resource not found")
	if err.StatusCode != 404 {
		t.Errorf("StatusCode = %d, want 404", err.StatusCode)
	}
	if err.Body != "resource not found" {
		t.Errorf("Body = %q, want %q", err.Body, "resource not found")
	}
	if err.Message != statusMessage(404) {
		t.Errorf("Message = %q", err.Message)
	}
}

func TestAPIError_ErrorsAs(t *testing.T) {
	err := newAPIError(429, "rate limited")
	var apiErr *APIError
	if !errors.As(err, &apiErr) {
		t.Fatal("errors.As failed to match *APIError")
	}
	if apiErr.StatusCode != 429 {
		t.Errorf("StatusCode = %d, want 429", apiErr.StatusCode)
	}
}
