# OMDb Integration - CineVerse API

## 📋 Overview

Integração completa do OMDb API no CineVerse com pattern de Adapter para facilitar a troca de providers de filmes no futuro.

## 🏗️ Architecture

### Adapter Pattern
Criamos uma interface `MovieProvider` que define o contrato para todos os provedores de dados de filmes:

```go
type MovieProvider interface {
    GetMovieByExternalID(id string) (*MovieDetails, error)
    GetMovieByTitle(title string, year string) (*MovieDetails, error)
    SearchMovies(query string, page int) (*SearchResults, error)
    GetProviderName() string
}
```

Isso permite que você troque facilmente entre OMDb, TMDb ou qualquer outro provider sem alterar o código do aplicativo.

## 📁 Files Created/Modified

### New Files:
- `internal/infrastructure/movie_provider.go` - Interface do adapter
- `internal/infrastructure/omdb.go` - Implementação do OMDb
- `internal/handler/http/omdb_handler.go` - HTTP handlers para OMDb
- `.env` - Variáveis de ambiente (incluindo OM DB_API_KEY)
- `start-server.sh` - Script para iniciar o servidor
- `test-omdb.sh` - Script para testar os endpoints
- `test_db.go` - Utilitário para testar conexão com BD

### Modified Files:
- `internal/config/config.go` - Adicionada configuração do OMDb
- `internal/server/server.go` - Integração do OMDb service e rotas
- `cmd/main.go` - Adicionado carregamento do .env com godotenv

## 🔑 API Key

Sua chave do OMDb (`83a81446`) está configurada no arquivo `.env`:

```bash
OMDB_API_KEY=000000
OMDB_BASE_URL=http://www.omdbapi.com/
```

## 🚀 How to Run

### 1. Start the Server

```bash
cd /home/hype/projects/cine/api_v2
./start-server.sh
```

### 2. Test the Endpoints

Em outro terminal:

```bash
cd /home/hype/projects/cine/api_v2
./test-omdb.sh
```

## 📡 Available Endpoints

### GET `/api/v1/omdb/test`
Testa a conexão com o OMDb (retorna dados do filme The Matrix)

**Response:**
```json
{
  "status": "success",
  "provider": "OMDb",
  "message": "Connection successful",
  "test_movie": "The Matrix",
  "test_imdbId": "tt0133093"
}
```

### GET `/api/v1/omdb/{imdbId}`
Busca filme por IMDb ID

**Example:**
```bash
curl "http://localhost:8080/api/v1/omdb/tt0133093"
```

**Response:**
```json
{
  "title": "The Matrix",
  "year": "1999",
  "runtime": "136 min",
  "genre": "Action, Sci-Fi",
  "director": "Lana Wachowski, Lilly Wachowski",
  "actors": "Keanu Reeves, Laurence Fishburne, Carrie-Anne Moss",
  "plot": "When a beautiful stranger leads computer hacker Neo to a forbidding underworld...",
  "poster": "https://m.media-amazon.com/images/...",
  "imdb_rating": "8.7",
  "imdb_votes": "2,050,000",
  "ratings": [
    {"source": "Internet Movie Database", "value": "8.7/10"},
    {"source": "Rotten Tomatoes", "value": "83%"},
    {"source": "Metacritic", "value": "73/100"}
  ],
  "provider": "OMDb"
}
```

### GET `/api/v1/omdb/title`
Busca filme por título

**Query Parameters:**
- `title` (required): Título do filme
- `year` (optional): Ano de lançamento

**Example:**
```bash
curl "http://localhost:8080/api/v1/omdb/title?title=Inception&year=2010"
```

### GET `/api/v1/omdb/search`
Busca filmes por query

**Query Parameters:**
- `q` (required): Termo de busca
- `page` (optional, default: 1): Número da página (1-100)

**Example:**
```bash
curl "http://localhost:8080/api/v1/omdb/search?q=Batman&page=1"
```

**Response:**
```json
{
  "results": [
    {
      "title": "Batman Begins",
      "year": "2005",
      "type": "movie",
      "poster": "https://...",
      "imdb_id": "tt0372784",
      "provider_id": "tt0372784"
    }
  ],
  "total_results": 563,
  "page": 1,
  "total_pages": 57,
  "provider": "OMDb"
}
```

### GET `/api/v1/omdb/search-by-type`
Busca filmes por query e tipo

**Query Parameters:**
- `q` (required): Termo de busca
- `type` (optional): `movie`, `series`, ou `episode`
- `page` (optional, default: 1): Número da página

**Example:**
```bash
curl "http://localhost:8080/api/v1/omdb/search-by-type?q=Star%20Wars&type=movie&page=1"
```

## 🧪 Manual Testing Examples

### Test 1: Connection
```bash
curl "http://localhost:8080/api/v1/omdb/test" | jq .
```

### Test 2: Get Specific Movie
```bash
# The Shawshank Redemption
curl "http://localhost:8080/api/v1/omdb/tt0111161" | jq .

# Interstellar
curl "http://localhost:8080/api/v1/omdb/tt0816692" | jq .
```

### Test 3: Search
```bash
# Search for Batman movies
curl "http://localhost:8080/api/v1/omdb/search?q=Batman" | jq .

# Search page 2
curl "http://localhost:8080/api/v1/omdb/search?q=Batman&page=2" | jq .
```

### Test 4: Search by Title
```bash
curl "http://localhost:8080/api/v1/omdb/title?title=The%20Dark%20Knight&year=2008" | jq .
```

## 🔄 Switching Providers

Para trocar de provider (ex: de OMDb para TMDb), você só precisa:

1. Implementar a interface `MovieProvider` para o novo provider
2. Atualizar a inicialização no `server.go`:

```go
// Instead of:
omdbService := infrastructure.NewOMDbService(s.config.OMDb.APIKey)

// Use:
tmdbService := infrastructure.NewTMDbService(s.config.TMDb.APIKey)
```

Todos os handlers e DTOs continuam funcionando sem alterações!

## 📊 Data Models

### MovieDetails
Estrutura unificada que funciona para qualquer provider:

```go
type MovieDetails struct {
    Title       string
    Year        string
    Plot        string
    Poster      string
    Genre       string
    Director    string
    Actors      string
    IMDbRating  string
    Ratings     []Rating
    Provider    string  // "OMDb", "TMDb", etc.
    // ... more fields
}
```

### SearchResults
```go
type SearchResults struct {
    Results      []SearchItem
    TotalResults int
    Page         int
    TotalPages   int
    Provider     string
}
```

## 🐛 Troubleshooting

### Database Connection Error
Se você ver `password authentication failed`:

1. Verifique se o PostgreSQL está rodando:
```bash
docker ps | grep postgres
```

2. Verifique a senha no `.env`:
```bash
DB_PASSWORD=cineverse123
```

### OMDb API Error
Se você ver erros do OMDb:

1. Verifique se a API key está correta no `.env`
2. Teste diretamente no navegador:
```
http://www.omdbapi.com/?apikey=83a81446&i=tt0133093
```

### Server Not Starting
```bash
# Rebuild
go build -o build/main ./cmd/main.go

# Check logs
tail -f /tmp/cine-api.log
```

## 📝 Next Steps

1. ✅ OMDb integrado com adapter pattern
2. 🔄 Adicionar cache Redis para respostas do OMDb
3. 🔄 Implementar rate limiting
4. 🔄 Adicionar mais providers (TMDb, etc.)
5. 🔄 Criar endpoints combinados que usam múltiplos providers

## 🎯 Benefits of this Architecture

- **Flexibilidade**: Troque de API facilmente
- **Testabilidade**: Mock providers para testes
- **Manutenibilidade**: Código organizado e desacoplado
- **Extensibilidade**: Adicione novos providers sem quebrar código existente
- **Unified Interface**: Mesma estrutura de dados independente do provider

---

**Author**: GitHub Copilot  
**Date**: November 7, 2025  
**Version**: 1.0.0
