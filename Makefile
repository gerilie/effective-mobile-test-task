all: build
	./subscription --env=dev

# Docker
docker-build:
	docker compose --env-file config/prod-postgres.env up --build -d
docker-up:
	docker compose --env-file config/prod-postgres.env up -d 
docker-down:
	docker compose --env-file config/prod-postgres.env down

# Migration
migrate-generate:
	goose create $(name) sql -env config/$(env)-migrate.env
migrate-up:
	goose up -env config/$(env)-migrate.env
migrate-down:
	goose down -env config/$(env)-migrate.env

# Build
build:
	go build -o subscription cmd/subscription/*.go

# Run
local:
	go run cmd/subscription/*.go --env=dev
lint:
	golangci-lint run -v
swagger:
	swag init -g cmd/subscription/main.go