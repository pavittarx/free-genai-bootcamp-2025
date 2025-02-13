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
CMD_DIR=.
PKG_DIR=./pkg
INTERNAL_DIR=./internal

# Linting configuration
LINT_CONFIG=.golangci.yml

.PHONY: all build test lint clean init run dev

# Default target
all: lint test build

# Initialize project dependencies
init:
	$(GOCMD) mod tidy
	$(GOLINT) run --install

# Build the application
build:
	mkdir -p bin
	$(GOBUILD) -o $(BINARY_PATH) .

# Run tests
test:
	$(GOTEST) ./... -v

# Run linter
lint:
	$(GOLINT) run --config=$(LINT_CONFIG) ./...

# Initialize database
db-init:
	./scripts/init_database.sh

# Run the application
run: build
	$(BINARY_PATH)

# Run with hot reloading
dev:
	$(HOME)/go/bin/air

# Clean build artifacts
clean:
	rm -rf bin
	$(GOCMD) clean -modcache
