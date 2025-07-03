# Build stage
FROM golang:1.23-alpine AS builder

# Set working directory
WORKDIR /app

# Copy go mod and sum files first to cache dependencies
COPY go.mod go.sum ./

# Download all dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the application
# The -ldflags="-s -w" strips debug information to make the binary smaller
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /app/bin/api ./cmd/api

# Final stage
FROM alpine:3.18

# Install CA certificates for HTTPS requests
RUN apk --no-cache add ca-certificates

# Set working directory
WORKDIR /app

# Copy the binary from builder
COPY --from=builder /app/bin/api /app/api

# Copy .env file (if you have one)
# COPY .env .

# Expose the port the app runs on
EXPOSE 8080

# Command to run the executable
CMD ["/app/api"]