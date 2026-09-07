package torrentclaw

import "fmt"

// APIError represents an error response from the TorrentClaw API.
type APIError struct {
	// StatusCode is the HTTP status code returned by the API.
	StatusCode int

	// Body contains the response body for client errors (4xx).
	// It is intentionally omitted for server errors (5xx) to avoid leaking internals.
	Body string

	// Message is a human-readable error description.
	Message string
}

// Error implements the error interface.
func (e *APIError) Error() string {
	if e.Body != "" {
		return fmt.Sprintf("torrentclaw: %s (HTTP %d): %s", e.Message, e.StatusCode, e.Body)
	}
	return fmt.Sprintf("torrentclaw: %s (HTTP %d)", e.Message, e.StatusCode)
}

// IsRetryable reports whether the error is likely transient and the request
// should be retried.
// 424 is here because the API reports an upstream DEBRID PROVIDER failure
// (TorBox/Real-Debrid down, unreachable, or unparseable) as 424 Failed
// Dependency rather than 502. That is a transient condition and deserves the
// same retry a 502 used to get; without this entry the switch to 424 would
// silently turn a retried hiccup into a hard error for the caller.
//
// The API changed because a 5xx means TorrentClaw itself failed, and emitting
// one for a third party's outage made the reverse proxy count the replica that
// answered as dead — ~4.4k/day from check-cache alone, which is what emptied the
// upstream pool and 502'd the whole site several times a day (Sept 2026).
func (e *APIError) IsRetryable() bool {
	switch e.StatusCode {
	case 424, 429, 500, 502, 503:
		return true
	default:
		return false
	}
}

// IsRateLimited reports whether the error is due to rate limiting (HTTP 429).
func (e *APIError) IsRateLimited() bool {
	return e.StatusCode == 429
}

// IsNotFound reports whether the error is due to a missing resource (HTTP 404).
func (e *APIError) IsNotFound() bool {
	return e.StatusCode == 404
}

// statusMessage returns a human-readable message for known HTTP status codes.
func statusMessage(code int) string {
	switch code {
	case 400:
		return "Bad request — check that all parameters are valid"
	case 401:
		return "Authentication required — provide a valid API key"
	case 403:
		return "Forbidden — insufficient permissions for this endpoint"
	case 404:
		return "Not found — the requested resource does not exist"
	case 424:
		return "The debrid provider is unavailable — not a TorrentClaw fault; retry shortly"
	case 429:
		return "Rate limit exceeded — wait before retrying"
	case 500:
		return "TorrentClaw server error"
	case 502:
		return "TorrentClaw is temporarily unavailable"
	case 503:
		return "TorrentClaw is under maintenance"
	default:
		return fmt.Sprintf("API request failed with status %d", code)
	}
}

// newAPIError creates an APIError from an HTTP status code and response body.
func newAPIError(statusCode int, body string) *APIError {
	return &APIError{
		StatusCode: statusCode,
		Body:       body,
		Message:    statusMessage(statusCode),
	}
}
