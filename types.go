package torrentclaw

// This file contains shared types referenced by multiple endpoint files.
// Domain-specific types (params, responses) live alongside their methods
// in the corresponding endpoint file (search.go, content.go, etc.).

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

// ─── Torrent file analysis ────

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
// Used by SearchResult and torrent download endpoints.
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
// Used by StreamingInfo in search results.
type StreamingProviderItem struct {
	ProviderID      int     `json:"providerId"`
	Name            string  `json:"name"`
	Logo            *string `json:"logo,omitempty"`
	Link            *string `json:"link,omitempty"`
	DisplayPriority int     `json:"displayPriority,omitempty"`
}

// StreamingInfo contains streaming availability grouped by type.
// Used by SearchResult when a country filter is applied.
type StreamingInfo struct {
	Flatrate []StreamingProviderItem `json:"flatrate,omitempty"`
	Rent     []StreamingProviderItem `json:"rent,omitempty"`
	Buy      []StreamingProviderItem `json:"buy,omitempty"`
	Free     []StreamingProviderItem `json:"free,omitempty"`
}
