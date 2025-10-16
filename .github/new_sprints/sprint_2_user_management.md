# 👤 Sprint 2: Gestão de Usuários e Perfis

**Objetivo:** Implementar funcionalidades completas de gestão de perfil e configurações de usuário

## 📋 Tarefas Principais

### 1. Serviços de Usuário
- [ ] **Criar service/user_service.go**:
  - `GetProfile(userID uuid.UUID) (*User, error)`
  - `GetPublicProfile(username string, viewerID *uuid.UUID) (*PublicProfile, error)`
  - `UpdateProfile(userID uuid.UUID, req UpdateProfileRequest) (*User, error)`
  - `UpdateSettings(userID uuid.UUID, req UpdateSettingsRequest) error`
  - `SearchUsers(query string, limit int, viewerID *uuid.UUID) ([]UserSearchResult, error)`
  - `CheckUsernameAvailability(username string, excludeUserID *uuid.UUID) (bool, error)`

### 2. Novos DTOs
- [ ] **Criar dto/user_dto.go**:
```go
type UpdateProfileRequest struct {
    DisplayName       *string `json:"display_name" validate:"omitempty,min=2,max=100"`
    Bio               *string `json:"bio" validate:"omitempty,max=500"`
    ProfilePictureURL *string `json:"profile_picture_url" validate:"omitempty,url"`
}

type UpdateSettingsRequest struct {
    Theme     *string `json:"theme" validate:"omitempty,oneof=light dark"`
    IsPrivate *bool   `json:"is_private"`
}

type PublicProfile struct {
    ID                uuid.UUID `json:"id"`
    Username          string    `json:"username"`
    DisplayName       string    `json:"display_name"`
    Bio               *string   `json:"bio,omitempty"`
    ProfilePictureURL *string   `json:"profile_picture_url,omitempty"`
    IsPrivate         bool      `json:"is_private"`
    CreatedAt         time.Time `json:"created_at"`
    // Stats opcionais
    ReviewCount       int  `json:"review_count,omitempty"`
    MovieListCount    int  `json:"movie_list_count,omitempty"`
    FollowerCount     int  `json:"follower_count,omitempty"`
    FollowingCount    int  `json:"following_count,omitempty"`
}

type UserSearchResult struct {
    ID                uuid.UUID `json:"id"`
    Username          string    `json:"username"`
    DisplayName       string    `json:"display_name"`
    ProfilePictureURL *string   `json:"profile_picture_url,omitempty"`
    IsPrivate         bool      `json:"is_private"`
}
```

### 3. Handlers HTTP
- [ ] **Criar handler/user_handler.go**:
  - `GET /users/me` - Perfil do usuário logado
  - `PUT /users/me` - Atualizar perfil
  - `PUT /users/me/settings` - Atualizar configurações
  - `GET /users/{username}` - Perfil público de usuário
  - `GET /users/search?q={query}` - Buscar usuários
  - `GET /users/username-available?username={username}` - Verificar disponibilidade

### 4. Middleware e Validações
- [ ] **Middleware de validação de perfil**:
  - Verificar se usuário tem permissão para ver perfil
  - Respeitar configurações de privacidade
  - Rate limiting para buscas

- [ ] **Validações específicas**:
  - Username: alfanumérico + underscore, 3-30 chars
  - Display name: 2-100 chars
  - Bio: máximo 500 chars
  - Avatar URL: validar URL e formato de imagem

### 5. Sistema de Privacidade
- [ ] **Lógica de visibilidade**:
```go
func (s *UserService) canViewProfile(targetUser *User, viewerID *uuid.UUID) bool {
    // Perfil público = sempre visível
    if !targetUser.IsPrivate {
        return true
    }
    
    // Owner sempre pode ver próprio perfil
    if viewerID != nil && *viewerID == targetUser.ID {
        return true
    }
    
    // Perfil privado = apenas amigos (implementar depois)
    return false
}
```

### 6. Cache de Perfis (Opcional)
- [ ] **Cache Redis** para perfis consultados frequentemente
- [ ] **TTL de 15 minutos** para dados de perfil
- [ ] **Invalidação** ao atualizar perfil

## 🔧 Endpoints da API

### Perfil Próprio
```http
GET /api/v1/users/me
Authorization: Bearer <token>

Response:
{
  "success": true,
  "data": {
    "id": "uuid",
    "username": "cinelover",
    "email": "user@example.com",
    "display_name": "Cine Lover",
    "bio": "Amante de filmes clássicos",
    "profile_picture_url": "https://...",
    "is_private": false,
    "email_verified": false,
    "theme": "dark",
    "created_at": "2025-10-16T10:00:00Z",
    "updated_at": "2025-10-16T10:30:00Z"
  }
}
```

### Atualizar Perfil
```http
PUT /api/v1/users/me
Authorization: Bearer <token>
Content-Type: application/json

{
  "display_name": "Cinema Enthusiast",
  "bio": "Passionate about cinema and storytelling",
  "profile_picture_url": "https://example.com/avatar.jpg"
}
```

### Configurações
```http
PUT /api/v1/users/me/settings
Authorization: Bearer <token>
Content-Type: application/json

{
  "theme": "dark",
  "is_private": true
}
```

### Perfil Público
```http
GET /api/v1/users/cinelover

Response:
{
  "success": true,
  "data": {
    "id": "uuid",
    "username": "cinelover", 
    "display_name": "Cine Lover",
    "bio": "Amante de filmes clássicos",
    "profile_picture_url": "https://...",
    "is_private": false,
    "created_at": "2025-10-16T10:00:00Z",
    "review_count": 15,
    "movie_list_count": 3,
    "follower_count": 42,
    "following_count": 38
  }
}
```

### Busca de Usuários
```http
GET /api/v1/users/search?q=cinema&limit=10
Authorization: Bearer <token>

Response:
{
  "success": true,
  "data": [
    {
      "id": "uuid",
      "username": "cinemaboy",
      "display_name": "Cinema Boy",
      "profile_picture_url": "https://...",
      "is_private": false
    }
  ]
}
```

### Verificar Username
```http
GET /api/v1/users/username-available?username=newuser

Response:
{
  "success": true,
  "data": {
    "username": "newuser",
    "available": true
  }
}
```

## 🎯 Funcionalidades Avançadas

### 1. Upload de Avatar (Futuro)
```go
// Para próxima sprint
POST /api/v1/users/me/avatar
Content-Type: multipart/form-data

// Integração com AWS S3 ou Cloudinary
```

### 2. Estatísticas de Perfil
- Contador de reviews
- Contador de listas de filmes
- Filmes assistidos este mês
- Gêneros favoritos

### 3. Perfil Responsivo
- Diferentes níveis de detalhe baseado no viewer
- Informações públicas vs privadas
- Amigos vs não-amigos

## 🧪 Testes Importantes

### Service Tests
```go
func TestUserService_GetProfile_Success(t *testing.T)
func TestUserService_GetProfile_NotFound(t *testing.T)
func TestUserService_UpdateProfile_Success(t *testing.T)
func TestUserService_UpdateProfile_ValidationError(t *testing.T)
func TestUserService_GetPublicProfile_PrivateUser(t *testing.T)
func TestUserService_SearchUsers_Success(t *testing.T)
```

### Handler Tests
```go
func TestUserHandler_GetMe_Success(t *testing.T)
func TestUserHandler_UpdateProfile_Success(t *testing.T)
func TestUserHandler_UpdateSettings_Success(t *testing.T)
func TestUserHandler_GetPublicProfile_Success(t *testing.T)
func TestUserHandler_SearchUsers_Success(t *testing.T)
```

## 🚨 Casos de Erro

### Perfil
- Perfil não encontrado: `404 - User not found`
- Perfil privado sem permissão: `403 - Private profile`
- Username já existe: `409 - Username already taken`

### Validação
- Display name muito curto: `400 - Display name too short`
- Bio muito longa: `400 - Bio exceeds maximum length`
- URL de avatar inválida: `400 - Invalid avatar URL`

## 🔐 Segurança e Privacidade

### Rate Limiting
- Busca de usuários: 10 req/min
- Atualização de perfil: 5 req/min
- Verificação de username: 20 req/min

### Dados Sensíveis
- Email nunca é exposto em perfis públicos
- Senhas nunca são retornadas
- IP e dados de sessão são privados

### Validação de Entrada
- Sanitização de bio e display_name
- Validação de URLs de avatar
- Prevenção de XSS

## 📊 Métricas

### Performance
- Tempo de resposta < 100ms para perfis
- Busca de usuários < 200ms
- Cache hit ratio > 80% (se implementado)

### Qualidade
- Cobertura de testes > 85%
- Zero vazamentos de dados sensíveis
- Validação completa de inputs

## 🎯 Critérios de Aceitação

- [ ] Usuário pode ver e editar seu perfil
- [ ] Configurações de privacidade funcionam
- [ ] Busca de usuários é rápida e precisa
- [ ] Perfis públicos respeitam configurações
- [ ] Username availability funciona
- [ ] Validações impedem dados inválidos
- [ ] Rate limiting protege contra abuso
- [ ] Testes passam com boa cobertura
- [ ] API documentada no Swagger

## ⏭️ Próxima Sprint

**Sprint 3: Sistema de Filmes e TMDb Integration**
- Integração com TMDb API
- Cache de filmes
- Busca e descoberta
- Informações detalhadas

---

**Tempo Estimado:** 2-3 dias
**Complexidade:** Média
**Prioridade:** ALTA
**Dependências:** Sprint 1 completa
