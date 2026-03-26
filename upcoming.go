package torrentclaw

import (
	"context"
	"net/url"
)

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
