# Build stage
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Install build dependencies
RUN apk add --no-cache git gcc musl-dev sqlite-dev

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build with CGO for SQLite
RUN CGO_ENABLED=1 GOOS=linux go build -a -installsuffix cgo -ldflags '-linkmode external -extldflags "-static"' -o git-migrator ./cmd/git-migrator

# Runtime stage
FROM alpine:latest

# Install runtime dependencies
RUN apk --no-cache add ca-certificates sqlite-libs

WORKDIR /app

# Copy binary from builder
COPY --from=builder /app/git-migrator /usr/local/bin/git-migrator

# Create directories for data
RUN mkdir -p /repos /config /data

# Expose web UI port
EXPOSE 8080

# Set default command
ENTRYPOINT ["git-migrator"]
CMD ["web", "--port", "8080"]
