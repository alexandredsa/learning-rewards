# Learning-Rewards

A microservices-based project that implements a catalog system with event processing and reward management capabilities. The project consists of three main services:

- **Catalog API**: A GraphQL API service that manages the catalog data
- **Event Processor**: A service that handles event processing and publishes events to Kafka
- **Reward Processor**: A service that processes events to manage and distribute rewards to users

## Architecture

The project uses a microservices architecture with the following components:

- PostgreSQL databases (one for each service)
- GraphQL API (using gqlgen)
- Event processing service with Kafka integration
- Reward processing service for managing user rewards
- Docker for containerization

### Event Processing Flow

1. Events are received via HTTP POST requests to the Event Processor
2. Events are validated and transformed into a standardized format
3. Events are published to Kafka topics for further processing
4. The Reward Processor consumes these events to manage user rewards
5. Other services can consume these events for their specific needs

## Prerequisites

- Docker and Docker Compose
- Make (optional, for using Makefile commands)

## Getting Started

1. Clone the repository:
```bash
git clone https://github.com/alexandredsa/learning-rewards.git
cd learning-rewards
```

2. Start all services using Docker Compose:
```bash
make
```

Or manually:
```bash
docker-compose up -d
```

This will:
- Build all service images
- Start the PostgreSQL databases
- Start Kafka and Kafka UI
- Start all services in the correct order

## Available Commands

- `make` - Build and run all services
- `make build` - Build all Docker images (no cache)
- `make dev` - Build and run in development mode (forces rebuild of all services)
- `make run` - Start all services using Docker Compose
- `make stop` - Stop all services
- `make clean` - Clean up Docker resources and remove unused images
- `make rebuild SERVICE=<name>` - Rebuild and restart a specific service
  - Available services: catalog-api, event-processor, reward-processor
- `make help` - Show all available commands
- `make stress-test` – Run Vegeta load test
- `make install-vegeta` – Install Vegeta CLI
- `make clean-stress-test` – Clean stress test results

## Development Workflow

### Starting Development

1. Start all services in development mode:
```bash
make dev
```
This will start all services with live logs and force rebuilds.

### Making Changes

When working on a specific service, you can rebuild just that service:

```bash
# Rebuild catalog-api
make rebuild SERVICE=catalog-api

# Rebuild event-processor
make rebuild SERVICE=event-processor

# Rebuild reward-processor
make rebuild SERVICE=reward-processor
```

### Viewing Logs

To view logs for a specific service:
```bash
docker-compose logs -f <service-name>
```

For example:
```bash
# View catalog-api logs
docker-compose logs -f catalog-api

# View event-processor logs
docker-compose logs -f event-processor

# View reward-processor logs
docker-compose logs -f reward-processor
```

### Clean Development Environment

To start with a clean environment:
```bash
make clean  # Stops all services and removes volumes
make dev    # Rebuilds and starts all services
```

## Services

### Catalog API
- Runs on port 8080
- GraphQL playground available at http://localhost:8080
- GraphQL endpoint at http://localhost:8080/query
- Uses PostgreSQL database 'catalog-api'

### Event Processor
- Runs on port 8081
- Health check endpoint at http://localhost:8081/health
- Events endpoint at http://localhost:8081/events
- Uses PostgreSQL database 'rewards'
- Publishes events to Kafka topic `learning-events`

### Reward Processor
- Runs on port 8082
- Health check endpoint at http://localhost:8082/health
- Consumes events from Kafka topic `learning-events`
- Uses PostgreSQL database 'rewards'

### PostgreSQL
- Single instance running on port 5432
- Contains two databases:
  - `catalog-api`: Used by the Catalog API service
  - `rewards`: Used by both Event Processor and Reward Processor services
- Credentials:
  - User: user
  - Password: pass

### Kafka
- Runs in KRaft mode (no ZooKeeper dependency)
- Broker available at:
  - Internal: `kafka:9092` (for container-to-container communication)
  - External: `localhost:29092` (for host machine access)
- Kafka UI available at http://localhost:9094

## Event Format

Events are published to Kafka in the following JSON format:

```json
{
    "id": "uuid",
    "user_id": "string",
    "event_type": "string",
    "course_id": "string",
    "timestamp": "ISO8601 datetime",
    "created_at": "ISO8601 datetime"
}
```

## Development

The project uses Docker Compose for service orchestration and environment configuration. All services are containerized and can be managed using Docker Compose commands.

### Environment Configuration

All environment variables are configured in the `docker-compose.yml` file. The services use the following environment variables:

#### Catalog API
- `DATABASE_DSN`: Database connection string (format: postgres://user:pass@postgres:5432/catalog-api?sslmode=disable)
- `PORT`: Service port (default: 8080)

#### Event Processor
- `DATABASE_DSN`: Database connection string (format: postgres://user:pass@postgres:5432/rewards?sslmode=disable)
- `PORT`: Service port (default: 8081)
- `KAFKA_BROKERS`: Kafka broker address (default: kafka:9092)

#### Reward Processor
- `DATABASE_DSN`: Database connection string (format: postgres://user:pass@postgres:5432/rewards?sslmode=disable)
- `PORT`: Service port (default: 8082)
- `KAFKA_BROKERS`: Kafka broker address (default: kafka:9092)
- `KAFKA_CONSUMER_TOPICS`: Kafka topic to consume from (default: learning-events)
- `KAFKA_CONSUMER_GROUP`: Kafka consumer group name (default: reward-processor)

## Testing the Event System

You can test the event system by sending a POST request to the Event Processor:

```bash
curl -X POST http://localhost:8081/events \
  -H "Content-Type: application/json" \
  -d '{
    "user_id": "user123",
    "event_type": "course_completed",
    "course_id": "course456",
    "timestamp": "2024-03-20T10:00:00Z"
  }'
```

You can monitor the events in Kafka using the Kafka UI at http://localhost:9094

## Docker Commands

Here are some useful Docker commands for development:

```bash
# View logs for all services
docker-compose logs -f

# View logs for a specific service
docker-compose logs -f <service-name>

# Rebuild and restart a specific service
docker-compose up -d --build --force-recreate <service-name>

# Access a service's shell
docker-compose exec <service-name> sh

# View running containers
docker-compose ps
```

## Troubleshooting

### Service Not Starting
If a service fails to start:
1. Check the logs: `docker-compose logs -f <service-name>`
2. Ensure all dependencies are running: `docker-compose ps`
3. Try rebuilding the service: `make rebuild SERVICE=<service-name>`

### Database Issues
If you encounter database issues:
1. Check if PostgreSQL is running: `docker-compose ps postgres`
2. View PostgreSQL logs: `docker-compose logs -f postgres`
3. Try cleaning and restarting: `make clean && make dev`

### Kafka Issues
If you encounter Kafka issues:
1. Check if Kafka is running: `docker-compose ps kafka`
2. View Kafka logs: `docker-compose logs -f kafka`
3. Check Kafka UI at http://localhost:9094
4. Try cleaning and restarting: `make clean && make dev`