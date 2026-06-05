# Stage 1: Builder
FROM golang:1.25-alpine AS builder

WORKDIR /build

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build application
RUN go build -o /app ./cmd/server

# Stage 2: Runtime
FROM alpine:3.19

RUN apk add --no-cache ca-certificates tzdata

WORKDIR /app

# Copy binary from builder
COPY --from=builder /app /app/server

# Copy migrations
COPY migrations /app/migrations

# Copy swagger docs
COPY docs /app/docs

EXPOSE 8080

CMD ["/app/server"]
