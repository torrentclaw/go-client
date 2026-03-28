package torrentclaw

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestStreamingTop(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/streaming-top" {
			t.Errorf("path = %q, want /api/v1/streaming-top", r.URL.Path)
		}
		q := r.URL.Query()
		if got := q.Get("service"); got != "netflix" {
			t.Errorf("service = %q, want netflix", got)
		}
		if got := q.Get("country"); got != "US" {
			t.Errorf("country = %q, want US", got)
		}
		if got := q.Get("show_type"); got != "movie" {
			t.Errorf("show_type = %q, want movie", got)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode([]StreamingTopItem{
			{Rank: 1, Title: "Top Movie", HasTorrents: true, MaxSeeders: 500},
			{Rank: 2, Title: "Second Movie", HasTorrents: false, MaxSeeders: 0},
		})
	}))
	defer srv.Close()

	c := NewClient(WithBaseURL(srv.URL), WithRetry(0, 0, 0))
	items, err := c.StreamingTop(context.Background(), StreamingTopParams{
		Service:  "netflix",
		Country:  "US",
		ShowType: "movie",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(items) != 2 {
		t.Fatalf("len(items) = %d, want 2", len(items))
	}
	if items[0].Rank != 1 {
		t.Errorf("Rank = %d, want 1", items[0].Rank)
	}
	if items[0].Title != "Top Movie" {
		t.Errorf("Title = %q, want Top Movie", items[0].Title)
	}
	if !items[0].HasTorrents {
		t.Error("HasTorrents should be true")
	}
}

func TestStreamingTop_DefaultParams(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		for _, key := range []string{"service", "country", "show_type"} {
			if q.Has(key) {
				t.Errorf("unexpected param %q", key)
			}
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode([]StreamingTopItem{})
	}))
	defer srv.Close()

	c := NewClient(WithBaseURL(srv.URL), WithRetry(0, 0, 0))
	items, err := c.StreamingTop(context.Background(), StreamingTopParams{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(items) != 0 {
		t.Errorf("len(items) = %d, want 0", len(items))
	}
}

func TestStreamingTop_ServerError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer srv.Close()

	c := NewClient(WithBaseURL(srv.URL), WithRetry(0, 0, 0))
	_, err := c.StreamingTop(context.Background(), StreamingTopParams{})
	if err == nil {
		t.Fatal("expected error")
	}
}
