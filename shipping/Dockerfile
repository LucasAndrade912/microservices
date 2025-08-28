FROM golang:1.24-alpine AS builder

WORKDIR /app

# Copy microservices-proto first
COPY microservices-proto/ ./microservices-proto/

# Copy shipping service files
COPY microservices/shipping/go.mod microservices/shipping/go.sum ./microservices/shipping/

# Set working directory to shipping service
WORKDIR /app/microservices/shipping

# Download dependencies
RUN go mod download

# Copy shipping service source code
COPY microservices/shipping/ .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o shipping ./cmd/

# Final stage
FROM alpine:latest

# Install ca-certificates for HTTPS
RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy the binary from builder stage
COPY --from=builder /app/microservices/shipping/shipping .

# Expose port
EXPOSE 3002

# Run the binary
CMD ["./shipping"]