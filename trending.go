package torrentclaw

import (
	"context"
	"net/url"
)

// TrendingParams holds the parameters for a trending content request.
type TrendingParams struct {
	// Period sets the time window: "daily", "weekly", "monthly" (default "daily").
	Period string

	// Limit sets the number of items (1-50, default 20).
	Limit int

	// Page is the page number (starts at 1).
	Page int

	// Locale sets the language for localized titles.
	Locale string
}

// TrendingItem represents a trending content item.
type TrendingItem struct {
	ID          int     `json:"id"`
	Title       string  `json:"title"`
	Year        *int    `json:"year,omitempty"`
	ContentType string  `json:"contentType"`
	PosterURL   *string `json:"posterUrl,omitempty"`
	RatingIMDb  *string `json:"ratingImdb,omitempty"`
	RatingTMDb  *string `json:"ratingTmdb,omitempty"`
	Overview    *string `json:"overview,omitempty"`
	MaxSeeders  int     `json:"maxSeeders"`
	ClickCount  int     `json:"clickCount"`
	TrendScore  int     `json:"trendScore"`
}

// TrendingResponse is the paginated response from the trending endpoint.
type TrendingResponse struct {
	Period   string         `json:"period"`
	Items    []TrendingItem `json:"items"`
	Total    int            `json:"total"`
	Page     int            `json:"page"`
	PageSize int            `json:"pageSize"`
}

// Trending returns trending content for the given time period.
func (c *Client) Trending(ctx context.Context, params TrendingParams) (*TrendingResponse, error) {
	q := url.Values{}
	addStringParam(q, "period", params.Period)
	addIntParam(q, "limit", params.Limit)
	addIntParam(q, "page", params.Page)
	addStringParam(q, "locale", params.Locale)

	var resp TrendingResponse
	if err := c.doJSON(ctx, "/api/v1/trending", q, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}
