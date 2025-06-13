# Build stage
FROM golang:1.24.4-alpine AS builder

# Add build arguments
ARG BUILD_DATE=unknown
ARG VCS_REF=unknown
ARG SERVICE=api  # Default to api service

WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build both binaries
RUN CGO_ENABLED=0 GOOS=linux go build -o reward-processor-api ./cmd/api/main.go
RUN CGO_ENABLED=0 GOOS=linux go build -o reward-processor-worker ./cmd/worker/main.go

# Final stage
FROM alpine:latest

# Add build arguments
ARG BUILD_DATE=unknown
ARG VCS_REF=unknown
ARG SERVICE=api  # Default to api service

WORKDIR /app

# Copy both binaries from builder
COPY --from=builder /app/reward-processor-api .
COPY --from=builder /app/reward-processor-worker .

# Use an entrypoint script to select which binary to run
COPY <<EOF /app/entrypoint.sh
#!/bin/sh
if [ "$SERVICE" = "api" ]; then
    exec ./reward-processor-api
elif [ "$SERVICE" = "worker" ]; then
    exec ./reward-processor-worker
else
    echo "Invalid SERVICE argument. Must be either 'api' or 'worker'"
    exit 1
fi
EOF

RUN chmod +x /app/entrypoint.sh

# Run the selected service
ENTRYPOINT ["/app/entrypoint.sh"] 