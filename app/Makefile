# Makefile for stormsensor

GO = go
GOFLAGS = -v

VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
GIT_COMMIT ?= $(shell git rev-parse HEAD 2>/dev/null || echo "unknown")
BUILD_TIME ?= $(shell date -u '+%Y-%m-%d %H:%M:%S UTC')
LDFLAGS = -ldflags "-s -w -X 'github.com/philipparndt/stormsensor/version.Version=$(VERSION)' -X 'github.com/philipparndt/stormsensor/version.GitCommit=$(GIT_COMMIT)' -X 'github.com/philipparndt/stormsensor/version.BuildTime=$(BUILD_TIME)'"

BUILD_DIR = build
CONFIG_DIR = $(PWD)/../production/config

BINARY_NAME = stormsensor
DOCKER_IMAGE_NAME = pharndt/stormmqtt:latest

.PHONY: build
build:
	@echo "Building $(BINARY_NAME)..."
	@$(GO) build $(GOFLAGS) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) .

.PHONY: version
version:
	@echo "Version: $(VERSION)"
	@echo "Git Commit: $(GIT_COMMIT)"
	@echo "Build Time: $(BUILD_TIME)"

.PHONY: run
run: build
	@echo "Running $(BINARY_NAME)..."
	@$(BUILD_DIR)/$(BINARY_NAME) $(CONFIG_DIR)/config.json

.PHONY: test
test:
	@echo "Running tests..."
	@$(GO) test $(GOFLAGS) ./...

.PHONY: docker
docker:
	@echo "Building Docker image..."
	@docker build -t $(DOCKER_IMAGE_NAME) .

.PHONY: clean
clean:
	@echo "Cleaning up..."
	@rm -rf $(BUILD_DIR)

.PHONY: deps
deps:
	@echo "Downloading Go dependencies..."
	@$(GO) mod download
	@$(GO) mod tidy
