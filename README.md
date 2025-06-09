# Learning-Rewards

A microservices-based project that implements a catalog system with event processing capabilities. The project consists of two main services:

- **Catalog API**: A GraphQL API service that manages the catalog data
- **Event Processor**: A service that handles event processing and publishes events to Kafka

## Architecture

The project uses a microservices architecture with the following components:

- PostgreSQL databases (one for each service)
- GraphQL API (using gqlgen)
- Event processing service with Kafka integration
- Docker for containerization

### Event Processing Flow

1. Events are received via HTTP POST requests to the Event Processor
2. Events are validated and transformed into a standardized format
3. Events are published to Kafka topics for further processing
4. Other services can consume these events for their specific needs

## Prerequisites

- Go 1.x
- Docker and Docker Compose
- Make

## Getting Started

1. Clone the repository:
```bash
git clone https://github.com/alexandredsa/learning-rewards.git
cd learning-rewards
```

2. Start the services:
```bash
make
```

This will:
- Start the PostgreSQL databases
- Start Kafka and Kafka UI
- Build both services
- Run the Catalog API on port 8080
- Run the Event Processor on port 8081

## Available Commands

- `make` - Build and run both applications
- `make build` - Build both applications
- `make run` - Start databases, Kafka, and run both applications
- `make stop` - Stop both applications, databases, and Kafka
- `make clean` - Clean up build artifacts
- `make help` - Show all available commands
- `make stress-test` – Run Vegeta load test
- `make install-vegeta` – Install Vegeta CLI
- `make clean-stress-test` – Clean stress test results

## Services

### Catalog API
- Runs on port 8080
- GraphQL playground available at http://localhost:8080
- GraphQL endpoint at http://localhost:8080/query
- Uses PostgreSQL database on port 5432

### Event Processor
- Runs on port 8081
- Health check endpoint at http://localhost:8081/health
- Events endpoint at http://localhost:8081/events
- Uses PostgreSQL database on port 5433
- Publishes events to Kafka topic `platform-events`

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

The project uses environment variables for configuration:

### Event Processor Environment Variables
- `DATABASE_DSN`: Database connection string
- `PORT`: Service port (defaults to 8081)
- `KAFKA_BROKERS`: Comma-separated list of Kafka broker addresses (defaults to "localhost:29092")
- `KAFKA_TOPIC`: Kafka topic name (defaults to "platform-events")
- `DEBUG`: Enable debug logging when set

### Catalog API Environment Variables
- `DATABASE_DSN`: Database connection string
- `PORT`: Service port (defaults to 8080)
- `ENV`: Environment (e.g., dev, prod)

To stop all services:
```bash
make stop
```

To clean up build artifacts:
```bash
make clean
```

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