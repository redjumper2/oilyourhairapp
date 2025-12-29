# Manual Testing Guide

Complete step-by-step guide to manually test the auth module integration.

## Prerequisites

### 1. Start All Services

```bash
cd /home/sparque/dev/oilyourhairapp

# Start all services
cd auth_module
docker compose up -d

# Wait for services to be healthy (~30 seconds)
docker ps

# Start auth-ui and test-domain
cd ..
docker compose up -d auth-ui test-domain
```

### 2. Verify Services Are Running

```bash
# Check all containers are running
docker ps

# Expected output:
# - auth-module-api (port 9090)
# - auth-module-mongodb (port 27017)
# - auth-ui (port 5173)
# - test-domain (port 8000)

# Test each service
curl http://localhost:9090/health     # Should return: {"status":"healthy"}
curl http://localhost:8000            # Should return HTML
curl http://localhost:5173            # Should return HTML
```

### 3. Initialize Test Domain

```bash
# Create testdomain.com and admin invitation
./init-test-domain.sh

# This creates:
# - Domain: testdomain.com
# - Admin user invitation
# - Outputs the invitation token
```

---

## Test 1: Invitation Flow (Customer User)

### Step 1: Create a Customer Invitation

Since we don't have an admin JWT yet, we'll create an invitation directly in MongoDB:

```bash
# Connect to MongoDB and create customer invitation
docker exec auth-module-mongodb mongosh auth_module --eval "
db.invitations.insertOne({
  domain: 'testdomain.com',
  token: 'test-customer-token-123',
  email: 'customer@example.com',
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
"
```

### Step 2: Verify the Invitation

```bash
# Verify invitation is valid
curl "http://localhost:9090/api/v1/auth/invitation/verify?token=test-customer-token-123" | jq .

# Expected output:
# {
#   "branding": {
#     "company_name": "Test Domain Store",
#     "logo_url": "",
#     "primary_color": "#000000"
#   },
#   "domain": "testdomain.com",
#   "email": "customer@example.com",
#   "role": "customer",
#   "expires_at": "...",
#   "time_remaining": "23h59m..."
# }
```

### Step 3: Accept the Invitation

```bash
# Accept invitation and create user account
curl -X POST "http://localhost:9090/api/v1/auth/invitation/accept" \
  -H "Content-Type: application/json" \
  -d '{
    "token": "test-customer-token-123",
    "email": "customer@example.com",
    "name": "Test Customer",
    "auth_provider": "magic_link"
  }' | jq .

# Expected output:
# {
#   "token": "eyJhbGciOiJI...",  (JWT token)
#   "user": {
#     "id": "...",
#     "email": "customer@example.com",
#     "domain": "testdomain.com",
#     "role": "customer",
#     "permissions": ["products.read", "orders.read"]
#   }
# }
```

### Step 4: Save the JWT Token

```bash
# Copy the JWT from the response above and save it
export JWT="eyJhbGciOiJI..."  # Paste the actual token here

# Or save to file
echo "eyJhbGciOiJI..." > /tmp/customer_jwt.txt
```

### Step 5: Test the JWT

```bash
# Verify JWT works with /auth/me endpoint
curl "http://localhost:9090/api/v1/auth/me" \
  -H "Authorization: Bearer $JWT" \
  -H "Host: testdomain.com" | jq .

# Expected output:
# {
#   "id": "...",
#   "email": "customer@example.com",
#   "domain": "testdomain.com",
#   "role": "customer",
#   "permissions": ["products.read", "orders.read"]
# }
```

### Step 6: Test Browser Integration

1. **Open test domain in browser:**
   ```
   http://localhost:8000
   ```
   - You should see "Login to view price" buttons
   - Products should be visible but prices hidden

2. **Simulate auth redirect:**
   - Open browser console (F12)
   - Paste the JWT from Step 4
   - Navigate to:
   ```
   http://localhost:8000#token=YOUR_JWT_HERE
   ```

3. **Verify authentication works:**
   - Page should reload
   - JWT should be stored in localStorage
   - User email should appear in header
   - Product prices should now be visible
   - Login button should change to Logout button

4. **Test logout:**
   - Click "Logout" button
   - JWT should be cleared from localStorage
   - Prices should be hidden again
   - Login button should reappear

---

## Test 2: Complete UI Flow (With Auth Portal)

### Step 1: Create Another Customer Invitation

```bash
# Create a new invitation with a different email
docker exec auth-module-mongodb mongosh auth_module --eval "
db.invitations.insertOne({
  domain: 'testdomain.com',
  token: 'test-customer-token-456',
  email: 'alice@example.com',
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
"
```

### Step 2: Test Invitation Acceptance via UI

1. **Open invitation URL in browser:**
   ```
   http://localhost:5173/invite?token=test-customer-token-456&redirect=http://localhost:8000
   ```

2. **Fill in the invitation form:**
   - Email should be pre-filled: `alice@example.com`
   - Enter name: `Alice Smith`
   - Choose auth provider: `magic_link`

3. **Submit the form:**
   - User account should be created
   - You'll receive a JWT
   - Browser should redirect to: `http://localhost:8000#token=JWT_HERE`

4. **Verify on test-domain:**
   - Should be automatically logged in
   - Prices should be visible
   - User email should show in header

---

## Test 3: Magic Link Flow (Requires SMTP)

**Note:** Magic link requires SMTP to be configured. If SMTP is not configured, you can still test by manually retrieving tokens from MongoDB.

### Step 1: Configure SMTP (Optional)

Edit `/home/sparque/dev/oilyourhairapp/.env`:

```bash
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_USER=your-email@gmail.com
SMTP_PASSWORD=your-app-password
EMAIL_FROM=noreply@testdomain.com
```

Restart services:
```bash
docker compose restart
```

### Step 2: Request Magic Link

**Via API:**
```bash
curl -X POST "http://localhost:9090/api/v1/auth/magic-link/request" \
  -H "Content-Type: application/json" \
  -H "Host: testdomain.com" \
  -d '{
    "email": "customer@example.com",
    "domain": "testdomain.com"
  }'
```

**Via Browser:**
1. Visit: `http://localhost:5173/login?domain=testdomain.com&redirect=http://localhost:8000`
2. Enter email: `customer@example.com`
3. Click "Send Magic Link"

### Step 3: Get Magic Link Token

**If SMTP is configured:**
- Check your email for the magic link

**If SMTP is NOT configured:**
```bash
# Get the latest magic link token from MongoDB
docker exec auth-module-mongodb mongosh auth_module --quiet --eval "
db.magic_link_tokens.find().sort({created_at:-1}).limit(1).pretty()
"

# Copy the 'token' field
```

### Step 4: Verify Magic Link

**Via Browser:**
```
http://localhost:5173/verify?token=TOKEN_HERE&redirect=http://localhost:8000
```

**Via API:**
```bash
curl "http://localhost:9090/api/v1/auth/magic-link/verify?token=TOKEN_HERE" | jq .
```

Should return JWT token and redirect to test-domain with token in hash.

---

## Test 4: Admin Operations

### Step 1: Accept Admin Invitation

```bash
# Get admin invitation token from init-test-domain.sh output
# Or retrieve from MongoDB:
docker exec auth-module-mongodb mongosh auth_module --quiet --eval "
db.invitations.findOne({role: 'admin', domain: 'testdomain.com'})
"

# Accept admin invitation
curl -X POST "http://localhost:9090/api/v1/auth/invitation/accept" \
  -H "Content-Type: application/json" \
  -d '{
    "token": "ADMIN_TOKEN_HERE",
    "email": "admin@testdomain.com",
    "name": "Admin User",
    "auth_provider": "magic_link"
  }' | jq .

# Save the JWT
export ADMIN_JWT="eyJhbGciOiJI..."
```

### Step 2: Create Customer Invitation via API

```bash
# Create invitation using admin JWT
curl -X POST "http://localhost:9090/api/v1/admin/users/invite" \
  -H "Authorization: Bearer $ADMIN_JWT" \
  -H "Host: testdomain.com" \
  -H "Content-Type: application/json" \
  -d '{
    "email": "bob@example.com",
    "role": "customer",
    "type": "email_with_qr"
  }' | jq .

# Response includes:
# - invitation_token
# - invitation_url
# - qr_code (base64 PNG)
```

### Step 3: Test Admin Endpoints

```bash
# List all users
curl "http://localhost:9090/api/v1/admin/users" \
  -H "Authorization: Bearer $ADMIN_JWT" \
  -H "Host: testdomain.com" | jq .

# Get domain settings
curl "http://localhost:9090/api/v1/admin/domain/settings" \
  -H "Authorization: Bearer $ADMIN_JWT" \
  -H "Host: testdomain.com" | jq .

# Get all permissions
curl "http://localhost:9090/api/v1/admin/permissions" \
  -H "Authorization: Bearer $ADMIN_JWT" \
  -H "Host: testdomain.com" | jq .

# Update domain branding
curl -X PUT "http://localhost:9090/api/v1/admin/domain/settings" \
  -H "Authorization: Bearer $ADMIN_JWT" \
  -H "Host: testdomain.com" \
  -H "Content-Type: application/json" \
  -d '{
    "branding": {
      "company_name": "My Awesome Store",
      "logo_url": "https://example.com/logo.png",
      "primary_color": "#FF5733"
    }
  }' | jq .
```

---

## Verification Checklist

### Services Health

- [ ] `curl http://localhost:9090/health` returns `{"status":"healthy"}`
- [ ] `curl http://localhost:8000` returns test-domain HTML
- [ ] `curl http://localhost:5173` returns auth-ui HTML
- [ ] `docker ps` shows all 4 containers running

### Database

```bash
# List all domains
docker exec auth-module-mongodb mongosh auth_module --quiet --eval "
db.domains.find({}, {domain: 1, name: 1}).pretty()
"

# List all users
docker exec auth-module-mongodb mongosh auth_module --quiet --eval "
db.users.find({}, {email: 1, role: 1, domain: 1}).pretty()
"

# List all invitations
docker exec auth-module-mongodb mongosh auth_module --quiet --eval "
db.invitations.find({}, {email: 1, role: 1, status: 1, token: 1}).pretty()
"
```

### API Endpoints

Test each endpoint:

- [ ] `POST /api/v1/auth/invitation/accept` - Create user from invitation
- [ ] `GET /api/v1/auth/invitation/verify` - Verify invitation validity
- [ ] `GET /api/v1/auth/me` - Get current user (with JWT)
- [ ] `POST /api/v1/auth/magic-link/request` - Request magic link
- [ ] `GET /api/v1/auth/magic-link/verify` - Verify magic link
- [ ] `GET /api/v1/admin/users` - List users (admin only)
- [ ] `POST /api/v1/admin/users/invite` - Create invitation (admin only)
- [ ] `GET /api/v1/admin/permissions` - List permissions

### Browser Testing

- [ ] test-domain shows "Login to view price" when not authenticated
- [ ] Clicking Login redirects to auth portal with correct domain parameter
- [ ] Auth portal shows correct branding (Test Domain Store)
- [ ] Accepting invitation creates user and redirects back with JWT
- [ ] test-domain extracts JWT from hash and stores in localStorage
- [ ] test-domain calls /auth/me with JWT
- [ ] Product prices become visible after authentication
- [ ] User email appears in header
- [ ] Logout clears JWT and hides prices

---

## Troubleshooting

### Services Won't Start

```bash
# Check what's using the ports
sudo lsof -i :9090
sudo lsof -i :27017
sudo lsof -i :5173
sudo lsof -i :8000

# Kill processes if needed
sudo kill <PID>

# Check logs
docker logs auth-module-api
docker logs auth-ui
docker logs test-domain
docker logs auth-module-mongodb
```

### JWT Not Working

```bash
# Decode JWT to check contents (without verification)
echo "YOUR_JWT" | cut -d. -f2 | base64 -d 2>/dev/null | jq .

# Check expiration time
# The 'exp' field should be in the future (Unix timestamp)

# Verify JWT_SECRET matches
docker exec auth-module-api env | grep JWT_SECRET
```

### Database Issues

```bash
# Connect to MongoDB
docker exec -it auth-module-mongodb mongosh auth_module

# Useful queries
db.users.countDocuments()
db.invitations.countDocuments()
db.domains.countDocuments()
db.magic_link_tokens.countDocuments()

# Clear all data (careful!)
db.users.deleteMany({})
db.invitations.deleteMany({})
db.magic_link_tokens.deleteMany({})
```

### CORS Issues

Check browser console. If you see CORS errors:

```bash
# Check API CORS configuration
docker exec auth-module-api cat /app/config.yaml | grep -A 5 cors

# Restart with updated config
docker compose restart auth-module-api
```

### Test Domain Permission Denied

```bash
# Fix file permissions
chmod -R 755 /home/sparque/dev/oilyourhairapp/test_domain
```

---

## Quick Test Script

Save this as `quick-test.sh`:

```bash
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
docker exec auth-module-mongodb mongosh auth_module --quiet --eval "
db.invitations.insertOne({
  domain: 'testdomain.com',
  token: 'quick-test-token',
  email: 'quicktest@example.com',
  role: 'customer',
  permissions: ['products.read'],
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
  -d '{
    "token": "quick-test-token",
    "email": "quicktest@example.com",
    "name": "Quick Test",
    "auth_provider": "magic_link"
  }' | jq -r '.token')

if [ -n "$JWT" ]; then
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
echo "JWT Token: $JWT"
echo ""
echo "Test in browser:"
echo "  http://localhost:8000#token=$JWT"
```

Run it:
```bash
chmod +x quick-test.sh
./quick-test.sh
```

---

## Next Steps

Once manual testing is complete, consider:

1. **Set up proper SMTP** for magic links in production
2. **Configure Google OAuth** for social login
3. **Add more domains** to test multi-tenancy
4. **Create automated tests** based on these manual steps
5. **Test with real customer scenarios**
6. **Security audit** of JWT handling and CORS
