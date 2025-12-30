# Build stage
FROM golang:1.21-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-s -w' -o auth-module .

# Final stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy the binary from builder
COPY --from=builder /app/auth-module .

# Copy config example (config.yaml will be mounted via docker-compose)
COPY --from=builder /app/config.yaml.example ./config.yaml.example

# Expose port
EXPOSE 8080

# Run the application
CMD ["./auth-module", "serve"]
