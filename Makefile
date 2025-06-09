.PHONY: all build run stop clean

# Default target
all: build run

# Build both applications
build:
	@echo "Building applications..."
	cd catalog-api && go build -o bin/catalog-api
	cd event-processor && go build -o bin/event-processor

# Start databases and run both applications
run:
	@echo "Starting databases..."
	docker-compose up -d
	@echo "Starting applications..."
	@echo "Starting Catalog API..."
	cd catalog-api && DATABASE_DSN="postgres://user:pass@localhost:5432/catalog?sslmode=disable" ENV=dev PORT=8080 ./bin/catalog-api & echo $$! > .pid.catalog-api
	@echo "Starting Event Processor..."
	cd event-processor && DATABASE_DSN="postgres://user:pass@localhost:5433/event_processor?sslmode=disable" ENV=dev PORT=8081 KAFKA_BROKERS="localhost:29092" KAFKA_TOPIC="platform-events" ./bin/event-processor & echo $$! > .pid.event-processor
	@echo "Applications are running!"
	@echo "Catalog API: http://localhost:8080"
	@echo "Event Processor: http://localhost:8081/health"

# Stop both applications and databases
stop:
	@echo "Stopping applications..."
	@for pid_file in catalog-api/.pid.catalog-api event-processor/.pid.event-processor; do \
		if [ -f "$$pid_file" ]; then \
			pid=$$(cat "$$pid_file"); \
			if ps -p $$pid > /dev/null; then \
				echo "Stopping process $$pid..."; \
				kill -15 $$pid 2>/dev/null || true; \
				sleep 1; \
				if ps -p $$pid > /dev/null; then \
					echo "Process $$pid still running, force killing..."; \
					kill -9 $$pid 2>/dev/null || true; \
				fi; \
			fi; \
			rm -f "$$pid_file"; \
		fi; \
	done
	@echo "Stopping databases..."
	docker-compose down

# Clean up build artifacts and ensure processes are stopped
clean: stop
	@echo "Cleaning up..."
	rm -rf catalog-api/bin event-processor/bin
	@# Double check if any processes are still running on our ports
	@for port in 8080 8081; do \
		pid=$$(lsof -ti:$$port 2>/dev/null); \
		if [ ! -z "$$pid" ]; then \
			echo "Found process $$pid still running on port $$port, killing..."; \
			kill -9 $$pid 2>/dev/null || true; \
		fi; \
	done

# Help command
help:
	@echo "Available commands:"
	@echo "  make        - Build and run both applications"
	@echo "  make build  - Build both applications"
	@echo "  make run    - Start databases and run both applications"
	@echo "  make stop   - Stop both applications and databases"
	@echo "  make clean  - Clean up build artifacts and ensure all processes are stopped"
	@echo "  make help   - Show this help message"

.PHONY: stress-test install-vegeta clean-stress-test

stress-test:
	$(MAKE) -C stress-test stress-test

install-vegeta:
	$(MAKE) -C stress-test install-vegeta

clean-stress-test:
	$(MAKE) -C stress-test clean-stress-test
