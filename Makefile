# Makefile for the Go Vector Search service

BINARY_NAME=go-vector-search

.PHONY: build run lint test-unit test-integration docker-build

build:
	@echo "Building the application..."
	@go build -o ./bin/$(BINARY_NAME) ./cmd/server

run:
	@echo "Running the application..."
	@go run ./cmd/server/main.go

lint:
	@echo "Linting the codebase..."
	@golangci-lint run

test-unit:
	@echo "Running unit tests..."
	@go test -v ./...

test-integration:
	@echo "Running integration tests..."
	@docker-compose -f test/docker-compose.yml down -v --remove-orphans
	@docker-compose -f test/docker-compose.yml up -d
	@go test -v ./test/... || (docker-compose -f test/docker-compose.yml logs && exit 1)
	@docker-compose -f test/docker-compose.yml down

docker-build:
	@echo "Building the Docker image..."
	@docker build -t $(BINARY_NAME):latest .
