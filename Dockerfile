# Builder
FROM --platform=$BUILDPLATFORM golang:1.23.5-alpine3.21 AS builder

RUN go install github.com/pressly/goose/v3/cmd/goose@latest

WORKDIR /usr/local/src/auth

COPY go.mod go.sum ./

RUN go mod download

COPY *.go ./

RUN CGO_ENABLED=0 GOOS=linux go build -o /usr/local/bin/auth

# Binary
FROM  --platform=$BUILDPLATFORM alpine:3.21

WORKDIR /usr/local/bin/

COPY --from=builder /go/bin/goose ./
COPY --from=builder /usr/local/bin/auth ./
COPY ./migrations ./migrations

EXPOSE 3000

CMD ["sh", "-c", "./goose up && ./auth"]
