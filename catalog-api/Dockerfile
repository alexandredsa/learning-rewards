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
RUN CGO_ENABLED=0 GOOS=linux go build -o catalog-api ./main.go

# Final stage
FROM alpine:latest

# Add build argument to force rebuilds
ARG BUILD_DATE=unknown
ARG VCS_REF=unknown

WORKDIR /app

# Copy the binary from builder
COPY --from=builder /app/catalog-api .
COPY --from=builder /app/config.yaml .

# Expose the application port
EXPOSE 8080

# Run the application
CMD ["./catalog-api"] 