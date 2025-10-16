# 🚀 Sprint 0: Limpeza e Fundação da API CineVerse v2

**Objetivo:** Reconstruir a base da API v2 do zero, mantendo apenas migrations e configurações essenciais

## 📋 Tarefas Principais

### 1. Limpeza Completa do Código
- [ ] **Remover todas as pastas de código**: `internal/` completa
- [ ] **Manter apenas**: 
  - `migrations/` (verificar se estão corretas)
  - `go.mod` e `go.sum`
  - `cmd/main.go` (será recriado)
  - `Dockerfile`
  - Arquivos de configuração

### 2. Verificação e Correção das Migrations
- [ ] **Verificar migrations existentes**:
  - `001_initial_schema.sql` ✓
  - `002_complete_rfc_implementation.sql` ✓
  - `003_social_features.sql` ✓
  - `004_convert_ids_to_uuid.sql` ✓
- [ ] **Corrigir inconsistências** se houver
- [ ] **Criar migration clean** se necessário

### 3. Estrutura Básica Clean Architecture
- [ ] **Criar estrutura de pastas**:
```
internal/
├── config/         # Configurações
├── domain/         # Entidades de negócio
├── repository/     # Acesso a dados
├── service/        # Lógica de negócio
├── handler/        # Controllers HTTP
├── middleware/     # Middlewares
└── server/         # Setup do servidor
```

### 4. Sistema de Configuração
- [ ] **Criar config/config.go**:
  - Configuração de banco de dados
  - Configuração de JWT
  - Configuração de servidor
  - Configuração de TMDb API
- [ ] **Variáveis de ambiente** estruturadas
- [ ] **Validação de configuração** na inicialização

### 5. Sistema de I18n Básico
- [ ] **Estrutura de i18n**:
  - Suporte a EN, PT, ES
  - Sistema de fallback
  - Middleware de detecção de idioma
- [ ] **Mensagens básicas** traduzidas
- [ ] **Integração com sistema de erros**

### 6. Logging Bonito e Informativo
- [ ] **Logger estruturado** com slog
- [ ] **Formatação colorida** para desenvolvimento
- [ ] **Logs de inicialização** informativos:
  - Banner da aplicação
  - Versão e ambiente
  - Configurações carregadas
  - Status da conexão com banco
  - Status da conexão com Redis (se configurado)
  - Porta do servidor
- [ ] **Logs de health check** básicos

### 7. Main.go Limpo e Organizado
- [ ] **Inicialização estruturada**:
  - Carregamento de configuração
  - Setup de logging
  - Conexão com banco de dados
  - Validação de migrations
  - Inicialização do servidor
- [ ] **Graceful shutdown** implementado
- [ ] **Health check endpoint** básico

### 8. Docker e Banco de Dados
- [ ] **Verificar docker-compose.yml**
- [ ] **Testar conexão com PostgreSQL**
- [ ] **Scripts de setup** se necessário

## 🎯 Resultado Esperado

Ao final desta sprint teremos:
1. ✅ **Base limpa** sem código problemático
2. ✅ **Migrations verificadas** e funcionando
3. ✅ **Sistema de configuração** robusto
4. ✅ **I18n básico** funcionando
5. ✅ **Documentação Swagger** funcionando
6. ✅ **Logging informativo** e bonito
7. ✅ **Aplicação rodando** com health check
8. ✅ **Estrutura preparada** para próximas sprints

## 🔧 Comandos Importantes

```bash
# Limpar build artifacts
cd api_v2
rm -rf build/ tmp/ docs/

# Rodar aplicação
docker-compose up api_v2

# Verificar logs
docker-compose logs api_v2

# Testar health check
curl http://localhost:8080/health
```

## 📝 Critérios de Aceitação

- [ ] Aplicação inicia sem erros
- [ ] Logs são informativos e bonitos
- [ ] Health check responde corretamente
- [ ] Banco de dados conecta sem problemas
- [ ] Migrations aplicam corretamente
- [ ] I18n responde com idiomas corretos
- [ ] Configuração carrega via environment variables

## ⏭️ Próxima Sprint

**Sprint 1: Autenticação e Usuários Básicos**
- Sistema de registro e login
- JWT tokens
- Middleware de autenticação
- Endpoints básicos de usuário

---

**Tempo Estimado:** 1-2 dias
**Complexidade:** Baixa
**Prioridade:** CRÍTICA
