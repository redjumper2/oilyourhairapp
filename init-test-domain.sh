#!/bin/bash

# Initialize test domain for testing
# This script creates testdomain.com and an admin invitation

echo "🔧 Initializing test domain..."
echo ""

# Wait for API to be healthy
echo "⏳ Waiting for auth API to be ready..."
until curl -sf http://localhost:9090/health > /dev/null; do
    echo "   Waiting for API..."
    sleep 2
done
echo "✅ API is ready!"
echo ""

# Create test domain
echo "📝 Creating testdomain.com..."
docker-compose exec -T auth-api ./auth-module domain create \
  --domain=testdomain.com \
  --name="Test Domain Store" \
  --admin-email=admin@testdomain.com

echo ""
echo "✅ Test domain initialized!"
echo ""
echo "🚀 Quick Start:"
echo "   1. Visit: http://localhost:8000 (test domain)"
echo "   2. Click 'Login' button"
echo "   3. Complete auth flow"
echo "   4. You'll be redirected back with JWT token"
echo ""
echo "📋 Services:"
echo "   - Test Domain:  http://localhost:8000"
echo "   - Auth Portal:  http://localhost:5173"
echo "   - Auth API:     http://localhost:9090"
echo "   - MongoDB:      localhost:27017"
echo ""
