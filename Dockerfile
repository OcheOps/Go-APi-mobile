# Build stage
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Install git for dependencies (if any needed, though we vendor or are pure go)
RUN apk add --no-cache git

# Cache dependencies
COPY go.mod go.sum ./
RUN go mod download

# Build
COPY . .
# CGO_ENABLED=0 is default for modernc.org/sqlite on pure go, but let's be explicit
RUN go mod tidy
RUN CGO_ENABLED=0 GOOS=linux go build -o fieldnotes .

# Runtime stage
FROM alpine:latest

WORKDIR /app

# Copy binary
COPY --from=builder /app/fieldnotes .

# Expose port
EXPOSE 8080

# Explicitly create a directory for data if we want to mount it
RUN mkdir -p /data

# Run
CMD ["./fieldnotes"]
