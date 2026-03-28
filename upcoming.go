package torrentclaw

import (
	"context"
	"net/url"
)

// UpcomingParams holds the parameters for an upcoming releases request.
type UpcomingParams struct {
	// Limit sets the number of items (1-50, default 24).
	Limit int

	// Page is the page number (starts at 1).
	Page int

	// Type filters by content type: "all", "movie", "show" (default "all").
	Type string

	// Locale sets the language for localized titles.
	Locale string
}

// UpcomingItem represents an upcoming release.
type UpcomingItem struct {
	ID             int      `json:"id"`
	TMDbID         *int     `json:"tmdbId,omitempty"`
	Title          string   `json:"title"`
	Year           *int     `json:"year,omitempty"`
	ContentType    string   `json:"contentType"`
	PosterURL      *string  `json:"posterUrl,omitempty"`
	BackdropURL    *string  `json:"backdropUrl,omitempty"`
	Overview       *string  `json:"overview,omitempty"`
	Genres         []string `json:"genres,omitempty"`
	RatingIMDb     *string  `json:"ratingImdb,omitempty"`
	RatingTMDb     *string  `json:"ratingTmdb,omitempty"`
	PopularityTMDb *float64 `json:"popularityTmdb,omitempty"`
	ReleaseDate    string   `json:"releaseDate"`
	HasTorrents    bool     `json:"hasTorrents"`
}

// UpcomingResponse is the paginated response from the upcoming endpoint.
type UpcomingResponse struct {
	Items    []UpcomingItem `json:"items"`
	Total    int            `json:"total"`
	Page     int            `json:"page"`
	PageSize int            `json:"pageSize"`
}

// Upcoming returns upcoming releases, optionally filtered by content type.
func (c *Client) Upcoming(ctx context.Context, params UpcomingParams) (*UpcomingResponse, error) {
	q := url.Values{}
	addIntParam(q, "limit", params.Limit)
	addIntParam(q, "page", params.Page)
	addStringParam(q, "type", params.Type)
	addStringParam(q, "locale", params.Locale)

	var resp UpcomingResponse
	if err := c.doJSON(ctx, "/api/v1/upcoming", q, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}
