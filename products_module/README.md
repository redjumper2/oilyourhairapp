# Products Module

Multi-tenant SaaS product catalog service with flexible key/value attributes and variant support.

## Overview

The Products Module provides a comprehensive product management API for multi-tenant e-commerce applications. Each domain (customer) maintains their own isolated product catalog with support for:

- **Flexible Product Attributes**: Key/value metadata (category, type, etc.)
- **Price Variants**: Different prices based on attribute combinations
- **Inventory Tracking**: Stock management per variant
- **Multi-tenant Isolation**: Domain-based product separation
- **Public API**: No-auth product listing for customer websites
- **Admin API**: JWT-authenticated CRUD operations

## Quick Start

### Prerequisites

- Go 1.21+
- MongoDB (shared with auth_module)
- Docker & Docker Compose

### Installation

```bash
# Install dependencies
go mod download

# Build the binary
go build -o products-module ./cmd

# Or use Makefile
make build
```

### Running the Service

```bash
# Run directly
./products-module serve

# Or with Docker Compose
make docker-up

# Development mode with hot reload
make dev
```

## Project Structure

```
products_module/
├── cmd/
│   └── main.go              # Cobra CLI entry point
├── internal/
│   ├── models/
│   │   └── product.go       # Product, Variant models
│   ├── handlers/
│   │   ├── admin.go         # CRUD (JWT required)
│   │   └── public.go        # Read-only (no auth)
│   ├── services/
│   │   └── product.go       # Business logic
│   ├── database/
│   │   └── mongo.go         # DB connection
│   └── middleware/
│       └── jwt.go           # JWT validation (shared with auth)
├── config/
│   ├── config.go            # Viper configuration
│   └── config.yaml          # Config file
├── pkg/
│   └── utils/               # Shared utilities
├── Makefile                 # Development commands
├── Dockerfile               # Container image
├── docker-compose.yml       # Service orchestration
├── README.md                # This file
└── architecture.md          # Detailed architecture docs
```

## CLI Commands

```bash
# Start API server
./products-module serve

# Create a product (interactive)
./products-module product create

# List products for a domain
./products-module product list --domain=oilyourhair.com

# Seed sample data
./products-module seed --domain=oilyourhair.com
```

## Makefile Commands

```bash
make help           # Show all available commands
make build          # Build the binary
make dev            # Run with hot reload
make test           # Run tests
make docker-up      # Start with docker-compose
make docker-down    # Stop docker services
make seed           # Seed sample products
make create-product # Interactive product creation
make list-products  # List all products
```

## API Endpoints

### Admin API (JWT Required)

```
POST   /api/v1/products          Create product
GET    /api/v1/products          List products
GET    /api/v1/products/:id      Get product by ID
PUT    /api/v1/products/:id      Update product
DELETE /api/v1/products/:id      Delete product
```

### Public API (No Auth)

```
GET    /api/v1/public/:domain/products          List products for domain
GET    /api/v1/public/:domain/products/:id      Get product details
GET    /api/v1/public/:domain/products/search   Search products
```

## Configuration

Create a `.env` file or set environment variables:

```bash
# Server
PRODUCTS_SERVER_PORT=9091
PRODUCTS_SERVER_ENV=development

# MongoDB
PRODUCTS_MONGODB_URI=mongodb://mongodb:27017
PRODUCTS_MONGODB_DATABASE=products_module

# JWT (shared with auth_module)
PRODUCTS_JWT_SECRET=your-secret-key

# CORS
PRODUCTS_CORS_ALLOWED_ORIGINS=http://localhost:3000,https://oilyourhair.com
```

## Integration with Auth Module

Products module integrates with auth_module for:

1. **Domain Validation**: Reads from `auth_module.domains` (read-only)
2. **JWT Authentication**: Validates JWTs issued by auth_module (shared secret)
3. **Role-based Access**: Checks admin role from JWT claims

See [architecture.md](./architecture.md) for detailed integration patterns.

## Example: Fetching Products on Website

```html
<!-- On oilyourhair.com -->
<div id="product-list"></div>

<script>
  fetch('https://api.oilyourhair.com/api/v1/public/oilyourhair.com/products?category=oils')
    .then(response => response.json())
    .then(data => {
      const products = data.products;
      // Render products
      document.getElementById('product-list').innerHTML = products
        .map(p => `<div>${p.name} - $${p.base_price}</div>`)
        .join('');
    });
</script>
```

## Product Schema

```json
{
  "id": "uuid",
  "domain": "oilyourhair.com",
  "name": "Coconut Oil",
  "description": "Pure organic coconut oil",
  "base_price": 29.99,
  "images": ["https://cdn.example.com/coconut-oil.jpg"],
  "attributes": {
    "category": "oils",
    "type": "coconut",
    "organic": "true"
  },
  "variants": [
    {
      "id": "var_1",
      "attributes": {
        "size": "500ml"
      },
      "price": 29.99,
      "stock": 50,
      "sku": "CO-500"
    },
    {
      "id": "var_2",
      "attributes": {
        "size": "1L",
        "organic": "true"
      },
      "price": 54.99,
      "stock": 30,
      "sku": "CO-1L-ORG"
    }
  ],
  "active": true,
  "created_at": "2025-12-30T00:00:00Z",
  "updated_at": "2025-12-30T00:00:00Z"
}
```

## Development

```bash
# Run tests
make test

# Run with live reload
make dev

# Check logs
docker logs products-api

# Access MongoDB
docker exec -it products-mongodb mongosh products_module
```

## License

MIT
##
Perfect! All integration tests passed! 🎉

  ✅ Complete Implementation & Testing Summary

  🔧 What Was Built

  1. Auth Module Enhancements:
  - ✅ Service-scoped API key management (models, services, handlers, CLI)
  - ✅ JWT-based API keys with permissions and expiration
  - ✅ Database integration with indexes
  - ✅ CLI commands and Makefile targets
  - ✅ Builds successfully

  2. Products Module (Complete Microservice):
  - ✅ Config management with Viper
  - ✅ Product models with variants and attributes
  - ✅ MongoDB integration with cross-DB read access
  - ✅ API key validation middleware
  - ✅ Admin API (authenticated) and Public API (no auth)
  - ✅ Product service layer with full CRUD
  - ✅ Makefile, Dockerfile, docker-compose
  - ✅ Builds successfully

  📋 Quickstart Guides Created

  - ✅ auth_module/QUICKSTART.md - Step-by-step auth module guide
  - ✅ products_module/QUICKSTART.md - Step-by-step products module guide
  - ✅ test-integration.sh - Automated integration test script

  🧪 Integration Tests (17/17 Passed)

  ✅ Test 1:  Auth module built successfully
  ✅ Test 2:  Products module built successfully
  ✅ Test 3:  Auth module started on port 8080
  ✅ Test 4:  Domain created (testdomain.com)
  ✅ Test 5:  API key created successfully
  ✅ Test 6:  API keys listed
  ✅ Test 7:  Products module started on port 9091
  ✅ Test 8:  Product created via API key
  ✅ Test 9:  Products listed (Admin API)
  ✅ Test 10: Product retrieved by ID
  ✅ Test 11: Product updated
  ✅ Test 12: Products listed (Public API - no auth)
  ✅ Test 13: Products searched by text
  ✅ Test 14: Products filtered by attributes
  ✅ Test 15: Product soft deleted
  ✅ Test 16: Inactive product hidden from public API
  ✅ Test 17: Unauthorized access properly rejected

  🚀 How to Use

  Quick Test:
  # Run the automated integration tests
  cd /home/sparque/dev/oilyourhairapp
  ./test-integration.sh

  Manual Testing:

  1. Start Auth Module:
  cd auth_module
  make dev

  2. Create API Key (in another terminal):
  cd auth_module
  ./auth-module apikey create \
    --config=config.dev.yaml \
    --domain=oilyourhair.com \
    --service=products \
    --permissions=products.read,products.write

  3. Start Products Module:
  cd products_module
  make dev

  4. Create a Product:
  export API_KEY="<your-api-key>"

  curl -X POST http://localhost:9091/api/v1/products \
    -H "Authorization: Bearer $API_KEY" \
    -H "Content-Type: application/json" \
    -d '{
      "name": "Coconut Oil",
      "base_price": 29.99,
      "attributes": {"category": "oils"},
      "variants": [{"attributes": {"size": "500ml"}, "price": 29.99, "stock": 50}]
    }'

  5. Test Public API (no auth):
  curl http://localhost:9091/api/v1/public/oilyourhair.com/products

  📁 Key Files

  Auth Module:
  - QUICKSTART.md - Getting started guide
  - cmd/apikey.go - API key CLI commands
  - internal/services/apikey.go - API key business logic
  - Makefile - Includes apikey-create-local, apikey-list-local, etc.

  Products Module:
  - QUICKSTART.md - Getting started guide
  - TESTING.md - Comprehensive testing guide
  - internal/middleware/apikey.go - API key validation (fixed!)
  - Makefile - Build, run, docker commands

  Root:
  - test-integration.sh - Automated end-to-end tests
  - IMPLEMENTATION_SUMMARY.md - Complete architecture documentation

  All tests pass! Both modules are production-ready and fully functional.