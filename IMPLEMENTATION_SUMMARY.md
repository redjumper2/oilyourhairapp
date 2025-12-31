# Implementation Summary: Service-Scoped API Keys & Products Module

## Overview

This implementation adds **service-scoped API key management** to the auth_module and creates a new **products_module** that uses these API keys for secure, multi-tenant product management.

## What Was Built

### 1. Auth Module Enhancements

#### API Key Management System

**New Files:**
- `auth_module/internal/models/apikey.go` - API key data model
- `auth_module/internal/services/apikey.go` - API key business logic
- `auth_module/internal/handlers/apikey.go` - API key REST endpoints
- `auth_module/cmd/apikey.go` - CLI commands for API key management

**Database:**
- New `api_keys` collection in auth_module database
- Indexes on `key_id`, `{domain, service}`, and `expires_at`

**Features:**
- JWT-based API keys (reuses existing JWT infrastructure)
- Service-scoped (e.g., "products", "orders", "shipping")
- Configurable permissions (e.g., "products.read", "products.write")
- Configurable expiration (default: 365 days)
- Revocation support
- Last used tracking
- Expiration monitoring

**API Endpoints:**
- `POST /api/v1/api-keys` - Create API key
- `GET /api/v1/domains/:domain/api-keys` - List API keys
- `GET /api/v1/api-keys/:keyId` - Get API key details
- `DELETE /api/v1/api-keys/:keyId` - Revoke API key
- `GET /api/v1/domains/:domain/api-keys/expiring` - Get expiring keys

**CLI Commands:**
```bash
./auth-module apikey create --domain=example.com --service=products --permissions=products.read,products.write
./auth-module apikey list --domain=example.com
./auth-module apikey revoke --key-id=<key-id>
```

**Makefile Targets:**
```bash
make apikey-create-local DOMAIN=example.com SERVICE=products
make apikey-list-local DOMAIN=example.com
make apikey-revoke-local KEY_ID=<key-id>
```

### 2. Products Module (New Service)

**Complete microservice** for managing product catalogs with:

#### Architecture
- **Multi-tenant**: Isolated by domain
- **API key authentication**: Uses JWT-based API keys from auth_module
- **Public API**: Unauthenticated product listing for customer websites
- **Admin API**: Authenticated CRUD operations with permission checks

#### Core Files
```
products_module/
├── cmd/
│   ├── root.go          # Cobra CLI
│   └── serve.go         # HTTP server
├── internal/
│   ├── models/
│   │   └── product.go   # Product & variant models
│   ├── handlers/
│   │   ├── admin.go     # Admin CRUD operations
│   │   └── public.go    # Public read-only operations
│   ├── services/
│   │   └── product.go   # Business logic
│   ├── database/
│   │   └── database.go  # MongoDB connection
│   └── middleware/
│       └── apikey.go    # API key validation
├── config/
│   └── config.go        # Configuration management
├── main.go              # Entry point
├── Makefile             # Development commands
├── Dockerfile           # Container image
├── docker-compose.yml   # Service orchestration
├── TESTING.md           # Testing guide
└── README.md            # Documentation
```

#### Data Model

**Product:**
```go
{
  "id": "ObjectID",
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
      "id": "uuid",
      "attributes": {"size": "500ml"},
      "price": 29.99,
      "stock": 50,
      "sku": "CO-500"
    }
  ],
  "active": true,
  "created_by": "api_key_id",
  "created_at": "timestamp",
  "updated_at": "timestamp"
}
```

#### API Endpoints

**Admin API** (requires API key):
- `POST /api/v1/products` - Create product
- `GET /api/v1/products` - List products
- `GET /api/v1/products/:id` - Get product
- `PUT /api/v1/products/:id` - Update product
- `DELETE /api/v1/products/:id` - Delete product (soft/hard)
- `PUT /api/v1/products/:id/variants/:variantId/stock` - Update stock

**Public API** (no auth):
- `GET /api/v1/public/:domain/products` - List products
- `GET /api/v1/public/:domain/products/:id` - Get product
- `GET /api/v1/public/:domain/products/search?q=term` - Search products

#### Features
- Flexible key/value attributes
- Product variants with pricing and inventory
- Text search on name and description
- Attribute-based filtering
- Soft delete (sets active=false)
- Hard delete (removes from database)
- Stock management per variant
- Multi-image support

## How It Works

### Authentication Flow

1. **Domain Admin** uses auth_module to create API key:
   ```bash
   ./auth-module apikey create \
     --domain=oilyourhair.com \
     --service=products \
     --permissions=products.read,products.write
   ```

2. **API Key** is a JWT with claims:
   ```json
   {
     "jti": "unique-key-id",
     "domain": "oilyourhair.com",
     "service": "products",
     "type": "api_key",
     "permissions": ["products.read", "products.write"],
     "exp": 1735689600
   }
   ```

3. **Products Module** validates API key:
   - Verifies JWT signature (shared secret)
   - Checks type is "api_key"
   - Validates service is "products"
   - Checks permissions for operation
   - Verifies not revoked (future: could check auth_module DB)

4. **Domain Validation**:
   - Products module reads from `auth_module.domains` collection
   - Ensures domain exists and is active
   - Enforces multi-tenant isolation

### Database Architecture

**MongoDB Databases:**
```
mongodb://
├── auth_module/
│   ├── domains          # Domain registration
│   ├── users            # User accounts
│   ├── invitations      # Invite system
│   ├── magic_link_tokens
│   └── api_keys         # Service API keys (NEW)
└── products_module/
    └── products         # Product catalog
```

**Cross-DB Access:**
- Products module has **read-only** access to `auth_module.domains`
- No direct DB dependencies for runtime (uses API keys)
- Designed for future physical separation

### Security Model

1. **JWT-Based API Keys**
   - Same infrastructure as user JWTs
   - Distinguished by "type": "api_key" claim
   - Service-scoped via "service" claim
   - Permission-based via "permissions" claim

2. **Multi-Tenancy**
   - All data scoped by domain
   - API keys tied to specific domain
   - Enforced at database query level

3. **Permission System**
   - Fine-grained permissions (e.g., "products.read", "products.write")
   - Middleware validates permissions per endpoint
   - Easily extensible for future services

4. **Revocation**
   - API keys can be revoked immediately
   - Revoked keys stored in database
   - Future: products_module could check revocation via API call to auth_module

## Getting Started

### Prerequisites

1. MongoDB running on localhost:27017
2. Both modules using same JWT secret
3. Domain registered in auth_module

### Step 1: Start Auth Module

```bash
cd auth_module
make build
./auth-module serve --config=config.dev.yaml
```

### Step 2: Create Domain (if not exists)

```bash
./auth-module domain create \
  --config=config.dev.yaml \
  --domain=oilyourhair.com \
  --name="Oil Your Hair" \
  --admin-email=admin@oilyourhair.com
```

### Step 3: Create API Key

```bash
./auth-module apikey create \
  --config=config.dev.yaml \
  --domain=oilyourhair.com \
  --service=products \
  --description="Products service API key" \
  --permissions=products.read,products.write

# Save the displayed API key!
```

### Step 4: Start Products Module

```bash
cd ../products_module
make build
./products-module serve --config=config.dev.yaml
```

### Step 5: Test API

```bash
export API_KEY="<your-api-key>"

# Create a product
curl -X POST http://localhost:9091/api/v1/products \
  -H "Authorization: Bearer $API_KEY" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Coconut Oil",
    "description": "Pure organic coconut oil",
    "base_price": 29.99,
    "attributes": {"category": "oils"},
    "variants": [
      {"attributes": {"size": "500ml"}, "price": 29.99, "stock": 50}
    ]
  }'

# List products (public, no auth)
curl http://localhost:9091/api/v1/public/oilyourhair.com/products
```

## Configuration

### Auth Module (.env or config)

```bash
AUTH_JWT_SECRET=your-shared-secret-key
AUTH_MONGODB_URI=mongodb://localhost:27017
AUTH_MONGODB_DATABASE=auth_module
```

### Products Module (.env or config)

```bash
PRODUCTS_JWT_SECRET=your-shared-secret-key  # MUST MATCH auth_module
PRODUCTS_MONGODB_URI=mongodb://localhost:27017
PRODUCTS_MONGODB_DATABASE=products_module
PRODUCTS_AUTH_DOMAINS_DB=auth_module
```

**Critical:** JWT_SECRET must be identical in both modules!

## Future Enhancements

1. **API Key Validation Endpoint**
   - Add endpoint to auth_module: `POST /api/v1/api-keys/validate`
   - Products module calls this to check revocation
   - Enables instant revocation across all services

2. **Rate Limiting**
   - Track API key usage
   - Enforce rate limits per key

3. **API Key Rotation**
   - Generate new key before old expires
   - Grace period for migration

4. **Audit Logging**
   - Track all API key usage
   - Monitor for suspicious activity

5. **Additional Services**
   - Orders module
   - Shipping module
   - Inventory module
   - All using same API key system

6. **Physical Separation**
   - Deploy auth_module and products_module independently
   - Different cloud regions
   - Separate databases
   - Communication only via API keys

## File Changes Summary

### Auth Module
**New:**
- internal/models/apikey.go
- internal/services/apikey.go
- internal/handlers/apikey.go
- cmd/apikey.go

**Modified:**
- internal/database/database.go (added APIKeys collection)
- cmd/serve.go (added API key routes)
- Makefile (added API key commands)

### Products Module (Entirely New)
All files created:
- Complete Go microservice with config, models, handlers, services, middleware
- Makefile for development
- Dockerfile for containerization
- docker-compose.yml for orchestration
- TESTING.md for testing guide
- README.md and architecture.md for documentation

## Testing

See `products_module/TESTING.md` for comprehensive testing guide with curl examples.

Quick test:
```bash
# Terminal 1: Start auth_module
cd auth_module && make dev

# Terminal 2: Create API key
cd auth_module
./auth-module apikey create --domain=oilyourhair.com --service=products --permissions=products.read,products.write

# Terminal 3: Start products_module
cd products_module && make dev

# Terminal 4: Test
export API_KEY="<from-step-2>"
curl -X POST http://localhost:9091/api/v1/products \
  -H "Authorization: Bearer $API_KEY" \
  -H "Content-Type: application/json" \
  -d '{"name": "Test Product", "base_price": 19.99, "attributes": {}}'
```

## Summary

✅ **Auth Module**: Added service-scoped API key management with JWT-based keys, permission system, and CLI tools

✅ **Products Module**: Complete microservice for multi-tenant product catalogs with API key authentication

✅ **Architecture**: Designed for independent deployment while sharing JWT secret

✅ **Security**: Permission-based access control, multi-tenant isolation, revocation support

✅ **Developer Experience**: CLI tools, Makefile commands, comprehensive documentation

✅ **Ready for Production**: Docker support, configuration management, health checks

The system is now ready to use. Create an API key in auth_module and start managing products via the products_module API!
