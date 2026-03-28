package torrentclaw

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestTrending(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/trending" {
			t.Errorf("path = %q, want /api/v1/trending", r.URL.Path)
		}
		q := r.URL.Query()
		if got := q.Get("period"); got != "weekly" {
			t.Errorf("period = %q, want weekly", got)
		}
		if got := q.Get("limit"); got != "10" {
			t.Errorf("limit = %q, want 10", got)
		}
		if got := q.Get("locale"); got != "es" {
			t.Errorf("locale = %q, want es", got)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(TrendingResponse{
			Period: "weekly",
			Items: []TrendingItem{
				{ID: 1, Title: "Trending Movie", ContentType: "movie", MaxSeeders: 1000, TrendScore: 95},
			},
			Total:    50,
			Page:     1,
			PageSize: 10,
		})
	}))
	defer srv.Close()

	c := NewClient(WithBaseURL(srv.URL), WithRetry(0, 0, 0))
	resp, err := c.Trending(context.Background(), TrendingParams{Period: "weekly", Limit: 10, Locale: "es"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.Period != "weekly" {
		t.Errorf("Period = %q, want weekly", resp.Period)
	}
	if len(resp.Items) != 1 {
		t.Fatalf("len(Items) = %d, want 1", len(resp.Items))
	}
	if resp.Items[0].TrendScore != 95 {
		t.Errorf("TrendScore = %d, want 95", resp.Items[0].TrendScore)
	}
}

func TestTrending_DefaultParams(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		for _, key := range []string{"period", "limit", "page", "locale"} {
			if q.Has(key) {
				t.Errorf("unexpected param %q", key)
			}
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(TrendingResponse{Period: "daily", Total: 0, Page: 1, PageSize: 20})
	}))
	defer srv.Close()

	c := NewClient(WithBaseURL(srv.URL), WithRetry(0, 0, 0))
	_, err := c.Trending(context.Background(), TrendingParams{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestTrending_ServerError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer srv.Close()

	c := NewClient(WithBaseURL(srv.URL), WithRetry(0, 0, 0))
	_, err := c.Trending(context.Background(), TrendingParams{})
	if err == nil {
		t.Fatal("expected error")
	}
}
