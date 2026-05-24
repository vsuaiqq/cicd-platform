FROM golang:1.25 AS builder
WORKDIR /app
COPY . .
WORKDIR /app/auth-service
ARG CACHEBUST=1
RUN go mod download && CGO_ENABLED=0 go build -o /server ./cmd/server

FROM debian:bookworm-slim
RUN apt-get update && apt-get install -y --no-install-recommends \
      ca-certificates wget \
    && rm -rf /var/lib/apt/lists/*
COPY --from=builder /server /server
EXPOSE 50051 8082
ENTRYPOINT ["/server"]
