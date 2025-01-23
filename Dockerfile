# Builder
FROM --platform=$BUILDPLATFORM golang:1.23.5-alpine3.21 AS builder

WORKDIR /usr/local/src/auth

COPY go.mod go.sum ./

RUN go mod download

COPY *.go ./

RUN go install github.com/pressly/goose/v3/cmd/goose@latest
RUN CGO_ENABLED=0 GOOS=linux go build -o /usr/local/bin/auth

# Binary
FROM  --platform=$BUILDPLATFORM alpine:3.21

WORKDIR /usr/local/bin/

COPY --from=builder /go/bin/goose ./
COPY --from=builder /usr/local/bin/auth ./
COPY ./migrations ./migrations

EXPOSE 3000

CMD ["/usr/local/bin/goose up && /usr/local/bin/auth"]
