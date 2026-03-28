package torrentclaw

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestPopular(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/popular" {
			t.Errorf("path = %q, want /api/v1/popular", r.URL.Path)
		}
		if got := r.URL.Query().Get("limit"); got != "5" {
			t.Errorf("limit = %q, want 5", got)
		}
		if got := r.URL.Query().Get("page"); got != "2" {
			t.Errorf("page = %q, want 2", got)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(PopularResponse{
			Items: []PopularItem{
				{ID: 1, Title: "The Matrix", ContentType: "movie", MaxSeeders: 500},
			},
			Total:    100,
			Page:     2,
			PageSize: 5,
		})
	}))
	defer srv.Close()

	c := NewClient(WithBaseURL(srv.URL), WithRetry(0, 0, 0))
	resp, err := c.Popular(context.Background(), PopularParams{Limit: 5, Page: 2})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.Total != 100 {
		t.Errorf("Total = %d, want 100", resp.Total)
	}
	if len(resp.Items) != 1 {
		t.Fatalf("len(Items) = %d, want 1", len(resp.Items))
	}
	if resp.Items[0].MaxSeeders != 500 {
		t.Errorf("MaxSeeders = %d, want 500", resp.Items[0].MaxSeeders)
	}
}

func TestPopular_DefaultParams(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		if q.Has("limit") {
			t.Errorf("unexpected limit param: %q", q.Get("limit"))
		}
		if q.Has("page") {
			t.Errorf("unexpected page param: %q", q.Get("page"))
		}
		if q.Has("locale") {
			t.Errorf("unexpected locale param: %q", q.Get("locale"))
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(PopularResponse{Total: 0, Page: 1, PageSize: 12})
	}))
	defer srv.Close()

	c := NewClient(WithBaseURL(srv.URL), WithRetry(0, 0, 0))
	_, err := c.Popular(context.Background(), PopularParams{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestPopular_WithLocale(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got := r.URL.Query().Get("locale"); got != "es" {
			t.Errorf("locale = %q, want es", got)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(PopularResponse{Total: 0, Page: 1, PageSize: 12})
	}))
	defer srv.Close()

	c := NewClient(WithBaseURL(srv.URL), WithRetry(0, 0, 0))
	_, err := c.Popular(context.Background(), PopularParams{Locale: "es"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRecent(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/recent" {
			t.Errorf("path = %q, want /api/v1/recent", r.URL.Path)
		}

		overview := "A great movie"
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(RecentResponse{
			Items: []RecentItem{
				{ID: 10, Title: "New Movie", ContentType: "movie", Overview: &overview, CreatedAt: "2025-01-15T10:00:00Z"},
			},
			Total:    50,
			Page:     1,
			PageSize: 12,
		})
	}))
	defer srv.Close()

	c := NewClient(WithBaseURL(srv.URL), WithRetry(0, 0, 0))
	resp, err := c.Recent(context.Background(), RecentParams{Limit: 12, Page: 1})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.Total != 50 {
		t.Errorf("Total = %d, want 50", resp.Total)
	}
	if len(resp.Items) != 1 {
		t.Fatalf("len(Items) = %d, want 1", len(resp.Items))
	}
	if resp.Items[0].CreatedAt != "2025-01-15T10:00:00Z" {
		t.Errorf("CreatedAt = %q", resp.Items[0].CreatedAt)
	}
	if resp.Items[0].Overview == nil || *resp.Items[0].Overview != "A great movie" {
		t.Errorf("Overview = %v, want 'A great movie'", resp.Items[0].Overview)
	}
}

func TestRecent_WithLocale(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got := r.URL.Query().Get("locale"); got != "fr" {
			t.Errorf("locale = %q, want fr", got)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(RecentResponse{Total: 0, Page: 1, PageSize: 12})
	}))
	defer srv.Close()

	c := NewClient(WithBaseURL(srv.URL), WithRetry(0, 0, 0))
	_, err := c.Recent(context.Background(), RecentParams{Locale: "fr"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestWatchProviders(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/content/42/watch-providers" {
			t.Errorf("path = %q", r.URL.Path)
		}
		if got := r.URL.Query().Get("country"); got != "US" {
			t.Errorf("country = %q, want US", got)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(WatchProvidersResponse{
			ContentID: 42,
			Country:   "US",
			Providers: WatchProviders{
				Flatrate: []WatchProviderItem{
					{ProviderID: 8, Name: "Netflix"},
				},
			},
			Attribution: "Powered by JustWatch",
		})
	}))
	defer srv.Close()

	c := NewClient(WithBaseURL(srv.URL), WithRetry(0, 0, 0))
	resp, err := c.WatchProviders(context.Background(), 42, "US")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.ContentID != 42 {
		t.Errorf("ContentID = %d, want 42", resp.ContentID)
	}
	if resp.Country != "US" {
		t.Errorf("Country = %q, want US", resp.Country)
	}
	if len(resp.Providers.Flatrate) != 1 {
		t.Fatalf("len(Flatrate) = %d, want 1", len(resp.Providers.Flatrate))
	}
	if resp.Providers.Flatrate[0].Name != "Netflix" {
		t.Errorf("Name = %q, want Netflix", resp.Providers.Flatrate[0].Name)
	}
}

func TestWatchProviders_WithVPN(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(WatchProvidersResponse{
			ContentID: 42,
			Country:   "AR",
			Providers: WatchProviders{},
			VPNSuggestion: &VPNSuggestion{
				AvailableIn:  []string{"US", "GB"},
				AffiliateURL: "https://example.com/vpn",
			},
			Attribution: "Powered by JustWatch",
		})
	}))
	defer srv.Close()

	c := NewClient(WithBaseURL(srv.URL), WithRetry(0, 0, 0))
	resp, err := c.WatchProviders(context.Background(), 42, "AR")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.VPNSuggestion == nil {
		t.Fatal("expected VPNSuggestion")
	}
	if len(resp.VPNSuggestion.AvailableIn) != 2 {
		t.Errorf("len(AvailableIn) = %d, want 2", len(resp.VPNSuggestion.AvailableIn))
	}
}

func TestCredits(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/content/42/credits" {
			t.Errorf("path = %q", r.URL.Path)
		}

		director := "Christopher Nolan"
		directorID := 525
		tmdbID := 6193
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(CreditsResponse{
			ContentID:      42,
			Director:       &director,
			DirectorTmdbID: &directorID,
			Cast: []CastMember{
				{TmdbID: &tmdbID, Name: "Leonardo DiCaprio", Character: "Cobb"},
			},
		})
	}))
	defer srv.Close()

	c := NewClient(WithBaseURL(srv.URL), WithRetry(0, 0, 0))
	resp, err := c.Credits(context.Background(), 42)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.ContentID != 42 {
		t.Errorf("ContentID = %d, want 42", resp.ContentID)
	}
	if resp.Director == nil || *resp.Director != "Christopher Nolan" {
		t.Errorf("Director = %v", resp.Director)
	}
	if resp.DirectorTmdbID == nil || *resp.DirectorTmdbID != 525 {
		t.Errorf("DirectorTmdbID = %v, want 525", resp.DirectorTmdbID)
	}
	if len(resp.Cast) != 1 {
		t.Fatalf("len(Cast) = %d, want 1", len(resp.Cast))
	}
	if resp.Cast[0].Name != "Leonardo DiCaprio" {
		t.Errorf("Name = %q", resp.Cast[0].Name)
	}
	if resp.Cast[0].TmdbID == nil || *resp.Cast[0].TmdbID != 6193 {
		t.Errorf("TmdbID = %v, want 6193", resp.Cast[0].TmdbID)
	}
}

func TestStats(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/stats" {
			t.Errorf("path = %q, want /api/v1/stats", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(StatsResponse{
			Content: ContentStats{
				Movies:       50000,
				Shows:        10000,
				TMDbEnriched: 45000,
			},
			Torrents: TorrentStats{
				Total:        200000,
				WithSeeders:  150000,
				Orphans:      5000,
				DailyAverage: 1200,
				BySource: map[string]int{
					"yts":  50000,
					"eztv": 30000,
				},
			},
			RecentIngestions: []IngestionRecord{
				{
					Source:    "yts",
					Status:    "completed",
					StartedAt: "2025-01-15T10:00:00Z",
					Fetched:   100,
					New:       10,
					Updated:   5,
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
	if resp.Content.Movies != 50000 {
		t.Errorf("Movies = %d, want 50000", resp.Content.Movies)
	}
	if resp.Torrents.Total != 200000 {
		t.Errorf("Total = %d, want 200000", resp.Torrents.Total)
	}
	if resp.Torrents.Orphans != 5000 {
		t.Errorf("Orphans = %d, want 5000", resp.Torrents.Orphans)
	}
	if resp.Torrents.DailyAverage != 1200 {
		t.Errorf("DailyAverage = %d, want 1200", resp.Torrents.DailyAverage)
	}
	if resp.Torrents.BySource["yts"] != 50000 {
		t.Errorf("BySource[yts] = %d, want 50000", resp.Torrents.BySource["yts"])
	}
	if len(resp.RecentIngestions) != 1 {
		t.Fatalf("len(RecentIngestions) = %d, want 1", len(resp.RecentIngestions))
	}
}

func TestTorrentDownloadURL(t *testing.T) {
	c := NewClient()
	got := c.TorrentDownloadURL("abc123def456")
	want := "https://torrentclaw.com/api/v1/torrent/abc123def456"
	if got != want {
		t.Errorf("TorrentDownloadURL = %q, want %q", got, want)
	}
}

func TestGetTorrentFile(t *testing.T) {
	torrentData := []byte{0x64, 0x38, 0x3A, 0x61} // fake torrent bytes
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/torrent/abc123" {
			t.Errorf("path = %q", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/x-bittorrent")
		w.Write(torrentData)
	}))
	defer srv.Close()

	c := NewClient(WithBaseURL(srv.URL), WithRetry(0, 0, 0))
	data, err := c.GetTorrentFile(context.Background(), "abc123")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(data) != len(torrentData) {
		t.Errorf("len(data) = %d, want %d", len(data), len(torrentData))
	}
}

func TestGetTorrentFile_NotFound(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("not found"))
	}))
	defer srv.Close()

	c := NewClient(WithBaseURL(srv.URL), WithRetry(0, 0, 0))
	_, err := c.GetTorrentFile(context.Background(), "nonexistent")
	if err == nil {
		t.Fatal("expected error")
	}
	apiErr, ok := err.(*APIError)
	if !ok {
		t.Fatalf("expected *APIError, got %T", err)
	}
	if !apiErr.IsNotFound() {
		t.Errorf("expected IsNotFound, got StatusCode=%d", apiErr.StatusCode)
	}
}

func TestPopular_ServerError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer srv.Close()

	c := NewClient(WithBaseURL(srv.URL), WithRetry(0, 0, 0))
	_, err := c.Popular(context.Background(), PopularParams{Limit: 10, Page: 1})
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestRecent_ServerError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer srv.Close()

	c := NewClient(WithBaseURL(srv.URL), WithRetry(0, 0, 0))
	_, err := c.Recent(context.Background(), RecentParams{Limit: 10, Page: 1})
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestWatchProviders_ServerError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer srv.Close()

	c := NewClient(WithBaseURL(srv.URL), WithRetry(0, 0, 0))
	_, err := c.WatchProviders(context.Background(), 42, "US")
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestWatchProviders_NoCountry(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Has("country") {
			t.Error("country param should not be set when empty")
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(WatchProvidersResponse{
			ContentID:   42,
			Country:     "US",
			Providers:   WatchProviders{},
			Attribution: "Powered by JustWatch",
		})
	}))
	defer srv.Close()

	c := NewClient(WithBaseURL(srv.URL), WithRetry(0, 0, 0))
	_, err := c.WatchProviders(context.Background(), 42, "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestCredits_ServerError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("not found"))
	}))
	defer srv.Close()

	c := NewClient(WithBaseURL(srv.URL), WithRetry(0, 0, 0))
	_, err := c.Credits(context.Background(), 999)
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestCredits_NilDirector(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(CreditsResponse{
			ContentID: 42,
			Director:  nil,
			Cast:      []CastMember{},
		})
	}))
	defer srv.Close()

	c := NewClient(WithBaseURL(srv.URL), WithRetry(0, 0, 0))
	resp, err := c.Credits(context.Background(), 42)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.Director != nil {
		t.Errorf("Director = %v, want nil", resp.Director)
	}
}
