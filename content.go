package torrentclaw

import (
	"context"
	"fmt"
	"net/url"
)

// Popular returns the most popular content ranked by community engagement.
func (c *Client) Popular(ctx context.Context, params PopularParams) (*PopularResponse, error) {
	q := url.Values{}
	addIntParam(q, "limit", params.Limit)
	addIntParam(q, "page", params.Page)
	addStringParam(q, "locale", params.Locale)

	var resp PopularResponse
	if err := c.doJSON(ctx, "/api/v1/popular", q, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// Recent returns the most recently added content.
func (c *Client) Recent(ctx context.Context, params RecentParams) (*RecentResponse, error) {
	q := url.Values{}
	addIntParam(q, "limit", params.Limit)
	addIntParam(q, "page", params.Page)
	addStringParam(q, "locale", params.Locale)

	var resp RecentResponse
	if err := c.doJSON(ctx, "/api/v1/recent", q, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// WatchProviders returns streaming/watch providers for a content item.
// The country parameter is an ISO 3166-1 code (e.g. "US", "ES"). Pass an
// empty string to use the server default ("US").
func (c *Client) WatchProviders(ctx context.Context, contentID int, country string) (*WatchProvidersResponse, error) {
	q := url.Values{}
	addStringParam(q, "country", country)

	path := fmt.Sprintf("/api/v1/content/%d/watch-providers", contentID)
	var resp WatchProvidersResponse
	if err := c.doJSON(ctx, path, q, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// Credits returns the director and top cast members for a content item.
func (c *Client) Credits(ctx context.Context, contentID int) (*CreditsResponse, error) {
	path := fmt.Sprintf("/api/v1/content/%d/credits", contentID)
	var resp CreditsResponse
	if err := c.doJSON(ctx, path, nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}
