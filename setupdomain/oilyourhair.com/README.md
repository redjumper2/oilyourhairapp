# OilYourHair.com Domain Setup

Complete, configurable end-to-end domain setup from JSON configuration files.

## Overview

This setup script creates a **complete end-to-end domain environment** including:

### Backend Setup
- ✅ Domain configuration and settings (auth, security, branding)
- ✅ Admin user with full permissions and roles
- ✅ API keys for service access
- ✅ Product catalog with comprehensive details (variants, inventory, SEO)
- ✅ Customer reviews (5-star ratings with metadata)

### Frontend Setup (Automatic!)
- ✅ **Branding configuration auto-generated** from domain config
- ✅ **Website built from template1** (store-template1)
- ✅ **All HTML/CSS/JS files generated** and ready to serve
- ✅ **Colors, logo, and site name applied** from branding config

### Result
**One command gives you a fully working e-commerce site!** Backend + Frontend ready to go.

  ✅ All 10 Steps Successful:
  1. Database Cleaned - All collections cleared
  2. Domain Created - oilyourhair.com with invitation token
  3. Admin User Created - with JWT token automatically
  4. API Key Created - for products service
  5. 6 Products Created - with variants, images, attributes
  6. 6 Reviews Created - all 5-star ratings
  7. Frontend Branding - Auto-generated from domain config
  8. Frontend Built - 16 files from template1
  9. Verification Passed - All checks successful

## Quick Start

### 1. Create Virtual Environment

```bash
# Create virtual environment
cd /home/sparque/dev/oilyourhairapp/setupdomain/oilyourhair.com
python3 -m venv venv

# Activate virtual environment
source venv/bin/activate  # On Linux/Mac
# or
venv\Scripts\activate  # On Windows
```

### 2. Install Dependencies

```bash
pip install -r requirements.txt
```

### 3. Run Setup

```bash
# Clean database and setup from scratch
python3 setup.py --clean

# Or just setup (keep existing data)
python3 setup.py
```

## What the Script Does

The setup script performs **10 automated steps** to create a complete working site:

### Backend Setup (Steps 1-7)
1. **Check Service Health** - Verifies auth and products APIs are running
2. **Create Domain** - Sets up domain in auth module with all settings
3. **Create Invite Token** - Generates token for user registration
4. **Create Admin User** - Creates admin account with roles and permissions
5. **Authenticate Admin** - Logs in and gets JWT token
6. **Create API Key** - Generates API key for products service
7. **Populate Products** - Creates 6 products with full details (variants, inventory, SEO)
8. **Populate Reviews** - Creates 6 five-star reviews with metadata

### Frontend Setup (Steps 8-9) ✨ NEW!
9. **Create Branding Config** - Auto-generates `branding.json` from domain config
   - Extracts colors, logo, site name from `config/domain.json`
   - Creates API endpoint URLs
   - Saves to `frontend/sites/{domain}/branding.json`

10. **Build Frontend** - Automatically builds website from template1
    - Runs `./build-site.sh {domain}`
    - Copies all files from `templates/store-template1/`
    - Applies branding (colors, logos, site name)
    - Generates output to `frontend/sites/{domain}/public/`
    - Shows file count on completion

### Verification (Step 10)
11. **Verify Setup** - Tests all endpoints and confirms everything works
    - Products API accessible
    - Reviews API accessible
    - Admin authentication works
    - Counts match expected values

### Success!
After completion, you get:
- ✅ Backend fully configured and populated
- ✅ Frontend built and branded
- ✅ **Ready to serve immediately** - just run a web server!

## Configuration Files

All settings are externalized in JSON files under `config/`:

### `config/domain.json`
Complete domain configuration including:
- Domain name and description
- Signup and authentication settings
- Password policies
- Session management
- Security settings (login attempts, lockout)
- Branding (colors, logo)
- Contact information

**Customizable Fields:**
```json
{
  "domain": "your domain.com",
  "name": "Display Name",
  "settings": {
    "allow_signups": true/false,
    "require_email_verification": true/false,
    "password_policy": {
      "min_length": 8,
      "require_uppercase": true,
      ...
    },
    ...
  }
}
```

### `config/admin.json`
Admin user configuration with:
- Basic credentials (email, password, name)
- Profile information
- Roles and permissions
- Metadata (employee ID, department)
- Preferences (notifications, etc.)
- Status flags

**Customizable Fields:**
```json
{
  "email": "admin@domain.com",
  "password": "SecurePassword123!",
  "roles": ["admin", "editor"],
  "permissions": [
    "products.read",
    "products.write",
    ...
  ]
}
```

### `config/products.json`
Product catalog with ALL product fields:
- Basic info (name, description, SKU, price)
- Categories and tags
- Attributes (size, weight, ingredients, etc.)
- Images with positioning
- Variants (different sizes/options)
- Inventory tracking
- SEO metadata
- Shipping details

**Example Product:**
```json
{
  "name": "Product Name",
  "sku": "SKU-CODE",
  "price": 29.99,
  "category": "Category",
  "variants": [
    {
      "name": "100ml",
      "sku": "SKU-100ML",
      "price": 29.99,
      "stock": 50
    }
  ],
  "inventory": {
    "track_inventory": true,
    "stock_quantity": 50,
    "low_stock_threshold": 10
  },
  "seo": {
    "meta_title": "SEO Title",
    "url_slug": "product-url-slug"
  }
}
```

### `config/reviews.json`
Customer reviews with:
- Product reference
- Reviewer information
- Rating (1-5 stars)
- Review text and title
- Verified purchase flag
- Helpful count
- Metadata (hair type, usage duration)
- Status (approved, pending, rejected)

### `config/services.json`
Service endpoints and configuration:
- Auth service URLs and endpoints
- Products service URLs and endpoints
- Database connection settings
- API key configurations

## Usage Options

### Normal Setup
```bash
python3 setup.py
```
Sets up domain using configuration files. Skips if data already exists.

### Clean Setup (Fresh Start)
```bash
python3 setup.py --clean
```
**WARNING:** Deletes ALL existing data for this domain, then sets up from scratch.

### Clean Only
```bash
python3 setup.py --clean-only
```
Only cleans the database without setting up. Useful for manual testing.

### Verify Only
```bash
python3 setup.py --verify-only
```
Only verifies the existing setup without making changes.

## What Gets Created

### Domain
- Domain: oilyourhair.com
- Name: Oil Your Hair
- Settings: As configured in domain.json

### Admin User
- Email: admin@oilyourhair.com
- Password: SecurePass123!
- Roles: admin, editor, moderator
- Full permissions for products and domain management

### Products (6 total)
1. Nourishing Argan Oil - $29.99 (Bestseller)
2. Hydrating Coconut Oil - $24.99 (Bestseller)
3. Strengthening Castor Oil - $19.99
4. Shine & Smooth Jojoba Oil - $27.99
5. Repair Hair Treatment - $34.99 (Bestseller, On Sale)
6. Moisturizing Hair Mask - $32.99

Each with:
- Full product details
- Variants (where applicable)
- Inventory tracking
- SEO metadata
- Shipping information

### Reviews (6 five-star reviews)
- Sarah Johnson - Nourishing Argan Oil ⭐⭐⭐⭐⭐
- Michael Chen - Hydrating Coconut Oil ⭐⭐⭐⭐⭐
- Emily Rodriguez - Moisturizing Hair Mask ⭐⭐⭐⭐⭐
- David Thompson - Strengthening Castor Oil ⭐⭐⭐⭐⭐
- Jessica Park - Repair Hair Treatment ⭐⭐⭐⭐⭐
- Ryan Martinez - Shine & Smooth Jojoba Oil ⭐⭐⭐⭐⭐

### Frontend Website
**Location:** `/home/sparque/dev/oilyourhairapp/frontend/sites/oilyourhair.com/`

**Files Created:**
- `branding.json` - Auto-generated from domain config
- `public/` - Complete website ready to serve
  - `index.html` - Homepage with featured products
  - `shop.html` - Products page with all 6 items
  - `about.html` - About page
  - `reviews.html` - Reviews page with all 6 reviews
  - `contact.html` - Contact form
  - `styles.css` - Branded with your colors
  - `cart.js`, `auth.js`, `navigation.js` - Full functionality
  - All images and assets

**Branding Applied:**
- Primary Color: #2E7D32 (from config)
- Secondary Color: #4CAF50 (from config)
- Site Name: "Oil Your Hair"
- Logo: From domain config
- API URLs: Auto-configured

### Complete Output Structure
```
frontend/sites/oilyourhair.com/
├── branding.json          ← Auto-generated from config
└── public/                ← Website files (ready to serve)
    ├── index.html
    ├── shop.html
    ├── about.html
    ├── reviews.html
    ├── contact.html
    ├── styles.css         ← Branded with your colors
    ├── cart.js
    ├── auth.js
    ├── branding.js
    ├── navigation.js
    └── banner.js
```

## Prerequisites

### Services Must Be Running

```bash
# Auth Module
cd /home/sparque/dev/oilyourhairapp/auth_module
./auth-module serve

# Products API
cd /home/sparque/dev/oilyourhairapp
docker compose up -d products-api
```

### Check Service Health
```bash
# Auth service
curl http://localhost:8080/health

# Products service
curl http://localhost:9091/health
```

## Customization

### To Create a New Domain

1. Copy this folder:
   ```bash
   cp -r oilyourhair.com/ example.com/
   cd example.com/
   ```

2. Edit configuration files:
   - `config/domain.json` - Change domain, name, branding
   - `config/admin.json` - Change admin credentials
   - `config/products.json` - Customize products
   - `config/reviews.json` - Customize reviews
   - `config/services.json` - Update if using different ports

3. Run setup:
   ```bash
   source venv/bin/activate
   python3 setup.py --clean
   ```

### Adding More Products

Edit `config/products.json` and add more product objects with all fields populated.

### Adding More Reviews

Edit `config/reviews.json` and add more review objects.

## Directory Structure

### of setup script

```
oilyourhair.com/
├── config/
│   ├── domain.json      # Domain configuration
│   ├── admin.json       # Admin user config
│   ├── products.json    # Product catalog
│   ├── reviews.json     # Customer reviews
│   └── services.json    # Service endpoints
├── venv/                # Virtual environment (git-ignored)
├── setup.py             # Main setup script
├── requirements.txt     # Python dependencies
└── README.md           # This file
```

### Directory structure of what gets created:

```
  /home/sparque/dev/oilyourhairapp/
  ├── setupdomain/oilyourhair.com/
  │   ├── config/               # All externalized configs
  │   │   ├── domain.json       # Domain settings
  │   │   ├── admin.json        # Admin user
  │   │   ├── products.json     # 6 products
  │   │   ├── reviews.json      # 6 reviews
  │   │   └── services.json     # API endpoints
  │   ├── setup.py              # Main script
  │   ├── run.sh                # Helper script
  │   └── README.md             # Documentation
  ├── frontend/sites/oilyourhair.com/
  │   ├── branding.json         # Auto-generated
  │   └── public/               # 16 website files ready to serve
  └── auth_module/
      └── config.local.yaml     # Local MongoDB config
```

## Verification

The script automatically verifies:
- ✅ All products are accessible
- ✅ Product counts match
- ✅ Bestsellers are correctly marked
- ✅ All reviews are accessible
- ✅ Review ratings are correct
- ✅ Admin can authenticate

## Next Steps After Setup

**The frontend is already built!** The setup script automatically:
- ✅ Created branding.json from domain config
- ✅ Built the website from template1
- ✅ Generated all HTML/CSS/JS files

### 1. Serve the Website
```bash
cd /home/sparque/dev/oilyourhairapp/frontend/sites/oilyourhair.com/public
python3 -m http.server 8000
```

### 2. Access Website
Open browser to: **http://localhost:8000**

### 3. Login as Admin
- Email: admin@oilyourhair.com
- Password: SecurePass123!

### 4. Test Features
- Browse products (6 products loaded)
- View reviews (6 five-star reviews)
- Test shopping cart
- Submit contact form
- Sign in with admin account

## Troubleshooting

### Service Not Running
```bash
# Check if services are up
curl http://localhost:8080/health  # Auth
curl http://localhost:9091/health  # Products

# Start services if needed
cd /home/sparque/dev/oilyourhairapp/auth_module && ./auth-module serve
cd /home/sparque/dev/oilyourhairapp && docker compose up -d products-api
```

### Clean Database Manually
```bash
# Connect to MongoDB and clean
docker exec auth-mongodb mongosh "mongodb://localhost:27017/auth_module" \
  --eval 'db.domains.deleteOne({domain: "oilyourhair.com"})'

docker exec auth-mongodb mongosh "mongodb://localhost:27017/products_module" \
  --eval 'db.products.deleteMany({domain: "oilyourhair.com"})'

docker exec auth-mongodb mongosh "mongodb://localhost:27017/products_module" \
  --eval 'db.reviews.deleteMany({domain: "oilyourhair.com"})'
```

### Virtual Environment Issues
```bash
# Recreate virtual environment
rm -rf venv/
python3 -m venv venv
source venv/bin/activate
pip install -r requirements.txt
```

## Configuration Reference

### All Available Domain Settings

See `config/domain.json` for complete example with:
- Password policies (length, complexity requirements)
- Session settings (token expiry, max sessions)
- Security settings (login attempts, lockout)
- Branding (colors, logos)
- Contact information

### All Available Product Fields

See `config/products.json` for complete example with:
- Basic info (name, description, SKU, price, cost)
- Categories and tags
- Attributes (size, weight, ingredients, certifications)
- Images (with alt text and positioning)
- Variants (different sizes/options with separate pricing/stock)
- Inventory (stock tracking, thresholds, backorder)
- SEO (meta tags, URL slugs)
- Shipping (weight, dimensions, requirements)

### All Available Review Fields

See `config/reviews.json` for complete example with:
- Product reference (name and SKU)
- Reviewer info (name, email)
- Rating (1-5 stars)
- Review content (title, text, highlight)
- Verification (verified purchase flag)
- Engagement (helpful count, would recommend)
- Metadata (hair type, usage duration, age range)
- Status (approved, pending, rejected, featured)

  🔑 Key Improvements Made:

  1. ✅ Added --no-invite flag to domain CLI command
  2. ✅ All CLI commands run inside Docker containers
  3. ✅ Proper token extraction from CLI output
  4. ✅ User creation via invitation/accept endpoint
  5. ✅ JWT token obtained automatically
  6. ✅ API key extraction fixed
  7. ✅ Product mapping to correct API format
  8. ✅ Attribute type conversion (bool → string)
  9. ✅ Verification checks all working