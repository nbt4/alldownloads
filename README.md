# 🚀 AllDownloads

[![Docker Pulls API](https://img.shields.io/docker/pulls/nbt4/alldownloads-api)](https://hub.docker.com/r/nbt4/alldownloads-api)
[![Docker Pulls Worker](https://img.shields.io/docker/pulls/nbt4/alldownloads-worker)](https://hub.docker.com/r/nbt4/alldownloads-worker)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)
[![Go Version](https://img.shields.io/badge/go-1.22+-blue.svg)](https://golang.org/)
[![Docker](https://img.shields.io/badge/docker-ready-blue.svg)](https://www.docker.com/)

> Self-hosted solution that fetches and serves the latest official download links for OS ISOs and common desktop applications with a modern, beautiful web interface.

## ✨ Features

- **🔒 Official Sources Only**: Fetches downloads directly from vendor websites and mirrors
- **🔄 Auto-Updated**: Scheduled refresh jobs keep everything current (6h default)
- **🎨 Modern UI**: Dark-themed interface with glassmorphism effects and smooth animations
- **🚀 Fast & Reliable**: Built with Go for the backend and Next.js for the frontend
- **📱 Responsive**: Works perfectly on desktop, tablet, and mobile devices
- **🔐 Secure**: Rate limiting, authentication, and security headers built-in
- **📊 Monitoring**: Prometheus metrics and structured logging
- **🐳 Easy Deployment**: Complete Docker Compose stack with reverse proxy

## 🏗️ Architecture

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│                 │    │                 │    │                 │
│   Next.js UI    │    │   Go API        │    │   Go Worker     │
│                 │    │                 │    │                 │
│ • React 18      │    │ • REST API      │    │ • Cron Jobs     │
│ • Tailwind CSS  │    │ • Authentication│    │ • Source Fetch  │
│ • Framer Motion │    │ • Rate Limiting │    │ • Data Updates  │
│ • PWA Ready     │    │ • Metrics       │    │ • Queue System  │
│                 │    │                 │    │                 │
└─────────────────┘    └─────────────────┘    └─────────────────┘
         │                       │                       │
         └───────────────────────┼───────────────────────┘
                                 │
    ┌────────────────────────────┼────────────────────────────┐
    │                            │                            │
    ▼                            ▼                            ▼
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│                 │    │                 │    │                 │
│   Caddy Proxy   │    │   PostgreSQL    │    │     Redis       │
│                 │    │                 │    │                 │
│ • HTTPS         │    │ • Product Data  │    │ • Job Queue     │
│ • Compression   │    │ • Versions      │    │ • Rate Limits   │
│ • Load Balance  │    │ • Jobs History  │    │ • Caching       │
│                 │    │                 │    │                 │
└─────────────────┘    └─────────────────┘    └─────────────────┘
```

## 🚀 Quick Start

### 🐳 One-Click Docker Deployment

**Deploy anywhere with Docker in one command:**

```bash
bash <(curl -sSL https://raw.githubusercontent.com/nbt4/alldownloads/main/deploy.sh)
```

### Manual Docker Deployment

```bash
# Create project directory
mkdir alldownloads && cd alldownloads

# Download production configuration
curl -sSL -o docker-compose.yml https://raw.githubusercontent.com/nbt4/alldownloads/main/docker-compose.prod.yml
curl -sSL -o .env https://raw.githubusercontent.com/nbt4/alldownloads/main/.env.example

# Edit configuration (set AUTH_TOKEN, DOMAIN, etc.)
nano .env

# Deploy
docker compose up -d
```

### Development Setup

1. **Clone the repository**
   ```bash
   git clone https://github.com/nbt4/alldownloads.git
   cd alldownloads
   ```

2. **Configure environment**
   ```bash
   cp .env.example .env
   # Edit .env with your settings
   ```

3. **Start the stack**
   ```bash
   docker compose up -d
   ```

4. **Access the application**
   - **Web Interface**: http://localhost:3000
   - **API Endpoints**: http://localhost:8080/api
   - **Metrics**: http://localhost:8080/metrics

### Development Setup

For local development without Docker:

```bash
# Start only database and cache
make dev

# Run database migrations
migrate -path ./migrations -database "postgres://alldl:alldl@localhost:5432/alldownloads?sslmode=disable" up

# Start the API server
go run cmd/api/main.go

# Start the worker (in another terminal)
go run cmd/worker/main.go

# Start the UI (in another terminal)
cd ui && npm install && npm run dev
```

## 📋 Supported Software (15+ Products)

### 🖥️ Operating Systems
- **Ubuntu** - Latest LTS and current releases (72 versions)
- **Debian** - Stable releases with checksums (2 versions)
- **Arch Linux** - Rolling release ISOs (1 version)
- **Kali Linux** - Security-focused distribution (1 version)
- **Windows 11** - Microsoft Media Creation Tool links (1 version)

### 🌐 Web Browsers
- **Chrome** - Google Chrome stable releases (4 versions)
- **Firefox** - Mozilla Firefox latest (4 versions)
- **Brave** - Privacy-focused browser (4 versions)

### 💻 Development Tools
- **Visual Studio Code** - Microsoft's popular editor (4 versions)
- **PowerShell** - Cross-platform automation (4 versions)
- **Notepad++** - Windows text editor (2 versions)

### 💬 Communication
- **Telegram Desktop** - Secure messaging (4 versions)
- **WhatsApp Desktop** - Meta's messaging app (2 versions)

### 🛠️ Utilities
- **Termius** - SSH client (3 versions)
- **Tailscale** - VPN service (3 versions)
- **Nextcloud Desktop** - File sync client (3 versions)

**All downloads fetched from official vendor sources with real version numbers and file sizes!**

## 🔧 Configuration

### Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `DOMAIN` | `localhost` | Domain name for the application |
| `PORT` | `8080` | API server port |
| `AUTH_TOKEN` | `change-me` | Bearer token for API authentication |
| `DB_URL` | `postgres://...` | PostgreSQL connection string |
| `REDIS_URL` | `redis://...` | Redis connection string |
| `REFRESH_CRON` | `@every 6h` | Schedule for automatic updates |
| `HTTP_TIMEOUT` | `15s` | HTTP client timeout |
| `MAX_CONCURRENT_FETCHES` | `6` | Max concurrent source fetches |
| `RATE_LIMIT_REQUESTS_PER_MINUTE` | `60` | API rate limiting |

### Custom Sources

Add new software sources by implementing the `Fetcher` interface:

```go
type Fetcher interface {
    Fetch(ctx context.Context) ([]*store.ProductVersion, error)
}
```

See `internal/sources/` for examples.

## 🔌 API Reference

### Get all products
```http
GET /api/products
```

### Get product details
```http
GET /api/products/{id}
```

### Trigger refresh (requires auth)
```http
POST /api/refresh
Authorization: Bearer {token}
```

### Health check
```http
GET /api/health
```

### Metrics (Prometheus format)
```http
GET /metrics
```

## 🛠️ Development

### Make Commands

```bash
make help        # Show available commands
make build       # Build Docker images
make up          # Start all services
make down        # Stop all services
make logs        # Show service logs
make clean       # Clean up containers and volumes
make lint        # Run code linters
make test        # Run test suites
make migrate     # Run database migrations
make dev         # Start development environment
make prod        # Start production environment with MinIO
```

### Running Tests

```bash
# Go tests
go test -v ./...

# UI tests
cd ui && npm test

# Integration tests with Docker
make test
```

### Code Quality

- **Go**: Uses `golangci-lint` for comprehensive linting
- **TypeScript**: ESLint with Next.js configuration
- **Security**: Gosec for Go security scanning
- **Dependencies**: Automated security updates via Dependabot

## 🚢 Deployment

### Production Deployment

1. **Set up environment**
   ```bash
   cp .env.example .env
   # Configure production values
   ```

2. **Deploy with storage**
   ```bash
   make prod
   ```

3. **Set up domain and SSL**
   - Update `Caddyfile` with your domain
   - Caddy handles automatic HTTPS with Let's Encrypt

### 🐳 Docker Hub Images

Pre-built images are available on Docker Hub:

```bash
docker pull nbt4/alldownloads-api:latest     # REST API service (36.6MB)
docker pull nbt4/alldownloads-worker:latest  # Background worker (24.2MB)
```

**Production-ready features:**
- ✅ Multi-architecture support (amd64)
- ✅ Security scanning included
- ✅ Alpine-based minimal images
- ✅ Health checks built-in
- ✅ Non-root user execution

## 📊 Monitoring

### Metrics

The application exposes Prometheus metrics at `/metrics`:

- HTTP request metrics (duration, status codes)
- Fetch job statistics
- Product and version counts
- Database connection health

### Logging

Structured JSON logging with configurable levels:

```json
{
  "level": "info",
  "timestamp": "2024-01-15T10:30:00Z",
  "request_id": "abc123",
  "message": "Product updated",
  "product_id": "ubuntu",
  "versions": 4
}
```

### Health Checks

- **API**: `GET /api/health`
- **Database**: Connection pooling with health checks
- **Redis**: Ping-based health monitoring
- **Docker**: Built-in healthcheck directives

## 🔒 Security

### Security Features

- **Rate Limiting**: Token bucket algorithm per IP
- **Authentication**: Bearer token for admin endpoints
- **CORS**: Configurable cross-origin policies
- **Headers**: Security headers via Caddy
- **Input Validation**: Request validation and sanitization
- **Dependencies**: Regular security updates

### Security Policy

- Only official vendor sources are supported
- No direct file downloads or proxying by default
- All external requests use proper User-Agent strings
- Regular security scanning in CI/CD pipeline

## 🤝 Contributing

We welcome contributions! Please see our [Contributing Guide](CONTRIBUTING.md) for details.

### Development Workflow

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests for new functionality
5. Run the test suite
6. Submit a pull request

### Code Style

- **Go**: Follow standard Go conventions and `gofmt`
- **TypeScript**: Use Prettier and ESLint configurations
- **Commits**: Use conventional commit messages

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- [Ubuntu](https://ubuntu.com/) for their reliable release infrastructure
- [Debian](https://www.debian.org/) for comprehensive package management
- [Arch Linux](https://archlinux.org/) for rolling release innovation
- [Caddy](https://caddyserver.com/) for excellent reverse proxy capabilities
- [Next.js](https://nextjs.org/) for the amazing React framework
- [Tailwind CSS](https://tailwindcss.com/) for utility-first styling

## 📞 Support

- **Issues**: [GitHub Issues](https://github.com/your-username/alldownloads/issues)
- **Discussions**: [GitHub Discussions](https://github.com/your-username/alldownloads/discussions)
- **Documentation**: [Project Wiki](https://github.com/your-username/alldownloads/wiki)

---

<div align="center">
  <strong>Made with ❤️ for the open source community</strong>
</div>
