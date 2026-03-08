.PHONY: build install clean init test

# Build variables
BINARY_NAME=bin/specforge
GO_CMD=/Users/gsortino/.local/share/mise/installs/go/1.26.0/bin/go
GO_BUILD=$(GO_CMD) build
GO_INSTALL=$(GO_CMD) install
MAIN_PATH=cmd/specforge

# Default target
all: build

# Build the binary
build:
	mkdir -p bin
	$(GO_BUILD) -o $(BINARY_NAME) ./$(MAIN_PATH)

# Build for all platforms
build-all:
	GOOS=darwin GOARCH=amd64 $(GO_BUILD) -o $(BINARY_NAME)-darwin-amd64 ./$(MAIN_PATH)
	GOOS=darwin GOARCH=arm64 $(GO_BUILD) -o $(BINARY_NAME)-darwin-arm64 ./$(MAIN_PATH)
	GOOS=linux GOARCH=amd64 $(GO_BUILD) -o $(BINARY_NAME)-linux-amd64 ./$(MAIN_PATH)

# Install the binary to GOPATH/bin
install:
	$(GO_INSTALL) ./$(MAIN_PATH)

# Clean build artifacts
clean:
	rm -rf bin/
	rm -f $(HOME)/bin/$(BINARY_NAME)

# Initialize SpecForge (install Claude Code commands)
init: build
	./$(BINARY_NAME) init

# Run tests
test:
	$(GO_CMD) test -v ./...

# Run linter (if golangci-lint is installed)
lint:
	@which golangci-lint > /dev/null 2>&1 && golangci-lint run || echo "golangci-lint not installed"

# Development: watch and rebuild
dev:
	$(GO_CMD) build -o $(BINARY_NAME) -w ./$(MAIN_PATH)

# Show help
help:
	@echo "SpecForge Build Targets"
	@echo ""
	@echo "  make build      Build the specforge binary"
	@echo "  make install    Install specforge to \$$GOPATH/bin"
	@echo "  make init       Build and run specforge init"
	@echo "  make clean      Remove build artifacts"
	@echo "  make test       Run tests"
	@echo "  make lint       Run linter"
	@echo "  make build-all  Build for all platforms"
