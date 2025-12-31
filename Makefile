.PHONY: help build start stop restart logs clean deploy dev-db status test-website
.PHONY: create-invite list-invites list-users create-domain list-domains delete-invite delete-user

# Default values
DOMAIN ?= oilyourhair.com
EMAIL ?=
NAME ?=
ROLE ?= customer
PERMISSIONS ?= products.read,orders.read

# Default target
help:
	@echo "OilYourHair App - Docker Management"
	@echo ""
	@echo "Service Management:"
	@echo "  make build                  - Build all Docker images"
	@echo "  make start                  - Start all services"
	@echo "  make stop                   - Stop all services"
	@echo "  make restart                - Restart all services"
	@echo "  make logs                   - View logs from all services"
	@echo "  make status                 - Show status of all containers"
	@echo "  make clean                  - Stop and remove all containers, networks, and volumes"
	@echo "  make deploy                 - Build and deploy all services"
	@echo ""
	@echo "Database Management:"
	@echo "  make dev-db                 - Access MongoDB shell"
	@echo ""
	@echo "Local Testing:"
	@echo "  make test-website           - Serve test website (products + auth) on :3000"
	@echo ""
	@echo "Invitation Management:"
	@echo "  make create-invite EMAIL=user@example.com [NAME=John] [DOMAIN=oilyourhair.com] [ROLE=customer]"
	@echo "                              - Create invitation for a user"
	@echo "  make list-invites [DOMAIN=oilyourhair.com]"
	@echo "                              - List all invitations for domain"
	@echo "  make delete-invite TOKEN=xxx"
	@echo "                              - Delete an invitation"
	@echo ""
	@echo "User Management:"
	@echo "  make list-users [DOMAIN=oilyourhair.com]"
	@echo "                              - List all users for domain"
	@echo "  make delete-user EMAIL=user@example.com"
	@echo "                              - Delete a user"
	@echo ""
	@echo "Domain Management:"
	@echo "  make create-domain DOMAIN=example.com NAME=CompanyName"
	@echo "                              - Create a new domain"
	@echo "  make list-domains           - List all domains"
	@echo ""

# Build all Docker images
build:
	@echo "Building all Docker images..."
	docker compose build

# Start all services
start:
	@echo "Starting all services..."
	docker compose up -d
	@echo ""
	@echo "Services started:"
	@echo "  - MongoDB:              Internal (auth-network)"
	@echo "  - Auth API:             http://localhost:9090"
	@echo "  - Auth UI:              http://localhost:5173"
	@echo "  - Test Domain:          http://localhost:8000"
	@echo "  - OilYourHair Frontend: http://localhost:8080"
	@echo ""

# Stop all services
stop:
	@echo "Stopping all services..."
	docker compose down

# Restart all services
restart: stop start

# View logs from all services
logs:
	docker compose logs -f

# View logs for specific service
logs-api:
	docker compose logs -f auth-api

logs-ui:
	docker compose logs -f auth-ui

logs-frontend:
	docker compose logs -f oilyourhair-frontend

logs-db:
	docker compose logs -f mongodb

# Show status of all containers
status:
	@echo "Container Status:"
	@docker compose ps
	@echo ""
	@echo "Health Checks:"
	@docker ps --filter "name=auth" --format "table {{.Names}}\t{{.Status}}"

# Clean up everything (containers, networks, volumes)
clean:
	@echo "WARNING: This will remove all containers, networks, and volumes!"
	@read -p "Are you sure? [y/N] " -n 1 -r; \
	echo; \
	if [[ $$REPLY =~ ^[Yy]$$ ]]; then \
		docker compose down -v; \
		echo "Cleanup complete."; \
	else \
		echo "Cleanup cancelled."; \
	fi

# Build and deploy all services
deploy: build start
	@echo "Deployment complete!"

# Access MongoDB shell
dev-db:
	docker exec -it auth-mongodb mongosh auth_module

# Rebuild specific service
rebuild-api:
	docker compose build auth-api
	docker compose up -d auth-api

rebuild-ui:
	docker compose build auth-ui
	docker compose up -d auth-ui

rebuild-frontend:
	docker compose build oilyourhair-frontend
	docker compose up -d oilyourhair-frontend

# Quick development cycle
dev: stop build start logs

# ============================================
# Invitation Management
# ============================================

# Create an invitation for a user
create-invite:
ifndef EMAIL
	@echo "Error: EMAIL is required"
	@echo "Usage: make create-invite EMAIL=user@example.com [NAME=John] [DOMAIN=oilyourhair.com] [ROLE=customer]"
	@exit 1
endif
	@echo "Creating invitation for $(EMAIL)..."
	@TIMESTAMP=$$(date +%s); \
	TOKEN="invite-$$TIMESTAMP"; \
	PERMS_ARRAY=$$(echo "$(PERMISSIONS)" | sed "s/,/', '/g"); \
	docker exec auth-mongodb mongosh auth_module --quiet --eval " \
		db.invitations.insertOne({ \
			domain: '$(DOMAIN)', \
			token: '$$TOKEN', \
			email: '$(EMAIL)', \
			role: '$(ROLE)', \
			permissions: ['$$PERMS_ARRAY'], \
			type: 'email_with_qr', \
			single_use: true, \
			uses_count: 0, \
			created_by: 'admin', \
			created_at: new Date(), \
			expires_at: new Date(Date.now() + 7*24*60*60*1000), \
			status: 'pending' \
		})" > /dev/null; \
	echo ""; \
	echo "✓ Invitation created successfully!"; \
	echo ""; \
	echo "  Email:  $(EMAIL)"; \
	echo "  Domain: $(DOMAIN)"; \
	echo "  Role:   $(ROLE)"; \
	echo "  Token:  $$TOKEN"; \
	echo ""; \
	echo "Invitation URL (Development):"; \
	echo "  http://localhost:5173/invite?token=$$TOKEN&redirect=http://localhost:8080"; \
	echo ""; \
	echo "Invitation URL (Production):"; \
	echo "  https://auth.$(DOMAIN)/invite?token=$$TOKEN&redirect=https://$(DOMAIN)"; \
	echo ""

# List all invitations for a domain
list-invites:
	@echo "Invitations for $(DOMAIN):"
	@echo ""
	@docker exec auth-mongodb mongosh auth_module --quiet --eval " \
		db.invitations.find({domain: '$(DOMAIN)'}).forEach(function(inv) { \
			print('Token: ' + inv.token); \
			print('Email: ' + inv.email); \
			print('Role: ' + inv.role); \
			print('Status: ' + inv.status); \
			print('Created: ' + inv.created_at); \
			print('Expires: ' + inv.expires_at); \
			print('Uses: ' + inv.uses_count); \
			print('---'); \
		})"

# Delete an invitation
delete-invite:
ifndef TOKEN
	@echo "Error: TOKEN is required"
	@echo "Usage: make delete-invite TOKEN=invite-123456"
	@exit 1
endif
	@echo "Deleting invitation with token: $(TOKEN)..."
	@docker exec auth-mongodb mongosh auth_module --quiet --eval " \
		db.invitations.deleteOne({token: '$(TOKEN)'})" > /dev/null
	@echo "✓ Invitation deleted"

# ============================================
# User Management
# ============================================

# List all users for a domain
list-users:
	@echo "Users for $(DOMAIN):"
	@echo ""
	@docker exec auth-mongodb mongosh auth_module --quiet --eval " \
		db.users.find({domain: '$(DOMAIN)'}).forEach(function(user) { \
			print('Email: ' + user.email); \
			print('Role: ' + user.role); \
			print('Permissions: ' + user.permissions.join(', ')); \
			print('Auth Provider: ' + user.auth_provider); \
			print('Created: ' + user.created_at); \
			print('---'); \
		})"

# Delete a user
delete-user:
ifndef EMAIL
	@echo "Error: EMAIL is required"
	@echo "Usage: make delete-user EMAIL=user@example.com"
	@exit 1
endif
	@echo "Deleting user: $(EMAIL)..."
	@docker exec auth-mongodb mongosh auth_module --quiet --eval " \
		db.users.deleteOne({email: '$(EMAIL)'})" > /dev/null
	@echo "✓ User deleted"

# ============================================
# Domain Management
# ============================================

# Create a new domain
create-domain:
ifndef DOMAIN
	@echo "Error: DOMAIN is required"
	@echo "Usage: make create-domain DOMAIN=example.com NAME=CompanyName"
	@exit 1
endif
ifndef NAME
	@echo "Error: NAME is required"
	@echo "Usage: make create-domain DOMAIN=example.com NAME=CompanyName"
	@exit 1
endif
	@echo "Creating domain: $(DOMAIN)..."
	@docker exec auth-api ./auth-module domain create \
		--domain $(DOMAIN) \
		--name "$(NAME)"
	@echo "✓ Domain created successfully!"

# List all domains
list-domains:
	@echo "All domains:"
	@echo ""
	@docker exec auth-api ./auth-module domain list

# ============================================
# Local Testing (Product Management + Auth)
# ============================================

# Serve the generic test website for product management and auth testing
test-website:
	@echo "🌐 Starting test website for product management + auth..."
	@echo ""
	@echo "This website tests both:"
	@echo "  - Product catalog and management (products_module)"
	@echo "  - Authentication and login (auth_module)"
	@echo ""
	@echo "Website available at:"
	@echo "  http://localhost:3000/index.html      - Homepage (bestselling products)"
	@echo "  http://localhost:3000/shop.html       - Product catalog with search/filters"
	@echo "  http://localhost:3000/admin/products.html - Admin product management"
	@echo ""
	@echo "Prerequisites:"
	@echo "  - Products module running on :9091"
	@echo "  - Auth module running on :9090"
	@echo ""
	@echo "Press Ctrl+C to stop the server"
	@echo ""
	@cd oilyourhair.com/public && python3 -m http.server 3000
