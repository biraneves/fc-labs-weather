# Variables
APP_NAME := weather-app
BINARY_NAME := svc
DOCKER_COMPOSE := docker compose

export IMAGE_TAG := $(shell git describe --tags --abbrev=0 2>/dev/null || echo "latest")

# Help
.DEFAULT_GOAL := help

.PHONY: help run test test-coverage docker-build docker-up docker-down docker-logs

help:
	@echo "Usage: make <target>"
	@echo ""
	@echo "Current image tag: $(IMAGE_TAG)"
	@echo ""
	@echo "Targets:"
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-20s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

run: ## Runs the application locally.
	@echo "--> Running application locally..."
	@go run ./cmd/server/

test: ## Runs all unit tests.
	@echo "--> Running unit tests..."
	@go test -v -race ./...

test-coverage: ## Runs tests and shows the coverage report in browser.
	@echo "--> Running tests and generating coverage report..."
	@go test -coverprofile=coverage.out ./...
	@echo "--> Opening coverage report in browser..."
	@go tool cover -html=coverage.out

docker-build: ## Builds production Docker image.
	@echo "--> Building production Docker image..."
	@$(DOCKER_COMPOSE) build

docker-up: ## Starts containers in dettached mode.
	@echo "--> Starting services with Docker Compose..."
	@$(DOCKER_COMPOSE) up -d --build

docker-down: ## Stops and removes the containers.
	@echo "--> Stopping services..."
	@$(DOCKER_COMPOSE) down

docker-logs: ## Shows application logs in real time.
	@echo "--> Tailing application logs..."
	@$(DOCKER_COMPOSE) logs -f $(APP_NAME)