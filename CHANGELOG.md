# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [0.1.0] - 2025-01-15

### Added

- Initial release of the TorrentClaw Go client library.
- `Search` — full-text search with advanced filtering (type, genre, year, quality, language, audio, HDR, sort, pagination, country, locale, availability).
- `Autocomplete` — title suggestions for search-as-you-type.
- `Popular` — trending content by community engagement.
- `Recent` — recently added movies and TV shows.
- `WatchProviders` — streaming availability (flatrate, rent, buy, free) with VPN suggestions.
- `Credits` — director and top cast members.
- `Stats` — aggregator statistics (content counts, torrent counts, ingestion history).
- `GetTorrentFile` — download raw `.torrent` file bytes.
- `TorrentDownloadURL` — construct download URL without making an HTTP call.
- Functional options pattern for client configuration (`WithAPIKey`, `WithBaseURL`, `WithTimeout`, `WithRetry`, `WithHTTPClient`, `WithUserAgent`).
- Exponential backoff retry for transient errors (429, 5xx).
- Custom `APIError` type with helper methods (`IsRetryable`, `IsRateLimited`, `IsNotFound`).
- Context support on all methods.
- Zero external dependencies (stdlib only).
- Comprehensive test suite with `httptest`.
- Example tests for godoc.
- CI workflow with GitHub Actions.
