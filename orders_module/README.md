# Orders Module

Order processing, payment handling, and stock management for the multi-tenant e-commerce platform.

## Features

- ✅ Order creation and management
- ✅ Stripe payment integration (test mode)
- ✅ Stock deduction on successful payment
- ✅ Stock transaction audit trail
- ✅ Multi-tenant support (domain isolation)
- ✅ Webhook handling for payment confirmations

## MVP Scope (Phase 1)

- Single currency (USD)
- Basic stock deduction (no reservation)
- Stripe test mode
- Simple order status tracking
- Guest and authenticated checkout

## Getting Started

### Prerequisites

- Go 1.24+
- MongoDB
- Stripe account (test mode)

### Installation

```bash
# Install dependencies
go mod download

# Copy config
cp config.dev.yaml config.yaml

# Add your Stripe test keys to config.yaml
```

### Configuration

Update `config.yaml` with your Stripe test keys:

```yaml
stripe:
  secret_key: "sk_test_YOUR_KEY_HERE"
  webhook_secret: "whsec_YOUR_WEBHOOK_SECRET"
  publishable_key: "pk_test_YOUR_KEY_HERE"
```

### Running

```bash
# Development
go run main.go serve

# Build
go build -o orders-module

# Run binary
./orders-module serve
```

Server will start on port 9092.

### API Endpoints

**Orders:**
- `POST /api/v1/orders` - Create order and payment intent
- `GET /api/v1/orders/:id` - Get order details
- `GET /api/v1/orders` - List user's orders

**Webhooks:**
- `POST /api/v1/webhooks/stripe` - Stripe payment webhook

**Admin:**
- `GET /api/v1/admin/orders` - List all orders
- `PATCH /api/v1/admin/orders/:id/status` - Update order status

### Testing with Stripe

Use these test card numbers:

```
Success:
4242 4242 4242 4242 - Visa (always succeeds)

Decline:
4000 0000 0000 0002 - Card declined

Any future expiry date, any 3-digit CVC
```

## Database Schema

See [DATABASE_SCHEMA.md](DATABASE_SCHEMA.md) for detailed schema documentation.

## Architecture

```
orders_module (port 9092)
├── Receives order creation requests
├── Creates Stripe payment intent
├── Listens for Stripe webhooks
├── Deducts stock from products_module
└── Stores orders in MongoDB
```

## Development Roadmap

**Phase 1 (MVP)** ← Current
- Basic order processing
- Stripe integration
- Simple stock deduction

**Phase 2 (Enhancements)**
- Stock reservation system
- Multi-currency support
- Advanced order management

**Phase 3 (Advanced)**
- Refund handling
- Discount codes
- Email notifications
