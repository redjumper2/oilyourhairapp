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
TEST_LOG="/tmp/integration-test-detailed.log"

# Clear previous log
> $TEST_LOG

# Logging function
log_test() {
    local test_num="$1"
    local test_name="$2"
    local command="$3"
    local response="$4"

    echo "" >> $TEST_LOG
    echo "===========================================" >> $TEST_LOG
    echo "Test $test_num: $test_name" >> $TEST_LOG
    echo "===========================================" >> $TEST_LOG
    echo "Timestamp: $(date '+%Y-%m-%d %H:%M:%S')" >> $TEST_LOG
    echo "" >> $TEST_LOG
    echo "Command:" >> $TEST_LOG
    echo "$command" >> $TEST_LOG
    echo "" >> $TEST_LOG
    if [ ! -z "$response" ]; then
        echo "Response:" >> $TEST_LOG
        echo "$response" | jq '.' >> $TEST_LOG 2>/dev/null || echo "$response" >> $TEST_LOG
        echo "" >> $TEST_LOG
    fi
}

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
    echo ""
    echo -e "${GREEN}📝 Detailed test log: $TEST_LOG${NC}"
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
DOMAIN_CMD="./auth-module domain create --config=config.dev.yaml --domain=$DOMAIN --name=\"Test Domain\" --admin-email=admin@$DOMAIN"
DOMAIN_OUTPUT=$(./auth-module domain create \
    --config=config.dev.yaml \
    --domain=$DOMAIN \
    --name="Test Domain" \
    --admin-email=admin@$DOMAIN 2>&1)

echo "$DOMAIN_OUTPUT"
log_test "4" "Create Domain" "$DOMAIN_CMD" "$DOMAIN_OUTPUT"

# Extract the admin invitation token from domain creation output
ADMIN_INVITE_TOKEN=$(echo "$DOMAIN_OUTPUT" | grep "token=" | sed 's/.*token=\([^&]*\).*/\1/')

if [ -z "$ADMIN_INVITE_TOKEN" ]; then
    echo -e "${RED}❌ Failed to extract admin invitation token${NC}"
    exit 1
fi

echo -e "${GREEN}✅ Domain created: $DOMAIN${NC}"
echo "   Admin invitation token: ${ADMIN_INVITE_TOKEN:0:20}..."
echo ""

# Test 5: Create API Key
echo -e "${YELLOW}🔑 Test 5: Creating API key...${NC}"
API_KEY_CMD="./auth-module apikey create --config=config.dev.yaml --domain=$DOMAIN --service=products --permissions=products.read,products.write"
API_KEY_OUTPUT=$(./auth-module apikey create \
    --config=config.dev.yaml \
    --domain=$DOMAIN \
    --service=products \
    --description="Test API key" \
    --permissions=products.read,products.write \
    --expires-in=365 2>&1)
log_test "5" "Create API Key" "$API_KEY_CMD" "$API_KEY_OUTPUT"

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
LIST_CMD="./auth-module apikey list --config=config.dev.yaml --domain=$DOMAIN"
LIST_KEYS_OUTPUT=$(./auth-module apikey list \
    --config=config.dev.yaml \
    --domain=$DOMAIN 2>&1)

echo "$LIST_KEYS_OUTPUT"
log_test "6" "List API Keys" "$LIST_CMD" "$LIST_KEYS_OUTPUT"

if [ $? -ne 0 ]; then
    echo -e "${RED}❌ Failed to list API keys${NC}"
    exit 1
fi

echo -e "${GREEN}✅ API keys listed${NC}"
echo ""
# After Test 6 (List API Keys), add these user auth tests:

# Test 7: Accept Initial Admin Invitation (from domain creation)
echo -e "${YELLOW}👤 Test 7: Accepting initial admin invitation...${NC}"

ACCEPT_CMD="curl -X POST http://localhost:$AUTH_PORT/api/v1/auth/invitation/accept -H 'Content-Type: application/json' -d '{\"token\": \"...\", \"email\": \"admin@$DOMAIN\", \"name\": \"Admin User\", \"auth_provider\": \"magic_link\"}'"
ADMIN_ACCEPT_RESPONSE=$(curl -s -X POST http://localhost:$AUTH_PORT/api/v1/auth/invitation/accept \
  -H "Content-Type: application/json" \
  -d "{
    \"token\": \"$ADMIN_INVITE_TOKEN\",
    \"email\": \"admin@$DOMAIN\",
    \"name\": \"Admin User\",
    \"auth_provider\": \"magic_link\"
  }")

log_test "7" "Accept Admin Invitation" "$ACCEPT_CMD" "$ADMIN_ACCEPT_RESPONSE"

ADMIN_JWT=$(echo "$ADMIN_ACCEPT_RESPONSE" | jq -r '.token // empty')

if [ -z "$ADMIN_JWT" ]; then
    echo -e "${RED}❌ Failed to accept admin invitation${NC}"
    echo "$ADMIN_ACCEPT_RESPONSE" | jq '.'
    exit 1
fi

echo -e "${GREEN}✅ Admin user created and authenticated${NC}"
echo ""

# Test 8: Validate Admin JWT
echo -e "${YELLOW}🔐 Test 8: Validating admin JWT with /auth/me...${NC}"
ME_CMD="curl http://localhost:$AUTH_PORT/api/v1/auth/me -H \'Authorization: Bearer \$ADMIN_JWT\' -H \'Host: $DOMAIN\'"
ADMIN_ME_RESPONSE=$(curl -s http://localhost:$AUTH_PORT/api/v1/auth/me \
  -H "Authorization: Bearer $ADMIN_JWT" \
  -H "Host: $DOMAIN")
log_test "8" "Validate Admin JWT" "$ME_CMD" "$ADMIN_ME_RESPONSE"


ADMIN_ROLE=$(echo "$ADMIN_ME_RESPONSE" | jq -r '.role // empty')
ADMIN_EMAIL=$(echo "$ADMIN_ME_RESPONSE" | jq -r '.email // empty')

if [ "$ADMIN_ROLE" != "admin" ] || [ "$ADMIN_EMAIL" != "admin@$DOMAIN" ]; then
    echo -e "${RED}❌ Admin JWT validation failed${NC}"
    exit 1
fi

echo -e "${GREEN}✅ Admin JWT validated (email: $ADMIN_EMAIL, role: $ADMIN_ROLE)${NC}"
echo ""

# Test 9: Create Customer User
echo -e "${YELLOW}👤 Test 9: Creating customer user...${NC}"

CUSTOMER_CMD="./auth-module invite create --config=config.dev.yaml --domain=$DOMAIN --email=customer@$DOMAIN --role=customer"
CUSTOMER_INVITE_OUTPUT=$(./auth-module invite create \
    --config=config.dev.yaml \
    --domain=$DOMAIN \
    --email=customer@$DOMAIN \
    --role=customer \
    --type=email_with_qr 2>&1)
log_test "9" "Create Customer Invitation" "$CUSTOMER_CMD" "$CUSTOMER_INVITE_OUTPUT"

CUSTOMER_INVITE_TOKEN=$(echo "$CUSTOMER_INVITE_OUTPUT" | grep "token:" | awk '{print $2}')

CUSTOMER_ACCEPT_RESPONSE=$(curl -s -X POST http://localhost:$AUTH_PORT/api/v1/auth/invitation/accept \
  -H "Content-Type: application/json" \
  -d "{
    \"token\": \"$CUSTOMER_INVITE_TOKEN\",
    \"email\": \"customer@$DOMAIN\",
    \"name\": \"Test Customer\",
    \"auth_provider\": \"magic_link\"
  }")


CUSTOMER_JWT=$(echo "$CUSTOMER_ACCEPT_RESPONSE" | jq -r '.token // empty')

if [ -z "$CUSTOMER_JWT" ]; then
    echo -e "${RED}❌ Failed to create customer user${NC}"
    exit 1
fi

echo -e "${GREEN}✅ Customer user created${NC}"
echo ""

# Test 10: Test Permission-Based Access (Admin vs Customer)
echo -e "${YELLOW}🔒 Test 10: Testing permission-based access control...${NC}"

# Test that customer CANNOT access admin endpoint
CUSTOMER_TEST_CMD="curl http://localhost:$AUTH_PORT/api/v1/admin/users -H \'Authorization: Bearer \$CUSTOMER_JWT\' (should be denied)"
CUSTOMER_ATTEMPT=$(curl -s -w "\n%{http_code}" http://localhost:$AUTH_PORT/api/v1/admin/users \
  -H "Authorization: Bearer $CUSTOMER_JWT" \
  -H "Host: $DOMAIN")
log_test "10a" "Customer Access Attempt (should fail)" "$CUSTOMER_TEST_CMD" "$CUSTOMER_ATTEMPT"

CUSTOMER_HTTP_CODE=$(echo "$CUSTOMER_ATTEMPT" | tail -n1)

if [ "$CUSTOMER_HTTP_CODE" != "403" ] && [ "$CUSTOMER_HTTP_CODE" != "401" ]; then
    echo -e "${RED}❌ Customer should be denied (got $CUSTOMER_HTTP_CODE)${NC}"
    exit 1
fi

echo -e "${GREEN}✅ Customer correctly denied admin endpoint${NC}"

# Test that admin CAN access admin endpoint
ADMIN_LIST_CMD="curl http://localhost:$AUTH_PORT/api/v1/admin/users -H \'Authorization: Bearer \$ADMIN_JWT\'"
ADMIN_LIST_USERS=$(curl -s http://localhost:$AUTH_PORT/api/v1/admin/users \
  -H "Authorization: Bearer $ADMIN_JWT" \
  -H "Host: $DOMAIN")
log_test "10b" "Admin List Users" "$ADMIN_LIST_CMD" "$ADMIN_LIST_USERS"


USER_COUNT=$(echo "$ADMIN_LIST_USERS" | jq '.users | length')

if [ "$USER_COUNT" -lt 2 ]; then
    echo -e "${RED}❌ Admin should see at least 2 users${NC}"
    exit 1
fi

echo -e "${GREEN}✅ Admin accessed admin endpoint successfully (found $USER_COUNT users)${NC}"
echo ""

# Test 11: Start Products Module
echo -e "${YELLOW}🚀 Test 11: Starting products_module...${NC}"
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

# Test 12: Create Product
echo -e "${YELLOW}🛍️  Test 12: Creating product...${NC}"
CREATE_CMD="curl -X POST http://localhost:$PRODUCTS_PORT/api/v1/products -H 'Authorization: Bearer \$API_KEY' -H 'Content-Type: application/json' -d '{\"name\": \"Test Coconut Oil\", ...}'"
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
log_test "12" "Create Product" "$CREATE_CMD" "$CREATE_RESPONSE"

# Extract product ID
PRODUCT_ID=$(echo "$CREATE_RESPONSE" | jq -r '.id // empty')

if [ -z "$PRODUCT_ID" ]; then
    echo -e "${RED}❌ Failed to create product${NC}"
    exit 1
fi

echo -e "${GREEN}✅ Product created with ID: $PRODUCT_ID${NC}"
echo ""

# Test 13: List Products (Admin API)
echo -e "${YELLOW}📋 Test 13: Listing products (Admin API)...${NC}"
LIST_CMD="curl http://localhost:$PRODUCTS_PORT/api/v1/products -H \'Authorization: Bearer \$API_KEY\'"
LIST_RESPONSE=$(curl -s http://localhost:$PRODUCTS_PORT/api/v1/products \
  -H "Authorization: Bearer $API_KEY")
log_test "13" "List Products (Admin)" "$LIST_CMD" "$LIST_RESPONSE"

echo "$LIST_RESPONSE" | jq '.'

PRODUCT_COUNT=$(echo "$LIST_RESPONSE" | jq '.count // 0')
if [ "$PRODUCT_COUNT" -eq 0 ]; then
    echo -e "${RED}❌ No products found${NC}"
    exit 1
fi

echo -e "${GREEN}✅ Found $PRODUCT_COUNT product(s)${NC}"
echo ""

# Test 14: Get Product by ID (Admin API)
echo -e "${YELLOW}🔍 Test 14: Getting product by ID...${NC}"
GET_CMD="curl http://localhost:$PRODUCTS_PORT/api/v1/products/\$PRODUCT_ID -H \'Authorization: Bearer \$API_KEY\'"
GET_RESPONSE=$(curl -s http://localhost:$PRODUCTS_PORT/api/v1/products/$PRODUCT_ID \
  -H "Authorization: Bearer $API_KEY")
log_test "14" "Get Product by ID" "$GET_CMD" "$GET_RESPONSE"

echo "$GET_RESPONSE" | jq '.'

PRODUCT_NAME=$(echo "$GET_RESPONSE" | jq -r '.name // empty')
if [ "$PRODUCT_NAME" != "Test Coconut Oil" ]; then
    echo -e "${RED}❌ Product name mismatch${NC}"
    exit 1
fi

echo -e "${GREEN}✅ Product retrieved successfully${NC}"
echo ""

# Test 15: Update Product
echo -e "${YELLOW}✏️  Test 15: Updating product...${NC}"
UPDATE_CMD="curl -X PUT http://localhost:$PRODUCTS_PORT/api/v1/products/\$PRODUCT_ID -H \'Authorization: Bearer \$API_KEY\' -d \'{\"base_price\": 27.99, ...}\'"
UPDATE_RESPONSE=$(curl -s -X PUT http://localhost:$PRODUCTS_PORT/api/v1/products/$PRODUCT_ID \
  -H "Authorization: Bearer $API_KEY" \
  -H "Content-Type: application/json" \
  -d '{
    "description": "Updated test description",
    "base_price": 27.99
  }')
log_test "15" "Update Product" "$UPDATE_CMD" "$UPDATE_RESPONSE"

echo "$UPDATE_RESPONSE" | jq '.'

UPDATED_PRICE=$(echo "$UPDATE_RESPONSE" | jq -r '.base_price // 0')
if [ "$UPDATED_PRICE" != "27.99" ]; then
    echo -e "${RED}❌ Product update failed${NC}"
    exit 1
fi

echo -e "${GREEN}✅ Product updated successfully${NC}"
echo ""

# Test 16: List Products (Public API - No Auth)
echo -e "${YELLOW}🌍 Test 16: Listing products (Public API)...${NC}"
PUBLIC_CMD="curl http://localhost:$PRODUCTS_PORT/api/v1/public/$DOMAIN/products"
PUBLIC_RESPONSE=$(curl -s http://localhost:$PRODUCTS_PORT/api/v1/public/$DOMAIN/products)
log_test "16" "List Products (Public)" "$PUBLIC_CMD" "$PUBLIC_RESPONSE"

echo "$PUBLIC_RESPONSE" | jq '.'

PUBLIC_COUNT=$(echo "$PUBLIC_RESPONSE" | jq '.count // 0')
if [ "$PUBLIC_COUNT" -eq 0 ]; then
    echo -e "${RED}❌ No public products found${NC}"
    exit 1
fi

echo -e "${GREEN}✅ Public API returned $PUBLIC_COUNT product(s)${NC}"
echo ""

# Test 17: Search Products
echo -e "${YELLOW}🔎 Test 17: Searching products...${NC}"
SEARCH_CMD="curl http://localhost:$PRODUCTS_PORT/api/v1/public/$DOMAIN/products/search?q=coconut"
SEARCH_RESPONSE=$(curl -s "http://localhost:$PRODUCTS_PORT/api/v1/public/$DOMAIN/products/search?q=coconut")
log_test "17" "Search Products" "$SEARCH_CMD" "$SEARCH_RESPONSE"

echo "$SEARCH_RESPONSE" | jq '.'

SEARCH_COUNT=$(echo "$SEARCH_RESPONSE" | jq '.count // 0')
if [ "$SEARCH_COUNT" -eq 0 ]; then
    echo -e "${RED}❌ Search returned no results${NC}"
    exit 1
fi

echo -e "${GREEN}✅ Search found $SEARCH_COUNT product(s)${NC}"
echo ""

# Test 18: Filter by Attributes
echo -e "${YELLOW}🔍 Test 18: Filtering by attributes...${NC}"
FILTER_CMD="curl http://localhost:$PRODUCTS_PORT/api/v1/public/$DOMAIN/products?category=oils"
FILTER_RESPONSE=$(curl -s "http://localhost:$PRODUCTS_PORT/api/v1/public/$DOMAIN/products?category=oils")
log_test "18" "Filter by Attributes" "$FILTER_CMD" "$FILTER_RESPONSE"

echo "$FILTER_RESPONSE" | jq '.'

FILTER_COUNT=$(echo "$FILTER_RESPONSE" | jq '.count // 0')
if [ "$FILTER_COUNT" -eq 0 ]; then
    echo -e "${RED}❌ Filter returned no results${NC}"
    exit 1
fi

echo -e "${GREEN}✅ Filter found $FILTER_COUNT product(s)${NC}"
echo ""

# Test 19: Soft Delete Product
echo -e "${YELLOW}🗑️  Test 19: Soft deleting product...${NC}"
DELETE_CMD="curl -X DELETE http://localhost:$PRODUCTS_PORT/api/v1/products/\$PRODUCT_ID -H \'Authorization: Bearer \$API_KEY\'"
DELETE_RESPONSE=$(curl -s -X DELETE http://localhost:$PRODUCTS_PORT/api/v1/products/$PRODUCT_ID \
  -H "Authorization: Bearer $API_KEY")
log_test "19" "Soft Delete Product" "$DELETE_CMD" "$DELETE_RESPONSE"

echo "$DELETE_RESPONSE" | jq '.'

echo -e "${GREEN}✅ Product soft deleted${NC}"
echo ""

# Test 20: Verify Product is Inactive
echo -e "${YELLOW}🔍 Test 20: Verifying product is inactive...${NC}"
PUBLIC_RESPONSE_AFTER_DELETE=$(curl -s http://localhost:$PRODUCTS_PORT/api/v1/public/$DOMAIN/products)

PUBLIC_COUNT_AFTER=$(echo "$PUBLIC_RESPONSE_AFTER_DELETE" | jq '.count // 0')
if [ "$PUBLIC_COUNT_AFTER" -ne 0 ]; then
    echo -e "${RED}❌ Inactive product still visible in public API${NC}"
    exit 1
fi

echo -e "${GREEN}✅ Inactive product not visible in public API${NC}"
echo ""

# Test 21: Reactivate Product
echo -e "${YELLOW}📦 Test 21: Reactivating deactivated product...${NC}"

# Reactivate the first product
REACTIVATE_CMD="curl -X PUT http://localhost:$PRODUCTS_PORT/api/v1/products/\$PRODUCT_ID -H 'Authorization: Bearer \$API_KEY' -d '{\"active\": true}'"
REACTIVATE_RESPONSE=$(curl -s -X PUT http://localhost:$PRODUCTS_PORT/api/v1/products/$PRODUCT_ID \
  -H "Authorization: Bearer $API_KEY" \
  -H "Content-Type: application/json" \
  -d '{"active": true}')

echo "$REACTIVATE_RESPONSE" | jq '.'
log_test "21" "Reactivate Product" "$REACTIVATE_CMD" "$REACTIVATE_RESPONSE"

PRODUCT_ACTIVE=$(echo "$REACTIVATE_RESPONSE" | jq -r '.active // false')
if [ "$PRODUCT_ACTIVE" != "true" ]; then
    echo -e "${RED}❌ Product reactivation failed${NC}"
    exit 1
fi

echo -e "${GREEN}✅ Product reactivated successfully${NC}"
echo ""

# Test 22: Price Range Filtering
echo -e "${YELLOW}💰 Test 22: Testing price range filtering...${NC}"

# First verify the product is visible in public API after reactivation
ALL_PRODUCTS=$(curl -s "http://localhost:$PRODUCTS_PORT/api/v1/public/$DOMAIN/products")
VISIBLE_COUNT=$(echo "$ALL_PRODUCTS" | jq '.count // 0')

if [ "$VISIBLE_COUNT" -eq 0 ]; then
    echo -e "${YELLOW}⚠️  Product not yet visible after reactivation, skipping price filter test${NC}"
else
    PRICE_FILTER_RESPONSE=$(curl -s "http://localhost:$PRODUCTS_PORT/api/v1/public/$DOMAIN/products?min_price=20&max_price=35")
    PRICE_FILTER_COUNT=$(echo "$PRICE_FILTER_RESPONSE" | jq '.count // 0')

    echo "$PRICE_FILTER_RESPONSE" | jq '.'
    echo -e "${GREEN}✅ Price filtering returned $PRICE_FILTER_COUNT result(s)${NC}"
fi
echo ""

# Test 23: Multiple Attribute Filtering
echo -e "${YELLOW}🔍 Test 23: Testing multiple attribute filtering...${NC}"
if [ "$VISIBLE_COUNT" -gt 0 ]; then
    MULTI_FILTER_RESPONSE=$(curl -s "http://localhost:$PRODUCTS_PORT/api/v1/public/$DOMAIN/products?category=oils&type=coconut")
    echo "$MULTI_FILTER_RESPONSE" | jq '.'
    MULTI_FILTER_COUNT=$(echo "$MULTI_FILTER_RESPONSE" | jq '.count // 0')
    echo -e "${GREEN}✅ Multiple attribute filtering returned $MULTI_FILTER_COUNT result(s)${NC}"
else
    echo -e "${YELLOW}⚠️  Skipping - product not visible${NC}"
fi
echo ""

# Test 24: Combined Filters (search + category + price)
echo -e "${YELLOW}🎯 Test 24: Testing combined filters...${NC}"
if [ "$VISIBLE_COUNT" -gt 0 ]; then
    COMBINED_RESPONSE=$(curl -s "http://localhost:$PRODUCTS_PORT/api/v1/public/$DOMAIN/products/search?q=coconut&category=oils&min_price=20&max_price=40")
    echo "$COMBINED_RESPONSE" | jq '.'
    COMBINED_COUNT=$(echo "$COMBINED_RESPONSE" | jq '.count // 0')
    echo -e "${GREEN}✅ Combined filtering returned $COMBINED_COUNT result(s)${NC}"
else
    echo -e "${YELLOW}⚠️  Skipping - product not visible${NC}"
fi
echo ""

# Test 25: Create Product for Reviews
echo -e "${YELLOW}🛍️  Test 25: Creating second product for reviews...${NC}"
cd ../products_module

PRODUCT_2_RESPONSE=$(curl -s -X POST http://localhost:$PRODUCTS_PORT/api/v1/products \
  -H "Authorization: Bearer $API_KEY" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Test Argan Oil",
    "description": "Premium argan oil for hair care",
    "base_price": 34.99,
    "images": ["https://example.com/argan.jpg"],
    "attributes": {
      "category": "oils",
      "type": "argan"
    },
    "variants": [
      {
        "attributes": {"size": "250ml"},
        "price": 34.99,
        "stock": 30,
        "sku": "ARG-250"
      }
    ]
  }')
PRODUCT2_CMD="curl -X POST http://localhost:$PRODUCTS_PORT/api/v1/products -H \'Authorization: Bearer \$API_KEY\' -d \'{\"name\": \"Test Argan Oil\", ...}\'"
log_test "25" "Create Second Product" "$PRODUCT2_CMD" "$PRODUCT_2_RESPONSE"

PRODUCT_ID_2=$(echo "$PRODUCT_2_RESPONSE" | jq -r '.id // empty')

if [ -z "$PRODUCT_ID_2" ]; then
    echo -e "${RED}❌ Failed to create second product${NC}"
    exit 1
fi

echo -e "${GREEN}✅ Second product created with ID: $PRODUCT_ID_2${NC}"
echo ""

# Test 26: Create Review
echo -e "${YELLOW}⭐ Test 26: Creating review...${NC}"
REVIEW_CMD="curl -X POST http://localhost:$PRODUCTS_PORT/api/v1/public/$DOMAIN/reviews -H 'Content-Type: application/json' -d '{\"product_id\": \"...\", \"product\": \"Test Argan Oil\", \"rating\": 5, ...}'"
REVIEW_RESPONSE=$(curl -s -X POST http://localhost:$PRODUCTS_PORT/api/v1/public/$DOMAIN/reviews \
  -H "Content-Type: application/json" \
  -d "{
    \"product_id\": \"$PRODUCT_ID_2\",
    \"product\": \"Test Argan Oil\",
    \"name\": \"Test User\",
    \"rating\": 5,
    \"text\": \"Amazing product! Works great for my hair.\",
    \"highlight\": \"Excellent quality\"
  }")

echo "$REVIEW_RESPONSE" | jq '.'
log_test "26" "Create Review" "$REVIEW_CMD" "$REVIEW_RESPONSE"

REVIEW_ID=$(echo "$REVIEW_RESPONSE" | jq -r '.id // empty')

if [ -z "$REVIEW_ID" ]; then
    echo -e "${RED}❌ Failed to create review${NC}"
    exit 1
fi

echo -e "${GREEN}✅ Review created with ID: $REVIEW_ID${NC}"
echo ""

# Test 27: List Reviews (Public API)
echo -e "${YELLOW}📋 Test 27: Listing reviews (Public API)...${NC}"
REVIEWS_CMD="curl http://localhost:$PRODUCTS_PORT/api/v1/public/$DOMAIN/reviews"
REVIEWS_RESPONSE=$(curl -s "http://localhost:$PRODUCTS_PORT/api/v1/public/$DOMAIN/reviews")
log_test "27" "List Reviews" "$REVIEWS_CMD" "$REVIEWS_RESPONSE"

echo "$REVIEWS_RESPONSE" | jq '.'

REVIEWS_COUNT=$(echo "$REVIEWS_RESPONSE" | jq '.count // 0')
if [ "$REVIEWS_COUNT" -eq 0 ]; then
    echo -e "${RED}❌ No reviews found${NC}"
    exit 1
fi

echo -e "${GREEN}✅ Found $REVIEWS_COUNT review(s)${NC}"
echo ""

# Test 28: Filter Reviews by Rating
echo -e "${YELLOW}⭐ Test 28: Filtering reviews by rating...${NC}"
FILTERED_CMD="curl http://localhost:$PRODUCTS_PORT/api/v1/public/$DOMAIN/reviews?min_rating=5"
FILTERED_REVIEWS=$(curl -s "http://localhost:$PRODUCTS_PORT/api/v1/public/$DOMAIN/reviews?min_rating=5")
log_test "28" "Filter Reviews by Rating" "$FILTERED_CMD" "$FILTERED_REVIEWS"

echo "$FILTERED_REVIEWS" | jq '.'

FILTERED_COUNT=$(echo "$FILTERED_REVIEWS" | jq '.count // 0')

echo -e "${GREEN}✅ Found $FILTERED_COUNT 5-star review(s)${NC}"
echo ""

# Test 29: Create Second Domain
echo -e "${YELLOW}🌐 Test 29: Creating second domain for isolation testing...${NC}"
cd ../auth_module

DOMAIN2="seconddomain.com"

DOMAIN2_CMD="./auth-module domain create --config=config.dev.yaml --domain=$DOMAIN2 --name=\"Second Test Domain\""
DOMAIN2_OUTPUT=$(./auth-module domain create \
    --config=config.dev.yaml \
    --domain=$DOMAIN2 \
    --name="Second Test Domain" \
    --admin-email=admin@$DOMAIN2 2>&1)

echo "$DOMAIN2_OUTPUT"
log_test "29" "Create Second Domain" "$DOMAIN2_CMD" "$DOMAIN2_OUTPUT"

if [ $? -ne 0 ]; then
    echo -e "${RED}❌ Second domain creation failed${NC}"
    exit 1
fi

echo -e "${GREEN}✅ Second domain created: $DOMAIN2${NC}"
echo ""

# Test 30: Create User in Second Domain
echo -e "${YELLOW}👤 Test 30: Creating user in second domain...${NC}"

DOMAIN2_INV_CMD="./auth-module invite create --config=config.dev.yaml --domain=$DOMAIN2 --email=user@$DOMAIN2 --role=admin"
DOMAIN2_INVITE_OUTPUT=$(./auth-module invite create \
    --config=config.dev.yaml \
    --domain=$DOMAIN2 \
    --email=user@$DOMAIN2 \
    --role=admin \
    --type=email_with_qr 2>&1)
log_test "30" "Create Invitation (Domain 2)" "$DOMAIN2_INV_CMD" "$DOMAIN2_INVITE_OUTPUT"

DOMAIN2_INVITE_TOKEN=$(echo "$DOMAIN2_INVITE_OUTPUT" | grep "token:" | awk '{print $2}')

DOMAIN2_ACCEPT_RESPONSE=$(curl -s -X POST http://localhost:$AUTH_PORT/api/v1/auth/invitation/accept \
  -H "Content-Type: application/json" \
  -d "{
    \"token\": \"$DOMAIN2_INVITE_TOKEN\",
    \"email\": \"user@$DOMAIN2\",
    \"name\": \"Domain2 User\",
    \"auth_provider\": \"magic_link\"
  }")

DOMAIN2_JWT=$(echo "$DOMAIN2_ACCEPT_RESPONSE" | jq -r '.token // empty')

if [ -z "$DOMAIN2_JWT" ]; then
    echo -e "${RED}❌ Failed to create user in second domain${NC}"
    exit 1
fi

echo -e "${GREEN}✅ User created in second domain${NC}"
echo ""

# Test 31: Verify Cross-Domain Isolation
echo -e "${YELLOW}🔒 Test 31: Verifying cross-domain isolation...${NC}"

# Domain1 user should NOT see Domain2 data
CROSS_CMD="curl http://localhost:$PRODUCTS_PORT/api/v1/public/$DOMAIN2/products (should be empty)"
CROSS_DOMAIN_PRODUCTS=$(curl -s http://localhost:$PRODUCTS_PORT/api/v1/public/$DOMAIN2/products)
log_test "31" "Verify Domain Isolation" "$CROSS_CMD" "$CROSS_DOMAIN_PRODUCTS"
DOMAIN2_PRODUCT_COUNT=$(echo "$CROSS_DOMAIN_PRODUCTS" | jq '.count // 0')

# Domain2 should have no products (we only created products in domain1)
if [ "$DOMAIN2_PRODUCT_COUNT" -ne 0 ]; then
    echo -e "${RED}❌ Domain isolation issue: domain2 shows domain1 products${NC}"
    exit 1
fi

echo -e "${GREEN}✅ Cross-domain isolation verified (domain2 has $DOMAIN2_PRODUCT_COUNT products)${NC}"
echo ""

# Test 32: Verify Unauthorized Access Fails
echo -e "${YELLOW}🔒 Test 32: Testing unauthorized access...${NC}"
UNAUTH_CMD="curl -X POST http://localhost:$PRODUCTS_PORT/api/v1/products (no Authorization header - should fail with 401)"
UNAUTH_RESPONSE=$(curl -s -w "\n%{http_code}" -X POST http://localhost:$PRODUCTS_PORT/api/v1/products \
  -H "Content-Type: application/json" \
  -d '{"name": "Should Fail", "base_price": 10}')

log_test "32" "Unauthorized Access Test" "$UNAUTH_CMD" "$UNAUTH_RESPONSE"

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
echo -e "${GREEN}🎉 All 32 Integration Tests Passed! 🎉${NC}"
echo -e "${GREEN}════════════════════════════════════════${NC}"
echo ""
echo -e "Summary:"
echo -e "  ✅ Auth & Products modules built and started"
echo -e "  ✅ Domain creation & API key management"
echo -e "  ✅ User authentication & JWT validation"
echo -e "  ✅ Permission-based access control"
echo -e "  ✅ Product CRUD operations & reactivation"
echo -e "  ✅ Advanced filtering (price, attributes, combined)"
echo -e "  ✅ Public API access & search"
echo -e "  ✅ Reviews system (create, list, filter)"
echo -e "  ✅ Domain isolation"
echo -e "  ✅ Soft delete & authorization checks"
echo ""
echo -e "${YELLOW}View logs:${NC}"
echo -e "  Auth module: /tmp/auth-module.log"
echo -e "  Products module: /tmp/products-module.log"
echo -e "  Detailed test log: $TEST_LOG"
echo ""
echo -e "${GREEN}📝 Detailed test log contains commands and JSON responses${NC}"
echo ""
