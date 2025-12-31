#!/bin/bash

# Integration Test Script for Auth Module + Products Module
# Tests the complete flow: domain creation → API key → product CRUD

set -e  # Exit on error

# Colors for output
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${YELLOW}🧪 Starting Integration Tests${NC}"
echo ""

# Configuration
AUTH_PORT=8080
PRODUCTS_PORT=9091
DOMAIN="testdomain.com"
JWT_SECRET="test-secret-key-for-integration-testing"

# Cleanup function
cleanup() {
    echo ""
    echo -e "${YELLOW}🧹 Cleaning up...${NC}"

    # Kill background processes
    if [ ! -z "$AUTH_PID" ]; then
        kill $AUTH_PID 2>/dev/null || true
    fi
    if [ ! -z "$PRODUCTS_PID" ]; then
        kill $PRODUCTS_PID 2>/dev/null || true
    fi

    # Wait for ports to be released
    sleep 2

    echo -e "${GREEN}✅ Cleanup complete${NC}"
}

# Set trap to cleanup on exit
trap cleanup EXIT

# Test 1: Build Auth Module
echo -e "${YELLOW}📦 Test 1: Building auth_module...${NC}"
cd auth_module
go build -o auth-module . || { echo -e "${RED}❌ Auth module build failed${NC}"; exit 1; }
echo -e "${GREEN}✅ Auth module built successfully${NC}"
echo ""

# Test 2: Build Products Module
echo -e "${YELLOW}📦 Test 2: Building products_module...${NC}"
cd ../products_module
go build -o products-module . || { echo -e "${RED}❌ Products module build failed${NC}"; exit 1; }
echo -e "${GREEN}✅ Products module built successfully${NC}"
echo ""

# Test 3: Start Auth Module
echo -e "${YELLOW}🚀 Test 3: Starting auth_module...${NC}"
cd ../auth_module

# Set environment variables
export AUTH_JWT_SECRET="$JWT_SECRET"
export AUTH_MONGODB_URI="mongodb://localhost:27017"
export AUTH_MONGODB_DATABASE="auth_module_test"
export AUTH_SERVER_PORT="$AUTH_PORT"

./auth-module serve --config=config.dev.yaml > /tmp/auth-module.log 2>&1 &
AUTH_PID=$!

# Wait for auth module to start
sleep 3

# Check if auth module is running
if ! curl -s http://localhost:$AUTH_PORT/health > /dev/null; then
    echo -e "${RED}❌ Auth module failed to start${NC}"
    cat /tmp/auth-module.log
    exit 1
fi

echo -e "${GREEN}✅ Auth module started on port $AUTH_PORT${NC}"
echo ""

# Test 4: Create Domain
echo -e "${YELLOW}🌐 Test 4: Creating domain...${NC}"
./auth-module domain create \
    --config=config.dev.yaml \
    --domain=$DOMAIN \
    --name="Test Domain" \
    --admin-email=admin@$DOMAIN || { echo -e "${RED}❌ Domain creation failed${NC}"; exit 1; }

echo -e "${GREEN}✅ Domain created: $DOMAIN${NC}"
echo ""

# Test 5: Create API Key
echo -e "${YELLOW}🔑 Test 5: Creating API key...${NC}"
API_KEY_OUTPUT=$(./auth-module apikey create \
    --config=config.dev.yaml \
    --domain=$DOMAIN \
    --service=products \
    --description="Test API key" \
    --permissions=products.read,products.write \
    --expires-in=365 2>&1)

# Extract API key from output (it's the JWT token line starting with eyJ)
API_KEY=$(echo "$API_KEY_OUTPUT" | grep "eyJ" | awk '{print $NF}')

echo "Extracted API Key: $API_KEY"

if [ -z "$API_KEY" ]; then
    echo -e "${RED}❌ Failed to extract API key${NC}"
    exit 1
fi

echo -e "${GREEN}✅ API key created${NC}"
echo ""

# Test 6: List API Keys
echo -e "${YELLOW}📋 Test 6: Listing API keys...${NC}"
./auth-module apikey list \
    --config=config.dev.yaml \
    --domain=$DOMAIN || { echo -e "${RED}❌ Failed to list API keys${NC}"; exit 1; }

echo -e "${GREEN}✅ API keys listed${NC}"
echo ""

# Test 7: Start Products Module
echo -e "${YELLOW}🚀 Test 7: Starting products_module...${NC}"
cd ../products_module

# Set environment variables
export PRODUCTS_JWT_SECRET="$JWT_SECRET"
export PRODUCTS_MONGODB_URI="mongodb://localhost:27017"
export PRODUCTS_MONGODB_DATABASE="products_module_test"
export PRODUCTS_SERVER_PORT="$PRODUCTS_PORT"
export PRODUCTS_AUTH_DOMAINS_DB="auth_module_test"

./products-module serve --config=config.dev.yaml > /tmp/products-module.log 2>&1 &
PRODUCTS_PID=$!

# Wait for products module to start
sleep 3

# Check if products module is running
if ! curl -s http://localhost:$PRODUCTS_PORT/health > /dev/null; then
    echo -e "${RED}❌ Products module failed to start${NC}"
    cat /tmp/products-module.log
    exit 1
fi

echo -e "${GREEN}✅ Products module started on port $PRODUCTS_PORT${NC}"
echo ""

# Test 8: Create Product
echo -e "${YELLOW}🛍️  Test 8: Creating product...${NC}"
CREATE_RESPONSE=$(curl -s -X POST http://localhost:$PRODUCTS_PORT/api/v1/products \
  -H "Authorization: Bearer $API_KEY" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Test Coconut Oil",
    "description": "Test product for integration testing",
    "base_price": 29.99,
    "images": ["https://example.com/test.jpg"],
    "attributes": {
      "category": "oils",
      "type": "coconut",
      "test": "true"
    },
    "variants": [
      {
        "attributes": {"size": "500ml"},
        "price": 29.99,
        "stock": 50,
        "sku": "TEST-500"
      }
    ]
  }')

echo "$CREATE_RESPONSE" | jq '.' || echo "$CREATE_RESPONSE"

# Extract product ID
PRODUCT_ID=$(echo "$CREATE_RESPONSE" | jq -r '.id // empty')

if [ -z "$PRODUCT_ID" ]; then
    echo -e "${RED}❌ Failed to create product${NC}"
    exit 1
fi

echo -e "${GREEN}✅ Product created with ID: $PRODUCT_ID${NC}"
echo ""

# Test 9: List Products (Admin API)
echo -e "${YELLOW}📋 Test 9: Listing products (Admin API)...${NC}"
LIST_RESPONSE=$(curl -s http://localhost:$PRODUCTS_PORT/api/v1/products \
  -H "Authorization: Bearer $API_KEY")

echo "$LIST_RESPONSE" | jq '.'

PRODUCT_COUNT=$(echo "$LIST_RESPONSE" | jq '.count // 0')
if [ "$PRODUCT_COUNT" -eq 0 ]; then
    echo -e "${RED}❌ No products found${NC}"
    exit 1
fi

echo -e "${GREEN}✅ Found $PRODUCT_COUNT product(s)${NC}"
echo ""

# Test 10: Get Product by ID (Admin API)
echo -e "${YELLOW}🔍 Test 10: Getting product by ID...${NC}"
GET_RESPONSE=$(curl -s http://localhost:$PRODUCTS_PORT/api/v1/products/$PRODUCT_ID \
  -H "Authorization: Bearer $API_KEY")

echo "$GET_RESPONSE" | jq '.'

PRODUCT_NAME=$(echo "$GET_RESPONSE" | jq -r '.name // empty')
if [ "$PRODUCT_NAME" != "Test Coconut Oil" ]; then
    echo -e "${RED}❌ Product name mismatch${NC}"
    exit 1
fi

echo -e "${GREEN}✅ Product retrieved successfully${NC}"
echo ""

# Test 11: Update Product
echo -e "${YELLOW}✏️  Test 11: Updating product...${NC}"
UPDATE_RESPONSE=$(curl -s -X PUT http://localhost:$PRODUCTS_PORT/api/v1/products/$PRODUCT_ID \
  -H "Authorization: Bearer $API_KEY" \
  -H "Content-Type: application/json" \
  -d '{
    "description": "Updated test description",
    "base_price": 27.99
  }')

echo "$UPDATE_RESPONSE" | jq '.'

UPDATED_PRICE=$(echo "$UPDATE_RESPONSE" | jq -r '.base_price // 0')
if [ "$UPDATED_PRICE" != "27.99" ]; then
    echo -e "${RED}❌ Product update failed${NC}"
    exit 1
fi

echo -e "${GREEN}✅ Product updated successfully${NC}"
echo ""

# Test 12: List Products (Public API - No Auth)
echo -e "${YELLOW}🌍 Test 12: Listing products (Public API)...${NC}"
PUBLIC_RESPONSE=$(curl -s http://localhost:$PRODUCTS_PORT/api/v1/public/$DOMAIN/products)

echo "$PUBLIC_RESPONSE" | jq '.'

PUBLIC_COUNT=$(echo "$PUBLIC_RESPONSE" | jq '.count // 0')
if [ "$PUBLIC_COUNT" -eq 0 ]; then
    echo -e "${RED}❌ No public products found${NC}"
    exit 1
fi

echo -e "${GREEN}✅ Public API returned $PUBLIC_COUNT product(s)${NC}"
echo ""

# Test 13: Search Products
echo -e "${YELLOW}🔎 Test 13: Searching products...${NC}"
SEARCH_RESPONSE=$(curl -s "http://localhost:$PRODUCTS_PORT/api/v1/public/$DOMAIN/products/search?q=coconut")

echo "$SEARCH_RESPONSE" | jq '.'

SEARCH_COUNT=$(echo "$SEARCH_RESPONSE" | jq '.count // 0')
if [ "$SEARCH_COUNT" -eq 0 ]; then
    echo -e "${RED}❌ Search returned no results${NC}"
    exit 1
fi

echo -e "${GREEN}✅ Search found $SEARCH_COUNT product(s)${NC}"
echo ""

# Test 14: Filter by Attributes
echo -e "${YELLOW}🔍 Test 14: Filtering by attributes...${NC}"
FILTER_RESPONSE=$(curl -s "http://localhost:$PRODUCTS_PORT/api/v1/public/$DOMAIN/products?category=oils")

echo "$FILTER_RESPONSE" | jq '.'

FILTER_COUNT=$(echo "$FILTER_RESPONSE" | jq '.count // 0')
if [ "$FILTER_COUNT" -eq 0 ]; then
    echo -e "${RED}❌ Filter returned no results${NC}"
    exit 1
fi

echo -e "${GREEN}✅ Filter found $FILTER_COUNT product(s)${NC}"
echo ""

# Test 15: Soft Delete Product
echo -e "${YELLOW}🗑️  Test 15: Soft deleting product...${NC}"
DELETE_RESPONSE=$(curl -s -X DELETE http://localhost:$PRODUCTS_PORT/api/v1/products/$PRODUCT_ID \
  -H "Authorization: Bearer $API_KEY")

echo "$DELETE_RESPONSE" | jq '.'

echo -e "${GREEN}✅ Product soft deleted${NC}"
echo ""

# Test 16: Verify Product is Inactive
echo -e "${YELLOW}🔍 Test 16: Verifying product is inactive...${NC}"
PUBLIC_RESPONSE_AFTER_DELETE=$(curl -s http://localhost:$PRODUCTS_PORT/api/v1/public/$DOMAIN/products)

PUBLIC_COUNT_AFTER=$(echo "$PUBLIC_RESPONSE_AFTER_DELETE" | jq '.count // 0')
if [ "$PUBLIC_COUNT_AFTER" -ne 0 ]; then
    echo -e "${RED}❌ Inactive product still visible in public API${NC}"
    exit 1
fi

echo -e "${GREEN}✅ Inactive product not visible in public API${NC}"
echo ""

# Test 17: Verify Unauthorized Access Fails
echo -e "${YELLOW}🔒 Test 17: Testing unauthorized access...${NC}"
UNAUTH_RESPONSE=$(curl -s -w "\n%{http_code}" -X POST http://localhost:$PRODUCTS_PORT/api/v1/products \
  -H "Content-Type: application/json" \
  -d '{"name": "Should Fail", "base_price": 10}')

HTTP_CODE=$(echo "$UNAUTH_RESPONSE" | tail -n1)
if [ "$HTTP_CODE" != "401" ]; then
    echo -e "${RED}❌ Expected 401 but got $HTTP_CODE${NC}"
    exit 1
fi

echo -e "${GREEN}✅ Unauthorized access properly rejected${NC}"
echo ""

# All tests passed!
echo ""
echo -e "${GREEN}════════════════════════════════════════${NC}"
echo -e "${GREEN}🎉 All Integration Tests Passed! 🎉${NC}"
echo -e "${GREEN}════════════════════════════════════════${NC}"
echo ""
echo -e "Summary:"
echo -e "  ✅ Auth module built and started"
echo -e "  ✅ Products module built and started"
echo -e "  ✅ Domain created"
echo -e "  ✅ API key created and listed"
echo -e "  ✅ Product CRUD operations"
echo -e "  ✅ Public API access"
echo -e "  ✅ Search and filtering"
echo -e "  ✅ Soft delete"
echo -e "  ✅ Authorization checks"
echo ""
echo -e "${YELLOW}View logs:${NC}"
echo -e "  Auth module: /tmp/auth-module.log"
echo -e "  Products module: /tmp/products-module.log"
echo ""
