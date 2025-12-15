# Build stage - Use latest Go with auto toolchain
FROM golang:alpine AS builder

WORKDIR /app

# Copy everything
COPY . .

# Build - let Go download the right toolchain automatically
RUN CGO_ENABLED=0 GOOS=linux go build \
    -mod=vendor \
    -ldflags="-s -w" \
    -o /app/bin/api \
    ./cmd/api

# Production stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /app

COPY --from=builder /app/bin/api .

EXPOSE 8080

CMD ["./api"]
