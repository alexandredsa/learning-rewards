# Reward Processor

A microservice that processes user learning events and triggers rewards based on defined rules.

## Features

- Consumes user events from Kafka topic `learning-events`
- Processes events against predefined rules
- Supports two types of rules:
  - `SINGLE_EVENT`: Triggers on a single matching event
  - `MILESTONE`: Tracks event counts and triggers when a target is reached
- Publishes reward events to Kafka topic `user-rewards`
- Persistent milestone tracking using PostgreSQL
- Graceful shutdown handling
- Structured logging

## Configuration

The service can be configured using environment variables:

### Kafka Configuration
- `KAFKA_BROKERS`: Comma-separated list of Kafka broker addresses (default: "localhost:29092")
- `KAFKA_CONSUMER_GROUP`: Kafka consumer group name (default: "reward-processor")
- `KAFKA_CONSUMER_TOPICS`: Comma-separated list of topics to consume (default: "learning-events")
- `KAFKA_PRODUCER_TOPIC`: Topic to publish reward events (default: "user-rewards")

### Database Configuration
- `DB_HOST`: PostgreSQL host address (default: "localhost")
- `DB_PORT`: PostgreSQL port (default: 5432)
- `DB_USER`: Database user (default: "postgres")
- `DB_PASSWORD`: Database password (default: "postgres")
- `DB_NAME`: Database name (default: "rewards")
- `DB_SSL_MODE`: SSL mode for database connection (default: "disable")

### Other Configuration
- `LOG_LEVEL`: Logging level (default: "info")
- `ENV`: Environment name, affects logging format (default: "development")

## Rule Types

### Single Event Rule

Triggers a reward when a single event matches the conditions:

```json
{
  "id": "rule-001",
  "type": "SINGLE_EVENT",
  "event_type": "COURSE_COMPLETED",
  "conditions": {
    "category": "MATH"
  },
  "reward": {
    "type": "BADGE",
    "description": "Finished a Math course"
  },
  "enabled": true
}
```

### Milestone Rule

Tracks event counts in the database and triggers when the target is reached:

```json
{
  "id": "rule-002",
  "type": "MILESTONE",
  "event_type": "COURSE_COMPLETED",
  "count": 5,
  "conditions": {
    "category": "MATH"
  },
  "reward": {
    "type": "POINTS",
    "amount": 100,
    "description": "Completed 5 math courses"
  },
  "enabled": true
}
```

## Event Schema

### Input Event (learning-events topic)

```json
{
  "user_id": "abc-123",
  "event_type": "COURSE_COMPLETED",
  "course_id": "course-xyz",
  "category": "MATH",
  "timestamp": "2025-06-03T14:00:00Z"
}
```

### Output Event (user-rewards topic)

```json
{
  "user_id": "abc-123",
  "rule_id": "rule-002",
  "reward": {
    "type": "POINTS",
    "amount": 100,
    "description": "Completed 5 math courses"
  },
  "timestamp": "2025-06-09T20:00:00Z"
}
```

## Local Development

1. Install Go 1.21 or later
2. Install PostgreSQL 12 or later
3. Create a database:
   ```sql
   CREATE DATABASE rewards;
   ```
4. Install dependencies:
   ```bash
   go mod download
   ```
5. Run the service:
   ```bash
   go run cmd/reward-processor/main.go
   ```

## Testing

Run the test suite:
```bash
go test ./...
```

The test suite includes:
- Unit tests for the rules engine
- Integration tests for database operations
- Mock tests for Kafka interactions
