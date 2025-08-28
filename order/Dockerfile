FROM golang:1.24-alpine AS builder

WORKDIR /app

# Copy microservices-proto first
COPY microservices-proto/ ./microservices-proto/

# Copy order service files
COPY microservices/order/go.mod microservices/order/go.sum ./microservices/order/

# Set working directory to order service
WORKDIR /app/microservices/order

# Download dependencies
RUN go mod download

# Copy order service source code
COPY microservices/order/ .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/

# Final stage
FROM alpine:latest

# Install ca-certificates for HTTPS
RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy the binary from builder stage
COPY --from=builder /app/microservices/order/main .

# Expose port
EXPOSE 3000

# Run the binary
CMD ["./main"]
