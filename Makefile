.PHONY: run build test swag

ifneq (,$(wildcard .env))
    include .env
    export
endif

MIGRATE=migrate -path ./migrations -database "postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=$(DB_SSLMODE)"

.PHONY: migrate-up
migrate-up:
	$(MIGRATE) up

.PHONY: migrate-down
migrate-down:
	$(MIGRATE) down

run:
	go run ./cmd/api

build:
	go build ./cmd/api

test:
	go test ./...

swag:
	go run github.com/swaggo/swag/cmd/swag@v1.16.6 init -g cmd/api/main.go -o docs/swagger

seed:
	go run ./cmd/seed

