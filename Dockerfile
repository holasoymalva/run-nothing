# Build stage
FROM golang:1.24-alpine AS builder

WORKDIR /src

COPY go.mod ./
COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /bin/run-nothing .

# Minimal runtime
FROM alpine:latest

COPY --from=builder /bin/run-nothing /usr/local/bin/run-nothing

ENTRYPOINT ["run-nothing"]
