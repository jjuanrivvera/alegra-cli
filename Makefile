# alegra-cli Makefile

BINARY_NAME := alegra
PKG         := github.com/jjuanrivvera/alegra-cli
MAIN        := ./cmd/alegra
BIN_DIR     := bin
INSTALL_DIR ?= /usr/local/bin

VERSION    ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo dev)
COMMIT     := $(shell git rev-parse --short HEAD 2>/dev/null || echo none)
BUILD_DATE := $(shell date -u +%Y-%m-%dT%H:%M:%SZ)

LDFLAGS := -ldflags "\
	-s -w \
	-X $(PKG)/internal/version.Version=$(VERSION) \
	-X $(PKG)/internal/version.Commit=$(COMMIT) \
	-X $(PKG)/internal/version.BuildDate=$(BUILD_DATE)"

GO ?= go

.PHONY: help
help: ## Show this help
	@grep -E '^[a-zA-Z0-9_-]+:.*?## .*$$' $(MAKEFILE_LIST) | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-18s\033[0m %s\n", $$1, $$2}'

.PHONY: build
build: ## Build the binary into bin/
	@mkdir -p $(BIN_DIR)
	$(GO) build $(LDFLAGS) -o $(BIN_DIR)/$(BINARY_NAME) $(MAIN)

.PHONY: install
install: build ## Install the binary to INSTALL_DIR
	install -m 0755 $(BIN_DIR)/$(BINARY_NAME) $(INSTALL_DIR)/$(BINARY_NAME)

.PHONY: uninstall
uninstall: ## Remove the installed binary
	rm -f $(INSTALL_DIR)/$(BINARY_NAME)

.PHONY: run
run: build ## Build and run (use ARGS="...")
	@$(BIN_DIR)/$(BINARY_NAME) $(ARGS)

.PHONY: test
test: ## Run all tests
	$(GO) test ./...

.PHONY: test-race
test-race: ## Run tests with the race detector
	$(GO) test -race ./...

.PHONY: test-coverage
test-coverage: ## Run tests with a coverage report (coverage.html)
	$(GO) test -coverprofile=coverage.out ./...
	$(GO) tool cover -html=coverage.out -o coverage.html
	@echo "coverage report: coverage.html"

.PHONY: fmt
fmt: ## Format the code
	$(GO) fmt ./...
	gofmt -s -w .

.PHONY: vet
vet: ## Run go vet
	$(GO) vet ./...

.PHONY: lint
lint: ## Run golangci-lint
	golangci-lint run

.PHONY: tidy
tidy: ## Tidy go.mod/go.sum
	$(GO) mod tidy

.PHONY: dev
dev: fmt vet build ## Format, vet, and build

.PHONY: check
check: fmt vet lint test ## Run the full local quality gate

.PHONY: docs-gen
docs-gen: ## Generate CLI reference docs from the command tree
	$(GO) run ./tools/gendocs

.PHONY: docs-serve
docs-serve: docs-gen ## Serve the docs site locally
	mkdocs serve

.PHONY: docs-build
docs-build: docs-gen ## Build the docs site
	mkdocs build --strict

.PHONY: snapshot
snapshot: ## Build a local release snapshot with goreleaser
	goreleaser release --snapshot --clean

.PHONY: setup-hooks
setup-hooks: ## Install git pre-commit hooks
	git config core.hooksPath .githooks
	@echo "git hooks installed (.githooks)"

.PHONY: clean
clean: ## Remove build artifacts
	rm -rf $(BIN_DIR) dist coverage.out coverage.html site
