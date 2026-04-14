.PHONY: help db-up db-down db-test-up db-test-down build run-metadata clean test test-manual test-unit test-integration

help:
	@echo "Auta Development Commands"
	@echo "========================="
	@echo ""
	@echo "Database:"
	@echo "  make db-up              - Start PostgreSQL in Docker"
	@echo "  make db-down            - Stop PostgreSQL"
	@echo "  make db-migrate         - Run database migrations"
	@echo "  make db-test-up         - Start test PostgreSQL database"
	@echo "  make db-test-down       - Stop test PostgreSQL database"
	@echo ""
	@echo "Building:"
	@echo "  make build              - Build all services"
	@echo "  make run-metadata       - Run metadata service (requires db-up)"
	@echo ""
	@echo "Testing:"
	@echo "  make test               - Run all tests"
	@echo "  make test-manual        - Show manual testing guide"
	@echo "  make test-unit          - Run unit tests only"
	@echo "  make test-integration   - Run integration tests (requires db-test-up)"
	@echo ""
	@echo "Other:"
	@echo "  make clean              - Clean build artifacts"

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

test-unit:
	go test -v -run "Test[A-Z]" ./internal/metadata

test-integration: db-test-up
	go test -v -run "TestIntegration" ./...

test-manual:
	@echo "=== Manual Testing Guide ==="
	@echo ""
	@echo "1. Start the service:"
	@echo "   make db-up"
	@echo "   make db-migrate"
	@echo "   make run-metadata"
	@echo ""
	@echo "2. In another terminal, run:"
	@echo "   bash scripts/manual-test.sh"
	@echo ""
	@echo "Or manually test individual endpoints:"
	@echo "   curl http://localhost:8000/health"

db-test-up:
	@echo "Setting up test database..."
	@docker-compose -f docker-compose.test.yml up -d postgres-test 2>/dev/null || \
	  (echo "Creating test database on localhost:5433..." && \
	   docker run -d --name auta-test-db \
	     -e POSTGRES_DB=auta_test \
	     -e POSTGRES_PASSWORD=auta_test_password \
	     -e POSTGRES_USER=auta \
	     -p 5433:5432 \
	     postgres:16-alpine 2>/dev/null || true)
	@sleep 3

db-test-down:
	@docker stop auta-test-db 2>/dev/null || true
	@docker rm auta-test-db 2>/dev/null || true
	rm -rf bin/
	go clean
