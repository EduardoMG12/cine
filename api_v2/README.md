# 🎬 CineVerse API v2

> Modern movie social network backend built with Go, following Clean Architecture principles

[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go)](https://golang.org/)
[![PostgreSQL](https://img.shields.io/badge/PostgreSQL-15-316192?style=flat&logo=postgresql)](https://www.postgresql.org/)
[![Redis](https://img.shields.io/badge/Redis-7-DC382D?style=flat&logo=redis)](https://redis.io/)

## 📋 Table of Contents

- [Features](#-features)
- [Architecture](#-architecture)
- [Getting Started](#-getting-started)
- [Configuration](#-configuration)
- [API Documentation](#-api-documentation)
- [Development](#-development)
- [Testing](#-testing)

## ✨ Features

### Authentication & Authorization
- ✅ JWT-based authentication
- ✅ Secure password hashing (bcrypt)
- ✅ Session management with Redis
- ✅ Multi-device logout support

### Movie Data Integration
- ✅ **OMDb API Integration** (Adapter Pattern)
- ✅ Search movies by title, IMDb ID
- ✅ Get detailed movie information
- ✅ Ratings from multiple sources
- 🔄 TMDb API (Coming Soon)

### Infrastructure
- ✅ Structured logging (slog)
- ✅ Graceful shutdown
- ✅ Health check endpoints
- ✅ CORS configuration
- ✅ Request logging middleware

## 🏗️ Architecture

This project follows **Clean Architecture** principles with clear separation of concerns:

```
api_v2/
├── cmd/                        # Application entrypoint
│   └── main.go                
├── internal/                  
│   ├── config/                # Configuration management
│   ├── domain/                # Business entities
│   ├── dto/                   # Data Transfer Objects
│   ├── repository/            # Data access layer
│   ├── usecase/               # Business logic
│   ├── handler/               # HTTP handlers
│   ├── infrastructure/        # External services
│   ├── middleware/            # HTTP middleware
│   └── server/                # Server setup
├── migrations/                # Database migrations
├── .env                       # Environment variables
└── README.md                 
```

### Design Patterns

#### Adapter Pattern (Movie Providers)
Easy switching between different movie data providers (OMDb, TMDb):

```go
type MovieProvider interface {
    GetMovieByExternalID(id string) (*MovieDetails, error)
    GetMovieByTitle(title, year string) (*MovieDetails, error)
    SearchMovies(query string, page int) (*SearchResults, error)
}
```

#### Repository Pattern (Data Access)
```go
type UserRepository interface {
    Create(ctx context.Context, user *domain.User) error
    FindByEmail(ctx context.Context, email string) (*domain.User, error)
}
```

## 🚀 Getting Started

### Prerequisites

- **Go** 1.21+
- **PostgreSQL** 15+
- **Redis** 7+

### Quick Start

```bash
# 1. Clone and navigate
git clone https://github.com/EduardoMG12/cine.git
cd cine/api_v2

# 2. Install dependencies
go mod download

# 3. Setup environment
cp .env.example .env
# Edit .env with your configuration

# 4. Start dependencies (from project root)
docker-compose up -d postgres redis

# 5. Start the server
./start-server.sh
```

The API will be available at `http://localhost:8080`

## ⚙️ Configuration

Environment variables in `.env`:

```bash
# Server
PORT=8080
ENVIRONMENT=development

# Database
DB_HOST=localhost
DB_PORT=5432
DB_NAME=cineverse
DB_USER=cineverse
DB_PASSWORD=cineverse123

# JWT
JWT_SECRET=your-secret-key
JWT_EXPIRATION=24h

# OMDb API
OMDB_API_KEY=your_omdb_key

# Redis
REDIS_HOST=localhost
REDIS_PORT=6379
```

## 📡 API Documentation

### Base URL
```
http://localhost:8080/api/v1
```

### Authentication

#### Register
```http
POST /api/v1/auth/register
Content-Type: application/json

{
  "email": "user@example.com",
  "password": "SecurePass123!",
  "username": "johndoe",
  "full_name": "John Doe"
}
```

#### Login
```http
POST /api/v1/auth/login

{
  "email": "user@example.com",
  "password": "SecurePass123!"
}
```

#### Get Current User
```http
GET /api/v1/auth/me
Authorization: Bearer {access_token}
```

### OMDb Movies

#### Test Connection
```http
GET /api/v1/omdb/test
```

#### Get Movie by IMDb ID
```bash
curl http://localhost:8080/api/v1/omdb/tt0133093
```

#### Search Movies
```bash
curl "http://localhost:8080/api/v1/omdb/search?q=Batman&page=1"
```

#### Get by Title
```bash
curl "http://localhost:8080/api/v1/omdb/title?title=Inception&year=2010"
```

### Health Check
```http
GET /health
```

See [API_DOCUMENTATION.md](./API_DOCUMENTATION.md) for complete endpoint reference.

## 🔧 Development

### Running Locally

```bash
# With hot reload (install air first)
go install github.com/cosmtrek/air@latest
air

# Without hot reload
go run cmd/main.go
```

### Code Quality

```bash
# Format
go fmt ./...

# Lint
golangci-lint run

# Vet
go vet ./...
```

### Testing Endpoints

```bash
# Use test script
./test-omdb.sh

# Or manually
curl http://localhost:8080/health
curl http://localhost:8080/api/v1/omdb/test
```

## 🧪 Testing

```bash
# Run all tests
go test ./...

# With coverage
go test -cover ./...
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out

# Specific package
go test ./internal/usecase/auth/...
```

## 🐳 Docker

```bash
# Build
docker build -t cineverse-api:latest .

# Run
docker run -p 8080:8080 --env-file .env cineverse-api:latest

# With compose
docker-compose up -d api_v2
```

## 📚 Additional Documentation

- [OMDb Integration Guide](./OMDB_INTEGRATION.md) - Detailed OMDb setup and usage
- [Architecture Decision Records](../docs/adr/) - Design decisions
- [API Documentation](./API_DOCUMENTATION.md) - Complete API reference

## 🤝 Contributing

1. Fork the repository
2. Create feature branch (`git checkout -b feature/amazing-feature`)
3. Commit changes (`git commit -m 'Add amazing feature'`)
4. Push to branch (`git push origin feature/amazing-feature`)
5. Open Pull Request

## 📝 License

MIT License - see [LICENSE](../LICENSE) file for details.

## 🙏 Acknowledgments

- [OMDb API](http://www.omdbapi.com/) for movie data
- [Chi](https://github.com/go-chi/chi) for routing
- Go community for amazing libraries

---

**Built with ❤️ using Go** | CineVerse © 2025
