package torrentclaw

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSearch(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/search" {
			t.Errorf("path = %q, want /api/v1/search", r.URL.Path)
		}
		q := r.URL.Query()
		if got := q.Get("q"); got != "inception" {
			t.Errorf("q = %q, want inception", got)
		}
		if got := q.Get("type"); got != "movie" {
			t.Errorf("type = %q, want movie", got)
		}
		if got := q.Get("year_min"); got != "2010" {
			t.Errorf("year_min = %q, want 2010", got)
		}
		if got := q.Get("quality"); got != "1080p" {
			t.Errorf("quality = %q, want 1080p", got)
		}
		if got := q.Get("lang"); got != "en" {
			t.Errorf("lang = %q, want en", got)
		}
		if got := q.Get("sort"); got != "seeders" {
			t.Errorf("sort = %q, want seeders", got)
		}
		if got := q.Get("page"); got != "1" {
			t.Errorf("page = %q, want 1", got)
		}
		if got := q.Get("limit"); got != "10" {
			t.Errorf("limit = %q, want 10", got)
		}
		if got := q.Get("country"); got != "US" {
			t.Errorf("country = %q, want US", got)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(SearchResponse{
			Total:    1,
			Page:     1,
			PageSize: 10,
			Results: []SearchResult{
				{
					ID:          42,
					ContentType: "movie",
					Title:       "Inception",
					HasTorrents: true,
					Torrents: []TorrentInfo{
						{
							InfoHash:  "abc123",
							RawTitle:  "Inception.2010.1080p",
							Seeders:   100,
							Leechers:  10,
							Source:    "yts",
							Languages: []string{"en"},
						},
					},
				},
			},
		})
	}))
	defer srv.Close()

	c := NewClient(WithBaseURL(srv.URL), WithRetry(0, 0, 0))
	resp, err := c.Search(context.Background(), SearchParams{
		Query:    "inception",
		Type:     "movie",
		YearMin:  2010,
		Quality:  "1080p",
		Language: "en",
		Sort:     "seeders",
		Page:     1,
		Limit:    10,
		Country:  "US",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.Total != 1 {
		t.Errorf("Total = %d, want 1", resp.Total)
	}
	if len(resp.Results) != 1 {
		t.Fatalf("len(Results) = %d, want 1", len(resp.Results))
	}
	r := resp.Results[0]
	if r.Title != "Inception" {
		t.Errorf("Title = %q, want Inception", r.Title)
	}
	if len(r.Torrents) != 1 {
		t.Fatalf("len(Torrents) = %d, want 1", len(r.Torrents))
	}
	if r.Torrents[0].Seeders != 100 {
		t.Errorf("Seeders = %d, want 100", r.Torrents[0].Seeders)
	}
}

func TestSearch_AllParams(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		checks := map[string]string{
			"q":            "test",
			"type":         "show",
			"genre":        "Comedy",
			"year_min":     "2020",
			"year_max":     "2025",
			"min_rating":   "7",
			"quality":      "2160p",
			"lang":         "es",
			"subs":         "en",
			"audio":        "atmos",
			"hdr":          "dolby_vision",
			"sort":         "rating",
			"page":         "2",
			"limit":        "25",
			"country":      "ES",
			"locale":       "es",
			"availability": "available",
			"verified":     "true",
			"season":       "3",
			"episode":      "5",
		}
		for key, want := range checks {
			if got := q.Get(key); got != want {
				t.Errorf("%s = %q, want %q", key, got, want)
			}
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(SearchResponse{Total: 0, Page: 2, PageSize: 25})
	}))
	defer srv.Close()

	c := NewClient(WithBaseURL(srv.URL), WithRetry(0, 0, 0))
	_, err := c.Search(context.Background(), SearchParams{
		Query:        "test",
		Type:         "show",
		Genre:        "Comedy",
		YearMin:      2020,
		YearMax:      2025,
		MinRating:    7,
		Quality:      "2160p",
		Language:     "es",
		Subs:         "en",
		Audio:        "atmos",
		HDR:          "dolby_vision",
		Sort:         "rating",
		Page:         2,
		Limit:        25,
		Country:      "ES",
		Locale:       "es",
		Availability: "available",
		Verified:     true,
		Season:       3,
		Episode:      5,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestSearch_EmptyOptionalParams(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		if got := q.Get("q"); got != "matrix" {
			t.Errorf("q = %q, want matrix", got)
		}
		for _, key := range []string{"type", "genre", "year_min", "year_max", "quality", "lang", "subs", "sort", "page", "limit", "verified", "season", "episode"} {
			if q.Has(key) {
				t.Errorf("unexpected param %q = %q", key, q.Get(key))
			}
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(SearchResponse{Total: 0, Page: 1, PageSize: 20})
	}))
	defer srv.Close()

	c := NewClient(WithBaseURL(srv.URL), WithRetry(0, 0, 0))
	_, err := c.Search(context.Background(), SearchParams{Query: "matrix"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestSearch_ParsedSeasonEpisode(t *testing.T) {
	season := 2
	episode := 5
	fuzzy := true
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(SearchResponse{
			Total:         1,
			Page:          1,
			PageSize:      20,
			ParsedSeason:  &season,
			ParsedEpisode: &episode,
			FuzzyMatch:    &fuzzy,
			Results:       []SearchResult{},
		})
	}))
	defer srv.Close()

	c := NewClient(WithBaseURL(srv.URL), WithRetry(0, 0, 0))
	resp, err := c.Search(context.Background(), SearchParams{Query: "breaking bad s02e05"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.ParsedSeason == nil || *resp.ParsedSeason != 2 {
		t.Errorf("ParsedSeason = %v, want 2", resp.ParsedSeason)
	}
	if resp.ParsedEpisode == nil || *resp.ParsedEpisode != 5 {
		t.Errorf("ParsedEpisode = %v, want 5", resp.ParsedEpisode)
	}
	if resp.FuzzyMatch == nil || !*resp.FuzzyMatch {
		t.Errorf("FuzzyMatch = %v, want true", resp.FuzzyMatch)
	}
}

func TestAutocomplete(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/autocomplete" {
			t.Errorf("path = %q, want /api/v1/autocomplete", r.URL.Path)
		}
		if got := r.URL.Query().Get("q"); got != "incep" {
			t.Errorf("q = %q, want incep", got)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"suggestions": []AutocompleteSuggestion{
				{ID: 1, Title: "Inception", ContentType: "movie"},
			},
		})
	}))
	defer srv.Close()

	c := NewClient(WithBaseURL(srv.URL), WithRetry(0, 0, 0))
	results, err := c.Autocomplete(context.Background(), AutocompleteParams{Query: "incep"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 1 {
		t.Fatalf("len(results) = %d, want 1", len(results))
	}
	if results[0].Title != "Inception" {
		t.Errorf("Title = %q, want Inception", results[0].Title)
	}
}

func TestAutocomplete_WithLocale(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got := r.URL.Query().Get("locale"); got != "es" {
			t.Errorf("locale = %q, want es", got)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{"suggestions": []AutocompleteSuggestion{}})
	}))
	defer srv.Close()

	c := NewClient(WithBaseURL(srv.URL), WithRetry(0, 0, 0))
	_, err := c.Autocomplete(context.Background(), AutocompleteParams{Query: "test", Locale: "es"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestAutocomplete_WithMovieCount(t *testing.T) {
	movieCount := 5
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"suggestions": []AutocompleteSuggestion{
				{ID: 1, Title: "Star Wars", ContentType: "collection", MovieCount: &movieCount},
			},
		})
	}))
	defer srv.Close()

	c := NewClient(WithBaseURL(srv.URL), WithRetry(0, 0, 0))
	results, err := c.Autocomplete(context.Background(), AutocompleteParams{Query: "star"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if results[0].MovieCount == nil || *results[0].MovieCount != 5 {
		t.Errorf("MovieCount = %v, want 5", results[0].MovieCount)
	}
}

func TestSearch_ServerError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer srv.Close()

	c := NewClient(WithBaseURL(srv.URL), WithRetry(0, 0, 0))
	_, err := c.Search(context.Background(), SearchParams{Query: "test"})
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

func TestAutocomplete_ServerError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer srv.Close()

	c := NewClient(WithBaseURL(srv.URL), WithRetry(0, 0, 0))
	_, err := c.Autocomplete(context.Background(), AutocompleteParams{Query: "test"})
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestAutocomplete_EmptyResults(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{"suggestions": []AutocompleteSuggestion{}})
	}))
	defer srv.Close()

	c := NewClient(WithBaseURL(srv.URL), WithRetry(0, 0, 0))
	results, err := c.Autocomplete(context.Background(), AutocompleteParams{Query: "zzzzz"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 0 {
		t.Errorf("len(results) = %d, want 0", len(results))
	}
}
