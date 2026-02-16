# Go version should match go.mod
GO_VERSION := 1.24
DOCKER_IMAGE := golang:$(GO_VERSION)
BINARY := jk

.PHONY: test test-linux test-all vet fmt check build build-all clean help

# --- Day-to-day ---

test: ## Run tests (native, current platform)
	go test ./...

vet: ## Run static analysis
	go vet ./...

fmt: ## Check formatting (prints badly-formatted files)
	@bad=$$(gofmt -l .); if [ -n "$$bad" ]; then echo "gofmt: needs formatting:"; echo "$$bad"; exit 1; fi

build: ## Build for current platform
	go build -o $(BINARY) .

# --- Pre-commit (the real gate) ---

test-linux: ## Run tests inside a Linux container
	docker run --rm -v "$(CURDIR)":/app -w /app $(DOCKER_IMAGE) go test ./...

test-all: test test-linux ## Run tests on macOS (native) + Linux (Docker)

build-all: ## Cross-compile for all platforms (compile check only)
	@# darwin: needs CGO_ENABLED=1 for libproc (Apple clang handles both archs)
	CGO_ENABLED=1 GOOS=darwin GOARCH=amd64 go build -o /dev/null . && echo "  ok darwin/amd64"
	CGO_ENABLED=1 GOOS=darwin GOARCH=arm64 go build -o /dev/null . && echo "  ok darwin/arm64"
	@# linux: pure Go, no cgo needed
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /dev/null . && echo "  ok linux/amd64"
	CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -o /dev/null . && echo "  ok linux/arm64"

check: fmt vet test-all build-all ## Full pre-commit gate: fmt, vet, test (both platforms), cross-compile
	@echo ""
	@echo "All checks passed."

# --- Cleanup ---

clean: ## Remove build artifacts
	rm -f $(BINARY)

# --- Help ---

help: ## Show this help
	@grep -E '^[a-z][a-z_-]+:.*## ' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*## "}; {printf "  %-14s %s\n", $$1, $$2}'
