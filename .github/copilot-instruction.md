# Guia de Arquitetatura e Padrões para o Projeto CineVerse

## 1. Visão Geral do Projeto

Você está trabalhando no **CineVerse**, uma rede social para cinéfilos. O objetivo é criar uma plataforma robusta, escalável e de fácil manutenção. As especificações funcionais detalhadas estão nos documentos RFC na pasta `.github/RFCs`, que servem como a fonte da verdade para todas as funcionalidades.

Sua tarefa é gerar código que siga estritamente as convenções, tecnologias e estruturas de pastas definidas neste guia.

## 2. Estrutura de Pastas do Monorepo

A organização do projeto é fundamental. Siga esta estrutura:

```
.
├── 📂 api/            # (LEGADO) API antiga em NestJS. NÃO UTILIZAR.
├── 📂 api_v2/          # ✅ API principal em Go. TODO o trabalho de backend é aqui.
├── 📂 flutter_app/     # ✅ Frontend multiplataforma em Flutter.
├── 📂 scripts/
│   └── setup.sh      # Script para configuração inicial do ambiente.
├── .lintstagedrc.json  # Configuração de pré-commit para qualidade de código.
├── docker-compose.yml  # Orquestra todos os serviços (Go, Flutter, Postgres, Redis).
└── README.md
```

-   **Hot-Reloading**: O ambiente `docker-compose.yml` está configurado para hot-reloading tanto no backend Go quanto no frontend Flutter.

## 3. Diretrizes de Desenvolvimento Backend (Go - `api_v2`)

A API deve ser desacoplada, testável e performática.

### 3.1. Filosofia Arquitetural Obrigatória: Clean Architecture

-   As dependências devem sempre apontar para o centro (domínio). A lógica de negócio não pode depender de detalhes de infraestrutura (framework web, banco de dados).
-   Utilize **Injeção de Dependência** para fornecer implementações (ex: repositórios) para as camadas de serviço.

### 3.2. Estrutura de Pastas (Layout Padrão)

Siga o [Standard Go Project Layout](https://github.com/golang-standards/project-layout):

-   **/cmd**: Contém o `main.go`. Responsável por inicializar configurações, dependências (banco de dados, etc.) e o servidor HTTP.
-   **/internal**: Contém todo o código-fonte principal da aplicação.
    -   **/internal/handler**: Camada de transporte (HTTP). Recebe requisições, valida DTOs e chama os serviços. **Use o framework `chi` aqui.**
    -   **/internal/service**: Camada de serviço. Orquestra a lógica de negócio e os casos de uso.
    -   **/internal/repository**: Camada de acesso a dados. Implementa a lógica de comunicação com o banco de dados e o cache.
    -   **/internal/domain**: Camada de domínio. Contém as entidades de negócio puras (structs) e as regras de negócio mais importantes.

### 3.3. Stack de Tecnologias e Bibliotecas (Padrão)

Utilize exclusivamente as seguintes tecnologias para suas respectivas finalidades:

| Finalidade             | Tecnologia/Biblioteca                                   | Motivo                                                      |
| :--------------------- | :------------------------------------------------------ | :---------------------------------------------------------- |
| **Banco de Dados** | **PostgreSQL** | Banco de dados relacional principal.                        |
| **Acesso a Dados (SQL)** | **`sqlx`** | Para mapear queries para structs Go de forma segura.        |
| **Router HTTP** | **`chi`** | Leve, idiomático e excelente para APIs RESTful.             |
| **Cache e Filas** | **Redis** | Para cache de dados e processamento assíncrono.             |
| **Configuração** | **`spf13/viper`** | Gerenciamento de configurações via arquivos e env vars.     |
| **Validação de DTOs** | **`go-playground/validator`** | Validação declarativa de structs com tags.                  |
| **Logging** | **`slog`** (nativo do Go)                               | Logging estruturado e performático.                         |
| **Qualidade de Código**| **`gofmt`** e **`go vet`** (via `lint-staged`)          | Formatação e análise estática padrão do Go.                 |

## 4. Diretrizes de Desenvolvimento Frontend (Flutter - `flutter_app`)

O aplicativo deve ser bem estruturado, reativo e preparado para múltiplas plataformas.

### 4.1. Filosofia Arquitetural Obrigatória: Feature-First

-   Organize o código em módulos de funcionalidades. Cada funcionalidade deve ser o mais autocontida possível.

### 4.2. Estrutura de Pastas

-   `lib/src/core`: Código compartilhado por toda a aplicação (tema, cliente API, modelos de dados, constantes, etc.).
-   `lib/src/features`: Cada subpasta aqui é uma funcionalidade (ex: `authentication`, `movie_details`, `match_session`).
    -   Dentro de cada feature, separe as camadas: `data` (repositórios, fontes de dados), `domain` (entidades, casos de uso) e `presentation` (widgets, telas, controllers/state).

### 4.3. Stack de Tecnologias e Bibliotecas (Padrão)

| Finalidade               | Tecnologia/Biblioteca                        | Motivo                                                                |
| :----------------------- | :------------------------------------------- | :-------------------------------------------------------------------- |
| **Gerenciamento de Estado** | **`flutter_riverpod`** | Solução reativa, compilável e segura para gerenciamento de estado.      |
| **Navegação** | **`go_router`** | Roteamento declarativo, ideal para deep linking e navegação complexa. |
| **Injeção de Dependência** | **`get_it`** | Service Locator para desacoplar a criação de dependências.            |
| **Cliente HTTP** | **`dio`** | Cliente HTTP poderoso com interceptors, cancelamento, etc.          |
| **Qualidade de Código** | **`flutter format`** e **`flutter analyze`** | Ferramentas padrão para formatação e análise estática do Dart.        |

## 5. Padrões de Código e Versionamento

Para manter a consistência e a qualidade do projeto, siga as regras abaixo.

### 5.1. Mensagens de Commit (Conventional Commits)

-   **Regra**: Todas as mensagens de commit devem seguir a especificação [**Conventional Commits**](https://www.conventionalcommits.org/).
-   **Estrutura**: `<type>(<scope>): <description>`
-   **Tipos Comuns**:
    -   `feat`: Uma nova funcionalidade.
    -   `fix`: Correção de um bug.
    -   `docs`: Alterações na documentação.
    -   `style`: Alterações de formatação, sem impacto na lógica.
    -   `refactor`: Refatoração de código que não corrige bug nem adiciona funcionalidade.
    -   `test`: Adição ou modificação de testes.
    -   `chore`: Manutenção do build, dependências, etc.
-   **Commits Atômicos**: Mantenha os commits pequenos e focados em uma única responsabilidade. Evite commits gigantes com múltiplas alterações não relacionadas.

### 5.2. Comentários no Código

-   **Regra Principal**: Evite comentários. O código deve ser limpo, claro e autoexplicativo através de nomes de variáveis, funções e classes bem escolhidos.
-   **Exceção**: Comentários são permitidos apenas em casos **estritamente necessários** para explicar algoritmos complexos ou lógicas de negócio não triviais que não podem ser simplificadas.
-   **Idioma**: Se um comentário for necessário, ele deve ser escrito **obrigatoriamente em inglês**.