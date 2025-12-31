# Products Module - Quick Start Guide

## Prerequisites

- Go 1.21+
- MongoDB running on localhost:27017
- Auth module running on localhost:8080
- API key created in auth module

## Step 1: Build the Module

```bash
make build
```

## Step 2: Configure

Edit `config.dev.yaml` and ensure JWT secret matches auth_module:

```yaml
jwt:
  secret: "your-secret-key-here"  # MUST match auth_module!
```

Or use environment variable:
```bash
export PRODUCTS_JWT_SECRET="your-secret-key-here"  # Same as AUTH_JWT_SECRET
```

## Step 3: Start the Server

```bash
./products-module serve --config=config.dev.yaml
```

Server will start on http://localhost:9091

## Step 4: Test Health Endpoint

```bash
curl http://localhost:9091/health
```

Expected response:
```json
{
  "status": "healthy",
  "env": "development"
}
```

## Step 5: Create Your First Product

Use the API key from auth_module:

```bash
export API_KEY="your-api-key-from-auth-module"

curl -X POST http://localhost:9091/api/v1/products \
  -H "Authorization: Bearer $API_KEY" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Coconut Oil",
    "description": "Pure organic coconut oil for hair care",
    "base_price": 29.99,
    "images": ["https://example.com/coconut-oil.jpg"],
    "attributes": {
      "category": "oils",
      "type": "coconut",
      "organic": "true"
    },
    "variants": [
      {
        "attributes": {"size": "500ml"},
        "price": 29.99,
        "stock": 50,
        "sku": "CO-500"
      },
      {
        "attributes": {"size": "1L"},
        "price": 54.99,
        "stock": 30,
        "sku": "CO-1L"
      }
    ]
  }'
```

## Step 6: List Products (Admin API)

```bash
curl http://localhost:9091/api/v1/products \
  -H "Authorization: Bearer $API_KEY"
```

## Step 7: List Products (Public API - No Auth)

```bash
curl http://localhost:9091/api/v1/public/oilyourhair.com/products
```

## Step 8: Search Products

```bash
curl "http://localhost:9091/api/v1/public/oilyourhair.com/products/search?q=coconut"
```

## Step 9: Filter by Attributes

```bash
curl "http://localhost:9091/api/v1/public/oilyourhair.com/products?category=oils&organic=true"
```

## Common Commands

```bash
# List all products
curl http://localhost:9091/api/v1/products \
  -H "Authorization: Bearer $API_KEY"

# Get specific product
curl http://localhost:9091/api/v1/products/<product-id> \
  -H "Authorization: Bearer $API_KEY"

# Update product
curl -X PUT http://localhost:9091/api/v1/products/<product-id> \
  -H "Authorization: Bearer $API_KEY" \
  -H "Content-Type: application/json" \
  -d '{"description": "Updated description", "base_price": 27.99}'

# Delete product (soft delete - sets active=false)
curl -X DELETE http://localhost:9091/api/v1/products/<product-id> \
  -H "Authorization: Bearer $API_KEY"

# Delete product (hard delete - removes from DB)
curl -X DELETE "http://localhost:9091/api/v1/products/<product-id>?hard=true" \
  -H "Authorization: Bearer $API_KEY"
```

## Troubleshooting

### 401 Unauthorized
- Verify API key is correct
- Check JWT_SECRET matches between auth_module and products_module

### 403 Forbidden
- API key doesn't have required permissions
- API key service must be "products"

### Domain not found or inactive
- Ensure domain is registered in auth_module
- Domain status must be "active"

## Next Steps

See `TESTING.md` for comprehensive testing guide with more examples.
