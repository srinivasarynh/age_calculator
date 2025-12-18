.PHONY: help build run test clean docker-up docker-down migrate sqlc

help:
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

build: 
	go build -o bin/server cmd/server/main.go

run: 
	go run cmd/server/main.go

test:
	go test -v ./...

test-coverage:
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out

clean: 
	rm -rf bin/
	rm -f coverage.out

deps: 
	go mod download
	go mod tidy

sqlc:
	sqlc generate

docker-build:
	docker build -t user-api:latest .

docker-up: 
	docker-compose up -d

docker-down:
	docker-compose down

docker-logs:
	docker-compose logs -f

migrate-up:
	migrate -path db/migrations -database "postgresql://postgres:postgres@localhost:5432/userdb?sslmode=disable" up

migrate-down: 
	migrate -path db/migrations -database "postgresql://postgres:postgres@localhost:5432/userdb?sslmode=disable" down

migrate-create:
	migrate create -ext sql -dir db/migrations -seq $(name)

.DEFAULT_GOAL := help
