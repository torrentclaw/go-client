package torrentclaw

import "context"

// HealthResponse contains the API health status.
type HealthResponse struct {
	Status    string  `json:"status"`
	Timestamp string  `json:"timestamp"`
	Uptime    int     `json:"uptime"`
	Database  *string `json:"database,omitempty"`
	Redis     *string `json:"redis,omitempty"`
}

// MirrorInfo represents an active TorrentClaw mirror instance.
type MirrorInfo struct {
	URL     string `json:"url"`
	Label   string `json:"label"`
	Primary bool   `json:"primary"`
}

// TorInfo represents a Tor (.onion) access point.
type TorInfo struct {
	URL   string `json:"url"`
	Label string `json:"label"`
}

// StatusChannel represents a status/announcement channel.
type StatusChannel struct {
	Label string `json:"label"`
	URL   string `json:"url"`
}

// MirrorsResponse contains the list of mirrors and access channels.
type MirrorsResponse struct {
	Mirrors  []MirrorInfo    `json:"mirrors"`
	Tor      *TorInfo        `json:"tor,omitempty"`
	Lite     *string         `json:"lite,omitempty"`
	Channels []StatusChannel `json:"channels,omitempty"`
}

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
