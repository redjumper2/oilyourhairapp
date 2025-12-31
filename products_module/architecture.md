# Products Module - Architecture Documentation

## Table of Contents

1. [System Overview](#system-overview)
2. [Architecture Principles](#architecture-principles)
3. [Integration with Auth Module](#integration-with-auth-module)
4. [Data Model](#data-model)
5. [API Design](#api-design)
6. [Security & Authentication](#security--authentication)
7. [Multi-Tenancy](#multi-tenancy)
8. [Design Decisions](#design-decisions)

## System Overview

The Products Module is a microservice designed for multi-tenant SaaS product catalog management. It operates independently but integrates with the Auth Module for domain validation and authentication.

### High-Level Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                       Customer Website                       │
│                    (e.g., oilyourhair.com)                   │
│                                                               │
│  ┌─────────────┐          ┌──────────────┐                  │
│  │  Admin UI   │          │ Public Pages │                  │
│  │ (Dashboard) │          │  (Shop page) │                  │
│  └──────┬──────┘          └──────┬───────┘                  │
│         │                        │                           │
│         │ JWT Auth              │ No Auth                   │
└─────────┼────────────────────────┼───────────────────────────┘
          │                        │
          ▼                        ▼
┌─────────────────────────────────────────────────────────────┐
│                     API Gateway / NGINX                      │
│              (Routes to appropriate service)                 │
└─────────────┬───────────────────────────┬───────────────────┘
              │                           │
              ▼                           ▼
    ┌─────────────────┐         ┌──────────────────┐
    │  Auth Module    │         │ Products Module  │
    │  Port: 9090     │         │  Port: 9091      │
    │                 │         │                  │
    │ - User Auth     │◄────────┤ - Product CRUD   │
    │ - Domain Mgmt   │ Reads   │ - Public API     │
    │ - JWT Issuance  │ Domains │ - Variants       │
    └────────┬────────┘         └────────┬─────────┘
             │                           │
             └───────────┬───────────────┘
                         ▼
              ┌─────────────────────┐
              │   MongoDB Server    │
              │                     │
              │  ┌───────────────┐  │
              │  │ auth_module   │  │
              │  │  - users      │  │
              │  │  - domains    │◄─┼─── Cross-DB Read
              │  │  - tokens     │  │
              │  └───────────────┘  │
              │                     │
              │  ┌───────────────┐  │
              │  │products_module│  │
              │  │  - products   │  │
              │  └───────────────┘  │
              └─────────────────────┘
```

## Architecture Principles

### 1. **Microservices Independence**
- Each module can be deployed, scaled, and updated independently
- No direct HTTP calls between services (except shared MongoDB)
- Communication via JWT claims (stateless)

### 2. **Single Source of Truth**
- Domains are managed ONLY in auth_module
- Products module reads domains (read-only, no writes)
- Avoids data duplication and sync issues

### 3. **Multi-Tenancy by Design**
- Every product belongs to a domain
- Complete data isolation between tenants
- Query-level tenant filtering

### 4. **Flexible Data Model**
- Key/value attributes instead of rigid schemas
- Variant-based pricing and inventory
- Schema-less MongoDB for extensibility

## Integration with Auth Module

### Domain Validation (Cross-Database Read)

Products module validates domains by reading from `auth_module.domains`:

```go
// products_module/internal/database/mongo.go

type DB struct {
    Client   *mongo.Client
    Products *mongo.Collection

    // Read-only access to auth_module domains
    AuthDomains *mongo.Collection
}

func NewDB(cfg *config.Config) (*DB, error) {
    client, err := mongo.Connect(ctx, options.Client().ApplyURI(cfg.MongoDB.URI))

    // Products database
    productsDB := client.Database("products_module")

    // Auth database (read-only)
    authDB := client.Database("auth_module")

    return &DB{
        Client:      client,
        Products:    productsDB.Collection("products"),
        AuthDomains: authDB.Collection("domains"), // Read-only
    }, nil
}
```

**Why this works:**
- ✅ No HTTP overhead
- ✅ Single source of truth (auth_module owns domains)
- ✅ Products module is read-only (can't corrupt auth data)
- ✅ Fast domain validation queries

### JWT-Based Authentication

Both modules share the same JWT secret, enabling stateless authentication:

```
┌─────────────────────────────────────────────────────────────┐
│  1. User logs in to Auth Module                             │
│     POST /api/v1/auth/magic-link/request                    │
└─────────────────────┬───────────────────────────────────────┘
                      ▼
┌─────────────────────────────────────────────────────────────┐
│  2. Auth Module issues JWT with claims:                     │
│     {                                                        │
│       "user_id": "123",                                      │
│       "email": "admin@oilyourhair.com",                      │
│       "domain": "oilyourhair.com",                           │
│       "role": "admin",                                       │
│       "exp": 1234567890                                      │
│     }                                                        │
│     Signed with: SHARED_JWT_SECRET                           │
└─────────────────────┬───────────────────────────────────────┘
                      ▼
┌─────────────────────────────────────────────────────────────┐
│  3. User calls Products API with JWT                        │
│     POST /api/v1/products                                   │
│     Authorization: Bearer <JWT>                             │
└─────────────────────┬───────────────────────────────────────┘
                      ▼
┌─────────────────────────────────────────────────────────────┐
│  4. Products Module validates JWT:                          │
│     - Verifies signature (using SHARED_JWT_SECRET)          │
│     - Checks expiration                                     │
│     - Extracts domain and role                              │
│     - Authorizes: role == "admin" for domain                │
└─────────────────────────────────────────────────────────────┘
```

**JWT Middleware (Shared Logic):**

```go
// products_module/internal/middleware/jwt.go

func JWTAuth(jwtSecret string) echo.MiddlewareFunc {
    return func(next echo.HandlerFunc) echo.HandlerFunc {
        return func(c echo.Context) error {
            // Extract token from Authorization header
            token := extractToken(c.Request().Header.Get("Authorization"))

            // Parse and validate
            claims := jwt.MapClaims{}
            _, err := jwt.ParseWithClaims(token, claims, func(t *jwt.Token) (interface{}, error) {
                return []byte(jwtSecret), nil
            })

            if err != nil {
                return c.JSON(401, map[string]string{"error": "Invalid token"})
            }

            // Store claims in context for handlers
            c.Set("user_id", claims["user_id"])
            c.Set("domain", claims["domain"])
            c.Set("role", claims["role"])

            return next(c)
        }
    }
}
```

### Configuration Integration

Both modules must share the JWT secret:

**auth_module/.env:**
```bash
AUTH_JWT_SECRET=my-super-secret-key-12345678901234567890
```

**products_module/.env:**
```bash
PRODUCTS_JWT_SECRET=my-super-secret-key-12345678901234567890  # SAME!
```

## Data Model

### Product Schema

```go
type Product struct {
    ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
    Domain      string            `bson:"domain" json:"domain"`
    Name        string            `bson:"name" json:"name"`
    Description string            `bson:"description" json:"description"`
    BasePrice   float64           `bson:"base_price" json:"base_price"`
    Images      []string          `bson:"images" json:"images"`

    // Flexible key/value attributes
    Attributes  map[string]string `bson:"attributes" json:"attributes"`

    // Variants with pricing and inventory
    Variants    []Variant         `bson:"variants" json:"variants"`

    Active      bool              `bson:"active" json:"active"`
    CreatedAt   time.Time         `bson:"created_at" json:"created_at"`
    UpdatedAt   time.Time         `bson:"updated_at" json:"updated_at"`
}

type Variant struct {
    ID         string            `bson:"id" json:"id"`
    Attributes map[string]string `bson:"attributes" json:"attributes"`
    Price      float64           `bson:"price" json:"price"`
    Stock      int               `bson:"stock" json:"stock"`
    SKU        string            `bson:"sku" json:"sku"`
}
```

### Example Documents

**Product with Variants:**
```json
{
  "_id": "507f1f77bcf86cd799439011",
  "domain": "oilyourhair.com",
  "name": "Coconut Oil",
  "description": "Pure organic coconut oil for hair and skin",
  "base_price": 29.99,
  "images": [
    "https://cdn.oilyourhair.com/products/coconut-oil-500ml.jpg",
    "https://cdn.oilyourhair.com/products/coconut-oil-1l.jpg"
  ],
  "attributes": {
    "category": "oils",
    "type": "coconut",
    "brand": "OilYourHair",
    "organic": "true"
  },
  "variants": [
    {
      "id": "var_co_500ml",
      "attributes": {
        "size": "500ml"
      },
      "price": 29.99,
      "stock": 50,
      "sku": "OYH-CO-500"
    },
    {
      "id": "var_co_1l",
      "attributes": {
        "size": "1L"
      },
      "price": 49.99,
      "stock": 30,
      "sku": "OYH-CO-1000"
    },
    {
      "id": "var_co_1l_org",
      "attributes": {
        "size": "1L",
        "organic": "certified"
      },
      "price": 64.99,
      "stock": 15,
      "sku": "OYH-CO-1000-ORG"
    }
  ],
  "active": true,
  "created_at": "2025-12-30T00:00:00Z",
  "updated_at": "2025-12-30T00:00:00Z"
}
```

### MongoDB Indexes

```javascript
// Products collection indexes
db.products.createIndex({ "domain": 1, "active": 1 })
db.products.createIndex({ "domain": 1, "attributes.category": 1 })
db.products.createIndex({ "domain": 1, "created_at": -1 })
db.products.createIndex({ "variants.sku": 1 }, { unique: true, sparse: true })
```

## API Design

### Admin Endpoints (JWT Required)

**Create Product:**
```http
POST /api/v1/products
Authorization: Bearer <JWT>
Content-Type: application/json

{
  "name": "Argan Oil",
  "description": "Moroccan argan oil",
  "base_price": 39.99,
  "images": ["https://..."],
  "attributes": {
    "category": "oils",
    "type": "argan"
  },
  "variants": [
    {
      "attributes": { "size": "100ml" },
      "price": 39.99,
      "stock": 100,
      "sku": "ARG-100"
    }
  ]
}
```

**List Products (with filters):**
```http
GET /api/v1/products?category=oils&active=true
Authorization: Bearer <JWT>
```

**Update Product:**
```http
PUT /api/v1/products/507f1f77bcf86cd799439011
Authorization: Bearer <JWT>

{
  "base_price": 34.99,
  "variants": [...]
}
```

### Public Endpoints (No Auth)

**List Products for Domain:**
```http
GET /api/v1/public/oilyourhair.com/products?category=oils&limit=20

Response:
{
  "products": [...],
  "total": 45,
  "page": 1,
  "per_page": 20
}
```

**Get Product Details:**
```http
GET /api/v1/public/oilyourhair.com/products/507f1f77bcf86cd799439011

Response:
{
  "id": "507f1f77bcf86cd799439011",
  "name": "Coconut Oil",
  "variants": [...]
}
```

**Search Products:**
```http
GET /api/v1/public/oilyourhair.com/products/search?q=coconut&category=oils
```

## Security & Authentication

### Admin Access Control

```go
func (h *ProductHandler) CreateProduct(c echo.Context) error {
    // JWT middleware already validated token
    userDomain := c.Get("domain").(string)
    userRole := c.Get("role").(string)

    // Check admin role
    if userRole != "admin" {
        return c.JSON(403, map[string]string{
            "error": "Admin access required",
        })
    }

    // Product will be created for user's domain
    product.Domain = userDomain

    // ... create product
}
```

### Public API - No Authentication Required

Public endpoints are read-only and filter by domain:

```go
func (h *PublicHandler) ListProducts(c echo.Context) error {
    domain := c.Param("domain")

    // Verify domain exists and is active
    if !h.verifyDomain(domain) {
        return c.JSON(404, map[string]string{
            "error": "Domain not found",
        })
    }

    // Query products for this domain only
    products, err := h.service.ListByDomain(domain, filters)

    return c.JSON(200, products)
}
```

### CORS Configuration

Products API allows CORS from customer domains:

```go
e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
    AllowOrigins: []string{
        "http://localhost:3000",
        "https://oilyourhair.com",
        "https://www.oilyourhair.com",
        // Add customer domains dynamically
    },
    AllowMethods: []string{"GET", "POST", "PUT", "DELETE"},
    AllowHeaders: []string{"Authorization", "Content-Type"},
}))
```

## Multi-Tenancy

### Tenant Isolation

Every query includes domain filter:

```go
// WRONG - Leaks data across domains
products, err := db.Products.Find(ctx, bson.M{})

// CORRECT - Domain isolation
products, err := db.Products.Find(ctx, bson.M{
    "domain": "oilyourhair.com",
    "active": true,
})
```

### Domain Validation

Before any operation, verify domain exists:

```go
func (s *ProductService) verifyDomain(domain string) error {
    var d models.Domain
    err := s.db.AuthDomains.FindOne(ctx, bson.M{
        "domain": domain,
        "status": "active",
    }).Decode(&d)

    if err != nil {
        return fmt.Errorf("domain not found or inactive")
    }

    return nil
}
```

## Design Decisions

### 1. **Why Cross-Database Read for Domains?**

**Decision:** Products module reads `auth_module.domains` directly

**Alternatives Considered:**
- ❌ HTTP calls to auth module → service dependency
- ❌ Duplicate domains → data sync issues
- ✅ Cross-DB read → no dependency, single source of truth

**Trade-offs:**
- Couples to MongoDB (acceptable for this use case)
- Products can't start if auth DB is down (mitigated by connection retry)
- Simpler than service mesh or event bus

### 2. **Why Key/Value Attributes Instead of Fixed Schema?**

**Decision:** Use `map[string]string` for product attributes

**Rationale:**
- Different product types need different attributes (e.g., oils vs electronics)
- Customers can add custom attributes without schema migration
- Easy filtering: `attributes.category = "oils"`
- MongoDB handles schema-less data efficiently

**Example:**
```json
// Oil product
{
  "attributes": {
    "category": "oils",
    "type": "coconut",
    "organic": "true"
  }
}

// Electronics product (different attributes!)
{
  "attributes": {
    "category": "electronics",
    "brand": "Sony",
    "warranty": "2 years"
  }
}
```

### 3. **Why Variants Instead of Separate Products?**

**Decision:** Variants are sub-entities of a product

**Rationale:**
- SKU management (500ml vs 1L are same product, different sizes)
- Simplified inventory (track stock per variant)
- Better UX (one product page with size/color selectors)
- Flexible pricing per variant combination

### 4. **Why Separate Public API?**

**Decision:** `/api/v1/public/:domain/products` (no auth) vs `/api/v1/products` (admin only)

**Rationale:**
- Customer websites shouldn't need authentication to show products
- Clear separation of concerns (public read vs admin write)
- Different rate limiting policies
- Easier to cache public endpoints

### 5. **Why Image URLs Instead of File Upload?**

**Decision:** Store image URLs, host images externally (Cloudflare R2, S3)

**Rationale:**
- Products module focuses on data, not file storage
- CDN for image optimization and delivery
- Easier to integrate with existing image management tools
- Can add image upload endpoint later if needed

## Future Considerations

### 1. **Event-Driven Architecture**

For more complex scenarios, consider events:

```
Product Created → Event Bus → Webhook to customer
Product Stock Low → Event → Notification service
```

### 2. **Caching Layer**

Add Redis for public API:

```
GET /public/oilyourhair.com/products?category=oils
  → Check Redis cache
  → If miss, query MongoDB
  → Cache result (TTL: 5 minutes)
```

### 3. **Search Engine**

For advanced search, integrate Elasticsearch:

```
Fulltext search, fuzzy matching, faceted filtering
```

### 4. **Product Recommendations**

ML-based product recommendations:

```
"Customers who bought X also bought Y"
```

### 5. **Inventory Sync**

Integrate with fulfillment services:

```
Order placed → Decrement stock → Sync to warehouse
```

## Conclusion

The Products Module follows microservices best practices while maintaining simplicity:

- **Independent deployment** yet integrated with auth
- **Multi-tenant by design** with strict data isolation
- **Flexible data model** for diverse product types
- **Secure authentication** via shared JWT
- **Public API** for easy website integration

This architecture provides a solid foundation for e-commerce product management while remaining extensible for future requirements.
