# Contributing to AllDownloads

Thank you for your interest in contributing to AllDownloads! This document provides guidelines and information for contributors.

## Table of Contents

- [Code of Conduct](#code-of-conduct)
- [Getting Started](#getting-started)
- [Development Setup](#development-setup)
- [Contributing Guidelines](#contributing-guidelines)
- [Pull Request Process](#pull-request-process)
- [Coding Standards](#coding-standards)
- [Testing](#testing)
- [Security](#security)

## Code of Conduct

This project adheres to a Code of Conduct that we expect all contributors to follow. Please read [CODE_OF_CONDUCT.md](CODE_OF_CONDUCT.md) for details.

## Getting Started

### Prerequisites

- Go 1.22 or later
- Node.js 20 or later
- Docker and Docker Compose
- Git

### Areas for Contribution

We welcome contributions in several areas:

1. **New Source Fetchers**: Add support for additional software products
2. **UI/UX Improvements**: Enhance the user interface and experience
3. **Performance Optimizations**: Improve efficiency and speed
4. **Documentation**: Help improve our documentation
5. **Bug Fixes**: Fix reported issues
6. **Testing**: Add or improve test coverage

## Development Setup

1. **Fork the repository** on GitHub

2. **Clone your fork**:
   ```bash
   git clone https://github.com/your-username/alldownloads.git
   cd alldownloads
   ```

3. **Set up the development environment**:
   ```bash
   make dev
   ```

4. **Install dependencies**:
   ```bash
   make deps
   ```

5. **Run the application locally**:
   ```bash
   # Terminal 1: API server
   go run cmd/api/main.go

   # Terminal 2: Worker
   go run cmd/worker/main.go

   # Terminal 3: UI
   cd ui && npm run dev
   ```

## Contributing Guidelines

### Reporting Issues

Before creating an issue, please:

1. **Search existing issues** to avoid duplicates
2. **Use the issue templates** provided
3. **Include relevant information**: OS, Go version, steps to reproduce
4. **Provide logs** when applicable

### Suggesting Features

For feature requests:

1. **Check existing issues** and discussions
2. **Describe the use case** clearly
3. **Explain the benefits** to users
4. **Consider backwards compatibility**

### Adding New Source Fetchers

When adding support for a new software product:

1. **Verify it's from official sources only**
2. **Implement the `Fetcher` interface**:
   ```go
   type Fetcher interface {
       Fetch(ctx context.Context) ([]*store.ProductVersion, error)
   }
   ```
3. **Add comprehensive tests**
4. **Update the seed data** in migrations
5. **Document the source** in README.md

Example structure:
```go
// internal/sources/mysoftware.go
package sources

import (
    "context"
    "github.com/your-username/alldownloads/internal/store"
)

type MySoftwareFetcher struct {
    client *HTTPClient
}

func NewMySoftwareFetcher() *MySoftwareFetcher {
    return &MySoftwareFetcher{
        client: NewHTTPClient(),
    }
}

func (f *MySoftwareFetcher) Fetch(ctx context.Context) ([]*store.ProductVersion, error) {
    // Implementation here
}
```

## Pull Request Process

1. **Create a feature branch** from `main`:
   ```bash
   git checkout -b feature/your-feature-name
   ```

2. **Make your changes** following our coding standards

3. **Add or update tests** as needed

4. **Run the test suite**:
   ```bash
   make test
   make lint
   ```

5. **Commit your changes** using conventional commits:
   ```bash
   git commit -m "feat: add support for new software X"
   ```

6. **Push to your fork**:
   ```bash
   git push origin feature/your-feature-name
   ```

7. **Open a Pull Request** with:
   - Clear title and description
   - Reference to related issues
   - Screenshots for UI changes
   - Breaking changes noted

### Commit Message Convention

We use [Conventional Commits](https://www.conventionalcommits.org/):

- `feat:` new features
- `fix:` bug fixes
- `docs:` documentation changes
- `style:` formatting changes
- `refactor:` code refactoring
- `test:` adding tests
- `chore:` maintenance tasks

Examples:
```
feat: add Firefox fetcher with version detection
fix: handle timeout errors in Ubuntu fetcher
docs: update API documentation for new endpoints
test: add integration tests for worker queue
```

## Coding Standards

### Go Code

- Follow standard Go conventions
- Use `gofmt` for formatting
- Run `golangci-lint` before committing
- Write meaningful variable and function names
- Add comments for exported functions
- Handle errors appropriately

### TypeScript/React Code

- Use TypeScript for all new code
- Follow the existing ESLint configuration
- Use Prettier for formatting
- Write semantic HTML
- Use Tailwind CSS for styling
- Follow React best practices

### Documentation

- Update README.md for new features
- Add inline comments for complex logic
- Include examples in documentation
- Update API documentation

## Testing

### Writing Tests

- **Unit tests**: Test individual functions
- **Integration tests**: Test component interactions
- **E2E tests**: Test complete workflows

### Go Testing

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -v -race -coverprofile=coverage.out ./...

# View coverage
go tool cover -html=coverage.out
```

### UI Testing

```bash
cd ui

# Run tests
npm test

# Run tests in watch mode
npm run test:watch
```

### Test Requirements

- **New features** must include tests
- **Bug fixes** should include regression tests
- **Maintain coverage** above 80%
- **Tests must be deterministic** and not flaky

## Security

### Security Guidelines

- **Never commit secrets** or credentials
- **Validate all inputs** from external sources
- **Use official sources only** for software downloads
- **Follow secure coding practices**
- **Report security issues** privately

### Reporting Security Issues

For security vulnerabilities, please email security@alldownloads.dev instead of creating a public issue.

## Getting Help

- **GitHub Discussions**: For general questions
- **GitHub Issues**: For bugs and feature requests
- **Documentation**: Check the project wiki
- **Code Review**: Ask questions in your PR

## Recognition

Contributors will be recognized in:
- README.md contributors section
- Release notes for significant contributions
- GitHub contributor statistics

Thank you for contributing to AllDownloads! 🚀