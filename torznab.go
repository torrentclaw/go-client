package torrentclaw

import (
	"context"
	"net/url"
)

// TorznabParams holds the parameters for a Torznab API request.
type TorznabParams struct {
	// T is the request type: "caps", "search", "tvsearch", "movie".
	T string

	// Q is the search query.
	Q string

	// IMDbID is an IMDb identifier (e.g. "tt1234567").
	IMDbID string

	// TMDbID is a TMDB identifier.
	TMDbID string

	// Season is the season number for TV searches.
	Season int

	// Ep is the episode number for TV searches.
	Ep int

	// Cat is a comma-separated list of category codes (2000=Movies, 5000=TV).
	Cat string

	// Limit sets the max results (1-100, default 50).
	Limit int

	// Offset is the 0-based pagination offset.
	Offset int
}

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
