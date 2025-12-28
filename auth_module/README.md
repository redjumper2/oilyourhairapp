# Auth Module

Multi-tenant authentication module with domain isolation, built with Go Echo framework, MongoDB, Cobra CLI, and Viper configuration.

## Features

- **Multi-tenant architecture** - Complete domain isolation (share nothing)
- **Flexible authentication** - Google OAuth and magic link email authentication
- **Invitation system** - Email invitations and QR code generation with tracking
- **Role-based access control** - Admin, editor, viewer, customer roles with permissions
- **CLI management** - Cobra CLI for domain and user management
- **Hierarchical configuration** - Viper with YAML/JSON/environment variables

## Project Structure

```
auth_module/
├── cmd/
│   ├── root.go           # Cobra root command
│   ├── serve.go          # Start HTTP API server
│   └── domain.go         # Domain management commands
├── internal/
│   ├── models/           # MongoDB models
│   │   ├── domain.go
│   │   ├── user.go
│   │   └── invitation.go
│   ├── handlers/         # Echo HTTP handlers
│   ├── middleware/       # Auth middleware
│   ├── services/         # Business logic
│   │   └── invitation.go
│   └── database/         # MongoDB connection
│       └── database.go
├── pkg/
│   └── utils/            # Utilities (JWT, etc.)
├── config/
│   └── config.go         # Viper configuration
├── config.yaml.example   # Example configuration
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
```

This will:
1. Create the domain in MongoDB
2. Set up default settings (Google + magic link auth enabled)
3. Create an admin invitation with QR code
4. Print the invitation URL and QR code data URL

**List all domains:**
```bash
./auth-module domain list
```

**Delete a domain:**
```bash
./auth-module domain delete --domain=oilyourhair.com
```
⚠️ This will delete the domain and all its users!

### Start API Server

```bash
./auth-module serve
```

Server will start on `http://localhost:8080` (configurable in config.yaml)

Health check: `GET /health`

## API Endpoints

### Authentication (Coming Soon)

- `POST /api/v1/auth/magic-link/request` - Request magic link
- `GET /api/v1/auth/magic-link/verify` - Verify magic link token
- `GET /api/v1/auth/google` - Initiate Google OAuth
- `GET /api/v1/auth/google/callback` - OAuth callback
- `GET /api/v1/auth/me` - Get current user

### Admin APIs (Coming Soon)

- `GET /api/v1/admin/domain/settings` - Get domain settings
- `PUT /api/v1/admin/domain/settings` - Update domain settings
- `GET /api/v1/admin/users` - List users
- `POST /api/v1/admin/users/invite` - Invite user
- `PUT /api/v1/admin/users/:id` - Update user
- `DELETE /api/v1/admin/users/:id` - Soft delete user

## Architecture

### Domain Isolation

Each domain is completely isolated:
- Separate user records per domain
- Domain determined by HTTP `Host` header
- JWT tokens include domain claim
- All queries filtered by domain

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
go test ./...
```

## Next Steps

- [ ] Implement magic link authentication handlers
- [ ] Implement Google OAuth with Goth
- [ ] Implement admin APIs
- [ ] Add JWT middleware
- [ ] Email service for sending magic links/invitations
- [ ] Frontend integration examples
- [ ] API documentation (Swagger)

## License

Proprietary
