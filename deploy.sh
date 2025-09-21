#!/bin/bash

# 🚀 AllDownloads - One-Click Deployment Script
# This script sets up AllDownloads on any Docker-enabled machine

set -e

echo "🚀 Starting AllDownloads deployment..."

# Check if Docker is installed
if ! command -v docker &> /dev/null; then
    echo "❌ Docker is not installed. Please install Docker first:"
    echo "   https://docs.docker.com/get-docker/"
    exit 1
fi

# Check if Docker Compose is available
if ! docker compose version &> /dev/null; then
    echo "❌ Docker Compose is not available. Please install Docker Compose:"
    echo "   https://docs.docker.com/compose/install/"
    exit 1
fi

# Create project directory
PROJECT_DIR="alldownloads"
if [ -d "$PROJECT_DIR" ]; then
    echo "📁 Directory $PROJECT_DIR already exists"
    read -p "Do you want to continue? This will overwrite existing files. (y/N): " -n 1 -r
    echo
    if [[ ! $REPLY =~ ^[Yy]$ ]]; then
        exit 1
    fi
else
    echo "📁 Creating project directory: $PROJECT_DIR"
    mkdir -p "$PROJECT_DIR"
fi

cd "$PROJECT_DIR"

# Download production compose file
echo "📥 Downloading production configuration..."
curl -sSL -o docker-compose.yml https://raw.githubusercontent.com/nbt4/alldownloads/main/docker-compose.prod.yml

# Create environment file if it doesn't exist
if [ ! -f ".env" ]; then
    echo "⚙️  Creating environment configuration..."
    cat > .env << 'EOF'
# AllDownloads Configuration
DOMAIN=localhost
API_PORT=8080
UI_PORT=3000
AUTH_TOKEN=alldownloads-$(openssl rand -hex 16 2>/dev/null || echo "change-this-secure-token")
BASE_URL=http://localhost
CORS_ORIGINS=http://localhost:3000,http://localhost:8080
REFRESH_CRON=@every 6h
HTTP_TIMEOUT=15s
MAX_CONCURRENT_FETCHES=6
RATE_LIMIT_REQUESTS_PER_MINUTE=60
LOG_LEVEL=info
LOG_FORMAT=json
EOF
    echo "✅ Created .env file with default configuration"
    echo "💡 You can edit .env to customize settings"
else
    echo "⚙️  Using existing .env file"
fi

# Pull latest images
echo "📦 Pulling latest Docker images..."
docker compose pull

# Start services
echo "🔄 Starting AllDownloads services..."
docker compose up -d

# Wait for services to be healthy
echo "⏳ Waiting for services to start..."
sleep 10

# Check service status
echo "📊 Checking service status..."
if docker compose ps | grep -q "Up"; then
    echo "✅ AllDownloads is running successfully!"
    echo ""
    echo "🌐 Access your AllDownloads installation:"
    echo "   Frontend: http://localhost:$(grep UI_PORT .env | cut -d'=' -f2 || echo 3000)"
    echo "   API:      http://localhost:$(grep API_PORT .env | cut -d'=' -f2 || echo 8080)"
    echo "   Health:   http://localhost:$(grep API_PORT .env | cut -d'=' -f2 || echo 8080)/api/health"
    echo ""
    echo "📚 Useful commands:"
    echo "   View logs:    docker compose logs -f"
    echo "   Stop:         docker compose down"
    echo "   Update:       docker compose pull && docker compose up -d"
    echo "   Reset data:   docker compose down -v"
    echo ""
    echo "🔧 Configuration file: .env"
    echo "📖 Documentation: https://github.com/nbt4/alldownloads"
else
    echo "❌ Some services failed to start. Check logs:"
    echo "   docker compose logs"
    exit 1
fi