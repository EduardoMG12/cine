# ⭐ Sprint 4: Sistema de Reviews e Avaliações

**Objetivo:** Implementar sistema completo de reviews de filmes com ratings e estatísticas

## 📋 Tarefas Principais

### 1. Entidade de Review
- [ ] **Criar domain/review.go**:
```go
type Review struct {
    ID        uuid.UUID `db:"id" json:"id"`
    UserID    uuid.UUID `db:"user_id" json:"user_id"`
    MovieID   uuid.UUID `db:"movie_id" json:"movie_id"`
    Rating    int       `db:"rating" json:"rating"`
    Content   *string   `db:"content" json:"content,omitempty"`
    CreatedAt time.Time `db:"created_at" json:"created_at"`
    UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
    
    // Campos relacionais (não no DB)
    User  *User  `json:"user,omitempty"`
    Movie *Movie `json:"movie,omitempty"`
}

type ReviewStats struct {
    AverageRating   float64 `json:"average_rating"`
    TotalReviews    int     `json:"total_reviews"`
    RatingDistribution map[int]int `json:"rating_distribution"`
}
```

### 2. DTOs de Review
- [ ] **Criar dto/review_dto.go**:
```go
type CreateReviewRequest struct {
    MovieID uuid.UUID `json:"movie_id" validate:"required"`
    Rating  int       `json:"rating" validate:"required,min=1,max=10"`
    Content *string   `json:"content" validate:"omitempty,max=2000"`
}

type UpdateReviewRequest struct {
    Rating  *int    `json:"rating" validate:"omitempty,min=1,max=10"`
    Content *string `json:"content" validate:"omitempty,max=2000"`
}

type ReviewResponse struct {
    Review
    UserInfo struct {
        Username          string  `json:"username"`
        DisplayName       string  `json:"display_name"`
        ProfilePictureURL *string `json:"profile_picture_url,omitempty"`
    } `json:"user_info"`
    MovieInfo struct {
        Title     string  `json:"title"`
        PosterURL *string `json:"poster_url,omitempty"`
        Year      *int    `json:"year,omitempty"`
    } `json:"movie_info"`
}

type ReviewListResponse struct {
    Reviews    []ReviewResponse `json:"reviews"`
    Pagination PaginationInfo   `json:"pagination"`
}

type UserReviewStats struct {
    TotalReviews     int     `json:"total_reviews"`
    AverageRating    float64 `json:"average_rating"`
    FavoriteGenres   []GenreStats `json:"favorite_genres"`
    ReviewsThisMonth int     `json:"reviews_this_month"`
    ReviewsThisYear  int     `json:"reviews_this_year"`
}

type GenreStats struct {
    Genre        string  `json:"genre"`
    Count        int     `json:"count"`
    AverageRating float64 `json:"average_rating"`
}
```

### 3. Repositório de Reviews
- [ ] **Criar repository/review_repository.go**:
  - `CreateReview(review *Review) error`
  - `GetReviewByID(id uuid.UUID) (*Review, error)`
  - `GetReviewByUserAndMovie(userID, movieID uuid.UUID) (*Review, error)`
  - `UpdateReview(review *Review) error`
  - `DeleteReview(id uuid.UUID) error`
  - `GetReviewsByUser(userID uuid.UUID, pagination *PaginationParams) ([]Review, error)`
  - `GetReviewsByMovie(movieID uuid.UUID, pagination *PaginationParams) ([]Review, error)`
  - `GetUserReviewStats(userID uuid.UUID) (*UserReviewStats, error)`
  - `GetMovieReviewStats(movieID uuid.UUID) (*ReviewStats, error)`
  - `GetLatestReviews(limit int) ([]Review, error)`

### 4. Serviço de Reviews
- [ ] **Criar service/review_service.go**:
  - `CreateReview(userID uuid.UUID, req CreateReviewRequest) (*ReviewResponse, error)`
  - `UpdateReview(userID uuid.UUID, reviewID uuid.UUID, req UpdateReviewRequest) (*ReviewResponse, error)`
  - `DeleteReview(userID uuid.UUID, reviewID uuid.UUID) error`
  - `GetUserReviews(userID uuid.UUID, viewerID *uuid.UUID, pagination *PaginationParams) (*ReviewListResponse, error)`
  - `GetMovieReviews(movieID uuid.UUID, pagination *PaginationParams) (*ReviewListResponse, error)`
  - `GetReviewByID(id uuid.UUID, viewerID *uuid.UUID) (*ReviewResponse, error)`
  - `GetUserStats(userID uuid.UUID) (*UserReviewStats, error)`
  - `GetLatestReviews(limit int) ([]ReviewResponse, error)`

### 5. Handlers HTTP
- [ ] **Criar handler/review_handler.go**:
  - `POST /reviews` - Criar review
  - `GET /reviews/{id}` - Ver review específico
  - `PUT /reviews/{id}` - Atualizar review próprio
  - `DELETE /reviews/{id}` - Deletar review próprio
  - `GET /users/{id}/reviews` - Reviews de um usuário
  - `GET /movies/{id}/reviews` - Reviews de um filme
  - `GET /users/{id}/stats` - Estatísticas do usuário
  - `GET /reviews/latest` - Reviews mais recentes

## 🔧 Endpoints da API

### Criar Review
```http
POST /api/v1/reviews
Authorization: Bearer <token>
Content-Type: application/json

{
  "movie_id": "uuid-do-filme",
  "rating": 9,
  "content": "Filme incrível! Christopher Nolan mais uma vez nos surpreende com uma narrativa complexa e visualmente deslumbrante."
}

Response:
{
  "success": true,
  "message": "Review created successfully",
  "data": {
    "id": "uuid-review",
    "user_id": "uuid-user",
    "movie_id": "uuid-filme",
    "rating": 9,
    "content": "Filme incrível! Christopher Nolan...",
    "created_at": "2025-10-16T10:00:00Z",
    "updated_at": "2025-10-16T10:00:00Z",
    "user_info": {
      "username": "cinelover",
      "display_name": "Cine Lover",
      "profile_picture_url": "https://..."
    },
    "movie_info": {
      "title": "Interstellar",
      "poster_url": "https://...",
      "year": 2014
    }
  }
}
```

### Reviews de um Filme
```http
GET /api/v1/movies/uuid-do-filme/reviews?page=1&limit=10&sort=newest

Response:
{
  "success": true,
  "data": {
    "reviews": [
      {
        "id": "uuid-review",
        "rating": 9,
        "content": "Excelente filme sobre viagem no tempo...",
        "created_at": "2025-10-16T10:00:00Z",
        "user_info": {
          "username": "moviecritic",
          "display_name": "Movie Critic",
          "profile_picture_url": "https://..."
        }
      }
    ],
    "pagination": {
      "page": 1,
      "limit": 10,
      "total": 156,
      "total_pages": 16,
      "has_next": true,
      "has_prev": false
    }
  }
}
```

### Reviews de um Usuário
```http
GET /api/v1/users/uuid-do-user/reviews?page=1&limit=20

Response:
{
  "success": true,
  "data": {
    "reviews": [
      {
        "id": "uuid-review",
        "rating": 8,
        "content": "Muito bom, recomendo!",
        "created_at": "2025-10-16T09:00:00Z",
        "movie_info": {
          "title": "Inception",
          "poster_url": "https://...",
          "year": 2010
        }
      }
    ],
    "pagination": {
      "page": 1,
      "limit": 20,
      "total": 47,
      "total_pages": 3
    }
  }
}
```

### Estatísticas do Usuário
```http
GET /api/v1/users/me/stats

Response:
{
  "success": true,
  "data": {
    "total_reviews": 47,
    "average_rating": 7.3,
    "reviews_this_month": 5,
    "reviews_this_year": 28,
    "favorite_genres": [
      {
        "genre": "Ficção científica",
        "count": 12,
        "average_rating": 8.1
      },
      {
        "genre": "Drama", 
        "count": 8,
        "average_rating": 7.8
      }
    ]
  }
}
```

### Atualizar Review
```http
PUT /api/v1/reviews/uuid-review
Authorization: Bearer <token>
Content-Type: application/json

{
  "rating": 10,
  "content": "Mudei de opinião, é uma obra-prima!"
}
```

## 🎯 Funcionalidades Avançadas

### 1. Sistema de Rating
- Escala de 1-10 (seguindo RFC)
- Validação de rating único por usuário por filme
- Cálculo de média automático
- Distribuição de ratings

### 2. Filtros e Ordenação
```go
type ReviewFilters struct {
    MinRating *int    `json:"min_rating"`
    MaxRating *int    `json:"max_rating"`
    HasContent *bool  `json:"has_content"` // Apenas reviews com texto
    SortBy    string  `json:"sort_by"`     // newest, oldest, rating_high, rating_low
}
```

### 3. Moderação de Conteúdo (Básico)
```go
func (s *ReviewService) moderateContent(content string) error {
    // Lista básica de palavras proibidas
    // Verificação de spam (muito texto repetido)
    // Validação de encoding (UTF-8)
    return nil
}
```

### 4. Estatísticas Avançadas
```go
type MovieReviewStats struct {
    AverageRating      float64            `json:"average_rating"`
    TotalReviews       int                `json:"total_reviews"`
    RatingDistribution map[string]int     `json:"rating_distribution"`
    ReviewsThisWeek    int                `json:"reviews_this_week"`
    TopReviews         []ReviewResponse   `json:"top_reviews"` // Mais curtidos (futuro)
}
```

### 5. Cache de Estatísticas
- Cache de stats de filme por 1 hora
- Cache de stats de usuário por 30 minutos
- Invalidação ao criar/atualizar review

## 🧪 Testes Importantes

### Service Tests
```go
func TestReviewService_CreateReview_Success(t *testing.T)
func TestReviewService_CreateReview_DuplicateReview(t *testing.T)
func TestReviewService_UpdateReview_Success(t *testing.T)
func TestReviewService_UpdateReview_NotOwner(t *testing.T)
func TestReviewService_DeleteReview_Success(t *testing.T)
func TestReviewService_GetUserStats_Success(t *testing.T)
```

### Repository Tests
```go
func TestReviewRepository_CreateReview_Success(t *testing.T)
func TestReviewRepository_GetReviewByUserAndMovie_Success(t *testing.T)
func TestReviewRepository_GetUserReviewStats_Success(t *testing.T)
func TestReviewRepository_GetMovieReviewStats_Success(t *testing.T)
```

### Handler Tests
```go
func TestReviewHandler_CreateReview_Success(t *testing.T)
func TestReviewHandler_CreateReview_Unauthorized(t *testing.T)
func TestReviewHandler_UpdateReview_NotOwner(t *testing.T)
func TestReviewHandler_GetMovieReviews_Success(t *testing.T)
```

## 🚨 Casos de Erro

### Criação de Review
- Filme não existe: `404 - Movie not found`
- Review já existe: `409 - Review already exists for this movie`
- Rating inválido: `400 - Rating must be between 1 and 10`
- Conteúdo muito longo: `400 - Review content too long`

### Edição/Exclusão
- Review não encontrado: `404 - Review not found`
- Não é o autor: `403 - You can only edit your own reviews`
- Review muito antigo: `403 - Cannot edit review older than 30 days`

### Validação
- Conteúdo inapropriado: `400 - Inappropriate content detected`
- Spam detectado: `429 - Too many similar reviews`

## 🔒 Regras de Negócio

### 1. Unicidade
- Um usuário pode ter apenas 1 review por filme
- Update substitui o review anterior
- Delete é permanente

### 2. Permissões
- Apenas o autor pode editar/deletar
- Reviews são públicos (respeitando privacidade do perfil)
- Usuários privados: apenas amigos veem reviews

### 3. Tempo Limite
- Reviews podem ser editados por até 30 dias
- Depois disso, apenas admin pode editar
- Histórico de edições (futuro)

### 4. Qualidade
- Mínimo de 10 caracteres para reviews com texto
- Máximo de 2000 caracteres
- Moderação básica automática

## 📊 Métricas de Qualidade

### Performance
- Listagem de reviews < 100ms
- Criação de review < 50ms
- Cálculo de stats < 200ms

### Engagement
- Taxa de reviews com conteúdo textual > 40%
- Distribuição de ratings balanceada
- Reviews por usuário ativo > 5

## 🎯 Critérios de Aceitação

- [ ] Usuário pode criar/editar/deletar reviews
- [ ] Sistema impede reviews duplicados
- [ ] Ratings são validados (1-10)
- [ ] Estatísticas são calculadas corretamente
- [ ] Paginação funciona nas listagens
- [ ] Filtros e ordenação funcionam
- [ ] Privacidade é respeitada
- [ ] Performance está dentro dos limites
- [ ] Testes passam com coverage > 85%
- [ ] API documentada no Swagger

## ⏭️ Próxima Sprint

**Sprint 5: Sistema de Listas de Filmes**
- Listas "Quero Assistir" e "Já Assisti"
- Listas personalizadas
- Compartilhamento de listas
- Estatísticas de listas

---

**Tempo Estimado:** 3-4 dias
**Complexidade:** Média
**Prioridade:** ALTA
**Dependências:** Sprint 1, 2 e 3 completas
