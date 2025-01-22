# Builder
FROM golang:1.23.5-alpine AS builder

WORKDIR /usr/local/src/auth

COPY go.mod go.sum ./

RUN go mod download

COPY *.go ./

RUN CGO_ENABLED=0 GOOS=linux go build -o auth

# Binary
FROM alpine:latest

WORKDIR /usr/local/bin/

COPY --from=builder /usr/local/src/auth/auth ./

EXPOSE 3000

CMD ["./auth"]
