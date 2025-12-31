#!/bin/bash

# Script to seed oilyourhair.com products into the products module
# Based on the existing product catalog

set -e

echo "🌿 Seeding oilyourhair.com Products"
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
echo "📦 Adding products to catalog..."
echo ""

# Product 1: Argan Oil Elixir
echo "Adding Argan Oil Elixir..."
curl -s -X POST http://localhost:9091/api/v1/products \
  -H "Authorization: Bearer $API_KEY" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Argan Oil Elixir",
    "description": "Pure Moroccan argan oil for deep nourishment and shine. Perfect for dry and damaged hair. Rich in vitamins and antioxidants.",
    "base_price": 45.99,
    "images": ["https://images.unsplash.com/photo-1571781926291-c477ebfd024b?ixlib=rb-4.0.3&auto=format&fit=crop&w=500&q=80"],
    "attributes": {
      "category": "oils",
      "type": "argan",
      "organic": "true",
      "badge": "Bestseller",
      "hair_type": "dry,damaged,normal",
      "features": "Organic,Cold-Pressed,Anti-Aging"
    },
    "variants": [
      {
        "attributes": {"size": "50ml"},
        "price": 45.99,
        "stock": 25,
        "sku": "AR-50"
      },
      {
        "attributes": {"size": "100ml"},
        "price": 79.99,
        "stock": 40,
        "sku": "AR-100"
      }
    ]
  }' > /dev/null
echo "✅ Added Argan Oil Elixir"

# Product 2: Coconut Miracle Oil
echo "Adding Coconut Miracle Oil..."
curl -s -X POST http://localhost:9091/api/v1/products \
  -H "Authorization: Bearer $API_KEY" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Coconut Miracle Oil",
    "description": "Virgin coconut oil for moisture and protection. Perfect for dry, curly, and normal hair types. Provides heat protection and deep conditioning.",
    "base_price": 29.99,
    "images": ["https://images.unsplash.com/photo-1599351431202-1e0f0137899a?ixlib=rb-4.0.3&auto=format&fit=crop&w=500&q=80"],
    "attributes": {
      "category": "oils",
      "type": "coconut",
      "organic": "true",
      "hair_type": "dry,curly,normal",
      "features": "Virgin,Moisturizing,Heat Protection"
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
      }
    ]
  }' > /dev/null
echo "✅ Added Coconut Miracle Oil"

# Product 3: Jojoba Scalp Treatment
echo "Adding Jojoba Scalp Treatment..."
curl -s -X POST http://localhost:9091/api/v1/products \
  -H "Authorization: Bearer $API_KEY" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Jojoba Scalp Treatment",
    "description": "Lightweight jojoba oil for scalp health and balance. Non-greasy formula perfect for oily and normal hair. Promotes healthy scalp environment.",
    "base_price": 38.99,
    "images": ["https://images.unsplash.com/photo-1556228720-195a672e8a03?ixlib=rb-4.0.3&auto=format&fit=crop&w=500&q=80"],
    "attributes": {
      "category": "treatments",
      "type": "jojoba",
      "organic": "true",
      "badge": "New",
      "hair_type": "oily,normal,damaged",
      "features": "Lightweight,Scalp Care,Non-Greasy"
    },
    "variants": [
      {
        "attributes": {"size": "100ml"},
        "price": 38.99,
        "stock": 35,
        "sku": "JO-100"
      },
      {
        "attributes": {"size": "200ml"},
        "price": 64.99,
        "stock": 20,
        "sku": "JO-200"
      }
    ]
  }' > /dev/null
echo "✅ Added Jojoba Scalp Treatment"

# Product 4: Rosemary Growth Serum
echo "Adding Rosemary Growth Serum..."
curl -s -X POST http://localhost:9091/api/v1/products \
  -H "Authorization: Bearer $API_KEY" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Rosemary Growth Serum",
    "description": "Stimulating rosemary oil blend for hair growth. Herbal formula that promotes circulation and strengthens follicles. Popular choice for growth support.",
    "base_price": 52.99,
    "images": ["https://images.unsplash.com/photo-1608248543803-ba4f8c70ae0b?ixlib=rb-4.0.3&auto=format&fit=crop&w=500&q=80"],
    "attributes": {
      "category": "treatments",
      "type": "rosemary",
      "organic": "true",
      "badge": "Popular",
      "hair_type": "dry,damaged,normal",
      "features": "Growth Formula,Stimulating,Herbal"
    },
    "variants": [
      {
        "attributes": {"size": "60ml"},
        "price": 52.99,
        "stock": 45,
        "sku": "RO-60"
      },
      {
        "attributes": {"size": "120ml"},
        "price": 89.99,
        "stock": 25,
        "sku": "RO-120"
      }
    ]
  }' > /dev/null
echo "✅ Added Rosemary Growth Serum"

# Product 5: Castor Oil Strength Booster
echo "Adding Castor Oil Strength Booster..."
curl -s -X POST http://localhost:9091/api/v1/products \
  -H "Authorization: Bearer $API_KEY" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Castor Oil Strength Booster",
    "description": "Thick castor oil for strengthening and thickness. Rich formula perfect for damaged, dry, and curly hair. Traditional strengthening treatment.",
    "base_price": 34.99,
    "images": ["https://images.unsplash.com/photo-1556228578-0d85b1a4d571?ixlib=rb-4.0.3&auto=format&fit=crop&w=500&q=80"],
    "attributes": {
      "category": "oils",
      "type": "castor",
      "organic": "false",
      "hair_type": "damaged,dry,curly",
      "features": "Strengthening,Thickening,Rich Formula"
    },
    "variants": [
      {
        "attributes": {"size": "250ml"},
        "price": 34.99,
        "stock": 35,
        "sku": "CA-250"
      },
      {
        "attributes": {"size": "500ml"},
        "price": 59.99,
        "stock": 20,
        "sku": "CA-500"
      }
    ]
  }' > /dev/null
echo "✅ Added Castor Oil Strength Booster"

# Product 6: Luxury Hair Oil Blend
echo "Adding Luxury Hair Oil Blend..."
curl -s -X POST http://localhost:9091/api/v1/products \
  -H "Authorization: Bearer $API_KEY" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Luxury Hair Oil Blend",
    "description": "Premium blend of 7 exotic oils for ultimate luxury. Our signature formula combines argan, jojoba, coconut, rosemary, lavender, and more for complete hair transformation.",
    "base_price": 89.99,
    "images": ["https://images.unsplash.com/photo-1590736969955-71cc94901144?ixlib=rb-4.0.3&auto=format&fit=crop&w=500&q=80"],
    "attributes": {
      "category": "oils",
      "type": "blend",
      "organic": "true",
      "badge": "Premium",
      "hair_type": "dry,damaged,normal,curly",
      "features": "Luxury Blend,7 Oils,Premium"
    },
    "variants": [
      {
        "attributes": {"size": "100ml"},
        "price": 89.99,
        "stock": 15,
        "sku": "LUX-100"
      },
      {
        "attributes": {"size": "200ml"},
        "price": 159.99,
        "stock": 10,
        "sku": "LUX-200"
      }
    ]
  }' > /dev/null
echo "✅ Added Luxury Hair Oil Blend"

echo ""
echo "✅ All 6 products added successfully!"
echo ""

# Get total products count
TOTAL=$(curl -s http://localhost:9091/api/v1/public/oilyourhair.com/products | grep -o '"count":[0-9]*' | cut -d':' -f2)
echo "📊 Total products in catalog: $TOTAL"
echo ""
echo "🌐 View products at:"
echo "   http://localhost:3000/shop.html"
echo ""
