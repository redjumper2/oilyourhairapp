# Auth Module

Multi-tenant authentication module with domain isolation, built with Go Echo framework, MongoDB, Cobra CLI, and Viper configuration.

## Features

- **Multi-tenant architecture** - Complete domain isolation (share nothing)
- **Flexible authentication** - Google OAuth and magic link email authentication
- **Invitation system** - Email invitations and QR code generation with tracking
- **Role-based access control** - Admin, editor, viewer, customer roles with granular permissions
- **Permission system** - 16 permissions across 6 categories with centralized registry
- **CLI management** - Cobra CLI for domain and user management
- **Hierarchical configuration** - Viper with YAML/JSON/environment variables

## Project Structure

```
auth_module/
├── cmd/
│   ├── root.go           # Cobra root command
│   ├── serve.go          # Start HTTP API server
│   ├── domain.go         # Domain management commands
│   └── permissions.go    # Permission management commands
├── internal/
│   ├── models/           # MongoDB models
│   │   ├── domain.go
│   │   ├── user.go
│   │   ├── invitation.go
│   │   └── permissions.go # Permission registry
│   ├── handlers/         # Echo HTTP handlers
│   │   ├── auth.go
│   │   ├── oauth.go
│   │   └── admin.go
│   ├── middleware/       # Auth middleware
│   ├── services/         # Business logic
│   │   ├── invitation.go
│   │   └── email.go
│   └── database/         # MongoDB connection
│       └── database.go
├── pkg/
│   └── utils/            # Utilities (JWT, etc.)
├── config/
│   └── config.go         # Viper configuration
├── docs/                 # Documentation
│   ├── API.md            # Complete API reference
│   ├── TESTING.md        # Testing guide
│   ├── GOOGLE_OAUTH_SETUP.md
│   └── LOCAL_DEV.md
├── config.yaml.example   # Example configuration
├── .env.example          # Docker environment variables
├── docker-compose.yml
├── Dockerfile
├── Makefile
├── go.mod
└── main.go
```

## Getting Started

### Prerequisites

- **Docker & Docker Compose** (recommended)
- OR Go 1.21+ and MongoDB 4.4+ (for local development)

### Quick Start with Docker (Recommended)

```bash
# 1. Initial setup (creates .env and config.yaml)
make setup

# 2. Edit .env file with your settings (especially JWT_SECRET!)
nano .env

# 3. Start all services (MongoDB + API)
make docker-up

# 4. Create your first domain
make domain-create DOMAIN=oilyourhair.com NAME="Oil Your Hair" EMAIL=admin@oilyourhair.com

# 5. Check health
make health
```

That's it! API is running at http://localhost:8080

### Configuration

**With Docker (uses environment variables):**
```bash
# Copy .env.example to .env
cp .env.example .env

# Edit .env with your settings
nano .env
```

**Without Docker (uses config.yaml):**
```bash
# Copy example config
cp config.yaml.example config.yaml

# Edit config.yaml
nano config.yaml
```

Required settings:
- `JWT_SECRET` / `jwt.secret` - Strong secret key for JWT tokens
- MongoDB connection (auto-configured in Docker)
- SMTP settings for email (optional for testing)
- Google OAuth credentials (optional)

Environment variables override config (prefix with `AUTH_`):
```bash
export AUTH_JWT_SECRET=your-secret-key
export AUTH_MONGODB_URI=mongodb://localhost:27017
```

## CLI Usage

### Domain Management

**Create a new domain:**
```bash
./auth-module domain create \
  --domain=oilyourhair.com \
  --name="Oil Your Hair" \
  --admin-email=admin@oilyourhair.com

# Or with Makefile:
make domain-create DOMAIN=oilyourhair.com NAME="Oil Your Hair" EMAIL=admin@oilyourhair.com
```

This will:
1. Create the domain in MongoDB
2. Set up default settings (Google + magic link auth enabled)
3. Create an admin invitation with QR code
4. Print the invitation URL and QR code data URL

**List all domains:**
```bash
./auth-module domain list
# Or: make domain-list
```

**Delete a domain:**
```bash
./auth-module domain delete --domain=oilyourhair.com
```
⚠️ This will delete the domain and all its users!

### Permission Management

**List all available permissions:**
```bash
./auth-module permissions list
```

**View permissions by role:**
```bash
./auth-module permissions roles
```

### Start API Server

```bash
./auth-module serve
```

Server will start on `http://localhost:8080` (configurable in config.yaml)

Health check: `GET /health`

## API Endpoints

All API endpoints are fully implemented and documented. See **[docs/API.md](docs/API.md)** for complete API reference.

### Authentication Endpoints (Public)

- `POST /api/v1/auth/magic-link/request` - Request magic link
- `GET /api/v1/auth/magic-link/verify` - Verify magic link token and get JWT
- `GET /api/v1/auth/invitation/verify` - Check invitation validity
- `POST /api/v1/auth/invitation/accept` - Accept invitation and create user
- `GET /api/v1/auth/google` - Initiate Google OAuth
- `GET /api/v1/auth/google/callback` - OAuth callback (browser redirect)
- `GET /api/v1/auth/google/callback/json` - OAuth callback (JSON)
- `GET /api/v1/auth/me` - Get current user (requires JWT)

### Admin APIs (Protected)

**Domain Settings:**
- `GET /api/v1/admin/domain/settings` - Get domain settings
- `PUT /api/v1/admin/domain/settings` - Update domain settings

**User Management:**
- `GET /api/v1/admin/users` - List users
- `POST /api/v1/admin/users/invite` - Invite user with QR code
- `PUT /api/v1/admin/users/:id` - Update user role/permissions
- `DELETE /api/v1/admin/users/:id` - Soft delete user

**Permissions:**
- `GET /api/v1/admin/permissions` - Get all permissions (grouped)
- `GET /api/v1/admin/permissions/roles` - Get permissions by role

## Architecture

### Domain Isolation

Each domain is completely isolated:
- Separate user records per domain
- Domain determined by HTTP `Host` header
- JWT tokens include domain claim
- All queries filtered by domain

### Permission System

Centralized permission registry with 16 permissions across 6 categories:
- **Domain Management**: settings.read, settings.write
- **User Management**: users.read, users.write, users.delete, users.invite
- **Product Management**: products.read, products.write
- **Order Management**: orders.read, orders.write
- **Inventory Management**: inventory.read, inventory.write
- **Shopping Cart**: cart.read, cart.write

**Four default roles:**
- **admin** (12 permissions) - Full domain control
- **editor** (5 permissions) - Product and inventory management
- **viewer** (3 permissions) - Read-only access
- **customer** (4 permissions) - Shopping and orders

### Authentication Flow

**Magic Link:**
1. User enters email on `oilyourhair.com`
2. System extracts domain from `Host` header
3. Generates token, sends email with link
4. User clicks link, verifies token
5. Issues JWT with domain in claims

**Google OAuth:**
1. User clicks "Sign in with Google" on `oilyourhair.com`
2. OAuth callback includes domain context
3. System checks if user exists for this domain
4. Creates user if new, issues JWT

### Invitation System

**Three invitation types:**

1. **Email invitation** - Specific user, sends email with link + QR code
2. **QR code (single-use)** - Generate QR for specific user, no email
3. **QR code (multi-use)** - Promotional QR code, trackable, reusable

**Expiry defaults (configurable):**
- Email invitations: 24 hours
- Single-use QR: 72 hours
- Multi-use promotional: 30 days

**Tracking metadata:**
- Promo codes
- Source attribution (booth, instagram, email)
- Referrer tracking
- Custom fields

## MongoDB Collections

### domains
- Whitelist of registered domains
- Settings (allowed auth providers, default role)
- Branding (logo, colors, company name)

### users
- User records (unique per email+domain)
- Role and permissions
- Soft delete support

### invitations
- Pending invitations
- QR code metadata
- Auto-deleted after expiry (TTL index)

### invitation_logs
- Audit trail for all invitations
- Analytics data (promo codes, sources)

### magic_link_tokens
- Temporary auth tokens
- Auto-deleted after 15 minutes (TTL index)

## Security

- Domain whitelist validation
- JWT domain verification (prevents cross-domain token reuse)
- Single-use tokens deleted after verification
- Configurable token expiry
- Soft delete for users (audit trail)
- Admin lockout prevention (can't delete last admin)

## Documentation

- **[API Reference](docs/API.md)** - Complete API documentation with examples
- **[Testing Guide](docs/TESTING.md)** - Step-by-step testing instructions
- **[Google OAuth Setup](docs/GOOGLE_OAUTH_SETUP.md)** - OAuth configuration
- **[Local Development](docs/LOCAL_DEV.md)** - Local dev setup without Docker

## Development

**Run in development mode:**
```bash
go run main.go serve
```

**Build for production:**
```bash
CGO_ENABLED=0 go build -o auth-module -ldflags="-s -w" .
```

**Run tests:**
```bash
# Run all tests
go test ./...

# Integration testing
# See docs/TESTING.md for complete testing guide
make health
make domain-create DOMAIN=test.com NAME="Test" EMAIL=admin@test.com
```

## Testing

See **[docs/TESTING.md](docs/TESTING.md)** for comprehensive testing guide covering:
- Magic link authentication
- Invitation system (email + QR codes)
- Google OAuth flow
- Admin APIs
- Permission system
- Domain isolation

Quick test:
```bash
# Start services
make docker-up

# Create domain
make domain-create DOMAIN=test.com NAME="Test" EMAIL=admin@test.com

# Health check
curl http://localhost:8080/health | jq
```

## License

Proprietary
