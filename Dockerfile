# app
FROM golang:alpine AS builder
WORKDIR /app
COPY . .

RUN apk add --no-cache make
RUN go install github.com/swaggo/swag/cmd/swag@latest
RUN go mod download

RUN make swagger
RUN make build

FROM alpine:latest AS runtime
WORKDIR /app

RUN apk add --no-cache ca-certificates tzdata
RUN addgroup -g 1000 appuser && adduser -D -u 1000 -G appuser appuser 

COPY --from=builder /app/subscription /app/subscription
COPY --from=builder /app/config /app/config
COPY --from=builder /app/docs /app/docs

RUN chown -R appuser:appuser /app
USER appuser

FROM runtime AS subscription
EXPOSE 8080
ENTRYPOINT ["sh", "-c", "/app/subscription --env=prod"]

# migrate
FROM golang:alpine AS migrate-builder
RUN go install github.com/pressly/goose/v3/cmd/goose@latest

FROM alpine:latest AS migrate-runtime
WORKDIR /app

RUN apk add --no-cache make 

COPY --from=migrate-builder /go/bin/goose /usr/local/bin/goose
COPY ./config /app/config
COPY ./migrations /app/migrations
COPY ./Makefile /app/Makefile

FROM migrate-runtime AS migrate
ENTRYPOINT ["sh", "-c", "goose up"]