.PHONY: help build build-web build-pushgateway build-all clean test install deps fmt lint docker docker-push

# Variables
BINARY_NAME=starlink_exporter
VERSION?=$(shell git describe --tags --always --dirty)
BUILD_TIME=$(shell date -u '+%Y-%m-%d_%H:%M:%S')
GIT_COMMIT=$(shell git rev-parse HEAD)
GO_VERSION=$(shell go version | cut -d' ' -f3)

LDFLAGS=-ldflags "-X main.Version=$(VERSION) -X main.BuildTime=$(BUILD_TIME) -X main.GitCommit=$(GIT_COMMIT) -X main.GoVersion=$(GO_VERSION)"

DOCKER_REGISTRY?=registry.example.com
DOCKER_IMAGE=$(DOCKER_REGISTRY)/starlink-exporter
DOCKER_TAG?=$(VERSION)

# Output directories
OUTPUT_DIR=bin
DIST_DIR=dist

# Supported platforms for cross-compilation
PLATFORMS=linux/amd64 linux/arm64 linux/arm/v7 darwin/amd64 darwin/arm64

help: ## Display this help screen
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

deps: ## Install dependencies
	@echo "Installing dependencies..."
	go mod download
	go mod tidy

fmt: ## Format code
	@echo "Formatting code..."
	gofmt -s -w .
	goimports -w .

lint: ## Run linter
	@echo "Running linter..."
	golangci-lint run ./...

test: ## Run tests
	@echo "Running tests..."
	go test -v -race -coverprofile=coverage.out ./...

clean: ## Clean build artifacts
	@echo "Cleaning build artifacts..."
	rm -rf $(OUTPUT_DIR) $(DIST_DIR) coverage.out

build: ## Build for current platform
	@echo "Building $(BINARY_NAME) for current platform..."
	mkdir -p $(OUTPUT_DIR)
	CGO_ENABLED=0 go build $(LDFLAGS) -o $(OUTPUT_DIR)/$(BINARY_NAME) ./cmd/starlink_exporter

build-web: ## Build web mode binary
	@echo "Building $(BINARY_NAME) web mode..."
	mkdir -p $(OUTPUT_DIR)
	CGO_ENABLED=0 go build $(LDFLAGS) -o $(OUTPUT_DIR)/$(BINARY_NAME)-web ./cmd/starlink_exporter

build-pushgateway: ## Build pushgateway mode binary
	@echo "Building $(BINARY_NAME) pushgateway mode..."
	mkdir -p $(OUTPUT_DIR)
	CGO_ENABLED=0 go build $(LDFLAGS) -o $(OUTPUT_DIR)/$(BINARY_NAME)-pushgateway ./cmd/starlink_exporter

# Build for all supported architectures
build-all: clean ## Build for all platforms
	@echo "Building for all platforms..."
	@mkdir -p $(DIST_DIR)
	@for platform in $(PLATFORMS); do \
		OS=$$(echo $$platform | cut -d/ -f1); \
		ARCH=$$(echo $$platform | cut -d/ -f2); \
		ARM=$$(echo $$platform | cut -d/ -f3); \
		if [ "$$ARM" != "$$ARCH" ]; then \
			ARCH="$$ARCH$$ARM"; \
		fi; \
		BINARY=$(DIST_DIR)/$(BINARY_NAME)-$$OS-$$ARCH; \
		if [ "$$OS" = "windows" ]; then BINARY="$$BINARY.exe"; fi; \
		echo "Building for $$OS/$$ARCH -> $$BINARY"; \
		CGO_ENABLED=0 GOOS=$$OS GOARCH=$$ARCH go build $(LDFLAGS) -o $$BINARY ./cmd/starlink_exporter; \
	done
	@echo "Build complete! Binaries in $(DIST_DIR)/"

build-docker: ## Build Docker image
	@echo "Building Docker image: $(DOCKER_IMAGE):$(DOCKER_TAG)"
	docker build -t $(DOCKER_IMAGE):$(DOCKER_TAG) \
		-t $(DOCKER_IMAGE):latest \
		--build-arg VERSION=$(VERSION) \
		--build-arg BUILD_TIME=$(BUILD_TIME) \
		--build-arg GIT_COMMIT=$(GIT_COMMIT) \
		.

docker-push: build-docker ## Push Docker image to registry
	@echo "Pushing Docker image to $(DOCKER_IMAGE):$(DOCKER_TAG)"
	docker push $(DOCKER_IMAGE):$(DOCKER_TAG)
	docker push $(DOCKER_IMAGE):latest

install: build ## Install binary to $GOBIN
	@echo "Installing $(BINARY_NAME) to $$GOBIN..."
	cp $(OUTPUT_DIR)/$(BINARY_NAME) $$(go env GOBIN)/$(BINARY_NAME)

run-web: build ## Run in web mode (requires dish at 192.168.100.1:9200)
	@echo "Running in web mode..."
	$(OUTPUT_DIR)/$(BINARY_NAME) -mode=web -source=live -listen=:9817

run-web-dummy: build ## Run in web mode with dummy metrics
	@echo "Running in web mode with dummy metrics..."
	$(OUTPUT_DIR)/$(BINARY_NAME) -mode=web -source=dummy -listen=:9817

run-pushgateway: build ## Run in pushgateway mode (requires pushgateway at localhost:9091)
	@echo "Running in pushgateway mode..."
	$(OUTPUT_DIR)/$(BINARY_NAME) -mode=pushgateway -source=live -pushgateway=http://localhost:9091 -interval=15s

run-pushgateway-dummy: build ## Run in pushgateway mode with dummy metrics
	@echo "Running in pushgateway mode with dummy metrics..."
	$(OUTPUT_DIR)/$(BINARY_NAME) -mode=pushgateway -source=dummy -pushgateway=http://localhost:9091 -interval=15s

# ARM compilation targets (Raspberry Pi)
build-arm64: ## Build for ARM64 (Raspberry Pi 4)
	@echo "Building for ARM64..."
	mkdir -p $(OUTPUT_DIR)
	CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build $(LDFLAGS) -o $(OUTPUT_DIR)/$(BINARY_NAME)-arm64 ./cmd/starlink_exporter

build-armv7: ## Build for ARMv7 (Raspberry Pi 3, Zero)
	@echo "Building for ARMv7..."
	mkdir -p $(OUTPUT_DIR)
	CGO_ENABLED=0 GOOS=linux GOARCH=arm GOARM=7 go build $(LDFLAGS) -o $(OUTPUT_DIR)/$(BINARY_NAME)-armv7 ./cmd/starlink_exporter

build-amd64: ## Build for AMD64
	@echo "Building for AMD64..."
	mkdir -p $(OUTPUT_DIR)
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o $(OUTPUT_DIR)/$(BINARY_NAME)-amd64 ./cmd/starlink_exporter

# Release target
release: build-all ## Create release artifacts
	@echo "Creating release artifacts..."
	mkdir -p $(DIST_DIR)
	cd $(DIST_DIR) && for f in *; do \
		if [ -f "$$f" ]; then \
			sha256sum "$$f" > "$$f.sha256"; \
		fi; \
	done
	@echo "Release artifacts ready in $(DIST_DIR)/"

.PHONY: version
version: ## Show version info
	@echo "Version: $(VERSION)"
	@echo "Build Time: $(BUILD_TIME)"
	@echo "Git Commit: $(GIT_COMMIT)"
	@echo "Go Version: $(GO_VERSION)"
