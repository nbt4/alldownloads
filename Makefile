.PHONY: help build up down logs clean lint test migrate dev prod

# Default target
help:
	@echo "Available targets:"
	@echo "  build     - Build all Docker images"
	@echo "  up        - Start all services"
	@echo "  down      - Stop all services"
	@echo "  logs      - Show logs from all services"
	@echo "  clean     - Remove all containers, images, and volumes"
	@echo "  lint      - Run linters for Go and Node.js"
	@echo "  test      - Run tests"
	@echo "  migrate   - Run database migrations"
	@echo "  dev       - Start development environment"
	@echo "  prod      - Start production environment with MinIO"

# Build all Docker images
build:
	@echo "Building Docker images..."
	docker-compose build --parallel

# Start all services
up:
	@echo "Starting AllDownloads stack..."
	docker-compose up -d
	@echo "Services started. Access the application at http://localhost:3000"
	@echo "API is available at http://localhost:8080"

# Stop all services
down:
	@echo "Stopping AllDownloads stack..."
	docker-compose down

# Show logs
logs:
	docker-compose logs -f

# Clean up everything
clean:
	@echo "Cleaning up..."
	docker-compose down -v --remove-orphans
	docker system prune -f
	docker volume prune -f

# Run linters
lint:
	@echo "Running Go linters..."
	golangci-lint run ./...
	@echo "Running Node.js linters..."
	cd ui && npm run lint

# Run tests
test:
	@echo "Running Go tests..."
	go test -v ./...
	@echo "Running Node.js tests..."
	cd ui && npm test

# Run database migrations
migrate:
	@echo "Running database migrations..."
	docker-compose exec api migrate -path ./migrations -database "$$DB_URL" up

# Development environment
dev:
	@echo "Starting development environment..."
	docker-compose -f docker-compose.yml up -d db cache
	@echo "Database and cache started. Run 'make migrate' to set up the database."
	@echo "Then start the API and worker locally with 'go run cmd/api/main.go' and 'go run cmd/worker/main.go'"
	@echo "Start the UI with 'cd ui && npm run dev'"

# Production environment with storage
prod:
	@echo "Starting production environment with MinIO..."
	docker-compose --profile storage up -d
	@echo "Full stack started with S3-compatible storage."
	@echo "MinIO console: http://localhost:9001"
	@echo "Application: http://localhost:3000"

# Install development dependencies
deps:
	@echo "Installing Go dependencies..."
	go mod download
	@echo "Installing Node.js dependencies..."
	cd ui && npm install

# Format code
fmt:
	@echo "Formatting Go code..."
	go fmt ./...
	@echo "Formatting Node.js code..."
	cd ui && npm run format

# Generate Go modules
mod:
	go mod tidy
	go mod vendor

# Build binaries locally
build-local:
	@echo "Building API binary..."
	go build -o bin/api ./cmd/api
	@echo "Building worker binary..."
	go build -o bin/worker ./cmd/worker

# Run security scans
security:
	@echo "Running security scans..."
	gosec ./...
	cd ui && npm audit

# Health check
health:
	@echo "Checking service health..."
	curl -f http://localhost:8080/api/health || echo "API not responding"
	curl -f http://localhost:3000 || echo "UI not responding"

# Backup database
backup:
	@echo "Creating database backup..."
	docker-compose exec db pg_dump -U alldl alldownloads > backup_$(shell date +%Y%m%d_%H%M%S).sql

# Restore database
restore:
	@echo "Restoring database from backup..."
	@read -p "Enter backup file path: " backup_file; \
	docker-compose exec -T db psql -U alldl alldownloads < $$backup_file