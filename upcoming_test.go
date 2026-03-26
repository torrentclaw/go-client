package torrentclaw

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestUpcoming(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/upcoming" {
			t.Errorf("path = %q, want /api/v1/upcoming", r.URL.Path)
		}
		q := r.URL.Query()
		if got := q.Get("type"); got != "movie" {
			t.Errorf("type = %q, want movie", got)
		}
		if got := q.Get("limit"); got != "10" {
			t.Errorf("limit = %q, want 10", got)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(UpcomingResponse{
			Items: []UpcomingItem{
				{ID: 1, Title: "Future Movie", ContentType: "movie", ReleaseDate: "2026-05-01", HasTorrents: false},
			},
			Total:    25,
			Page:     1,
			PageSize: 10,
		})
	}))
	defer srv.Close()

	c := NewClient(WithBaseURL(srv.URL), WithRetry(0, 0, 0))
	resp, err := c.Upcoming(context.Background(), UpcomingParams{Type: "movie", Limit: 10})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.Total != 25 {
		t.Errorf("Total = %d, want 25", resp.Total)
	}
	if len(resp.Items) != 1 {
		t.Fatalf("len(Items) = %d, want 1", len(resp.Items))
	}
	if resp.Items[0].ReleaseDate != "2026-05-01" {
		t.Errorf("ReleaseDate = %q", resp.Items[0].ReleaseDate)
	}
	if resp.Items[0].HasTorrents {
		t.Error("HasTorrents should be false for upcoming")
	}
}

func TestUpcoming_ServerError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer srv.Close()

	c := NewClient(WithBaseURL(srv.URL), WithRetry(0, 0, 0))
	_, err := c.Upcoming(context.Background(), UpcomingParams{})
	if err == nil {
		t.Fatal("expected error")
	}
}
