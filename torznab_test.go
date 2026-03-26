package torrentclaw

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestTorznab_Search(t *testing.T) {
	xmlResp := `<?xml version="1.0" encoding="UTF-8"?><rss><channel><item><title>Test</title></item></channel></rss>`
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/torznab" {
			t.Errorf("path = %q, want /api/v1/torznab", r.URL.Path)
		}
		q := r.URL.Query()
		if got := q.Get("t"); got != "search" {
			t.Errorf("t = %q, want search", got)
		}
		if got := q.Get("q"); got != "inception" {
			t.Errorf("q = %q, want inception", got)
		}
		if got := q.Get("limit"); got != "50" {
			t.Errorf("limit = %q, want 50", got)
		}

		w.Header().Set("Content-Type", "application/xml")
		w.Write([]byte(xmlResp))
	}))
	defer srv.Close()

	c := NewClient(WithBaseURL(srv.URL), WithRetry(0, 0, 0))
	data, err := c.Torznab(context.Background(), TorznabParams{
		T:     "search",
		Q:     "inception",
		Limit: 50,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(string(data), "<title>Test</title>") {
		t.Errorf("response does not contain expected XML")
	}
}

func TestTorznab_TVSearch(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		if got := q.Get("t"); got != "tvsearch" {
			t.Errorf("t = %q, want tvsearch", got)
		}
		if got := q.Get("imdbid"); got != "tt1234567" {
			t.Errorf("imdbid = %q, want tt1234567", got)
		}
		if got := q.Get("season"); got != "3" {
			t.Errorf("season = %q, want 3", got)
		}
		if got := q.Get("ep"); got != "5" {
			t.Errorf("ep = %q, want 5", got)
		}
		w.Write([]byte("<rss/>"))
	}))
	defer srv.Close()

	c := NewClient(WithBaseURL(srv.URL), WithRetry(0, 0, 0))
	_, err := c.Torznab(context.Background(), TorznabParams{
		T:      "tvsearch",
		IMDbID: "tt1234567",
		Season: 3,
		Ep:     5,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestTorznabCaps(t *testing.T) {
	capsXML := `<?xml version="1.0"?><caps><server title="TorrentClaw"/></caps>`
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got := r.URL.Query().Get("t"); got != "caps" {
			t.Errorf("t = %q, want caps", got)
		}
		w.Write([]byte(capsXML))
	}))
	defer srv.Close()

	c := NewClient(WithBaseURL(srv.URL), WithRetry(0, 0, 0))
	data, err := c.TorznabCaps(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(string(data), "TorrentClaw") {
		t.Error("caps should contain TorrentClaw")
	}
}

func TestTorznab_ServerError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte("pro tier required"))
	}))
	defer srv.Close()

	c := NewClient(WithBaseURL(srv.URL), WithRetry(0, 0, 0))
	_, err := c.Torznab(context.Background(), TorznabParams{T: "search", Q: "test"})
	if err == nil {
		t.Fatal("expected error")
	}
	apiErr, ok := err.(*APIError)
	if !ok {
		t.Fatalf("expected *APIError, got %T", err)
	}
	if apiErr.StatusCode != 403 {
		t.Errorf("StatusCode = %d, want 403", apiErr.StatusCode)
	}
}
