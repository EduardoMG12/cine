# 🔐 Sprint 1: Autenticação e Usuários

**Objetivo:** Implementar sistema completo de autenticação com JWT e gestão básica de usuários

## 📋 Tarefas Principais

### 1. Domínio e Entidades
- [ ] **Criar domain/user.go**:
```go
type User struct {
    ID                 uuid.UUID `db:"id" json:"id"`
    Username          string    `db:"username" json:"username"`
    Email             string    `db:"email" json:"email"`
    DisplayName       string    `db:"display_name" json:"display_name"`
    Bio               *string   `db:"bio" json:"bio,omitempty"`
    ProfilePictureURL *string   `db:"profile_picture_url" json:"profile_picture_url,omitempty"`
    PasswordHash      string    `db:"password_hash" json:"-"`
    IsPrivate         bool      `db:"is_private" json:"is_private"`
    EmailVerified     bool      `db:"email_verified" json:"email_verified"`
    Theme             string    `db:"theme" json:"theme"`
    CreatedAt         time.Time `db:"created_at" json:"created_at"`
    UpdatedAt         time.Time `db:"updated_at" json:"updated_at"`
}
```

- [ ] **Criar domain/auth.go**:
```go
type LoginRequest struct {
    Email    string `json:"email" validate:"required,email"`
    Password string `json:"password" validate:"required,min=8"`
}

type RegisterRequest struct {
    Username    string `json:"username" validate:"required,min=3,max=30"`
    Email       string `json:"email" validate:"required,email"`
    DisplayName string `json:"display_name" validate:"required,min=2,max=100"`
    Password    string `json:"password" validate:"required,min=8"`
}

type AuthResponse struct {
    Token string `json:"token"`
    User  User   `json:"user"`
}
```

### 2. Repositório de Usuários
- [ ] **Criar repository/user_repository.go**:
  - `CreateUser(user *User) error`
  - `GetUserByID(id uuid.UUID) (*User, error)`
  - `GetUserByEmail(email string) (*User, error)`
  - `GetUserByUsername(username string) (*User, error)`
  - `UpdateUser(user *User) error`
  - `DeleteUser(id uuid.UUID) error`

### 3. Repositório de Sessões
- [ ] **Criar repository/session_repository.go**:
  - `CreateSession(session *UserSession) error`
  - `GetSessionByToken(token string) (*UserSession, error)`
  - `DeleteSession(token string) error`
  - `DeleteUserSessions(userID uuid.UUID) error`

### 4. Serviços de Autenticação
- [ ] **Criar service/auth_service.go**:
  - `Register(req RegisterRequest) (*AuthResponse, error)`
  - `Login(req LoginRequest) (*AuthResponse, error)`
  - `ValidateToken(token string) (*User, error)`
  - `Logout(token string) error`
  - `LogoutAll(userID uuid.UUID) error`

- [ ] **Criar service/password_service.go**:
  - `HashPassword(password string) (string, error)`
  - `ComparePassword(hash, password string) bool`

- [ ] **Criar service/jwt_service.go**:
  - `GenerateToken(userID uuid.UUID) (string, error)`
  - `ValidateToken(token string) (*Claims, error)`
  - `ParseToken(token string) (*Claims, error)`

### 5. Handlers HTTP
- [ ] **Criar handler/auth_handler.go**:
  - `POST /auth/register` - Registro de usuário
  - `POST /auth/login` - Login do usuário
  - `POST /auth/logout` - Logout (requer auth)
  - `POST /auth/logout-all` - Logout de todas as sessões
  - `GET /auth/me` - Dados do usuário logado

### 6. Middleware de Autenticação
- [ ] **Criar middleware/auth_middleware.go**:
  - Middleware para rotas protegidas
  - Extração de JWT do header Authorization
  - Validação de token
  - Injeção de contexto do usuário

### 7. DTOs e Validação
- [ ] **Criar dto/auth_dto.go**:
  - Estruturas para requests e responses
  - Tags de validação
  - Serialização JSON

### 8. Testes
- [ ] **Testes unitários** para services
- [ ] **Testes de integração** para handlers
- [ ] **Testes de repositório** com mock DB

## 🔧 Endpoints da API

### Autenticação
```http
POST /api/v1/auth/register
Content-Type: application/json

{
  "username": "cinelover",
  "email": "user@example.com", 
  "display_name": "Cine Lover",
  "password": "strongpassword"
}
```

```http
POST /api/v1/auth/login
Content-Type: application/json

{
  "email": "user@example.com",
  "password": "strongpassword"
}
```

```http
GET /api/v1/auth/me
Authorization: Bearer <jwt_token>
```

```http
POST /api/v1/auth/logout
Authorization: Bearer <jwt_token>
```

### Resposta Padrão de Sucesso
```json
{
  "success": true,
  "message": "User registered successfully",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIs...",
    "user": {
      "id": "123e4567-e89b-12d3-a456-426614174000",
      "username": "cinelover",
      "email": "user@example.com",
      "display_name": "Cine Lover",
      "is_private": false,
      "email_verified": false,
      "theme": "light",
      "created_at": "2025-10-16T10:00:00Z"
    }
  }
}
```

## 🛡️ Segurança

### Password Hashing
- Usar `bcrypt` com cost 12
- Nunca retornar password_hash nas APIs

### JWT Configuration
```go
type JWTClaims struct {
    UserID    uuid.UUID `json:"user_id"`
    Username  string    `json:"username"`
    Email     string    `json:"email"`
    jwt.RegisteredClaims
}

// Token expira em 24 horas
expirationTime := time.Now().Add(24 * time.Hour)
```

### Headers de Segurança
- `Authorization: Bearer <token>`
- CORS configurado corretamente
- Rate limiting básico

## 📊 Estrutura do Banco

### Tabelas Utilizadas
- `users` - Dados dos usuários
- `user_sessions` - Sessões ativas (opcional para logout)

### Índices Importantes
```sql
-- Já existem nas migrations
CREATE UNIQUE INDEX idx_users_email ON users(email);
CREATE UNIQUE INDEX idx_users_username ON users(username);
CREATE INDEX idx_user_sessions_token ON user_sessions(token);
```

## 🧪 Testes Principais

### Service Tests
```go
func TestAuthService_Register_Success(t *testing.T)
func TestAuthService_Register_DuplicateEmail(t *testing.T)
func TestAuthService_Login_Success(t *testing.T)
func TestAuthService_Login_InvalidPassword(t *testing.T)
```

### Handler Tests
```go
func TestAuthHandler_Register_Success(t *testing.T)
func TestAuthHandler_Login_ValidationError(t *testing.T)
func TestAuthHandler_Me_Unauthorized(t *testing.T)
```

## 🎯 Critérios de Aceitação

- [ ] Usuário pode se registrar com dados válidos
- [ ] Sistema valida email e username únicos
- [ ] Login retorna JWT válido
- [ ] Middleware protege rotas autenticadas
- [ ] Logout invalida o token
- [ ] Senhas são hasheadas com bcrypt
- [ ] Todas as validações funcionam
- [ ] Testes passam com coverage > 80%
- [ ] API documentada no Swagger
- [ ] Suporte a múltiplos idiomas (i18n)

## 🚨 Casos de Erro

### Registro
- Email já existe: `400 - Email already registered`
- Username já existe: `400 - Username already taken`
- Senha fraca: `400 - Password too weak`
- Dados inválidos: `400 - Validation error`

### Login
- Email não existe: `401 - Invalid credentials`
- Senha incorreta: `401 - Invalid credentials`
- Conta não verificada: `403 - Email not verified` (futuro)

### Autenticação
- Token inválido: `401 - Invalid token`
- Token expirado: `401 - Token expired`
- Token ausente: `401 - Authorization required`

## ⏭️ Próxima Sprint

**Sprint 2: Gestão de Usuários e Perfis**
- Atualização de perfil
- Upload de avatar
- Configurações de privacidade
- Busca de usuários

---

**Tempo Estimado:** 3-4 dias
**Complexidade:** Média-Alta
**Prioridade:** CRÍTICA
**Dependências:** Sprint 0 completa
