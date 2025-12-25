FROM golang:1.23-alpine AS builder

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./
COPY services/inventory/go.mod services/inventory/go.sum ./services/inventory/

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build
WORKDIR /app/services/inventory
RUN CGO_ENABLED=0 GOOS=linux go build -o /inventory ./cmd/inventory

FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy binary
COPY --from=builder /inventory .

# Copy configs
COPY services/inventory/configs ./configs

EXPOSE 8000 9000

CMD ["./inventory", "-conf", "./configs"]
