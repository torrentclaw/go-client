package torrentclaw

import (
	"context"
	"net/url"
)

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
