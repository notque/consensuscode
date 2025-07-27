# Multi-stage build for collective efficiency
FROM golang:1.21-alpine AS builder

# Install collective build dependencies
RUN apk add --no-cache git ca-certificates tzdata

# Set working directory
WORKDIR /collective

# Copy go mod files first for better caching
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build with collective optimizations
ARG VERSION=dev
ARG BUILD_DATE
RUN CGO_ENABLED=0 GOOS=linux go build \
    -ldflags="-w -s -X main.version=${VERSION} -X main.buildDate=${BUILD_DATE}" \
    -o bluesky-collective \
    ./cmd/bluesky-collective

# Final stage - minimal collective container
FROM scratch

# Copy from builder
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /collective/bluesky-collective /bluesky-collective

# Non-root user for collective security
USER 65534:65534

# Expose collective API port
EXPOSE 8080

# Health check for collective monitoring
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD ["/bluesky-collective", "health"]

# Run collective service
ENTRYPOINT ["/bluesky-collective"]
CMD ["serve"]