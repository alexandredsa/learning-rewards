.PHONY: all build run stop clean dev rebuild

# Default target
all: build run

# Build all services using Docker Compose with no cache
build:
	@echo "Building Docker images (no cache)..."
	docker-compose build --no-cache

# Development mode - rebuild and run
dev:
	@echo "Building and running in development mode..."
	docker-compose up --build --force-recreate

# Rebuild a specific service
# Usage: make rebuild SERVICE=catalog-api
rebuild:
	@if [ -z "$(SERVICE)" ]; then \
		echo "Error: SERVICE argument is required"; \
		echo "Usage: make rebuild SERVICE=<service-name>"; \
		echo "Available services: catalog-api, event-processor, reward-processor-api, reward-processor-worker"; \
		exit 1; \
	fi
	@echo "Rebuilding $(SERVICE)..."
	docker-compose up -d --build --force-recreate $(SERVICE)

# Start all services using Docker Compose
run:
	@echo "Starting all services..."
	docker-compose up -d
	@echo "Services are running!"
	@echo "Catalog API: http://localhost:8080"
	@echo "Event Processor: http://localhost:8081/health"
	@echo "Reward Processor API: http://localhost:8082/health"
	@echo "Reward Processor Worker: http://localhost:8083/health"
	@echo "Kafka UI: http://localhost:9094"

# Stop all services
stop:
	@echo "Stopping all services..."
	docker-compose down

# Clean up Docker resources
clean: stop
	@echo "Cleaning up Docker resources..."
	docker-compose down -v
	@echo "Removing unused Docker images..."
	docker image prune -f

# Help command
help:
	@echo "Available commands:"
	@echo "  make              - Build and run all services"
	@echo "  make build        - Build all Docker images (no cache)"
	@echo "  make dev          - Build and run in development mode (forces rebuild)"
	@echo "  make run          - Start all services using Docker Compose"
	@echo "  make stop         - Stop all services"
	@echo "  make clean        - Clean up Docker resources and remove unused images"
	@echo "  make rebuild SERVICE=<name> - Rebuild and restart a specific service"
	@echo "                              Available services: catalog-api, event-processor, reward-processor-api, reward-processor-worker"
	@echo "  make help         - Show this help message"

.PHONY: stress-test install-vegeta clean-stress-test

stress-test:
	$(MAKE) -C stress-test stress-test

install-vegeta:
	$(MAKE) -C stress-test install-vegeta

clean-stress-test:
	$(MAKE) -C stress-test clean-stress-test
