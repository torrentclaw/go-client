# TorrentClaw Go Client

[![Go Reference](https://pkg.go.dev/badge/github.com/torrentclaw/go-client.svg)](https://pkg.go.dev/github.com/torrentclaw/go-client)
[![CI](https://github.com/torrentclaw/go-client/actions/workflows/ci.yml/badge.svg)](https://github.com/torrentclaw/go-client/actions/workflows/ci.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/torrentclaw/go-client)](https://goreportcard.com/report/github.com/torrentclaw/go-client)
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)
[![Discord](https://img.shields.io/badge/Discord-coming%20soon-7289da)](https://torrentclaw.com)

Go client library for the [TorrentClaw](https://torrentclaw.com) API.

## About TorrentClaw

[TorrentClaw](https://torrentclaw.com) is a torrent search engine that aggregates movies and TV shows from **30+ international sources** into a single, clean API. No ads, no tracking.

- **Quality scoring** — torrents ranked by real quality metrics, not just seeders
- **TMDB metadata** — posters, ratings, genres, cast, and crew for every result
- **Watch providers** — see where content is streaming (Netflix, Amazon, Disney+, etc.)
- **Semantic search** — natural language queries in 11+ languages
- **TrueSpec verification** — verify actual media specs from info hashes

### Links

| | |
|---|---|
| **Website** | [torrentclaw.com](https://torrentclaw.com) |
| **API Docs (OpenAPI)** | [torrentclaw.com/api/openapi.json](https://torrentclaw.com/api/openapi.json) |
| **Discord** | Coming soon |
| **GitHub** | [github.com/torrentclaw](https://github.com/torrentclaw) |

## Ecosystem

TorrentClaw is more than just an API. It's a growing ecosystem of tools:

| Project | Language | Description |
|---|---|---|
| [torrentclaw-go-client](https://github.com/torrentclaw/go-client) | Go | API client library **(this repo)** |
| torrentclaw-cli | Go | Command-line interface *(coming soon)* |
| [torrentclaw-mcp](https://github.com/torrentclaw/torrentclaw-mcp) | TypeScript | MCP server for AI agents |
| [truespec](https://github.com/torrentclaw/truespec) | Go | Verify real media specs from info hashes |
| [torrentclaw-skill](https://github.com/torrentclaw/torrentclaw-skill) | - | Agent skill for OpenClaw |

## Features

- **Zero external dependencies** — stdlib only
- **Context support** on all methods
- **Functional options** for client configuration
- **Exponential backoff retry** for transient errors (429, 5xx)
- **Custom error types** with helper methods (`IsRetryable`, `IsRateLimited`, `IsNotFound`)
- **Full API coverage** — search, autocomplete, popular, recent, watch providers, credits, stats, torrent download

## Installation

```bash
go get github.com/torrentclaw/go-client
```

Requires Go 1.22 or later.

## Quick Start

```go
package main

import (
    "context"
    "fmt"
    "log"

    torrentclaw "github.com/torrentclaw/go-client"
)

func main() {
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
        fmt.Printf("%s (%d)\n", result.Title, *result.Year)
        for _, t := range result.Torrents {
            fmt.Printf("  %s — %d seeders\n", t.InfoHash, t.Seeders)
        }
    }
}
```

## Client Configuration

```go
client := torrentclaw.NewClient(
    torrentclaw.WithAPIKey("your-api-key"),           // API key (X-API-Key header)
    torrentclaw.WithBaseURL("https://custom.example"), // Custom base URL
    torrentclaw.WithTimeout(30 * time.Second),         // HTTP timeout (default: 15s)
    torrentclaw.WithUserAgent("my-app/1.0"),           // Custom User-Agent
    torrentclaw.WithHTTPClient(customHTTPClient),      // Custom *http.Client
    torrentclaw.WithRetry(5, 2*time.Second, 60*time.Second), // Retry policy
)
```

## API Reference

### Search

```go
resp, err := client.Search(ctx, torrentclaw.SearchParams{
    Query:        "Breaking Bad",
    Type:         "show",           // "movie" or "show"
    Genre:        "Drama",          // exact genre match
    YearMin:      2008,
    YearMax:      2013,
    MinRating:    8.0,              // 0-10
    Quality:      "1080p",          // "480p", "720p", "1080p", "2160p"
    Language:     "en",             // ISO 639 code
    Audio:        "atmos",          // audio codec substring
    HDR:          "dolby_vision",   // HDR format
    Sort:         "seeders",        // "relevance", "seeders", "year", "rating", "added"
    Page:         1,
    Limit:        20,
    Country:      "US",             // includes streaming availability
    Locale:       "es",             // localized titles/overviews
    Availability: "available",      // "all", "available", "unavailable"
})
```

### Autocomplete

```go
suggestions, err := client.Autocomplete(ctx, "incep")
```

### Popular Content

```go
resp, err := client.Popular(ctx, 10, 1) // limit, page (0 = server default)
```

### Recent Content

```go
resp, err := client.Recent(ctx, 12, 1)
```

### Watch Providers

```go
resp, err := client.WatchProviders(ctx, contentID, "US")
// resp.Providers.Flatrate — subscription (Netflix, Disney+, etc.)
// resp.Providers.Rent     — available for rental
// resp.Providers.Buy      — available for purchase
// resp.Providers.Free     — free with ads
// resp.VPNSuggestion      — available in other countries
```

### Credits

```go
credits, err := client.Credits(ctx, contentID)
// credits.Director — director name
// credits.Cast     — top 10 cast members
```

### Stats

```go
stats, err := client.Stats(ctx)
// stats.Content.Movies, stats.Content.Shows
// stats.Torrents.Total, stats.Torrents.BySource
// stats.RecentIngestions
```

### Torrent File

```go
// Get the download URL (no HTTP call)
url := client.TorrentDownloadURL("abc123...")

// Download the .torrent file
data, err := client.GetTorrentFile(ctx, "abc123...")
```

## Error Handling

All API errors are returned as `*torrentclaw.APIError`:

```go
resp, err := client.Search(ctx, params)
if err != nil {
    var apiErr *torrentclaw.APIError
    if errors.As(err, &apiErr) {
        fmt.Printf("HTTP %d: %s\n", apiErr.StatusCode, apiErr.Message)

        if apiErr.IsRateLimited() {
            // Handle 429
        }
        if apiErr.IsNotFound() {
            // Handle 404
        }
        if apiErr.IsRetryable() {
            // 429, 500, 502, 503 — retries are automatic by default
        }
    }
}
```

Transient errors (429, 500, 502, 503) are automatically retried with exponential backoff. Configure the retry policy with `WithRetry`:

```go
// Disable retries
client := torrentclaw.NewClient(torrentclaw.WithRetry(0, 0, 0))

// Custom: 5 retries, starting at 2s, capped at 60s
client := torrentclaw.NewClient(torrentclaw.WithRetry(5, 2*time.Second, 60*time.Second))
```

## Development

```bash
# Install git hooks (requires lefthook)
make install-hooks

# Run tests
make test

# Run linter
make lint

# Run tests with coverage
make coverage

# Format code
make fmt

# Check formatting (CI-friendly, no write)
make check

# Run all checks
make all
```

See [CONTRIBUTING.md](CONTRIBUTING.md) for full setup instructions including lefthook installation.

## Contributing

Contributions are welcome! Please read our [Contributing Guide](CONTRIBUTING.md) before submitting a pull request.

## Reporting Bugs

Found a bug? [Open an issue](https://github.com/torrentclaw/go-client/issues/new?labels=bug&template=bug_report.md) on GitHub with:

- A clear description of the problem
- Steps to reproduce
- Expected vs actual behavior
- Go version and OS

## Requesting Features

Have an idea? [Open a feature request](https://github.com/torrentclaw/go-client/issues/new?labels=enhancement&template=feature_request.md) on GitHub. We'd love to hear from you.

## License

[MIT](LICENSE)

---

<p align="center">
  Made with ❤️ by the <a href="https://github.com/torrentclaw">TorrentClaw</a> community
  <br>
  <a href="https://torrentclaw.com">Website</a> · <a href="https://github.com/torrentclaw">GitHub</a>
</p>
