package torrentclaw

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestTorrentDownloadURL_CustomBase(t *testing.T) {
	c := NewClient(WithBaseURL("https://custom.example.com"))
	got := c.TorrentDownloadURL("deadbeef")
	want := "https://custom.example.com/api/v1/torrent/deadbeef"
	if got != want {
		t.Errorf("TorrentDownloadURL = %q, want %q", got, want)
	}
}

func TestGetTorrentFile_ServerError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer srv.Close()

	c := NewClient(WithBaseURL(srv.URL), WithRetry(0, 0, 0))
	_, err := c.GetTorrentFile(context.Background(), "abc123")
	if err == nil {
		t.Fatal("expected error")
	}
	apiErr, ok := err.(*APIError)
	if !ok {
		t.Fatalf("expected *APIError, got %T", err)
	}
	if apiErr.StatusCode != 500 {
		t.Errorf("StatusCode = %d, want 500", apiErr.StatusCode)
	}
}

func TestGetTorrentFile_RetryOn503(t *testing.T) {
	attempts := 0
	torrentData := []byte{0xDE, 0xAD, 0xBE, 0xEF}
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		attempts++
		if attempts < 3 {
			w.WriteHeader(http.StatusServiceUnavailable)
			return
		}
		w.Header().Set("Content-Type", "application/x-bittorrent")
		w.Write(torrentData)
	}))
	defer srv.Close()

	c := NewClient(
		WithBaseURL(srv.URL),
		WithRetry(3, 1*time.Millisecond, 10*time.Millisecond),
	)
	data, err := c.GetTorrentFile(context.Background(), "abc123")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(data) != len(torrentData) {
		t.Errorf("len(data) = %d, want %d", len(data), len(torrentData))
	}
	if attempts != 3 {
		t.Errorf("attempts = %d, want 3", attempts)
	}
}

func TestGetTorrentFile_RetryExhausted(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadGateway)
	}))
	defer srv.Close()

	c := NewClient(
		WithBaseURL(srv.URL),
		WithRetry(1, 1*time.Millisecond, 10*time.Millisecond),
	)
	_, err := c.GetTorrentFile(context.Background(), "abc123")
	if err == nil {
		t.Fatal("expected error after retries exhausted")
	}
	apiErr, ok := err.(*APIError)
	if !ok {
		t.Fatalf("expected *APIError, got %T", err)
	}
	if apiErr.StatusCode != 502 {
		t.Errorf("StatusCode = %d, want 502", apiErr.StatusCode)
	}
}

func TestGetTorrentFile_ContextCanceled(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(1 * time.Second)
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	c := NewClient(WithBaseURL(srv.URL), WithRetry(0, 0, 0))
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	_, err := c.GetTorrentFile(ctx, "abc123")
	if err == nil {
		t.Fatal("expected error for canceled context")
	}
}
