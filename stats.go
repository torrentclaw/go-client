package torrentclaw

import "context"

// Stats returns aggregator statistics including content counts, torrent
// counts by source, and recent ingestion history.
func (c *Client) Stats(ctx context.Context) (*StatsResponse, error) {
	var resp StatsResponse
	if err := c.doJSON(ctx, "/api/v1/stats", nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}
