.DEFAULT_GOAL := help

## help: Show the help docs.
.PHONY: help
help:
	@echo "Options:\n"
	@sed -n 's|^##||p' ${PWD}/Makefile

## build: Build the CLI binaries with goreleaser.
.PHONY: build
build:
	@echo "Building with goreleaser..."
	@goreleaser release --snapshot --rm-dist

## format: Run format.
.PHONY: format
format:
	@echo "Running formatters..."
	@goimports -w --local github.com/dwin ./.

## lint: Run linter.
.PHONY: lint
lint:
	@echo "Running linter..."
	@golangci-lint run

## test: Run tests.
.PHONY: test
test:
	@echo "Running tests with coverage..."
	@go test ./... -race -coverprofile=coverage.out -covermode=atomic


	