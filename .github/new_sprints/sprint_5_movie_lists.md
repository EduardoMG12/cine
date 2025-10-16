# 📋 Sprint 5: Sistema de Listas de Filmes

**Objetivo:** Implementar sistema completo de listas de filmes (Quero Assistir, Já Assisti, Listas Personalizadas)

## 📋 Tarefas Principais

### 1. Entidades de Domínio
- [ ] **Criar domain/movie_list.go**:
```go
type MovieList struct {
    ID          uuid.UUID `db:"id" json:"id"`
    UserID      uuid.UUID `db:"user_id" json:"user_id"`
    Name        string    `db:"name" json:"name"`
    IsDefault   bool      `db:"is_default" json:"is_default"`
    CreatedAt   time.Time `db:"created_at" json:"created_at"`
    UpdatedAt   time.Time `db:"updated_at" json:"updated_at"`
    
    // Campos relacionais
    User    *User             `json:"user,omitempty"`
    Entries []MovieListEntry `json:"entries,omitempty"`
    Stats   *ListStats       `json:"stats,omitempty"`
}

type MovieListEntry struct {
    ID           uuid.UUID `db:"id" json:"id"`
    MovieListID  uuid.UUID `db:"movie_list_id" json:"movie_list_id"`
    MovieID      uuid.UUID `db:"movie_id" json:"movie_id"`
    AddedAt      time.Time `db:"added_at" json:"added_at"`
    
    // Campos relacionais
    Movie *Movie `json:"movie,omitempty"`
}

type ListStats struct {
    TotalMovies    int                `json:"total_movies"`
    GenreStats     []GenreCount       `json:"genre_stats"`
    AverageRating  float64            `json:"average_rating"`
    TotalRuntime   int                `json:"total_runtime"` // em minutos
    YearStats      map[string]int     `json:"year_stats"`
}

type GenreCount struct {
    Genre string `json:"genre"`
    Count int    `json:"count"`
}
```

### 2. DTOs de Listas
- [ ] **Criar dto/movie_list_dto.go**:
```go
type CreateListRequest struct {
    Name string `json:"name" validate:"required,min=1,max=100"`
}

type UpdateListRequest struct {
    Name *string `json:"name" validate:"omitempty,min=1,max=100"`
}

type AddMovieToListRequest struct {
    MovieID uuid.UUID `json:"movie_id" validate:"required"`
}

type MovieListResponse struct {
    MovieList
    MovieCount int `json:"movie_count"`
    UserInfo   struct {
        Username          string  `json:"username"`
        DisplayName       string  `json:"display_name"`
        ProfilePictureURL *string `json:"profile_picture_url,omitempty"`
    } `json:"user_info,omitempty"`
}

type MovieListDetailResponse struct {
    MovieListResponse
    Movies []MovieWithDetails `json:"movies"`
}

type MovieWithDetails struct {
    MovieListEntry
    MovieInfo struct {
        Title       string    `json:"title"`
        PosterURL   *string   `json:"poster_url,omitempty"`
        ReleaseDate *string   `json:"release_date,omitempty"`
        Genres      []string  `json:"genres"`
        Runtime     *int      `json:"runtime,omitempty"`
    } `json:"movie_info"`
    UserReview *struct {
        Rating  int     `json:"rating"`
        Content *string `json:"content,omitempty"`
    } `json:"user_review,omitempty"`
}

type DefaultListsResponse struct {
    WatchList   *MovieListResponse `json:"watch_list"`
    WatchedList *MovieListResponse `json:"watched_list"`
}
```

### 3. Repositório de Listas
- [ ] **Criar repository/movie_list_repository.go**:
  - `CreateList(list *MovieList) error`
  - `GetListByID(id uuid.UUID) (*MovieList, error)`
  - `GetUserLists(userID uuid.UUID) ([]MovieList, error)`
  - `GetDefaultLists(userID uuid.UUID) (*DefaultListsResponse, error)`
  - `UpdateList(list *MovieList) error`
  - `DeleteList(id uuid.UUID) error`
  - `AddMovieToList(listID, movieID uuid.UUID) error`
  - `RemoveMovieFromList(listID, movieID uuid.UUID) error`
  - `GetListMovies(listID uuid.UUID, pagination *PaginationParams) ([]MovieListEntry, error)`
  - `GetListStats(listID uuid.UUID) (*ListStats, error)`
  - `IsMovieInList(listID, movieID uuid.UUID) (bool, error)`

### 4. Serviço de Listas
- [ ] **Criar service/movie_list_service.go**:
  - `CreateList(userID uuid.UUID, req CreateListRequest) (*MovieListResponse, error)`
  - `GetUserLists(userID uuid.UUID, viewerID *uuid.UUID) ([]MovieListResponse, error)`
  - `GetListDetail(listID uuid.UUID, viewerID *uuid.UUID) (*MovieListDetailResponse, error)`
  - `UpdateList(userID uuid.UUID, listID uuid.UUID, req UpdateListRequest) (*MovieListResponse, error)`
  - `DeleteList(userID uuid.UUID, listID uuid.UUID) error`
  - `AddMovieToWatchList(userID uuid.UUID, movieID uuid.UUID) error`
  - `AddMovieToWatchedList(userID uuid.UUID, movieID uuid.UUID) error`
  - `MoveToWatched(userID uuid.UUID, movieID uuid.UUID) error`
  - `RemoveMovieFromList(userID uuid.UUID, listID uuid.UUID, movieID uuid.UUID) error`
  - `GetDefaultLists(userID uuid.UUID) (*DefaultListsResponse, error)`
  - `CreateDefaultLists(userID uuid.UUID) error`

### 5. Handlers HTTP
- [ ] **Criar handler/movie_list_handler.go**:
  - `GET /lists` - Listar listas do usuário logado
  - `POST /lists` - Criar nova lista personalizada
  - `GET /lists/{id}` - Detalhes de uma lista
  - `PUT /lists/{id}` - Atualizar lista
  - `DELETE /lists/{id}` - Deletar lista
  - `POST /lists/{id}/movies` - Adicionar filme à lista
  - `DELETE /lists/{id}/movies/{movieId}` - Remover filme da lista
  - `GET /users/{id}/lists` - Listas públicas de um usuário
  - `POST /movies/{id}/want-to-watch` - Adicionar à lista "Quero Assistir"
  - `POST /movies/{id}/watched` - Adicionar à lista "Já Assisti"
  - `PUT /movies/{id}/move-to-watched` - Mover de "Quero" para "Assistido"

## 🔧 Endpoints da API

### Listas do Usuário
```http
GET /api/v1/lists
Authorization: Bearer <token>

Response:
{
  "success": true,
  "data": [
    {
      "id": "uuid-lista",
      "name": "Quero Assistir",
      "is_default": true,
      "created_at": "2025-10-16T10:00:00Z",
      "movie_count": 15
    },
    {
      "id": "uuid-lista-2", 
      "name": "Já Assisti",
      "is_default": true,
      "created_at": "2025-10-16T10:00:00Z",
      "movie_count": 47
    },
    {
      "id": "uuid-lista-3",
      "name": "Melhores de Terror",
      "is_default": false,
      "created_at": "2025-10-16T11:00:00Z",
      "movie_count": 8
    }
  ]
}
```

### Detalhes de uma Lista
```http
GET /api/v1/lists/uuid-lista?page=1&limit=20

Response:
{
  "success": true,
  "data": {
    "id": "uuid-lista",
    "name": "Quero Assistir",
    "is_default": true,
    "created_at": "2025-10-16T10:00:00Z",
    "movie_count": 15,
    "stats": {
      "total_movies": 15,
      "genre_stats": [
        {
          "genre": "Ficção científica",
          "count": 6
        },
        {
          "genre": "Drama", 
          "count": 4
        }
      ],
      "total_runtime": 1847,
      "year_stats": {
        "2023": 5,
        "2022": 3,
        "2021": 7
      }
    },
    "movies": [
      {
        "id": "uuid-entry",
        "movie_id": "uuid-filme",
        "added_at": "2025-10-15T14:30:00Z",
        "movie_info": {
          "title": "Dune: Part Two",
          "poster_url": "https://...",
          "release_date": "2024-03-01",
          "genres": ["Ficção científica", "Aventura"],
          "runtime": 166
        }
      }
    ]
  }
}
```

### Criar Lista Personalizada
```http
POST /api/v1/lists
Authorization: Bearer <token>
Content-Type: application/json

{
  "name": "Filmes para Maratona"
}

Response:
{
  "success": true,
  "message": "List created successfully",
  "data": {
    "id": "uuid-nova-lista",
    "name": "Filmes para Maratona",
    "is_default": false,
    "created_at": "2025-10-16T12:00:00Z",
    "movie_count": 0
  }
}
```

### Adicionar Filme à Lista "Quero Assistir"
```http
POST /api/v1/movies/uuid-filme/want-to-watch
Authorization: Bearer <token>

Response:
{
  "success": true,
  "message": "Movie added to Want to Watch list"
}
```

### Mover Filme para "Já Assisti"
```http
PUT /api/v1/movies/uuid-filme/move-to-watched
Authorization: Bearer <token>

Response:
{
  "success": true,
  "message": "Movie moved to Watched list"
}
```

### Adicionar Filme à Lista Personalizada
```http
POST /api/v1/lists/uuid-lista/movies
Authorization: Bearer <token>
Content-Type: application/json

{
  "movie_id": "uuid-filme"
}
```

## 🎯 Funcionalidades Especiais

### 1. Listas Padrão (Default Lists)
```go
// Criadas automaticamente no primeiro login
const (
    DefaultWatchListName   = "Quero Assistir"
    DefaultWatchedListName = "Já Assisti"
)

func (s *MovieListService) ensureDefaultLists(userID uuid.UUID) error {
    // Verifica se já existem, se não, cria
}
```

### 2. Sistema de Movimentação
- Adicionar filme diretamente a "Quero Assistir"
- Mover de "Quero Assistir" para "Já Assisti"
- Adicionar direto a "Já Assisti" (para filmes já vistos)
- Remove automaticamente de uma lista ao adicionar na outra

### 3. Estatísticas Avançadas
```go
func (r *MovieListRepository) GetListStats(listID uuid.UUID) (*ListStats, error) {
    // Calcular:
    // - Total de filmes
    // - Distribuição por gênero  
    // - Tempo total (runtime)
    // - Anos de lançamento
    // - Rating médio (se tiver reviews)
}
```

### 4. Controle de Privacidade
- Listas seguem configuração de privacidade do usuário
- Usuários privados: apenas amigos veem listas
- Listas podem ser individualmente privadas (futuro)

### 5. Validações de Negócio
```go
// Um filme não pode estar na mesma lista mais de uma vez
// Listas padrão não podem ser deletadas
// Limite de 50 listas personalizadas por usuário
// Limite de 1000 filmes por lista
```

## 🧪 Testes Importantes

### Service Tests
```go
func TestMovieListService_CreateList_Success(t *testing.T)
func TestMovieListService_AddMovieToWatchList_Success(t *testing.T)
func TestMovieListService_MoveToWatched_Success(t *testing.T)
func TestMovieListService_DeleteDefaultList_Error(t *testing.T)
func TestMovieListService_AddDuplicateMovie_Error(t *testing.T)
func TestMovieListService_GetListStats_Success(t *testing.T)
```

### Repository Tests
```go
func TestMovieListRepository_CreateList_Success(t *testing.T)
func TestMovieListRepository_AddMovieToList_Success(t *testing.T)
func TestMovieListRepository_GetListMovies_WithPagination(t *testing.T)
func TestMovieListRepository_IsMovieInList_Success(t *testing.T)
```

### Handler Tests
```go
func TestMovieListHandler_GetUserLists_Success(t *testing.T)
func TestMovieListHandler_CreateList_Success(t *testing.T)
func TestMovieListHandler_AddToWatchList_Success(t *testing.T)
func TestMovieListHandler_MoveToWatched_Success(t *testing.T)
```

## 🚨 Casos de Erro

### Criação/Edição
- Nome muito curto/longo: `400 - Invalid list name`
- Limite de listas: `400 - Maximum lists limit reached`
- Nome duplicado: `409 - List name already exists`

### Adição de Filmes
- Filme já na lista: `409 - Movie already in list`
- Lista não encontrada: `404 - List not found`
- Filme não existe: `404 - Movie not found`
- Limite de filmes: `400 - List is full`

### Exclusão
- Tentar deletar lista padrão: `403 - Cannot delete default list`
- Lista não é do usuário: `403 - Not list owner`

## 🔒 Regras de Negócio

### 1. Listas Padrão
- São criadas automaticamente no registro
- Não podem ser deletadas
- Podem ser renomeadas
- São sempre privadas conforme configuração do usuário

### 2. Listas Personalizadas  
- Máximo 50 por usuário
- Nome único por usuário
- Podem ser deletadas (com confirmação)
- Máximo 1000 filmes por lista

### 3. Movimentação de Filmes
- Filme pode estar em múltiplas listas personalizadas
- Mas apenas em uma das listas padrão (Quero OU Já Assisti)
- Mover para "Já Assisti" remove de "Quero Assistir"

### 4. Privacidade
- Seguem configuração geral do usuário
- Usuários privados = apenas amigos veem
- Usuários públicos = todos veem

## 📊 Métricas de Engagement

### Performance
- Listagem de listas < 50ms
- Detalhes de lista < 100ms  
- Adição de filme < 30ms

### Usuário
- Média de listas por usuário ativo > 3
- Filmes por lista "Quero Assistir" > 10
- Taxa de conversão "Quero" → "Assistido" > 30%

## 🎯 Critérios de Aceitação

- [ ] Listas padrão são criadas automaticamente
- [ ] Usuário pode criar listas personalizadas
- [ ] Sistema de movimentação funciona corretamente
- [ ] Estatísticas são calculadas precisamente
- [ ] Paginação funciona nas listagens
- [ ] Duplicatas são prevenidas
- [ ] Privacidade é respeitada
- [ ] Performance está dentro dos limites
- [ ] Validações impedem uso indevido
- [ ] Testes passam com coverage > 85%
- [ ] API documentada no Swagger

## ⏭️ Próxima Sprint

**Sprint 6: Sistema Social - Amizades e Seguidores**
- Envio de pedidos de amizade
- Sistema de seguir usuários
- Feed de atividades
- Notificações básicas

---

**Tempo Estimado:** 3-4 dias
**Complexidade:** Média
**Prioridade:** ALTA  
**Dependências:** Sprint 1, 2 e 3 completas
