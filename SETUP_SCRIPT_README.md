# OilYourHair Setup Script

Automated setup script that recreates the complete oilyourhair.com configuration with all products and reviews.

## What This Script Does

This Python script automates the entire setup process:

1. ✅ **Checks service health** - Verifies auth and products APIs are running
2. ✅ **Creates domain** - Sets up oilyourhair.com domain
3. ✅ **Creates invite token** - Generates token for user registration
4. ✅ **Creates admin user** - Sets up admin@oilyourhair.com
5. ✅ **Generates API key** - Creates key for products management
6. ✅ **Populates products** - Adds 6 hair care products
7. ✅ **Populates reviews** - Adds 6 five-star reviews
8. ✅ **Verifies setup** - Tests everything works correctly

## Prerequisites

### 1. Install Python Dependencies

```bash
pip install requests colorama
```

### 2. Start Required Services

```bash
# Start MongoDB
cd /home/sparque/dev/oilyourhairapp/auth_module
make dev-db

# Start Auth Module (in one terminal)
cd /home/sparque/dev/oilyourhairapp/auth_module
./auth-module serve

# Start Products API (in another terminal)
cd /home/sparque/dev/oilyourhairapp
docker compose up -d products-api
```

## Usage

### Basic Usage

```bash
cd /home/sparque/dev/oilyourhairapp
python3 setup_oilyourhair.py
```

### Expected Output

The script will display colored output for each step:

```
================================================================================
  OilYourHair.com Setup Script
  Automated setup and verification
================================================================================

================================================================================
STEP 0: Checking Service Health
================================================================================
✓ Auth Module is healthy
✓ Products API is healthy

================================================================================
STEP 1: Creating Domain
================================================================================
ℹ Creating domain: oilyourhair.com
✓ Domain created: oilyourhair.com
Domain ID: 507f1f77bcf86cd799439011

================================================================================
STEP 2: Creating Invite Token
================================================================================
ℹ Creating invite token for oilyourhair.com
✓ Invite token created
Token: p0afXqqyJlcgGqHyoZdry95...

... (continues for all steps)

================================================================================
  ✓ SETUP COMPLETE!
================================================================================

Summary:
  Domain: oilyourhair.com
  Admin: admin@oilyourhair.com
  Products: 6
  Reviews: 6

Next Steps:
  1. Build frontend: cd frontend && ./build-site.sh oilyourhair.com
  2. Access site: http://localhost:8000
  3. Login with: admin@oilyourhair.com
```

## What Gets Created

### Admin User
- **Email:** admin@oilyourhair.com
- **Password:** SecurePass123!
- **Name:** Admin User

### Products (6 total)
1. **Nourishing Argan Oil** - $29.99 (Bestseller)
2. **Hydrating Coconut Oil** - $24.99 (Bestseller)
3. **Strengthening Castor Oil** - $19.99
4. **Shine & Smooth Jojoba Oil** - $27.99
5. **Repair Hair Treatment** - $34.99 (Bestseller)
6. **Moisturizing Hair Mask** - $32.99

### Reviews (6 five-star reviews)
- Sarah Johnson - Nourishing Argan Oil
- Michael Chen - Hydrating Coconut Oil
- Emily Rodriguez - Moisturizing Hair Mask
- David Thompson - Strengthening Castor Oil
- Jessica Park - Repair Hair Treatment
- Ryan Martinez - Shine & Smooth Jojoba Oil

## Verification

The script automatically verifies:
- All products are accessible via public API
- All reviews are accessible via public API
- Admin can authenticate successfully
- Bestseller products are marked correctly
- All reviews are 5-star ratings

## Troubleshooting

### Error: "Auth Module is not accessible"

**Solution:** Start the auth module:
```bash
cd /home/sparque/dev/oilyourhairapp/auth_module
./auth-module serve
```

### Error: "Products API is not accessible"

**Solution:** Start the products API:
```bash
cd /home/sparque/dev/oilyourhairapp
docker compose up -d products-api
```

### Error: "Domain already exists"

**Solution:** The script will continue if domain exists. If you want a fresh start:
```bash
# WARNING: This deletes all data for the domain
docker exec auth-mongodb mongosh "mongodb://localhost:27017/auth_module" \
  --eval 'db.domains.deleteOne({domain: "oilyourhair.com"})'

docker exec auth-mongodb mongosh "mongodb://localhost:27017/products_module" \
  --eval 'db.products.deleteMany({domain: "oilyourhair.com"})'

docker exec auth-mongodb mongosh "mongodb://localhost:27017/products_module" \
  --eval 'db.reviews.deleteMany({domain: "oilyourhair.com"})'
```

### Error: "User already exists"

**Solution:** The script will continue if user exists. You can still login with the existing credentials.

### Error: "Failed to create invite token"

**Solution:** The script will prompt you to manually paste the invite token. Run this command in another terminal:
```bash
cd /home/sparque/dev/oilyourhairapp/auth_module
./auth-module domain create-invite oilyourhair.com
```

Then paste the token when prompted.

## Manual Verification

After the script completes, you can manually verify:

### Check Products
```bash
curl "http://localhost:9091/api/v1/public/oilyourhair.com/products" | jq '.count'
# Should return: 6
```

### Check Reviews
```bash
curl "http://localhost:9091/api/v1/public/oilyourhair.com/reviews" | jq '.count'
# Should return: 6
```

### Check Admin Login
```bash
curl -X POST "http://localhost:8080/api/v1/auth/login" \
  -H "Content-Type: application/json" \
  -d '{
    "email": "admin@oilyourhair.com",
    "password": "SecurePass123!"
  }' | jq '.token'
# Should return a JWT token
```

## Next Steps After Setup

1. **Build the frontend:**
   ```bash
   cd /home/sparque/dev/oilyourhairapp/frontend
   ./build-site.sh oilyourhair.com
   ```

2. **Serve the site locally:**
   ```bash
   cd sites/oilyourhair.com/public
   python3 -m http.server 8000
   ```

3. **Access the website:**
   Open browser to: http://localhost:8000

4. **Login as admin:**
   - Email: admin@oilyourhair.com
   - Password: SecurePass123!

## Script Architecture

### Functions

- `check_service_health()` - Verifies APIs are running
- `create_domain()` - Creates domain via REST API
- `create_invite_token()` - Calls CLI to generate invite token
- `create_admin_user()` - Registers admin via signup API
- `login_admin()` - Authenticates and gets JWT token
- `create_api_key()` - Calls CLI to generate API key
- `get_products_data()` - Returns product definitions
- `create_products()` - Posts products to API
- `get_reviews_data()` - Returns review definitions
- `create_reviews()` - Posts reviews to API
- `verify_setup()` - Tests all endpoints

### Color Coding

- 🟢 **Green (✓)** - Success
- 🔴 **Red (✗)** - Error
- 🟡 **Yellow (ℹ)** - Information
- 🟣 **Magenta** - Data values
- 🔵 **Cyan** - Step headers

## Extending the Script

To add more products or reviews, edit these functions:

### Add Products
Edit `get_products_data()` function around line 250:
```python
def get_products_data() -> List[Dict]:
    return [
        {
            "domain": DOMAIN,
            "name": "Your Product Name",
            "description": "Description here",
            "price": 39.99,
            # ... more fields
        },
        # Add more products here
    ]
```

### Add Reviews
Edit `get_reviews_data()` function around line 370:
```python
def get_reviews_data() -> List[Dict]:
    return [
        {
            "product": "Product Name",
            "name": "Reviewer Name",
            "rating": 5,
            "text": "Review text",
            "highlight": "Key highlight"
        },
        # Add more reviews here
    ]
```

## License

This script is part of the OilYourHair platform.

## Support

For issues or questions, refer to the main USER_GUIDE.md or contact the development team.
