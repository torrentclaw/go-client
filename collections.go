package torrentclaw

import (
	"context"
	"fmt"
	"net/url"
)

// CollectionListParams holds the parameters for listing collections.
type CollectionListParams struct {
	// Limit sets the number of items (1-48, default 24).
	Limit int

	// Page is the page number (starts at 1).
	Page int

	// Locale sets the language for localized names.
	Locale string
}

// CollectionListItem represents a movie collection in a list.
type CollectionListItem struct {
	ID           int     `json:"id"`
	TMDbID       int     `json:"tmdbId"`
	Name         string  `json:"name"`
	PosterURL    *string `json:"posterUrl,omitempty"`
	BackdropURL  *string `json:"backdropUrl,omitempty"`
	MovieCount   int     `json:"movieCount"`
	TotalSeeders int     `json:"totalSeeders"`
	PartCount    int     `json:"partCount"`
}

// CollectionListResponse is the paginated response from the collections endpoint.
type CollectionListResponse struct {
	Items    []CollectionListItem `json:"items"`
	Total    int                  `json:"total"`
	Page     int                  `json:"page"`
	PageSize int                  `json:"pageSize"`
}

// CollectionDetail contains full details about a movie collection.
type CollectionDetail struct {
	ID           int           `json:"id"`
	TMDbID       int           `json:"tmdbId"`
	Name         string        `json:"name"`
	PosterURL    *string       `json:"posterUrl,omitempty"`
	BackdropURL  *string       `json:"backdropUrl,omitempty"`
	MovieCount   int           `json:"movieCount"`
	TotalSeeders int           `json:"totalSeeders"`
	PartCount    int           `json:"partCount"`
	Overview     *string       `json:"overview,omitempty"`
	Movies       []PopularItem `json:"movies"`
}

// Collections returns a paginated list of movie collections (sagas).
func (c *Client) Collections(ctx context.Context, params CollectionListParams) (*CollectionListResponse, error) {
	q := url.Values{}
	addIntParam(q, "limit", params.Limit)
	addIntParam(q, "page", params.Page)
	addStringParam(q, "locale", params.Locale)

	var resp CollectionListResponse
	if err := c.doJSON(ctx, "/api/v1/collections", q, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// CollectionByID returns full details for a single movie collection.
func (c *Client) CollectionByID(ctx context.Context, id int, locale string) (*CollectionDetail, error) {
	q := url.Values{}
	addStringParam(q, "locale", locale)

	path := fmt.Sprintf("/api/v1/collections/%d", id)
	var resp CollectionDetail
	if err := c.doJSON(ctx, path, q, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}
