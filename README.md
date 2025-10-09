# 🎬 CineVerse

> A modern social network for movie enthusiasts built with Flutter and Go

CineVerse is a comprehensive social platform designed for movie lovers to discover, rate, review, and discuss films. Built with modern technologies and following clean architecture principles, it offers a seamless experience across web and mobile platforms.

## 🏗️ Architecture Overview

This project follows a **monorepo** structure with clean architecture principles:

```
📂 CineVerse/
├── 📂 api_v2/          # ✅ Go Backend (Primary API)
│   ├── cmd/            # Main application entry point
│   ├── internal/       # Private application code
│   │   ├── domain/     # Business entities and rules
│   │   ├── handler/    # HTTP handlers (Chi router)
│   │   ├── service/    # Business logic layer
│   │   └── repository/ # Data access layer
│   └── migrations/     # Database migrations
├── 📂 flutter_app/     # ✅ Flutter Frontend (Multi-platform)
│   ├── lib/src/        # Application source code
│   │   ├── core/       # Shared utilities and configs
│   │   └── features/   # Feature-based modules
│   └── Dockerfile*     # Web & Android containers
├── 📂 api/             # ❌ Legacy NestJS API (Deprecated)
└── 📂 scripts/         # Setup and utility scripts
```

## 🚀 Tech Stack

### Backend (Go - `api_v2`)
- **Framework**: Chi (lightweight HTTP router)
- **Database**: PostgreSQL + SQLx
- **Cache**: Redis
- **Config**: Viper
- **Validation**: go-playground/validator
- **Logging**: slog (structured logging)

### Frontend (Flutter - `flutter_app`)
- **State Management**: Riverpod
- **Navigation**: go_router
- **HTTP Client**: Dio
- **Dependency Injection**: get_it
- **Platforms**: Web, Android, iOS

### Infrastructure
- **Containerization**: Docker + Docker Compose
- **Database**: PostgreSQL 15
- **Cache**: Redis 7
- **Development**: Hot-reloading enabled

## 📋 Prerequisites

Before you begin, ensure you have installed:

- **Docker** (v20.10+)
- **Docker Compose** (v2.0+)
- **Git**
- **Node.js** (v16+) - for development tools

### Optional (for local development)
- **Go** (v1.21+)
- **Flutter SDK** (v3.24+)

## 🎯 Quick Start

### 1. Clone the Repository
```bash
git clone https://github.com/EduardoMG12/cine.git
cd cine
```

### 2. Setup Development Environment
```bash
# Run the automated setup script
./scripts/setup.sh
```

This script will:
- Install Node.js dependencies (husky, lint-staged)
- Setup pre-commit hooks
- Build Docker images
- Initialize PostgreSQL with migrations
- Start all services
- Optionally setup Android development environment

### 3. Manual Setup (Alternative)
If you prefer manual setup:

```bash
# Install development dependencies
npm install

# Start infrastructure services
docker-compose up -d postgres redis

# Start the Go API
docker-compose up -d api_v2

# Start the Flutter web app
docker-compose up -d flutter_app

# Optional: Start Android development environment
docker-compose up -d flutter_android
```

## 🌐 Access Your Application

After setup, you can access:

| Service | URL | Description |
|---------|-----|-------------|
| **Flutter Web App** | http://localhost:3000 | Main application interface |
| **Go API** | http://localhost:8080 | REST API endpoints |
| **API Health Check** | http://localhost:8080/health | API status |
| **Android Studio** | http://localhost:6080 | Web-based Android development |
| **PostgreSQL** | localhost:5432 | Database (user: `cineverse`, db: `cineverse`) |
| **Redis** | localhost:6379 | Cache and sessions |

## 📱 Android Development

CineVerse includes a complete Android development environment with Android Studio running in Docker.

### Starting Android Environment
```bash
# Start the Android development container
docker-compose up -d flutter_android

# Access Android Studio in your browser
# URL: http://localhost:6080
# Password: cineverse
```

### Features:
- ✅ Android Studio IDE in browser
- ✅ Android SDK and build tools
- ✅ Pre-configured Android emulator
- ✅ Flutter Android development ready
- ✅ Hot reload support

### Building for Android
```bash
# Connect to the Android container
docker-compose exec flutter_android bash

# Build APK
flutter build apk

# Build App Bundle
flutter build appbundle

# Run on emulator
flutter run -d android
```

## 🔧 Development Workflow

### Code Quality
The project uses automated code quality tools:

```bash
# Format code (runs automatically on commit)
npm run lint-staged

# Manual formatting
# Go code
gofmt -w api_v2/
go vet ./api_v2/...

# Flutter code
flutter format flutter_app/
flutter analyze flutter_app/
```

### Database Migrations
```bash
# Run migrations
docker-compose exec api_v2 go run cmd/migrate.go up

# Create new migration
docker-compose exec api_v2 go run cmd/migrate.go create migration_name
```

### Logs and Debugging
```bash
# View logs for all services
docker-compose logs -f

# View specific service logs
docker-compose logs -f api_v2
docker-compose logs -f flutter_app
docker-compose logs -f flutter_android

# View database logs
docker-compose logs -f postgres
```

## 🧪 Testing

### Backend Tests (Go)
```bash
# Run unit tests
docker-compose exec api_v2 go test ./...

# Run tests with coverage
docker-compose exec api_v2 go test -cover ./...

# Run integration tests
docker-compose exec api_v2 go test -tags=integration ./...
```

### Frontend Tests (Flutter)
```bash
# Run Flutter tests
docker-compose exec flutter_app flutter test

# Run with coverage
docker-compose exec flutter_app flutter test --coverage
```

## 🚢 Production Deployment

### Building Production Images
```bash
# Build optimized production images
docker-compose -f docker-compose.prod.yml build

# Deploy to production
docker-compose -f docker-compose.prod.yml up -d
```

### Environment Variables
Create a `.env` file for production:

```env
# Database
CINE_DATABASE_URL=postgres://user:password@host:5432/database
CINE_REDIS_ADDR=redis:6379
CINE_REDIS_PASSWORD=your_redis_password

# API
CINE_PORT=8080
CINE_ENV=production

# Security
JWT_SECRET=your_jwt_secret_here
API_KEY=your_api_key_here
```

## 🤝 Contributing

We welcome contributions! Please see our [Contributing Guidelines](.github/CONTRIBUTING.md).

### Development Rules
1. **Follow Conventional Commits**: `feat:`, `fix:`, `docs:`, etc.
2. **Write Tests**: Ensure new features have adequate test coverage
3. **Code Quality**: All code must pass linting and formatting checks
4. **Architecture**: Follow clean architecture principles
5. **Documentation**: Update documentation for new features

### Commit Message Format
```
<type>(<scope>): <description>

[optional body]

[optional footer]
```

Examples:
- `feat(auth): add user registration endpoint`
- `fix(ui): resolve login form validation issue`
- `docs(api): update authentication documentation`

## 📚 API Documentation

### Authentication Endpoints
```
POST /api/auth/register    # User registration
POST /api/auth/login       # User login
POST /api/auth/refresh     # Refresh access token
DELETE /api/auth/logout    # User logout
```

### User Endpoints
```
GET    /api/users/profile  # Get user profile
PUT    /api/users/profile  # Update user profile
GET    /api/users/:id      # Get user by ID
```

### Movie Endpoints
```
GET    /api/movies         # List movies
GET    /api/movies/:id     # Get movie details
POST   /api/movies/:id/rating  # Rate a movie
GET    /api/movies/:id/reviews # Get movie reviews
```

## 🛠️ Utility Commands

### Reset Development Environment
```bash
# Complete reset (removes all data)
docker-compose down -v
docker system prune -f
./scripts/setup.sh
```

### Backup Database
```bash
# Create database backup
docker-compose exec postgres pg_dump -U cineverse cineverse > backup.sql

# Restore database backup
docker-compose exec -T postgres psql -U cineverse cineverse < backup.sql
```

### Performance Monitoring
```bash
# Monitor container resources
docker stats

# Monitor specific service
docker stats cineverse-api cineverse-flutter cineverse-postgres
```

## 🐛 Troubleshooting

### Common Issues

**Flutter app not starting:**
```bash
# Rebuild Flutter container
docker-compose build --no-cache flutter_app
docker-compose up flutter_app
```

**API connection issues:**
```bash
# Check API health
curl http://localhost:8080/health

# Restart API service
docker-compose restart api_v2
```

**Database connection problems:**
```bash
# Check PostgreSQL status
docker-compose exec postgres pg_isready -U cineverse

# Reset database
docker-compose down postgres
docker volume rm cineverse_postgres-data
docker-compose up -d postgres
```

**Android Studio not accessible:**
```bash
# Restart Android container
docker-compose restart flutter_android

# Check container logs
docker-compose logs flutter_android
```

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 👥 Team

- **Eduardo MG** - *Project Lead & Backend Developer* - [@EduardoMG12](https://github.com/EduardoMG12)

## 🙏 Acknowledgments

- Flutter team for the amazing framework
- Go community for excellent libraries
- Docker for containerization simplicity
- All contributors who make this project better

---

**Ready to build the future of movie social networking?** 🍿✨

For detailed technical documentation, check our [Wiki](../../wiki) or the `.github/RFCs` directory for feature specifications.