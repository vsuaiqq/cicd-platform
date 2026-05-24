FROM golang:1.25 AS builder
WORKDIR /app
COPY . .
WORKDIR /app/api-gateway
RUN go mod download && CGO_ENABLED=0 go build -o /server ./cmd/server

FROM debian:bookworm-slim
RUN apt-get update && apt-get install -y --no-install-recommends \
      ca-certificates wget \
    && rm -rf /var/lib/apt/lists/*
COPY --from=builder /server /server
EXPOSE 8080
ENTRYPOINT ["/server"]
