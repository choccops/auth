# Builder
FROM arm64v8/golang:1.23.5-alpine3.21 AS builder

WORKDIR /usr/local/src/auth

COPY go.mod go.sum ./

RUN go mod download

COPY *.go ./

RUN CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -o auth

# Binary
FROM arm64v8/alpine:3.21

WORKDIR /usr/local/bin/

COPY --from=builder /usr/local/src/auth/auth ./

EXPOSE 3000

CMD ["./auth"]
