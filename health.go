package torrentclaw

import "context"

// Health returns the API health status. This endpoint does not require authentication.
func (c *Client) Health(ctx context.Context) (*HealthResponse, error) {
	var resp HealthResponse
	if err := c.doJSON(ctx, "/api/health", nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// Mirrors returns the list of active TorrentClaw mirrors and access channels.
// This endpoint does not require authentication.
func (c *Client) Mirrors(ctx context.Context) (*MirrorsResponse, error) {
	var resp MirrorsResponse
	if err := c.doJSON(ctx, "/api/mirrors", nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}
