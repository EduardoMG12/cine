# 🎬 CineVerse API v2

CineVerse is a comprehensive social network platform for movie enthusiasts built with Go. It provides a robust backend API for movie discovery, reviews, social interactions, and personalized recommendations.

## 📋 Table of Contents

- [Features](#features)
- [Architecture](#architecture)
- [Getting Started](#getting-started)
- [API Documentation](#api-documentation)
- [Configuration](#configuration)
- [Database Setup](#database-setup)
- [Development](#development)
- [Testing](#testing)
- [Deployment](#deployment)
- [Contributing](#contributing)

## ✨ Features

### 🔐 Authentication & User Management
- JWT-based authentication with session management
- Email confirmation and verification
- Secure password reset functionality
- User profiles with privacy settings
- Theme preferences (light/dark mode)

### 🎥 Movie System
- Integration with The Movie Database (TMDb) API
- Advanced movie search and filtering
- Comprehensive movie information (cast, crew, ratings, etc.)
- Intelligent caching system for optimal performance

### ⭐ Review System
- Rate movies (1-10 scale)
- Write detailed text reviews
- View reviews by movie or user
- Update and manage personal reviews

### 📝 Movie Lists
- "Want to Watch" and "Watched" default lists
- Create custom movie lists
- Move movies between lists
- Share lists with other users

### 📧 Email System
- SMTP integration for transactional emails
- Beautiful HTML email templates
- Email confirmation for new registrations
- Password reset via secure email links

## 🏗️ Architecture

The API follows Clean Architecture principles with clear separation of concerns:

```
api_v2/
├── cmd/                    # Application entry points
├── internal/
│   ├── auth/              # Authentication utilities
│   ├── config/            # Configuration management
│   ├── domain/            # Business entities and interfaces
│   ├── dto/               # Data Transfer Objects
│   ├── handler/           # HTTP request handlers
│   ├── middleware/        # HTTP middleware
│   ├── repository/        # Data access layer
│   ├── service/           # Business logic layer
│   ├── server/            # Server setup and routing
│   └── utils/             # Utility functions
├── migrations/            # Database migrations
└── docs/                  # API documentation
```

### Technology Stack

- **Language**: Go 1.21+
- **Web Framework**: Chi Router
- **Database**: PostgreSQL with SQLX
- **Cache**: Redis
- **Authentication**: JWT
- **Email**: SMTP with HTML templates
- **Documentation**: Swagger/OpenAPI 3.0
- **External APIs**: The Movie Database (TMDb)

## 🚀 Getting Started

### Prerequisites

- Go 1.21 or higher
- PostgreSQL 15+
- Redis 7+
- SMTP server credentials
- TMDb API key

### Installation

1. **Clone the repository**
```bash
git clone https://github.com/EduardoMG12/cine.git
cd cine/api_v2
```

2. **Install dependencies**
```bash
go mod download
```

3. **Set up environment variables**
```bash
cp .env.example .env
# Edit .env with your configuration
```

4. **Run database migrations**
```bash
# Install golang-migrate if not installed
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

# Run migrations
migrate -path migrations -database "postgres://username:password@localhost:5432/cineverse?sslmode=disable" up
```

5. **Start the server**
```bash
go run cmd/main.go
```

The API will be available at `http://localhost:8080`

## 📚 API Documentation

### Swagger UI
Access the interactive API documentation at: `http://localhost:8080/swagger/index.html`

### Generate Swagger Documentation
```bash
# Install swag CLI tool
go install github.com/swaggo/swag/cmd/swag@latest

# Generate docs
swag init -g cmd/main.go -o docs/
```

### Core Endpoints

#### Authentication
```
POST   /api/v1/auth/register        # User registration
POST   /api/v1/auth/login           # User login
POST   /api/v1/auth/confirm-email   # Confirm email address
POST   /api/v1/auth/forgot-password # Request password reset
POST   /api/v1/auth/reset-password  # Reset password with token
```

#### Users
```
GET    /api/v1/users/{id}           # Get user profile
PUT    /api/v1/users/me             # Update own profile
PUT    /api/v1/users/me/settings    # Update user settings
```

#### Movies
```
GET    /api/v1/movies/search        # Search movies
GET    /api/v1/movies/popular       # Get popular movies
GET    /api/v1/movies/genre/{genre} # Get movies by genre
GET    /api/v1/movies/{id}          # Get movie details
GET    /api/v1/movies/external/{id} # Get movie by external ID
```

#### Reviews
```
POST   /api/v1/reviews              # Create review
GET    /api/v1/reviews/{id}         # Get specific review
GET    /api/v1/reviews/movie/{id}   # Get movie reviews
GET    /api/v1/reviews/user/{id}    # Get user reviews
PUT    /api/v1/reviews/{id}         # Update review
DELETE /api/v1/reviews/{id}         # Delete review
```

#### Movie Lists
```
GET    /api/v1/movie-lists          # Get user's lists
POST   /api/v1/movie-lists          # Create new list
GET    /api/v1/movie-lists/{id}     # Get specific list
PUT    /api/v1/movie-lists/{id}     # Update list
DELETE /api/v1/movie-lists/{id}     # Delete list
POST   /api/v1/movie-lists/want-to-watch   # Add to want-to-watch
POST   /api/v1/movie-lists/watched         # Add to watched
POST   /api/v1/movie-lists/move-to-watched # Move to watched
```

## ⚙️ Configuration

### Environment Variables

```bash
# Server Configuration
SERVER_HOST=0.0.0.0
SERVER_PORT=8080

# Database Configuration
DATABASE_URL=postgres://username:password@localhost:5432/cineverse?sslmode=disable
DATABASE_MAX_OPEN_CONNS=25
DATABASE_MAX_IDLE_CONNS=5
DATABASE_CONN_MAX_LIFETIME=300

# Redis Configuration
REDIS_ADDR=localhost:6379
REDIS_PASSWORD=
REDIS_DB=0

# JWT Configuration
JWT_SECRET=your-super-secret-jwt-key-change-in-production
JWT_EXPIRATION=24

# Email Configuration
EMAIL_SMTP_HOST=smtp.gmail.com
EMAIL_SMTP_PORT=587
EMAIL_SMTP_USERNAME=your-email@gmail.com
EMAIL_SMTP_PASSWORD=your-app-password
EMAIL_FROM_EMAIL=noreply@cineverse.com
EMAIL_FROM_NAME=CineVerse

# TMDb Configuration
TMDB_API_KEY=your-tmdb-api-key
TMDB_BASE_URL=https://api.themoviedb.org/3
```

### Configuration File (config.yaml)

```yaml
server:
  host: "0.0.0.0"
  port: "8080"

database:
  url: "postgres://username:password@localhost:5432/cineverse?sslmode=disable"
  max_open_conns: 25
  max_idle_conns: 5
  conn_max_lifetime: 300

redis:
  addr: "localhost:6379"
  password: ""
  db: 0

jwt:
  secret: "your-jwt-secret"
  expiration: 24

email:
  smtp_host: "smtp.gmail.com"
  smtp_port: 587
  smtp_username: "your-email@gmail.com"
  smtp_password: "your-password"
  from_email: "noreply@cineverse.com"
  from_name: "CineVerse"

tmdb:
  api_key: "your-tmdb-api-key"
  base_url: "https://api.themoviedb.org/3"
```

## 🗄️ Database Setup

### PostgreSQL Schema

The application uses the following main entities:

- **users**: User accounts and profiles
- **movies**: Movie information from TMDb
- **reviews**: User reviews and ratings
- **movie_lists**: Custom movie lists
- **movie_list_entries**: Movies in lists
- **email_verification_tokens**: Email confirmation tokens
- **password_reset_tokens**: Password reset tokens
- **user_sessions**: Active user sessions

### Migration Commands

```bash
# Create new migration
migrate create -ext sql -dir migrations -seq migration_name

# Apply migrations
migrate -path migrations -database $DATABASE_URL up

# Rollback migrations
migrate -path migrations -database $DATABASE_URL down

# Check migration status
migrate -path migrations -database $DATABASE_URL version
```

## 🛠️ Development

### Project Structure

```
internal/
├── auth/                   # JWT and password handling
├── config/                 # Configuration loading and validation
├── domain/                 # Business entities and interfaces
│   ├── user.go            # User domain model and interfaces
│   ├── movie.go           # Movie domain model and interfaces
│   ├── review.go          # Review domain model and interfaces
│   └── email.go           # Email service interfaces
├── dto/                   # Request/Response DTOs
│   ├── auth.go            # Authentication DTOs
│   ├── movie.go           # Movie and list DTOs
│   └── review.go          # Review DTOs
├── handler/               # HTTP handlers
│   ├── auth_handler.go    # Authentication endpoints
│   ├── user_handler.go    # User management endpoints
│   ├── movie_handler.go   # Movie endpoints
│   ├── movie_list_handler.go # Movie list endpoints
│   └── review_handler.go  # Review endpoints
├── middleware/            # HTTP middleware
│   └── auth.go           # Authentication middleware
├── repository/           # Data access layer
│   ├── user_repository.go    # User data access
│   ├── movie_repository.go   # Movie data access
│   └── review_repository.go  # Review data access
├── service/             # Business logic layer
│   ├── auth_service.go     # Authentication logic
│   ├── user_service.go     # User business logic
│   ├── movie_service.go    # Movie business logic
│   ├── review_service.go   # Review business logic
│   ├── email_service.go    # Email sending logic
│   └── tmdb_service.go     # TMDb API integration
├── server/              # Server setup
│   ├── server.go         # Server initialization
│   └── routes.go         # Route configuration
└── utils/              # Utility functions
    └── http.go          # HTTP utilities
```

### Code Standards

- Follow Go best practices and idioms
- Use meaningful variable and function names
- Write comprehensive unit tests
- Document public APIs with godoc
- Use structured logging with `log/slog`
- Implement proper error handling
- Follow Clean Architecture principles

### Adding New Features

1. **Define the domain model** in `internal/domain/`
2. **Create DTOs** in `internal/dto/`
3. **Implement repository** in `internal/repository/`
4. **Add business logic** in `internal/service/`
5. **Create HTTP handlers** in `internal/handler/`
6. **Add routes** in `internal/server/routes.go`
7. **Write tests** for all layers
8. **Update documentation**

## 🧪 Testing

### Running Tests

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run tests with detailed coverage report
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out

# Run specific package tests
go test ./internal/service/...

# Run tests with verbose output
go test -v ./...
```

### Test Structure

```
├── internal/
│   ├── service/
│   │   ├── user_service.go
│   │   └── user_service_test.go
│   ├── repository/
│   │   ├── user_repository.go
│   │   └── user_repository_test.go
│   └── handler/
│       ├── auth_handler.go
│       └── auth_handler_test.go
```

## 🚢 Deployment

### Docker Deployment

```dockerfile
# Build stage
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o main cmd/main.go

# Runtime stage
FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/main .
COPY --from=builder /app/migrations ./migrations
CMD ["./main"]
```

### Docker Compose

```yaml
version: '3.8'
services:
  api:
    build: .
    ports:
      - "8080:8080"
    environment:
      - DATABASE_URL=postgres://postgres:password@db:5432/cineverse?sslmode=disable
      - REDIS_ADDR=redis:6379
    depends_on:
      - db
      - redis

  db:
    image: postgres:15
    environment:
      POSTGRES_DB: cineverse
      POSTGRES_USER: postgres
      POSTGRES_PASSWORD: password
    ports:
      - "5432:5432"
    volumes:
      - postgres_data:/var/lib/postgresql/data

  redis:
    image: redis:7-alpine
    ports:
      - "6379:6379"

volumes:
  postgres_data:
```

### Production Considerations

- Use strong JWT secrets
- Configure HTTPS/TLS
- Set up database backups
- Monitor application metrics
- Configure log aggregation
- Set up health checks
- Use environment-specific configurations
- Implement rate limiting
- Set up monitoring and alerting

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Make your changes following the code standards
4. Write tests for new functionality
5. Update documentation as needed
6. Commit your changes (`git commit -m 'feat: add amazing feature'`)
7. Push to the branch (`git push origin feature/amazing-feature`)
8. Open a Pull Request

### Commit Convention

We use [Conventional Commits](https://conventionalcommits.org/):

```
feat: add new feature
fix: bug fix
docs: documentation changes
style: formatting changes
refactor: code refactoring
test: add tests
chore: maintenance tasks
```

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🔗 Links

- [TMDb API Documentation](https://developers.themoviedb.org/3)
- [Go Documentation](https://golang.org/doc/)
- [Chi Router](https://go-chi.io/)
- [PostgreSQL Documentation](https://www.postgresql.org/docs/)
- [Redis Documentation](https://redis.io/documentation)

## 📞 Support

If you have any questions or need help, please:

1. Check the [documentation](#api-documentation)
2. Search existing [issues](https://github.com/EduardoMG12/cine/issues)
3. Create a new issue with detailed information
4. Contact the development team

---

**Made with ❤️ by the CineVerse Team**
