# Makefile for the Go Vector Search service

BINARY_NAME=go-vector-search

.PHONY: help build run lint test-unit test-integration docker-build

.DEFAULT_GOAL := help

help:
	@echo "Usage: make [command]"
	@echo ""
	@echo "Available commands:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-20s\033[0m %s\n", $$1, $$2}'

build: ## Build the application binary
	@echo "Building the application..."
	@go build -o ./bin/$(BINARY_NAME) ./cmd/server

run: ## Run the application
	@echo "Running the application..."
	@go run ./cmd/server/main.go

lint: ## Lint the codebase
	@echo "Linting the codebase..."
	@golangci-lint run

test-unit: ## Run unit tests
	@echo "Running unit tests..."
	@docker-compose -f test/docker-compose.yml down -v --remove-orphans
	@docker-compose -f test/docker-compose.yml up -d
	@go test -v ./... || (docker-compose -f test/docker-compose.yml logs && exit 1)
	@docker-compose -f test/docker-compose.yml down

test-integration: ## Run integration tests
	@echo "Running integration tests..."
	@docker-compose -f test/docker-compose.yml down -v --remove-orphans
	@docker-compose -f test/docker-compose.yml up -d
	@go test -v ./test/... || (docker-compose -f test/docker-compose.yml logs && exit 1)
	@docker-compose -f test/docker-compose.yml down

docker-build: ## Build the Docker image
	@echo "Building the Docker image..."
	@docker build -t $(BINARY_NAME):latest .
