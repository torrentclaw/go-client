package torrentclaw

import "context"

// DebridCheckCache checks which info hashes are cached in the user's debrid
// service. Requires a Pro tier API key.
func (c *Client) DebridCheckCache(ctx context.Context, provider, debridKey string, infoHashes []string) (*DebridCheckCacheResponse, error) {
	body := DebridCheckCacheRequest{InfoHashes: infoHashes}
	headers := map[string]string{
		headerDebridProvider: provider,
		headerDebridKey:      debridKey,
	}

	var resp DebridCheckCacheResponse
	if err := c.doPost(ctx, "/api/v1/debrid/check-cache", body, &resp, headers); err != nil {
		return nil, err
	}
	return &resp, nil
}

// DebridAddMagnet adds a magnet link to the user's debrid service.
// Requires a Pro tier API key.
func (c *Client) DebridAddMagnet(ctx context.Context, provider, debridKey, infoHash string) (*DebridAddMagnetResponse, error) {
	body := DebridAddMagnetRequest{InfoHash: infoHash}
	headers := map[string]string{
		headerDebridProvider: provider,
		headerDebridKey:      debridKey,
	}

	var resp DebridAddMagnetResponse
	if err := c.doPost(ctx, "/api/v1/debrid/add-magnet", body, &resp, headers); err != nil {
		return nil, err
	}
	return &resp, nil
}
