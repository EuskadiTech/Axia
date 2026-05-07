# Multi-stage build for axia4 (Tallarin)
FROM golang:1.26.1-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git ca-certificates tzdata

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the application
ARG VERSION=dev
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags "-X main.appVersion=${VERSION} -s -w" -o axia4 .

# Final stage
FROM alpine:latest

# Install runtime dependencies
RUN apk add --no-cache ca-certificates tzdata postgresql-client imagemagick ghostscript

# Create app directory and user
RUN addgroup -g 1000 axia4 && \
    adduser -u 1000 -G axia4 -s /bin/sh -D axia4 && \
    mkdir -p /app && \
    chown axia4:axia4 /app

WORKDIR /app

# Copy binary and config files
COPY --from=builder /app/axia4 /app/
COPY --from=builder /app/config_template.json /app/
COPY --from=builder /app/config_dedicated.json /app/
COPY --from=builder /app/config_portable.json /app/
COPY --from=builder /app/config_docker.json /app/

# Set permissions
RUN chmod +x /app/axia4

# Switch to non-root user
USER axia4

# Set environment variable to skip fsync for better performance in Docker
# Container storage drivers (overlay, aufs) can make fsync very slow
ENV REI3_SKIP_FSYNC=true

# Expose default ports
EXPOSE 80 443

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD wget --no-verbose --tries=1 --spider http://localhost:80/ || exit 1

CMD ["./axia4", "-run"]