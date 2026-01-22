.PHONY: help dev build up down logs clean install

help:
	@echo "Aether Panel - Makefile Commands"
	@echo ""
	@echo "  make dev       - Start development environment"
	@echo "  make build     - Build Docker images"
	@echo "  make up        - Start all services"
	@echo "  make down      - Stop all services"
	@echo "  make logs      - View logs"
	@echo "  make clean     - Remove all containers and volumes"
	@echo "  make install   - Run installation script"

dev:
	docker compose up -d postgres redis
	cd backend && go run ./cmd/api &
	cd frontend && npm run dev

build:
	docker compose build

up:
	docker compose up -d

down:
	docker compose down

logs:
	docker compose logs -f

clean:
	docker compose down -v --remove-orphans
	docker system prune -f

install:
	chmod +x install.sh
	sudo ./install.sh

# Backend commands
backend-build:
	cd backend && go build -o bin/aether-api ./cmd/api

backend-test:
	cd backend && go test ./...

backend-lint:
	cd backend && golangci-lint run

# Frontend commands
frontend-install:
	cd frontend && npm install

frontend-build:
	cd frontend && npm run build

frontend-dev:
	cd frontend && npm run dev

# Database commands
db-migrate:
	docker exec aether_api /app/aether-api migrate

db-seed:
	docker exec aether_api /app/aether-api seed

db-backup:
	docker exec aether_postgres pg_dump -U aether aether_panel > backup_$$(date +%Y%m%d_%H%M%S).sql

db-restore:
	docker exec -i aether_postgres psql -U aether aether_panel < $(FILE)
