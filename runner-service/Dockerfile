FROM golang:1.25 AS builder

WORKDIR /src

COPY shared/ shared/
COPY runner-service/ runner-service/

WORKDIR /src/runner-service
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o /runner ./cmd/server/...

FROM debian:bookworm-slim
RUN apt-get update && apt-get install -y --no-install-recommends \
      git openssh-client ca-certificates wget docker.io \
    && rm -rf /var/lib/apt/lists/*

COPY --from=builder /runner /runner

ENTRYPOINT ["/runner"]
