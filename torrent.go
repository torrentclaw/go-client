package torrentclaw

import (
	"context"
	"fmt"
)

// TorrentDownloadURL returns the URL for downloading a .torrent file by its
// info hash. This method does not make an HTTP request.
func (c *Client) TorrentDownloadURL(infoHash string) string {
	return fmt.Sprintf("%s/api/v1/torrent/%s", c.baseURL, infoHash)
}

// GetTorrentFile downloads the .torrent file for the given info hash and
// returns the raw bytes.
func (c *Client) GetTorrentFile(ctx context.Context, infoHash string) ([]byte, error) {
	path := fmt.Sprintf("/api/v1/torrent/%s", infoHash)
	return c.doRaw(ctx, path, nil)
}
