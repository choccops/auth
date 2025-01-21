# Builder
FROM golang:1.23.5-alpine AS builder

WORKDIR /usr/local/src/auth

COPY go.mod go.sum ./

RUN go mod download

COPY *.go ./

RUN go install github.com/pressly/goose/v3/cmd/goose@latest
RUN CGO_ENABLED=0 GOOS=linux go build -o auth

# Go
FROM alpine:latest

WORKDIR /usr/local/bin/

COPY --from=builder /usr/local/src/auth/auth ./
COPY --from=builder /go/bin/goose ./

EXPOSE 3000

CMD ["./auth"]
