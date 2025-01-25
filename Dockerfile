# Builder
FROM --platform=$BUILDPLATFORM golang:1.23.5-alpine3.21 AS builder

WORKDIR /usr/local/src/auth

COPY go.mod go.sum ./

RUN go mod download

COPY *.go ./
COPY ./internal ./internal
COPY ./migrations ./migrations

RUN CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -o /usr/local/bin/auth

# Binary
FROM --platform=$BUILDPLATFORM alpine:3.21

WORKDIR /usr/local/bin/

COPY --from=builder /usr/local/bin/auth ./

EXPOSE 3000

CMD ["./auth"]
