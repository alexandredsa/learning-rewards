# Event Processor

A Go service that processes learning platform events and publishes them to Kafka. This service acts as an event gateway, receiving events via HTTP and publishing them to a Kafka topic for further processing by other services.

## Features

- HTTP API for receiving events
- Kafka event publishing
- Health check endpoint
- Configurable via environment variables
- Graceful shutdown handling

## Prerequisites

- Go 1.21 or later
- Access to a Kafka broker (see root project for Kafka setup)

## Quick Start

1. Clone the repository:
```bash
git clone https://github.com/yourusername/event-processor.git
cd event-processor
```

2. Install dependencies:
```bash
go mod download
```

3. Run the service:
```bash
go run main.go
```

The service will start on port 8081 by default and connect to Kafka at localhost:9092.

## Environment Variables

The service can be configured using the following environment variables:

| Variable | Description | Default |
|----------|-------------|---------|
| `PORT` | HTTP server port | `8081` |
| `KAFKA_BROKERS` | Comma-separated list of Kafka broker addresses | `localhost:9092` |
| `KAFKA_TOPIC` | Kafka topic name for events | `platform-events` |

Example:
```bash
export PORT=8080
export KAFKA_BROKERS=kafka1:9092,kafka2:9092
export KAFKA_TOPIC=learning-events
go run main.go
```

## API Endpoints

### POST /events
Publishes a learning event to Kafka.

Request body:
```json
{
    "user_id": "user123",
    "event_type": "course_completed",
    "course_id": "course456",
    "timestamp": "2024-03-20T10:00:00Z"
}
```

Response:
- 202 Accepted: Event was successfully published
- 400 Bad Request: Invalid request body
- 500 Internal Server Error: Failed to publish event

### GET /health
Health check endpoint.

Response:
- 200 OK: Service is healthy


## Development

### Running Tests
```bash
go test ./...
```

### Building
```bash
go build -o bin/event-processor
```

## License

MIT License 