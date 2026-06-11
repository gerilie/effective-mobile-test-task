.PHONY: all clean

all: lint test build
clean:
	rm -f subscription

# App
.PHONY: dev-up dev-down prod-up prod-down

dev-up:
	docker compose -p subscription-dev --env-file .env.dev \
	-f docker-compose.yml -f docker-compose.dev.yml up --build -d
dev-down:
	docker compose -p subscription-dev --env-file .env.dev \
	-f docker-compose.yml -f docker-compose.dev.yml down
prod-up:
	docker compose -p subscription-prod --env-file .env.prod \
	-f docker-compose.yml -f docker-compose.prod.yml up --build -d
prod-down:
	docker compose -p subscription-prod --env-file .env.prod \
	-f docker-compose.yml -f docker-compose.prod.yml down

# Tools
.PHONY: build lint test go-gen migrate-gen swagger git-hooks

build:
	go build -o subscription cmd/subscription/*.go
lint:
	golangci-lint run -v
gen:
	go generate ./...

test:
	go test ./... -v -cover -race -tags=integration
test-unit:
	go test ./... -v -cover -race -short
test-integration:
	go test ./... -v -race -tags=integration

migrate-gen:
	goose create $(name) sql --dir ./migrations
swagger:
	swag init -g cmd/subscription/main.go
git-hooks:
	prek install
