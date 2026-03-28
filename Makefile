.PHONY: all test lint coverage clean fmt vet check install-hooks

all: fmt vet lint test

## Run all tests
test:
	go test -v -race -count=1 ./...

## Run linter (requires golangci-lint)
lint:
	golangci-lint run ./...

## Run tests with coverage report
coverage:
	go test -race -coverprofile=coverage.out -covermode=atomic ./...
	go tool cover -func=coverage.out
	go tool cover -html=coverage.out -o coverage.html

## Format code
fmt:
	gofmt -s -w .

## Check formatting (no write, exits non-zero if unformatted)
check:
	@test -z "$$(gofmt -l .)" || { echo "Files not formatted:"; gofmt -l .; exit 1; }

## Run go vet
vet:
	go vet ./...

## Install lefthook git hooks
install-hooks:
	lefthook install

## Remove generated files
clean:
	rm -f coverage.out coverage.html

# Release with goreleaser (future):
#   1. Install goreleaser: https://goreleaser.com/install/
#   2. Create .goreleaser.yml config
#   3. Tag a version: git tag -a v0.x.0 -m "release v0.x.0"
#   4. Run: goreleaser release --clean
