.PHONY: build install test clean run

# Version variables
VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
BUILD_DATE := $(shell date -u '+%Y-%m-%d_%H:%M:%S')
GIT_COMMIT := $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")

LDFLAGS := -X github.com/alameen/lazycommands/internal/version.Version=$(VERSION) \
           -X github.com/alameen/lazycommands/internal/version.BuildDate=$(BUILD_DATE) \
           -X github.com/alameen/lazycommands/internal/version.GitCommit=$(GIT_COMMIT)

build:
	@echo "Building lazycommands $(VERSION)..."
	@mkdir -p bin
	@go build -ldflags "$(LDFLAGS)" -o bin/lazycommands main.go
	@echo "Build complete: bin/lazycommands"

install:
	@echo "Installing lazycommands $(VERSION)..."
	@go install -ldflags "$(LDFLAGS)"
	@echo "Installed to $(shell go env GOPATH)/bin/lazycommands"

test:
	@echo "Running tests..."
	@go test ./...

run:
	@go run main.go $(ARGS)

clean:
	@echo "Cleaning build artifacts..."
	@rm -rf bin/
	@echo "Clean complete"

# Cross-compilation examples
build-all:
	@echo "Building for multiple platforms (version $(VERSION))..."
	@mkdir -p bin
	GOOS=linux GOARCH=amd64 go build -ldflags "$(LDFLAGS)" -o bin/lazycommands-linux-amd64 main.go
	GOOS=darwin GOARCH=amd64 go build -ldflags "$(LDFLAGS)" -o bin/lazycommands-darwin-amd64 main.go
	GOOS=darwin GOARCH=arm64 go build -ldflags "$(LDFLAGS)" -o bin/lazycommands-darwin-arm64 main.go
	GOOS=windows GOARCH=amd64 go build -ldflags "$(LDFLAGS)" -o bin/lazycommands-windows-amd64.exe main.go
	@echo "Cross-compilation complete"
