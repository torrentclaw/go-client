package torrentclaw

import (
	"context"
	"net/url"
)

// Torznab queries the Torznab-compatible API and returns the raw XML response.
// Requires a Pro tier API key.
func (c *Client) Torznab(ctx context.Context, params TorznabParams) ([]byte, error) {
	q := url.Values{}
	addStringParam(q, "t", params.T)
	addStringParam(q, "q", params.Q)
	addStringParam(q, "imdbid", params.IMDbID)
	addStringParam(q, "tmdbid", params.TMDbID)
	addIntParam(q, "season", params.Season)
	addIntParam(q, "ep", params.Ep)
	addStringParam(q, "cat", params.Cat)
	addIntParam(q, "limit", params.Limit)
	addIntParam(q, "offset", params.Offset)

	return c.doRaw(ctx, "/api/v1/torznab", q)
}

// TorznabCaps returns the Torznab capabilities XML.
// Requires a Pro tier API key.
func (c *Client) TorznabCaps(ctx context.Context) ([]byte, error) {
	return c.Torznab(ctx, TorznabParams{T: "caps"})
}
