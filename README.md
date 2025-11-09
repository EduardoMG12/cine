<div align="center">

<!-- Space for IFPR Logo -->
<img src="https://via.placeholder.com/200x100/003366/FFFFFF?text=IFPR+Logo" alt="IFPR Logo" width="200"/>

# 🎬 CineVerse

### Social Network Platform for Movie Enthusiasts

[![Go](https://img.shields.io/badge/Go-1.21+-00ADD8?style=for-the-badge&logo=go&logoColor=white)](https://go.dev/)
[![Flutter](https://img.shields.io/badge/Flutter-3.0+-02569B?style=for-the-badge&logo=flutter&logoColor=white)](https://flutter.dev/)
[![PostgreSQL](https://img.shields.io/badge/PostgreSQL-15-336791?style=for-the-badge&logo=postgresql&logoColor=white)](https://www.postgresql.org/)
[![Docker](https://img.shields.io/badge/Docker-Ready-2496ED?style=for-the-badge&logo=docker&logoColor=white)](https://www.docker.com/)
[![License](https://img.shields.io/badge/License-MIT-green?style=for-the-badge)](LICENSE)

**Academic Project - IFPR Campus Palmas-PR**  
**Software Engineering - 2025/1**

[🇺🇸 English](#english) | [🇧🇷 Português](#português)

</div>

---

<a name="english"></a>

<details open>
<summary><h2>🇺🇸 ENGLISH VERSION</h2></summary>

<details>
<summary><h3>📋 Table of Contents</h3></summary>

- [About the Project](#about-en)
- [Academic Context](#academic-en)
- [Architecture](#architecture-en)
- [Technologies](#technologies-en)
- [Project Structure](#structure-en)
- [Features](#features-en)
- [Getting Started](#getting-started-en)
- [Development Status](#status-en)
- [Roadmap](#roadmap-en)
- [Team](#team-en)
- [License](#license-en)

</details>

<details>
<summary><h3>📖 About the Project</h3></summary>

<a name="about-en"></a>

CineVerse is a social network platform designed for movie enthusiasts to discover, share, and discuss films. The project implements a modern architecture with:

- **RESTful API Backend** built with Go
- **Flutter Mobile App** (in development)
- **PostgreSQL Database** for data persistence
- **Redis Cache** for performance optimization
- **OMDb API Integration** for movie data

**Key Differentiators:**
- 🔗 **Chain of Responsibility Pattern** for multi-provider movie data fetching
- 🎯 **Clean Architecture** with clear separation of concerns
- 🔐 **JWT Authentication** with stateless session management
- 📊 **Comprehensive API Documentation** with Swagger/OpenAPI
- 🌐 **Internationalization** support (English/Portuguese)

</details>

<details>
<summary><h3>🎓 Academic Context</h3></summary>

<a name="academic-en"></a>

**Institution:** Federal Institute of Paraná (IFPR) - Campus Palmas-PR  
**Course:** Software Engineering  
**Semester:** 2025/1  
**Project Type:** Academic Development Project

**Objectives:**
- Apply software engineering best practices
- Implement clean architecture and design patterns
- Develop a full-stack application with modern technologies
- Create comprehensive technical documentation
- Practice agile development methodologies

**Evaluation Criteria:**
- Code quality and organization
- Architecture and design patterns
- API documentation
- Testing coverage
- Version control and collaboration

</details>

<details>
<summary><h3>🏗️ Architecture</h3></summary>

<a name="architecture-en"></a>

#### System Overview

```
┌─────────────────┐
│  Flutter App    │  ◄── Mobile Application (In Development)
│  (Frontend)     │
└────────┬────────┘
         │ HTTP/JSON
         ▼
┌─────────────────┐
│   Go API        │  ◄── RESTful API (Chi Router)
│   (Backend)     │      • Authentication & Authorization
└────────┬────────┘      • Business Logic
         │                • Data Validation
         ├──────────┬─────────┬─────────┐
         ▼          ▼         ▼         ▼
    ┌────────┐ ┌────────┐ ┌────────┐ ┌────────┐
    │  DB    │ │ Redis  │ │ OMDb   │ │ Email  │
    │  PG15  │ │ Cache  │ │  API   │ │ SMTP   │
    └────────┘ └────────┘ └────────┘ └────────┘
```

#### Backend Layers (Clean Architecture)

<details>
<summary><strong>Domain Layer</strong></summary>

- **Entities:** User, Movie, Session, Friendship
- **Interfaces:** Repository contracts
- **Business rules:** Independent of frameworks

</details>

<details>
<summary><strong>Use Case Layer</strong></summary>

- **User Operations:** Register, Login, Profile Management
- **Movie Operations:** Search, Details, Favorites
- **Social Features:** Friends, Following, Matching
- **Auth Operations:** Token generation, Session management

</details>

<details>
<summary><strong>Infrastructure Layer</strong></summary>

- **HTTP Handlers:** Request/Response mapping
- **Repositories:** Database implementations
- **External Services:** OMDb API, Email service
- **Middleware:** Auth, Logging, CORS, Rate limiting

</details>

<details>
<summary><strong>Movie Data Chain</strong></summary>

```
Request → OMDb API → Database Cache → 404
          (80-90%)    (10-20%)         (<1%)
```

The system uses a **Chain of Responsibility** pattern:
1. First tries OMDb API (external data)
2. Falls back to database (cached data)
3. Returns 404 if not found

**Auto-Save:** Movies fetched from OMDb are automatically saved to the database, reducing future API calls by 80-90%.

</details>

</details>

<details>
<summary><h3>💻 Technologies</h3></summary>

<a name="technologies-en"></a>

#### Backend
| Technology | Version | Purpose |
|------------|---------|---------|
| Go | 1.21+ | Primary language |
| Chi Router | v5 | HTTP routing |
| PostgreSQL | 15 | Main database |
| SQLx | Latest | SQL toolkit |
| Redis | 7 | Caching layer |
| JWT | - | Authentication |
| Bcrypt | - | Password hashing |
| Swagger | 2.0 | API documentation |

#### Frontend (In Development)
| Technology | Version | Purpose |
|------------|---------|---------|
| Flutter | 3.0+ | Mobile framework |
| Dart | Latest | Programming language |
| Provider | Latest | State management (planned) |

#### DevOps
| Technology | Purpose |
|------------|---------|
| Docker | Containerization |
| Docker Compose | Multi-container orchestration |
| Git | Version control |
| GitHub | Code hosting |

</details>

<details>
<summary><h3>📁 Project Structure</h3></summary>

<a name="structure-en"></a>

```
cine/
├── api_v2/                    # Go Backend (Current)
│   ├── cmd/
│   │   └── main.go           # Application entry point
│   ├── internal/
│   │   ├── config/           # Configuration management
│   │   ├── domain/           # Business entities
│   │   ├── dto/              # Data Transfer Objects
│   │   ├── handler/http/     # HTTP handlers
│   │   ├── middleware/       # HTTP middleware
│   │   ├── repository/       # Data access layer
│   │   ├── service/          # Business logic
│   │   ├── usecase/          # Application use cases
│   │   ├── infrastructure/   # External integrations
│   │   └── i18n/             # Internationalization
│   ├── migrations/           # Database migrations
│   ├── docs/                 # Swagger documentation
│   └── README.md             # Backend documentation
│
├── flutter_app/              # Flutter Frontend (In Development)
│   ├── lib/
│   │   └── src/              # Application source
│   └── README.md
│
├── scripts/                  # Utility scripts
│   ├── setup.sh              # Environment setup
│   └── go-lint.sh            # Code linting
│
├── docker-compose.yml        # Container orchestration
└── README.md                 # This file
```

</details>

<details>
<summary><h3>✨ Features</h3></summary>

<a name="features-en"></a>

<details>
<summary><strong>✅ Implemented Features</strong></summary>

#### Authentication & Authorization
- ✅ User registration with validation
- ✅ Login with JWT token generation
- ✅ Password hashing with bcrypt
- ✅ Session management
- ✅ Protected routes with middleware
- ✅ Token refresh mechanism

#### User Management
- ✅ User profiles with display names
- ✅ Username uniqueness validation
- ✅ Email validation
- ✅ Profile retrieval
- ✅ User search

#### Movie Integration
- ✅ OMDb API integration
- ✅ Movie search with pagination
- ✅ Movie details by ID
- ✅ Chain of Responsibility pattern
- ✅ Database caching (48h TTL)
- ✅ Auto-save after API fetch
- ✅ Provider tracking (OMDb, Database)

#### Infrastructure
- ✅ PostgreSQL database with migrations
- ✅ Redis caching support
- ✅ Docker containerization
- ✅ Health check endpoint
- ✅ CORS configuration
- ✅ Request logging
- ✅ Error handling
- ✅ Swagger API documentation (14 endpoints)

</details>

<details>
<summary><strong>🚧 In Development</strong></summary>

#### Mobile Application
- 🚧 Flutter app structure
- 🚧 UI/UX design
- 🚧 API integration
- 🚧 State management
- 🚧 Offline support

#### Enhanced Features
- 🚧 Email confirmation
- 🚧 Password reset flow
- 🚧 User profile editing
- 🚧 Avatar upload
- 🚧 Social features (friends, following)
- 🚧 Movie lists (watched, favorites, wishlist)
- 🚧 Movie ratings and reviews
- 🚧 Notification system
- 🚧 Real-time chat

</details>

</details>

<details>
<summary><h3>🚀 Getting Started</h3></summary>

<a name="getting-started-en"></a>

<details>
<summary><strong>Prerequisites</strong></summary>

- Go 1.21 or higher
- PostgreSQL 15
- Redis 7 (optional but recommended)
- Docker & Docker Compose (for containerized setup)
- OMDb API Key (get free at [omdbapi.com](https://www.omdbapi.com/apikey.aspx))

</details>

<details>
<summary><strong>Using Docker (Recommended)</strong></summary>

1. **Clone the repository**
```bash
git clone https://github.com/EduardoMG12/cine.git
cd cine
```

2. **Configure environment**
```bash
cp api_v2/.env.example api_v2/.env
# Edit api_v2/.env with your OMDb API key and database credentials
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

5. **Access the API**
- API: http://localhost:8080
- Swagger UI: http://localhost:8080/swagger/index.html
- Health Check: http://localhost:8080/health

</details>

<details>
<summary><strong>Local Development Setup</strong></summary>

1. **Install Go dependencies**
```bash
cd api_v2
go mod download
```

2. **Configure database**
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

4. **Set environment variables**
```bash
export DB_HOST=localhost
export DB_PORT=5432
export DB_USER=cineverse
export DB_PASSWORD=your_password
export DB_NAME=cineverse
export OMDB_API_KEY=your_omdb_key
export JWT_SECRET=your_jwt_secret
```

5. **Run the API**
```bash
go run ./cmd/main.go
```

</details>

<details>
<summary><strong>Testing the API</strong></summary>

#### Register a new user
```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "johndoe",
    "email": "john@example.com",
    "display_name": "John Doe",
    "password": "SecurePass123!"
  }'
```

#### Login
```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "johndoe",
    "password": "SecurePass123!"
  }'
```

#### Search movies
```bash
curl "http://localhost:8080/api/v1/movies/search?q=matrix&page=1"
```

#### Get movie details
```bash
curl "http://localhost:8080/api/v1/movies/tt0133093"
```

</details>

</details>

<details>
<summary><h3>📊 Development Status</h3></summary>

<a name="status-en"></a>

#### Backend Progress
```
████████████████░░░░  80% Complete
```
- ✅ Core API structure
- ✅ Authentication system
- ✅ Movie integration (OMDb)
- ✅ Database design
- ✅ API documentation
- 🚧 Social features
- 🚧 Notification system
- 🚧 Advanced matching

#### Frontend Progress
```
███░░░░░░░░░░░░░░░░░  15% Complete
```
- ✅ Project structure
- ✅ Basic navigation
- 🚧 Authentication screens
- 🚧 Home feed
- 🚧 Movie details
- 🚧 User profiles
- 🚧 Social interactions

#### Documentation Progress
```
████████████░░░░░░░░  60% Complete
```
- ✅ API documentation (Swagger)
- ✅ Backend README
- ✅ Root README
- ✅ Architecture guide
- 🚧 User manual
- 🚧 Deployment guide
- 🚧 Contributing guide

#### Testing Progress
```
████░░░░░░░░░░░░░░░░  20% Complete
```
- ✅ Manual endpoint testing
- 🚧 Unit tests
- 🚧 Integration tests
- 🚧 E2E tests
- 🚧 Load testing

</details>

<details>
<summary><h3>🗺️ Roadmap</h3></summary>

<a name="roadmap-en"></a>

<details>
<summary><strong>Sprint 1 - Foundation</strong> ✅ Complete</summary>

- [x] Project structure setup
- [x] Database schema design
- [x] Authentication system
- [x] Basic user management
- [x] OMDb integration
- [x] API documentation

</details>

<details>
<summary><strong>Sprint 2 - Core Features</strong> 🚧 In Progress</summary>

- [ ] Email confirmation system
- [ ] Password reset flow
- [ ] User profile editing
- [ ] Custom movie lists (watched, favorites, wishlist)
- [ ] Friend system (add, remove, list)
- [ ] Basic notification system

</details>

<details>
<summary><strong>Sprint 3 - Social Features</strong> 📅 Planned</summary>

- [ ] Movie matching algorithm
- [ ] Shared watchlists
- [ ] Movie recommendations
- [ ] User following system
- [ ] Activity feed
- [ ] Social sharing

</details>

<details>
<summary><strong>Sprint 4 - Mobile App</strong> 📅 Planned</summary>

- [ ] Flutter app authentication
- [ ] Home feed implementation
- [ ] Movie search and details
- [ ] User profiles
- [ ] Social interactions
- [ ] Offline support

</details>

<details>
<summary><strong>Sprint 5 - Advanced Features</strong> 📅 Planned</summary>

- [ ] Real-time chat system
- [ ] Movie reviews and ratings
- [ ] Advanced search filters
- [ ] Personalized recommendations (ML)
- [ ] Watchlist notifications
- [ ] Integration with streaming platforms

</details>

<details>
<summary><strong>Sprint 6 - Polish & Deploy</strong> 📅 Planned</summary>

- [ ] Comprehensive testing
- [ ] Performance optimization
- [ ] Security audit
- [ ] Production deployment
- [ ] Monitoring and logging
- [ ] User documentation

</details>

</details>

<details>
<summary><h3>👥 Team</h3></summary>

<a name="team-en"></a>

**Project Lead & Developer:** Charles Eduardo Mello Guimaraes  and Willian Fragata
**Institution:** IFPR Campus Palmas-PR  
**Course:** Information System  
**GitHub:** [@EduardoMG12](https://github.com/EduardoMG12)

**Advisor:** Alexis Kang 
**Course:** Software Engineering - IFPR

</details>

<details>
<summary><h3>📄 License</h3></summary>

<a name="license-en"></a>

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

**Academic Use:** This project is part of an academic assignment at IFPR Campus Palmas-PR. Feel free to use it as a reference for learning purposes.

</details>

<details>
<summary><h3>📚 Additional Resources</h3></summary>

- [Backend Documentation](api_v2/README.md)
- [Architecture Guide](api_v2/ARCHITECTURE.md)
- [OMDb Integration](api_v2/OMDB_INTEGRATION.md)
- [API Swagger Documentation](http://localhost:8080/swagger/index.html)
- [Database Diagram](database_diagram.dbml)

</details>

</details>

---

<a name="português"></a>

<details>
<summary><h2>🇧🇷 VERSÃO EM PORTUGUÊS</h2></summary>

<details>
<summary><h3>📋 Índice</h3></summary>

- [Sobre o Projeto](#sobre-pt)
- [Contexto Acadêmico](#academico-pt)
- [Arquitetura](#arquitetura-pt)
- [Tecnologias](#tecnologias-pt)
- [Estrutura do Projeto](#estrutura-pt)
- [Funcionalidades](#funcionalidades-pt)
- [Primeiros Passos](#inicio-pt)
- [Status do Desenvolvimento](#status-pt)
- [Roadmap](#roadmap-pt)
- [Equipe](#equipe-pt)
- [Licença](#licenca-pt)

</details>

<details>
<summary><h3>📖 Sobre o Projeto</h3></summary>

<a name="sobre-pt"></a>

CineVerse é uma plataforma de rede social projetada para entusiastas de cinema descobrirem, compartilharem e discutirem filmes. O projeto implementa uma arquitetura moderna com:

- **Backend API RESTful** construído com Go
- **Aplicativo Mobile Flutter** (em desenvolvimento)
- **Banco de Dados PostgreSQL** para persistência de dados
- **Cache Redis** para otimização de performance
- **Integração com API OMDb** para dados de filmes

**Diferenciais:**
- 🔗 **Padrão Chain of Responsibility** para busca de dados de filmes multi-provedor
- 🎯 **Clean Architecture** com clara separação de responsabilidades
- 🔐 **Autenticação JWT** com gerenciamento de sessão stateless
- 📊 **Documentação Completa da API** com Swagger/OpenAPI
- 🌐 **Suporte a Internacionalização** (Inglês/Português)

</details>

<details>
<summary><h3>🎓 Contexto Acadêmico</h3></summary>

<a name="academico-pt"></a>

**Instituição:** Instituto Federal do Paraná (IFPR) - Campus Palmas-PR  
**Curso:** sistema de informações  
**Semestre:** 2025/1  
**Tipo de Projeto:** Projeto de Desenvolvimento Acadêmico

**Objetivos:**
- Aplicar boas práticas de engenharia de software
- Implementar arquitetura limpa e padrões de design
- Desenvolver uma aplicação full-stack com tecnologias modernas
- Criar documentação técnica abrangente
- Praticar metodologias ágeis de desenvolvimento

**Critérios de Avaliação:**
- Qualidade e organização do código
- Arquitetura e padrões de design
- Documentação da API
- Cobertura de testes
- Controle de versão e colaboração

</details>

<details>
<summary><h3>🏗️ Arquitetura</h3></summary>

<a name="arquitetura-pt"></a>

#### Visão Geral do Sistema

```
┌─────────────────┐
│  App Flutter    │  ◄── Aplicação Mobile (Em Desenvolvimento)
│  (Frontend)     │
└────────┬────────┘
         │ HTTP/JSON
         ▼
┌─────────────────┐
│   API Go        │  ◄── API RESTful (Chi Router)
│   (Backend)     │      • Autenticação & Autorização
└────────┬────────┘      • Lógica de Negócio
         │                • Validação de Dados
         ├──────────┬─────────┬─────────┐
         ▼          ▼         ▼         ▼
    ┌────────┐ ┌────────┐ ┌────────┐ ┌────────┐
    │  DB    │ │ Redis  │ │ OMDb   │ │ Email  │
    │  PG15  │ │ Cache  │ │  API   │ │ SMTP   │
    └────────┘ └────────┘ └────────┘ └────────┘
```

#### Camadas do Backend (Clean Architecture)

<details>
<summary><strong>Camada de Domínio</strong></summary>

- **Entidades:** User, Movie, Session, Friendship
- **Interfaces:** Contratos de repositórios
- **Regras de negócio:** Independentes de frameworks

</details>

<details>
<summary><strong>Camada de Casos de Uso</strong></summary>

- **Operações de Usuário:** Registro, Login, Gerenciamento de Perfil
- **Operações de Filme:** Busca, Detalhes, Favoritos
- **Recursos Sociais:** Amigos, Seguindo, Matching
- **Operações de Auth:** Geração de token, Gerenciamento de sessão

</details>

<details>
<summary><strong>Camada de Infraestrutura</strong></summary>

- **Handlers HTTP:** Mapeamento de Request/Response
- **Repositórios:** Implementações de banco de dados
- **Serviços Externos:** API OMDb, Serviço de email
- **Middleware:** Auth, Logging, CORS, Rate limiting

</details>

<details>
<summary><strong>Cadeia de Dados de Filmes</strong></summary>

```
Requisição → API OMDb → Cache DB → 404
              (80-90%)    (10-20%)   (<1%)
```

O sistema usa o padrão **Chain of Responsibility**:
1. Primeiro tenta a API OMDb (dados externos)
2. Fallback para banco de dados (dados em cache)
3. Retorna 404 se não encontrado

**Auto-Save:** Filmes buscados da OMDb são automaticamente salvos no banco, reduzindo chamadas de API futuras em 80-90%.

</details>

</details>

<details>
<summary><h3>💻 Tecnologias</h3></summary>

<a name="tecnologias-pt"></a>

#### Backend
| Tecnologia | Versão | Propósito |
|------------|---------|---------|
| Go | 1.21+ | Linguagem principal |
| Chi Router | v5 | Roteamento HTTP |
| PostgreSQL | 15 | Banco de dados principal |
| SQLx | Latest | Toolkit SQL |
| Redis | 7 | Camada de cache |
| JWT | - | Autenticação |
| Bcrypt | - | Hash de senhas |
| Swagger | 2.0 | Documentação da API |

#### Frontend (Em Desenvolvimento)
| Tecnologia | Versão | Propósito |
|------------|---------|---------|
| Flutter | 3.0+ | Framework mobile |
| Dart | Latest | Linguagem de programação |
| Provider | Latest | Gerenciamento de estado (planejado) |

#### DevOps
| Tecnologia | Propósito |
|------------|---------|
| Docker | Containerização |
| Docker Compose | Orquestração multi-container |
| Git | Controle de versão |
| GitHub | Hospedagem de código |

</details>

<details>
<summary><h3>📁 Estrutura do Projeto</h3></summary>

<a name="estrutura-pt"></a>

```
cine/
├── api_v2/                    # Backend Go (Atual)
│   ├── cmd/
│   │   └── main.go           # Ponto de entrada da aplicação
│   ├── internal/
│   │   ├── config/           # Gerenciamento de configuração
│   │   ├── domain/           # Entidades de negócio
│   │   ├── dto/              # Data Transfer Objects
│   │   ├── handler/http/     # Handlers HTTP
│   │   ├── middleware/       # Middleware HTTP
│   │   ├── repository/       # Camada de acesso a dados
│   │   ├── service/          # Lógica de negócio
│   │   ├── usecase/          # Casos de uso da aplicação
│   │   ├── infrastructure/   # Integrações externas
│   │   └── i18n/             # Internacionalização
│   ├── migrations/           # Migrações do banco de dados
│   ├── docs/                 # Documentação Swagger
│   └── README.md             # Documentação do backend
│
├── flutter_app/              # Frontend Flutter (Em Desenvolvimento)
│   ├── lib/
│   │   └── src/              # Código fonte da aplicação
│   └── README.md
│
├── scripts/                  # Scripts utilitários
│   ├── setup.sh              # Configuração do ambiente
│   └── go-lint.sh            # Linting de código
│
├── docker-compose.yml        # Orquestração de containers
└── README.md                 # Este arquivo
```

</details>

<details>
<summary><h3>✨ Funcionalidades</h3></summary>

<a name="funcionalidades-pt"></a>

<details>
<summary><strong>✅ Funcionalidades Implementadas</strong></summary>

#### Autenticação & Autorização
- ✅ Registro de usuário com validação
- ✅ Login com geração de token JWT
- ✅ Hash de senhas com bcrypt
- ✅ Gerenciamento de sessão
- ✅ Rotas protegidas com middleware
- ✅ Mecanismo de refresh de token

#### Gerenciamento de Usuários
- ✅ Perfis de usuário com nomes de exibição
- ✅ Validação de unicidade de username
- ✅ Validação de email
- ✅ Recuperação de perfil
- ✅ Busca de usuários

#### Integração de Filmes
- ✅ Integração com API OMDb
- ✅ Busca de filmes com paginação
- ✅ Detalhes de filme por ID
- ✅ Padrão Chain of Responsibility
- ✅ Cache em banco de dados (TTL 48h)
- ✅ Auto-save após busca na API
- ✅ Rastreamento de provedor (OMDb, Database)

#### Infraestrutura
- ✅ Banco de dados PostgreSQL com migrações
- ✅ Suporte a cache Redis
- ✅ Containerização Docker
- ✅ Endpoint de health check
- ✅ Configuração CORS
- ✅ Logging de requisições
- ✅ Tratamento de erros
- ✅ Documentação da API Swagger (14 endpoints)

</details>

<details>
<summary><strong>🚧 Em Desenvolvimento</strong></summary>

#### Aplicação Mobile
- 🚧 Estrutura do app Flutter
- 🚧 Design UI/UX
- 🚧 Integração com API
- 🚧 Gerenciamento de estado
- 🚧 Suporte offline

#### Recursos Avançados
- 🚧 Sistema de confirmação de email
- 🚧 Fluxo de reset de senha
- 🚧 Edição de perfil de usuário
- 🚧 Upload de avatar
- 🚧 Recursos sociais (amigos, seguindo)
- 🚧 Listas de filmes (assistidos, favoritos, wishlist)
- 🚧 Avaliações e reviews de filmes
- 🚧 Sistema de notificações
- 🚧 Chat em tempo real

</details>

</details>

<details>
<summary><h3>🚀 Primeiros Passos</h3></summary>

<a name="inicio-pt"></a>

<details>
<summary><strong>Pré-requisitos</strong></summary>

- Go 1.21 ou superior
- PostgreSQL 15
- Redis 7 (opcional mas recomendado)
- Docker & Docker Compose (para setup containerizado)
- Chave da API OMDb (obtenha grátis em [omdbapi.com](https://www.omdbapi.com/apikey.aspx))

</details>

<details>
<summary><strong>Usando Docker (Recomendado)</strong></summary>

1. **Clone o repositório**
```bash
git clone https://github.com/EduardoMG12/cine.git
cd cine
```

2. **Configure o ambiente**
```bash
cp api_v2/.env.example api_v2/.env
# Edite api_v2/.env com sua chave da API OMDb e credenciais do banco
```

3. **Inicie os serviços**
```bash
docker-compose up -d
```

4. **Execute as migrações**
```bash
docker-compose exec api psql -U cineverse -d cineverse -f /app/migrations/001_clean_initial_schema.sql
docker-compose exec api psql -U cineverse -d cineverse -f /app/migrations/002_add_provider_and_sync.sql
```

5. **Acesse a API**
- API: http://localhost:8080
- Swagger UI: http://localhost:8080/swagger/index.html
- Health Check: http://localhost:8080/health

</details>

<details>
<summary><strong>Setup de Desenvolvimento Local</strong></summary>

1. **Instale as dependências Go**
```bash
cd api_v2
go mod download
```

2. **Configure o banco de dados**
```bash
psql -U postgres
CREATE DATABASE cineverse;
CREATE USER cineverse WITH PASSWORD 'sua_senha';
GRANT ALL PRIVILEGES ON DATABASE cineverse TO cineverse;
\q
```

3. **Execute as migrações**
```bash
psql -U cineverse -d cineverse -f migrations/001_clean_initial_schema.sql
psql -U cineverse -d cineverse -f migrations/002_add_provider_and_sync.sql
```

4. **Configure as variáveis de ambiente**
```bash
export DB_HOST=localhost
export DB_PORT=5432
export DB_USER=cineverse
export DB_PASSWORD=sua_senha
export DB_NAME=cineverse
export OMDB_API_KEY=sua_chave_omdb
export JWT_SECRET=seu_jwt_secret
```

5. **Execute a API**
```bash
go run ./cmd/main.go
```

</details>

<details>
<summary><strong>Testando a API</strong></summary>

#### Registrar um novo usuário
```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "johndoe",
    "email": "john@example.com",
    "display_name": "John Doe",
    "password": "SecurePass123!"
  }'
```

#### Fazer login
```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "johndoe",
    "password": "SecurePass123!"
  }'
```

#### Buscar filmes
```bash
curl "http://localhost:8080/api/v1/movies/search?q=matrix&page=1"
```

#### Obter detalhes de um filme
```bash
curl "http://localhost:8080/api/v1/movies/tt0133093"
```

</details>

</details>

<details>
<summary><h3>📊 Status do Desenvolvimento</h3></summary>

<a name="status-pt"></a>

#### Progresso do Backend
```
████████████████░░░░  80% Completo
```
- ✅ Estrutura principal da API
- ✅ Sistema de autenticação
- ✅ Integração de filmes (OMDb)
- ✅ Design do banco de dados
- ✅ Documentação da API
- 🚧 Recursos sociais
- 🚧 Sistema de notificações
- 🚧 Matching avançado

#### Progresso do Frontend
```
███░░░░░░░░░░░░░░░░░  15% Completo
```
- ✅ Estrutura do projeto
- ✅ Navegação básica
- 🚧 Telas de autenticação
- 🚧 Feed principal
- 🚧 Detalhes de filmes
- 🚧 Perfis de usuário
- 🚧 Interações sociais

#### Progresso da Documentação
```
████████████░░░░░░░░  60% Completo
```
- ✅ Documentação da API (Swagger)
- ✅ README do Backend
- ✅ README Raiz
- ✅ Guia de arquitetura
- 🚧 Manual do usuário
- 🚧 Guia de deployment
- 🚧 Guia de contribuição

#### Progresso de Testes
```
████░░░░░░░░░░░░░░░░  20% Completo
```
- ✅ Testes manuais de endpoints
- 🚧 Testes unitários
- 🚧 Testes de integração
- 🚧 Testes E2E
- 🚧 Testes de carga

</details>

<details>
<summary><h3>🗺️ Roadmap</h3></summary>

<a name="roadmap-pt"></a>

<details>
<summary><strong>Sprint 1 - Fundação</strong> ✅ Completo</summary>

- [x] Setup da estrutura do projeto
- [x] Design do schema do banco de dados
- [x] Sistema de autenticação
- [x] Gerenciamento básico de usuários
- [x] Integração OMDb
- [x] Documentação da API

</details>

<details>
<summary><strong>Sprint 2 - Recursos Principais</strong> 🚧 Em Progresso</summary>

- [ ] Sistema de confirmação de email
- [ ] Fluxo de reset de senha
- [ ] Edição de perfil de usuário
- [ ] Listas personalizadas de filmes (assistidos, favoritos, wishlist)
- [ ] Sistema de amigos (adicionar, remover, listar)
- [ ] Sistema básico de notificações

</details>

<details>
<summary><strong>Sprint 3 - Recursos Sociais</strong> 📅 Planejado</summary>

- [ ] Algoritmo de matching de filmes
- [ ] Listas de assistir compartilhadas
- [ ] Recomendações de filmes
- [ ] Sistema de following
- [ ] Feed de atividades
- [ ] Compartilhamento social

</details>

<details>
<summary><strong>Sprint 4 - App Mobile</strong> 📅 Planejado</summary>

- [ ] Autenticação no app Flutter
- [ ] Implementação do feed principal
- [ ] Busca e detalhes de filmes
- [ ] Perfis de usuário
- [ ] Interações sociais
- [ ] Suporte offline

</details>

<details>
<summary><strong>Sprint 5 - Recursos Avançados</strong> 📅 Planejado</summary>

- [ ] Sistema de chat em tempo real
- [ ] Reviews e avaliações de filmes
- [ ] Filtros avançados de busca
- [ ] Recomendações personalizadas (ML)
- [ ] Notificações de watchlist
- [ ] Integração com plataformas de streaming

</details>

<details>
<summary><strong>Sprint 6 - Finalização & Deploy</strong> 📅 Planejado</summary>

- [ ] Testes abrangentes
- [ ] Otimização de performance
- [ ] Auditoria de segurança
- [ ] Deploy em produção
- [ ] Monitoramento e logging
- [ ] Documentação do usuário

</details>

</details>

<details>
<summary><h3>👥 Equipe</h3></summary>

<a name="equipe-pt"></a>

**Líder do Projeto & Desenvolvedor:** Charles Eduardo Mello Guimaraes e Willian Fragata
**Instituição:** IFPR Campus Palmas-PR  
**Curso:** Sistema de informações  
**GitHub:** [@EduardoMG12](https://github.com/EduardoMG12)

**Orientador:** Alexis Kang  
**Curso:** Engenharia de Software - IFPR

</details>

<details>
<summary><h3>📄 Licença</h3></summary>

<a name="licenca-pt"></a>

Este projeto está licenciado sob a Licença MIT - veja o arquivo [LICENSE](LICENSE) para detalhes.

**Uso Acadêmico:** Este projeto faz parte de um trabalho acadêmico do IFPR Campus Palmas-PR. Sinta-se livre para usá-lo como referência para fins de aprendizado.

</details>

<details>
<summary><h3>📚 Recursos Adicionais</h3></summary>

- [Documentação do Backend](api_v2/README.md)
- [Guia de Arquitetura](api_v2/ARCHITECTURE.md)
- [Integração OMDb](api_v2/OMDB_INTEGRATION.md)
- [Documentação Swagger da API](http://localhost:8080/swagger/index.html)
- [Diagrama do Banco de Dados](database_diagram.dbml)

</details>

</details>

---

<div align="center">

**Made with ❤️ at IFPR Campus Palmas-PR**

</div>
