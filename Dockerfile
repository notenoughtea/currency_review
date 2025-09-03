FROM golang:1.24.4-alpine AS builder
ENV CGO_ENABLED=0 GOTOOLCHAIN=auto
WORKDIR /app
RUN apk add --no-cache ca-certificates protoc git
COPY go.mod go.sum ./
RUN go mod download
RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
RUN go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
ENV PATH="/root/go/bin:${PATH}"
COPY . .
RUN protoc -I proto \
  --go_out=./pkg --go_opt=paths=source_relative \
  --go-grpc_out=./pkg --go-grpc_opt=paths=source_relative \
  proto/currency.proto
RUN go build -o /bin/currency ./currency/cmd/currency
RUN go build -o /bin/cron ./currency/cmd/cron
RUN go build -o /bin/migrator ./currency/cmd/migrator
RUN go build -o /bin/gateway ./gateway/cmd/gateway

FROM alpine:3.20
WORKDIR /app
RUN apk add --no-cache ca-certificates netcat-openbsd
COPY --from=builder /bin/currency /bin/cron /bin/migrator /bin/gateway /bin/
COPY config.yaml /app/config.yaml
ENV CONFIG_PATH=/app/config.yaml
