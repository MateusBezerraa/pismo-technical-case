.PHONY: help run test test-cover build docker docker-run clean fmt vet tidy

# Variables
BINARY := bin/api
DOCKER_IMAGE := pismo-technical-case
PORT := 8080

help:
	@echo "Available commands:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2}'

run: ## Execute the API locally
	go run ./cmd/api

test: ## Run all tests with race detector
	go test ./... -race -count=1

test-cover: ## Runs tests + prints per-package coverage
	go test ./... -cover -count=1

test-cover-total: ## Runs tests + prints total coverage (includes main.go)
	go test ./... -coverprofile=coverage.out -coverpkg=./...
	go tool cover -func=coverage.out | tail -1

test-html: ## Opens coverage report in browser
	go test ./... -coverprofile=coverage.out -coverpkg=./...
	go tool cover -html=coverage.out

build: ## Compile the binary
	go build -trimpath -ldflags="-s -w" -o $(BINARY) ./cmd/api

docker: ## Build the Docker image
	docker build -t $(DOCKER_IMAGE) .

docker-run: ## Run the container (mapping port 8080)
	docker run --rm -p $(PORT):$(PORT) -v pismo-data:/data $(DOCKER_IMAGE)

fmt: ## Format the code
	gofmt -s -w .

vet: ## Run go vet
	go vet ./...

tidy: ## Organize dependencies
	go mod tidy

clean: ## Remove build artifacts
	rm -rf bin/ coverage.out coverage.html pismo.db

docker-clean: ## Removes the Docker volume (wipes data)
	docker volume rm pismo-data

