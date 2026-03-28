package torrentclaw

import (
	"context"
	"fmt"
	"net/url"
)

// ─── Popular ────

// PopularParams holds the parameters for a popular content request.
type PopularParams struct {
	// Limit sets the number of items (1-24, default 12).
	Limit int

	// Page is the page number (starts at 1).
	Page int

	// Locale sets the language for localized titles.
	Locale string
}

// PopularItem represents a popular content item.
type PopularItem struct {
	ID          int      `json:"id"`
	Title       string   `json:"title"`
	Year        *int     `json:"year,omitempty"`
	ContentType string   `json:"contentType"`
	PosterURL   *string  `json:"posterUrl,omitempty"`
	RatingIMDb  *string  `json:"ratingImdb,omitempty"`
	RatingTMDb  *string  `json:"ratingTmdb,omitempty"`
	Overview    *string  `json:"overview,omitempty"`
	MaxSeeders  int      `json:"maxSeeders"`
	Genres      []string `json:"genres,omitempty"`
	BestQuality *string  `json:"bestQuality,omitempty"`
	HasHDR      *bool    `json:"hasHdr,omitempty"`
	TopInfoHash *string  `json:"topInfoHash,omitempty"`
}

// PopularResponse is the paginated response from the popular endpoint.
type PopularResponse struct {
	Items    []PopularItem `json:"items"`
	Total    int           `json:"total"`
	Page     int           `json:"page"`
	PageSize int           `json:"pageSize"`
}

// ─── Recent ────

// RecentParams holds the parameters for a recent content request.
type RecentParams struct {
	// Limit sets the number of items (1-24, default 12).
	Limit int

	// Page is the page number (starts at 1).
	Page int

	// Locale sets the language for localized titles.
	Locale string
}

// RecentItem represents a recently added content item.
type RecentItem struct {
	ID          int     `json:"id"`
	Title       string  `json:"title"`
	Year        *int    `json:"year,omitempty"`
	ContentType string  `json:"contentType"`
	PosterURL   *string `json:"posterUrl,omitempty"`
	RatingIMDb  *string `json:"ratingImdb,omitempty"`
	RatingTMDb  *string `json:"ratingTmdb,omitempty"`
	Overview    *string `json:"overview,omitempty"`
	CreatedAt   string  `json:"createdAt"`
}

// RecentResponse is the paginated response from the recent endpoint.
type RecentResponse struct {
	Items    []RecentItem `json:"items"`
	Total    int          `json:"total"`
	Page     int          `json:"page"`
	PageSize int          `json:"pageSize"`
}

// ─── Watch Providers ────

// WatchProviderItem represents a single watch/streaming provider.
type WatchProviderItem struct {
	ProviderID      int     `json:"providerId"`
	Name            string  `json:"name"`
	Logo            *string `json:"logo,omitempty"`
	Link            *string `json:"link,omitempty"`
	DisplayPriority int     `json:"displayPriority,omitempty"`
}

// WatchProviders contains streaming availability grouped by access type.
type WatchProviders struct {
	Flatrate []WatchProviderItem `json:"flatrate,omitempty"`
	Rent     []WatchProviderItem `json:"rent,omitempty"`
	Buy      []WatchProviderItem `json:"buy,omitempty"`
	Free     []WatchProviderItem `json:"free,omitempty"`
}

// VPNSuggestion is included when content is available in other countries
// via subscription but not in the user's country.
type VPNSuggestion struct {
	AvailableIn  []string `json:"availableIn"`
	AffiliateURL string   `json:"affiliateUrl"`
}

// WatchProvidersResponse contains watch providers for a content item.
type WatchProvidersResponse struct {
	ContentID     int            `json:"contentId"`
	Country       string         `json:"country"`
	Providers     WatchProviders `json:"providers"`
	VPNSuggestion *VPNSuggestion `json:"vpnSuggestion,omitempty"`
	Attribution   string         `json:"attribution"`
}

// ─── Credits ────

// CastMember represents an actor in a movie or TV show.
type CastMember struct {
	TmdbID     *int    `json:"tmdbId,omitempty"`
	Name       string  `json:"name"`
	Character  string  `json:"character"`
	ProfileURL *string `json:"profileUrl,omitempty"`
}

// CreditsResponse contains the director and cast for a content item.
type CreditsResponse struct {
	ContentID      int          `json:"contentId"`
	Director       *string      `json:"director,omitempty"`
	DirectorTmdbID *int         `json:"directorTmdbId,omitempty"`
	Cast           []CastMember `json:"cast"`
}

// ─── Methods ────

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
