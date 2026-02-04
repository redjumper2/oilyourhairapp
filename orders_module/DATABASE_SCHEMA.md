# Orders Module - Database Schema (MVP)

## Overview
MongoDB collections for order processing, payments, and stock management.

Database: `orders_module`

---

## Collections

### 1. `orders`

Main collection storing all order information.

```javascript
{
  // Identity
  _id: ObjectId("..."),
  order_number: "ORD-2026-00001",  // Sequential, user-friendly
  domain: "oilyourhair.com",        // Multi-tenant support

  // Customer Info
  customer: {
    user_id: "abc123",              // From auth_module (null for guest)
    email: "customer@example.com",
    name: "John Doe"
  },

  // Order Items (embedded for simplicity in MVP)
  items: [
    {
      product_id: "prod_123",
      product_name: "Luxury Hair Oil Blend",  // Snapshot (products can change)
      variant_id: "var_456",
      variant_sku: "LUX-100",
      variant_attributes: {
        size: "100ml"
      },
      quantity: 2,
      unit_price: 89.99,              // USD (base currency)
      total: 179.98
    }
  ],

  // Pricing (MVP: USD only, simple calculation)
  currency: "USD",
  subtotal: 179.98,
  tax: 0.00,                          // MVP: no tax calculation
  shipping: 0.00,                     // MVP: free shipping
  total: 179.98,

  // Payment (Stripe)
  payment: {
    provider: "stripe",
    payment_intent_id: "pi_...",      // Stripe payment intent ID
    status: "pending",                 // pending | succeeded | failed | refunded
    amount: 17998,                     // Stripe uses cents (179.98 * 100)
    currency: "usd",
    client_secret: "pi_...secret_..."  // For frontend confirmation
  },

  // Addresses (MVP: simple structure)
  shipping_address: {
    name: "John Doe",
    address_line1: "123 Main St",
    address_line2: "Apt 4B",           // Optional
    city: "New York",
    state: "NY",
    postal_code: "10001",
    country: "US",
    phone: "+1234567890"               // Optional
  },

  billing_address: {
    // Same structure as shipping_address
    // Can be same as shipping or different
  },

  // Order Status
  status: "pending",  // pending | paid | processing | shipped | delivered | cancelled | refunded

  // Timestamps
  created_at: ISODate("2026-01-05T10:00:00Z"),
  updated_at: ISODate("2026-01-05T10:00:00Z"),
  paid_at: null,      // Set when payment succeeds

  // Metadata
  notes: "",          // Optional customer notes
  admin_notes: ""     // Optional admin notes
}
```

**Indexes:**
```javascript
db.orders.createIndex({ "order_number": 1 }, { unique: true })
db.orders.createIndex({ "domain": 1, "customer.user_id": 1 })
db.orders.createIndex({ "domain": 1, "status": 1 })
db.orders.createIndex({ "payment.payment_intent_id": 1 })
db.orders.createIndex({ "created_at": -1 })
```

---

### 2. `order_counters`

Used to generate sequential order numbers per domain.

```javascript
{
  _id: "oilyourhair.com",  // Domain name
  sequence: 1,              // Next order number
  year: 2026                // Reset counter each year
}
```

**Usage:**
- Use MongoDB's `findOneAndUpdate` with `$inc` for atomic counter increment
- Format: `ORD-{year}-{sequence:05d}` → `ORD-2026-00001`

---

### 3. `stock_transactions`

Track all stock changes for audit trail.

```javascript
{
  _id: ObjectId("..."),
  domain: "oilyourhair.com",
  product_id: "prod_123",
  variant_id: "var_456",

  // Transaction details
  type: "sale",        // sale | restock | adjustment | return
  quantity: -2,        // Negative for decrease, positive for increase

  // Stock levels
  stock_before: 15,
  stock_after: 13,

  // Reference
  order_id: "ord_789",     // If related to order
  order_number: "ORD-2026-00001",

  // Metadata
  created_by: "system",    // system | user_id
  reason: "Order payment confirmed",
  created_at: ISODate("2026-01-05T10:00:00Z")
}
```

**Indexes:**
```javascript
db.stock_transactions.createIndex({ "domain": 1, "product_id": 1, "variant_id": 1 })
db.stock_transactions.createIndex({ "order_id": 1 })
db.stock_transactions.createIndex({ "created_at": -1 })
```

---

## Order Status Flow (MVP)

```
pending → paid → processing → shipped → delivered
   ↓        ↓
cancelled  refunded
```

**Status Definitions:**
- `pending` - Order created, payment not yet completed
- `paid` - Payment successful, order confirmed
- `processing` - Order being prepared (manual admin update)
- `shipped` - Order shipped (manual admin update)
- `delivered` - Order delivered (manual admin update)
- `cancelled` - Order cancelled before payment
- `refunded` - Payment refunded after successful payment

---

## Payment Status Flow

```
pending → succeeded
   ↓
 failed
   ↓
refunded (after succeeded)
```

---

## Stock Management Flow (MVP - Simple)

1. **Order Created** - No stock change (payment pending)
2. **Payment Succeeded** (webhook)
   - Deduct stock from variant
   - Create stock_transaction record
   - Update order status to `paid`
3. **Payment Failed**
   - Update order status
   - No stock changes needed

---

## MVP Simplifications

**What we're NOT implementing yet:**
- ❌ Stock reservation (Phase 2)
- ❌ Multi-currency (Phase 2)
- ❌ Tax calculation (Phase 2+)
- ❌ Shipping cost calculation (Phase 2+)
- ❌ Discount codes (Phase 3)
- ❌ Partial refunds (Phase 3)
- ❌ Order line item updates (Phase 3)

**What we ARE implementing:**
- ✅ Basic order creation
- ✅ Stripe payment integration
- ✅ Simple stock deduction on payment
- ✅ Order status tracking
- ✅ Order history
- ✅ Stock audit trail

---

## Next Steps

1. Create `orders_module` Go project structure
2. Define Go structs matching these schemas
3. Implement MongoDB connection
4. Create order service with CRUD operations
5. Integrate Stripe SDK
