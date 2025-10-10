# RFC-003: Funcionalidades Restantes e Roadmap CineVerse API

**Autor:** CineVerse Development Team  
**Status:** Proposta  
**Data de Criação:** 2025-10-10  
**Versão:** 1.0  
**Baseado em:** RFC-001 (Análise de Gap de Implementação)

---

## 1. Resumo Executivo

Este RFC define as funcionalidades restantes a serem implementadas na API CineVerse v2, priorizadas com base na análise do estado atual da implementação. Com **80% das funcionalidades core já implementadas**, este documento estabelece o roadmap para completar a plataforma social de cinéfilos.

### Status Atual (Atualizado: 2025-10-10)
- ✅ **24 features implementadas e funcionando** (77%)
- 🟡 **2 features com infraestrutura pronta** (7%)
- ❌ **5 features não iniciadas** (16%)
- 🆕 **Sistema de i18n completo implementado**
- 🆕 **Infraestrutura social preparada**

---

## 2. Funcionalidades Pendentes (Por Prioridade)

### 🔥 **PRIORIDADE CRÍTICA - Completar Funcionalidades Básicas**

#### ✅ P0.1 - Sistema de Usuários Finalizado (RF-01.9 & RF-01.10)
**Status:** ✅ CONCLUÍDO  
**Descrição:** Endpoints de sessões e preferências do usuário implementados  
**Tempo Real:** 1 dia (2025-10-10)  

**Tarefas Concluídas:**
- [x] Implementar endpoint `GET /users/me/sessions` (listar sessões ativas)
- [x] Implementar endpoint `DELETE /users/me/sessions/{sessionId}` (revogar sessão específica)
- [x] Implementar endpoint `DELETE /users/me/sessions` (logout de todas as sessões)
- [x] Implementar endpoint `PUT /users/me/settings` (atualizar tema/preferências)
- [x] Documentar no Swagger com anotações completas
- [ ] Adicionar testes para os novos endpoints (pendente)

**Endpoints Necessários:**
```
GET    /users/me/sessions           # Listar sessões ativas
DELETE /users/me/sessions/{id}      # Revogar sessão específica
DELETE /users/me/sessions           # Logout completo
PUT    /users/me/settings           # Atualizar preferências
```

---

### 🚀 **PRIORIDADE ALTA - Funcionalidades Sociais Core**

#### P1.1 - Sistema de Amizade (RF-03.1)
**Status:** 🟡 INFRAESTRUTURA PRONTA  
**Descrição:** Sistema completo de amizade com pedidos e aprovações  
**Tempo Estimado:** 1 semana (handlers HTTP apenas)  
**Dependências:** ✅ Completo

**Tarefas Concluídas:**
- [x] Criar entidade `Friendship` no domínio
- [x] Implementar repository para amizades
- [x] Criar service layer para lógica de amizade
- [x] Migração de banco de dados aplicada
- [x] Sistema de i18n para mensagens multilíngues

**Tarefas Restantes:**
- [ ] Implementar handlers HTTP
- [ ] Implementar DTOs para requests/responses
- [ ] Escrever testes completos
- [ ] Documentar API no Swagger

**Endpoints:**
```
POST   /users/{id}/friend-request    # Enviar pedido de amizade
POST   /users/friend-requests/{id}/accept    # Aceitar pedido
POST   /users/friend-requests/{id}/decline   # Recusar pedido
GET    /users/me/friends            # Listar amigos
GET    /users/me/friend-requests    # Listar pedidos pendentes
DELETE /users/{id}/friendship       # Remover amizade
```

#### P1.2 - Sistema de Seguidores (RF-03.2)
**Status:** 🟡 INFRAESTRUTURA PRONTA  
**Descrição:** Permitir seguir usuários sem amizade mútua  
**Tempo Estimado:** 1 semana (handlers HTTP apenas)  
**Dependências:** ✅ Completo

**Tarefas Concluídas:**
- [x] Criar entidade `Follow` no domínio
- [x] Implementar repository para follows
- [x] Criar service layer com validações
- [x] Controle de privacidade integrado
- [x] Migração de banco de dados aplicada

**Tarefas Restantes:**
- [ ] Implementar handlers HTTP
- [ ] Implementar DTOs
- [ ] Testes automatizados
- [ ] Documentação Swagger

**Endpoints:**
```
POST   /users/{id}/follow           # Seguir usuário
DELETE /users/{id}/follow           # Deixar de seguir
GET    /users/{id}/followers        # Listar seguidores
GET    /users/{id}/following        # Listar seguindo
GET    /users/me/feed               # Feed de atividades
```

#### P1.3 - Sistema de Posts (RF-03.3)
**Status:** Não implementado  
**Descrição:** Posts de usuários com controle de visibilidade  
**Tempo Estimado:** 1-2 semanas  
**Dependências:** Sistemas de amizade e seguidores

**Tarefas:**
- [ ] Criar entidade `Post` no domínio
- [ ] Implementar repository com paginação
- [ ] Service layer com controle de privacidade
- [ ] Handlers HTTP completos
- [ ] Sistema de feed personalizado
- [ ] DTOs para posts e feeds
- [ ] Testes de integração
- [ ] Documentação completa

**Endpoints:**
```
POST   /posts                       # Criar post
GET    /posts/{id}                  # Ver post específico
PUT    /posts/{id}                  # Editar post
DELETE /posts/{id}                  # Deletar post
GET    /users/{id}/posts            # Posts do usuário
GET    /feed                        # Feed personalizado
```

---

### ⭐ **PRIORIDADE MÉDIA - Funcionalidades Avançadas**

#### P2.1 - Sistema de Match de Filmes (RF-04.*)
**Status:** Não implementado  
**Descrição:** Sistema completo de matching de filmes entre usuários  
**Tempo Estimado:** 2-3 semanas  
**Dependências:** Sistema social completo

**Funcionalidades:**
- **RF-04.1:** Sessões de match multi-usuário
- **RF-04.2:** Sugestões baseadas em preferências
- **RF-04.3:** Sistema de like/dislike
- **RF-04.4:** Detecção automática de matches
- **RF-04.5:** Workflow pós-match

**Tarefas:**
- [ ] Modelar entidades: `MatchSession`, `MatchParticipant`, `MatchInteraction`
- [ ] Algoritmo de sugestões de filmes
- [ ] Sistema de votação em tempo real
- [ ] Detecção de matches automática
- [ ] WebSocket para interações real-time
- [ ] Sistema de convites
- [ ] Interface de resultados
- [ ] Testes de cenários complexos

**Endpoints:**
```
POST   /match-sessions/start        # Iniciar sessão de match
POST   /match-sessions/{id}/invite  # Convidar usuários
GET    /match-sessions/{id}/suggestions  # Obter sugestões
POST   /match-sessions/{id}/interact     # Like/dislike filme
GET    /match-sessions/{id}/matches      # Ver matches encontrados
POST   /match-sessions/{id}/finalize     # Finalizar sessão
```

#### P2.2 - Sistema de Notificações (RF-06.1)
**Status:** Não implementado  
**Descrição:** Notificações em tempo real para interações sociais  
**Tempo Estimado:** 1-2 semanas  
**Dependências:** Sistema social básico

**Tarefas:**
- [ ] Infraestrutura de notificações (WebSocket/SSE)
- [ ] Sistema de templates de notificação
- [ ] Persistência de notificações
- [ ] Preferências de notificação do usuário
- [ ] Push notifications (futuro)
- [ ] Sistema de agregação de notificações
- [ ] Marcação como lida/não lida

**Endpoints:**
```
GET    /notifications               # Listar notificações
PUT    /notifications/{id}/read     # Marcar como lida
DELETE /notifications/{id}          # Deletar notificação
PUT    /users/me/notification-settings  # Configurações
WebSocket: /ws/notifications       # Notificações real-time
```

---

### 🎯 **PRIORIDADE BAIXA - Melhorias e Otimizações**

#### P3.1 - Funcionalidades Avançadas de Lista
**Tempo Estimado:** 1 semana

**Tarefas:**
- [ ] Listas colaborativas entre amigos
- [ ] Sistema de tags personalizadas
- [ ] Ordenação avançada de listas
- [ ] Exportação/importação de listas
- [ ] Listas temáticas automáticas

#### P3.2 - Analytics e Insights
**Tempo Estimado:** 2 semanas

**Tarefas:**
- [ ] Estatísticas pessoais de filmes assistidos
- [ ] Insights de preferências do usuário
- [ ] Relatórios de atividade
- [ ] Comparação com amigos
- [ ] Trending topics personalizados

#### P3.3 - Moderação de Conteúdo
**Tempo Estimado:** 1-2 semanas

**Tarefas:**
- [ ] Sistema de reports/denúncias
- [ ] Moderação automática de conteúdo
- [ ] Bloqueio de usuários
- [ ] Filtros de conteúdo
- [ ] Dashboard de moderação

---

## 3. Especificações Técnicas

### 3.1 Novos Modelos de Domínio

```go
// Sistema de Amizade
type Friendship struct {
    UserID1   int       `db:"user_id_1"`
    UserID2   int       `db:"user_id_2"` 
    Status    string    `db:"status"`     // pending, accepted, blocked
    CreatedAt time.Time `db:"created_at"`
    UpdatedAt time.Time `db:"updated_at"`
}

// Sistema de Seguidores
type Follow struct {
    FollowerID  int       `db:"follower_id"`
    FollowingID int       `db:"following_id"`
    CreatedAt   time.Time `db:"created_at"`
}

// Sistema de Posts
type Post struct {
    ID         int       `db:"id"`
    UserID     int       `db:"user_id"`
    Content    string    `db:"content"`
    Visibility string    `db:"visibility"` // public, friends, private
    CreatedAt  time.Time `db:"created_at"`
    UpdatedAt  time.Time `db:"updated_at"`
}

// Sistema de Match
type MatchSession struct {
    ID          int       `db:"id"`
    HostUserID  int       `db:"host_user_id"`
    Status      string    `db:"status"`     // active, completed, cancelled
    CreatedAt   time.Time `db:"created_at"`
    CompletedAt *time.Time `db:"completed_at"`
}

type MatchParticipant struct {
    SessionID int `db:"session_id"`
    UserID    int `db:"user_id"`
    JoinedAt  time.Time `db:"joined_at"`
}

type MatchInteraction struct {
    SessionID int  `db:"session_id"`
    UserID    int  `db:"user_id"`
    MovieID   int  `db:"movie_id"`
    Liked     bool `db:"liked"`
    CreatedAt time.Time `db:"created_at"`
}

// Sistema de Notificações
type Notification struct {
    ID        int       `db:"id"`
    UserID    int       `db:"user_id"`
    Type      string    `db:"type"`       // friend_request, match_found, etc.
    Content   string    `db:"content"`
    Data      JSON      `db:"data"`       // Additional metadata
    Read      bool      `db:"read"`
    CreatedAt time.Time `db:"created_at"`
}
```

### 3.2 Arquitetura de Real-time

**WebSocket Implementation:**
```go
// WebSocket Hub for real-time features
type Hub struct {
    clients    map[int]*Client  // userID -> client
    broadcast  chan []byte
    register   chan *Client
    unregister chan *Client
}

// Message types for real-time communication
type WSMessage struct {
    Type string      `json:"type"`
    Data interface{} `json:"data"`
}

// Types: "notification", "match_update", "friend_request", etc.
```

### 3.3 Migrações de Banco de Dados

**Novas tabelas necessárias:**
```sql
-- Sistema de Amizade
CREATE TABLE friendships (
    user_id_1 INTEGER REFERENCES users(id),
    user_id_2 INTEGER REFERENCES users(id),
    status VARCHAR(20) DEFAULT 'pending',
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    PRIMARY KEY (user_id_1, user_id_2)
);

-- Sistema de Seguidores
CREATE TABLE follows (
    follower_id INTEGER REFERENCES users(id),
    following_id INTEGER REFERENCES users(id),
    created_at TIMESTAMP DEFAULT NOW(),
    PRIMARY KEY (follower_id, following_id)
);

-- Sistema de Posts
CREATE TABLE posts (
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(id),
    content TEXT NOT NULL,
    visibility VARCHAR(20) DEFAULT 'public',
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- Sistema de Match (continua...)
```

---

## 4. Cronograma de Implementação

### **✅ Sprint 1 (Semana 1): Base Finalizada - CONCLUÍDO**
- ✅ Sistema de sessões de usuário completo
- ✅ Configurações de usuário implementadas
- ✅ Documentação Swagger completa
- ✅ Docker Compose simplificado (sem networks customizadas)
- ✅ API funcionando em http://localhost:8080

### **Sprint 2-3 (Semanas 2-3): Sistema Social Básico**
- 🚀 Sistema de amizade completo
- 🚀 Sistema de seguidores
- 🚀 Notificações básicas

### **Sprint 4-5 (Semanas 4-5): Posts e Feed**
- ⭐ Sistema de posts
- ⭐ Feed personalizado
- ⭐ Controles de privacidade

### **Sprint 6-8 (Semanas 6-8): Match de Filmes**
- ⭐ Sistema de matching completo
- ⭐ WebSocket para real-time
- ⭐ Algoritmos de sugestão

### **Sprint 9+ (Semanas 9+): Polimento**
- 🎯 Otimizações de performance
- 🎯 Analytics avançados
- 🎯 Funcionalidades premium

---

## 5. Métricas de Sucesso

### 5.1 Métricas Técnicas
- **Cobertura de Testes:** Manter > 80%
- **Performance API:** Response time < 200ms (95th percentile)
- **Uptime:** > 99.9%
- **Documentação:** 100% endpoints documentados no Swagger

### 5.2 Métricas de Produto
- **Sistema Social:** Taxa de engajamento entre usuários
- **Match de Filmes:** Taxa de sucesso de matches
- **Notificações:** Taxa de abertura e engajamento
- **Adoção:** Crescimento de usuários ativos

---

## 6. Riscos e Mitigações

### 6.1 Riscos Técnicos
- **Complexidade do Real-time:** Mitigação via WebSocket bem testado
- **Performance do Feed:** Mitigação via caching inteligente
- **Escalabilidade Social:** Mitigação via arquitetura assíncrona

### 6.2 Riscos de Produto
- **Adoção das Funcionalidades Sociais:** Mitigação via UX/UI intuitivos
- **Spam/Abuso:** Mitigação via moderação automática
- **Performance com Muitos Usuários:** Mitigação via otimizações de BD

---

## 7. Considerações de Implementação

### 7.1 Padrões a Seguir
- Manter Clean Architecture existente
- Seguir padrões de nomenclatura estabelecidos
- Implementar testes para cada nova funcionalidade
- Documentar completamente no Swagger
- Usar conventional commits

### 7.2 Tecnologias Adicionais
- **WebSocket:** Para funcionalidades real-time
- **Background Jobs:** Redis/Queue para processamento assíncrono
- **Caching Avançado:** Para feeds e notificações
- **Rate Limiting:** Para proteção de APIs sociais

---

## 8. Conclusão

Este RFC estabelece um roadmap claro para completar a implementação da CineVerse API v2. Com as **funcionalidades core já implementadas (71%)**, o foco agora é nas **funcionalidades sociais e de matching** que transformarão a plataforma em uma verdadeira rede social de cinéfilos.

**Próximos Passos Imediatos:**
1. ✅ Endpoints de sessão/preferências (P0.1) - CONCLUÍDO
2. 🚀 Implementar sistema de amizade (P1.1) - PRÓXIMO
3. 🚀 Desenvolver sistema de seguidores (P1.2) - DEPOIS DE P1.1

**Meta:** Ter uma **API social completa e funcional** em **8-10 semanas**, mantendo os altos padrões de qualidade, documentação e testes já estabelecidos no projeto.
