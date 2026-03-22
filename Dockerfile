# ---------- Build Stage ----------
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Install dependencies
RUN apk add --no-cache git protobuf

# Install protoc plugins
RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
RUN go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

ENV PATH="$PATH:/go/bin"

# Cache dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy project
COPY . .

# Generate proto files
RUN protoc \
    --go_out=. \
    --go-grpc_out=. \
    proto/*.proto

# Build optimized static binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -ldflags="-s -w" -o file-service ./cmd/server


# ---------- Runtime Stage ----------
FROM alpine:3.19

WORKDIR /app

# Create non-root user + uploads directory with permission
RUN adduser -D appuser && \
    mkdir -p /app/uploads && \
    chown -R appuser:appuser /app

# Copy binary
COPY --from=builder /app/file-service .

# Ensure ownership
RUN chown appuser:appuser /app/file-service

# Switch to non-root user
USER appuser

EXPOSE 50051

CMD ["./file-service"]