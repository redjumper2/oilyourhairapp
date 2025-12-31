#!/bin/bash

# Test script for oilyourhair.com website integration
# This script creates sample products and opens the website

set -e

echo "🌿 Oil Your Hair - Website Test"
echo ""

# Check if products module is running
if ! curl -s http://localhost:9091/health > /dev/null 2>&1; then
    echo "❌ Products module is not running on port 9091"
    echo "   Please start it with: cd ../products_module && ./products-module serve --config=config.dev.yaml"
    exit 1
fi

echo "✅ Products module is running"
echo ""

# Check for API key
if [ -z "$API_KEY" ]; then
    echo "⚠️  API_KEY environment variable not set"
    echo ""
    echo "To create an API key, run:"
    echo "  cd ../auth_module"
    echo "  ./auth-module apikey create \\"
    echo "    --config=config.dev.yaml \\"
    echo "    --domain=oilyourhair.com \\"
    echo "    --service=products \\"
    echo "    --permissions=products.read,products.write"
    echo ""
    echo "Then export it:"
    echo "  export API_KEY='your-api-key-here'"
    echo ""
    exit 1
fi

echo "✅ API key found"
echo ""

# Create sample products
echo "📦 Creating sample products..."
echo ""

# Product 1: Coconut Oil
echo "Creating Coconut Oil..."
curl -s -X POST http://localhost:9091/api/v1/products \
  -H "Authorization: Bearer $API_KEY" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Organic Coconut Oil",
    "description": "Pure, cold-pressed coconut oil for deep conditioning and hair growth. Rich in vitamins and fatty acids.",
    "base_price": 29.99,
    "images": ["https://images.unsplash.com/photo-1587334206443-0c0e1f5bf7cd?w=500"],
    "attributes": {
      "category": "oils",
      "type": "coconut",
      "organic": "true"
    },
    "variants": [
      {
        "attributes": {"size": "250ml"},
        "price": 29.99,
        "stock": 50,
        "sku": "CO-250"
      },
      {
        "attributes": {"size": "500ml"},
        "price": 49.99,
        "stock": 30,
        "sku": "CO-500"
      },
      {
        "attributes": {"size": "1L"},
        "price": 89.99,
        "stock": 15,
        "sku": "CO-1L"
      }
    ]
  }' > /dev/null

# Product 2: Argan Oil
echo "Creating Argan Oil..."
curl -s -X POST http://localhost:9091/api/v1/products \
  -H "Authorization: Bearer $API_KEY" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Moroccan Argan Oil",
    "description": "Premium argan oil from Morocco. Perfect for frizzy and damaged hair. Adds shine and softness.",
    "base_price": 39.99,
    "images": ["https://images.unsplash.com/photo-1608571423902-eed4a5ad8108?w=500"],
    "attributes": {
      "category": "oils",
      "type": "argan",
      "organic": "true"
    },
    "variants": [
      {
        "attributes": {"size": "100ml"},
        "price": 39.99,
        "stock": 40,
        "sku": "AR-100"
      },
      {
        "attributes": {"size": "250ml"},
        "price": 79.99,
        "stock": 25,
        "sku": "AR-250"
      }
    ]
  }' > /dev/null

# Product 3: Hair Treatment
echo "Creating Deep Conditioning Treatment..."
curl -s -X POST http://localhost:9091/api/v1/products \
  -H "Authorization: Bearer $API_KEY" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Deep Conditioning Treatment",
    "description": "Intensive repair treatment for dry and damaged hair. Use weekly for best results.",
    "base_price": 24.99,
    "images": ["https://images.unsplash.com/photo-1571875257727-256c39da42af?w=500"],
    "attributes": {
      "category": "treatments",
      "type": "conditioner"
    },
    "variants": [
      {
        "attributes": {"size": "200ml"},
        "price": 24.99,
        "stock": 60,
        "sku": "TR-200"
      }
    ]
  }' > /dev/null

# Product 4: Castor Oil
echo "Creating Castor Oil..."
curl -s -X POST http://localhost:9091/api/v1/products \
  -H "Authorization: Bearer $API_KEY" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Jamaican Black Castor Oil",
    "description": "Authentic Jamaican black castor oil for hair growth and thickness. Traditional formula.",
    "base_price": 19.99,
    "images": ["https://images.unsplash.com/photo-1596462502278-27bfdc403348?w=500"],
    "attributes": {
      "category": "oils",
      "type": "castor",
      "organic": "false"
    },
    "variants": [
      {
        "attributes": {"size": "250ml"},
        "price": 19.99,
        "stock": 35,
        "sku": "CA-250"
      },
      {
        "attributes": {"size": "500ml"},
        "price": 34.99,
        "stock": 20,
        "sku": "CA-500"
      }
    ]
  }' > /dev/null

echo ""
echo "✅ Sample products created!"
echo ""

# Get products count
PRODUCTS=$(curl -s http://localhost:9091/api/v1/public/oilyourhair.com/products | grep -o '"count":[0-9]*' | cut -d':' -f2)

echo "📊 Products in catalog: $PRODUCTS"
echo ""

# Start a simple HTTP server
echo "🚀 Starting web server on http://localhost:3000"
echo ""
echo "   Shop page: http://localhost:3000/shop.html"
echo "   Admin page: http://localhost:3000/admin/products.html"
echo ""
echo "Press Ctrl+C to stop the server"
echo ""

cd public
python3 -m http.server 3000
