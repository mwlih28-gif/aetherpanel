# Build stage
FROM golang:1.22-alpine AS builder

WORKDIR /app

# Install dependencies
RUN apk add --no-cache git ca-certificates tzdata

# Copy source code
COPY . .

# Copy config file
COPY ../config.example.yaml /app/config.yaml

# Download dependencies and build
RUN go mod tidy && CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o /app/aether-api ./cmd/api

# Final stage
FROM alpine:3.19

WORKDIR /app

# Install runtime dependencies
RUN apk add --no-cache ca-certificates tzdata

# Copy binary and config from builder
COPY --from=builder /app/aether-api /app/aether-api
COPY --from=builder /app/config.yaml /app/config.yaml

# Create non-root user
RUN adduser -D -g '' aether
USER aether

# Expose port
EXPOSE 8080

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1

# Run the application
CMD ["/app/aether-api"]
