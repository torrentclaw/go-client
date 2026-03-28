package torrentclaw

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCollections(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/collections" {
			t.Errorf("path = %q, want /api/v1/collections", r.URL.Path)
		}
		if got := r.URL.Query().Get("limit"); got != "12" {
			t.Errorf("limit = %q, want 12", got)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(CollectionListResponse{
			Items: []CollectionListItem{
				{ID: 1, TMDbID: 10, Name: "Star Wars", MovieCount: 9, TotalSeeders: 5000, PartCount: 9},
			},
			Total:    50,
			Page:     1,
			PageSize: 12,
		})
	}))
	defer srv.Close()

	c := NewClient(WithBaseURL(srv.URL), WithRetry(0, 0, 0))
	resp, err := c.Collections(context.Background(), CollectionListParams{Limit: 12})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.Total != 50 {
		t.Errorf("Total = %d, want 50", resp.Total)
	}
	if len(resp.Items) != 1 {
		t.Fatalf("len(Items) = %d, want 1", len(resp.Items))
	}
	if resp.Items[0].Name != "Star Wars" {
		t.Errorf("Name = %q, want Star Wars", resp.Items[0].Name)
	}
	if resp.Items[0].MovieCount != 9 {
		t.Errorf("MovieCount = %d, want 9", resp.Items[0].MovieCount)
	}
}

func TestCollectionByID(t *testing.T) {
	overview := "The epic saga"
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/collections/10" {
			t.Errorf("path = %q, want /api/v1/collections/10", r.URL.Path)
		}
		if got := r.URL.Query().Get("locale"); got != "es" {
			t.Errorf("locale = %q, want es", got)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(CollectionDetail{
			ID:           1,
			TMDbID:       10,
			Name:         "Star Wars",
			MovieCount:   9,
			TotalSeeders: 5000,
			PartCount:    9,
			Overview:     &overview,
			Movies: []PopularItem{
				{ID: 100, Title: "A New Hope", ContentType: "movie", MaxSeeders: 800},
			},
		})
	}))
	defer srv.Close()

	c := NewClient(WithBaseURL(srv.URL), WithRetry(0, 0, 0))
	resp, err := c.CollectionByID(context.Background(), 10, "es")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.Name != "Star Wars" {
		t.Errorf("Name = %q, want Star Wars", resp.Name)
	}
	if resp.Overview == nil || *resp.Overview != "The epic saga" {
		t.Errorf("Overview = %v", resp.Overview)
	}
	if len(resp.Movies) != 1 {
		t.Fatalf("len(Movies) = %d, want 1", len(resp.Movies))
	}
	if resp.Movies[0].Title != "A New Hope" {
		t.Errorf("Title = %q, want A New Hope", resp.Movies[0].Title)
	}
}

func TestCollectionByID_NotFound(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"error":"Collection not found"}`))
	}))
	defer srv.Close()

	c := NewClient(WithBaseURL(srv.URL), WithRetry(0, 0, 0))
	_, err := c.CollectionByID(context.Background(), 999, "")
	if err == nil {
		t.Fatal("expected error")
	}
	apiErr, ok := err.(*APIError)
	if !ok {
		t.Fatalf("expected *APIError, got %T", err)
	}
	if !apiErr.IsNotFound() {
		t.Errorf("expected 404, got %d", apiErr.StatusCode)
	}
}

func TestCollections_ServerError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer srv.Close()

	c := NewClient(WithBaseURL(srv.URL), WithRetry(0, 0, 0))
	_, err := c.Collections(context.Background(), CollectionListParams{})
	if err == nil {
		t.Fatal("expected error")
	}
}
