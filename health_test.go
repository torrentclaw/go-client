package torrentclaw

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHealth(t *testing.T) {
	dbStatus := "connected"
	redisStatus := "connected"
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/health" {
			t.Errorf("path = %q, want /api/health", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(HealthResponse{
			Status:    "ok",
			Timestamp: "2026-03-26T10:00:00Z",
			Uptime:    86400,
			Database:  &dbStatus,
			Redis:     &redisStatus,
		})
	}))
	defer srv.Close()

	c := NewClient(WithBaseURL(srv.URL), WithRetry(0, 0, 0))
	resp, err := c.Health(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.Status != "ok" {
		t.Errorf("Status = %q, want ok", resp.Status)
	}
	if resp.Uptime != 86400 {
		t.Errorf("Uptime = %d, want 86400", resp.Uptime)
	}
	if resp.Database == nil || *resp.Database != "connected" {
		t.Errorf("Database = %v, want connected", resp.Database)
	}
}

func TestHealth_ServerError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer srv.Close()

	c := NewClient(WithBaseURL(srv.URL), WithRetry(0, 0, 0))
	_, err := c.Health(context.Background())
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestMirrors(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/mirrors" {
			t.Errorf("path = %q, want /api/mirrors", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(MirrorsResponse{
			Mirrors: []MirrorInfo{
				{URL: "https://torrentclaw.com", Label: "Primary", Primary: true},
				{URL: "https://tc2.example.com", Label: "Mirror 1", Primary: false},
			},
			Tor: &TorInfo{URL: "http://example.onion", Label: ".onion"},
			Channels: []StatusChannel{
				{Label: "Telegram", URL: "https://t.me/torrentclaw"},
			},
		})
	}))
	defer srv.Close()

	c := NewClient(WithBaseURL(srv.URL), WithRetry(0, 0, 0))
	resp, err := c.Mirrors(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(resp.Mirrors) != 2 {
		t.Fatalf("len(Mirrors) = %d, want 2", len(resp.Mirrors))
	}
	if !resp.Mirrors[0].Primary {
		t.Error("first mirror should be primary")
	}
	if resp.Tor == nil {
		t.Fatal("expected Tor")
	}
	if resp.Tor.URL != "http://example.onion" {
		t.Errorf("Tor.URL = %q", resp.Tor.URL)
	}
	if len(resp.Channels) != 1 {
		t.Fatalf("len(Channels) = %d, want 1", len(resp.Channels))
	}
}

func TestMirrors_ServerError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer srv.Close()

	c := NewClient(WithBaseURL(srv.URL), WithRetry(0, 0, 0))
	_, err := c.Mirrors(context.Background())
	if err == nil {
		t.Fatal("expected error")
	}
}
