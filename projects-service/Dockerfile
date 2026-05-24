FROM golang:1.25 AS builder
WORKDIR /app

COPY projects-service/go.mod projects-service/go.sum ./
COPY shared /shared
RUN go mod download

COPY projects-service/ ./
RUN CGO_ENABLED=0 go build -o /server ./cmd/server

FROM debian:bookworm-slim
RUN apt-get update && apt-get install -y --no-install-recommends \
      ca-certificates git openssh-client wget \
    && rm -rf /var/lib/apt/lists/*
COPY --from=builder /server /server
EXPOSE 8084 50052
ENTRYPOINT ["/server"]
