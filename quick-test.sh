#!/bin/bash

set -e

echo "🧪 Quick Integration Test"
echo ""

# Colors
GREEN='\033[0;32m'
RED='\033[0;31m'
NC='\033[0m'

# Test health
echo -n "Testing API health... "
if curl -sf http://localhost:9090/health > /dev/null; then
    echo -e "${GREEN}✓${NC}"
else
    echo -e "${RED}✗${NC}"
    exit 1
fi

# Test test-domain
echo -n "Testing test-domain... "
if curl -sf http://localhost:8000 > /dev/null; then
    echo -e "${GREEN}✓${NC}"
else
    echo -e "${RED}✗${NC}"
    exit 1
fi

# Test auth-ui
echo -n "Testing auth-ui... "
if curl -sf http://localhost:5173 > /dev/null; then
    echo -e "${GREEN}✓${NC}"
else
    echo -e "${RED}✗${NC}"
    exit 1
fi

# Create test invitation
echo -n "Creating test invitation... "
RANDOM_TOKEN="quick-test-$(date +%s)"
docker exec auth-module-mongodb mongosh auth_module --quiet --eval "
db.invitations.insertOne({
  domain: 'testdomain.com',
  token: '$RANDOM_TOKEN',
  email: 'quicktest@example.com',
  role: 'customer',
  permissions: ['products.read', 'orders.read'],
  type: 'email_with_qr',
  single_use: true,
  uses_count: 0,
  created_by: 'system',
  created_at: new Date(),
  expires_at: new Date(Date.now() + 24*60*60*1000),
  status: 'pending'
})
" > /dev/null && echo -e "${GREEN}✓${NC}" || echo -e "${RED}✗${NC}"

# Accept invitation
echo -n "Accepting invitation... "
JWT=$(curl -sf -X POST "http://localhost:9090/api/v1/auth/invitation/accept" \
  -H "Content-Type: application/json" \
  -d "{
    \"token\": \"$RANDOM_TOKEN\",
    \"email\": \"quicktest@example.com\",
    \"name\": \"Quick Test\",
    \"auth_provider\": \"magic_link\"
  }" | jq -r '.token')

if [ -n "$JWT" ] && [ "$JWT" != "null" ]; then
    echo -e "${GREEN}✓${NC}"
else
    echo -e "${RED}✗${NC}"
    exit 1
fi

# Test JWT
echo -n "Testing JWT with /auth/me... "
USER=$(curl -sf "http://localhost:9090/api/v1/auth/me" \
  -H "Authorization: Bearer $JWT" \
  -H "Host: testdomain.com" | jq -r '.email')

if [ "$USER" = "quicktest@example.com" ]; then
    echo -e "${GREEN}✓${NC}"
else
    echo -e "${RED}✗${NC}"
    exit 1
fi

echo ""
echo -e "${GREEN}✅ All tests passed!${NC}"
echo ""
echo "User created: quicktest@example.com"
echo "JWT Token saved to: /tmp/quick_test_jwt.txt"
echo ""
echo "$JWT" > /tmp/quick_test_jwt.txt
echo "Test in browser:"
echo "  http://localhost:8000#token=$JWT"
echo ""
