# Multi-stage build for Network Sniffing Service

# Build stage
FROM golang:1.21-alpine AS builder

# Install git and ca-certificates (needed for go mod download)
RUN apk add --no-cache git ca-certificates tzdata

# Set working directory
WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o network-sniffer ./cmd/server

# Final stage
FROM alpine:latest

# Install ca-certificates for HTTPS requests
RUN apk --no-cache add ca-certificates tzdata

# Set working directory
WORKDIR /app

# Copy binary from builder stage
COPY --from=builder /app/network-sniffer .

# Copy environment file
COPY --from=builder /app/.env.development .env.development

# Set default environment variables
ENV STORAGE_MAX_SIZE=1000
ENV SNIFFING_INTERVAL=5s
ENV SERVER_PORT=8080
ENV SERVER_SHUTDOWN_TIMEOUT=30s

# Expose port
EXPOSE 8080

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1

# Run the application
CMD ["./network-sniffer"]
