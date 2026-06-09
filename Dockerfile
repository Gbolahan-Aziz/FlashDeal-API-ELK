FROM golang:1.22 AS build
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# CGO=0 by default (Render, no SQLite). Pass CGO=1 via docker-compose for local SQLite support.
ARG CGO=0
RUN CGO_ENABLED=${CGO} GOOS=linux go build -o /flash-api ./cmd/api

# Minimal runtime image
FROM gcr.io/distroless/base-debian12
COPY --from=build /flash-api /flash-api

ENV PORT=8080 SERVICE_NAME=flash-api ENV=dev

EXPOSE 8080
ENTRYPOINT ["/flash-api"]
