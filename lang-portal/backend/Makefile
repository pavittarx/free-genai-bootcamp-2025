# Makefile for Language Portal Backend

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOTEST=$(GOCMD) test
GOLINT=golangci-lint

# Binary output
BINARY_NAME=lang-portal
BINARY_PATH=./bin/$(BINARY_NAME)

# Directories
CMD_DIR=./cmd/server
PKG_DIR=./pkg
INTERNAL_DIR=./internal

# Linting configuration
LINT_CONFIG=.golangci.yml

.PHONY: all build test lint clean init run

# Default target
all: lint test build

# Initialize project dependencies
init:
	$(GOCMD) mod tidy
	$(GOLINT) run --install

# Build the application
build:
	mkdir -p bin
	$(GOBUILD) -o $(BINARY_PATH) $(CMD_DIR)

# Run tests
test:
	$(GOTEST) ./... -v

# Run linter
lint:
	$(GOLINT) run --config=$(LINT_CONFIG) ./...

# Initialize database
db-init:
	./scripts/init_db.sh

# Run the application
run: build
	$(BINARY_PATH)

# Clean build artifacts
clean:
	rm -rf bin
	$(GOCMD) clean -modcache

# Development workflow
dev: lint test run
