# Quick Start Guide

Get the auth module running in 5 minutes with Docker.

## Step 1: Setup

```bash
cd auth_module
make setup
```

This creates `.env` and `config.yaml` from examples.

## Step 2: Configure

Edit `.env` and set a strong JWT secret:

```bash
nano .env
```

**Minimum required:**
```env
JWT_SECRET=your-very-strong-secret-key-here-change-this
```

**Optional (for email features):**
```env
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_USER=your-email@gmail.com
SMTP_PASSWORD=your-app-password
EMAIL_FROM=noreply@yourdomain.com
```

## Step 3: Start Services

```bash
make docker-up
```

This starts:
- MongoDB (localhost:27017)
- Auth API (localhost:8080)

Wait for "✅ API is ready" message.

## Step 4: Create First Domain

```bash
make domain-create \
  DOMAIN=oilyourhair.com \
  NAME="Oil Your Hair" \
  EMAIL=admin@oilyourhair.com
```

You'll see:
```
✅ Domain created: oilyourhair.com
✅ Admin invitation created for: admin@oilyourhair.com
📧 Invitation URL: http://localhost:3000/invite?token=abc123...
📱 QR Code: data:image/png;base64,...
⏰ Expires: 2025-12-29T...
```

## Step 5: Test the API

```bash
# Check health
curl http://localhost:8080/health

# Request magic link (replace with your domain)
curl -X POST http://localhost:8080/api/v1/auth/magic-link/request \
  -H "Host: oilyourhair.com" \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com"}'
```

## Common Commands

```bash
# View logs
make docker-logs

# List domains
make domain-list

# Stop services
make docker-down

# Restart services
make docker-restart

# Clean everything
make docker-clean
```

## Debug Mode

Start with Mongo Express (web UI for MongoDB):

```bash
make docker-up-debug
```

Access Mongo Express at http://localhost:8081 (admin/admin)

## Next Steps

1. **Read API documentation:** `cat API.md`
2. **Build a frontend** that calls these endpoints
3. **Configure email** for magic links and invitations
4. **Set up Google OAuth** (optional, see README.md)

## Troubleshooting

**MongoDB not connecting:**
```bash
# Check if MongoDB is healthy
docker ps
make docker-logs-mongo
```

**API not starting:**
```bash
# Check API logs
make docker-logs-api

# Verify JWT_SECRET is set
cat .env | grep JWT_SECRET
```

**Port conflicts:**
```bash
# Check what's using the ports
lsof -i :8080  # API
lsof -i :27017 # MongoDB

# Change ports in docker-compose.yml if needed
```

**Reset everything:**
```bash
make docker-clean
make docker-up
```

## Production Deployment

Before deploying to production:

1. ✅ Change `JWT_SECRET` to a strong random value
2. ✅ Configure real SMTP settings
3. ✅ Set `AUTH_SERVER_ENV=production`
4. ✅ Use managed MongoDB (MongoDB Atlas, etc.)
5. ✅ Enable HTTPS/TLS
6. ✅ Set proper CORS origins
7. ✅ Configure Google OAuth credentials
8. ✅ Review security settings in `config.yaml`

## Development Without Docker

If you prefer local development:

```bash
# Install dependencies
make deps

# Start MongoDB locally
sudo systemctl start mongodb

# Build and run
make build
./auth-module serve
```

See README.md for full documentation.
