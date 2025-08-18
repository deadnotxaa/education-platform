FROM golang:1.21-alpine AS builder

WORKDIR /app

# Copy go mod file if it exists (for dependencies)
COPY go.* ./

# Download dependencies if go.mod exists
RUN if [ -f go.mod ]; then go mod download; fi

# Copy source code
COPY main.go .

# Build the application
RUN go build -o client main.go

# Final stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy the binary from builder stage
COPY --from=builder /app/client .

# Run the client
CMD ["./client"]
