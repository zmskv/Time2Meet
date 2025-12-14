.PHONY: run build test swag

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

