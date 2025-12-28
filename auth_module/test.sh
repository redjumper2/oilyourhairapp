#!/bin/bash

# Auth Module Integration Test Script
# Tests basic functionality of the auth module

set -e

API_URL="http://localhost:8080/api/v1"
DOMAIN="testdomain.com"
GREEN='\033[0;32m'
RED='\033[0;31m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${BLUE}🧪 Auth Module Integration Test${NC}"
echo ""

# Test 1: Health Check
echo -e "${BLUE}1️⃣ Testing health endpoint...${NC}"
HEALTH=$(curl -s http://localhost:8080/health | jq -r '.status')
if [ "$HEALTH" = "healthy" ]; then
    echo -e "${GREEN}✅ Health check passed${NC}"
else
    echo -e "${RED}❌ Health check failed${NC}"
    exit 1
fi
echo ""

# Test 2: Request Magic Link
echo -e "${BLUE}2️⃣ Testing magic link request...${NC}"
RESPONSE=$(curl -s -X POST $API_URL/auth/magic-link/request \
  -H "Host: $DOMAIN" \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com"}')

MESSAGE=$(echo $RESPONSE | jq -r '.message // .error')
if [[ $MESSAGE == *"Magic link sent"* ]]; then
    echo -e "${GREEN}✅ Magic link request passed${NC}"
elif [[ $MESSAGE == *"Domain not found"* ]]; then
    echo -e "${RED}❌ Domain not found. Create domain first:${NC}"
    echo -e "   make domain-create DOMAIN=$DOMAIN NAME=\"Test\" EMAIL=admin@$DOMAIN"
    exit 1
else
    echo -e "${RED}❌ Magic link request failed: $MESSAGE${NC}"
    exit 1
fi
echo ""

# Test 3: Test domain isolation (should fail)
echo -e "${BLUE}3️⃣ Testing domain isolation...${NC}"
RESPONSE=$(curl -s -X POST $API_URL/auth/magic-link/request \
  -H "Host: nonexistent.com" \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com"}')

ERROR=$(echo $RESPONSE | jq -r '.error')
if [[ $ERROR == *"not found"* ]] || [[ $ERROR == *"inactive"* ]]; then
    echo -e "${GREEN}✅ Domain isolation working (rejected invalid domain)${NC}"
else
    echo -e "${RED}❌ Domain isolation failed${NC}"
    exit 1
fi
echo ""

# Test 4: Test without auth (should fail)
echo -e "${BLUE}4️⃣ Testing protected endpoint without auth...${NC}"
RESPONSE=$(curl -s $API_URL/auth/me -H "Host: $DOMAIN")
ERROR=$(echo $RESPONSE | jq -r '.error')

if [[ $ERROR == *"Missing authorization"* ]] || [[ $ERROR == *"Unauthorized"* ]]; then
    echo -e "${GREEN}✅ Auth protection working (rejected unauthenticated request)${NC}"
else
    echo -e "${RED}❌ Auth protection failed${NC}"
    exit 1
fi
echo ""

# Test 5: Test admin endpoint without auth (should fail)
echo -e "${BLUE}5️⃣ Testing admin endpoint without auth...${NC}"
RESPONSE=$(curl -s $API_URL/admin/users -H "Host: $DOMAIN")
ERROR=$(echo $RESPONSE | jq -r '.error')

if [[ $ERROR == *"Missing authorization"* ]] || [[ $ERROR == *"Unauthorized"* ]]; then
    echo -e "${GREEN}✅ Admin protection working (rejected unauthenticated request)${NC}"
else
    echo -e "${RED}❌ Admin protection failed${NC}"
    exit 1
fi
echo ""

echo -e "${GREEN}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo -e "${GREEN}✅ All automated tests passed!${NC}"
echo -e "${GREEN}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo ""
echo -e "${BLUE}📋 Manual tests remaining:${NC}"
echo "1. Accept invitation and get JWT token"
echo "2. Test authenticated endpoints with JWT"
echo "3. Test admin APIs with admin JWT"
echo "4. Test invitation creation with QR codes"
echo ""
echo "See TESTING.md for detailed manual testing steps"
