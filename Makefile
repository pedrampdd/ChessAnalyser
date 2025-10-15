.PHONY: build run test clean docker-build docker-run help

# Variables
BINARY_NAME=chess-analyzer
DOCKER_IMAGE=chess-analyzer
DOCKER_TAG=latest

# Default target
help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

build: ## Build the application
	@echo "Building $(BINARY_NAME)..."
	go build -o $(BINARY_NAME) .

run: ## Run the application
	@echo "Running $(BINARY_NAME)..."
	go run .

test: ## Run tests
	@echo "Running tests..."
	go test -v ./...

test-coverage: ## Run tests with coverage
	@echo "Running tests with coverage..."
	go test -v -cover ./...

clean: ## Clean build artifacts
	@echo "Cleaning..."
	go clean
	rm -f $(BINARY_NAME)

docker-build: ## Build Docker image
	@echo "Building Docker image..."
	docker build -t $(DOCKER_IMAGE):$(DOCKER_TAG) .

docker-run: ## Run Docker container
	@echo "Running Docker container..."
	docker run -p 8080:8080 $(DOCKER_IMAGE):$(DOCKER_TAG)

docker-compose-up: ## Start with docker-compose
	@echo "Starting with docker-compose..."
	docker-compose up --build

docker-compose-down: ## Stop docker-compose
	@echo "Stopping docker-compose..."
	docker-compose down

lint: ## Run linter
	@echo "Running linter..."
	golangci-lint run

fmt: ## Format code
	@echo "Formatting code..."
	go fmt ./...

vet: ## Run go vet
	@echo "Running go vet..."
	go vet ./...

mod-tidy: ## Tidy go modules
	@echo "Tidying go modules..."
	go mod tidy

mod-download: ## Download go modules
	@echo "Downloading go modules..."
	go mod download

dev: ## Run in development mode with hot reload
	@echo "Running in development mode..."
	air

install-tools: ## Install development tools
	@echo "Installing development tools..."
	go install github.com/cosmtrek/air@latest
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

check: fmt vet lint test ## Run all checks
