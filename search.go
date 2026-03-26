package torrentclaw

import (
	"context"
	"net/url"
)

// Search queries the TorrentClaw search endpoint with the given parameters.
// The Query field in params is required.
func (c *Client) Search(ctx context.Context, params SearchParams) (*SearchResponse, error) {
	q := url.Values{}
	q.Set("q", params.Query)
	addStringParam(q, "type", params.Type)
	addStringParam(q, "genre", params.Genre)
	addIntParam(q, "year_min", params.YearMin)
	addIntParam(q, "year_max", params.YearMax)
	addFloatParam(q, "min_rating", params.MinRating)
	addStringParam(q, "quality", params.Quality)
	addStringParam(q, "lang", params.Language)
	addStringParam(q, "subs", params.Subs)
	addStringParam(q, "audio", params.Audio)
	addStringParam(q, "hdr", params.HDR)
	addStringParam(q, "sort", params.Sort)
	addIntParam(q, "page", params.Page)
	addIntParam(q, "limit", params.Limit)
	addStringParam(q, "country", params.Country)
	addStringParam(q, "locale", params.Locale)
	addStringParam(q, "availability", params.Availability)
	addBoolParam(q, "verified", params.Verified)
	addIntParam(q, "season", params.Season)
	addIntParam(q, "episode", params.Episode)

	var resp SearchResponse
	if err := c.doJSON(ctx, "/api/v1/search", q, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// Autocomplete returns title suggestions for the given query prefix.
func (c *Client) Autocomplete(ctx context.Context, params AutocompleteParams) ([]AutocompleteSuggestion, error) {
	q := url.Values{}
	q.Set("q", params.Query)
	addStringParam(q, "locale", params.Locale)

	var resp struct {
		Suggestions []AutocompleteSuggestion `json:"suggestions"`
	}
	if err := c.doJSON(ctx, "/api/v1/autocomplete", q, &resp); err != nil {
		return nil, err
	}
	return resp.Suggestions, nil
}
