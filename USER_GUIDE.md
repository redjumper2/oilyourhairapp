# OilYourHair Platform - User Guide

Complete guide for setting up, configuring, and testing the multi-tenant e-commerce platform.

---

## Part 1: Initial Setup & Configuration

### Prerequisites
- Docker and Docker Compose installed
- MongoDB running in Docker
- Go 1.24+ installed
- Basic understanding of REST APIs

---

## Step 1: Setup Auth Module

### 1.1 Start the Auth Module

```bash
cd /home/sparque/dev/oilyourhairapp/auth_module
make dev-db  # Start MongoDB if not running
./auth-module serve
```

The auth module will start on `http://localhost:8080`

### 1.2 Create a Domain

```bash
# Create domain for oilyourhair.com
curl -X POST "http://localhost:8080/api/v1/domains" \
  -H "Content-Type: application/json" \
  -d '{
    "domain": "oilyourhair.com",
    "name": "Oil Your Hair",
    "settings": {
      "allow_signups": true,
      "require_email_verification": false
    }
  }'
```

**Expected Response:**
```json
{
  "id": "...",
  "domain": "oilyourhair.com",
  "name": "Oil Your Hair",
  "created_at": "..."
}
```

### 1.3 List All Domains

```bash
./auth-module domain list
```

### 1.4 Create an Invite Token

```bash
./auth-module domain create-invite oilyourhair.com
```

**Save the token** - you'll use it to create your admin user.

### 1.5 Create Admin User

```bash
# Use the invite token from above
INVITE_TOKEN="your-invite-token-here"

curl -X POST "http://localhost:8080/api/v1/auth/signup" \
  -H "Content-Type: application/json" \
  -d "{
    \"email\": \"admin@oilyourhair.com\",
    \"password\": \"SecurePass123!\",
    \"first_name\": \"Admin\",
    \"last_name\": \"User\",
    \"invite_token\": \"$INVITE_TOKEN\"
  }"
```

### 1.6 Login as Admin

```bash
curl -X POST "http://localhost:8080/api/v1/auth/login" \
  -H "Content-Type: application/json" \
  -d '{
    "email": "admin@oilyourhair.com",
    "password": "SecurePass123!"
  }'
```

**Save the JWT token** from the response - you'll need it for authenticated requests.

---

## Step 2: Setup Products Module

### 2.1 Start Products Module

```bash
cd /home/sparque/dev/oilyourhairapp
docker compose up -d products-api
```

The products API will start on `http://localhost:9091`

### 2.2 Create API Key for Products

```bash
cd /home/sparque/dev/oilyourhairapp/auth_module

# Create API key with products permissions
./auth-module apikey create oilyourhair.com \
  --service products \
  --permissions products.read,products.write,products.delete
```

**Save the API key** - you'll use it to manage products.

### 2.3 Test Health Check

```bash
curl http://localhost:9091/health
```

**Expected Response:**
```json
{
  "status": "healthy",
  "env": "development"
}
```

### 2.4 Create a Test Product

```bash
API_KEY="your-api-key-here"

curl -X POST "http://localhost:9091/api/v1/products" \
  -H "Authorization: Bearer $API_KEY" \
  -H "Content-Type: application/json" \
  -d '{
    "domain": "oilyourhair.com",
    "name": "Nourishing Argan Oil",
    "description": "Pure organic argan oil for beautiful, healthy hair",
    "price": 29.99,
    "category": "Hair Oils",
    "attributes": {
      "size": "100ml",
      "ingredients": "100% Pure Argan Oil",
      "origin": "Morocco"
    },
    "images": [
      "https://images.unsplash.com/photo-1571781926291-c477ebfd024b"
    ],
    "active": true,
    "bestseller": true
  }'
```

### 2.5 List All Products

```bash
# Public endpoint (no auth required)
curl "http://localhost:9091/api/v1/public/oilyourhair.com/products"
```

### 2.6 Get Product by ID

```bash
PRODUCT_ID="product-id-from-above"

curl "http://localhost:9091/api/v1/public/oilyourhair.com/products/$PRODUCT_ID"
```

---

## Step 3: Test Reviews Feature

### 3.1 Create a Review (Public - No Auth)

```bash
curl -X POST "http://localhost:9091/api/v1/public/oilyourhair.com/reviews" \
  -H "Content-Type: application/json" \
  -d '{
    "product": "Nourishing Argan Oil",
    "name": "Sarah Johnson",
    "rating": 5,
    "text": "Amazing product! My hair has never been softer.",
    "highlight": "Best hair oil ever!"
  }'
```

### 3.2 List All Reviews

```bash
curl "http://localhost:9091/api/v1/public/oilyourhair.com/reviews"
```

### 3.3 Filter Reviews by Product

```bash
curl "http://localhost:9091/api/v1/public/oilyourhair.com/reviews?product_id=$PRODUCT_ID"
```

---

## Step 4: Test Contact Form

### 4.1 Submit Contact Form

```bash
curl -X POST "http://localhost:9091/api/v1/public/oilyourhair.com/contact" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "John Doe",
    "email": "john@example.com",
    "subject": "Product Inquiry",
    "message": "I would like to know more about your products."
  }'
```

**Expected Response:**
```json
{
  "success": true,
  "message": "Thank you for contacting us! We'll get back to you soon.",
  "id": "..."
}
```

### 4.2 View Contact Submissions (Direct MongoDB)

```bash
docker exec auth-mongodb mongosh "mongodb://localhost:27017/products_module" \
  --eval 'db.contacts.find({domain: "oilyourhair.com"}).toArray()' \
  --quiet
```

---

## Step 5: Test Frontend

### 5.1 Build Frontend for oilyourhair.com

```bash
cd /home/sparque/dev/oilyourhairapp/frontend
./build-site.sh oilyourhair.com
```

### 5.2 Serve Frontend Locally

```bash
# Using Python
cd sites/oilyourhair.com/public
python3 -m http.server 8000

# Or using nginx (if configured)
sudo nginx -t && sudo nginx -s reload
```

### 5.3 Access the Website

Open browser to: `http://localhost:8000` (or configured domain)

**Test the following:**
- ✅ Homepage loads
- ✅ Products page shows products
- ✅ Reviews page shows reviews
- ✅ Contact form submits successfully
- ✅ Login/Signup works
- ✅ Cart functionality works

---

## Part 2: Multi-Tenant Setup (Second Domain)

### Goal
Set up `example.com` as a second website using the same template and backend.

---

## Step 6: Configure Second Domain

### 6.1 Create Domain in Auth Module

```bash
curl -X POST "http://localhost:8080/api/v1/domains" \
  -H "Content-Type: application/json" \
  -d '{
    "domain": "example.com",
    "name": "Example Store",
    "settings": {
      "allow_signups": true,
      "require_email_verification": false
    }
  }'
```

### 6.2 Create Invite for Example.com

```bash
./auth-module domain create-invite example.com
```

### 6.3 Create Admin User for Example.com

```bash
INVITE_TOKEN="new-invite-token-here"

curl -X POST "http://localhost:8080/api/v1/auth/signup" \
  -H "Content-Type: application/json" \
  -d "{
    \"email\": \"admin@example.com\",
    \"password\": \"SecurePass123!\",
    \"first_name\": \"Example\",
    \"last_name\": \"Admin\",
    \"invite_token\": \"$INVITE_TOKEN\"
  }"
```

### 6.4 Create API Key for Example.com Products

```bash
./auth-module apikey create example.com \
  --service products \
  --permissions products.read,products.write,products.delete
```

**Save the new API key**

---

## Step 7: Add Products for Example.com

### 7.1 Create Test Products

```bash
API_KEY_EXAMPLE="api-key-for-example.com"

curl -X POST "http://localhost:9091/api/v1/products" \
  -H "Authorization: Bearer $API_KEY_EXAMPLE" \
  -H "Content-Type: application/json" \
  -d '{
    "domain": "example.com",
    "name": "Premium Shampoo",
    "description": "Luxury shampoo for all hair types",
    "price": 24.99,
    "category": "Shampoos",
    "active": true,
    "bestseller": true
  }'
```

### 7.2 Verify Domain Isolation

```bash
# Should only show example.com products
curl "http://localhost:9091/api/v1/public/example.com/products"

# Should only show oilyourhair.com products
curl "http://localhost:9091/api/v1/public/oilyourhair.com/products"
```

**Important:** Each domain's data is completely isolated!

---

## Step 8: Setup Frontend for Example.com

### 8.1 Create Branding Configuration

```bash
mkdir -p /home/sparque/dev/oilyourhairapp/frontend/sites/example.com
```

Create branding file:

```bash
cat > /home/sparque/dev/oilyourhairapp/frontend/sites/example.com/branding.json <<EOF
{
  "siteName": "Example Store",
  "brandName": "Example",
  "tagline": "Your Trusted Store",
  "primaryColor": "#3B82F6",
  "primaryColorLight": "#60A5FA",
  "logo": "https://via.placeholder.com/150x50?text=Example",
  "domain": "example.com",
  "apiUrl": "http://localhost:9091",
  "authUrl": "http://localhost:8080"
}
EOF
```

### 8.2 Build Frontend for Example.com

```bash
cd /home/sparque/dev/oilyourhairapp/frontend
./build-site.sh example.com
```

Output will be in: `sites/example.com/public/`

### 8.3 Configure Nginx (Optional)

Add to nginx config:

```nginx
server {
    listen 80;
    server_name example.com;

    root /home/sparque/dev/oilyourhairapp/frontend/sites/example.com/public;
    index index.html;

    location / {
        try_files $uri $uri/ =404;
    }
}
```

Reload nginx:
```bash
sudo nginx -s reload
```

---

## Step 9: Test Example.com Features

### 9.1 Test Products API

```bash
# List products
curl "http://localhost:9091/api/v1/public/example.com/products"

# Should return only example.com products
```

### 9.2 Test Reviews

```bash
# Create review for example.com
curl -X POST "http://localhost:9091/api/v1/public/example.com/reviews" \
  -H "Content-Type: application/json" \
  -d '{
    "product": "Premium Shampoo",
    "name": "Jane Smith",
    "rating": 5,
    "text": "Excellent shampoo!",
    "highlight": "Love it!"
  }'

# List reviews
curl "http://localhost:9091/api/v1/public/example.com/reviews"
```

### 9.3 Test Contact Form

```bash
curl -X POST "http://localhost:9091/api/v1/public/example.com/contact" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Test User",
    "email": "test@example.com",
    "subject": "Test",
    "message": "Testing contact form for example.com"
  }'
```

### 9.4 Test Authentication

```bash
# Login to example.com
curl -X POST "http://localhost:8080/api/v1/auth/login" \
  -H "Content-Type: application/json" \
  -d '{
    "email": "admin@example.com",
    "password": "SecurePass123!"
  }'
```

### 9.5 Access Example.com Frontend

Open browser to: `http://example.com` (or localhost with port)

**Verify:**
- ✅ Site shows "Example Store" branding
- ✅ Different colors (blue instead of green)
- ✅ Shows only example.com products
- ✅ Reviews are separate from oilyourhair.com
- ✅ Login works with example.com users
- ✅ Contact form submits to example.com domain

---

## Step 10: Verify Complete Domain Isolation

### 10.1 Database Verification

```bash
# Check products are domain-isolated
docker exec auth-mongodb mongosh "mongodb://localhost:27017/products_module" \
  --eval 'db.products.find({}, {domain: 1, name: 1}).toArray()' \
  --quiet

# Should see:
# - Products with domain: "oilyourhair.com"
# - Products with domain: "example.com"
```

### 10.2 Reviews Isolation

```bash
# oilyourhair.com reviews
curl "http://localhost:9091/api/v1/public/oilyourhair.com/reviews" | jq '.count'

# example.com reviews (should be different count)
curl "http://localhost:9091/api/v1/public/example.com/reviews" | jq '.count'
```

### 10.3 User Isolation

```bash
# Login with oilyourhair.com user should NOT work on example.com API calls
# And vice versa - JWT tokens are domain-scoped
```

---

## Troubleshooting

### Products Not Showing

```bash
# Check MongoDB connection
docker ps | grep mongo

# Check products in DB
docker exec auth-mongodb mongosh "mongodb://localhost:27017/products_module" \
  --eval 'db.products.countDocuments({domain: "oilyourhair.com"})'

# Check API logs
docker logs products-api --tail 50
```

### Auth Not Working

```bash
# Check auth module is running
curl http://localhost:8080/health

# Verify domain exists
./auth-module domain list

# Check JWT token is valid
curl -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  http://localhost:8080/api/v1/auth/me
```

### Frontend Not Loading

```bash
# Rebuild site
cd /home/sparque/dev/oilyourhairapp/frontend
./build-site.sh oilyourhair.com

# Check output exists
ls -la sites/oilyourhair.com/public/

# Verify branding.json exists
cat sites/oilyourhair.com/branding.json
```

### CORS Issues

If accessing from different domain:
- Update products_module config to allow your domain in CORS
- Check browser console for CORS errors

---

## Quick Reference

### Important URLs

| Service | URL |
|---------|-----|
| Auth API | http://localhost:8080 |
| Products API | http://localhost:9091 |
| MongoDB | mongodb://localhost:27017 |
| Frontend (local) | http://localhost:8000 |

### Database Collections

| Collection | Database | Purpose |
|------------|----------|---------|
| domains | auth_module | Domain configurations |
| users | auth_module | User accounts (per domain) |
| products | products_module | Products (domain-isolated) |
| reviews | products_module | Reviews (domain-isolated) |
| contacts | products_module | Contact submissions (domain-isolated) |

### Common Commands

```bash
# List domains
./auth-module domain list

# Create invite
./auth-module domain create-invite DOMAIN

# Create API key
./auth-module apikey create DOMAIN --service SERVICE --permissions PERMS

# Build frontend
./build-site.sh DOMAIN

# View logs
docker logs products-api
docker logs auth-mongodb

# Restart services
docker compose restart products-api
```

---

## Summary

You now have:
1. ✅ Auth module with domain management
2. ✅ Products API with multi-tenant support
3. ✅ Reviews system (domain-isolated)
4. ✅ Contact form system (domain-isolated)
5. ✅ Template-based frontend that works for any domain
6. ✅ Two working examples: oilyourhair.com and example.com

**Each domain is completely isolated:**
- Separate users
- Separate products
- Separate reviews
- Separate contact submissions
- Separate branding/styling

**But shares:**
- Same codebase
- Same backend services
- Same template system
- Same API infrastructure

Perfect for a multi-tenant SaaS e-commerce platform! 🎉
