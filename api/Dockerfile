# Multi-stage Docker build for QuizNinja API
# Optimized for Google Cloud Run deployment

# Stage 1: Build the Go application
FROM golang:1.23.1-alpine AS builder

# Install build dependencies
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
# CGO_ENABLED=0: Build static binary
# -ldflags="-w -s": Strip debug info to reduce size
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o main .

# Stage 2: Runtime image with Nginx
FROM nginx:alpine

# Set working directory
WORKDIR /app

# Install supervisor to manage multiple processes
RUN apk add --no-cache supervisor

# Copy the compiled binary from builder
COPY --from=builder /app/main /app/main

# Copy migration files from builder
COPY --from=builder /app/database/migrations /app/database/migrations

# Copy nginx configuration
COPY nginx.conf /etc/nginx/nginx.conf

# Create supervisor configuration
RUN mkdir -p /etc/supervisor.d
COPY <<EOF /etc/supervisor.d/supervisord.ini
[supervisord]
nodaemon=true
logfile=/dev/null
logfile_maxbytes=0
loglevel=info

[program:nginx]
command=nginx -g "daemon off;"
stdout_logfile=/dev/stdout
stdout_logfile_maxbytes=0
stderr_logfile=/dev/stderr
stderr_logfile_maxbytes=0
autostart=true
autorestart=true
priority=10

[program:quizninja-api]
command=/app/main
stdout_logfile=/dev/stdout
stdout_logfile_maxbytes=0
stderr_logfile=/dev/stderr
stderr_logfile_maxbytes=0
autostart=true
autorestart=true
priority=20
environment=PORT="8081"
EOF

# Expose port 8080 (Cloud Run requirement)
EXPOSE 8080

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1

# Run supervisor to manage both nginx and the Go app
CMD ["/usr/bin/supervisord", "-c", "/etc/supervisor.d/supervisord.ini"]
