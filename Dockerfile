# Builder
FROM --platform=$BUILDPLATFORM golang:1.23.5-alpine3.21 AS builder

WORKDIR /usr/local/src/auth

COPY go.mod go.sum ./

RUN go mod download

COPY *.go ./

RUN CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -o /usr/local/bin/auth

# Binary
FROM --platform=$BUILDPLATFORM alpine:3.21

WORKDIR /usr/local/bin/

RUN apk add --no-cache curl && \
    curl -L https://github.com/pressly/goose/releases/download/v3.24.1/goose_darwin_arm64 -o goose && \
    chmod +x goose

COPY --from=builder /usr/local/bin/auth ./

EXPOSE 3000

CMD ["sh", "-c", "./goose && ./auth"]
