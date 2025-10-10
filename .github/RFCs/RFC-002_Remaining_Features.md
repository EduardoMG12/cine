# RFC-002: Funcionalidades Pendentes da API CineVerse

**Autor:** AI Assistant (baseado na análise do RFC-001)
**Status:** Proposta
**Data de Criação:** 2025-10-10
**Versão:** 1.0

---

## 1. Resumo (Abstract)

Este documento especifica as funcionalidades ainda não implementadas da API CineVerse (RFC-001), priorizando-as por impacto e complexidade. Baseado na análise do código atual, identifica o que falta para completar a especificação original.

## 2. Status Atual da Implementação

### ✅ **FUNCIONALIDADES COMPLETAMENTE IMPLEMENTADAS:**

#### RF-01: Autenticação e Usuários (80% completo)
- ✅ Registro e login com JWT
- ✅ Gestão de perfis e configurações
- ✅ Controle de sessões ativas
- ✅ Sistema de privacidade

#### RF-02: Gerenciamento de Filmes (100% completo)
- ✅ Listas "Quero Assistir" e "Já Assisti"
- ✅ Sistema de reviews e ratings
- ✅ Listas personalizadas
- ✅ Integração com TMDb
- ✅ Sistema de cache otimizado

#### RF-05: Descoberta de Filmes (100% completo)
- ✅ Busca por título, gênero, popularidade
- ✅ Cache inteligente com TTL
- ✅ Integração externa funcionando

---

## 3. Funcionalidades Pendentes (Por Prioridade)

### 🔴 **PRIORIDADE ALTA - Funcionalidades Core Ausentes**

#### P1.1 - Sistema de E-mail (RF-01.6 & RF-01.7)
**Descrição:** Implementar confirmação de e-mail e recuperação de senha
**Complexidade:** Média
**Justificativa:** Funcionalidade de segurança essencial

**Tarefas:**
- [ ] Configurar serviço de e-mail (SMTP/SendGrid)
- [ ] Implementar envio de e-mail de confirmação no registro
- [ ] Criar sistema de tokens para confirmação/reset de senha
- [ ] Implementar fluxo "esqueci minha senha"
- [ ] Adicionar templates de e-mail profissionais
- [ ] Testes automatizados para fluxos de e-mail

**Endpoints:**
```
POST /auth/confirm-email     # Confirmar e-mail com token
POST /auth/forgot-password   # Solicitar reset de senha
POST /auth/reset-password    # Resetar senha com token
POST /auth/resend-confirmation # Reenviar e-mail de confirmação
```

#### P1.2 - Sistema Social Básico (RF-03.1 & RF-03.2)
**Descrição:** Implementar amizades e sistema de seguir usuários
**Complexidade:** Alta
**Justificativa:** Base para funcionalidades sociais e match de filmes

**Tarefas:**
- [ ] Implementar repositórios de Friendship e Follow
- [ ] Criar serviços para operações sociais
- [ ] Desenvolver handlers HTTP para amizade/seguir
- [ ] Sistema de notificações básico
- [ ] Configurar rotas protegidas
- [ ] Testes de integração social

**Endpoints:**
```
# Amizades
POST /users/{id}/friend-request    # Enviar pedido de amizade
POST /users/friend-requests/{id}/accept  # Aceitar pedido
POST /users/friend-requests/{id}/decline # Recusar pedido
DELETE /users/friends/{id}         # Remover amigo
GET /users/me/friends              # Listar amigos
GET /users/me/friend-requests      # Pedidos pendentes

# Seguir
POST /users/{id}/follow            # Seguir usuário
DELETE /users/{id}/unfollow        # Parar de seguir
GET /users/{id}/followers          # Seguidores
GET /users/{id}/following          # Seguindo
```

### 🟡 **PRIORIDADE MÉDIA - Funcionalidades Avançadas**

#### P2.1 - Sistema de Posts (RF-03.3)
**Descrição:** Posts de usuários com controle de visibilidade
**Complexidade:** Média
**Dependências:** P1.2 (Sistema Social)

**Tarefas:**
- [ ] Implementar repositório de Posts
- [ ] Criar serviço com lógica de visibilidade
- [ ] Desenvolver handlers para CRUD de posts
- [ ] Sistema de feed baseado em amigos/seguindo
- [ ] Integração com sistema de privacidade

**Endpoints:**
```
POST /posts                 # Criar post
GET /posts/feed            # Feed personalizado
GET /users/{id}/posts      # Posts do usuário
PUT /posts/{id}            # Editar post
DELETE /posts/{id}         # Deletar post
```

#### P2.2 - Match de Filmes (RF-04.1 até RF-04.4)
**Descrição:** Sistema colaborativo para escolha de filmes
**Complexidade:** Alta
**Dependências:** P1.2 (Sistema Social)

**Tarefas:**
- [ ] Implementar repositório de MatchSession
- [ ] Algoritmo de sugestão baseado em preferências
- [ ] Sistema de interações (like/dislike)
- [ ] Detecção automática de matches
- [ ] Notificações em tempo real
- [ ] Interface para gestão de sessões

**Endpoints:**
```
POST /match-sessions/start           # Iniciar sessão
POST /match-sessions/{id}/invite     # Convidar usuários
GET /match-sessions/{id}/suggestions # Obter sugestões
POST /match-sessions/{id}/interact   # Registrar like/dislike
GET /match-sessions/{id}/matches     # Ver matches encontrados
POST /match-sessions/{id}/finish     # Finalizar sessão
```

### 🟢 **PRIORIDADE BAIXA - Funcionalidades Complementares**

#### P3.1 - Sistema de Filas Assíncronas (RNF-07)
**Descrição:** Processamento assíncrono para e-mails e notificações
**Complexidade:** Média

**Tarefas:**
- [ ] Configurar Redis como message broker
- [ ] Implementar workers para tarefas assíncronas
- [ ] Sistema de retry para falhas
- [ ] Monitoramento de filas
- [ ] Integração com sistema de e-mail

#### P3.2 - Sistema de Notificações Avançado (RF-06.1)
**Descrição:** Notificações em tempo real via WebSocket/SSE
**Complexidade:** Alta
**Dependências:** P1.2, P2.1, P2.2

**Tarefas:**
- [ ] Escolher tecnologia (WebSocket vs SSE)
- [ ] Implementar servidor de notificações
- [ ] Sistema de persistência de notificações
- [ ] Notificações push para mobile
- [ ] Preferências de notificação do usuário

**Endpoints:**
```
GET /notifications          # Listar notificações
PUT /notifications/{id}/read # Marcar como lida
DELETE /notifications/{id}   # Deletar notificação
PUT /notifications/settings  # Configurar preferências
WebSocket: /ws/notifications # Tempo real
```

---

## 4. Cronograma Sugerido

### **Sprint 1-2 (2 semanas):** Sistema de E-mail
- Implementar confirmação de e-mail e recuperação de senha
- Configurar templates e testes

### **Sprint 3-5 (3 semanas):** Sistema Social Básico
- Amizades e sistema de seguir
- Base para funcionalidades colaborativas

### **Sprint 6-7 (2 semanas):** Sistema de Posts
- Posts com controle de visibilidade
- Feed personalizado

### **Sprint 8-10 (3 semanas):** Match de Filmes
- Funcionalidade principal diferencial
- Sistema colaborativo complexo

### **Sprint 11-12 (2 semanas):** Infraestrutura Avançada
- Sistema de filas
- Notificações em tempo real

---

## 5. Considerações Técnicas

### **Arquitetura Existente:**
- ✅ Clean Architecture bem implementada
- ✅ Sistema de autenticação robusto
- ✅ Cache e integração externa funcionando
- ✅ Testes e documentação estruturados

### **Tecnologias a Adicionar:**
- **E-mail:** SendGrid ou AWS SES
- **Filas:** Redis Pub/Sub ou RabbitMQ
- **Tempo Real:** WebSockets com gorilla/websocket
- **Templates:** html/template para e-mails

### **Banco de Dados:**
- Revisar migrations para novas tabelas sociais
- Índices para performance em consultas sociais
- Considerar particionamento para notificações

---

## 6. Métricas de Sucesso

### **Sistema de E-mail:**
- Taxa de confirmação de e-mail > 80%
- Tempo de entrega < 30 segundos
- Taxa de reset de senha bem-sucedido > 90%

### **Sistema Social:**
- Tempo de resposta para operações sociais < 200ms
- Suporte a 10k+ conexões sociais por usuário
- 99.9% de consistência em relacionamentos

### **Match de Filmes:**
- Sessões simultâneas suportadas > 1000
- Tempo para detectar match < 1 segundo
- Taxa de satisfação com sugestões > 75%

---

## 7. Próximos Passos

1. **Revisar e aprovar este RFC**
2. **Criar issues detalhados para P1.1 (Sistema de E-mail)**
3. **Configurar ambiente de desenvolvimento para e-mails**
4. **Implementar testes automatizados para novos fluxos**
5. **Atualizar documentação da API**

---

**Observação:** Este RFC considera a arquitetura atual bem implementada como base sólida. As funcionalidades pendentes seguem os mesmos padrões de qualidade e organização já estabelecidos no projeto.
