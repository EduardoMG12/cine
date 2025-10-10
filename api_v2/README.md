# CineVerse API v2

[![Go Version](https://img.shields.io/badge/Go-1.19+-blue.svg)](https://golang.org/)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)

## Overview

CineVerse is a social network platform for movie enthusiasts built with Go. It provides a comprehensive REST API for movie discovery, user reviews, social features, and collaborative movie selection through an innovative "movie matching" system.

## 🚀 Features

### ✅ **Implemented Features**

#### Authentication & User Management
- **User Registration**: Complete signup flow with email verification
- **Login System**: JWT-based authentication with secure session management  
- **Email Confirmation**: Token-based email verification system
- **Password Reset**: Secure password reset flow with time-limited tokens
- **User Profiles**: Profile management with privacy settings

#### Movie Management
- **Movie Data**: Integration with TMDb API for comprehensive movie information
- **Movie Lists**: "Want to Watch" and "Watched" lists management
- **Movie Search**: Search movies by title, genre, and other criteria
- **Caching Strategy**: Intelligent caching with TTL for optimal performance

#### Review System
- **Movie Reviews**: Full CRUD operations for user movie reviews
- **Rating System**: 1-10 star rating system for movies
- **Review Comments**: Text-based reviews with rich content support

#### Technical Infrastructure
- **Clean Architecture**: Domain-driven design with clear separation of concerns
- **Database**: PostgreSQL with proper migrations and indexing
- **Email Service**: SMTP email service with HTML templates
- **Security**: Password hashing, input validation, and JWT authentication
- **Validation**: Comprehensive input validation using go-playground/validator
- **Error Handling**: Structured error responses with proper HTTP status codes

### 🔄 **In Progress / Planned Features**

#### Social Features
- User following system
- Friendship management (send/accept/decline)
- User posts and activity feeds
- Review interactions (likes/dislikes)

#### Movie Matching System
- Collaborative movie selection sessions
- Multi-user preference matching
- Real-time voting system
- Match notifications

#### Advanced Features
- Real-time notifications
- Enhanced search and filtering
- Recommendation engine
- Advanced user preferences

## 📁 Project Structure

```
api_v2/
├── cmd/
│   └── main.go                 # Application entry point
├── internal/
│   ├── config/                 # Configuration management
│   ├── domain/                 # Business entities and interfaces
│   │   ├── email.go
│   │   ├── movie.go
│   │   ├── review.go
│   │   └── user.go
│   ├── dto/                    # Data Transfer Objects
│   │   ├── auth.go
│   │   └── movie.go
│   ├── handler/                # HTTP handlers
│   │   ├── auth_handler.go
│   │   ├── movie_handler.go
│   │   ├── movie_list_handler.go
│   │   ├── profile_handler.go
│   │   ├── review_handler.go
│   │   └── user_handler.go
│   ├── middleware/             # HTTP middleware
│   │   └── auth.go
│   ├── repository/             # Data access layer
│   │   ├── movie_repository.go
│   │   ├── review_repository.go
│   │   ├── user_repository.go
│   │   └── user_session_repository.go
│   ├── service/                # Business logic layer
│   │   ├── email_service.go
│   │   ├── movie_service.go
│   │   ├── review_service.go
│   │   ├── tmdb_service.go
│   │   ├── user_service.go
│   │   └── user_session_service.go
│   ├── server/                 # Server setup and routing
│   │   ├── routes.go
│   │   └── server.go
│   └── utils/                  # Utility functions
│       └── http.go
├── migrations/                 # Database migrations
│   ├── 001_initial_schema.sql
│   └── 002_complete_rfc_implementation.sql
└── build/                      # Build artifacts
```

## 🔧 Technology Stack

- **Language**: Go 1.19+
- **Framework**: Chi Router v5
- **Database**: PostgreSQL 15+ with sqlx
- **Cache**: Redis 7+
- **Authentication**: JWT tokens
- **Email**: SMTP with HTML templates
- **External APIs**: The Movie Database (TMDb)
- **Validation**: go-playground/validator/v10
- **Configuration**: Viper
- **Logging**: slog (structured logging)

## 🚀 Quick Start

### Prerequisites

- Go 1.19 or higher
- PostgreSQL 15+
- Redis 7+
- TMDb API key

### Installation

1. **Clone the repository**
```bash
git clone https://github.com/EduardoMG12/cine.git
cd cine/api_v2
```

2. **Install dependencies**
```bash
go mod tidy
```

3. **Set up environment variables**
```bash
# Copy example config
cp config.example.yaml config.yaml

# Edit configuration with your settings
nano config.yaml
```

4. **Run database migrations**
```bash
# Make sure PostgreSQL is running
psql -U postgres -f migrations/001_initial_schema.sql
psql -U postgres -f migrations/002_complete_rfc_implementation.sql
```

5. **Start the server**
```bash
go run cmd/main.go
```

The API will be available at `http://localhost:8080`

## ⚙️ Configuration

### Environment Variables

```yaml
server:
  port: "8080"
  host: "0.0.0.0"

database:
  url: "postgres://user:password@localhost:5432/cineverse?sslmode=disable"
  max_open_conns: 25
  max_idle_conns: 5
  conn_max_lifetime: 300

redis:
  addr: "localhost:6379"
  password: ""
  db: 0

jwt:
  secret: "your-secret-key-change-in-production"
  expiration: 24 # hours

email:
  smtp_host: "smtp.gmail.com"
  smtp_port: 587
  smtp_username: "your-email@gmail.com"
  smtp_password: "your-app-password"
  from_email: "noreply@cineverse.com"
  from_name: "CineVerse"

tmdb:
  api_key: "your-tmdb-api-key"
  base_url: "https://api.themoviedb.org/3"
```

### TMDb API Setup

1. Create an account at [The Movie Database](https://www.themoviedb.org/)
2. Request an API key from your account settings
3. Add the API key to your configuration

### Email Configuration

For Gmail SMTP:
1. Enable 2-factor authentication on your Google account
2. Generate an app-specific password
3. Use this app password in the configuration

## 📚 API Documentation

### Authentication Endpoints

#### Register User
```http
POST /api/v1/auth/register
Content-Type: application/json

{
  "username": "johndoe",
  "email": "john@example.com",
  "password": "SecurePass123",
  "display_name": "John Doe"
}
```

#### Login
```http
POST /api/v1/auth/login
Content-Type: application/json

{
  "email": "john@example.com",
  "password": "SecurePass123"
}
```

#### Confirm Email
```http
POST /api/v1/auth/confirm-email
Content-Type: application/json

{
  "token": "email-confirmation-token"
}
```

#### Forgot Password
```http
POST /api/v1/auth/forgot-password
Content-Type: application/json

{
  "email": "john@example.com"
}
```

#### Reset Password
```http
POST /api/v1/auth/reset-password
Content-Type: application/json

{
  "token": "reset-token",
  "new_password": "NewSecurePass123"
}
```

### Movie Endpoints

#### Search Movies
```http
GET /api/v1/movies/search?query=inception&page=1
Authorization: Bearer <jwt-token>
```

#### Get Movie Details
```http
GET /api/v1/movies/external/550
Authorization: Bearer <jwt-token>
```

#### Get Popular Movies
```http
GET /api/v1/movies/popular?page=1
Authorization: Bearer <jwt-token>
```

### Review Endpoints

#### Create Review
```http
POST /api/v1/reviews
Authorization: Bearer <jwt-token>
Content-Type: application/json

{
  "movie_id": 1,
  "rating": 9,
  "content": "Amazing movie with great plot and characters!"
}
```

#### Get Movie Reviews
```http
GET /api/v1/reviews/movie/1
Authorization: Bearer <jwt-token>
```

#### Update Review
```http
PUT /api/v1/reviews/123
Authorization: Bearer <jwt-token>
Content-Type: application/json

{
  "rating": 10,
  "content": "Actually, this is a masterpiece!"
}
```

### User Profile Endpoints

#### Get Current User Profile
```http
GET /api/v1/users/me
Authorization: Bearer <jwt-token>
```

#### Update Profile
```http
PUT /api/v1/users/me
Authorization: Bearer <jwt-token>
Content-Type: application/json

{
  "display_name": "John Updated",
  "bio": "Movie enthusiast and reviewer"
}
```

### Movie List Endpoints

#### Add to Want to Watch
```http
POST /api/v1/movie-lists/want-to-watch
Authorization: Bearer <jwt-token>
Content-Type: application/json

{
  "movie_id": 1
}
```

#### Move to Watched
```http
POST /api/v1/movie-lists/move-to-watched
Authorization: Bearer <jwt-token>
Content-Type: application/json

{
  "movie_external_id": "550"
}
```

## 🔒 Security Features

- **Password Hashing**: bcrypt with salt
- **JWT Authentication**: Secure token-based authentication
- **Input Validation**: Comprehensive validation on all inputs
- **SQL Injection Prevention**: Parameterized queries
- **Rate Limiting**: (Planned for production)
- **HTTPS Support**: TLS/SSL ready
- **CORS Configuration**: Configurable cross-origin policies

## 🧪 Testing

### Running Tests

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run specific package tests
go test ./internal/service/...
```

### Test Structure

Tests are organized following Go conventions:
- Unit tests for services and repositories
- Integration tests for handlers
- Table-driven tests for complex scenarios

## 📊 Database Schema

### Key Tables

- **users**: User accounts and profiles
- **user_sessions**: Active user sessions
- **movies**: Cached movie data from TMDb
- **reviews**: User movie reviews and ratings
- **movie_lists**: User movie lists (Want to Watch, Watched)
- **movie_list_entries**: Entries in movie lists
- **email_verification_tokens**: Email confirmation tokens
- **password_reset_tokens**: Password reset tokens

## 🚦 Development Workflow

### Code Standards

- **Clean Architecture**: Domain → Service → Repository → Handler
- **Error Handling**: Always handle errors explicitly
- **Validation**: Validate all inputs
- **Logging**: Use structured logging (slog)
- **Comments**: Minimal, focus on why, not what
- **Naming**: Descriptive names that express intent

### Git Workflow

- **Conventional Commits**: Use conventional commit format
- **Feature Branches**: Create branches for new features
- **Code Review**: All changes require review
- **Testing**: Ensure tests pass before merging

## 🔧 Deployment

### Docker Support

```dockerfile
# Dockerfile example (to be added)
FROM golang:1.19-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o main cmd/main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/main .
CMD ["./main"]
```

### Production Considerations

- Set strong JWT secrets
- Configure proper database connections
- Enable HTTPS/TLS
- Set up monitoring and logging
- Configure Redis for caching
- Set up email service properly

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'feat: add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🆘 Support

- **Documentation**: Check this README and code comments
- **Issues**: Open an issue on GitHub
- **Email**: Contact the development team

## 🗺️ Roadmap

### Phase 1 (Current) ✅
- [x] Authentication system
- [x] User management
- [x] Movie data integration
- [x] Review system
- [x] Email confirmation
- [x] Password reset

### Phase 2 (Next)
- [ ] Social features (following, friends)
- [ ] Movie matching system
- [ ] Real-time notifications
- [ ] Advanced search

### Phase 3 (Future)
- [ ] Recommendation engine
- [ ] Mobile API optimizations
- [ ] Performance enhancements
- [ ] Advanced analytics

---

**CineVerse API v2** - Building the future of movie social networking 🎬
