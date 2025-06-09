# Learning-Rewards

A microservices-based project that implements a catalog system with event processing capabilities. The project consists of two main services:

- **Catalog API**: A GraphQL API service that manages the catalog data
- **Event Processor**: A service that handles event processing and data transformation

## Architecture

The project uses a microservices architecture with the following components:

- PostgreSQL databases (one for each service)
- GraphQL API (using gqlgen)
- Event processing service
- Docker for containerization

## Prerequisites

- Go 1.x
- Docker and Docker Compose
- Make

## Getting Started

1. Clone the repository:
```bash
git clone https://github.com/yourusername/learning-rewards.git
cd learning-rewards
```

2. Start the services:
```bash
make
```

This will:
- Start the PostgreSQL databases
- Build both services
- Run the Catalog API on port 8080
- Run the Event Processor on port 8081

## Available Commands

- `make` - Build and run both applications
- `make build` - Build both applications
- `make run` - Start databases and run both applications
- `make stop` - Stop both applications and databases
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

## Development

The project uses environment variables for configuration:

- `DATABASE_DSN`: Database connection string
- `PORT`: Service port (defaults to 8080 for Catalog API and 8081 for Event Processor)
- `ENV`: Environment (e.g., dev, prod)

To stop all services:
```bash
make stop
```

To clean up build artifacts:
```bash
make clean
```