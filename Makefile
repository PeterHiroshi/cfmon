# cfmon - Cloudflare Workers/Containers CLI Monitor
# Makefile for building, testing, and managing the project

BINARY_NAME := cfmon
GO := /usr/local/go/bin/go
VERSION := $(shell git describe --tags --abbrev=0 2>/dev/null || echo "v0.1.0")
BUILD_TIME := $(shell date -u '+%Y-%m-%d_%H:%M:%S')
GIT_COMMIT := $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
LDFLAGS := -ldflags "-X main.Version=$(VERSION) -X main.BuildTime=$(BUILD_TIME) -X main.GitCommit=$(GIT_COMMIT)"

# Go packages
PACKAGES := $(shell $(GO) list ./...)
INTERNAL_PACKAGES := $(shell $(GO) list ./internal/...)

# Colors for output
RESET := \033[0m
GREEN := \033[32m
YELLOW := \033[33m
BLUE := \033[34m
RED := \033[31m

.PHONY: all build test lint install clean help coverage integration-test fmt vet

# Default target
all: fmt vet lint test build ## Format, vet, lint, test, and build

build: ## Build the binary
	@echo "$(BLUE)Building $(BINARY_NAME)...$(RESET)"
	@$(GO) build $(LDFLAGS) -o $(BINARY_NAME) .
	@echo "$(GREEN)✓ Build complete: ./$(BINARY_NAME)$(RESET)"

test: ## Run unit tests
	@echo "$(BLUE)Running tests...$(RESET)"
	@$(GO) test -v -count=1 -race -cover $(PACKAGES)
	@echo "$(GREEN)✓ Tests passed$(RESET)"

coverage: ## Run tests with coverage report
	@echo "$(BLUE)Running tests with coverage...$(RESET)"
	@$(GO) test -v -count=1 -race -coverprofile=coverage.out $(PACKAGES)
	@$(GO) tool cover -html=coverage.out -o coverage.html
	@echo "$(GREEN)✓ Coverage report generated: coverage.html$(RESET)"
	@echo "$(YELLOW)Coverage summary:$(RESET)"
	@$(GO) tool cover -func=coverage.out | grep total | awk '{print "Total coverage: " $$3}'

integration-test: build ## Run integration tests
	@echo "$(BLUE)Running integration tests...$(RESET)"
	@$(GO) test -v -count=1 ./test/integration/... -tags=integration
	@echo "$(GREEN)✓ Integration tests passed$(RESET)"

lint: ## Run golangci-lint
	@echo "$(BLUE)Running linters...$(RESET)"
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run ./...; \
		echo "$(GREEN)✓ Linting complete$(RESET)"; \
	else \
		echo "$(YELLOW)⚠ golangci-lint not installed. Install with:$(RESET)"; \
		echo "  curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $$(go env GOPATH)/bin"; \
		echo "$(YELLOW)Running basic go vet instead...$(RESET)"; \
		$(GO) vet $(PACKAGES); \
	fi

fmt: ## Format code with gofmt
	@echo "$(BLUE)Formatting code...$(RESET)"
	@$(GO) fmt $(PACKAGES)
	@echo "$(GREEN)✓ Code formatted$(RESET)"

vet: ## Run go vet
	@echo "$(BLUE)Running go vet...$(RESET)"
	@$(GO) vet $(PACKAGES)
	@echo "$(GREEN)✓ Vet complete$(RESET)"

install: build ## Install the binary to GOPATH/bin
	@echo "$(BLUE)Installing $(BINARY_NAME) to $(shell $(GO) env GOPATH)/bin...$(RESET)"
	@$(GO) install $(LDFLAGS)
	@echo "$(GREEN)✓ Installed to $(shell $(GO) env GOPATH)/bin/$(BINARY_NAME)$(RESET)"

clean: ## Clean build artifacts
	@echo "$(BLUE)Cleaning build artifacts...$(RESET)"
	@rm -f $(BINARY_NAME)
	@rm -f coverage.out coverage.html
	@rm -rf dist/
	@echo "$(GREEN)✓ Clean complete$(RESET)"

deps: ## Download and verify dependencies
	@echo "$(BLUE)Downloading dependencies...$(RESET)"
	@$(GO) mod download
	@$(GO) mod verify
	@echo "$(GREEN)✓ Dependencies ready$(RESET)"

tidy: ## Run go mod tidy
	@echo "$(BLUE)Tidying modules...$(RESET)"
	@$(GO) mod tidy
	@echo "$(GREEN)✓ Modules tidied$(RESET)"

release: ## Build release binaries with goreleaser (requires goreleaser)
	@echo "$(BLUE)Building release binaries...$(RESET)"
	@if command -v goreleaser >/dev/null 2>&1; then \
		goreleaser release --snapshot --skip-publish --clean; \
		echo "$(GREEN)✓ Release binaries built in ./dist/$(RESET)"; \
	else \
		echo "$(RED)✗ goreleaser not installed$(RESET)"; \
		echo "Install with: go install github.com/goreleaser/goreleaser@latest"; \
		exit 1; \
	fi

run: build ## Build and run the binary
	@echo "$(BLUE)Running $(BINARY_NAME)...$(RESET)"
	@./$(BINARY_NAME)

dev: ## Run with hot reload (requires air)
	@if command -v air >/dev/null 2>&1; then \
		air; \
	else \
		echo "$(YELLOW)⚠ air not installed. Install with: go install github.com/cosmtrek/air@latest$(RESET)"; \
		echo "$(YELLOW)Running without hot reload...$(RESET)"; \
		$(GO) run .; \
	fi

check: fmt vet lint test ## Run all checks (fmt, vet, lint, test)
	@echo "$(GREEN)✓ All checks passed$(RESET)"

help: ## Display this help message
	@echo "$(BLUE)cfmon Makefile$(RESET)"
	@echo "$(YELLOW)Usage:$(RESET)"
	@echo "  make [target]"
	@echo ""
	@echo "$(YELLOW)Available targets:$(RESET)"
	@awk 'BEGIN {FS = ":.*##"; printf ""} /^[a-zA-Z_-]+:.*?##/ { printf "  $(GREEN)%-18s$(RESET) %s\n", $$1, $$2 } /^##@/ { printf "\n$(YELLOW)%s$(RESET)\n", substr($$0, 5) } ' $(MAKEFILE_LIST)
	@echo ""
	@echo "$(YELLOW)Examples:$(RESET)"
	@echo "  make build          # Build the binary"
	@echo "  make test           # Run tests"
	@echo "  make install        # Install to GOPATH/bin"
	@echo "  make all            # Run all checks and build"

.DEFAULT_GOAL := help