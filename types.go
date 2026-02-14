// Package torrentclaw provides a Go client for the TorrentClaw API,
// a torrent search engine that aggregates movies and TV shows from 30+
// international sources with TMDB metadata enrichment.
package torrentclaw

// TorrentInfo contains metadata about a single torrent.
type TorrentInfo struct {
	InfoHash     string   `json:"infoHash"`
	Quality      *string  `json:"quality,omitempty"`
	Codec        *string  `json:"codec,omitempty"`
	SourceType   *string  `json:"sourceType,omitempty"`
	SizeBytes    *string  `json:"sizeBytes,omitempty"`
	Seeders      int      `json:"seeders"`
	Leechers     int      `json:"leechers"`
	MagnetURL    *string  `json:"magnetUrl,omitempty"`
	Source       string   `json:"source"`
	QualityScore *int     `json:"qualityScore,omitempty"`
	UploadedAt   *string  `json:"uploadedAt,omitempty"`
	Languages    []string `json:"languages"`
	AudioCodec   *string  `json:"audioCodec,omitempty"`
	HDRType      *string  `json:"hdrType,omitempty"`
	ReleaseGroup *string  `json:"releaseGroup,omitempty"`
	IsProper     *bool    `json:"isProper,omitempty"`
	IsRepack     *bool    `json:"isRepack,omitempty"`
	IsRemastered *bool    `json:"isRemastered,omitempty"`
}

// StreamingProviderItem represents a streaming service provider.
type StreamingProviderItem struct {
	ProviderID      int     `json:"providerId"`
	Name            string  `json:"name"`
	Logo            *string `json:"logo,omitempty"`
	Link            *string `json:"link,omitempty"`
	DisplayPriority int     `json:"displayPriority,omitempty"`
}

// StreamingInfo contains streaming availability grouped by type.
type StreamingInfo struct {
	Flatrate []StreamingProviderItem `json:"flatrate,omitempty"`
	Rent     []StreamingProviderItem `json:"rent,omitempty"`
	Buy      []StreamingProviderItem `json:"buy,omitempty"`
	Free     []StreamingProviderItem `json:"free,omitempty"`
}

// SearchResult represents a single content result from a search.
type SearchResult struct {
	ID            int            `json:"id"`
	IMDbID        *string        `json:"imdbId,omitempty"`
	TMDbID        *string        `json:"tmdbId,omitempty"`
	ContentType   string         `json:"contentType"`
	Title         string         `json:"title"`
	TitleOriginal *string        `json:"titleOriginal,omitempty"`
	Year          *int           `json:"year,omitempty"`
	Overview      *string        `json:"overview,omitempty"`
	PosterURL     *string        `json:"posterUrl,omitempty"`
	BackdropURL   *string        `json:"backdropUrl,omitempty"`
	Genres        []string       `json:"genres,omitempty"`
	RatingIMDb    *string        `json:"ratingImdb,omitempty"`
	RatingTMDb    *string        `json:"ratingTmdb,omitempty"`
	HasTorrents   bool           `json:"hasTorrents"`
	Torrents      []TorrentInfo  `json:"torrents"`
	Streaming     *StreamingInfo `json:"streaming,omitempty"`
}

// SearchResponse is the paginated response from the search endpoint.
type SearchResponse struct {
	Total    int            `json:"total"`
	Page     int            `json:"page"`
	PageSize int            `json:"pageSize"`
	Results  []SearchResult `json:"results"`
}

// SearchParams holds the parameters for a search request.
type SearchParams struct {
	// Query is the search query (required, max 200 characters).
	Query string

	// Type filters by content type: "movie" or "show".
	Type string

	// Genre filters by genre (exact match), e.g. "Action", "Comedy".
	Genre string

	// YearMin sets the minimum release year.
	YearMin int

	// YearMax sets the maximum release year.
	YearMax int

	// MinRating sets the minimum IMDb/TMDB rating threshold (0-10).
	MinRating float64

	// Quality filters by video resolution: "480p", "720p", "1080p", "2160p".
	Quality string

	// Language filters by torrent audio language (ISO 639 code, e.g. "en", "es").
	Language string

	// Audio filters by audio codec (substring match, e.g. "aac", "atmos").
	Audio string

	// HDR filters by HDR format (e.g. "hdr10", "dolby_vision").
	HDR string

	// Sort sets the sort order: "relevance", "seeders", "year", "rating", "added".
	Sort string

	// Page is the page number (starts at 1).
	Page int

	// Limit sets the number of results per page (1-50).
	Limit int

	// Country is an ISO 3166-1 country code to include streaming availability.
	Country string

	// Locale sets the language for title/overview translations (e.g. "es", "fr").
	Locale string

	// Availability filters by torrent availability: "all", "available", "unavailable".
	Availability string
}

// PopularItem represents a popular content item.
type PopularItem struct {
	ID          int     `json:"id"`
	Title       string  `json:"title"`
	Year        *int    `json:"year,omitempty"`
	ContentType string  `json:"contentType"`
	PosterURL   *string `json:"posterUrl,omitempty"`
	RatingIMDb  *string `json:"ratingImdb,omitempty"`
	RatingTMDb  *string `json:"ratingTmdb,omitempty"`
	ClickCount  int     `json:"clickCount"`
}

// PopularResponse is the paginated response from the popular endpoint.
type PopularResponse struct {
	Items    []PopularItem `json:"items"`
	Total    int           `json:"total"`
	Page     int           `json:"page"`
	PageSize int           `json:"pageSize"`
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
	CreatedAt   string  `json:"createdAt"`
}

// RecentResponse is the paginated response from the recent endpoint.
type RecentResponse struct {
	Items    []RecentItem `json:"items"`
	Total    int          `json:"total"`
	Page     int          `json:"page"`
	PageSize int          `json:"pageSize"`
}

// CastMember represents an actor in a movie or TV show.
type CastMember struct {
	Name       string  `json:"name"`
	Character  string  `json:"character"`
	ProfileURL *string `json:"profileUrl,omitempty"`
}

// CreditsResponse contains the director and cast for a content item.
type CreditsResponse struct {
	ContentID int          `json:"contentId"`
	Director  *string      `json:"director,omitempty"`
	Cast      []CastMember `json:"cast"`
}

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

// ContentStats contains counts for movies, shows, and TMDB enrichment.
type ContentStats struct {
	Movies       int `json:"movies"`
	Shows        int `json:"shows"`
	TMDbEnriched int `json:"tmdbEnriched"`
}

// TorrentStats contains aggregate torrent statistics.
type TorrentStats struct {
	Total       int            `json:"total"`
	WithSeeders int            `json:"withSeeders"`
	BySource    map[string]int `json:"bySource"`
}

// IngestionRecord represents a recent data ingestion event.
type IngestionRecord struct {
	Source      string  `json:"source"`
	Status      string  `json:"status"`
	StartedAt   string  `json:"startedAt"`
	CompletedAt *string `json:"completedAt,omitempty"`
	Fetched     int     `json:"fetched"`
	New         int     `json:"new"`
	Updated     int     `json:"updated"`
}

// StatsResponse contains aggregator statistics.
type StatsResponse struct {
	Content          ContentStats      `json:"content"`
	Torrents         TorrentStats      `json:"torrents"`
	RecentIngestions []IngestionRecord `json:"recentIngestions"`
}

// AutocompleteResponse contains autocomplete suggestions.
type AutocompleteResponse []AutocompleteSuggestion

// AutocompleteSuggestion represents a single autocomplete result.
type AutocompleteSuggestion struct {
	ID          int     `json:"id"`
	Title       string  `json:"title"`
	Year        *int    `json:"year,omitempty"`
	ContentType string  `json:"contentType"`
	PosterURL   *string `json:"posterUrl,omitempty"`
}
