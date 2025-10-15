# Build stage
FROM golang:1.23.4-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git make

# Set working directory
WORKDIR /build

# Copy go module files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the MCP server binary with version information
ARG VERSION=dev
ARG GIT_COMMIT=unknown
ARG BUILD_DATE=unknown
RUN go build \
    -ldflags="-X github.com/chronosphereio/chronosphere-mcp/pkg/version.Version=${VERSION} \
              -X github.com/chronosphereio/chronosphere-mcp/pkg/version.BuildDate=${BUILD_DATE} \
              -X github.com/chronosphereio/chronosphere-mcp/pkg/version.GitCommit=${GIT_COMMIT}" \
    -o /build/bin/chronomcp \
    ./mcp-server

# Runtime stage
FROM alpine:3.20

# Install runtime dependencies
RUN apk add --no-cache ca-certificates && \
    addgroup -g 1001 chronosphere && \
    adduser -D -s /bin/sh -u 1001 -G chronosphere chronosphere

# Set working directory
WORKDIR /app

# Copy binary from builder stage
COPY --from=builder /build/bin/chronomcp /app/chronomcp

# Copy default configuration file (optional - can be overridden via volume mount)
COPY config.http.yaml /app/config.yaml

# Change ownership to non-root user
RUN chown -R chronosphere:chronosphere /app

# Switch to non-root user
USER chronosphere

# Expose HTTP MCP server port (default from config.http.yaml)
EXPOSE 8081

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD pgrep chronomcp || exit 1

# Set default command
# Users can override config file via: docker run -v $(pwd)/custom-config.yaml:/app/config.yaml
ENTRYPOINT ["/app/chronomcp"]
CMD ["--config", "/app/config.yaml"]
