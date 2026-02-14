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

	// Create a client with custom options.
	_ = torrentclaw.NewClient(
		torrentclaw.WithAPIKey("your-api-key"),
		torrentclaw.WithTimeout(30*time.Second),
		torrentclaw.WithRetry(5, 2*time.Second, 60*time.Second),
	)
}

func ExampleClient_Search() {
	client := torrentclaw.NewClient()

	resp, err := client.Search(context.Background(), torrentclaw.SearchParams{
		Query:    "Breaking Bad",
		Type:     "show",
		YearMin:  2008,
		Language: "en",
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

	resp, err := client.Popular(context.Background(), 10, 1)
	if err != nil {
		log.Fatal(err)
	}

	for _, item := range resp.Items {
		fmt.Printf("%s — %d clicks\n", item.Title, item.ClickCount)
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

	suggestions, err := client.Autocomplete(context.Background(), "incep")
	if err != nil {
		log.Fatal(err)
	}

	for _, s := range suggestions {
		fmt.Printf("%s (%s)\n", s.Title, s.ContentType)
	}
}

func ExampleClient_Recent() {
	client := torrentclaw.NewClient()

	resp, err := client.Recent(context.Background(), 10, 1)
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
