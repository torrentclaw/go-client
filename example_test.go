package torrentclaw_test

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/torrentclaw/torrentclaw-go-client"
)

func Example() {
	client := torrentclaw.NewClient(
		torrentclaw.WithAPIKey("your-api-key"),
	)

	resp, err := client.Search(context.Background(), torrentclaw.SearchParams{
		Query:   "Inception",
		Type:    "movie",
		Quality: "1080p",
		Sort:    "seeders",
	})
	if err != nil {
		log.Fatal(err)
	}

	for _, result := range resp.Results {
		fmt.Printf("%s (%s)\n", result.Title, result.ContentType)
		for _, t := range result.Torrents {
			fmt.Printf("  %s — %d seeders\n", t.InfoHash, t.Seeders)
		}
	}
}

func ExampleNewClient() {
	// Create a client with default settings.
	_ = torrentclaw.NewClient()

	// Create a client with API key.
	_ = torrentclaw.NewClient(
		torrentclaw.WithAPIKey("your-api-key"),
		torrentclaw.WithTimeout(30*time.Second),
		torrentclaw.WithRetry(5, 2*time.Second, 60*time.Second),
	)

	// Create a client with bearer token.
	_ = torrentclaw.NewClient(
		torrentclaw.WithBearerToken("your-bearer-token"),
	)
}

func ExampleClient_Search() {
	client := torrentclaw.NewClient()

	resp, err := client.Search(context.Background(), torrentclaw.SearchParams{
		Query:    "Breaking Bad",
		Type:     "show",
		YearMin:  2008,
		Language: "en",
		Season:   1,
		Episode:  5,
		Page:     1,
		Limit:    5,
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Found %d results\n", resp.Total)
}

func ExampleClient_Popular() {
	client := torrentclaw.NewClient()

	resp, err := client.Popular(context.Background(), torrentclaw.PopularParams{
		Limit:  10,
		Page:   1,
		Locale: "es",
	})
	if err != nil {
		log.Fatal(err)
	}

	for _, item := range resp.Items {
		fmt.Printf("%s — %d seeders\n", item.Title, item.MaxSeeders)
	}
}

func ExampleClient_Credits() {
	client := torrentclaw.NewClient()

	credits, err := client.Credits(context.Background(), 42)
	if err != nil {
		log.Fatal(err)
	}

	if credits.Director != nil {
		fmt.Printf("Director: %s\n", *credits.Director)
	}
	for _, member := range credits.Cast {
		fmt.Printf("%s as %s\n", member.Name, member.Character)
	}
}

func ExampleClient_TorrentDownloadURL() {
	client := torrentclaw.NewClient()

	url := client.TorrentDownloadURL("abc123def456789012345678901234567890abcd")
	fmt.Println(url)
	// Output: https://torrentclaw.com/api/v1/torrent/abc123def456789012345678901234567890abcd
}

func ExampleClient_Autocomplete() {
	client := torrentclaw.NewClient()

	suggestions, err := client.Autocomplete(context.Background(), torrentclaw.AutocompleteParams{
		Query:  "incep",
		Locale: "es",
	})
	if err != nil {
		log.Fatal(err)
	}

	for _, s := range suggestions {
		fmt.Printf("%s (%s)\n", s.Title, s.ContentType)
	}
}

func ExampleClient_Recent() {
	client := torrentclaw.NewClient()

	resp, err := client.Recent(context.Background(), torrentclaw.RecentParams{
		Limit: 10,
		Page:  1,
	})
	if err != nil {
		log.Fatal(err)
	}

	for _, item := range resp.Items {
		fmt.Printf("%s — added %s\n", item.Title, item.CreatedAt)
	}
}

func ExampleClient_WatchProviders() {
	client := torrentclaw.NewClient()

	resp, err := client.WatchProviders(context.Background(), 42, "US")
	if err != nil {
		log.Fatal(err)
	}

	for _, p := range resp.Providers.Flatrate {
		fmt.Printf("Stream on: %s\n", p.Name)
	}
	if resp.VPNSuggestion != nil {
		fmt.Printf("Also available via VPN in: %v\n", resp.VPNSuggestion.AvailableIn)
	}
}

func ExampleClient_Stats() {
	client := torrentclaw.NewClient()

	stats, err := client.Stats(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Movies: %d, Shows: %d\n", stats.Content.Movies, stats.Content.Shows)
	fmt.Printf("Torrents: %d (with seeders: %d)\n", stats.Torrents.Total, stats.Torrents.WithSeeders)
}

func ExampleClient_GetTorrentFile() {
	client := torrentclaw.NewClient()

	data, err := client.GetTorrentFile(context.Background(), "abc123def456")
	if err != nil {
		log.Fatal(err)
	}

	// data contains the raw .torrent file bytes
	_ = data
}

func ExampleWithRetry() {
	// Disable retries entirely.
	_ = torrentclaw.NewClient(torrentclaw.WithRetry(0, 0, 0))

	// Custom: 5 retries, starting at 2s, capped at 60s.
	_ = torrentclaw.NewClient(
		torrentclaw.WithRetry(5, 2*time.Second, 60*time.Second),
	)
}

func ExampleClient_Trending() {
	client := torrentclaw.NewClient()

	resp, err := client.Trending(context.Background(), torrentclaw.TrendingParams{
		Period: "weekly",
		Limit:  10,
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Trending (%s):\n", resp.Period)
	for _, item := range resp.Items {
		fmt.Printf("  %s — score %d\n", item.Title, item.TrendScore)
	}
}

func ExampleClient_Upcoming() {
	client := torrentclaw.NewClient()

	resp, err := client.Upcoming(context.Background(), torrentclaw.UpcomingParams{
		Type:  "movie",
		Limit: 10,
	})
	if err != nil {
		log.Fatal(err)
	}

	for _, item := range resp.Items {
		fmt.Printf("%s — releases %s\n", item.Title, item.ReleaseDate)
	}
}

func ExampleClient_Collections() {
	client := torrentclaw.NewClient()

	resp, err := client.Collections(context.Background(), torrentclaw.CollectionListParams{
		Limit: 12,
	})
	if err != nil {
		log.Fatal(err)
	}

	for _, c := range resp.Items {
		fmt.Printf("%s (%d movies)\n", c.Name, c.MovieCount)
	}
}

func ExampleClient_CollectionByID() {
	client := torrentclaw.NewClient()

	detail, err := client.CollectionByID(context.Background(), 10, "es")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%s — %d movies\n", detail.Name, detail.MovieCount)
	for _, m := range detail.Movies {
		fmt.Printf("  %s\n", m.Title)
	}
}

func ExampleClient_StreamingTop() {
	client := torrentclaw.NewClient()

	items, err := client.StreamingTop(context.Background(), torrentclaw.StreamingTopParams{
		Service:  "netflix",
		Country:  "US",
		ShowType: "movie",
	})
	if err != nil {
		log.Fatal(err)
	}

	for _, item := range items {
		fmt.Printf("#%d %s (torrents: %v)\n", item.Rank, item.Title, item.HasTorrents)
	}
}

func ExampleClient_DebridCheckCache() {
	client := torrentclaw.NewClient(torrentclaw.WithAPIKey("your-pro-key"))

	resp, err := client.DebridCheckCache(context.Background(),
		"real-debrid", "your-debrid-key",
		[]string{"abc123...", "def456..."},
	)
	if err != nil {
		log.Fatal(err)
	}

	for hash, cached := range resp.Cached {
		fmt.Printf("%s: cached=%v\n", hash, cached)
	}
}

func ExampleClient_DebridAddMagnet() {
	client := torrentclaw.NewClient(torrentclaw.WithAPIKey("your-pro-key"))

	resp, err := client.DebridAddMagnet(context.Background(),
		"real-debrid", "your-debrid-key", "abc123...",
	)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Added: %s (cached: %v)\n", resp.Name, resp.Cached)
}

func ExampleClient_Torznab() {
	client := torrentclaw.NewClient(torrentclaw.WithAPIKey("your-pro-key"))

	// Get capabilities
	caps, err := client.TorznabCaps(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	_ = caps // XML bytes

	// Search for a movie
	data, err := client.Torznab(context.Background(), torrentclaw.TorznabParams{
		T:     "movie",
		Q:     "inception",
		Limit: 25,
	})
	if err != nil {
		log.Fatal(err)
	}
	_ = data // XML bytes
}

func ExampleClient_Health() {
	client := torrentclaw.NewClient()

	health, err := client.Health(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Status: %s, Uptime: %ds\n", health.Status, health.Uptime)
}

func ExampleClient_Mirrors() {
	client := torrentclaw.NewClient()

	resp, err := client.Mirrors(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	for _, m := range resp.Mirrors {
		fmt.Printf("%s — %s (primary: %v)\n", m.Label, m.URL, m.Primary)
	}
}
