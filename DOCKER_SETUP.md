# Complete Multi-Tenant Auth System - Docker Setup

End-to-end authentication system with API, UI portal, and test domain.

## Architecture

```
┌──────────────┐
│ Test Domain  │  http://localhost:8000 (Sample customer site)
└──────┬───────┘
       │ Redirects to auth portal
       ▼
┌──────────────┐
│  Auth Portal │  http://localhost:5173 (Svelte UI)
│  (auth_ui)   │  Login, Invitations, OAuth
└──────┬───────┘
       │ Calls API
       ▼
┌──────────────┐
│   Auth API   │  http://localhost:9090 (Go backend)
│ (auth_module)│  JWT, Users, Domains
└──────┬───────┘
       │
       ▼
┌──────────────┐
│   MongoDB    │  localhost:27017
└──────────────┘
```

## Quick Start

```bash
# 1. Start all services
docker-compose up -d

# 2. Wait ~30 seconds, then initialize test domain
./init-test-domain.sh

# 3. Visit test domain
open http://localhost:8000

# 4. Click Login and test the auth flow!
```

## Services

| Service | Port | Description |
|---------|------|-------------|
| **auth-api** | 9090 | Go backend API |
| **auth-ui** | 5173 | Svelte auth portal |
| **test-domain** | 8000 | Sample customer site |
| **mongodb** | 27017 | Database |

## Complete Test Flow

1. **Visit**: http://localhost:8000
   - See products (prices hidden)
   - Click "Login" button

2. **Redirected to**: http://localhost:5173/login?domain=testdomain.com&redirect=http://localhost:8000
   - Shows testdomain.com branding
   - Enter email for magic link

3. **Get magic link token**:
   ```bash
   docker-compose exec mongodb mongosh auth_module --eval \
     "db.magic_link_tokens.find().sort({created_at:-1}).limit(1).pretty()"
   ```

4. **Verify magic link**:
   ```bash
   open "http://localhost:5173/verify?token=TOKEN_HERE&redirect=http://localhost:8000"
   ```

5. **Redirected back** to: http://localhost:8000#token=JWT
   - Token extracted and stored
   - Now logged in!
   - Prices visible

## CLI Management

```bash
# List domains
docker-compose exec auth-api ./auth-module domain list

# Create invitation
docker-compose exec auth-api ./auth-module invite create \
  --email=user@example.com \
  --role=editor \
  --domain=testdomain.com

# List users
docker-compose exec auth-api ./auth-module users list --domain=testdomain.com
```

## Development

### Auth UI changes
```bash
cd auth_ui_module
npm install
npm run dev  # Runs outside Docker for hot-reload
```

### Auth API changes
```bash
cd auth_module
docker-compose up -d --build auth-api
```

### Test Domain changes
Edit `test_domain/` files, refresh browser

## Logs

```bash
docker-compose logs -f auth-api
docker-compose logs -f auth-ui
docker-compose logs -f test-domain
```

## Cleanup

```bash
# Stop everything
docker-compose down

# Remove data
docker-compose down -v
```

## Why Port 9090?

Port 8080 is used by oilyourhair.com production nginx with cloudflared.
