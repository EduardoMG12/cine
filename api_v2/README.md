<div align="center">

# 🎬 CineVerse API

### RESTful Backend for Movie Social Network

[![Go](https://img.shields.io/badge/Go-1.21+-00ADD8?style=for-the-badge&logo=go&logoColor=white)](https://go.dev/)
[![PostgreSQL](https://img.shields.io/badge/PostgreSQL-15-336791?style=for-the-badge&logo=postgresql&logoColor=white)](https://www.postgresql.org/)
[![Redis](https://img.shields.io/badge/Redis-7-DC382D?style=for-the-badge&logo=redis&logoColor=white)](https://redis.io/)
[![Swagger](https://img.shields.io/badge/Swagger-OpenAPI-85EA2D?style=for-the-badge&logo=swagger&logoColor=black)](https://swagger.io/)
[![Docker](https://img.shields.io/badge/Docker-Ready-2496ED?style=for-the-badge&logo=docker&logoColor=white)](https://www.docker.com/)

[🇺🇸 English](#english) | [🇧🇷 Português](#português)

</div>

---

<a name="english"></a>

<details open>
<summary><h2>🇺🇸 ENGLISH VERSION</h2></summary>

<details>
<summary><h3>📋 Table of Contents</h3></summary>

- [About](#about-en)
- [Features](#features-en)
- [Architecture](#architecture-en)
- [API Documentation](#api-docs-en)
- [Getting Started](#getting-started-en)
- [Configuration](#configuration-en)
- [Project Structure](#structure-en)
- [Development](#development-en)
- [Testing](#testing-en)
- [Deployment](#deployment-en)

</details>

<details>
<summary><h3>📖 About</h3></summary>

<a name="about-en"></a>

CineVerse API is a robust RESTful backend built with Go, implementing Clean Architecture principles and designed for scalability and maintainability.

**Key Features:**
- 🔐 **Stateless JWT Authentication** with bcrypt password hashing
- 🎬 **Multi-Provider Movie Data** using Chain of Responsibility pattern
- 📊 **Comprehensive API Documentation** with Swagger/OpenAPI
- 🗄️ **PostgreSQL Database** with migration support
- ⚡ **Redis Caching** for performance optimization
- 🌐 **Internationalization** (English/Portuguese)
- 🐳 **Docker Support** for easy deployment

**Technical Highlights:**
- Clean Architecture with clear layer separation
- Repository pattern for data access
- Use case driven design
- Dependency injection
- Structured logging with slog
- Error handling with custom error types
- Input validation and sanitization

</details>

<details>
<summary><h3>✨ Features</h3></summary>

<a name="features-en"></a>

<details>
<summary><strong>Authentication & Authorization</strong></summary>

#### User Registration
- Username validation (alphanumeric, 3-30 chars)
- Email format validation
- Strong password requirements
- Bcrypt password hashing (cost 12)
- Duplicate username/email detection

#### User Login
- Credential validation
- JWT token generation
- Configurable token expiration
- Stateless session management

#### Protected Routes
- JWT middleware for route protection
- Token validation and parsing
- User context injection
- Automatic error responses

</details>

<details>
<summary><strong>User Management</strong></summary>

#### Profile Operations
- Get user by ID
- Get user by username
- Update profile information
- Display name management
- User search functionality

#### Data Validation
- Username uniqueness
- Email format validation
- Required field validation
- Data sanitization

</details>

<details>
<summary><strong>Movie Integration</strong></summary>

#### Chain of Responsibility
```
Request → OMDb API → Database Cache → 404
          (80-90%)    (10-20%)         (<1%)
```

**Advantages:**
- Reduces API calls by 80-90%
- Automatic data persistence
- 48-hour cache TTL
- Provider tracking (OMDb, Database)
- Fallback mechanism for reliability

#### Movie Operations
- Search movies with pagination
- Get movie details by ID
- Auto-save fetched movies
- Provider and sync tracking
- Genre and metadata support

#### OMDb Integration
- Direct OMDb API access
- Search with filters
- Detailed movie information
- Poster and rating data
- Configurable API key

</details>

<details>
<summary><strong>Infrastructure</strong></summary>

#### Database
- PostgreSQL 15 with UUID primary keys
- Connection pooling
- Migration system
- JSONB support for flexible data
- Full-text search capabilities

#### Caching
- Redis integration
- Configurable TTL
- Optional cache support
- Connection error handling

#### API Documentation
- Swagger/OpenAPI 2.0
- Interactive API explorer
- Request/response examples
- Schema definitions
- Authentication documentation

#### Middleware
- CORS configuration
- Request logging
- Auth validation
- Error handling
- Recovery from panics

</details>

</details>

<details>
<summary><h3>🏗️ Architecture</h3></summary>

<a name="architecture-en"></a>

<details>
<summary><strong>Clean Architecture Layers</strong></summary>

```
┌─────────────────────────────────────────┐
│           HTTP Handlers                  │  ← Presentation Layer
│  (Request/Response, DTOs, Validation)    │
└──────────────┬──────────────────────────┘
               │
┌──────────────▼──────────────────────────┐
│           Use Cases                      │  ← Application Layer
│  (Business Logic, Orchestration)         │
└──────────────┬──────────────────────────┘
               │
┌──────────────▼──────────────────────────┐
│           Domain                         │  ← Domain Layer
│  (Entities, Interfaces, Business Rules)  │
└──────────────┬──────────────────────────┘
               │
┌──────────────▼──────────────────────────┐
│      Infrastructure                      │  ← Infrastructure Layer
│  (Database, External APIs, Services)     │
└──────────────────────────────────────────┘
```

</details>

<details>
<summary><strong>Directory Structure</strong></summary>

```
api_v2/
├── cmd/
│   └── main.go                 # Application entry point
├── internal/
│   ├── config/
│   │   └── config.go           # Configuration management
│   ├── domain/
│   │   ├── auth.go             # Auth entity & interfaces
│   │   ├── movie.go            # Movie entity & interfaces
│   │   └── user.go             # User entity & interfaces
│   ├── dto/
│   │   ├── auth_dto.go         # Auth data transfer objects
│   │   ├── movie_dto.go        # Movie data transfer objects
│   │   └── user_dto.go         # User data transfer objects
│   ├── handler/
│   │   └── http/
│   │       ├── auth_handler.go # Auth HTTP handlers
│   │       ├── movie_handler.go# Movie HTTP handlers
│   │       ├── omdb_handler.go # OMDb HTTP handlers
│   │       └── user_handler.go # User HTTP handlers
│   ├── infrastructure/
│   │   ├── jwt.go              # JWT token service
│   │   ├── movie_provider.go   # Movie provider interface
│   │   └── movie_fetcher_chain.go # Chain of Responsibility
│   ├── middleware/
│   │   ├── auth.go             # JWT auth middleware
│   │   ├── cors.go             # CORS middleware
│   │   └── logger.go           # Request logging
│   ├── repository/
│   │   ├── movie_repository.go # Movie data access
│   │   ├── session_repository.go# Session data access
│   │   └── user_repository.go  # User data access
│   ├── service/
│   │   └── omdb_service.go     # OMDb API client
│   ├── usecase/
│   │   ├── auth/               # Auth use cases
│   │   ├── movie/              # Movie use cases
│   │   └── user/               # User use cases
│   └── i18n/
│       ├── i18n.go             # Internationalization setup
│       └── locales/
│           ├── en.json         # English translations
│           └── pt.json         # Portuguese translations
├── migrations/
│   ├── 001_clean_initial_schema.sql  # Initial schema
│   └── 002_add_provider_and_sync.sql # Provider tracking
├── docs/
│   ├── docs.go                 # Generated Swagger docs
│   ├── swagger.json            # OpenAPI JSON spec
│   └── swagger.yaml            # OpenAPI YAML spec
└── README.md                   # This file
```

</details>

<details>
<summary><strong>Design Patterns</strong></summary>

#### Chain of Responsibility
Used for movie data fetching with automatic fallback:
- **OMDbMovieFetcher**: Primary data source
- **DatabaseMovieFetcher**: Fallback cache
- **Auto-save**: Persist OMDb data to database

#### Repository Pattern
Abstracts data access logic:
- Interface-based contracts
- Easy testing with mocks
- Database implementation independence

#### Dependency Injection
Services and repositories injected via constructors:
- Loose coupling between layers
- Testability
- Flexibility to swap implementations

#### DTO Pattern
Separates internal entities from API contracts:
- Input validation
- Response formatting
- Versioning support

</details>

<details>
<summary><strong>Database Schema</strong></summary>

#### Users Table
```sql
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    username VARCHAR(30) UNIQUE NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    display_name VARCHAR(100),
    password_hash TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

#### Movies Table
```sql
CREATE TABLE movies (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    external_api_id VARCHAR(50) UNIQUE,
    title VARCHAR(255) NOT NULL,
    overview TEXT,
    poster_url TEXT,
    genres JSONB,
    adult BOOLEAN DEFAULT false,
    provider VARCHAR(50),
    last_sync_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

#### Sessions Table
```sql
CREATE TABLE sessions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    token TEXT NOT NULL,
    expires_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

</details>

</details>

<details>
<summary><h3>📚 API Documentation</h3></summary>

<a name="api-docs-en"></a>

<details>
<summary><strong>Swagger/OpenAPI</strong></summary>

The API is fully documented using Swagger/OpenAPI 2.0.

**Access Swagger UI:**
```
http://localhost:8080/swagger/index.html
```

**Features:**
- Interactive API explorer
- Try-it-out functionality
- Request/response examples
- Schema definitions
- Authentication documentation

**Endpoints:**
- 14 documented endpoints
- Complete request/response schemas
- Error code documentation
- Example payloads

</details>

<details>
<summary><strong>Authentication Endpoints</strong></summary>

#### POST /api/v1/auth/register
Register a new user account.

**Request:**
```json
{
  "username": "johndoe",
  "email": "john@example.com",
  "display_name": "John Doe",
  "password": "SecurePass123!"
}
```

**Response (201):**
```json
{
  "success": true,
  "message": "User registered successfully",
  "data": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "username": "johndoe",
    "email": "john@example.com",
    "display_name": "John Doe"
  }
}
```

#### POST /api/v1/auth/login
Authenticate and receive JWT token.

**Request:**
```json
{
  "username": "johndoe",
  "password": "SecurePass123!"
}
```

**Response (200):**
```json
{
  "success": true,
  "message": "Login successful",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIs...",
    "user": {
      "id": "550e8400-e29b-41d4-a716-446655440000",
      "username": "johndoe",
      "display_name": "John Doe"
    }
  }
}
```

#### POST /api/v1/auth/logout
Invalidate current session (requires authentication).

</details>

<details>
<summary><strong>Movie Endpoints</strong></summary>

#### GET /api/v1/movies/search
Search movies using Chain of Responsibility (OMDb → Database).

**Parameters:**
- `q` (query string, required): Search term
- `page` (integer, optional): Page number (default: 1)

**Example:**
```bash
curl "http://localhost:8080/api/v1/movies/search?q=matrix&page=1"
```

**Response (200):**
```json
{
  "success": true,
  "message": "Movies found",
  "data": [
    {
      "id": "550e8400-e29b-41d4-a716-446655440000",
      "external_api_id": "tt0133093",
      "title": "The Matrix",
      "overview": "1999",
      "poster_url": "https://m.media-amazon.com/...",
      "genres": null,
      "adult": false,
      "provider": "OMDb",
      "last_sync_at": "2025-11-08T18:40:19Z"
    }
  ]
}
```

#### GET /api/v1/movies/{id}
Get movie details by ID (OMDb → Database chain).

**Example:**
```bash
curl "http://localhost:8080/api/v1/movies/tt0133093"
```

</details>

<details>
<summary><strong>OMDb Direct Endpoints</strong></summary>

#### GET /api/v1/omdb/{imdb_id}
Get movie details directly from OMDb API.

**Example:**
```bash
curl "http://localhost:8080/api/v1/omdb/tt0133093"
```

#### GET /api/v1/omdb/search
Search movies directly on OMDb API.

**Parameters:**
- `q` (query string, required): Search term
- `page` (integer, optional): Page number

**Example:**
```bash
curl "http://localhost:8080/api/v1/omdb/search?q=batman&page=1"
```

</details>

<details>
<summary><strong>User Endpoints</strong></summary>

#### GET /api/v1/users/me
Get current authenticated user profile (requires authentication).

**Headers:**
```
Authorization: Bearer <token>
```

**Response (200):**
```json
{
  "success": true,
  "data": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "username": "johndoe",
    "email": "john@example.com",
    "display_name": "John Doe",
    "created_at": "2025-11-08T15:30:00Z"
  }
}
```

#### GET /api/v1/users/{username}
Get user profile by username.

</details>

<details>
<summary><strong>Health Check</strong></summary>

#### GET /health
Check API health status.

**Response (200):**
```json
{
  "status": "healthy",
  "service": "cineverse-api"
}
```

</details>

<details>
<summary><strong>Error Responses</strong></summary>

All endpoints return consistent error responses:

**400 Bad Request:**
```json
{
  "success": false,
  "message": "Invalid input data",
  "error": "username must be between 3 and 30 characters"
}
```

**401 Unauthorized:**
```json
{
  "success": false,
  "message": "Authentication required",
  "error": "Missing or invalid token"
}
```

**404 Not Found:**
```json
{
  "success": false,
  "message": "Resource not found",
  "error": "Movie not found"
}
```

**500 Internal Server Error:**
```json
{
  "success": false,
  "message": "Internal server error",
  "error": "Database connection failed"
}
```

</details>

</details>

<details>
<summary><h3>🚀 Getting Started</h3></summary>

<a name="getting-started-en"></a>

<details>
<summary><strong>Prerequisites</strong></summary>

- **Go** 1.21 or higher
- **PostgreSQL** 15
- **Redis** 7 (optional but recommended)
- **OMDb API Key** (free at [omdbapi.com](https://www.omdbapi.com/apikey.aspx))
- **Docker & Docker Compose** (for containerized setup)

</details>

<details>
<summary><strong>Installation (Docker)</strong></summary>

1. **Clone repository**
```bash
git clone https://github.com/EduardoMG12/cine.git
cd cine/api_v2
```

2. **Configure environment**
```bash
cp .env.example .env
# Edit .env with your configuration
```

3. **Start services**
```bash
docker-compose up -d
```

4. **Run migrations**
```bash
docker-compose exec api psql -U cineverse -d cineverse -f /app/migrations/001_clean_initial_schema.sql
docker-compose exec api psql -U cineverse -d cineverse -f /app/migrations/002_add_provider_and_sync.sql
```

5. **Verify health**
```bash
curl http://localhost:8080/health
```

</details>

<details>
<summary><strong>Installation (Local)</strong></summary>

1. **Install dependencies**
```bash
go mod download
```

2. **Setup PostgreSQL**
```bash
psql -U postgres
CREATE DATABASE cineverse;
CREATE USER cineverse WITH PASSWORD 'your_password';
GRANT ALL PRIVILEGES ON DATABASE cineverse TO cineverse;
\q
```

3. **Run migrations**
```bash
psql -U cineverse -d cineverse -f migrations/001_clean_initial_schema.sql
psql -U cineverse -d cineverse -f migrations/002_add_provider_and_sync.sql
```

4. **Configure environment**
```bash
export DB_HOST=localhost
export DB_PORT=5432
export DB_USER=cineverse
export DB_PASSWORD=your_password
export DB_NAME=cineverse
export DB_SSLMODE=disable

export OMDB_API_KEY=your_omdb_key
export JWT_SECRET=your_jwt_secret_minimum_32_chars

export SERVER_PORT=8080
export SERVER_TIMEOUT=30
```

5. **Run the API**
```bash
go run ./cmd/main.go
```

6. **Access Swagger**
```
http://localhost:8080/swagger/index.html
```

</details>

<details>
<summary><strong>Quick Test</strong></summary>

```bash
# Health check
curl http://localhost:8080/health

# Register user
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "email": "test@example.com",
    "display_name": "Test User",
    "password": "SecurePass123!"
  }'

# Login
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "password": "SecurePass123!"
  }'

# Search movies
curl "http://localhost:8080/api/v1/movies/search?q=matrix&page=1"

# Get movie details
curl "http://localhost:8080/api/v1/movies/tt0133093"
```

</details>

</details>

<details>
<summary><h3>⚙️ Configuration</h3></summary>

<a name="configuration-en"></a>

<details>
<summary><strong>Environment Variables</strong></summary>

#### Database Configuration
```bash
DB_HOST=localhost           # Database host
DB_PORT=5432                # Database port
DB_USER=cineverse           # Database user
DB_PASSWORD=password        # Database password
DB_NAME=cineverse           # Database name
DB_SSLMODE=disable          # SSL mode (disable, require, verify-ca, verify-full)
```

#### Redis Configuration
```bash
REDIS_HOST=localhost        # Redis host
REDIS_PORT=6379             # Redis port
REDIS_PASSWORD=             # Redis password (empty for no auth)
REDIS_DB=0                  # Redis database number
```

#### OMDb Configuration
```bash
OMDB_API_KEY=your_key       # OMDb API key (required)
OMDB_BASE_URL=http://www.omdbapi.com/  # OMDb base URL
```

#### JWT Configuration
```bash
JWT_SECRET=your_secret      # JWT signing secret (min 32 chars)
JWT_EXPIRATION=24h          # Token expiration (24h, 7d, etc)
```

#### Server Configuration
```bash
SERVER_PORT=8080            # HTTP server port
SERVER_TIMEOUT=30           # Request timeout in seconds
SERVER_HOST=0.0.0.0         # Server bind address
```

</details>

<details>
<summary><strong>Configuration File</strong></summary>

The application uses `internal/config/config.go` for centralized configuration management.

**Loading Priority:**
1. Environment variables
2. Default values
3. Configuration file (if implemented)

**Example usage:**
```go
cfg := config.Load()
db := setupDatabase(cfg.Database)
omdb := omdb.NewService(cfg.OMDb.APIKey)
```

</details>

</details>

<details>
<summary><h3>📁 Project Structure</h3></summary>

<a name="structure-en"></a>

<details>
<summary><strong>Layer Details</strong></summary>

#### Domain Layer (`internal/domain/`)
Contains business entities and repository interfaces.

**Files:**
- `auth.go`: Auth-related domain types
- `user.go`: User entity and UserRepository interface
- `movie.go`: Movie entity and MovieRepository interface

**Principles:**
- No external dependencies
- Pure business logic
- Framework-independent

#### Use Case Layer (`internal/usecase/`)
Orchestrates business operations using domain entities.

**Directories:**
- `auth/`: Registration, login, logout
- `user/`: Profile management
- `movie/`: Search, details, favorites

**Principles:**
- Single responsibility
- Dependency inversion
- Testable without infrastructure

#### Infrastructure Layer (`internal/infrastructure/`)
Implements external integrations.

**Files:**
- `jwt.go`: JWT token service
- `movie_provider.go`: Movie provider interface
- `movie_fetcher_chain.go`: Chain of Responsibility

**Principles:**
- Implements domain interfaces
- Handles external communication
- Error handling and retry logic

#### Handler Layer (`internal/handler/http/`)
HTTP request/response handling.

**Files:**
- `auth_handler.go`: Auth endpoints
- `user_handler.go`: User endpoints
- `movie_handler.go`: Movie endpoints
- `omdb_handler.go`: OMDb direct access

**Principles:**
- Request validation
- DTO conversion
- HTTP status code handling

</details>

</details>

<details>
<summary><h3>🧪 Testing</h3></summary>

<a name="testing-en"></a>

<details>
<summary><strong>Running Tests</strong></summary>

```bash
# Run all tests
go test ./...

# Run with coverage
go test -cover ./...

# Run specific package
go test ./internal/usecase/auth/...

# Verbose output
go test -v ./...

# Generate coverage report
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

</details>

<details>
<summary><strong>Manual Testing</strong></summary>

#### Test Authentication Flow
```bash
# 1. Register
REGISTER_RESPONSE=$(curl -s -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "email": "test@example.com",
    "display_name": "Test User",
    "password": "SecurePass123!"
  }')
echo $REGISTER_RESPONSE

# 2. Login
LOGIN_RESPONSE=$(curl -s -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "password": "SecurePass123!"
  }')
TOKEN=$(echo $LOGIN_RESPONSE | jq -r '.data.token')
echo "Token: $TOKEN"

# 3. Get Profile
curl -s http://localhost:8080/api/v1/users/me \
  -H "Authorization: Bearer $TOKEN" | jq '.'

# 4. Logout
curl -s -X POST http://localhost:8080/api/v1/auth/logout \
  -H "Authorization: Bearer $TOKEN" | jq '.'
```

#### Test Movie Chain
```bash
# Search movies (OMDb → Database chain)
curl -s "http://localhost:8080/api/v1/movies/search?q=matrix&page=1" | jq '.data[0]'

# Get movie details (OMDb → Database chain)
curl -s "http://localhost:8080/api/v1/movies/tt0133093" | jq '.data'

# Direct OMDb access
curl -s "http://localhost:8080/api/v1/omdb/tt0133093" | jq '.data'

# OMDb search
curl -s "http://localhost:8080/api/v1/omdb/search?q=batman&page=1" | jq '.data.results[0]'
```

</details>

</details>

<details>
<summary><h3>🚀 Deployment</h3></summary>

<a name="deployment-en"></a>

<details>
<summary><strong>Docker Production</strong></summary>

1. **Build image**
```bash
docker build -t cineverse-api:latest .
```

2. **Run container**
```bash
docker run -d \
  --name cineverse-api \
  -p 8080:8080 \
  -e DB_HOST=postgres \
  -e DB_PASSWORD=secure_password \
  -e OMDB_API_KEY=your_key \
  -e JWT_SECRET=your_secret \
  cineverse-api:latest
```

3. **Docker Compose**
```yaml
version: '3.8'
services:
  api:
    image: cineverse-api:latest
    ports:
      - "8080:8080"
    environment:
      - DB_HOST=postgres
      - DB_PASSWORD=${DB_PASSWORD}
      - OMDB_API_KEY=${OMDB_API_KEY}
      - JWT_SECRET=${JWT_SECRET}
    depends_on:
      - postgres
      - redis
```

</details>

<details>
<summary><strong>Production Checklist</strong></summary>

- [ ] Enable HTTPS/TLS
- [ ] Set strong JWT secret (32+ chars)
- [ ] Use strong database passwords
- [ ] Enable PostgreSQL SSL
- [ ] Configure CORS properly
- [ ] Set up monitoring (Prometheus, Grafana)
- [ ] Configure logging (structured logs)
- [ ] Set up backups (database, configs)
- [ ] Enable rate limiting
- [ ] Configure reverse proxy (Nginx, Traefik)
- [ ] Set up health checks
- [ ] Configure environment-specific settings
- [ ] Enable Redis password protection
- [ ] Set up CI/CD pipeline
- [ ] Configure secrets management

</details>

</details>

<details>
<summary><h3>📚 Additional Resources</h3></summary>

- [Main Project README](../README.md)
- [Architecture Documentation](ARCHITECTURE.md)
- [OMDb Integration Guide](OMDB_INTEGRATION.md)
- [Database Diagram](../database_diagram.dbml)
- [Swagger Documentation](http://localhost:8080/swagger/index.html)

</details>

</details>

---

<a name="português"></a>

<details>
<summary><h2>🇧🇷 VERSÃO EM PORTUGUÊS</h2></summary>

<details>
<summary><h3>📋 Índice</h3></summary>

- [Sobre](#sobre-pt)
- [Funcionalidades](#funcionalidades-pt)
- [Arquitetura](#arquitetura-pt)
- [Documentação da API](#api-docs-pt)
- [Primeiros Passos](#inicio-pt)
- [Configuração](#configuracao-pt)
- [Estrutura do Projeto](#estrutura-pt)
- [Desenvolvimento](#desenvolvimento-pt)
- [Testes](#testes-pt)
- [Implantação](#implantacao-pt)

</details>

<details>
<summary><h3>📖 Sobre</h3></summary>

<a name="sobre-pt"></a>

A API CineVerse é um backend RESTful robusto construído com Go, implementando princípios de Clean Architecture e projetado para escalabilidade e manutenibilidade.

**Recursos Principais:**
- 🔐 **Autenticação JWT Stateless** com hash de senha bcrypt
- 🎬 **Dados de Filmes Multi-Provedor** usando padrão Chain of Responsibility
- 📊 **Documentação Completa da API** com Swagger/OpenAPI
- 🗄️ **Banco de Dados PostgreSQL** com suporte a migrações
- ⚡ **Cache Redis** para otimização de performance
- 🌐 **Internacionalização** (Inglês/Português)
- 🐳 **Suporte Docker** para fácil implantação

**Destaques Técnicos:**
- Clean Architecture com clara separação de camadas
- Padrão Repository para acesso a dados
- Design orientado a casos de uso
- Injeção de dependências
- Logging estruturado com slog
- Tratamento de erros com tipos customizados
- Validação e sanitização de entrada

</details>

<details>
<summary><h3>✨ Funcionalidades</h3></summary>

<a name="funcionalidades-pt"></a>

<details>
<summary><strong>Autenticação & Autorização</strong></summary>

#### Registro de Usuário
- Validação de username (alfanumérico, 3-30 caracteres)
- Validação de formato de email
- Requisitos de senha forte
- Hash de senha com bcrypt (custo 12)
- Detecção de duplicação de username/email

#### Login de Usuário
- Validação de credenciais
- Geração de token JWT
- Expiração configurável de token
- Gerenciamento de sessão stateless

#### Rotas Protegidas
- Middleware JWT para proteção de rotas
- Validação e parsing de token
- Injeção de contexto de usuário
- Respostas automáticas de erro

</details>

<details>
<summary><strong>Gerenciamento de Usuários</strong></summary>

#### Operações de Perfil
- Obter usuário por ID
- Obter usuário por username
- Atualizar informações de perfil
- Gerenciamento de nome de exibição
- Funcionalidade de busca de usuário

#### Validação de Dados
- Unicidade de username
- Validação de formato de email
- Validação de campos obrigatórios
- Sanitização de dados

</details>

<details>
<summary><strong>Integração de Filmes</strong></summary>

#### Chain of Responsibility
```
Requisição → API OMDb → Cache DB → 404
              (80-90%)    (10-20%)   (<1%)
```

**Vantagens:**
- Reduz chamadas de API em 80-90%
- Persistência automática de dados
- TTL de cache de 48 horas
- Rastreamento de provedor (OMDb, Database)
- Mecanismo de fallback para confiabilidade

#### Operações de Filme
- Buscar filmes com paginação
- Obter detalhes de filme por ID
- Auto-save de filmes buscados
- Rastreamento de provedor e sincronização
- Suporte a gêneros e metadados

#### Integração OMDb
- Acesso direto à API OMDb
- Busca com filtros
- Informações detalhadas de filmes
- Dados de poster e avaliação
- Chave de API configurável

</details>

<details>
<summary><strong>Infraestrutura</strong></summary>

#### Banco de Dados
- PostgreSQL 15 com chaves primárias UUID
- Connection pooling
- Sistema de migração
- Suporte a JSONB para dados flexíveis
- Capacidades de busca full-text

#### Cache
- Integração Redis
- TTL configurável
- Suporte opcional de cache
- Tratamento de erros de conexão

#### Documentação da API
- Swagger/OpenAPI 2.0
- Explorador interativo de API
- Exemplos de request/response
- Definições de schema
- Documentação de autenticação

#### Middleware
- Configuração CORS
- Logging de requisições
- Validação de auth
- Tratamento de erros
- Recuperação de panics

</details>

</details>

<details>
<summary><h3>🏗️ Arquitetura</h3></summary>

<a name="arquitetura-pt"></a>

<details>
<summary><strong>Camadas da Clean Architecture</strong></summary>

```
┌─────────────────────────────────────────┐
│           Handlers HTTP                  │  ← Camada de Apresentação
│  (Request/Response, DTOs, Validação)     │
└──────────────┬──────────────────────────┘
               │
┌──────────────▼──────────────────────────┐
│           Casos de Uso                   │  ← Camada de Aplicação
│  (Lógica de Negócio, Orquestração)      │
└──────────────┬──────────────────────────┘
               │
┌──────────────▼──────────────────────────┐
│           Domínio                        │  ← Camada de Domínio
│  (Entidades, Interfaces, Regras)        │
└──────────────┬──────────────────────────┘
               │
┌──────────────▼──────────────────────────┐
│      Infraestrutura                      │  ← Camada de Infraestrutura
│  (Banco de Dados, APIs Externas)        │
└──────────────────────────────────────────┘
```

</details>

<details>
<summary><strong>Estrutura de Diretórios</strong></summary>

```
api_v2/
├── cmd/
│   └── main.go                 # Ponto de entrada da aplicação
├── internal/
│   ├── config/
│   │   └── config.go           # Gerenciamento de configuração
│   ├── domain/
│   │   ├── auth.go             # Entidade & interfaces de Auth
│   │   ├── movie.go            # Entidade & interfaces de Movie
│   │   └── user.go             # Entidade & interfaces de User
│   ├── dto/
│   │   ├── auth_dto.go         # Data transfer objects de Auth
│   │   ├── movie_dto.go        # Data transfer objects de Movie
│   │   └── user_dto.go         # Data transfer objects de User
│   ├── handler/
│   │   └── http/
│   │       ├── auth_handler.go # Handlers HTTP de Auth
│   │       ├── movie_handler.go# Handlers HTTP de Movie
│   │       ├── omdb_handler.go # Handlers HTTP de OMDb
│   │       └── user_handler.go # Handlers HTTP de User
│   ├── infrastructure/
│   │   ├── jwt.go              # Serviço de token JWT
│   │   ├── movie_provider.go   # Interface de provedor de filme
│   │   └── movie_fetcher_chain.go # Chain of Responsibility
│   ├── middleware/
│   │   ├── auth.go             # Middleware de auth JWT
│   │   ├── cors.go             # Middleware CORS
│   │   └── logger.go           # Logging de requisições
│   ├── repository/
│   │   ├── movie_repository.go # Acesso a dados de Movie
│   │   ├── session_repository.go# Acesso a dados de Session
│   │   └── user_repository.go  # Acesso a dados de User
│   ├── service/
│   │   └── omdb_service.go     # Cliente da API OMDb
│   ├── usecase/
│   │   ├── auth/               # Casos de uso de Auth
│   │   ├── movie/              # Casos de uso de Movie
│   │   └── user/               # Casos de uso de User
│   └── i18n/
│       ├── i18n.go             # Setup de internacionalização
│       └── locales/
│           ├── en.json         # Traduções em inglês
│           └── pt.json         # Traduções em português
├── migrations/
│   ├── 001_clean_initial_schema.sql  # Schema inicial
│   └── 002_add_provider_and_sync.sql # Rastreamento de provedor
├── docs/
│   ├── docs.go                 # Docs Swagger gerados
│   ├── swagger.json            # Spec OpenAPI JSON
│   └── swagger.yaml            # Spec OpenAPI YAML
└── README.md                   # Este arquivo
```

</details>

<details>
<summary><strong>Padrões de Design</strong></summary>

#### Chain of Responsibility
Usado para busca de dados de filmes com fallback automático:
- **OMDbMovieFetcher**: Fonte de dados primária
- **DatabaseMovieFetcher**: Cache de fallback
- **Auto-save**: Persistir dados OMDb no banco

#### Padrão Repository
Abstrai lógica de acesso a dados:
- Contratos baseados em interface
- Fácil teste com mocks
- Independência de implementação de banco

#### Injeção de Dependências
Serviços e repositórios injetados via construtores:
- Baixo acoplamento entre camadas
- Testabilidade
- Flexibilidade para trocar implementações

#### Padrão DTO
Separa entidades internas de contratos de API:
- Validação de entrada
- Formatação de resposta
- Suporte a versionamento

</details>

<details>
<summary><strong>Schema do Banco de Dados</strong></summary>

#### Tabela Users
```sql
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    username VARCHAR(30) UNIQUE NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    display_name VARCHAR(100),
    password_hash TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

#### Tabela Movies
```sql
CREATE TABLE movies (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    external_api_id VARCHAR(50) UNIQUE,
    title VARCHAR(255) NOT NULL,
    overview TEXT,
    poster_url TEXT,
    genres JSONB,
    adult BOOLEAN DEFAULT false,
    provider VARCHAR(50),
    last_sync_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

#### Tabela Sessions
```sql
CREATE TABLE sessions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    token TEXT NOT NULL,
    expires_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

</details>

</details>

<details>
<summary><h3>📚 Documentação da API</h3></summary>

<a name="api-docs-pt"></a>

<details>
<summary><strong>Swagger/OpenAPI</strong></summary>

A API é totalmente documentada usando Swagger/OpenAPI 2.0.

**Acessar Swagger UI:**
```
http://localhost:8080/swagger/index.html
```

**Recursos:**
- Explorador interativo de API
- Funcionalidade try-it-out
- Exemplos de request/response
- Definições de schema
- Documentação de autenticação

**Endpoints:**
- 14 endpoints documentados
- Schemas completos de request/response
- Documentação de códigos de erro
- Payloads de exemplo

</details>

<details>
<summary><strong>Endpoints de Autenticação</strong></summary>

#### POST /api/v1/auth/register
Registrar uma nova conta de usuário.

**Requisição:**
```json
{
  "username": "johndoe",
  "email": "john@example.com",
  "display_name": "John Doe",
  "password": "SecurePass123!"
}
```

**Resposta (201):**
```json
{
  "success": true,
  "message": "User registered successfully",
  "data": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "username": "johndoe",
    "email": "john@example.com",
    "display_name": "John Doe"
  }
}
```

#### POST /api/v1/auth/login
Autenticar e receber token JWT.

**Requisição:**
```json
{
  "username": "johndoe",
  "password": "SecurePass123!"
}
```

**Resposta (200):**
```json
{
  "success": true,
  "message": "Login successful",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIs...",
    "user": {
      "id": "550e8400-e29b-41d4-a716-446655440000",
      "username": "johndoe",
      "display_name": "John Doe"
    }
  }
}
```

#### POST /api/v1/auth/logout
Invalidar sessão atual (requer autenticação).

</details>

<details>
<summary><strong>Endpoints de Filmes</strong></summary>

#### GET /api/v1/movies/search
Buscar filmes usando Chain of Responsibility (OMDb → Database).

**Parâmetros:**
- `q` (query string, obrigatório): Termo de busca
- `page` (integer, opcional): Número da página (padrão: 1)

**Exemplo:**
```bash
curl "http://localhost:8080/api/v1/movies/search?q=matrix&page=1"
```

**Resposta (200):**
```json
{
  "success": true,
  "message": "Movies found",
  "data": [
    {
      "id": "550e8400-e29b-41d4-a716-446655440000",
      "external_api_id": "tt0133093",
      "title": "The Matrix",
      "overview": "1999",
      "poster_url": "https://m.media-amazon.com/...",
      "genres": null,
      "adult": false,
      "provider": "OMDb",
      "last_sync_at": "2025-11-08T18:40:19Z"
    }
  ]
}
```

#### GET /api/v1/movies/{id}
Obter detalhes de filme por ID (cadeia OMDb → Database).

**Exemplo:**
```bash
curl "http://localhost:8080/api/v1/movies/tt0133093"
```

</details>

<details>
<summary><strong>Endpoints Diretos OMDb</strong></summary>

#### GET /api/v1/omdb/{imdb_id}
Obter detalhes de filme diretamente da API OMDb.

**Exemplo:**
```bash
curl "http://localhost:8080/api/v1/omdb/tt0133093"
```

#### GET /api/v1/omdb/search
Buscar filmes diretamente na API OMDb.

**Parâmetros:**
- `q` (query string, obrigatório): Termo de busca
- `page` (integer, opcional): Número da página

**Exemplo:**
```bash
curl "http://localhost:8080/api/v1/omdb/search?q=batman&page=1"
```

</details>

<details>
<summary><strong>Endpoints de Usuário</strong></summary>

#### GET /api/v1/users/me
Obter perfil do usuário autenticado atual (requer autenticação).

**Headers:**
```
Authorization: Bearer <token>
```

**Resposta (200):**
```json
{
  "success": true,
  "data": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "username": "johndoe",
    "email": "john@example.com",
    "display_name": "John Doe",
    "created_at": "2025-11-08T15:30:00Z"
  }
}
```

#### GET /api/v1/users/{username}
Obter perfil de usuário por username.

</details>

<details>
<summary><strong>Health Check</strong></summary>

#### GET /health
Verificar status de saúde da API.

**Resposta (200):**
```json
{
  "status": "healthy",
  "service": "cineverse-api"
}
```

</details>

<details>
<summary><strong>Respostas de Erro</strong></summary>

Todos os endpoints retornam respostas de erro consistentes:

**400 Bad Request:**
```json
{
  "success": false,
  "message": "Invalid input data",
  "error": "username must be between 3 and 30 characters"
}
```

**401 Unauthorized:**
```json
{
  "success": false,
  "message": "Authentication required",
  "error": "Missing or invalid token"
}
```

**404 Not Found:**
```json
{
  "success": false,
  "message": "Resource not found",
  "error": "Movie not found"
}
```

**500 Internal Server Error:**
```json
{
  "success": false,
  "message": "Internal server error",
  "error": "Database connection failed"
}
```

</details>

</details>

<details>
<summary><h3>🚀 Primeiros Passos</h3></summary>

<a name="inicio-pt"></a>

<details>
<summary><strong>Pré-requisitos</strong></summary>

- **Go** 1.21 ou superior
- **PostgreSQL** 15
- **Redis** 7 (opcional mas recomendado)
- **Chave da API OMDb** (grátis em [omdbapi.com](https://www.omdbapi.com/apikey.aspx))
- **Docker & Docker Compose** (para setup containerizado)

</details>

<details>
<summary><strong>Instalação (Docker)</strong></summary>

1. **Clonar repositório**
```bash
git clone https://github.com/EduardoMG12/cine.git
cd cine/api_v2
```

2. **Configurar ambiente**
```bash
cp .env.example .env
# Edite .env com sua configuração
```

3. **Iniciar serviços**
```bash
docker-compose up -d
```

4. **Executar migrações**
```bash
docker-compose exec api psql -U cineverse -d cineverse -f /app/migrations/001_clean_initial_schema.sql
docker-compose exec api psql -U cineverse -d cineverse -f /app/migrations/002_add_provider_and_sync.sql
```

5. **Verificar saúde**
```bash
curl http://localhost:8080/health
```

</details>

<details>
<summary><strong>Instalação (Local)</strong></summary>

1. **Instalar dependências**
```bash
go mod download
```

2. **Configurar PostgreSQL**
```bash
psql -U postgres
CREATE DATABASE cineverse;
CREATE USER cineverse WITH PASSWORD 'sua_senha';
GRANT ALL PRIVILEGES ON DATABASE cineverse TO cineverse;
\q
```

3. **Executar migrações**
```bash
psql -U cineverse -d cineverse -f migrations/001_clean_initial_schema.sql
psql -U cineverse -d cineverse -f migrations/002_add_provider_and_sync.sql
```

4. **Configurar ambiente**
```bash
export DB_HOST=localhost
export DB_PORT=5432
export DB_USER=cineverse
export DB_PASSWORD=sua_senha
export DB_NAME=cineverse
export DB_SSLMODE=disable

export OMDB_API_KEY=sua_chave_omdb
export JWT_SECRET=seu_jwt_secret_minimo_32_chars

export SERVER_PORT=8080
export SERVER_TIMEOUT=30
```

5. **Executar a API**
```bash
go run ./cmd/main.go
```

6. **Acessar Swagger**
```
http://localhost:8080/swagger/index.html
```

</details>

<details>
<summary><strong>Teste Rápido</strong></summary>

```bash
# Health check
curl http://localhost:8080/health

# Registrar usuário
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "email": "test@example.com",
    "display_name": "Test User",
    "password": "SecurePass123!"
  }'

# Login
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "password": "SecurePass123!"
  }'

# Buscar filmes
curl "http://localhost:8080/api/v1/movies/search?q=matrix&page=1"

# Obter detalhes de filme
curl "http://localhost:8080/api/v1/movies/tt0133093"
```

</details>

</details>

<details>
<summary><h3>⚙️ Configuração</h3></summary>

<a name="configuracao-pt"></a>

<details>
<summary><strong>Variáveis de Ambiente</strong></summary>

#### Configuração do Banco de Dados
```bash
DB_HOST=localhost           # Host do banco
DB_PORT=5432                # Porta do banco
DB_USER=cineverse           # Usuário do banco
DB_PASSWORD=password        # Senha do banco
DB_NAME=cineverse           # Nome do banco
DB_SSLMODE=disable          # Modo SSL (disable, require, verify-ca, verify-full)
```

#### Configuração Redis
```bash
REDIS_HOST=localhost        # Host Redis
REDIS_PORT=6379             # Porta Redis
REDIS_PASSWORD=             # Senha Redis (vazio para sem auth)
REDIS_DB=0                  # Número do banco Redis
```

#### Configuração OMDb
```bash
OMDB_API_KEY=sua_chave      # Chave da API OMDb (obrigatório)
OMDB_BASE_URL=http://www.omdbapi.com/  # URL base OMDb
```

#### Configuração JWT
```bash
JWT_SECRET=seu_secret       # Secret de assinatura JWT (mín 32 chars)
JWT_EXPIRATION=24h          # Expiração do token (24h, 7d, etc)
```

#### Configuração do Servidor
```bash
SERVER_PORT=8080            # Porta do servidor HTTP
SERVER_TIMEOUT=30           # Timeout de requisição em segundos
SERVER_HOST=0.0.0.0         # Endereço de bind do servidor
```

</details>

<details>
<summary><strong>Arquivo de Configuração</strong></summary>

A aplicação usa `internal/config/config.go` para gerenciamento centralizado de configuração.

**Prioridade de Carregamento:**
1. Variáveis de ambiente
2. Valores padrão
3. Arquivo de configuração (se implementado)

**Exemplo de uso:**
```go
cfg := config.Load()
db := setupDatabase(cfg.Database)
omdb := omdb.NewService(cfg.OMDb.APIKey)
```

</details>

</details>

<details>
<summary><h3>📁 Estrutura do Projeto</h3></summary>

<a name="estrutura-pt"></a>

<details>
<summary><strong>Detalhes das Camadas</strong></summary>

#### Camada de Domínio (`internal/domain/`)
Contém entidades de negócio e interfaces de repositórios.

**Arquivos:**
- `auth.go`: Tipos de domínio relacionados a Auth
- `user.go`: Entidade User e interface UserRepository
- `movie.go`: Entidade Movie e interface MovieRepository

**Princípios:**
- Sem dependências externas
- Lógica de negócio pura
- Independente de framework

#### Camada de Caso de Uso (`internal/usecase/`)
Orquestra operações de negócio usando entidades de domínio.

**Diretórios:**
- `auth/`: Registro, login, logout
- `user/`: Gerenciamento de perfil
- `movie/`: Busca, detalhes, favoritos

**Princípios:**
- Responsabilidade única
- Inversão de dependência
- Testável sem infraestrutura

#### Camada de Infraestrutura (`internal/infrastructure/`)
Implementa integrações externas.

**Arquivos:**
- `jwt.go`: Serviço de token JWT
- `movie_provider.go`: Interface de provedor de filme
- `movie_fetcher_chain.go`: Chain of Responsibility

**Princípios:**
- Implementa interfaces de domínio
- Lida com comunicação externa
- Tratamento de erros e lógica de retry

#### Camada de Handler (`internal/handler/http/`)
Tratamento de request/response HTTP.

**Arquivos:**
- `auth_handler.go`: Endpoints de Auth
- `user_handler.go`: Endpoints de User
- `movie_handler.go`: Endpoints de Movie
- `omdb_handler.go`: Acesso direto OMDb

**Princípios:**
- Validação de requisição
- Conversão de DTO
- Tratamento de código de status HTTP

</details>

</details>

<details>
<summary><h3>🧪 Testes</h3></summary>

<a name="testes-pt"></a>

<details>
<summary><strong>Executando Testes</strong></summary>

```bash
# Executar todos os testes
go test ./...

# Executar com cobertura
go test -cover ./...

# Executar pacote específico
go test ./internal/usecase/auth/...

# Saída verbosa
go test -v ./...

# Gerar relatório de cobertura
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

</details>

<details>
<summary><strong>Testes Manuais</strong></summary>

#### Testar Fluxo de Autenticação
```bash
# 1. Registrar
REGISTER_RESPONSE=$(curl -s -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "email": "test@example.com",
    "display_name": "Test User",
    "password": "SecurePass123!"
  }')
echo $REGISTER_RESPONSE

# 2. Login
LOGIN_RESPONSE=$(curl -s -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "password": "SecurePass123!"
  }')
TOKEN=$(echo $LOGIN_RESPONSE | jq -r '.data.token')
echo "Token: $TOKEN"

# 3. Obter Perfil
curl -s http://localhost:8080/api/v1/users/me \
  -H "Authorization: Bearer $TOKEN" | jq '.'

# 4. Logout
curl -s -X POST http://localhost:8080/api/v1/auth/logout \
  -H "Authorization: Bearer $TOKEN" | jq '.'
```

#### Testar Cadeia de Filmes
```bash
# Buscar filmes (cadeia OMDb → Database)
curl -s "http://localhost:8080/api/v1/movies/search?q=matrix&page=1" | jq '.data[0]'

# Obter detalhes de filme (cadeia OMDb → Database)
curl -s "http://localhost:8080/api/v1/movies/tt0133093" | jq '.data'

# Acesso direto OMDb
curl -s "http://localhost:8080/api/v1/omdb/tt0133093" | jq '.data'

# Busca OMDb
curl -s "http://localhost:8080/api/v1/omdb/search?q=batman&page=1" | jq '.data.results[0]'
```

</details>

</details>

<details>
<summary><h3>🚀 Implantação</h3></summary>

<a name="implantacao-pt"></a>

<details>
<summary><strong>Docker Produção</strong></summary>

1. **Construir imagem**
```bash
docker build -t cineverse-api:latest .
```

2. **Executar container**
```bash
docker run -d \
  --name cineverse-api \
  -p 8080:8080 \
  -e DB_HOST=postgres \
  -e DB_PASSWORD=senha_segura \
  -e OMDB_API_KEY=sua_chave \
  -e JWT_SECRET=seu_secret \
  cineverse-api:latest
```

3. **Docker Compose**
```yaml
version: '3.8'
services:
  api:
    image: cineverse-api:latest
    ports:
      - "8080:8080"
    environment:
      - DB_HOST=postgres
      - DB_PASSWORD=${DB_PASSWORD}
      - OMDB_API_KEY=${OMDB_API_KEY}
      - JWT_SECRET=${JWT_SECRET}
    depends_on:
      - postgres
      - redis
```

</details>

<details>
<summary><strong>Checklist de Produção</strong></summary>

- [ ] Habilitar HTTPS/TLS
- [ ] Definir JWT secret forte (32+ chars)
- [ ] Usar senhas fortes de banco de dados
- [ ] Habilitar SSL do PostgreSQL
- [ ] Configurar CORS adequadamente
- [ ] Configurar monitoramento (Prometheus, Grafana)
- [ ] Configurar logging (logs estruturados)
- [ ] Configurar backups (banco, configs)
- [ ] Habilitar rate limiting
- [ ] Configurar reverse proxy (Nginx, Traefik)
- [ ] Configurar health checks
- [ ] Configurar settings específicos de ambiente
- [ ] Habilitar proteção de senha Redis
- [ ] Configurar pipeline CI/CD
- [ ] Configurar gerenciamento de secrets

</details>

</details>

<details>
<summary><h3>📚 Recursos Adicionais</h3></summary>

- [README Principal do Projeto](../README.md)
- [Documentação de Arquitetura](ARCHITECTURE.md)
- [Guia de Integração OMDb](OMDB_INTEGRATION.md)
- [Diagrama do Banco de Dados](../database_diagram.dbml)
- [Documentação Swagger](http://localhost:8080/swagger/index.html)

</details>

</details>

---

<div align="center">

**Desenvolvido para IFPR Campus Palmas-PR**

**CineVerse API** • **v1.0** • **2025**

</div>
