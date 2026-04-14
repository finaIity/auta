.PHONY: help db-up db-down build run-metadata clean test

help:
	@echo "Auta Development Commands"
	@echo "========================="
	@echo "make db-up           - Start PostgreSQL in Docker"
	@echo "make db-down         - Stop PostgreSQL"
	@echo "make db-migrate      - Run database migrations"
	@echo "make build           - Build all services"
	@echo "make run-metadata    - Run metadata service"
	@echo "make test            - Run tests"
	@echo "make clean           - Clean build artifacts"

db-up:
	docker-compose up -d postgres
	@echo "Waiting for PostgreSQL to be ready..."
	@sleep 3

db-down:
	docker-compose down

db-migrate: db-up
	@echo "Running migrations..."
	psql -h localhost -U auta -d auta -f migrations/001_initial_schema.sql

build:
	go build -o bin/metadata-service ./cmd/metadata-service

run-metadata: build
	DATABASE_URL="postgres://auta:auta_dev_password@localhost/auta?sslmode=disable" ./bin/metadata-service

test:
	go test -v ./...

clean:
	rm -rf bin/
	go clean
