# Local Development Guide

Run the auth module code locally with MongoDB in Docker for easy debugging and development.

## Quick Start

```bash
# 1. Download Go dependencies
make deps

# 2. Start MongoDB in Docker + Build + Run server
make dev
```

That's it! Server is running at http://localhost:8080

## Step-by-Step Setup

### 1. Install Go

```bash
# Check if Go is installed
go version

# If not installed (Ubuntu/Debian):
sudo apt update
sudo apt install golang-go

# Verify
go version  # Should show 1.21+
```

### 2. Download Dependencies

```bash
cd auth_module
make deps
```

### 3. Start MongoDB Only

```bash
# Start MongoDB in Docker
make dev-db

# Or with Mongo Express web UI for debugging
make dev-db-debug
```

MongoDB is now running on `localhost:27017`

### 4. Build the Application

```bash
make build
```

This creates the `auth-module` binary.

### 5. Run the Server

```bash
make run
```

Server starts on http://localhost:8080 with `config.dev.yaml`

## Development Workflow

### Option 1: Automated (Recommended)

```bash
# Stop everything, start fresh MongoDB, build, and run
make dev
```

### Option 2: Manual Control

```bash
# Terminal 1: Start MongoDB
make dev-db

# Terminal 2: Build and run
make build
make run

# Make code changes...
# Ctrl+C to stop server
make run  # Run again to test changes
```

### Option 3: Run Without Make

```bash
# Start MongoDB
docker compose -f docker-compose.dev.yml up -d

# Build
go build -o auth-module .

# Run with specific config
./auth-module serve --config=config.dev.yaml

# Or with CLI commands
./auth-module domain list --config=config.dev.yaml
```

## Managing Domains Locally

```bash
# Build first
make build

# Create domain
make domain-create-local \
  DOMAIN=testdomain.com \
  NAME="Test Domain" \
  EMAIL=admin@testdomain.com

# List domains
make domain-list-local

# Delete domain
make domain-delete-local DOMAIN=testdomain.com
```

## Debugging

### VSCode

Create `.vscode/launch.json`:

```json
{
  "version": "0.2.0",
  "configurations": [
    {
      "name": "Launch Server",
      "type": "go",
      "request": "launch",
      "mode": "auto",
      "program": "${workspaceFolder}",
      "args": ["serve", "--config=config.dev.yaml"],
      "env": {},
      "preLaunchTask": ""
    },
    {
      "name": "Domain Create",
      "type": "go",
      "request": "launch",
      "mode": "auto",
      "program": "${workspaceFolder}",
      "args": [
        "domain", "create",
        "--config=config.dev.yaml",
        "--domain=test.com",
        "--name=Test",
        "--admin-email=admin@test.com"
      ]
    }
  ]
}
```

Set breakpoints in your code, then press F5 to debug!

### GoLand / IntelliJ

1. Right-click `main.go` → Debug 'go build'
2. Edit configuration → Add program arguments: `serve --config=config.dev.yaml`
3. Set breakpoints and debug

### Delve (Command Line)

```bash
# Install delve
go install github.com/go-delve/delve/cmd/dlv@latest

# Debug
dlv debug . -- serve --config=config.dev.yaml

# Set breakpoint
(dlv) break main.main
(dlv) continue
```

## Testing Locally

```bash
# Start server in background
make dev &

# Wait for startup
sleep 3

# Create test domain
make domain-create-local \
  DOMAIN=localhost \
  NAME="Local" \
  EMAIL=admin@localhost

# Test API
curl -X POST http://localhost:8080/api/v1/auth/magic-link/request \
  -H "Host: localhost" \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com"}'

# Stop server
killall auth-module
```

## Viewing MongoDB Data

### Option 1: Mongo Shell

```bash
# Connect to MongoDB
docker exec -it auth-module-mongodb-dev mongosh auth_module

# List collections
show collections

# View domains
db.domains.find().pretty()

# View users
db.users.find().pretty()

# Exit
exit
```

### Option 2: Mongo Express (Web UI)

```bash
# Start with debug profile
make dev-db-debug

# Open browser
open http://localhost:8081
```

Login: `admin` / `admin`

## Configuration

### Using config.dev.yaml

The default local config file with sensible defaults:

```yaml
server:
  port: 8080
mongodb:
  uri: mongodb://localhost:27017
  database: auth_module
jwt:
  secret: dev-secret-key-for-local-testing-only
```

### Environment Variable Overrides

```bash
# Override JWT secret
export AUTH_JWT_SECRET=my-custom-secret

# Override MongoDB URI
export AUTH_MONGODB_URI=mongodb://localhost:27017

# Run
make run
```

Environment variables prefixed with `AUTH_` override config file values.

## Hot Reload (Optional)

For automatic rebuilds on code changes:

```bash
# Install air
go install github.com/cosmtrek/air@latest

# Create .air.toml
cat > .air.toml <<EOF
root = "."
tmp_dir = "tmp"

[build]
  cmd = "go build -o ./tmp/auth-module ."
  bin = "./tmp/auth-module serve --config=config.dev.yaml"
  include_ext = ["go", "yaml"]
  exclude_dir = ["tmp", "vendor"]
  delay = 1000
EOF

# Run with hot reload
air
```

Now code changes automatically rebuild and restart the server!

## Troubleshooting

### "connection refused" to MongoDB

```bash
# Check if MongoDB is running
docker ps | grep mongodb-dev

# If not, start it
make dev-db

# Check logs
docker logs auth-module-mongodb-dev
```

### "command not found: go"

```bash
# Install Go
sudo apt update
sudo apt install golang-go

# Verify
go version
```

### Build errors

```bash
# Clean and rebuild
make clean
make deps
make build
```

### Port 8080 already in use

```bash
# Find what's using port 8080
lsof -i :8080

# Kill it or change port in config.dev.yaml:
# server:
#   port: 8081
```

### MongoDB port 27017 in use

```bash
# Check what's using it
lsof -i :27017

# Stop local MongoDB if running
sudo systemctl stop mongod

# Or change port in docker-compose.dev.yml
```

## Cleanup

```bash
# Stop MongoDB
make dev-db-down

# Clean build artifacts
make clean

# Clean everything including database
make clean-all
```

## Production Build

Before deploying:

```bash
# Build optimized binary
CGO_ENABLED=0 go build -ldflags="-s -w" -o auth-module .

# Test
./auth-module serve --config=config.yaml
```

## Next Steps

- Make code changes in your editor
- Set breakpoints for debugging
- Test API endpoints with curl or Postman
- View data in Mongo Express
- Write tests in `*_test.go` files
- Run `make test` to execute tests

Happy coding! 🚀
