# Products Module Makefile

.PHONY: help build test run docker-build docker-up docker-down clean

# Variables
APP_NAME = products-module
DOCKER_COMPOSE = docker compose
GO = go

# Default target
help:
	@echo "📦 Products Module - Available Commands"
	@echo ""
	@echo "Local Development:"
	@echo "  make build            Build the binary"
	@echo "  make run              Run the server locally"
	@echo "  make dev              Build and run locally"
	@echo "  make deps             Download Go dependencies"
	@echo "  make clean            Clean build artifacts"
	@echo ""
	@echo "Testing:"
	@echo "  make test             Run Go tests"
	@echo "  make health           Check API health"
	@echo ""
	@echo "Docker:"
	@echo "  make docker-build     Build Docker image"
	@echo "  make docker-up        Start all services"
	@echo "  make docker-down      Stop all services"
	@echo "  make docker-logs      View logs"
	@echo "  make docker-shell     Shell into container"
	@echo "  make docker-clean     Remove containers and volumes"
	@echo ""

#------------------------------------
# Development (Local)
#------------------------------------

# Build the binary
build:
	@echo "🔨 Building $(APP_NAME)..."
	$(GO) build -o $(APP_NAME) .
	@echo "✅ Build complete: ./$(APP_NAME)"

# Run Go tests
test:
	@echo "🧪 Running Go tests..."
	$(GO) test -v ./...

# Download dependencies
deps:
	@echo "📦 Downloading dependencies..."
	$(GO) mod download
	$(GO) mod tidy
	@echo "✅ Dependencies ready"

# Run locally (requires MongoDB running)
run: build
	@echo "🚀 Starting server locally..."
	./$(APP_NAME) serve --config=config.dev.yaml

# Full local development workflow
dev: build run

# Clean build artifacts
clean:
	@echo "🧹 Cleaning build artifacts..."
	rm -f $(APP_NAME)
	rm -rf vendor/

#------------------------------------
# Docker
#------------------------------------

docker-build:
	@echo "🐳 Building Docker image..."
	$(DOCKER_COMPOSE) build

docker-up:
	@echo "🐳 Starting services..."
	$(DOCKER_COMPOSE) up -d
	@echo "✅ Services started"
	@echo "   API: http://localhost:9091"

docker-down:
	@echo "🛑 Stopping services..."
	$(DOCKER_COMPOSE) down

docker-restart:
	@echo "🔄 Restarting services..."
	$(DOCKER_COMPOSE) restart

docker-logs:
	@echo "📜 Viewing logs (Ctrl+C to exit)..."
	$(DOCKER_COMPOSE) logs -f

docker-logs-api:
	@echo "📜 Viewing API logs (Ctrl+C to exit)..."
	$(DOCKER_COMPOSE) logs -f products-api

docker-shell:
	@echo "🐚 Opening shell in API container..."
	docker exec -it products-api sh

docker-clean:
	@echo "🗑️  Removing containers and volumes..."
	$(DOCKER_COMPOSE) down -v
	docker system prune -f

#------------------------------------
# Health Checks
#------------------------------------

health:
	@echo "🏥 Checking service health..."
	@curl -s http://localhost:9091/health | jq . || echo "❌ API not responding"

status:
	@echo "📊 Service Status:"
	@$(DOCKER_COMPOSE) ps

#------------------------------------
# Quick Start
#------------------------------------

quick-start: docker-up
	@echo ""
	@echo "🎉 Products Module is running!"
	@echo ""
	@echo "Next steps:"
	@echo "  1. Create an API key in auth_module"
	@echo "  2. Use the API key to create products"
	@echo ""
	@echo "View logs: make docker-logs"
