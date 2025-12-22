.PHONY: build test clean docker-build docker-push run-api run-worker lint fmt

# Variables
API_BINARY=bin/api
WORKER_BINARY=bin/worker
DOCKER_REGISTRY=your-registry
API_IMAGE=dab-api
WORKER_IMAGE=dab-worker
VERSION?=latest

# Build
build: build-api build-worker

build-api:
	@echo "Building API..."
	@go build -o $(API_BINARY) cmd/api/main.go

build-worker:
	@echo "Building Worker..."
	@go build -o $(WORKER_BINARY) cmd/worker/main.go

# Run
run-api:
	@go run cmd/api/main.go

run-worker:
	@go run cmd/worker/main.go

# Test
test:
	@echo "Running tests..."
	@go test -v -cover ./...

test-coverage:
	@go test -v -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html

# Docker
docker-build: docker-build-api docker-build-worker

docker-build-api:
	@docker build -f deployments/docker/api.Dockerfile -t $(API_IMAGE):$(VERSION) .

docker-build-worker:
	@docker build -f deployments/docker/worker.Dockerfile -t $(WORKER_IMAGE):$(VERSION) .

docker-push:
	@docker tag $(API_IMAGE):$(VERSION) $(DOCKER_REGISTRY)/$(API_IMAGE):$(VERSION)
	@docker tag $(WORKER_IMAGE):$(VERSION) $(DOCKER_REGISTRY)/$(WORKER_IMAGE):$(VERSION)
	@docker push $(DOCKER_REGISTRY)/$(API_IMAGE):$(VERSION)
	@docker push $(DOCKER_REGISTRY)/$(WORKER_IMAGE):$(VERSION)

# Development
dev:
	@docker-compose up -d

down:
	@docker-compose down

logs:
	@docker-compose logs -f

# Code Quality
lint:
	@golangci-lint run

fmt:
	@go fmt ./...

vet:
	@go vet ./...

# Clean
clean:
	@rm -rf bin/
	@rm -f coverage.out coverage.html

# Dependencies
deps:
	@go mod download
	@go mod tidy

# Help
help:
	@echo "Available targets:"
	@echo "  build          - Build both services"
	@echo "  test           - Run tests"
	@echo "  docker-build   - Build Docker images"
	@echo "  dev            - Start development environment"
	@echo "  lint           - Run linter"
	@echo "  clean          - Clean build artifacts"