package torrentclaw

import (
	"context"
	"net/url"
)

// StreamingTopParams holds the parameters for a streaming top 10 request.
type StreamingTopParams struct {
	// Service is the streaming service: "netflix", "prime", "disney", "hbo", "apple".
	Service string

	// Country is an ISO 3166-1 country code (default "US").
	Country string

	// ShowType is the content type: "movie" or "series" (default "movie").
	ShowType string

	// Locale sets the language for localized titles.
	Locale string
}

// StreamingTopItem represents a ranked item in a streaming top list.
type StreamingTopItem struct {
	Rank          int      `json:"rank"`
	Title         string   `json:"title"`
	IMDbID        *string  `json:"imdbId,omitempty"`
	TMDbID        *int     `json:"tmdbId,omitempty"`
	ContentType   *string  `json:"contentType,omitempty"`
	Year          *int     `json:"year,omitempty"`
	Overview      *string  `json:"overview,omitempty"`
	Rating        *string  `json:"rating,omitempty"`
	PosterURL     *string  `json:"posterUrl,omitempty"`
	BackdropURL   *string  `json:"backdropUrl,omitempty"`
	Genres        []string `json:"genres,omitempty"`
	StreamingLink *string  `json:"streamingLink,omitempty"`
	ContentID     *int     `json:"contentId,omitempty"`
	HasTorrents   bool     `json:"hasTorrents"`
	MaxSeeders    int      `json:"maxSeeders"`
}

// StreamingTop returns the top-ranked content for a streaming service.
// The response is an array of ranked items (not a paginated wrapper).
func (c *Client) StreamingTop(ctx context.Context, params StreamingTopParams) ([]StreamingTopItem, error) {
	q := url.Values{}
	addStringParam(q, "service", params.Service)
	addStringParam(q, "country", params.Country)
	addStringParam(q, "show_type", params.ShowType)
	addStringParam(q, "locale", params.Locale)

	var resp []StreamingTopItem
	if err := c.doJSON(ctx, "/api/v1/streaming-top", q, &resp); err != nil {
		return nil, err
	}
	return resp, nil
}
