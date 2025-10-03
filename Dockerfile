# Multi-stage build for REI3 (Tallarin)
FROM golang:1.24.4-alpine AS builder

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
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags "-X main.appVersion=${VERSION} -s -w" -o r3 .

# Final stage
FROM alpine:latest

# Install runtime dependencies
RUN apk add --no-cache ca-certificates tzdata postgresql-client imagemagick ghostscript

# Create app directory and user
RUN addgroup -g 1000 rei3 && \
    adduser -u 1000 -G rei3 -s /bin/sh -D rei3 && \
    mkdir -p /app && \
    chown rei3:rei3 /app

WORKDIR /app

# Copy binary and config files
COPY --from=builder /app/r3 /app/
COPY --from=builder /app/config_template.json /app/
COPY --from=builder /app/config_dedicated.json /app/
COPY --from=builder /app/config_portable.json /app/
COPY --from=builder /app/config_docker.json /app/

# Set permissions
RUN chmod +x /app/r3

# Switch to non-root user
USER rei3

# Expose default ports
EXPOSE 80 443

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD wget --no-verbose --tries=1 --spider http://localhost:80/ || exit 1

CMD ["./r3", "-run"]