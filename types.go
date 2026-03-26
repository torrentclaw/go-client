// Package torrentclaw provides a Go client for the TorrentClaw API,
// a torrent search engine that aggregates movies and TV shows from 30+
// international sources with TMDB metadata enrichment.
package torrentclaw

// ─── TrueSpec media metadata ────

// AudioTrack represents a single audio track detected by TrueSpec scanning.
type AudioTrack struct {
	Lang     string `json:"lang"`
	Codec    string `json:"codec"`
	Channels int    `json:"channels"`
	Title    string `json:"title,omitempty"`
	Default  bool   `json:"default,omitempty"`
}

// SubtitleTrack represents a single subtitle track detected by TrueSpec scanning.
type SubtitleTrack struct {
	Lang    string `json:"lang"`
	Codec   string `json:"codec"`
	Title   string `json:"title,omitempty"`
	Forced  bool   `json:"forced,omitempty"`
	Default bool   `json:"default,omitempty"`
}

// VideoInfo contains video stream metadata detected by TrueSpec scanning.
type VideoInfo struct {
	Codec     string   `json:"codec"`
	Width     int      `json:"width"`
	Height    int      `json:"height"`
	BitDepth  *int     `json:"bitDepth,omitempty"`
	HDR       *string  `json:"hdr,omitempty"`
	FrameRate *float64 `json:"frameRate,omitempty"`
	Profile   *string  `json:"profile,omitempty"`
}

// TorrentFileInfo describes a single file inside a torrent.
type TorrentFileInfo struct {
	Path string `json:"path"`
	Size int64  `json:"size"`
	Ext  string `json:"ext"`
}

// SuspiciousFileInfo describes a file flagged as potentially dangerous.
type SuspiciousFileInfo struct {
	Path           string  `json:"path"`
	Size           int64   `json:"size"`
	Ext            string  `json:"ext"`
	Reason         string  `json:"reason"`
	VTPermalink    *string `json:"vtPermalink,omitempty"`
	VTDetections   *int    `json:"vtDetections,omitempty"`
	VTTotalEngines *int    `json:"vtTotalEngines,omitempty"`
	VTStatus       *string `json:"vtStatus,omitempty"`
}

// TorrentFilesInfo contains the file listing and threat analysis for a torrent.
type TorrentFilesInfo struct {
	Total       int                  `json:"total"`
	TotalSize   int64                `json:"totalSize"`
	VideoFiles  []TorrentFileInfo    `json:"videoFiles"`
	AudioFiles  []TorrentFileInfo    `json:"audioFiles"`
	SubFiles    []TorrentFileInfo    `json:"subFiles"`
	ImageFiles  []TorrentFileInfo    `json:"imageFiles"`
	OtherFiles  []TorrentFileInfo    `json:"otherFiles"`
	Suspicious  []SuspiciousFileInfo `json:"suspicious"`
	ThreatLevel string               `json:"threatLevel"`
}

// ─── Torrent ────

// TorrentInfo contains metadata about a single torrent.
type TorrentInfo struct {
	InfoHash          string            `json:"infoHash"`
	RawTitle          string            `json:"rawTitle"`
	Quality           *string           `json:"quality,omitempty"`
	Codec             *string           `json:"codec,omitempty"`
	SourceType        *string           `json:"sourceType,omitempty"`
	SizeBytes         *string           `json:"sizeBytes,omitempty"`
	Seeders           int               `json:"seeders"`
	Leechers          int               `json:"leechers"`
	MagnetURL         *string           `json:"magnetUrl,omitempty"`
	Source            string            `json:"source"`
	QualityScore      *int              `json:"qualityScore,omitempty"`
	UploadedAt        *string           `json:"uploadedAt,omitempty"`
	Languages         []string          `json:"languages"`
	AudioCodec        *string           `json:"audioCodec,omitempty"`
	AudioTracks       []AudioTrack      `json:"audioTracks,omitempty"`
	SubtitleTracks    []SubtitleTrack   `json:"subtitleTracks,omitempty"`
	VideoInfo         *VideoInfo        `json:"videoInfo,omitempty"`
	ScanStatus        *string           `json:"scanStatus,omitempty"`
	ThreatLevel       *string           `json:"threatLevel,omitempty"`
	TorrentFiles      *TorrentFilesInfo `json:"torrentFiles,omitempty"`
	HDRType           *string           `json:"hdrType,omitempty"`
	ReleaseGroup      *string           `json:"releaseGroup,omitempty"`
	IsProper          *bool             `json:"isProper,omitempty"`
	IsRepack          *bool             `json:"isRepack,omitempty"`
	IsRemastered      *bool             `json:"isRemastered,omitempty"`
	Season            *int              `json:"season,omitempty"`
	Episode           *int              `json:"episode,omitempty"`
	SubtitleLanguages []string          `json:"subtitleLanguages"`
}

// ─── Streaming ────

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

// ─── Search ────

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
	ContentURL    string         `json:"contentUrl"`
	HasTorrents   bool           `json:"hasTorrents"`
	Torrents      []TorrentInfo  `json:"torrents"`
	Streaming     *StreamingInfo `json:"streaming,omitempty"`
}

// SearchResponse is the paginated response from the search endpoint.
type SearchResponse struct {
	Total         int            `json:"total"`
	Page          int            `json:"page"`
	PageSize      int            `json:"pageSize"`
	ParsedSeason  *int           `json:"parsedSeason,omitempty"`
	ParsedEpisode *int           `json:"parsedEpisode,omitempty"`
	FuzzyMatch    *bool          `json:"fuzzyMatch,omitempty"`
	Results       []SearchResult `json:"results"`
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

	// Subs filters by subtitle language (ISO 639 code, e.g. "en", "es").
	Subs string

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

	// Verified filters to only TrueSpec verified releases.
	Verified bool

	// Season filters by TV show season number.
	Season int

	// Episode filters by TV show episode number.
	Episode int
}

// ─── Autocomplete ────

// AutocompleteParams holds the parameters for an autocomplete request.
type AutocompleteParams struct {
	// Query is the search prefix (required, 2-200 characters).
	Query string

	// Locale sets the language for localized suggestions.
	Locale string
}

// AutocompleteSuggestion represents a single autocomplete result.
type AutocompleteSuggestion struct {
	ID          int     `json:"id"`
	Title       string  `json:"title"`
	Year        *int    `json:"year,omitempty"`
	ContentType string  `json:"contentType"`
	PosterURL   *string `json:"posterUrl,omitempty"`
	MovieCount  *int    `json:"movieCount,omitempty"`
}

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

// ─── Stats ────

// ContentStats contains counts for movies, shows, and TMDB enrichment.
type ContentStats struct {
	Movies       int `json:"movies"`
	Shows        int `json:"shows"`
	TMDbEnriched int `json:"tmdbEnriched"`
}

// TorrentStats contains aggregate torrent statistics.
type TorrentStats struct {
	Total        int            `json:"total"`
	WithSeeders  int            `json:"withSeeders"`
	Orphans      int            `json:"orphans"`
	DailyAverage int            `json:"dailyAverage"`
	BySource     map[string]int `json:"bySource"`
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

// ─── Trending ────

// TrendingParams holds the parameters for a trending content request.
type TrendingParams struct {
	// Period sets the time window: "daily", "weekly", "monthly" (default "daily").
	Period string

	// Limit sets the number of items (1-50, default 20).
	Limit int

	// Page is the page number (starts at 1).
	Page int

	// Locale sets the language for localized titles.
	Locale string
}

// TrendingItem represents a trending content item.
type TrendingItem struct {
	ID          int     `json:"id"`
	Title       string  `json:"title"`
	Year        *int    `json:"year,omitempty"`
	ContentType string  `json:"contentType"`
	PosterURL   *string `json:"posterUrl,omitempty"`
	RatingIMDb  *string `json:"ratingImdb,omitempty"`
	RatingTMDb  *string `json:"ratingTmdb,omitempty"`
	Overview    *string `json:"overview,omitempty"`
	MaxSeeders  int     `json:"maxSeeders"`
	ClickCount  int     `json:"clickCount"`
	TrendScore  int     `json:"trendScore"`
}

// TrendingResponse is the paginated response from the trending endpoint.
type TrendingResponse struct {
	Period   string         `json:"period"`
	Items    []TrendingItem `json:"items"`
	Total    int            `json:"total"`
	Page     int            `json:"page"`
	PageSize int            `json:"pageSize"`
}

// ─── Upcoming ────

// UpcomingParams holds the parameters for an upcoming releases request.
type UpcomingParams struct {
	// Limit sets the number of items (1-50, default 24).
	Limit int

	// Page is the page number (starts at 1).
	Page int

	// Type filters by content type: "all", "movie", "show" (default "all").
	Type string

	// Locale sets the language for localized titles.
	Locale string
}

// UpcomingItem represents an upcoming release.
type UpcomingItem struct {
	ID             int      `json:"id"`
	TMDbID         *int     `json:"tmdbId,omitempty"`
	Title          string   `json:"title"`
	Year           *int     `json:"year,omitempty"`
	ContentType    string   `json:"contentType"`
	PosterURL      *string  `json:"posterUrl,omitempty"`
	BackdropURL    *string  `json:"backdropUrl,omitempty"`
	Overview       *string  `json:"overview,omitempty"`
	Genres         []string `json:"genres,omitempty"`
	RatingIMDb     *string  `json:"ratingImdb,omitempty"`
	RatingTMDb     *string  `json:"ratingTmdb,omitempty"`
	PopularityTMDb *float64 `json:"popularityTmdb,omitempty"`
	ReleaseDate    string   `json:"releaseDate"`
	HasTorrents    bool     `json:"hasTorrents"`
}

// UpcomingResponse is the paginated response from the upcoming endpoint.
type UpcomingResponse struct {
	Items    []UpcomingItem `json:"items"`
	Total    int            `json:"total"`
	Page     int            `json:"page"`
	PageSize int            `json:"pageSize"`
}

// ─── Collections ────

// CollectionListParams holds the parameters for listing collections.
type CollectionListParams struct {
	// Limit sets the number of items (1-48, default 24).
	Limit int

	// Page is the page number (starts at 1).
	Page int

	// Locale sets the language for localized names.
	Locale string
}

// CollectionListItem represents a movie collection in a list.
type CollectionListItem struct {
	ID           int     `json:"id"`
	TMDbID       int     `json:"tmdbId"`
	Name         string  `json:"name"`
	PosterURL    *string `json:"posterUrl,omitempty"`
	BackdropURL  *string `json:"backdropUrl,omitempty"`
	MovieCount   int     `json:"movieCount"`
	TotalSeeders int     `json:"totalSeeders"`
	PartCount    int     `json:"partCount"`
}

// CollectionListResponse is the paginated response from the collections endpoint.
type CollectionListResponse struct {
	Items    []CollectionListItem `json:"items"`
	Total    int                  `json:"total"`
	Page     int                  `json:"page"`
	PageSize int                  `json:"pageSize"`
}

// CollectionDetail contains full details about a movie collection.
type CollectionDetail struct {
	ID           int           `json:"id"`
	TMDbID       int           `json:"tmdbId"`
	Name         string        `json:"name"`
	PosterURL    *string       `json:"posterUrl,omitempty"`
	BackdropURL  *string       `json:"backdropUrl,omitempty"`
	MovieCount   int           `json:"movieCount"`
	TotalSeeders int           `json:"totalSeeders"`
	PartCount    int           `json:"partCount"`
	Overview     *string       `json:"overview,omitempty"`
	Movies       []PopularItem `json:"movies"`
}

// ─── Streaming Top ────

// StreamingTopParams holds the parameters for a streaming top 10 request.
type StreamingTopParams struct {
	// Service is the streaming service: "netflix", "prime", "disney", "hbo", "apple".
	Service string

	// Country is an ISO 3166-1 country code (default "US").
	Country string

	// ShowType is the content type: "movie" or "series" (default "movie").
	ShowType string

	// Locale sets the language for localized titles.
	Locale string
}

// StreamingTopItem represents a ranked item in a streaming top list.
type StreamingTopItem struct {
	Rank        int      `json:"rank"`
	Title       string   `json:"title"`
	IMDbID      *string  `json:"imdbId,omitempty"`
	TMDbID      *int     `json:"tmdbId,omitempty"`
	ContentType *string  `json:"contentType,omitempty"`
	Year        *int     `json:"year,omitempty"`
	Overview    *string  `json:"overview,omitempty"`
	Rating      *string  `json:"rating,omitempty"`
	PosterURL   *string  `json:"posterUrl,omitempty"`
	BackdropURL *string  `json:"backdropUrl,omitempty"`
	Genres      []string `json:"genres,omitempty"`
	StreamingLink *string `json:"streamingLink,omitempty"`
	ContentID   *int     `json:"contentId,omitempty"`
	HasTorrents bool     `json:"hasTorrents"`
	MaxSeeders  int      `json:"maxSeeders"`
}

// ─── Debrid ────

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

// ─── Torznab ────

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

// ─── Health & Mirrors ────

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
