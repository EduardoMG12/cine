# 🎬 CineVerse

> A modern social network for movie enthusiasts built with Flutter and Go

**CineVerse** is a comprehensive social platform designed for movie lovers to discover, rate, review, and discuss films. Built with modern technologies and following clean architecture principles, it offers a seamless experience across web and mobile platforms.

[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go)](https://golang.org/)
[![Flutter](https://img.shields.io/badge/Flutter-3.24+-02569B?style=flat&logo=flutter)](https://flutter.dev/)
[![PostgreSQL](https://img.shields.io/badge/PostgreSQL-15-316192?style=flat&logo=postgresql)](https://www.postgresql.org/)
[![License](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)

## 🎯 Quick Start

```bash
# Clone repository
git clone https://github.com/EduardoMG12/cine.git
cd cine

# Run setup script
./scripts/setup.sh

# Or manually with Docker
docker-compose up -d
```

**Access the app**: http://localhost:3000  
**API Health**: http://localhost:8080/health

## ✨ Features

### 🎥 Movie Discovery
- Search movies by title, genre, actor
- Get detailed information (cast, crew, ratings)
- OMDb API integration (TMDb coming soon)
- Smart caching for performance

### 👤 User Profiles
- Secure JWT authentication
- Profile customization
- Watch history tracking
- Multi-device session management

### ⭐ Social Features
- Rate and review movies
- Create custom watchlists
- Follow other users (coming soon)
- Activity feed (coming soon)

### 📱 Multi-Platform
- ✅ Web (PWA)
- ✅ Android
- 🔄 iOS (coming soon)
- 🔄 Desktop (coming soon)

## 🏗️ Architecture

This project follows a **monorepo** structure with clean architecture:

```
cine/
├── api_v2/              # 🎯 Go Backend (PRIMARY)
│   ├── cmd/             # Application entrypoint
│   ├── internal/        # Business logic
│   │   ├── domain/      # Entities & interfaces
│   │   ├── usecase/     # Business rules
│   │   ├── handler/     # HTTP handlers
│   │   ├── repository/  # Data access
│   │   └── infrastructure/ # External services
│   ├── migrations/      # Database schema
│   └── README.md        # Backend docs
├── flutter_app/         # 📱 Flutter Frontend
│   ├── lib/src/         # App source code
│   │   ├── core/        # Shared utilities
│   │   └── features/    # Feature modules
│   └── Dockerfile       # Web & Android build
├── api/                 # ⚠️  Legacy NestJS (deprecated)
└── scripts/             # Automation scripts
```

### Tech Stack

**Backend (api_v2)**
- Go 1.21+ with Chi router
- PostgreSQL 15 + SQLx
- Redis 7 for caching
- JWT authentication
- OMDb/TMDb integration

**Frontend (flutter_app)**
- Flutter 3.24+
- Riverpod (state management)
- Dio (HTTP client)
- go_router (navigation)

**Infrastructure**
- Docker & Docker Compose
- PostgreSQL, Redis
- Development hot-reload

## 📋 Prerequisites

**Required:**
- Docker 20.10+ & Docker Compose 2.0+
- Git

**Optional (for local development):**
- Go 1.21+
- Flutter SDK 3.24+
- Node.js 16+ (for dev tools)

## 🚀 Getting Started

### Automated Setup (Recommended)

```bash
# 1. Clone repository
git clone https://github.com/EduardoMG12/cine.git
cd cine

# 2. Run automated setup
./scripts/setup.sh
```

This script will:
- Install Node.js dependencies & Git hooks
- Build Docker images
- Initialize PostgreSQL with migrations
- Start all services
- Setup Android environment (optional)

### Manual Setup

```bash
# 1. Install development dependencies
npm install

# 2. Start infrastructure
docker-compose up -d postgres redis

# 3. Start backend
docker-compose up -d api_v2

# 4. Start frontend
docker-compose up -d flutter_app

# 5. (Optional) Android development
docker-compose up -d flutter_android
```

### Environment Configuration

Create `.env` file in project root:

```bash
# Database
DB_HOST=postgres
DB_PORT=5432
DB_NAME=cineverse
DB_USER=cineverse
DB_PASSWORD=cineverse123

# API
PORT=8080
JWT_SECRET=your-super-secret-key
ENVIRONMENT=development

# External APIs
OMDB_API_KEY=your_omdb_key
TMDB_API_KEY=your_tmdb_key

# Redis
REDIS_HOST=redis
REDIS_PORT=6379
```

## 🌐 Services & Ports

After running the setup, access:

| Service | URL | Description |
|---------|-----|-------------|
| **Web App** | http://localhost:3000 | Flutter PWA |
| **API** | http://localhost:8080 | Go REST API |
| **Health** | http://localhost:8080/health | API status |
| **Android Studio** | http://localhost:6080 | Browser IDE (password: `cineverse`) |
| **PostgreSQL** | localhost:5432 | Database |
| **Redis** | localhost:6379 | Cache |

## � Development

### Backend Development

```bash
# Navigate to backend
cd api_v2

# Start with hot reload (install air first)
go install github.com/cosmtrek/air@latest
air

# Or manually
go run cmd/main.go

# Run tests
go test ./...

# Format code
go fmt ./...
```

See [api_v2/README.md](./api_v2/README.md) for detailed backend documentation.

### Frontend Development

```bash
# Navigate to frontend
cd flutter_app

# Install dependencies
flutter pub get

# Run on web
flutter run -d chrome

# Run on Android
flutter run -d android

# Run tests
flutter test

# Format code
flutter format .
```

### Code Quality

The project uses automated code quality tools:

```bash
# Format all code (runs on commit)
npm run lint-staged

# Backend formatting
cd api_v2
go fmt ./...
go vet ./...

# Frontend formatting
cd flutter_app
flutter format .
flutter analyze
```

### Database Management

```bash
# Connect to database
docker exec -it cineverse-postgres psql -U cineverse -d cineverse

# View tables
\dt

# Backup database
docker exec cineverse-postgres pg_dump -U cineverse cineverse > backup.sql

# Restore database
docker exec -i cineverse-postgres psql -U cineverse cineverse < backup.sql
```

### Viewing Logs

```bash
# All services
docker-compose logs -f

# Specific service
docker-compose logs -f api_v2
docker-compose logs -f flutter_app
docker-compose logs -f postgres

# Follow logs
docker-compose logs --tail=100 -f api_v2
```

## 🧪 Testing

### Backend Tests

```bash
cd api_v2

# Run all tests
go test ./...

# With coverage
go test -cover ./...
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out

# Integration tests
go test -tags=integration ./...

# Specific package
go test ./internal/usecase/auth/...
```

### Frontend Tests

```bash
cd flutter_app

# Run tests
flutter test

# With coverage
flutter test --coverage

# Integration tests
flutter test integration_test/
```

### API Testing

```bash
# Use provided test script
cd api_v2
./test-omdb.sh

# Or manually with curl
curl http://localhost:8080/health
curl http://localhost:8080/api/v1/omdb/test
curl "http://localhost:8080/api/v1/omdb/search?q=Batman"
```

## 🚢 Production Deployment

### Building for Production

```bash
# Backend
cd api_v2
CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o build/main ./cmd/main.go

# Frontend (Web)
cd flutter_app
flutter build web --release

# Frontend (Android)
flutter build apk --release
flutter build appbundle --release
```

### Docker Production Build

```bash
# Build optimized images
docker-compose -f docker-compose.prod.yml build

# Deploy
docker-compose -f docker-compose.prod.yml up -d
```

### Production Checklist

- [ ] Set strong `JWT_SECRET`
- [ ] Configure `ENVIRONMENT=production`
- [ ] Enable PostgreSQL SSL (`DB_SSL_MODE=require`)
- [ ] Set Redis password
- [ ] Configure CORS for production domains
- [ ] Setup SSL/TLS certificates
- [ ] Configure monitoring & logging
- [ ] Setup database backups
- [ ] Enable rate limiting
- [ ] Configure reverse proxy (nginx/caddy)

## 🤝 Contributing

We welcome contributions! Please follow these guidelines:

### Development Rules

1. **Follow Conventional Commits**: `feat:`, `fix:`, `docs:`, etc.
2. **Write Tests**: Ensure new features have test coverage
3. **Code Quality**: Pass all linting and formatting checks
4. **Clean Architecture**: Follow established patterns
5. **Documentation**: Update docs for new features

### Commit Format

```
<type>(<scope>): <description>

[optional body]
[optional footer]
```

**Examples:**
- `feat(auth): add user registration endpoint`
- `fix(ui): resolve login form validation`
- `docs(api): update authentication guide`

### Pull Request Process

1. Fork the repository
2. Create feature branch (`git checkout -b feature/amazing-feature`)
3. Make your changes
4. Run tests and linting
5. Commit with conventional format
6. Push to your fork
7. Open a Pull Request

## 🐛 Troubleshooting

### Common Issues

**Port already in use:**
```bash
# Find process using port
lsof -i :8080
# Kill process
kill -9 <PID>
```

**Database connection errors:**
```bash
# Check PostgreSQL status
docker exec cineverse-postgres pg_isready -U cineverse

# Restart database
docker-compose restart postgres

# Reset database (WARNING: deletes data)
docker-compose down postgres
docker volume rm cine_postgres-data
docker-compose up -d postgres
```

**Flutter app not starting:**
```bash
# Rebuild container
docker-compose build --no-cache flutter_app
docker-compose up flutter_app

# Check logs
docker-compose logs flutter_app
```

**API connection issues:**
```bash
# Check health
curl http://localhost:8080/health

# Restart API
docker-compose restart api_v2

# Rebuild API
docker-compose build api_v2
docker-compose up -d api_v2
```

### Reset Development Environment

```bash
# Complete reset (removes all data)
docker-compose down -v
docker system prune -f
./scripts/setup.sh
```

## 📚 Documentation

- **[Backend README](./api_v2/README.md)** - Complete Go API documentation
- **[OMDb Integration](./api_v2/OMDB_INTEGRATION.md)** - Movie API integration guide
- **[Architecture Decisions](./docs/adr/)** - ADRs and design decisions
- **[Contributing Guide](./.github/CONTRIBUTING.md)** - How to contribute

## 📡 API Endpoints

Quick reference for main endpoints:

### Authentication
```
POST   /api/v1/auth/register    # User registration
POST   /api/v1/auth/login       # User login
GET    /api/v1/auth/me          # Get current user
POST   /api/v1/auth/logout      # Logout
```

### Movies (OMDb)
```
GET    /api/v1/omdb/test                # Test connection
GET    /api/v1/omdb/{imdbId}            # Get movie by IMDb ID
GET    /api/v1/omdb/title               # Get movie by title
GET    /api/v1/omdb/search              # Search movies
GET    /api/v1/omdb/search-by-type      # Search by type
```

### Health
```
GET    /health                           # API health check
```

See [api_v2/README.md](./api_v2/README.md) for complete API documentation.

## 🛠️ Utility Commands

### Performance Monitoring

```bash
# Monitor all containers
docker stats

# Monitor specific services
docker stats cineverse-api cineverse-postgres cineverse-redis

# Check resource usage
docker system df
```

### Cleanup

```bash
# Stop all services
docker-compose down

# Remove volumes (WARNING: deletes data)
docker-compose down -v

# Clean system
docker system prune -a --volumes

# Remove specific volume
docker volume rm cine_postgres-data
```

### Database Operations

```bash
# Export database
docker exec cineverse-postgres pg_dump -U cineverse cineverse > backup_$(date +%Y%m%d).sql

# Import database
docker exec -i cineverse-postgres psql -U cineverse cineverse < backup.sql

# Connect to PostgreSQL CLI
docker exec -it cineverse-postgres psql -U cineverse -d cineverse

# View all tables
docker exec cineverse-postgres psql -U cineverse -d cineverse -c "\dt"
```

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 👥 Team

- **Eduardo MG** - *Project Lead & Developer* - [@EduardoMG12](https://github.com/EduardoMG12)

## 🙏 Acknowledgments

- Flutter team for the amazing framework
- Go community for excellent libraries
- [OMDb API](http://www.omdbapi.com/) for movie data
- [TMDb](https://www.themoviedb.org/) for additional movie information
- Docker for containerization simplicity
- All contributors who make this project better

## 📞 Support

- **Issues**: [GitHub Issues](https://github.com/EduardoMG12/cine/issues)
- **Documentation**: Check the [api_v2/README.md](./api_v2/README.md)
- **Email**: eduardo@cineverse.app

---

<div align="center">

**Ready to build the future of movie social networking?** 🍿✨

Made with ❤️ by the CineVerse team

[![Star on GitHub](https://img.shields.io/github/stars/EduardoMG12/cine?style=social)](https://github.com/EduardoMG12/cine)

</div>