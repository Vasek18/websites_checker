# Multi-stage build for Go application
FROM golang:1.21-alpine AS builder

# Install git and ca-certificates for dependency downloads
RUN apk add --no-cache git ca-certificates

# Set working directory
WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the monitor application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o monitor ./cmd/monitor

# Build the migrate application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o migrate ./cmd/migrate

# Final stage
FROM alpine:latest

# Install ca-certificates for HTTPS requests
RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy binaries from builder stage
COPY --from=builder /app/monitor .
COPY --from=builder /app/migrate .

# Copy configuration files # todo probably remove
COPY urls.yaml .

# Create directory for logs # todo probably remove
RUN mkdir -p /var/log/monitor

# Default command runs the monitor
CMD ["./monitor"]