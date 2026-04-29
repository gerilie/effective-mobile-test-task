# App
dev-up:
	docker compose \
	-p subscription-dev \
	--env-file .env.dev \
	-f docker-compose.yml \
	-f docker-compose.dev.yml \
	up --build -d
dev-down:
	docker compose \
	-p subscription-dev \
	--env-file .env.dev \
	-f docker-compose.yml \
	-f docker-compose.dev.yml \
	down

prod-up:
	docker compose \
	-p subscription-prod \
	--env-file .env.prod \
	-f docker-compose.yml \
	-f docker-compose.prod.yml \
	up --build -d
prod-down:
	docker compose \
	-p subscription-prod \
	--env-file .env.prod \
	-f docker-compose.yml \
	-f docker-compose.prod.yml \
	down

# Tools
build:
	go build -o subscription cmd/subscription/*.go
lint:
	golangci-lint run -v
test:
	go test ./... -v -cover
swagger:
	swag init -g cmd/subscription/main.go