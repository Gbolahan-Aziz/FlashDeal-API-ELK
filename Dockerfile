FROM golang:1.22 AS build
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -o /flash-api ./cmd/api

FROM debian:bookworm-slim


RUN apt-get update && apt-get install -y ca-certificates && rm -rf /var/lib/apt/lists/*

RUN useradd -m -u 1001 appuser

RUN mkdir -p /data

RUN chown -R appuser:appuser /data

WORKDIR /

COPY --from=build /flash-api /flash-api

ENV PORT=8080 SERVICE_NAME=flash-api ENV=dev DB_PATH=/data/flashdeals.db

EXPOSE 8080

USER appuser

ENTRYPOINT ["/flash-api"]