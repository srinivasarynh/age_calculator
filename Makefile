.PHONY: help build run test clean deps sqlc \
	docker-build docker-up docker-down docker-logs \
	migrate-up migrate-down migrate-create migrate-docker

-include .env
export

APP_NAME := user-api
BIN_DIR := bin
BIN_FILE := $(BIN_DIR)/server

DB_URL_LOCAL := postgresql://postgres:postgres@localhost:5432/userdb?sslmode=disable
DB_URL_DOCKER := postgresql://postgres:postgres@postgres:5432/userdb?sslmode=disable

help: 
	@echo 'Usage: make [target]'
	@echo ''
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-20s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

build:
	go build -o $(BIN_FILE) cmd/server/main.go

run: 
	go run cmd/server/main.go

test:
	go test -v ./...

test-coverage:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out

clean: 
	rm -rf $(BIN_DIR) coverage.out

deps:
	go mod download
	go mod tidy

sqlc:
	sqlc generate

docker-build:
	docker build -t $(APP_NAME):latest .

docker-up: 
	docker compose up -d --build

docker-down:
	docker compose down

docker-logs:
	docker compose logs -f


migrate-up: 
	migrate -path db/migrations -database "$(DB_URL_LOCAL)" up

migrate-down:
	@read -p "Are you sure? This will rollback DB (y/N): " ans; \
	if [ "$$ans" = "y" ]; then \
		migrate -path db/migrations -database "$(DB_URL_LOCAL)" down; \
	else \
		echo "Aborted"; \
	fi

migrate-create:
	migrate create -ext sql -dir db/migrations -seq $(name)


migrate-docker:
	docker run --rm \
		--network container:user_api_postgres\
		-v $(PWD)/db/migrations:/migrations \
		migrate/migrate \
		-path=/migrations \
		-database "postgresql://postgres:postgres@postgres:5432/userdb?sslmode=disable" up

.DEFAULT_GOAL := help
