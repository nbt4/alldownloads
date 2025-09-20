# AllDownloads

[![CI/CD Pipeline](https://github.com/your-username/alldownloads/actions/workflows/ci.yml/badge.svg)](https://github.com/your-username/alldownloads/actions/workflows/ci.yml)
[![Docker Pulls](https://img.shields.io/docker/pulls/your-username/alldownloads)](https://hub.docker.com/r/your-username/alldownloads)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/your-username/alldownloads)](https://goreportcard.com/report/github.com/your-username/alldownloads)

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

### Prerequisites

- Docker and Docker Compose
- Git

### Installation

1. **Clone the repository**
   ```bash
   git clone https://github.com/your-username/alldownloads.git
   cd alldownloads
   ```

2. **Configure environment**
   ```bash
   cp .env.example .env
   # Edit .env with your settings
   ```

3. **Start the stack**
   ```bash
   make up
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

## 📋 Supported Software

### Operating Systems
- **Ubuntu** - Latest LTS and current releases
- **Debian** - Stable releases with checksums
- **Arch Linux** - Rolling release ISOs
- **Kali Linux** - Security-focused distribution
- **Windows 11** - Microsoft Media Creation Tool links

### Applications
- **Web Browsers**: Chrome, Firefox, Brave
- **Development**: Visual Studio Code, PowerShell
- **Communication**: Telegram Desktop, WhatsApp Desktop
- **Tools**: Termius, Tailscale, Nextcloud Desktop, Notepad++

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

### Container Registry

Pre-built images are available on GitHub Container Registry:

```bash
docker pull ghcr.io/your-username/alldownloads-api:latest
docker pull ghcr.io/your-username/alldownloads-worker:latest
docker pull ghcr.io/your-username/alldownloads-ui:latest
```

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
