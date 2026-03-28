package torrentclaw

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestStats_ServerError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer srv.Close()

	c := NewClient(WithBaseURL(srv.URL), WithRetry(0, 0, 0))
	_, err := c.Stats(context.Background())
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

func TestStats_EmptyIngestions(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(StatsResponse{
			Content:          ContentStats{Movies: 100, Shows: 50, TMDbEnriched: 80},
			Torrents:         TorrentStats{Total: 500, WithSeeders: 300, BySource: map[string]int{}},
			RecentIngestions: []IngestionRecord{},
		})
	}))
	defer srv.Close()

	c := NewClient(WithBaseURL(srv.URL), WithRetry(0, 0, 0))
	resp, err := c.Stats(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.Content.TMDbEnriched != 80 {
		t.Errorf("TMDbEnriched = %d, want 80", resp.Content.TMDbEnriched)
	}
	if resp.Torrents.WithSeeders != 300 {
		t.Errorf("WithSeeders = %d, want 300", resp.Torrents.WithSeeders)
	}
	if len(resp.RecentIngestions) != 0 {
		t.Errorf("len(RecentIngestions) = %d, want 0", len(resp.RecentIngestions))
	}
}

func TestStats_WithCompletedAt(t *testing.T) {
	completedAt := "2025-01-15T10:05:00Z"
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(StatsResponse{
			Content:  ContentStats{},
			Torrents: TorrentStats{BySource: map[string]int{}},
			RecentIngestions: []IngestionRecord{
				{
					Source:      "yts",
					Status:      "completed",
					StartedAt:   "2025-01-15T10:00:00Z",
					CompletedAt: &completedAt,
					Fetched:     100,
					New:         10,
					Updated:     5,
				},
			},
		})
	}))
	defer srv.Close()

	c := NewClient(WithBaseURL(srv.URL), WithRetry(0, 0, 0))
	resp, err := c.Stats(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(resp.RecentIngestions) != 1 {
		t.Fatalf("len(RecentIngestions) = %d, want 1", len(resp.RecentIngestions))
	}
	rec := resp.RecentIngestions[0]
	if rec.CompletedAt == nil || *rec.CompletedAt != completedAt {
		t.Errorf("CompletedAt = %v, want %q", rec.CompletedAt, completedAt)
	}
}
