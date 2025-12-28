# Auth Module Makefile

.PHONY: help build test run docker-build docker-up docker-down docker-logs clean dev domain-create domain-list

# Variables
APP_NAME = auth-module
DOCKER_COMPOSE = docker-compose
GO = go

# Default target
help:
	@echo "📦 Auth Module - Available Commands"
	@echo ""
	@echo "Development:"
	@echo "  make dev              Run server locally (requires MongoDB running)"
	@echo "  make build            Build the binary"
	@echo "  make test             Run Go tests"
	@echo "  make clean            Clean build artifacts"
	@echo ""
	@echo "Testing:"
	@echo "  make test-integration Run automated integration tests"
	@echo "  make test-manual      View manual testing guide"
	@echo "  make health           Check API health"
	@echo "  make status           Show service status"
	@echo ""
	@echo "Docker:"
	@echo "  make docker-build     Build Docker image"
	@echo "  make docker-up        Start all services (MongoDB + API)"
	@echo "  make docker-down      Stop all services"
	@echo "  make docker-restart   Restart all services"
	@echo "  make docker-logs      View logs"
	@echo "  make docker-shell     Shell into API container"
	@echo "  make docker-clean     Remove containers and volumes"
	@echo ""
	@echo "Docker with Debug Tools:"
	@echo "  make docker-up-debug  Start with Mongo Express (localhost:8081)"
	@echo ""
	@echo "Domain Management (Docker):"
	@echo "  make domain-create DOMAIN=example.com NAME=\"Example\" EMAIL=admin@example.com"
	@echo "  make domain-list      List all domains"
	@echo ""
	@echo "Setup:"
	@echo "  make setup            Initial setup (copy configs)"
	@echo "  make deps             Download Go dependencies"

#------------------------------------
# Development (Local)
#------------------------------------

build:
	@echo "🔨 Building $(APP_NAME)..."
	$(GO) build -o $(APP_NAME) .
	@echo "✅ Build complete"

test:
	@echo "🧪 Running tests..."
	$(GO) test -v ./...

dev: build
	@echo "🚀 Starting development server..."
	./$(APP_NAME) serve

deps:
	@echo "📦 Downloading dependencies..."
	$(GO) mod download
	$(GO) mod tidy

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
	@echo "   API: http://localhost:8080"
	@echo "   MongoDB: localhost:27017"
	@make docker-wait-healthy

docker-up-debug:
	@echo "🐳 Starting services with debug tools..."
	$(DOCKER_COMPOSE) --profile debug up -d
	@echo "✅ Services started"
	@echo "   API: http://localhost:8080"
	@echo "   MongoDB: localhost:27017"
	@echo "   Mongo Express: http://localhost:8081 (admin/admin)"
	@make docker-wait-healthy

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
	$(DOCKER_COMPOSE) logs -f auth-module

docker-logs-mongo:
	@echo "📜 Viewing MongoDB logs (Ctrl+C to exit)..."
	$(DOCKER_COMPOSE) logs -f mongodb

docker-shell:
	@echo "🐚 Opening shell in API container..."
	docker exec -it auth-module-api sh

docker-mongo-shell:
	@echo "🐚 Opening MongoDB shell..."
	docker exec -it auth-module-mongodb mongosh auth_module

docker-clean:
	@echo "🗑️  Removing containers and volumes..."
	$(DOCKER_COMPOSE) down -v
	docker system prune -f

docker-wait-healthy:
	@echo "⏳ Waiting for services to be healthy..."
	@timeout 60 sh -c 'until docker exec auth-module-mongodb mongosh --quiet --eval "db.adminCommand(\"ping\")" > /dev/null 2>&1; do sleep 2; done' && echo "✅ MongoDB is ready"
	@timeout 60 sh -c 'until curl -sf http://localhost:8080/health > /dev/null 2>&1; do sleep 2; done' && echo "✅ API is ready"

#------------------------------------
# Domain Management (via Docker)
#------------------------------------

domain-create:
	@if [ -z "$(DOMAIN)" ] || [ -z "$(NAME)" ] || [ -z "$(EMAIL)" ]; then \
		echo "❌ Error: DOMAIN, NAME, and EMAIL are required"; \
		echo "Usage: make domain-create DOMAIN=example.com NAME=\"Example\" EMAIL=admin@example.com"; \
		exit 1; \
	fi
	@echo "🌐 Creating domain: $(DOMAIN)"
	docker exec -it auth-module-api ./auth-module domain create \
		--domain=$(DOMAIN) \
		--name="$(NAME)" \
		--admin-email=$(EMAIL)

domain-list:
	@echo "📋 Listing domains..."
	docker exec -it auth-module-api ./auth-module domain list

domain-delete:
	@if [ -z "$(DOMAIN)" ]; then \
		echo "❌ Error: DOMAIN is required"; \
		echo "Usage: make domain-delete DOMAIN=example.com"; \
		exit 1; \
	fi
	@echo "⚠️  Deleting domain: $(DOMAIN)"
	docker exec -it auth-module-api ./auth-module domain delete --domain=$(DOMAIN)

#------------------------------------
# Setup
#------------------------------------

setup:
	@echo "⚙️  Setting up project..."
	@if [ ! -f .env ]; then \
		cp .env.example .env; \
		echo "✅ Created .env file (please edit with your settings)"; \
	else \
		echo "ℹ️  .env already exists"; \
	fi
	@if [ ! -f config.yaml ]; then \
		cp config.yaml.example config.yaml; \
		echo "✅ Created config.yaml file"; \
	else \
		echo "ℹ️  config.yaml already exists"; \
	fi
	@echo "✅ Setup complete!"
	@echo ""
	@echo "Next steps:"
	@echo "  1. Edit .env with your settings (especially JWT_SECRET)"
	@echo "  2. Run: make docker-up"
	@echo "  3. Create a domain: make domain-create DOMAIN=example.com NAME=\"Example\" EMAIL=admin@example.com"

#------------------------------------
# Testing
#------------------------------------

test-integration:
	@echo "🧪 Running integration tests..."
	@./test.sh

test-manual:
	@echo "📖 Opening manual testing guide..."
	@cat TESTING.md

#------------------------------------
# Health Checks
#------------------------------------

health:
	@echo "🏥 Checking service health..."
	@curl -s http://localhost:8080/health | jq . || echo "❌ API not responding"

status:
	@echo "📊 Service Status:"
	@$(DOCKER_COMPOSE) ps

#------------------------------------
# Quick Start
#------------------------------------

quick-start: setup docker-up
	@echo ""
	@echo "🎉 Auth Module is running!"
	@echo ""
	@echo "Create your first domain:"
	@echo "  make domain-create DOMAIN=localhost NAME=\"Local Dev\" EMAIL=admin@localhost"
	@echo ""
	@echo "View API docs: cat API.md"
	@echo "View logs: make docker-logs"
