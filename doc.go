// Package torrentclaw provides a Go client for the TorrentClaw API,
// a torrent search engine that aggregates movies and TV shows from 30+
// international sources with TMDB metadata enrichment.
//
// Each file in this package is self-contained: it holds the types (params,
// response structs) and methods for a single API resource. Shared types
// that are referenced across multiple resources live in types.go.
//
// Usage:
//
//	client := torrentclaw.NewClient(
//		torrentclaw.WithAPIKey("your-api-key"),
//	)
//
//	resp, err := client.Search(context.Background(), torrentclaw.SearchParams{
//		Query: "Inception",
//		Type:  "movie",
//	})
package torrentclaw
