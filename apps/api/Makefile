.PHONY: run build test clean help dev dev-init air-init

# Default target
.DEFAULT_GOAL := help

# Run the application
run:
	@echo "Starting application..."
	go run cmd/api/main.go

# Build the application
build:
	@echo "Building application..."
	go build -o bin/api cmd/api/main.go
	@echo "Binary created at: bin/api"

# Run tests
test:
	@echo "Running tests..."
	go test -v ./...

# Run tests with coverage
test-coverage:
	@echo "Running tests with coverage..."
	go test -cover -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Install dependencies
deps:
	@echo "Installing dependencies..."
	go mod download
	go mod tidy

# Format code
fmt:
	@echo "Formatting code..."
	go fmt ./...

# Lint code
lint:
	@echo "Linting code..."
	golangci-lint run

# Clean build artifacts
clean:
	@echo "Cleaning..."
	rm -rf bin/
	rm -rf tmp/
	rm -f coverage.out coverage.html
	rm -f build-errors.log
	@echo "Clean complete"

# Database migrations
migrate-up:
	@echo "Running migrations..."
	go run cmd/migrate/main.go -action=up

migrate-down:
	@echo "Rolling back migrations..."
	go run cmd/migrate/main.go -action=down

migrate-fresh:
	@echo "Running fresh migration (drop all + migrate)..."
	go run cmd/migrate/main.go -action=fresh

migrate-seed:
	@echo "Seeding database with default data..."
	@echo "  • Admin role"
	@echo "  • User role"
	@echo "  • Admin user: admin@fahmi.com / Test@1234"
	go run cmd/migrate/main.go -action=seed

migrate-status:
	@echo "Checking migration status..."
	go run cmd/migrate/main.go -action=status

# Hot reload for development (requires air)
# Note: Migrations are now run automatically on server startup
dev:
	@echo "Starting with hot reload (Air)..."
	@if [ ! -f ~/go/bin/air ]; then \
		echo "Air not installed. Installing..."; \
		go install github.com/air-verse/air@latest; \
		echo "✓ Air installed successfully"; \
	fi
	@~/go/bin/air

# Development without auto-migrations (run migrations manually first)
dev-no-migrate:
	@echo "Starting without auto-migrations..."
	@echo "Note: Make sure to run 'make migrate-up' first if needed"
	@~/go/bin/air

# Initialize Air with default config (run once)
air-init:
	@echo "Initializing Air configuration..."
	@if [ ! -f .air.toml ]; then \
		air init; \
	else \
		echo ".air.toml already exists"; \
	fi

# Development with fresh start (clean + migrate + dev)
dev-init: clean migrate-fresh
	@echo "Starting fresh development environment..."
	@make dev

# Install development tools
install-tools:
	@echo "Installing development tools..."
	go install github.com/air-verse/air@latest
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install github.com/swaggo/swag/cmd/swag@latest
	@echo "Tools installed successfully"

# Generate Swagger documentation
swagger:
	@echo "Generating Swagger documentation..."
	~/go/bin/swag init -g cmd/api/main.go -o docs
	@echo "Swagger docs generated successfully"

# Docker commands
docker-build:
	@echo "Building Docker image..."
	docker-compose build

docker-up:
	@echo "Starting services..."
	docker-compose up -d

docker-down:
	@echo "Stopping services..."
	docker-compose down

docker-restart:
	@echo "Restarting services..."
	docker-compose restart

docker-logs:
	docker-compose logs -f

docker-logs-api:
	docker-compose logs -f api

docker-logs-nginx:
	docker-compose logs -f nginx

docker-clean:
	@echo "Cleaning Docker resources..."
	docker-compose down -v
	docker system prune -f

docker-rebuild:
	@echo "Rebuilding and restarting..."
	docker-compose down
	docker-compose build --no-cache
	docker-compose up -d

# Database operations (Docker)
docker-migrate-up:
	@echo "Running migrations in container..."
	docker-compose exec api ./api migrate up

docker-shell-api:
	docker-compose exec api sh

docker-shell-db:
	docker-compose exec postgres psql -U bloguser -d blogdb

docker-db-backup:
	@echo "Backing up database..."
	docker-compose exec -T postgres pg_dump -U bloguser blogdb > backup_$$(date +%Y%m%d_%H%M%S).sql

health-check:
	@echo "Checking service health..."
	@curl -f http://localhost/health && echo "✓ Nginx/API: UP" || echo "✗ Nginx/API: DOWN"
	@curl -f http://localhost:8080/health && echo "✓ API Direct: UP" || echo "✗ API Direct: DOWN"

# Local development with Docker dependencies only
local-dev:
	@echo "Starting local development..."
	docker-compose up -d postgres mailhog
	@echo "Waiting for PostgreSQL to be ready..."
	@sleep 5
	@echo "Starting API locally..."
	go run cmd/api/main.go

# Help command
help:
	@echo "Available commands:"
	@echo ""
	@echo "Application:"
	@echo "  make run             - Run the application"
	@echo "  make build           - Build the application binary"
	@echo "  make dev             - Run with hot reload (Air - auto-install if needed)"
	@echo "  make dev-init        - Clean + migrate + run with hot reload"
	@echo "  make local-dev       - Run API locally with Docker deps (PostgreSQL + MailHog)"
	@echo ""
	@echo "Database Migrations (Local):"
	@echo "  make migrate-up      - Run database migrations"
	@echo "  make migrate-down    - Rollback all migrations"
	@echo "  make migrate-fresh   - Drop all tables and re-migrate"
	@echo "  make migrate-seed    - Seed database with default roles and admin user (admin@fahmi.com / Test@1234)"
	@echo "  make migrate-status  - Check migration status"
	@echo ""
	@echo "Testing:"
	@echo "  make test            - Run tests"
	@echo "  make test-coverage   - Run tests with coverage report"
	@echo ""
	@echo "Development:"
	@echo "  make deps            - Install dependencies"
	@echo "  make fmt             - Format code"
	@echo "  make lint            - Lint code"
	@echo "  make clean           - Clean build artifacts and temp files"
	@echo "  make install-tools   - Install development tools (air, golangci-lint, swag)"
	@echo "  make swagger         - Generate Swagger API documentation"
	@echo "  make air-init        - Initialize Air configuration file"
	@echo ""
	@echo "Docker Compose:"
	@echo "  make docker-build    - Build Docker images"
	@echo "  make docker-up       - Start all services (API, PostgreSQL, Nginx, MailHog)"
	@echo "  make docker-down     - Stop all services"
	@echo "  make docker-restart  - Restart all services"
	@echo "  make docker-rebuild  - Rebuild and restart (no cache)"
	@echo "  make docker-clean    - Remove containers, volumes, and prune system"
	@echo "  make docker-logs     - View logs from all services"
	@echo "  make docker-logs-api - View API logs only"
	@echo "  make docker-logs-nginx - View Nginx logs only"
	@echo ""
	@echo "Docker Database:"
	@echo "  make docker-migrate-up - Run migrations in container"
	@echo "  make docker-shell-api  - Open shell in API container"
	@echo "  make docker-shell-db   - Open PostgreSQL shell"
	@echo "  make docker-db-backup  - Backup database to SQL file"
	@echo ""
	@echo "Health & Monitoring:"
	@echo "  make health-check    - Check health of all services"
	@echo ""
	@echo "  make help            - Show this help message"
