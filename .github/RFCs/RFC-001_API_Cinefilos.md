# RFC-001: Especificação da API para Rede Social de Cinéfilos

**Autor:** Hype
**Status:** Proposta
**Data de Criação:** 2025-10-08
**Versão:** 1.1

---

## 1. Resumo (Abstract)

Este documento descreve os requisitos e a arquitetura para uma nova API RESTful a ser desenvolvida em Go. A plataforma será uma rede social focada em cinéfilos, permitindo que os usuários descubram, avaliem, salvem e discutam filmes. Além das funcionalidades padrão de listas e avaliações, a API incluirá recursos sociais avançados como seguir usuários, perfis personalizáveis e uma ferramenta interativa para ajudar amigos a escolherem um filme para assistir juntos.

## 2. Motivação

A dificuldade em escolher um filme que agrade a todos em um grupo é um problema comum. Além disso, os cinéfilos buscam plataformas onde possam não apenas registrar suas atividades, mas também interagir com uma comunidade que compartilha de seus interesses. Esta API visa suprir essa necessidade, centralizando o gerenciamento de filmes assistidos, a interação social e a descoberta de novos conteúdos de forma colaborativa e divertida.

## 3. Requisitos Funcionais (RF)

### RF-01: Módulo de Autenticação e Usuários
- **RF-01.1:** O sistema deve permitir o registro de novos usuários através de um endpoint, solicitando no mínimo: nome, e-mail e senha.
- **RF-01.2:** O sistema deve permitir o login de usuários registrados, retornando um token de autenticação (ex: JWT) para acesso a rotas protegidas.
- **RF-01.3:** O usuário deve poder visualizar e editar os dados do seu próprio perfil (nome, bio, foto de perfil, etc.).
- **RF-01.4:** Deve haver um endpoint público para buscar os dados de um perfil de usuário (respeitando as configurações de privacidade).
- **RF-01.5:** O usuário deve ter a opção de tornar seu perfil privado. Se o perfil for privado, suas atividades (listas, reviews) só serão visíveis para seus amigos/seguidores aprovados.
- **`[NOVO]` RF-01.6:** Após o registro, o sistema deve enviar um e-mail de confirmação para o usuário. A conta só será totalmente ativada após a verificação do e-mail.
- **`[NOVO]` RF-01.7:** O sistema deve permitir que o usuário inicie um fluxo de "esqueci minha senha", que enviará um link de redefinição para o seu e-mail.
- **`[NOVO]` RF-01.8:** A API deve registrar informações de segurança a cada login, como endereço de IP, User-Agent (navegador/sistema operacional).
- **`[NOVO]` RF-01.9:** O usuário deve poder visualizar suas sessões ativas (dispositivos logados) e ter a opção de deslogar de sessões específicas ou de todas as outras.
- **`[NOVO]` RF-01.10:** O usuário deve poder escolher um tema de preferência (ex: "light", "dark") que ficará salvo em seu perfil para ser consumido pelo front-end.

### RF-02: Gerenciamento de Filmes e Listas
- **RF-02.1:** O usuário deve poder adicionar um filme à sua lista de "Quero Assistir".
- **RF-02.2:** O usuário deve poder mover um filme da lista "Quero Assistir" para a lista "Já Assisti".
- **RF-02.3:** O usuário deve poder dar uma nota (ex: de 1 a 10) para os filmes que já assistiu.
- **RF-02.4:** O usuário deve poder escrever um review textual para um filme que já assistiu.
- **RF-02.5 (Sugestão):** O usuário deve poder criar listas personalizadas (ex: "Meus Top 10 de Terror", "Filmes para o Fim de Semana").

### RF-03: Módulo Social
- **RF-03.1:** O usuário deve poder enviar e aceitar/recusar pedidos de amizade.
- **RF-03.2:** O sistema deve implementar um modelo de "seguir", onde um usuário pode seguir outro para ver suas atividades públicas no feed, sem necessidade de amizade mútua.
- **RF-03.3:** O usuário deve poder fazer posts em seu perfil, que podem ser públicos (visíveis para todos) ou privados (visíveis apenas para amigos/seguidores).

### RF-04: Módulo de "Match de Filme"
- **`[ATUALIZADO]` RF-04.1:** Um usuário (o "anfitrião") deve poder iniciar uma sessão de "match de filme" e convidar **um ou mais amigos**.
- **RF-04.2:** O sistema deve gerar uma lista de filmes sugeridos com base nos gêneros preferidos de **todos os usuários** na sessão.
- **RF-04.3:** Para cada filme sugerido, os usuários poderão indicar se gostaram (swipe para a direita) ou não (swipe para a esquerda). A API deve registrar essas interações.
- **`[ATUALIZADO]` RF-04.4:** Quando **todos os usuários** na sessão indicarem que gostaram do mesmo filme, a API deve registrar um "match" e notificar todos os participantes.
- **RF-04.5:** Após um match, os usuários podem decidir assistir ao filme (adicionando-o à lista "Quero Assistir", por exemplo) ou continuar a busca.

### RF-05: Módulo de Descoberta e Informações
- **RF-05.1:** A API deve se integrar com um serviço externo (como o The Movie Database - TMDb) para obter informações detalhadas dos filmes (pôster, sinopse, elenco, etc.).
- **RF-05.2:** Devem existir endpoints para buscar filmes por título, gênero, ator ou diretor.
- **RF-05.3:** A API deve fornecer endpoints para listar filmes populares, em alta e lançamentos recentes.

### RF-06: Sistema de Notificações (Sugestão)
- **RF-06.1:** O sistema deve notificar o usuário quando:
    - Receber um novo pedido de amizade/seguidor.
    - Alguém interagir com seu review ou post (curtir, comentar).
    - Ocorrer um "match de filme".

## 4. Requisitos Não Funcionais (RNF)

- **RNF-01 (Desempenho):** A API deve ter um tempo de resposta médio inferior a 500ms para 95% das requisições.
- **RNF-02 (Segurança):** Senhas devem ser armazenadas utilizando algoritmos de hash seguros (ex: bcrypt). Todas as rotas que manipulam dados de usuário devem ser protegidas e requerer um token de autenticação válido.
- **RNF-03 (Escalabilidade):** A arquitetura deve ser pensada para suportar um aumento no número de usuários e requisições sem degradação significativa do serviço.
- **`[ATUALIZADO]` RNF-04 (Documentação):** A API deve ser bem documentada. Devem ser definidos DTOs (Data Transfer Objects) para todas as requisições e respostas. Toda a documentação de rotas e DTOs deve ser escrita **em inglês**, preferencialmente utilizando um padrão como OpenAPI (Swagger) para facilitar o consumo por parte dos desenvolvedores front-end.
- **RNF-05 (Confiabilidade):** A API deve ter um uptime de no mínimo 99.5%.
- **`[NOVO]` RNF-06 (Estratégia de Cache de Dados):** Para requisições de dados de filmes, o sistema deve primeiro consultar o banco de dados local. Se o filme não for encontrado ou os dados estiverem expirados, a API deverá buscar as informações na API externa (ex: TMDb), salvá-las no banco de dados local com um TTL (Time To Live - tempo de validade, ex: 24 horas) e, então, retornar a resposta. Isso otimiza a performance e reduz a dependência de serviços externos.
- **`[NOVO]` RNF-07 (Processamento Assíncrono):** Tarefas que podem demorar e não necessitam de resposta imediata (como envio de e-mails de notificação, e-mail de boas-vindas, etc.) devem ser processadas de forma assíncrona utilizando um sistema de filas (Queues), para não bloquear a thread principal da aplicação e garantir uma resposta rápida ao usuário.
- **`[NOVO]` RNF-08 (Tecnologias Recomendadas):** Recomenda-se o uso de **Redis** para gerenciamento de cache rápido (complementar ao RNF-06) e como message broker para o sistema de filas (RNF-07).

## 5. Modelo de Dados Preliminar (Entidades) `[ATUALIZADO]`

- **User:** `id`, `name`, `email`, `password_hash`, `bio`, `profile_picture_url`, `is_private`, `created_at`, **`email_verified` (boolean)**, **`theme` (string)**.
- **Movie:** `id`, `external_api_id`, `title`, `release_date`, `poster_url`, `genres`, **`cache_expires_at` (datetime)**.
- **UserSession:** **`[NOVO]`** `id`, `user_id`, `token`, `ip_address`, `user_agent`, `created_at`, `expires_at`.
- **Review:** `id`, `user_id`, `movie_id`, `rating`, `content`, `created_at`.
- **MovieList:** `id`, `user_id`, `name`.
- **MovieListEntry:** `id`, `movie_list_id`, `movie_id`.
- **Friendship:** `user_id_1`, `user_id_2`, `status`.
- **Follow:** `follower_id`, `following_id`.
- **Post:** `id`, `user_id`, `content`, `visibility`, `created_at`.
- **MatchSession:** `id`, `host_user_id`, `status` (active, finished), `created_at`.
- **MatchSessionParticipant:** **`[NOVO]`** `session_id`, `user_id`.
- **MatchInteraction:** `session_id`, `user_id`, `movie_id`, `liked` (boolean).

## 6. Arquitetura da API (Endpoints Sugeridos) `[ATUALIZADO]`

**Autenticação e Segurança:**
- `POST /auth/register`
- `POST /auth/login`
- **`[NOVO]`** `POST /auth/confirm-email`
- **`[NOVO]`** `POST /auth/forgot-password`
- **`[NOVO]`** `POST /auth/reset-password`

**Usuários e Configurações:**
- `GET /users/{userId}`
- `PUT /users/me`
- **`[NOVO]`** `PUT /users/me/settings` (Body: `{ "theme": "dark" }`)
- `POST /users/{userId}/follow`
- `GET /users/me/sessions`
- **`[NOVO]`** `DELETE /users/me/sessions/{sessionId}`

**Filmes e Listas:**
- `GET /movies/search?q={title}`
- `GET /movies/popular`
- `POST /users/me/lists/want-to-watch` (Body: `{ "movieId": "..." }`)
- `POST /users/me/lists/watched` (Body: `{ "movieId": "..." }`)

**Reviews e Notas:**
- `POST /movies/{movieId}/reviews` (Body: `{ "rating": 9, "content": "..." }`)
- `GET /movies/{movieId}/reviews`

**Match de Filmes:**
- **`[ATUALIZADO]`** `POST /match-sessions/start` (Body: `{ "guestUserIds": ["...", "..."] }`)
- `GET /match-sessions/{sessionId}/suggestions`
- `POST /match-sessions/{sessionId}/interact` (Body: `{ "movieId": "...", "liked": true }`)

## 7. Questões em Aberto

- Qual será o provedor de dados de filmes (TMDb, IMDb API, etc.)?
- Como as notificações em tempo real serão entregues (WebSockets, Server-Sent Events)?
- Qual será a infraestrutura para o sistema de filas (RabbitMQ, SQS, ou o próprio Redis Pub/Sub)?