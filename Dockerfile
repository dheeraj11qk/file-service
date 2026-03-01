# ---------- Build Stage ----------
FROM golang:1.24.0-alpine AS builder

WORKDIR /app

# Install dependencies
RUN apk add --no-cache git protobuf

# Install protoc plugins
RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
RUN go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

ENV PATH="$PATH:/go/bin"

# Copy go files
COPY go.mod go.sum ./
RUN go mod download

# Copy entire project
COPY . .

# Generate proto files
RUN protoc \
    --go_out=. \
    --go-grpc_out=. \
    proto/*.proto

# Build binary
RUN go build -o main ./cmd/server


# ---------- Runtime Stage ----------
FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/main .

EXPOSE 50051

CMD ["./main"]