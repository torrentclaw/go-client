package torrentclaw

import "context"

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

// Stats returns aggregator statistics including content counts, torrent
// counts by source, and recent ingestion history.
func (c *Client) Stats(ctx context.Context) (*StatsResponse, error) {
	var resp StatsResponse
	if err := c.doJSON(ctx, "/api/v1/stats", nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}
