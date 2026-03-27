package torrentclaw

import "context"

// DebridCheckCacheRequest is the request body for checking debrid cache.
type DebridCheckCacheRequest struct {
	InfoHashes []string `json:"infoHashes"`
}

// DebridCheckCacheResponse contains the cache status for each info hash.
type DebridCheckCacheResponse struct {
	Cached map[string]bool `json:"cached"`
}

// DebridAddMagnetRequest is the request body for adding a magnet to debrid.
type DebridAddMagnetRequest struct {
	InfoHash string `json:"infoHash"`
}

// DebridAddMagnetResponse contains the result of adding a magnet.
type DebridAddMagnetResponse struct {
	ID     string `json:"id"`
	Cached bool   `json:"cached"`
	Name   string `json:"name,omitempty"`
}

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
