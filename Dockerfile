# App
FROM golang:alpine AS dev
WORKDIR /app

RUN apk add --no-cache make git
RUN go install github.com/air-verse/air@latest

COPY . .

CMD ["air"]

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

ENV TZ=UTC

RUN addgroup -g 1000 appuser && adduser -D -u 1000 -G appuser appuser 

COPY --from=builder /app/subscription /app/subscription
COPY --from=builder /app/docs /app/docs

RUN chown -R appuser:appuser /app
USER appuser

FROM runtime AS subscription
EXPOSE 8080
ENTRYPOINT ["/app/subscription"]
CMD ["--env", "prod"]

# migrate
FROM golang:alpine AS migrate-builder
RUN go install github.com/pressly/goose/v3/cmd/goose@latest

FROM alpine:latest AS migrate-runtime
WORKDIR /app

COPY --from=migrate-builder /go/bin/goose /usr/local/bin/goose
COPY ./migrations /app/migrations

FROM migrate-runtime AS migrate
ENTRYPOINT ["sh", "-c", "goose up"]