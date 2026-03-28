package torrentclaw

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestDebridCheckCache(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("method = %q, want POST", r.Method)
		}
		if r.URL.Path != "/api/v1/debrid/check-cache" {
			t.Errorf("path = %q", r.URL.Path)
		}
		if got := r.Header.Get(headerDebridProvider); got != "real-debrid" {
			t.Errorf("X-Debrid-Provider = %q, want real-debrid", got)
		}
		if got := r.Header.Get(headerDebridKey); got != "my-debrid-key" {
			t.Errorf("X-Debrid-Key = %q, want my-debrid-key", got)
		}

		body, _ := io.ReadAll(r.Body)
		var req DebridCheckCacheRequest
		if err := json.Unmarshal(body, &req); err != nil {
			t.Fatalf("failed to parse request body: %v", err)
		}
		if len(req.InfoHashes) != 2 {
			t.Errorf("len(InfoHashes) = %d, want 2", len(req.InfoHashes))
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(DebridCheckCacheResponse{
			Cached: map[string]bool{
				"hash1": true,
				"hash2": false,
			},
		})
	}))
	defer srv.Close()

	c := NewClient(WithBaseURL(srv.URL), WithRetry(0, 0, 0))
	resp, err := c.DebridCheckCache(context.Background(), "real-debrid", "my-debrid-key", []string{"hash1", "hash2"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !resp.Cached["hash1"] {
		t.Error("hash1 should be cached")
	}
	if resp.Cached["hash2"] {
		t.Error("hash2 should not be cached")
	}
}

func TestDebridAddMagnet(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("method = %q, want POST", r.Method)
		}
		if r.URL.Path != "/api/v1/debrid/add-magnet" {
			t.Errorf("path = %q", r.URL.Path)
		}
		if got := r.Header.Get(headerDebridProvider); got != "alldebrid" {
			t.Errorf("X-Debrid-Provider = %q, want alldebrid", got)
		}

		body, _ := io.ReadAll(r.Body)
		var req DebridAddMagnetRequest
		if err := json.Unmarshal(body, &req); err != nil {
			t.Fatalf("failed to parse request body: %v", err)
		}
		if req.InfoHash != "abc123" {
			t.Errorf("InfoHash = %q, want abc123", req.InfoHash)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(DebridAddMagnetResponse{
			ID:     "torrent-id-1",
			Cached: true,
			Name:   "Test Movie",
		})
	}))
	defer srv.Close()

	c := NewClient(WithBaseURL(srv.URL), WithRetry(0, 0, 0))
	resp, err := c.DebridAddMagnet(context.Background(), "alldebrid", "key123", "abc123")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.ID != "torrent-id-1" {
		t.Errorf("ID = %q, want torrent-id-1", resp.ID)
	}
	if !resp.Cached {
		t.Error("Cached should be true")
	}
	if resp.Name != "Test Movie" {
		t.Errorf("Name = %q, want Test Movie", resp.Name)
	}
}

func TestDebridCheckCache_Unauthorized(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(`{"error":"Invalid API key"}`))
	}))
	defer srv.Close()

	c := NewClient(WithBaseURL(srv.URL), WithRetry(0, 0, 0))
	_, err := c.DebridCheckCache(context.Background(), "real-debrid", "bad-key", []string{"hash1"})
	if err == nil {
		t.Fatal("expected error")
	}
	apiErr, ok := err.(*APIError)
	if !ok {
		t.Fatalf("expected *APIError, got %T", err)
	}
	if apiErr.StatusCode != 401 {
		t.Errorf("StatusCode = %d, want 401", apiErr.StatusCode)
	}
}

func TestDebridAddMagnet_ServerError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadGateway)
	}))
	defer srv.Close()

	c := NewClient(WithBaseURL(srv.URL), WithRetry(0, 0, 0))
	_, err := c.DebridAddMagnet(context.Background(), "real-debrid", "key", "hash")
	if err == nil {
		t.Fatal("expected error")
	}
}
