FROM golang:1.25-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o conntest .

FROM alpine:3.22
RUN addgroup -S conntest && adduser -S conntest -G conntest
USER conntest
COPY --from=builder /app/conntest /usr/local/bin/conntest
ENTRYPOINT ["conntest"]
