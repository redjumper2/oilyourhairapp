# Testing Guide

## Prerequisites

1. **Auth Module Running**: Ensure auth_module is running on port 8080
2. **MongoDB Running**: Shared MongoDB instance
3. **Domain Created**: Domain registered in auth_module (e.g., oilyourhair.com)
4. **Admin User**: Admin user exists for the domain

## Step 1: Create API Key

First, create an API key for the products service in auth_module:

```bash
cd ../auth_module

# Build auth_module
make build

# Create API key for products service
./auth-module apikey create \
  --config=config.dev.yaml \
  --domain=oilyourhair.com \
  --service=products \
  --description="Products API key for testing" \
  --permissions=products.read,products.write \
  --expires-in=365
```

**Save the API key** that is displayed. You'll need it for all subsequent requests.

Example output:
```
================================================================================
API Key: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
================================================================================
```

## Step 2: Start Products Module

```bash
cd ../products_module

# Build and run
make build
./products-module serve --config=config.dev.yaml
```

## Step 3: Test Health Endpoint

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

## Step 4: Create a Product (Admin API)

```bash
export API_KEY="your-api-key-from-step-1"

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
        "attributes": {
          "size": "500ml"
        },
        "price": 29.99,
        "stock": 50,
        "sku": "CO-500"
      },
      {
        "attributes": {
          "size": "1L"
        },
        "price": 54.99,
        "stock": 30,
        "sku": "CO-1L"
      }
    ]
  }'
```

## Step 5: List Products (Admin API)

```bash
curl -X GET http://localhost:9091/api/v1/products \
  -H "Authorization: Bearer $API_KEY"
```

## Step 6: Get Product by ID (Admin API)

```bash
# Get the product ID from the list response
export PRODUCT_ID="<product-id>"

curl -X GET http://localhost:9091/api/v1/products/$PRODUCT_ID \
  -H "Authorization: Bearer $API_KEY"
```

## Step 7: Update Product (Admin API)

```bash
curl -X PUT http://localhost:9091/api/v1/products/$PRODUCT_ID \
  -H "Authorization: Bearer $API_KEY" \
  -H "Content-Type: application/json" \
  -d '{
    "description": "Updated description for coconut oil",
    "base_price": 27.99
  }'
```

## Step 8: Update Stock (Admin API)

```bash
# Get variant ID from product details
export VARIANT_ID="<variant-id>"

curl -X PUT http://localhost:9091/api/v1/products/$PRODUCT_ID/variants/$VARIANT_ID/stock \
  -H "Authorization: Bearer $API_KEY" \
  -H "Content-Type: application/json" \
  -d '{
    "stock": 100
  }'
```

## Step 9: Test Public API (No Auth Required)

List products publicly:
```bash
curl http://localhost:9091/api/v1/public/oilyourhair.com/products
```

Get product by ID:
```bash
curl http://localhost:9091/api/v1/public/oilyourhair.com/products/$PRODUCT_ID
```

Search products:
```bash
curl "http://localhost:9091/api/v1/public/oilyourhair.com/products/search?q=coconut"
```

## Step 10: Filter Products by Attributes

Admin API:
```bash
curl "http://localhost:9091/api/v1/products?category=oils&type=coconut" \
  -H "Authorization: Bearer $API_KEY"
```

Public API:
```bash
curl "http://localhost:9091/api/v1/public/oilyourhair.com/products?category=oils&organic=true"
```

## Step 11: Soft Delete Product

```bash
curl -X DELETE http://localhost:9091/api/v1/products/$PRODUCT_ID \
  -H "Authorization: Bearer $API_KEY"
```

This sets `active: false` on the product.

## Step 12: Hard Delete Product

```bash
curl -X DELETE "http://localhost:9091/api/v1/products/$PRODUCT_ID?hard=true" \
  -H "Authorization: Bearer $API_KEY"
```

This permanently removes the product from the database.

## Common Issues

### 401 Unauthorized
- Check that your API key is valid and not expired
- Verify the JWT_SECRET matches between auth_module and products_module
- Ensure you're using "Bearer" prefix in Authorization header

### 403 Forbidden
- API key doesn't have required permissions
- API key is for a different service (not "products")
- API key domain doesn't match the requested resource

### 404 Not Found
- Product doesn't exist or is inactive (for public API)
- Invalid product ID format

### Domain not found or inactive
- The domain specified in the API key isn't registered in auth_module
- Domain status is not "active"

## JavaScript Example (Website Integration)

```html
<div id="products"></div>

<script>
const domain = 'oilyourhair.com';
const apiUrl = 'http://localhost:9091';

// Fetch products
fetch(`${apiUrl}/api/v1/public/${domain}/products?category=oils`)
  .then(response => response.json())
  .then(data => {
    const products = data.products;
    const html = products.map(p => `
      <div class="product">
        <h3>${p.name}</h3>
        <p>${p.description}</p>
        <p>Price: $${p.base_price}</p>
        ${p.variants.map(v => `
          <div class="variant">
            Size: ${v.attributes.size} - $${v.price}
            (${v.stock} in stock)
          </div>
        `).join('')}
      </div>
    `).join('');
    document.getElementById('products').innerHTML = html;
  });
</script>
```
