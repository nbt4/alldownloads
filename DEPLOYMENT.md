# 🚀 AllDownloads - Production Deployment Guide

**Self-hosted solution for fetching latest official download links for OS ISOs and desktop applications**

## 📋 Quick Start

### Prerequisites
- Docker & Docker Compose installed
- 2GB+ RAM available
- Internet connection

### 1-Command Deployment
```bash
# Download and run
curl -L https://raw.githubusercontent.com/your-repo/alldownloads/main/docker-compose.prod.yml | docker compose -f - up -d
```

## 🏗️ Manual Setup

### 1. Create Project Directory
```bash
mkdir alldownloads && cd alldownloads
```

### 2. Download Configuration Files
```bash
# Download production docker-compose
wget https://raw.githubusercontent.com/your-repo/alldownloads/main/docker-compose.prod.yml

# Download environment template
wget https://raw.githubusercontent.com/your-repo/alldownloads/main/.env.example -O .env
```

### 3. Configure Environment
```bash
# Edit configuration
nano .env

# Key settings to change:
# - AUTH_TOKEN: Generate a secure random token
# - DOMAIN: Your domain name
# - BASE_URL: Your public URL
# - CORS_ORIGINS: Allowed frontend origins
```

### 4. Deploy
```bash
# Start all services
docker compose -f docker-compose.prod.yml up -d

# Check status
docker compose -f docker-compose.prod.yml ps

# View logs
docker compose -f docker-compose.prod.yml logs -f
```

## 🌐 Service URLs

- **API**: `http://localhost:8080`
- **Frontend**: `http://localhost:3000`
- **Health Check**: `http://localhost:8080/api/health`
- **Metrics**: `http://localhost:8080/metrics`

## 📦 Available Docker Images

All images are available on Docker Hub under `nbt4/`:

- `nbt4/alldownloads-api:latest` - REST API service
- `nbt4/alldownloads-worker:latest` - Background worker

## 🔧 Configuration Options

### Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `AUTH_TOKEN` | - | **Required** API authentication token |
| `DOMAIN` | localhost | Your domain name |
| `API_PORT` | 8080 | API service port |
| `UI_PORT` | 3000 | Frontend port |
| `BASE_URL` | http://localhost | Public base URL |
| `CORS_ORIGINS` | localhost origins | Comma-separated allowed origins |
| `REFRESH_CRON` | @every 6h | Worker refresh schedule |
| `HTTP_TIMEOUT` | 15s | HTTP request timeout |
| `MAX_CONCURRENT_FETCHES` | 6 | Max parallel fetches |
| `RATE_LIMIT_REQUESTS_PER_MINUTE` | 60 | API rate limit |
| `LOG_LEVEL` | info | Logging level (debug/info/warn/error) |

### Cron Schedule Examples
```bash
REFRESH_CRON=@every 1h        # Every hour
REFRESH_CRON=@every 30m       # Every 30 minutes
REFRESH_CRON=0 */6 * * *      # Every 6 hours (cron format)
REFRESH_CRON=0 0 * * *        # Daily at midnight
```

## 🏭 Production Setup

### Reverse Proxy (Nginx)
```nginx
server {
    listen 80;
    server_name yourdomain.com;

    # Frontend
    location / {
        proxy_pass http://localhost:3000;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
    }

    # API
    location /api/ {
        proxy_pass http://localhost:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
    }
}
```

### SSL with Let's Encrypt
```bash
# Install certbot
sudo apt install certbot python3-certbot-nginx

# Get certificate
sudo certbot --nginx -d yourdomain.com

# Update environment
echo "BASE_URL=https://yourdomain.com" >> .env
echo "CORS_ORIGINS=https://yourdomain.com" >> .env
```

### Resource Requirements

**Minimum:**
- 1 CPU core
- 2GB RAM
- 5GB disk space

**Recommended:**
- 2 CPU cores
- 4GB RAM
- 20GB disk space

## 📊 Monitoring

### Health Checks
```bash
# API health
curl http://localhost:8080/api/health

# Database connection
curl http://localhost:8080/api/health/db

# Cache connection
curl http://localhost:8080/api/health/cache
```

### Metrics
Prometheus metrics available at `/metrics` endpoint:
- HTTP request duration/count
- Active connections
- Database connection pool stats
- Job queue metrics

### Logs
```bash
# All services
docker compose logs -f

# Specific service
docker compose logs -f api
docker compose logs -f worker

# Follow with timestamps
docker compose logs -f --timestamps
```

## 🔄 Updates

### Update Images
```bash
# Pull latest images
docker compose pull

# Restart with new images
docker compose up -d

# Clean old images
docker image prune
```

### Database Migrations
```bash
# Backup database first
docker exec alldownloads-db pg_dump -U alldl alldownloads > backup.sql

# Apply migrations (if any)
docker compose restart api
```

## 🛠️ Troubleshooting

### Common Issues

**API not responding:**
```bash
docker compose logs api
# Check CORS_ORIGINS configuration
# Verify AUTH_TOKEN is set
```

**Worker not fetching:**
```bash
docker compose logs worker
# Check REFRESH_CRON schedule
# Verify internet connectivity
```

**Database connection errors:**
```bash
docker compose logs db
# Check disk space
# Verify database health
```

### Reset Everything
```bash
# Stop and remove all
docker compose down -v

# Remove all data
docker volume prune

# Start fresh
docker compose up -d
```

## 🌍 Supported Software

**Browsers:** Chrome, Firefox, Brave
**Editors:** VS Code, Notepad++
**Operating Systems:** Ubuntu, Debian, Arch Linux, Kali Linux, Windows 11
**Communication:** Telegram, WhatsApp
**Tools:** PowerShell, Tailscale, Nextcloud, Termius

## 📝 API Usage

### Get All Products
```bash
curl "http://localhost:8080/api/products"
```

### Get Product Downloads
```bash
curl "http://localhost:8080/api/products/chrome"
```

### Authentication
```bash
curl -H "Authorization: Bearer YOUR_TOKEN" "http://localhost:8080/api/admin/refresh"
```

## 📞 Support

- **Issues**: Report bugs and request features
- **Documentation**: Full API documentation available
- **Community**: Join discussions

---

**Made with ❤️ for the developer community**