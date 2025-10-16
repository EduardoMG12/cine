# 👥 Sprint 6: Sistema Social - Amizades e Seguidores

**Objetivo:** Implementar funcionalidades sociais básicas (amizades bidirecionais e seguidores unidirecionais)

## 📋 Tarefas Principais

### 1. Entidades Sociais
- [ ] **Criar domain/friendship.go**:
```go
type Friendship struct {
    UserID1   uuid.UUID `db:"user_id_1" json:"user_id_1"`
    UserID2   uuid.UUID `db:"user_id_2" json:"user_id_2"`
    Status    string    `db:"status" json:"status"`
    CreatedAt time.Time `db:"created_at" json:"created_at"`
    UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}

const (
    FriendshipStatusPending  = "pending"
    FriendshipStatusAccepted = "accepted"
    FriendshipStatusDeclined = "declined"
    FriendshipStatusBlocked  = "blocked"
)

type Follow struct {
    FollowerID  uuid.UUID `db:"follower_id" json:"follower_id"`
    FollowingID uuid.UUID `db:"following_id" json:"following_id"`
    CreatedAt   time.Time `db:"created_at" json:"created_at"`
}

type SocialStats struct {
    FriendsCount    int `json:"friends_count"`
    FollowersCount  int `json:"followers_count"`
    FollowingCount  int `json:"following_count"`
    PendingRequests int `json:"pending_requests"`
}
```

### 2. DTOs Sociais
- [ ] **Criar dto/social_dto.go**:
```go
type FriendshipRequest struct {
    UserID uuid.UUID `json:"user_id" validate:"required"`
}

type FriendshipResponse struct {
    UserID1       uuid.UUID `json:"user_id_1"`
    UserID2       uuid.UUID `json:"user_id_2"`
    Status        string    `json:"status"`
    CreatedAt     time.Time `json:"created_at"`
    OtherUserInfo UserBasicInfo `json:"other_user_info"`
}

type FollowResponse struct {
    FollowerID  uuid.UUID     `json:"follower_id"`
    FollowingID uuid.UUID     `json:"following_id"`
    CreatedAt   time.Time     `json:"created_at"`
    UserInfo    UserBasicInfo `json:"user_info"`
}

type UserBasicInfo struct {
    ID                uuid.UUID `json:"id"`
    Username          string    `json:"username"`
    DisplayName       string    `json:"display_name"`
    ProfilePictureURL *string   `json:"profile_picture_url,omitempty"`
    IsPrivate         bool      `json:"is_private"`
}

type SocialConnectionStatus struct {
    IsFriend        bool   `json:"is_friend"`
    FriendshipStatus *string `json:"friendship_status,omitempty"` // pending, accepted
    IsFollowing     bool   `json:"is_following"`
    IsFollower      bool   `json:"is_follower"`
    CanSendRequest  bool   `json:"can_send_request"`
    CanFollow       bool   `json:"can_follow"`
}

type FriendRequestsResponse struct {
    Sent     []FriendshipResponse `json:"sent"`
    Received []FriendshipResponse `json:"received"`
}
```

### 3. Repositórios Sociais
- [ ] **Criar repository/friendship_repository.go**:
  - `CreateFriendship(friendship *Friendship) error`
  - `GetFriendship(userID1, userID2 uuid.UUID) (*Friendship, error)`
  - `UpdateFriendshipStatus(userID1, userID2 uuid.UUID, status string) error`
  - `DeleteFriendship(userID1, userID2 uuid.UUID) error`
  - `GetUserFriends(userID uuid.UUID) ([]uuid.UUID, error)`
  - `GetFriendRequests(userID uuid.UUID) (*FriendRequestsResponse, error)`
  - `GetFriendsCount(userID uuid.UUID) (int, error)`

- [ ] **Criar repository/follow_repository.go**:
  - `CreateFollow(follow *Follow) error`
  - `DeleteFollow(followerID, followingID uuid.UUID) error`
  - `GetFollow(followerID, followingID uuid.UUID) (*Follow, error)`
  - `GetFollowers(userID uuid.UUID, pagination *PaginationParams) ([]uuid.UUID, error)`
  - `GetFollowing(userID uuid.UUID, pagination *PaginationParams) ([]uuid.UUID, error)`
  - `GetFollowersCount(userID uuid.UUID) (int, error)`
  - `GetFollowingCount(userID uuid.UUID) (int, error)`
  - `IsFollowing(followerID, followingID uuid.UUID) (bool, error)`

### 4. Serviços Sociais
- [ ] **Criar service/social_service.go**:
  - `SendFriendRequest(fromUserID, toUserID uuid.UUID) error`
  - `AcceptFriendRequest(userID, requesterID uuid.UUID) error`
  - `DeclineFriendRequest(userID, requesterID uuid.UUID) error`
  - `RemoveFriend(userID, friendID uuid.UUID) error`
  - `BlockUser(userID, blockedUserID uuid.UUID) error`
  - `FollowUser(followerID, followingID uuid.UUID) error`
  - `UnfollowUser(followerID, followingID uuid.UUID) error`
  - `GetSocialStats(userID uuid.UUID) (*SocialStats, error)`
  - `GetConnectionStatus(viewerID, targetID uuid.UUID) (*SocialConnectionStatus, error)`
  - `GetFriendRequests(userID uuid.UUID) (*FriendRequestsResponse, error)`

### 5. Handlers HTTP
- [ ] **Criar handler/social_handler.go**:
  - `POST /users/{id}/friend-request` - Enviar pedido de amizade
  - `POST /friends/requests/{id}/accept` - Aceitar pedido
  - `POST /friends/requests/{id}/decline` - Recusar pedido
  - `DELETE /friends/{id}` - Remover amizade
  - `POST /users/{id}/block` - Bloquear usuário
  - `GET /friends` - Listar amigos
  - `GET /friends/requests` - Pedidos de amizade
  - `POST /users/{id}/follow` - Seguir usuário
  - `DELETE /users/{id}/follow` - Deixar de seguir
  - `GET /users/{id}/followers` - Seguidores do usuário
  - `GET /users/{id}/following` - Quem o usuário segue
  - `GET /users/{id}/social-status` - Status da conexão social

## 🔧 Endpoints da API

### Enviar Pedido de Amizade
```http
POST /api/v1/users/uuid-usuario/friend-request
Authorization: Bearer <token>

Response:
{
  "success": true,
  "message": "Friend request sent successfully",
  "data": {
    "status": "pending",
    "sent_at": "2025-10-16T12:00:00Z"
  }
}
```

### Listar Pedidos de Amizade
```http
GET /api/v1/friends/requests
Authorization: Bearer <token>

Response:
{
  "success": true,
  "data": {
    "received": [
      {
        "user_id_1": "uuid-remetente",
        "user_id_2": "uuid-meu",
        "status": "pending",
        "created_at": "2025-10-16T11:30:00Z",
        "other_user_info": {
          "id": "uuid-remetente",
          "username": "moviefan123",
          "display_name": "Movie Fan",
          "profile_picture_url": "https://...",
          "is_private": false
        }
      }
    ],
    "sent": [
      {
        "user_id_1": "uuid-meu",
        "user_id_2": "uuid-destinatario", 
        "status": "pending",
        "created_at": "2025-10-16T10:00:00Z",
        "other_user_info": {
          "id": "uuid-destinatario",
          "username": "cinephile",
          "display_name": "Cinephile",
          "profile_picture_url": null,
          "is_private": true
        }
      }
    ]
  }
}
```

### Aceitar Pedido de Amizade
```http
POST /api/v1/friends/requests/uuid-remetente/accept
Authorization: Bearer <token>

Response:
{
  "success": true,
  "message": "Friend request accepted",
  "data": {
    "status": "accepted",
    "accepted_at": "2025-10-16T12:30:00Z"
  }
}
```

### Seguir Usuário
```http
POST /api/v1/users/uuid-usuario/follow
Authorization: Bearer <token>

Response:
{
  "success": true,
  "message": "User followed successfully",
  "data": {
    "following_id": "uuid-usuario",
    "followed_at": "2025-10-16T13:00:00Z"
  }
}
```

### Status da Conexão Social
```http
GET /api/v1/users/uuid-usuario/social-status
Authorization: Bearer <token>

Response:
{
  "success": true,
  "data": {
    "is_friend": false,
    "friendship_status": "pending",
    "is_following": true,
    "is_follower": false,
    "can_send_request": false,
    "can_follow": false
  }
}
```

### Listar Amigos
```http
GET /api/v1/friends?page=1&limit=20

Response:
{
  "success": true,
  "data": {
    "friends": [
      {
        "id": "uuid-amigo",
        "username": "bestfriend",
        "display_name": "Best Friend",
        "profile_picture_url": "https://...",
        "friendship_since": "2025-09-15T10:00:00Z"
      }
    ],
    "pagination": {
      "page": 1,
      "limit": 20,
      "total": 15,
      "total_pages": 1
    }
  }
}
```

### Seguidores do Usuário
```http
GET /api/v1/users/uuid-usuario/followers?page=1&limit=20

Response:
{
  "success": true,
  "data": {
    "followers": [
      {
        "id": "uuid-seguidor",
        "username": "follower1",
        "display_name": "Follower One",
        "profile_picture_url": "https://...",
        "following_since": "2025-10-10T14:00:00Z"
      }
    ],
    "pagination": {
      "page": 1,
      "limit": 20,
      "total": 42,
      "total_pages": 3
    }
  }
}
```

## 🎯 Lógica de Negócio Complexa

### 1. Sistema de Amizade Bidirecional
```go
func (s *SocialService) SendFriendRequest(fromUserID, toUserID uuid.UUID) error {
    // 1. Verificar se não são a mesma pessoa
    if fromUserID == toUserID {
        return errors.New("cannot send friend request to yourself")
    }
    
    // 2. Verificar se já não são amigos
    existing, _ := s.friendshipRepo.GetFriendship(fromUserID, toUserID)
    if existing != nil && existing.Status == "accepted" {
        return errors.New("users are already friends")
    }
    
    // 3. Verificar se não há pedido pendente
    if existing != nil && existing.Status == "pending" {
        return errors.New("friend request already sent")
    }
    
    // 4. Verificar se não está bloqueado
    if existing != nil && existing.Status == "blocked" {
        return errors.New("cannot send friend request to blocked user")
    }
    
    // 5. Verificar privacidade (se usuário é privado, precisa ser seguidor primeiro)
    targetUser, err := s.userRepo.GetUserByID(toUserID)
    if err != nil {
        return err
    }
    
    if targetUser.IsPrivate {
        isFollowing, _ := s.followRepo.IsFollowing(fromUserID, toUserID)
        if !isFollowing {
            return errors.New("must follow user first to send friend request")
        }
    }
    
    // 6. Criar amizade com status pending
    return s.friendshipRepo.CreateFriendship(&Friendship{
        UserID1: fromUserID,
        UserID2: toUserID,
        Status:  "pending",
    })
}
```

### 2. Sistema de Seguir Unidirecional
```go
func (s *SocialService) FollowUser(followerID, followingID uuid.UUID) error {
    // 1. Verificar se não é a mesma pessoa
    if followerID == followingID {
        return errors.New("cannot follow yourself")
    }
    
    // 2. Verificar se já não está seguindo
    existing, _ := s.followRepo.GetFollow(followerID, followingID)
    if existing != nil {
        return errors.New("already following this user")
    }
    
    // 3. Verificar se não está bloqueado
    blocked := s.IsBlocked(followerID, followingID)
    if blocked {
        return errors.New("cannot follow blocked user")
    }
    
    // 4. Criar follow
    return s.followRepo.CreateFollow(&Follow{
        FollowerID:  followerID,
        FollowingID: followingID,
    })
}
```

### 3. Controle de Privacidade
```go
func (s *SocialService) CanViewUserContent(viewerID *uuid.UUID, targetUserID uuid.UUID) bool {
    // Usuário vendo próprio perfil
    if viewerID != nil && *viewerID == targetUserID {
        return true
    }
    
    targetUser, _ := s.userRepo.GetUserByID(targetUserID)
    if targetUser == nil {
        return false
    }
    
    // Perfil público
    if !targetUser.IsPrivate {
        return true
    }
    
    // Perfil privado - verificar se é amigo
    if viewerID != nil {
        friendship, _ := s.friendshipRepo.GetFriendship(*viewerID, targetUserID)
        if friendship != nil && friendship.Status == "accepted" {
            return true
        }
    }
    
    return false
}
```

## 🧪 Testes Críticos

### Casos de Teste de Amizade
```go
func TestSocialService_SendFriendRequest_Success(t *testing.T)
func TestSocialService_SendFriendRequest_ToSelf_Error(t *testing.T)
func TestSocialService_SendFriendRequest_AlreadyFriends_Error(t *testing.T)
func TestSocialService_SendFriendRequest_AlreadyPending_Error(t *testing.T)
func TestSocialService_AcceptFriendRequest_Success(t *testing.T)
func TestSocialService_AcceptFriendRequest_NotPending_Error(t *testing.T)
```

### Casos de Teste de Follow
```go
func TestSocialService_FollowUser_Success(t *testing.T)
func TestSocialService_FollowUser_Self_Error(t *testing.T)
func TestSocialService_FollowUser_AlreadyFollowing_Error(t *testing.T)
func TestSocialService_UnfollowUser_Success(t *testing.T)
```

### Casos de Privacidade
```go
func TestSocialService_CanViewUserContent_PublicProfile_Success(t *testing.T)
func TestSocialService_CanViewUserContent_PrivateProfile_NotFriend_Denied(t *testing.T)
func TestSocialService_CanViewUserContent_PrivateProfile_Friend_Allowed(t *testing.T)
```

## 🚨 Casos de Erro

### Amizades
- Enviar para si mesmo: `400 - Cannot send friend request to yourself`
- Já são amigos: `409 - Users are already friends`
- Pedido já existe: `409 - Friend request already pending`
- Usuário bloqueado: `403 - Cannot send request to blocked user`
- Usuário não existe: `404 - User not found`

### Seguidores
- Seguir a si mesmo: `400 - Cannot follow yourself`
- Já seguindo: `409 - Already following this user`
- Usuário bloqueado: `403 - Cannot follow blocked user`

### Permissões
- Aceitar pedido não direcionado: `403 - Not authorized to accept this request`
- Remover amigo inexistente: `404 - Friendship not found`

## 🔒 Regras de Privacidade

### 1. Perfis Privados
- Apenas amigos podem ver conteúdo completo
- Seguidores podem ver informações básicas
- Para enviar pedido de amizade em perfil privado, precisa seguir primeiro

### 2. Bloqueio
- Usuário bloqueado não pode:
  - Enviar pedidos de amizade
  - Seguir o usuário
  - Ver perfil ou conteúdo
- Bloqueio é unidirecional
- Remove amizade e follow existentes

### 3. Visibilidade de Listas Sociais
- Amigos: sempre visível
- Seguidores/Seguindo: visível se perfil público
- Perfil privado: apenas amigos veem listas sociais

## 📊 Métricas Sociais

### Engagement
- Taxa de aceitação de pedidos > 70%
- Média de amigos por usuário ativo > 5
- Taxa follow-back entre usuários > 40%

### Performance
- Listagem de amigos < 100ms
- Envio de pedido < 50ms
- Verificação de status < 30ms

## 🎯 Critérios de Aceitação

- [ ] Sistema de amizade bidirecional funciona
- [ ] Sistema de seguir unidirecional funciona
- [ ] Pedidos de amizade são gerenciados corretamente
- [ ] Controles de privacidade funcionam
- [ ] Bloqueio impede interações
- [ ] Status de conexão é preciso
- [ ] Performance está dentro dos limites
- [ ] Validações impedem casos inválidos
- [ ] Testes passam com coverage > 85%
- [ ] API documentada no Swagger

## ⏭️ Próxima Sprint

**Sprint 7: Sistema de Posts e Feed Social**
- Posts com controle de visibilidade
- Feed personalizado baseado em amigos/seguidores
- Interações básicas (futuro: curtidas/comentários)

---

**Tempo Estimado:** 4-5 dias
**Complexidade:** Alta
**Prioridade:** MÉDIA
**Dependências:** Sprint 1 e 2 completas
