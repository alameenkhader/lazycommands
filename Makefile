.PHONY: build install test clean run

build:
	@echo "Building lazycommands..."
	@go build -o bin/lazycommands main.go
	@echo "Build complete: bin/lazycommands"

install:
	@echo "Installing lazycommands..."
	@go install
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
	@echo "Building for multiple platforms..."
	@mkdir -p bin
	GOOS=linux GOARCH=amd64 go build -o bin/lazycommands-linux-amd64 main.go
	GOOS=darwin GOARCH=amd64 go build -o bin/lazycommands-darwin-amd64 main.go
	GOOS=darwin GOARCH=arm64 go build -o bin/lazycommands-darwin-arm64 main.go
	GOOS=windows GOARCH=amd64 go build -o bin/lazycommands-windows-amd64.exe main.go
	@echo "Cross-compilation complete"
