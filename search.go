package torrentclaw

import (
	"context"
	"net/url"
)

// SearchParams holds the parameters for a search request.
type SearchParams struct {
	// Query is the search query (required, max 200 characters).
	Query string

	// Type filters by content type: "movie" or "show".
	Type string

	// Genre filters by genre (exact match), e.g. "Action", "Comedy".
	Genre string

	// YearMin sets the minimum release year.
	YearMin int

	// YearMax sets the maximum release year.
	YearMax int

	// MinRating sets the minimum IMDb/TMDB rating threshold (0-10).
	MinRating float64

	// Quality filters by video resolution: "480p", "720p", "1080p", "2160p".
	Quality string

	// Language filters by torrent audio language (ISO 639 code, e.g. "en", "es").
	Language string

	// Subs filters by subtitle language (ISO 639 code, e.g. "en", "es").
	Subs string

	// Audio filters by audio codec (substring match, e.g. "aac", "atmos").
	Audio string

	// HDR filters by HDR format (e.g. "hdr10", "dolby_vision").
	HDR string

	// Sort sets the sort order: "relevance", "seeders", "year", "rating", "added".
	Sort string

	// Page is the page number (starts at 1).
	Page int

	// Limit sets the number of results per page (1-50).
	Limit int

	// Country is an ISO 3166-1 country code to include streaming availability.
	Country string

	// Locale sets the language for title/overview translations (e.g. "es", "fr").
	Locale string

	// Availability filters by torrent availability: "all", "available", "unavailable".
	Availability string

	// Verified filters to only TrueSpec verified releases.
	Verified bool

	// Season filters by TV show season number.
	Season int

	// Episode filters by TV show episode number.
	Episode int
}

// SearchResult represents a single content result from a search.
type SearchResult struct {
	ID            int            `json:"id"`
	IMDbID        *string        `json:"imdbId,omitempty"`
	TMDbID        *string        `json:"tmdbId,omitempty"`
	ContentType   string         `json:"contentType"`
	Title         string         `json:"title"`
	TitleOriginal *string        `json:"titleOriginal,omitempty"`
	Year          *int           `json:"year,omitempty"`
	Overview      *string        `json:"overview,omitempty"`
	PosterURL     *string        `json:"posterUrl,omitempty"`
	BackdropURL   *string        `json:"backdropUrl,omitempty"`
	Genres        []string       `json:"genres,omitempty"`
	RatingIMDb    *string        `json:"ratingImdb,omitempty"`
	RatingTMDb    *string        `json:"ratingTmdb,omitempty"`
	ContentURL    string         `json:"contentUrl"`
	HasTorrents   bool           `json:"hasTorrents"`
	Torrents      []TorrentInfo  `json:"torrents"`
	Streaming     *StreamingInfo `json:"streaming,omitempty"`
}

// SearchResponse is the paginated response from the search endpoint.
type SearchResponse struct {
	Total         int            `json:"total"`
	Page          int            `json:"page"`
	PageSize      int            `json:"pageSize"`
	ParsedSeason  *int           `json:"parsedSeason,omitempty"`
	ParsedEpisode *int           `json:"parsedEpisode,omitempty"`
	FuzzyMatch    *bool          `json:"fuzzyMatch,omitempty"`
	Results       []SearchResult `json:"results"`
}

// AutocompleteParams holds the parameters for an autocomplete request.
type AutocompleteParams struct {
	// Query is the search prefix (required, 2-200 characters).
	Query string

	// Locale sets the language for localized suggestions.
	Locale string
}

// AutocompleteSuggestion represents a single autocomplete result.
type AutocompleteSuggestion struct {
	ID          int     `json:"id"`
	Title       string  `json:"title"`
	Year        *int    `json:"year,omitempty"`
	ContentType string  `json:"contentType"`
	PosterURL   *string `json:"posterUrl,omitempty"`
	MovieCount  *int    `json:"movieCount,omitempty"`
}

// Search queries the TorrentClaw search endpoint with the given parameters.
// The Query field in params is required.
func (c *Client) Search(ctx context.Context, params SearchParams) (*SearchResponse, error) {
	q := url.Values{}
	q.Set("q", params.Query)
	addStringParam(q, "type", params.Type)
	addStringParam(q, "genre", params.Genre)
	addIntParam(q, "year_min", params.YearMin)
	addIntParam(q, "year_max", params.YearMax)
	addFloatParam(q, "min_rating", params.MinRating)
	addStringParam(q, "quality", params.Quality)
	addStringParam(q, "lang", params.Language)
	addStringParam(q, "subs", params.Subs)
	addStringParam(q, "audio", params.Audio)
	addStringParam(q, "hdr", params.HDR)
	addStringParam(q, "sort", params.Sort)
	addIntParam(q, "page", params.Page)
	addIntParam(q, "limit", params.Limit)
	addStringParam(q, "country", params.Country)
	addStringParam(q, "locale", params.Locale)
	addStringParam(q, "availability", params.Availability)
	addBoolParam(q, "verified", params.Verified)
	addIntParam(q, "season", params.Season)
	addIntParam(q, "episode", params.Episode)

	var resp SearchResponse
	if err := c.doJSON(ctx, "/api/v1/search", q, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// Autocomplete returns title suggestions for the given query prefix.
func (c *Client) Autocomplete(ctx context.Context, params AutocompleteParams) ([]AutocompleteSuggestion, error) {
	q := url.Values{}
	q.Set("q", params.Query)
	addStringParam(q, "locale", params.Locale)

	var resp struct {
		Suggestions []AutocompleteSuggestion `json:"suggestions"`
	}
	if err := c.doJSON(ctx, "/api/v1/autocomplete", q, &resp); err != nil {
		return nil, err
	}
	return resp.Suggestions, nil
}
