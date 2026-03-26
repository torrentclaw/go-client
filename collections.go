package torrentclaw

import (
	"context"
	"fmt"
	"net/url"
)

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
