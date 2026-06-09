# Build stage
FROM golang:1.24-alpine AS builder
WORKDIR /app

# Copy go.mod first
COPY go.mod ./
# Copy go.sum if it exists, otherwise generate it
COPY go.sum* ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o server ./cmd/server

# Runtime stage
FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/server .
EXPOSE 8080
CMD ["./server"]