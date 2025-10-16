# 📱 Sprint 7: Sistema de Posts e Feed Social

**Objetivo:** Implementar sistema de posts de usuários com controle de visibilidade e feed personalizado

## 📋 Tarefas Principais

### 1. Entidade de Posts
- [ ] **Criar domain/post.go**:
```go
type Post struct {
    ID         uuid.UUID `db:"id" json:"id"`
    UserID     uuid.UUID `db:"user_id" json:"user_id"`
    Content    string    `db:"content" json:"content"`
    Visibility string    `db:"visibility" json:"visibility"`
    CreatedAt  time.Time `db:"created_at" json:"created_at"`
    UpdatedAt  time.Time `db:"updated_at" json:"updated_at"`
    
    // Campos relacionais
    Author *User `json:"author,omitempty"`
}

const (
    PostVisibilityPublic  = "public"
    PostVisibilityFriends = "friends"
    PostVisibilityPrivate = "private"
)

type FeedItem struct {
    Post
    AuthorInfo UserBasicInfo `json:"author_info"`
    CanEdit    bool          `json:"can_edit"`
    CanDelete  bool          `json:"can_delete"`
}

type FeedResponse struct {
    Items      []FeedItem     `json:"items"`
    Pagination PaginationInfo `json:"pagination"`
}
```

### 2. DTOs de Posts
- [ ] **Criar dto/post_dto.go**:
```go
type CreatePostRequest struct {
    Content    string `json:"content" validate:"required,min=1,max=2000"`
    Visibility string `json:"visibility" validate:"required,oneof=public friends private"`
}

type UpdatePostRequest struct {
    Content    *string `json:"content" validate:"omitempty,min=1,max=2000"`
    Visibility *string `json:"visibility" validate:"omitempty,oneof=public friends private"`
}

type PostResponse struct {
    Post
    AuthorInfo UserBasicInfo `json:"author_info"`
    CanEdit    bool          `json:"can_edit"`
    CanDelete  bool          `json:"can_delete"`
    Stats      *PostStats    `json:"stats,omitempty"`
}

type PostStats struct {
    LikesCount    int `json:"likes_count"`
    CommentsCount int `json:"comments_count"`
    SharesCount   int `json:"shares_count"`
}

type FeedFilters struct {
    UserID     *uuid.UUID `json:"user_id"`     // Posts de usuário específico
    Visibility *string    `json:"visibility"`  // Filtrar por visibilidade
    Since      *time.Time `json:"since"`       // Posts desde data específica
    Until      *time.Time `json:"until"`       // Posts até data específica
}
```

### 3. Repositório de Posts
- [ ] **Criar repository/post_repository.go**:
  - `CreatePost(post *Post) error`
  - `GetPostByID(id uuid.UUID) (*Post, error)`
  - `UpdatePost(post *Post) error`
  - `DeletePost(id uuid.UUID) error`
  - `GetUserPosts(userID uuid.UUID, pagination *PaginationParams) ([]Post, error)`
  - `GetFeedPosts(userID uuid.UUID, friendIDs []uuid.UUID, pagination *PaginationParams) ([]Post, error)`
  - `GetPublicPosts(pagination *PaginationParams) ([]Post, error)`
  - `GetPostsCount(userID uuid.UUID, visibility *string) (int, error)`

### 4. Serviço de Posts
- [ ] **Criar service/post_service.go**:
  - `CreatePost(userID uuid.UUID, req CreatePostRequest) (*PostResponse, error)`
  - `GetPost(postID uuid.UUID, viewerID *uuid.UUID) (*PostResponse, error)`
  - `UpdatePost(userID uuid.UUID, postID uuid.UUID, req UpdatePostRequest) (*PostResponse, error)`
  - `DeletePost(userID uuid.UUID, postID uuid.UUID) error`
  - `GetUserPosts(userID uuid.UUID, viewerID *uuid.UUID, pagination *PaginationParams) (*FeedResponse, error)`
  - `GetUserFeed(userID uuid.UUID, pagination *PaginationParams) (*FeedResponse, error)`
  - `GetPublicFeed(viewerID *uuid.UUID, pagination *PaginationParams) (*FeedResponse, error)`

### 5. Serviço de Feed
- [ ] **Criar service/feed_service.go**:
```go
type FeedService struct {
    postRepo       repository.PostRepository
    socialService  SocialService
    userRepo       repository.UserRepository
}

func (s *FeedService) BuildPersonalizedFeed(userID uuid.UUID, pagination *PaginationParams) (*FeedResponse, error) {
    // 1. Obter amigos do usuário
    friends, err := s.socialService.GetUserFriends(userID)
    if err != nil {
        return nil, err
    }
    
    // 2. Obter usuários que o user segue
    following, err := s.socialService.GetUserFollowing(userID)
    if err != nil {
        return nil, err
    }
    
    // 3. Combinar IDs e buscar posts
    feedUserIDs := append(friends, following...)
    feedUserIDs = append(feedUserIDs, userID) // Incluir posts próprios
    
    posts, err := s.postRepo.GetFeedPosts(userID, feedUserIDs, pagination)
    if err != nil {
        return nil, err
    }
    
    // 4. Converter para feed items com permissões
    return s.buildFeedResponse(posts, userID)
}
```

### 6. Handlers HTTP
- [ ] **Criar handler/post_handler.go**:
  - `POST /posts` - Criar post
  - `GET /posts/{id}` - Ver post específico
  - `PUT /posts/{id}` - Editar post próprio
  - `DELETE /posts/{id}` - Deletar post próprio
  - `GET /posts/feed` - Feed personalizado do usuário
  - `GET /posts/public` - Feed público (descoberta)
  - `GET /users/{id}/posts` - Posts de um usuário específico

## 🔧 Endpoints da API

### Criar Post
```http
POST /api/v1/posts
Authorization: Bearer <token>
Content-Type: application/json

{
  "content": "Acabei de assistir Interstellar pela 5ª vez e ainda me emociono com aquela cena final! 🚀✨ #Nolan #SciFi",
  "visibility": "public"
}

Response:
{
  "success": true,
  "message": "Post created successfully",
  "data": {
    "id": "uuid-post",
    "user_id": "uuid-user",
    "content": "Acabei de assistir Interstellar...",
    "visibility": "public",
    "created_at": "2025-10-16T15:00:00Z",
    "updated_at": "2025-10-16T15:00:00Z",
    "author_info": {
      "id": "uuid-user",
      "username": "cinelover",
      "display_name": "Cine Lover",
      "profile_picture_url": "https://..."
    },
    "can_edit": true,
    "can_delete": true
  }
}
```

### Feed Personalizado
```http
GET /api/v1/posts/feed?page=1&limit=10
Authorization: Bearer <token>

Response:
{
  "success": true,
  "data": {
    "items": [
      {
        "id": "uuid-post-1",
        "user_id": "uuid-friend",
        "content": "Alguém mais acha que o final de Inception é confuso? 🤔",
        "visibility": "public",
        "created_at": "2025-10-16T14:30:00Z",
        "author_info": {
          "id": "uuid-friend",
          "username": "moviecritic",
          "display_name": "Movie Critic",
          "profile_picture_url": "https://..."
        },
        "can_edit": false,
        "can_delete": false
      },
      {
        "id": "uuid-post-2",
        "user_id": "uuid-user",
        "content": "Minha lista de filmes para assistir chegou a 100! 📽️",
        "visibility": "friends",
        "created_at": "2025-10-16T13:15:00Z",
        "author_info": {
          "id": "uuid-user",
          "username": "cinelover",
          "display_name": "Cine Lover",
          "profile_picture_url": "https://..."
        },
        "can_edit": true,
        "can_delete": true
      }
    ],
    "pagination": {
      "page": 1,
      "limit": 10,
      "total": 45,
      "total_pages": 5,
      "has_next": true,
      "has_prev": false
    }
  }
}
```

### Posts de um Usuário
```http
GET /api/v1/users/uuid-usuario/posts?page=1&limit=20&visibility=public
Authorization: Bearer <token>

Response:
{
  "success": true,
  "data": {
    "items": [
      {
        "id": "uuid-post",
        "content": "Que filme incrível! Recomendo a todos.",
        "visibility": "public",
        "created_at": "2025-10-16T12:00:00Z",
        "author_info": {
          "username": "moviefan",
          "display_name": "Movie Fan"
        },
        "can_edit": false,
        "can_delete": false
      }
    ],
    "pagination": {
      "page": 1,
      "limit": 20,
      "total": 12,
      "total_pages": 1
    }
  }
}
```

### Editar Post
```http
PUT /api/v1/posts/uuid-post
Authorization: Bearer <token>
Content-Type: application/json

{
  "content": "Acabei de assistir Interstellar pela 5ª vez e ainda me emociono! Atualização: A trilha sonora do Hans Zimmer é perfeita! 🚀✨🎵",
  "visibility": "public"
}
```

### Feed Público (Descoberta)
```http
GET /api/v1/posts/public?page=1&limit=15

Response:
{
  "success": true,
  "data": {
    "items": [
      {
        "id": "uuid-post",
        "content": "Top 5 filmes de 2024 que não podem faltar na sua lista! 🎬",
        "visibility": "public",
        "created_at": "2025-10-16T16:00:00Z",
        "author_info": {
          "username": "topmovies",
          "display_name": "Top Movies",
          "profile_picture_url": "https://..."
        },
        "can_edit": false,
        "can_delete": false
      }
    ]
  }
}
```

## 🎯 Lógica de Visibilidade

### 1. Controle de Acesso
```go
func (s *PostService) CanViewPost(post *Post, viewerID *uuid.UUID) bool {
    // Post público = todos podem ver
    if post.Visibility == PostVisibilityPublic {
        return true
    }
    
    // Autor sempre pode ver próprios posts
    if viewerID != nil && *viewerID == post.UserID {
        return true
    }
    
    // Post privado = apenas o autor
    if post.Visibility == PostVisibilityPrivate {
        return false
    }
    
    // Post para amigos = verificar amizade
    if post.Visibility == PostVisibilityFriends && viewerID != nil {
        return s.socialService.AreFriends(*viewerID, post.UserID)
    }
    
    return false
}
```

### 2. Algoritmo de Feed
```go
func (s *FeedService) BuildFeed(userID uuid.UUID) (*FeedResponse, error) {
    var allPosts []Post
    
    // 1. Posts próprios (todas as visibilidades)
    myPosts, _ := s.postRepo.GetUserPosts(userID, nil)
    allPosts = append(allPosts, myPosts...)
    
    // 2. Posts públicos de quem segue
    following, _ := s.socialService.GetUserFollowing(userID)
    for _, followedID := range following {
        posts, _ := s.postRepo.GetUserPosts(followedID, &PostFilters{
            Visibility: &PostVisibilityPublic,
        })
        allPosts = append(allPosts, posts...)
    }
    
    // 3. Posts de amigos (públicos + friends)
    friends, _ := s.socialService.GetUserFriends(userID)
    for _, friendID := range friends {
        posts, _ := s.postRepo.GetUserPosts(friendID, &PostFilters{
            Visibility: []string{PostVisibilityPublic, PostVisibilityFriends},
        })
        allPosts = append(allPosts, posts...)
    }
    
    // 4. Ordenar por data (mais recentes primeiro)
    sort.Slice(allPosts, func(i, j int) bool {
        return allPosts[i].CreatedAt.After(allPosts[j].CreatedAt)
    })
    
    return s.convertToFeedItems(allPosts, userID), nil
}
```

### 3. Sistema de Permissões
```go
func (s *PostService) GetPostPermissions(post *Post, viewerID *uuid.UUID) (bool, bool) {
    canEdit := viewerID != nil && *viewerID == post.UserID
    canDelete := canEdit // Mesma regra por enquanto
    
    // Futuramente: moderadores podem deletar posts inapropriados
    
    return canEdit, canDelete
}
```

## 🧪 Testes Importantes

### Service Tests
```go
func TestPostService_CreatePost_Success(t *testing.T)
func TestPostService_CreatePost_InvalidVisibility(t *testing.T)
func TestPostService_UpdatePost_Success(t *testing.T)
func TestPostService_UpdatePost_NotOwner(t *testing.T)
func TestPostService_GetPost_PublicVisible(t *testing.T)
func TestPostService_GetPost_PrivateNotVisible(t *testing.T)
func TestPostService_GetPost_FriendsOnlyVisible(t *testing.T)
```

### Feed Tests
```go
func TestFeedService_PersonalizedFeed_IncludesOwnPosts(t *testing.T)
func TestFeedService_PersonalizedFeed_IncludesFriendsPosts(t *testing.T)
func TestFeedService_PersonalizedFeed_ExcludesPrivatePosts(t *testing.T)
func TestFeedService_PublicFeed_OnlyPublicPosts(t *testing.T)
```

### Repository Tests
```go
func TestPostRepository_CreatePost_Success(t *testing.T)
func TestPostRepository_GetFeedPosts_Success(t *testing.T)
func TestPostRepository_GetUserPosts_FilterByVisibility(t *testing.T)
```

## 🚨 Casos de Erro

### Criação de Posts
- Conteúdo vazio: `400 - Post content cannot be empty`
- Conteúdo muito longo: `400 - Post content too long (max 2000 characters)`
- Visibilidade inválida: `400 - Invalid visibility setting`

### Edição/Exclusão
- Post não encontrado: `404 - Post not found`
- Não é o autor: `403 - You can only edit your own posts`
- Post muito antigo: `403 - Cannot edit post older than 24 hours`

### Visualização
- Post privado: `403 - This post is private`
- Post apenas para amigos: `403 - This post is visible to friends only`

## 🔒 Regras de Negócio

### 1. Criação de Posts
- Conteúdo: mínimo 1, máximo 2000 caracteres
- Visibilidade obrigatória
- Todos os usuários podem criar posts

### 2. Edição de Posts
- Apenas o autor pode editar
- Limite de 24 horas após criação
- Histórico de edições (futuro)

### 3. Exclusão de Posts
- Apenas o autor pode deletar
- Exclusão é permanente
- Moderadores podem deletar (futuro)

### 4. Feed
- Ordenação cronológica reversa
- Filtra por permissões de visualização
- Paginação obrigatória
- Cache de 5 minutos (futuro)

## 📊 Métricas de Engagement

### Criação de Conteúdo
- Posts por usuário ativo > 2/semana
- Taxa de posts públicos > 60%
- Comprimento médio de posts > 50 caracteres

### Feed
- Tempo médio de carregamento < 200ms
- Taxa de atualização do feed > 3x/dia
- Engagement com posts do feed > 30%

## 🎯 Funcionalidades Futuras (Próximas Sprints)

### 1. Interações com Posts
- Sistema de curtidas (likes)
- Comentários em posts
- Compartilhamentos
- Reações diversas

### 2. Mídia em Posts
- Upload de imagens
- GIFs e vídeos
- Links com preview
- Hashtags e menções

### 3. Feed Inteligente
- Algoritmo baseado em relevância
- Posts promovidos
- Filtros avançados
- Sugestões personalizadas

## 🎯 Critérios de Aceitação

- [ ] Usuário pode criar posts com diferentes visibilidades
- [ ] Posts são filtrados corretamente no feed
- [ ] Sistema de permissões funciona
- [ ] Feed personalizado mostra conteúdo relevante
- [ ] Edição/exclusão respeitam regras de negócio
- [ ] Performance está dentro dos limites
- [ ] Validações impedem conteúdo inválido
- [ ] Paginação funciona corretamente
- [ ] Testes passam com coverage > 85%
- [ ] API documentada no Swagger

## ⏭️ Próxima Sprint

**Sprint 8: Sistema de Match de Filmes**
- Sessões colaborativas de escolha de filmes
- Algoritmo de sugestões baseado em preferências
- Sistema de voting (like/dislike)
- WebSocket para tempo real

---

**Tempo Estimado:** 3-4 dias
**Complexidade:** Média-Alta
**Prioridade:** MÉDIA
**Dependências:** Sprint 1, 2 e 6 completas
