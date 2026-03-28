# Contributing to TorrentClaw Go Client

Thank you for your interest in contributing! This guide will help you get started.

## Getting Started

1. **Fork** the repository on GitHub
2. **Clone** your fork locally:
   ```bash
   git clone https://github.com/YOUR-USERNAME/torrentclaw-go-client.git
   cd torrentclaw-go-client
   ```
3. **Create a branch** for your change:
   ```bash
   git checkout -b feature/my-feature
   ```
4. **Make your changes**, write tests, and ensure everything passes
5. **Commit** with a clear message (see [Commit Messages](#commit-messages))
6. **Push** to your fork and [open a Pull Request](https://github.com/torrentclaw/go-client/compare)

## Development Setup

You need **Go 1.22+** installed.

### Git Hooks (Lefthook)

This project uses [Lefthook](https://github.com/evilmartians/lefthook) to run pre-commit checks and validate commit messages automatically.

```bash
# Install lefthook (pick one):
brew install lefthook          # macOS
go install github.com/evilmartians/lefthook@latest  # Go
npm install -g lefthook        # npm

# Activate hooks in your local clone:
make install-hooks
# or: lefthook install
```

Once installed, every commit will automatically:
- **pre-commit**: check `gofmt`, run `go vet`, and run `golangci-lint` (if installed)
- **commit-msg**: validate the message follows [Conventional Commits](#commit-messages)

### Make Targets

```bash
make test       # Run tests
make coverage   # Run tests with coverage (90% threshold enforced)
make lint       # Run golangci-lint
make fmt        # Format code (gofmt -s -w)
make check      # Verify formatting (no write, CI-friendly)
make vet        # Run go vet
make all        # fmt + vet + lint + test
make install-hooks  # Install lefthook git hooks
```

## Code Style

- Run `gofmt` on all code (or `make fmt`)
- Run `golangci-lint` (or `make lint`)
- This project has **zero external dependencies** — keep it that way. Only use the Go standard library.
- Follow existing patterns in the codebase:
  - Functional options for configuration (`WithXxx`)
  - `context.Context` as the first parameter on all public methods
  - Custom error types with helper methods

## Running Tests

```bash
# All tests
make test

# Specific test
go test -run TestSearch -v ./...

# With coverage report
make coverage
```

Tests run in CI against Go 1.22, 1.23, and 1.24. A minimum of **90% code coverage** is enforced.

## Commit Messages

This project enforces [Conventional Commits](https://www.conventionalcommits.org/) via a git hook. Format:

```
<type>[optional scope]: <description>
```

Allowed types: `feat`, `fix`, `docs`, `test`, `chore`, `refactor`, `ci`, `style`, `perf`, `build`

Examples:

```
feat: add support for filtering by audio codec
fix(client): handle nil response body on 204
docs: update Quick Start example
test: add edge case tests for retry logic
chore: update CI matrix to Go 1.24
refactor: extract retry logic into helper
```

## Pull Request Guidelines

- Keep PRs focused — one feature or fix per PR
- Include tests for new functionality
- Update documentation if the public API changes
- Ensure all CI checks pass before requesting review
- Link related issues in the PR description

## Reporting Bugs

[Open an issue](https://github.com/torrentclaw/go-client/issues/new?labels=bug) with:

- **Description** — what went wrong
- **Steps to reproduce** — minimal code or commands to trigger the bug
- **Expected behavior** — what you expected to happen
- **Actual behavior** — what actually happened
- **Environment** — Go version, OS, client version

## Requesting Features

[Open an issue](https://github.com/torrentclaw/go-client/issues/new?labels=enhancement) with:

- **Problem** — what are you trying to solve?
- **Proposed solution** — how do you think it should work?
- **Alternatives considered** — other approaches you thought about

## Code of Conduct

This project follows the [Contributor Covenant v2.1](https://www.contributor-covenant.org/version/2/1/code_of_conduct/).

In short:

- **Be respectful** — treat everyone with dignity regardless of background or experience level
- **Be constructive** — focus on what's best for the project and community
- **Be collaborative** — welcome newcomers, help others learn
- **No harassment** — unacceptable behavior includes trolling, insults, and unwelcome attention

Violations can be reported to the project maintainers. All complaints will be reviewed and investigated promptly and fairly.

## Questions?

If you're unsure about something, [open a discussion](https://github.com/torrentclaw/go-client/issues) or reach out on Discord (coming soon).

---

Thank you for helping make TorrentClaw better!
