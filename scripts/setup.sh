#!/bin/bash

# CineVerse Development Environment Setup Script
# Este script configura o ambiente de desenvolvimento completo

set -e

echo "🎬 CineVerse - Configurando ambiente de desenvolvimento..."

# Cores para output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Função para logging colorido
log_info() {
    echo -e "${BLUE}ℹ️  $1${NC}"
}

log_success() {
    echo -e "${GREEN}✅ $1${NC}"
}

log_warning() {
    echo -e "${YELLOW}⚠️  $1${NC}"
}

log_error() {
    echo -e "${RED}❌ $1${NC}"
}

# Verificar se Docker está instalado
if ! command -v docker &> /dev/null; then
    log_error "Docker não está instalado. Por favor, instale o Docker primeiro."
    exit 1
fi

# Verificar se Docker Compose está instalado
if ! command -v docker-compose &> /dev/null && ! docker compose version &> /dev/null; then
    log_error "Docker Compose não está instalado. Por favor, instale o Docker Compose primeiro."
    exit 1
fi

# Verificar se Node.js está instalado (para lint-staged)
if ! command -v node &> /dev/null; then
    log_warning "Node.js não está instalado. Instalando dependências com npm pode falhar."
fi

log_info "Configurando hooks de pré-commit..."

# Instalar dependências do Node.js se existir package.json
if [ -f "package.json" ]; then
    if command -v npm &> /dev/null; then
        npm install
        log_success "Dependências Node.js instaladas"
        
        # Configurar husky se disponível
        if [ -d "node_modules/.bin" ] && [ -f "node_modules/.bin/husky" ]; then
            npx husky install
            echo 'npx lint-staged' > .husky/pre-commit
            chmod +x .husky/pre-commit
            log_success "Hooks de pré-commit configurados"
        fi
    else
        log_warning "npm não encontrado. Pule a configuração de hooks se não precisar."
    fi
fi

log_info "Parando containers existentes..."
docker-compose down -v 2>/dev/null || true

log_info "Limpando volumes antigos..."
docker volume prune -f 2>/dev/null || true

log_info "Construindo imagens Docker..."
docker-compose build --no-cache

log_info "Inicializando serviços..."
docker-compose up -d postgres redis

log_info "Aguardando PostgreSQL inicializar..."
sleep 10

# Verificar se PostgreSQL está rodando
log_info "Verificando conexão com PostgreSQL..."
max_attempts=30
attempt=1

while [ $attempt -le $max_attempts ]; do
    if docker-compose exec -T postgres pg_isready -U cineverse &> /dev/null; then
        log_success "PostgreSQL está funcionando!"
        break
    fi
    
    if [ $attempt -eq $max_attempts ]; then
        log_error "PostgreSQL não conseguiu inicializar após $max_attempts tentativas"
        docker-compose logs postgres
        exit 1
    fi
    
    log_info "Tentativa $attempt/$max_attempts - aguardando PostgreSQL..."
    sleep 2
    ((attempt++))
done

log_info "Iniciando API Go..."
docker-compose up -d api_v2

log_info "Aguardando API inicializar..."
sleep 5

# Verificar se a API está rodando
log_info "Verificando API..."
max_attempts=15
attempt=1

while [ $attempt -le $max_attempts ]; do
    if curl -f -s http://localhost:8080/health &> /dev/null; then
        log_success "API Go está funcionando!"
        break
    fi
    
    if [ $attempt -eq $max_attempts ]; then
        log_error "API não conseguiu inicializar após $max_attempts tentativas"
        docker-compose logs api_v2
        exit 1
    fi
    
    log_info "Tentativa $attempt/$max_attempts - aguardando API..."
    sleep 2
    ((attempt++))
done

log_success "🎉 Ambiente configurado com sucesso!"
echo ""
echo "📋 Serviços disponíveis:"
echo "  • API Go:      http://localhost:8080"
echo "  • Flutter App: http://localhost:3000"
echo "  • PostgreSQL:  localhost:5432 (user: cineverse, db: cineverse)"
echo "  • Redis:       localhost:6379"
if [[ $start_android =~ ^[Yy]$ ]]; then
    echo "  • Android Studio: http://localhost:6080 (senha: cineverse)"
fi
echo ""
echo "🔧 Comandos úteis:"
echo "  • Ver logs:           docker-compose logs -f [service]"
echo "  • Parar serviços:     docker-compose down"
echo "  • Rebuild:            docker-compose build --no-cache"
echo "  • Reset completo:     docker-compose down -v && ./scripts/setup.sh"
echo ""
echo "📚 Próximos passos:"
echo "  1. Abra http://localhost:3000 para ver o Flutter app"
echo "  2. Teste a API em http://localhost:8080/health"
echo "  3. Comece a desenvolver! Os arquivos têm hot-reloading ativo."
echo ""
log_success "Happy coding! 🚀"