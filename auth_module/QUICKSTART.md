# Auth Module - Quick Start Guide

## Prerequisites

- Go 1.21+
- MongoDB running on localhost:27017

## Step 1: Build the Module

```bash
make build
```

## Step 2: Configure

Edit `config.dev.yaml` and set your JWT secret:

```yaml
jwt:
  secret: "your-secret-key-here"  # Change this!
```

Or use environment variable:
```bash
export AUTH_JWT_SECRET="your-secret-key-here"
```

## Step 3: Start the Server

```bash
./auth-module serve --config=config.dev.yaml
```

Server will start on http://localhost:8080

## Step 4: Create a Domain

In another terminal:

```bash
./auth-module domain create \
  --config=config.dev.yaml \
  --domain=oilyourhair.com \
  --name="Oil Your Hair" \
  --admin-email=admin@oilyourhair.com
```

## Step 5: Create an API Key for Products Service

```bash
./auth-module apikey create \
  --config=config.dev.yaml \
  --domain=oilyourhair.com \
  --service=products \
  --description="Products service API key" \
  --permissions=products.read,products.write \
  --expires-in=365
```

**IMPORTANT:** Save the API key that is displayed! You'll need it for the products module.

## Step 6: List API Keys

```bash
./auth-module apikey list \
  --config=config.dev.yaml \
  --domain=oilyourhair.com
```

## Step 7: Test Health Endpoint

```bash
curl http://localhost:8080/health
```

## Using Makefile Commands

```bash
# Create domain
make domain-create-local DOMAIN=oilyourhair.com NAME="Oil Your Hair" EMAIL=admin@oilyourhair.com

# List domains
make domain-list-local

# Create API key
make apikey-create-local DOMAIN=oilyourhair.com SERVICE=products PERMISSIONS="products.read,products.write"

# List API keys
make apikey-list-local DOMAIN=oilyourhair.com
```

## Next Steps

Now you can use the API key to access the products module. See `../products_module/QUICKSTART.md`
