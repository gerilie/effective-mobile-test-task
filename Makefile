# App
dev-up:
	@env=dev $(MAKE) docker-build
prod-up:
	@env=prod $(MAKE) docker-build
host-up:
	$(MAKE) host-db-up
	$(MAKE) host-migrate-up
	go run cmd/subscription/*.go --env=host

# Docker
docker-build:
	docker compose --env-file config/$(env)-postgres.env up --build -d
docker-up:
	docker compose --env-file config/$(env)-postgres.env up -d 
docker-down:
	docker compose --env-file config/$(env)-postgres.env down

# Host
host-db-up:
	docker compose --env-file config/host-postgres.env up -d postgres-host
host-db-down:
	docker compose --env-file config/host-postgres.env down
host-migrate-gen:
	goose create $(name) sql -env config/host-migrate.env
host-migrate-up:
	goose up sql -env config/host-migrate.env
host-migrate-down:
	goose down sql -env config/host-migrate.env

# Tools
build:
	go build -o subscription cmd/subscription/*.go
lint:
	golangci-lint run -v
test:
	go test ./... -v -cover
swagger:
	swag init -g cmd/subscription/main.go