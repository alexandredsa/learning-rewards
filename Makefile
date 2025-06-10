.PHONY: all build run stop clean

# Default target
all: build run

# Build all services using Docker Compose
build:
	@echo "Building Docker images..."
	docker-compose build

# Start all services using Docker Compose
run:
	@echo "Starting all services..."
	docker-compose up -d
	@echo "Services are running!"
	@echo "Catalog API: http://localhost:8080"
	@echo "Event Processor: http://localhost:8081/health"
	@echo "Reward Processor: http://localhost:8082/health"
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
	@echo "  make        - Build and run all services"
	@echo "  make build  - Build all Docker images"
	@echo "  make run    - Start all services using Docker Compose"
	@echo "  make stop   - Stop all services"
	@echo "  make clean  - Clean up Docker resources and remove unused images"
	@echo "  make help   - Show this help message"

.PHONY: stress-test install-vegeta clean-stress-test

stress-test:
	$(MAKE) -C stress-test stress-test

install-vegeta:
	$(MAKE) -C stress-test install-vegeta

clean-stress-test:
	$(MAKE) -C stress-test clean-stress-test
