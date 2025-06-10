# Build stage
FROM golang:1.24.4-alpine AS builder

# Add build argument to force rebuilds
ARG BUILD_DATE=unknown
ARG VCS_REF=unknown

WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o reward-processor ./main.go

# Final stage
FROM alpine:latest

# Add build argument to force rebuilds
ARG BUILD_DATE=unknown
ARG VCS_REF=unknown

WORKDIR /app

# Copy the binary from builder
COPY --from=builder /app/reward-processor .

# Run the application
CMD ["./reward-processor"] 