.PHONY: help build run test clean deps sqlc \
	docker-build docker-up docker-down docker-logs \
	migrate-up migrate-down migrate-create migrate-docker

# Load env vars
-include .env
export

APP_NAME := user-api
BIN_DIR := bin
BIN_FILE := $(BIN_DIR)/server

DB_URL_LOCAL := postgresql://postgres:postgres@localhost:5432/userdb?sslmode=disable
DB_URL_DOCKER := postgresql://postgres:postgres@postgres:5432/userdb?sslmode=disable

help: ## Show help
	@echo 'Usage: make [target]'
	@echo ''
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-20s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

build: ## Build binary
	go build -o $(BIN_FILE) cmd/server/main.go

run: ## Run locally
	go run cmd/server/main.go

test: ## Run tests
	go test -v ./...

test-coverage: ## Run tests with coverage
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out

clean: ## Clean artifacts
	rm -rf $(BIN_DIR) coverage.out

deps: ## Download dependencies
	go mod download
	go mod tidy

sqlc: ## Generate sqlc code
	sqlc generate

docker-build: ## Build docker image
	docker build -t $(APP_NAME):latest .

docker-up: ## Start docker stack
	docker compose up -d --build

docker-down: ## Stop docker stack
	docker compose down

docker-logs: ## Follow docker logs
	docker compose logs -f

# -------------------------
# Migrations (local)
# -------------------------

migrate-up: ## Run migrations locally
	migrate -path db/migrations -database "$(DB_URL_LOCAL)" up

migrate-down: ## Rollback migrations locally (CONFIRM!)
	@read -p "Are you sure? This will rollback DB (y/N): " ans; \
	if [ "$$ans" = "y" ]; then \
		migrate -path db/migrations -database "$(DB_URL_LOCAL)" down; \
	else \
		echo "Aborted"; \
	fi

migrate-create: ## Create new migration (name=xxx)
	migrate create -ext sql -dir db/migrations -seq $(name)

# -------------------------
# Migrations (Docker)
# -------------------------

migrate-docker:
	docker run --rm \
		--network user-api_default \
		-v $(PWD)/db/migrations:/migrations \
		migrate/migrate \
		-path=/migrations \
		-database "postgresql://postgres:postgres@postgres:5432/userdb?sslmode=disable" up

.DEFAULT_GOAL := help
