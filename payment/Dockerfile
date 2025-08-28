FROM golang:1.24-alpine AS builder

WORKDIR /app

# Copy microservices-proto first
COPY microservices-proto/ ./microservices-proto/

# Copy payment service files
COPY microservices/payment/go.mod microservices/payment/go.sum ./microservices/payment/

# Set working directory to payment service
WORKDIR /app/microservices/payment

# Download dependencies
RUN go mod download

# Copy payment service source code
COPY microservices/payment/ .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o payment ./cmd/

# Final stage
FROM alpine:latest

# Install ca-certificates for HTTPS
RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy the binary from builder stage
COPY --from=builder /app/microservices/payment/payment .

# Expose port
EXPOSE 3001

# Run the binary
CMD ["./payment"]