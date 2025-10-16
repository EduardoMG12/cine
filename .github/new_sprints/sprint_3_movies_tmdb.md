# 🎬 Sprint 3: Sistema de Filmes e Integração TMDb

**Objetivo:** Implementar sistema completo de filmes com integração à TMDb API e cache inteligente

## 📋 Tarefas Principais

### 1. Entidades de Domínio
- [ ] **Criar domain/movie.go**:
```go
type Movie struct {
    ID              uuid.UUID  `db:"id" json:"id"`
    ExternalAPIID   string     `db:"external_api_id" json:"external_api_id"`
    Title           string     `db:"title" json:"title"`
    Overview        *string    `db:"overview" json:"overview,omitempty"`
    ReleaseDate     *time.Time `db:"release_date" json:"release_date,omitempty"`
    PosterURL       *string    `db:"poster_url" json:"poster_url,omitempty"`
    BackdropURL     *string    `db:"backdrop_url" json:"backdrop_url,omitempty"`
    Genres          []string   `db:"genres" json:"genres"`
    Runtime         *int       `db:"runtime" json:"runtime,omitempty"`
    VoteAverage     *float64   `db:"vote_average" json:"vote_average,omitempty"`
    VoteCount       *int       `db:"vote_count" json:"vote_count,omitempty"`
    Adult           bool       `db:"adult" json:"adult"`
    CacheExpiresAt  time.Time  `db:"cache_expires_at" json:"-"`
    CreatedAt       time.Time  `db:"created_at" json:"created_at"`
    UpdatedAt       time.Time  `db:"updated_at" json:"updated_at"`
}
```

### 2. DTOs para TMDb
- [ ] **Criar dto/movie_dto.go**:
```go
// TMDb API Response structures
type TMDbMovieResponse struct {
    ID               int     `json:"id"`
    Title           string  `json:"title"`
    Overview        string  `json:"overview"`
    ReleaseDate     string  `json:"release_date"`
    PosterPath      string  `json:"poster_path"`
    BackdropPath    string  `json:"backdrop_path"`
    Genres          []TMDbGenre `json:"genres"`
    Runtime         int     `json:"runtime"`
    VoteAverage     float64 `json:"vote_average"`
    VoteCount       int     `json:"vote_count"`
    Adult           bool    `json:"adult"`
}

type TMDbGenre struct {
    ID   int    `json:"id"`
    Name string `json:"name"`
}

type TMDbSearchResponse struct {
    Page         int                 `json:"page"`
    Results      []TMDbMovieResponse `json:"results"`
    TotalPages   int                 `json:"total_pages"`
    TotalResults int                 `json:"total_results"`
}

// API Request/Response DTOs
type MovieSearchRequest struct {
    Query    string `json:"query" validate:"required,min=1"`
    Page     int    `json:"page" validate:"min=1,max=500"`
    Language string `json:"language" validate:"omitempty,oneof=en pt es"`
}

type MovieDetailsResponse struct {
    Movie
    // Campos adicionais se necessário
    Credits *MovieCredits `json:"credits,omitempty"`
}

type MovieCredits struct {
    Cast []CastMember `json:"cast"`
    Crew []CrewMember `json:"crew"`
}
```

### 3. Repositório de Filmes
- [ ] **Criar repository/movie_repository.go**:
  - `CreateMovie(movie *Movie) error`
  - `GetMovieByID(id uuid.UUID) (*Movie, error)`
  - `GetMovieByExternalID(externalID string) (*Movie, error)`
  - `UpdateMovie(movie *Movie) error`
  - `SearchMovies(query string, limit int) ([]Movie, error)`
  - `GetExpiredMovies(limit int) ([]Movie, error)`
  - `GetPopularMovies(limit int) ([]Movie, error)`

### 4. Cliente TMDb
- [ ] **Criar service/tmdb_client.go**:
```go
type TMDbClient struct {
    apiKey     string
    baseURL    string
    httpClient *http.Client
}

func (c *TMDbClient) SearchMovies(query string, page int, language string) (*TMDbSearchResponse, error)
func (c *TMDbClient) GetMovieDetails(movieID string, language string) (*TMDbMovieResponse, error)  
func (c *TMDbClient) GetPopularMovies(page int, language string) (*TMDbSearchResponse, error)
func (c *TMDbClient) GetTrendingMovies(timeWindow string, language string) (*TMDbSearchResponse, error)
```

### 5. Serviço de Filmes
- [ ] **Criar service/movie_service.go**:
  - `SearchMovies(req MovieSearchRequest) (*MovieSearchResponse, error)`
  - `GetMovieDetails(id uuid.UUID) (*MovieDetailsResponse, error)`
  - `GetPopularMovies(language string, limit int) ([]Movie, error)`
  - `GetTrendingMovies(language string, limit int) ([]Movie, error)`
  - `RefreshMovieCache(externalID string) (*Movie, error)`

### 6. Sistema de Cache Inteligente
- [ ] **Estratégia de cache**:
```go
// 1. Consulta no banco local primeiro
movie, err := movieRepo.GetMovieByExternalID(externalID)
if err != nil || movie.CacheExpiresAt.Before(time.Now()) {
    // 2. Buscar na TMDb API
    tmdbMovie, err := tmdbClient.GetMovieDetails(externalID, language)
    if err != nil {
        return nil, err
    }
    
    // 3. Converter e salvar no banco
    movie = convertTMDbToMovie(tmdbMovie)
    movie.CacheExpiresAt = time.Now().Add(24 * time.Hour) // TTL 24h
    
    if existingMovie != nil {
        movieRepo.UpdateMovie(movie)
    } else {
        movieRepo.CreateMovie(movie)
    }
}

return movie, nil
```

### 7. Handlers HTTP
- [ ] **Criar handler/movie_handler.go**:
  - `GET /movies/search?q={query}&page={page}` - Buscar filmes
  - `GET /movies/{id}` - Detalhes do filme
  - `GET /movies/popular?limit={limit}` - Filmes populares
  - `GET /movies/trending?limit={limit}` - Em alta
  - `POST /movies/{id}/refresh` - Forçar atualização do cache

## 🔧 Endpoints da API

### Busca de Filmes
```http
GET /api/v1/movies/search?q=interstellar&page=1&language=pt
Accept-Language: pt-BR

Response:
{
  "success": true,
  "data": {
    "movies": [
      {
        "id": "uuid",
        "external_api_id": "157336",
        "title": "Interstellar",
        "overview": "As reservas naturais da Terra se esgotaram...",
        "release_date": "2014-11-07T00:00:00Z",
        "poster_url": "https://image.tmdb.org/t/p/w500/gEU2QniE6E77NI6lCU6MxlNBvIx.jpg",
        "backdrop_url": "https://image.tmdb.org/t/p/w1280/...",
        "genres": ["Drama", "Ficção científica"],
        "runtime": 169,
        "vote_average": 8.4,
        "vote_count": 35847,
        "adult": false,
        "created_at": "2025-10-16T10:00:00Z",
        "updated_at": "2025-10-16T10:00:00Z"
      }
    ],
    "pagination": {
      "page": 1,
      "total_pages": 1,
      "total_results": 1,
      "has_next": false,
      "has_prev": false
    }
  }
}
```

### Detalhes do Filme
```http
GET /api/v1/movies/uuid-do-filme

Response:
{
  "success": true,
  "data": {
    "id": "uuid",
    "external_api_id": "157336",
    "title": "Interstellar",
    "overview": "As reservas naturais da Terra se esgotaram...",
    "release_date": "2014-11-07T00:00:00Z",
    "poster_url": "https://image.tmdb.org/t/p/w500/gEU2QniE6E77NI6lCU6MxlNBvIx.jpg",
    "genres": ["Drama", "Ficção científica"],
    "runtime": 169,
    "vote_average": 8.4,
    "vote_count": 35847,
    "adult": false,
    "credits": {
      "cast": [
        {
          "name": "Matthew McConaughey",
          "character": "Joseph Cooper",
          "profile_path": "/path/to/photo.jpg"
        }
      ],
      "crew": [
        {
          "name": "Christopher Nolan",
          "job": "Director",
          "profile_path": "/path/to/photo.jpg"
        }
      ]
    }
  }
}
```

### Filmes Populares
```http
GET /api/v1/movies/popular?limit=20&language=pt

Response:
{
  "success": true,
  "data": [
    {
      "id": "uuid",
      "title": "Top Gun: Maverick",
      "poster_url": "https://...",
      "vote_average": 8.3,
      "release_date": "2022-05-27T00:00:00Z"
    }
  ]
}
```

### Filmes em Alta
```http
GET /api/v1/movies/trending?limit=20&time_window=week

Response:
{
  "success": true,
  "data": [
    // Array de filmes similar ao popular
  ]
}
```

## ⚙️ Configuração TMDb

### Environment Variables
```bash
TMDB_API_KEY=your_tmdb_api_key_here
TMDB_BASE_URL=https://api.themoviedb.org/3
TMDB_IMAGE_BASE_URL=https://image.tmdb.org/t/p/
TMDB_CACHE_TTL=24h
TMDB_RATE_LIMIT=40 # requests per second
```

### Configuração no config.go
```go
type TMDbConfig struct {
    APIKey      string `mapstructure:"api_key"`
    BaseURL     string `mapstructure:"base_url"`
    ImageBaseURL string `mapstructure:"image_base_url"`
    CacheTTL    time.Duration `mapstructure:"cache_ttl"`
    RateLimit   int    `mapstructure:"rate_limit"`
}
```

## 🚀 Funcionalidades Avançadas

### 1. Rate Limiting para TMDb
```go
type RateLimiter struct {
    limiter *rate.Limiter
}

func (rl *RateLimiter) Wait(ctx context.Context) error {
    return rl.limiter.Wait(ctx)
}
```

### 2. Cache Redis (Opcional)
- Cache de buscas frequentes por 1 hora
- Cache de filmes populares por 6 horas
- Invalidação inteligente

### 3. Fallback Strategy
```go
// Se TMDb API falhar, retornar dados do cache local (mesmo expirado)
// Log do erro para monitoramento
// Resposta com warning header
```

### 4. Transformação de Dados
```go
func convertTMDbToMovie(tmdb *TMDbMovieResponse) *Movie {
    genres := make([]string, len(tmdb.Genres))
    for i, g := range tmdb.Genres {
        genres[i] = g.Name
    }
    
    return &Movie{
        ExternalAPIID: strconv.Itoa(tmdb.ID),
        Title:        tmdb.Title,
        Overview:     &tmdb.Overview,
        Genres:       genres,
        // ... outros campos
        CacheExpiresAt: time.Now().Add(24 * time.Hour),
    }
}
```

## 🧪 Testes Importantes

### Service Tests
```go
func TestMovieService_SearchMovies_Success(t *testing.T)
func TestMovieService_SearchMovies_CacheHit(t *testing.T)
func TestMovieService_SearchMovies_TMDbAPIError(t *testing.T)
func TestMovieService_GetMovieDetails_Success(t *testing.T)
func TestMovieService_RefreshCache_Success(t *testing.T)
```

### Integration Tests
```go
func TestTMDbClient_SearchMovies_Success(t *testing.T)
func TestTMDbClient_GetMovieDetails_Success(t *testing.T)
func TestTMDbClient_RateLimit_Respected(t *testing.T)
```

### Repository Tests
```go
func TestMovieRepository_CreateMovie_Success(t *testing.T)
func TestMovieRepository_GetExpiredMovies_Success(t *testing.T)
```

## 🚨 Casos de Erro

### TMDb API
- API key inválida: `500 - External API configuration error`
- Rate limit excedido: `429 - Too many requests`
- Filme não encontrado: `404 - Movie not found`
- Timeout da API: `503 - External service unavailable`

### Cache
- Cache expirado + API indisponível: Return cached data with warning
- Dados corrompidos: Refresh from API

### Validação
- Query de busca vazia: `400 - Search query required`
- Página inválida: `400 - Invalid page number`

## 📊 Métricas de Performance

### Cache Efficiency
- Cache hit ratio > 70%
- Tempo médio de resposta com cache < 50ms
- Tempo médio de resposta sem cache < 300ms

### TMDb API
- Rate limiting respeitado (40 req/s)
- Timeout configurado (5 segundos)
- Retry strategy para falhas temporárias

## 🎯 Critérios de Aceitação

- [ ] Busca de filmes funciona com TMDb API
- [ ] Cache local funciona corretamente
- [ ] Detalhes de filmes são completos
- [ ] Filmes populares/trending são atualizados
- [ ] Rate limiting protege a API key
- [ ] Fallback funciona em caso de falha
- [ ] Performance está dentro dos limites
- [ ] Testes passam com coverage > 80%
- [ ] API documentada no Swagger
- [ ] Suporte a múltiplos idiomas

## ⏭️ Próxima Sprint

**Sprint 4: Sistema de Reviews e Avaliações**
- Reviews de filmes
- Sistema de ratings
- Estatísticas pessoais
- Feed de reviews

---

**Tempo Estimado:** 3-4 dias
**Complexidade:** Média-Alta
**Prioridade:** ALTA
**Dependências:** Sprint 1 e 2 completas
