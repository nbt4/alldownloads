#!/bin/bash

# 🐛 AllDownloads Debug Deployment Script
# Löst häufige 403 Forbidden Probleme

echo "🔍 Starting AllDownloads debug deployment..."

# Check Docker permissions
echo "👤 Checking Docker permissions..."
if ! docker ps >/dev/null 2>&1; then
    echo "❌ Docker permission denied. Trying with sudo..."
    DOCKER_CMD="sudo docker"
    COMPOSE_CMD="sudo docker compose"
else
    echo "✅ Docker permissions OK"
    DOCKER_CMD="docker"
    COMPOSE_CMD="docker compose"
fi

# Create proxy network if it doesn't exist
echo "🌐 Creating proxy network..."
$DOCKER_CMD network create proxy 2>/dev/null || echo "Network proxy already exists"

# Stop any existing containers
echo "🛑 Stopping existing containers..."
$COMPOSE_CMD down 2>/dev/null || true

# Check if .env exists, create debug version
if [ ! -f ".env" ]; then
    echo "⚙️ Creating debug .env file..."
    cp .env.debug .env 2>/dev/null || cat > .env << 'EOF'
# Debug Configuration
DOMAIN=localhost
API_PORT=9780
UI_PORT=9779
AUTH_TOKEN=debug-token-12345
BASE_URL=http://localhost
CORS_ORIGINS=http://localhost:9779,http://127.0.0.1:9779,http://0.0.0.0:9779
REFRESH_CRON=@every 6h
HTTP_TIMEOUT=15s
MAX_CONCURRENT_FETCHES=6
RATE_LIMIT_REQUESTS_PER_MINUTE=1000
LOG_LEVEL=debug
LOG_FORMAT=json
EOF
fi

# Pull latest images
echo "📦 Pulling latest images..."
$COMPOSE_CMD pull

# Start with debug logging
echo "🚀 Starting AllDownloads with debug configuration..."
$COMPOSE_CMD up -d

# Wait for services
echo "⏳ Waiting for services to start..."
sleep 15

# Check service status
echo "📊 Checking service status..."
$COMPOSE_CMD ps

# Test connectivity
echo "🧪 Testing connectivity..."
echo "Testing API health..."
curl -f http://localhost:9780/api/health 2>/dev/null && echo "✅ API accessible" || echo "❌ API not accessible"

echo "Testing UI..."
curl -f http://localhost:9779 2>/dev/null && echo "✅ UI accessible" || echo "❌ UI not accessible"

# Show logs if there are issues
echo "📋 Recent logs:"
echo "=== API Logs ==="
$COMPOSE_CMD logs --tail 5 api
echo "=== UI Logs ==="
$COMPOSE_CMD logs --tail 5 ui

echo ""
echo "🌐 Access URLs:"
echo "   Frontend: http://localhost:9779"
echo "   API:      http://localhost:9780"
echo "   Health:   http://localhost:9780/api/health"
echo ""
echo "🔧 If still getting 403, check logs with:"
echo "   $COMPOSE_CMD logs -f"