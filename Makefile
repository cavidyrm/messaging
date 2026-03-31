# Makefile

migrate:
	@echo "Running migrations..."
	@go run cmd/migrate/main.go

up:
	docker-compose up -d

down:
	docker-compose down -v

build:
	docker-compose build

logs:
	docker-compose logs -f app

.PHONY: migrate up down build logs
