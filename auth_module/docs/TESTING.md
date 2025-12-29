# Testing Guide

Complete guide to testing the auth module from scratch.

## Prerequisites

- Docker & Docker Compose installed
- `curl` or Postman
- Optional: `jq` for pretty JSON output (`sudo apt install jq`)

## Part 1: Setup & Start Services

### Step 1: Initial Setup

```bash
cd auth_module

# Create config files
make setup

# Edit .env and set a strong JWT secret
nano .env
```

**Required in `.env`:**
```bash
JWT_SECRET=my-super-secret-key-for-testing-12345
```

**Optional (for email testing):**
```bash
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_USER=your-email@gmail.com
SMTP_PASSWORD=your-app-password
EMAIL_FROM=noreply@yourdomain.com
```

### Step 2: Start Services

```bash
# Start MongoDB + API
make docker-up

# Wait for "✅ API is ready"
# This may take 30-60 seconds on first run
```

**Verify services are running:**
```bash
make status

# Should show:
# auth-module-api      Up
# auth-module-mongodb  Up
```

### Step 3: Health Check

```bash
make health

# Or:
curl http://localhost:8080/health | jq

# Expected:
# {
#   "status": "healthy",
#   "env": "development"
# }
```

✅ **If you see this, the system is running!**

## Part 2: Domain Management

### Create First Domain

```bash
make domain-create \
  DOMAIN=testdomain.com \
  NAME="Test Domain" \
  EMAIL=admin@testdomain.com
```

**Expected output:**
```
✅ Domain created: testdomain.com
✅ Admin invitation created for: admin@testdomain.com
📧 Invitation URL: http://localhost:3000/invite?token=abc123...
📱 QR Code: data:image/png;base64,...
⏰ Expires: 2025-12-29T...
```

**Save the token** from the invitation URL - you'll need it!

### List Domains

```bash
make domain-list

# Expected:
# DOMAIN              NAME            STATUS      CREATED
# testdomain.com      Test Domain     active      2025-12-28 10:00
```

### Create Second Domain (for testing isolation)

```bash
make domain-create \
  DOMAIN=anotherdomain.com \
  NAME="Another Domain" \
  EMAIL=admin@anotherdomain.com
```

## Part 3: Test Magic Link Authentication

### Request Magic Link

```bash
curl -X POST http://localhost:8080/api/v1/auth/magic-link/request \
  -H "Host: testdomain.com" \
  -H "Content-Type: application/json" \
  -d '{"email":"user@example.com"}' | jq

# Expected:
# {
#   "message": "Magic link sent to your email"
# }
```

**Note:** Without SMTP configured, the email won't actually send, but the token is created.

### Get Magic Link Token (for testing)

Since email isn't configured, we'll query MongoDB directly:

```bash
make docker-mongo-shell

# In MongoDB shell:
use auth_module
db.magic_link_tokens.find().sort({created_at: -1}).limit(1).pretty()

# Copy the "token" value
# Type: exit
```

### Verify Magic Link

```bash
# Replace TOKEN_HERE with actual token from MongoDB
curl "http://localhost:8080/api/v1/auth/magic-link/verify?token=TOKEN_HERE" | jq

# Expected:
# {
#   "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
#   "user": {
#     "id": "507f1f77bcf86cd799439011",
#     "email": "user@example.com",
#     "domain": "testdomain.com",
#     "role": "customer",
#     "permissions": ["products.read", "cart.read", "cart.write", "orders.read"]
#   }
# }
```

**Save the JWT token** - you'll use it for authenticated requests!

## Part 4: Test Invitation System

### Verify Invitation (from domain creation)

```bash
# Use the token from domain creation step
curl "http://localhost:8080/api/v1/auth/invitation/verify?token=INVITATION_TOKEN_HERE" | jq

# Expected:
# {
#   "invitation_id": "...",
#   "role": "admin",
#   "domain": "testdomain.com",
#   "expires_at": "2025-12-29T...",
#   "time_remaining": "23h59m",
#   "email": "admin@testdomain.com",
#   "branding": {
#     "company_name": "Test Domain",
#     "primary_color": "#000000",
#     "logo_url": ""
#   }
# }
```

### Accept Invitation

```bash
curl -X POST http://localhost:8080/api/v1/auth/invitation/accept \
  -H "Content-Type: application/json" \
  -d '{
    "token": "INVITATION_TOKEN_HERE",
    "email": "admin@testdomain.com",
    "auth_provider": "magic_link",
    "provider_id": ""
  }' | jq

# Expected:
# {
#   "token": "eyJhbGc...",
#   "user": {
#     "id": "...",
#     "email": "admin@testdomain.com",
#     "domain": "testdomain.com",
#     "role": "admin",
#     "permissions": ["domain.settings.read", "domain.settings.write", ...]
#   }
# }
```

**Save this admin JWT token** - you'll need it for admin API testing!

## Part 5: Test Authenticated Endpoints

### Get Current User

```bash
# Replace JWT_TOKEN with your actual token
export JWT_TOKEN="eyJhbGc..."

curl http://localhost:8080/api/v1/auth/me \
  -H "Authorization: Bearer $JWT_TOKEN" \
  -H "Host: testdomain.com" | jq

# Expected:
# {
#   "id": "...",
#   "email": "admin@testdomain.com",
#   "domain": "testdomain.com",
#   "role": "admin",
#   "permissions": [...]
# }
```

### Test Domain Isolation

Try using the token on a different domain:

```bash
curl http://localhost:8080/api/v1/auth/me \
  -H "Authorization: Bearer $JWT_TOKEN" \
  -H "Host: anotherdomain.com" | jq

# Expected:
# {
#   "error": "Token domain mismatch"
# }
```

✅ **Domain isolation is working!**

## Part 6: Test Admin APIs

Use the admin JWT token for these tests.

### Get Domain Settings

```bash
export ADMIN_JWT="your-admin-jwt-here"

curl http://localhost:8080/api/v1/admin/domain/settings \
  -H "Authorization: Bearer $ADMIN_JWT" \
  -H "Host: testdomain.com" | jq

# Expected:
# {
#   "domain": "testdomain.com",
#   "name": "Test Domain",
#   "status": "active",
#   "settings": {
#     "allowed_auth_providers": ["google", "magic_link"],
#     "default_role": "customer",
#     "require_email_verification": true
#   },
#   "branding": {...}
# }
```

### Update Domain Settings

```bash
curl -X PUT http://localhost:8080/api/v1/admin/domain/settings \
  -H "Authorization: Bearer $ADMIN_JWT" \
  -H "Host: testdomain.com" \
  -H "Content-Type: application/json" \
  -d '{
    "branding": {
      "company_name": "Test Domain Updated",
      "primary_color": "#FF5722",
      "logo_url": "https://example.com/logo.png",
      "support_email": "support@testdomain.com"
    }
  }' | jq

# Expected:
# {
#   "message": "Domain settings updated successfully"
# }
```

### List Users

```bash
curl http://localhost:8080/api/v1/admin/users \
  -H "Authorization: Bearer $ADMIN_JWT" \
  -H "Host: testdomain.com" | jq

# Expected:
# {
#   "users": [
#     {
#       "id": "...",
#       "email": "admin@testdomain.com",
#       "role": "admin",
#       "auth_provider": "magic_link",
#       "created_at": "...",
#       "last_login": "..."
#     }
#   ],
#   "count": 1
# }
```

### Create Invitation (QR Code)

```bash
curl -X POST http://localhost:8080/api/v1/admin/users/invite \
  -H "Authorization: Bearer $ADMIN_JWT" \
  -H "Host: testdomain.com" \
  -H "Content-Type: application/json" \
  -d '{
    "email": "editor@testdomain.com",
    "role": "editor",
    "type": "email_with_qr",
    "single_use": true,
    "expires_in_hours": 48,
    "promo_code": "WELCOME20",
    "discount_percent": 20
  }' | jq

# Expected:
# {
#   "invitation_id": "...",
#   "url": "http://localhost:3000/invite?token=...",
#   "token": "...",
#   "expires_at": "...",
#   "qr_code": "data:image/png;base64,iVBORw0KG..."
# }
```

**The QR code is a base64 data URL!** You can:
- Save it to an HTML file and view in browser
- Use an online base64 image decoder
- Pass it to frontend to display

### Create Multi-Use Promotional QR Code

```bash
curl -X POST http://localhost:8080/api/v1/admin/users/invite \
  -H "Authorization: Bearer $ADMIN_JWT" \
  -H "Host: testdomain.com" \
  -H "Content-Type: application/json" \
  -d '{
    "role": "customer",
    "type": "qr_code",
    "single_use": false,
    "max_uses": 100,
    "expires_in_hours": 720,
    "promo_code": "TRADESHOW2025",
    "source": "booth",
    "discount_percent": 15
  }' | jq

# Expected: Same response with QR code
```

### Update User Role

First, get a user ID from the list users endpoint, then:

```bash
curl -X PUT http://localhost:8080/api/v1/admin/users/USER_ID_HERE \
  -H "Authorization: Bearer $ADMIN_JWT" \
  -H "Host: testdomain.com" \
  -H "Content-Type: application/json" \
  -d '{
    "role": "viewer"
  }' | jq

# Expected:
# {
#   "message": "User updated successfully"
# }
```

### Delete User (Soft Delete)

```bash
curl -X DELETE http://localhost:8080/api/v1/admin/users/USER_ID_HERE \
  -H "Authorization: Bearer $ADMIN_JWT" \
  -H "Host: testdomain.com" | jq

# Expected:
# {
#   "message": "User deleted successfully"
# }
```

**Test admin protection:**

Try to delete yourself:
```bash
# Should fail with: "Cannot delete yourself"
```

### Get All Permissions

View all available permissions grouped by category:

```bash
curl http://localhost:8080/api/v1/admin/permissions \
  -H "Authorization: Bearer $ADMIN_JWT" \
  -H "Host: testdomain.com" | jq

# Expected:
# {
#   "groups": [
#     {
#       "name": "Domain Management",
#       "description": "Manage domain settings and branding",
#       "permissions": ["domain.settings.read", "domain.settings.write"]
#     },
#     {
#       "name": "User Management",
#       "description": "Manage users, roles, and invitations",
#       "permissions": ["users.read", "users.write", "users.delete", "users.invite"]
#     },
#     ...
#   ],
#   "total": 16
# }
```

**Use case:** Build a frontend permission selector with these grouped permissions.

### Get Permissions by Role

View what permissions each role has:

```bash
curl http://localhost:8080/api/v1/admin/permissions/roles \
  -H "Authorization: Bearer $ADMIN_JWT" \
  -H "Host: testdomain.com" | jq

# Expected:
# {
#   "admin": {
#     "count": 12,
#     "permissions": ["domain.settings.read", "domain.settings.write", ...]
#   },
#   "editor": {
#     "count": 5,
#     "permissions": ["products.read", "products.write", ...]
#   },
#   "viewer": {
#     "count": 3,
#     "permissions": ["products.read", "orders.read", "inventory.read"]
#   },
#   "customer": {
#     "count": 4,
#     "permissions": ["products.read", "cart.read", "cart.write", "orders.read"]
#   }
# }
```

**Use case:** Display role comparison table in admin UI.

### Create Invitation with Custom Permissions

Override default role permissions with custom selection:

```bash
curl -X POST http://localhost:8080/api/v1/admin/users/invite \
  -H "Authorization: Bearer $ADMIN_JWT" \
  -H "Host: testdomain.com" \
  -H "Content-Type: application/json" \
  -d '{
    "email": "custom@testdomain.com",
    "role": "viewer",
    "type": "email_with_qr",
    "permissions": ["products.read", "products.write", "orders.read"]
  }' | jq

# Expected: Invitation created with custom permissions
# User will have products.write even though "viewer" role normally doesn't
```

## Part 7: Test Google OAuth (Optional)

If you've configured Google OAuth (see `GOOGLE_OAUTH_SETUP.md`):

### Test OAuth Flow

```bash
# Open in browser
open http://localhost:8080/api/v1/auth/google

# Or with curl (will show redirect URL)
curl -v http://localhost:8080/api/v1/auth/google \
  -H "Host: testdomain.com" 2>&1 | grep Location
```

You should be redirected to Google for sign-in.

After signing in, you'll be redirected to:
```
http://localhost:3000/auth/callback?token=eyJhbGc...
```

## Part 8: Integration Testing Script

Save this as `test.sh`:

```bash
#!/bin/bash

API_URL="http://localhost:8080/api/v1"
DOMAIN="testdomain.com"

echo "🧪 Auth Module Integration Test"
echo ""

# 1. Health check
echo "1️⃣ Testing health endpoint..."
curl -s $API_URL/../health | jq -r '.status'

# 2. Request magic link
echo "2️⃣ Requesting magic link..."
curl -s -X POST $API_URL/auth/magic-link/request \
  -H "Host: $DOMAIN" \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com"}' | jq -r '.message'

# 3. Get invitation token (you need to create domain first)
echo "3️⃣ Testing invitation verify..."
# (Manual: replace with actual token)
# curl -s "$API_URL/auth/invitation/verify?token=TOKEN" | jq

echo ""
echo "✅ Basic tests complete!"
echo "Run full tests manually using commands in TESTING.md"
```

Make executable and run:
```bash
chmod +x test.sh
./test.sh
```

## Part 9: View Data in MongoDB

### Using Mongo Shell

```bash
make docker-mongo-shell

# View all domains
use auth_module
db.domains.find().pretty()

# View all users
db.users.find().pretty()

# View invitations
db.invitations.find().pretty()

# View invitation logs
db.invitation_logs.find().pretty()

# Exit
exit
```

### Using Mongo Express (Web UI)

```bash
# Start with debug profile
make docker-up-debug

# Open browser to:
open http://localhost:8081

# Login: admin / admin
# Database: auth_module
```

You can browse all collections visually!

## Part 10: Test Cleanup

### Reset Everything

```bash
# Stop and remove all data
make docker-clean

# Start fresh
make docker-up

# Recreate domain
make domain-create DOMAIN=testdomain.com NAME="Test Domain" EMAIL=admin@testdomain.com
```

## Common Issues & Solutions

### "Connection refused" on port 8080

**Problem:** API not running

**Solution:**
```bash
make docker-logs-api
# Check for errors
make docker-restart
```

### "Domain not found"

**Problem:** Using Host header for non-existent domain

**Solution:**
```bash
make domain-list
# Use a domain from the list
# Or create new domain: make domain-create ...
```

### "Invalid or expired token"

**Problem:** JWT token expired (default 24 hours) or wrong secret

**Solution:**
```bash
# Get a new token by logging in again
# Or check JWT_SECRET in .env matches
```

### "Insufficient permissions"

**Problem:** User role doesn't have permission for admin endpoints

**Solution:**
```bash
# Use admin JWT token (from accepting admin invitation)
# Regular users created via magic link have "customer" role
```

## Next Steps

1. ✅ All tests passing? Great!
2. Build a frontend (React/Vue) to integrate with these APIs
3. Set up email (SMTP) to test magic link emails
4. Configure Google OAuth to test social login
5. Deploy to production

## Production Testing Checklist

Before deploying to production:

- [ ] Change `JWT_SECRET` to strong random value
- [ ] Configure real SMTP settings
- [ ] Set up Google OAuth with production URLs
- [ ] Test with real domain names
- [ ] Enable HTTPS
- [ ] Test all auth flows end-to-end
- [ ] Test role-based access control
- [ ] Test domain isolation
- [ ] Load test with multiple concurrent users
- [ ] Review security headers and CORS

---

**Happy Testing! 🎉**

For questions or issues, check:
- `API.md` - Complete API reference
- `README.md` - Full documentation
- `QUICKSTART.md` - Quick start guide
