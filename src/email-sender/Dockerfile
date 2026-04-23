
FROM golang:1.26.2-alpine AS builder

RUN apk add --no-cache \
    git \
    ca-certificates \
    gcc \
    musl-dev \
    librdkafka-dev \
    pkgconf

WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 \
    go build \
      -tags musl \
      -trimpath \
      -ldflags="-s -w" \
      -o /email-sender \
      ./cmd/service

FROM alpine:3.23

RUN apk add --no-cache curl ca-certificates

COPY --from=builder /email-sender /email-sender

USER 65532:65532

ENTRYPOINT ["/email-sender"]