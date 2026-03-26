package torrentclaw

import (
	"context"
	"net/url"
)

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
