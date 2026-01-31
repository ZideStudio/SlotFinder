#!/bin/bash

# Dev Container Helper Script
# Ce script facilite la gestion du dev container SlotFinder

set -e

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Functions
print_header() {
    echo -e "${BLUE}================================${NC}"
    echo -e "${BLUE}$1${NC}"
    echo -e "${BLUE}================================${NC}"
}

print_success() {
    echo -e "${GREEN}✓ $1${NC}"
}

print_error() {
    echo -e "${RED}✗ $1${NC}"
}

print_info() {
    echo -e "${YELLOW}ℹ $1${NC}"
}

check_docker() {
    if ! command -v docker &> /dev/null; then
        print_error "Docker n'est pas installé"
        exit 1
    fi
    print_success "Docker est installé"
}

check_docker_compose() {
    if ! command -v docker-compose &> /dev/null; then
        print_error "Docker Compose n'est pas installé"
        exit 1
    fi
    print_success "Docker Compose est installé"
}

status() {
    print_header "Status des Services"
    docker-compose -f docker-compose.dev.yml ps
}

logs() {
    print_header "Logs des Services"
    if [ -z "$1" ]; then
        docker-compose -f docker-compose.dev.yml logs -f
    else
        docker-compose -f docker-compose.dev.yml logs -f "$1"
    fi
}

rebuild() {
    print_header "Reconstruction du Dev Container"
    docker-compose -f docker-compose.dev.yml stop devtools
    docker-compose -f docker-compose.dev.yml build --no-cache devtools
    docker-compose -f docker-compose.dev.yml up -d devtools
    print_success "Dev Container reconstruit"
}

restart() {
    print_header "Redémarrage des Services"
    docker-compose -f docker-compose.dev.yml restart
    print_success "Services redémarrés"
}

clean() {
    print_header "Nettoyage"
    print_info "Arrêt de tous les services..."
    docker-compose -f docker-compose.dev.yml down

    read -p "Voulez-vous supprimer les volumes ? (y/N) " -n 1 -r
    echo
    if [[ $REPLY =~ ^[Yy]$ ]]; then
        docker-compose -f docker-compose.dev.yml down -v
        print_success "Services et volumes supprimés"
    else
        print_success "Services arrêtés (volumes conservés)"
    fi
}

setup() {
    print_header "Configuration Initiale"

    check_docker
    check_docker_compose

    print_info "Installation des dépendances frontend..."
    if [ -d "./front" ]; then
        docker-compose -f docker-compose.dev.yml run --rm devtools sh -c "cd /workspace/front && npm install"
        print_success "Dépendances frontend installées"
    else
        print_error "Dossier front/ non trouvé"
    fi

    print_info "Vérification des modules Go..."
    if [ -d "./back" ]; then
        docker-compose -f docker-compose.dev.yml run --rm devtools sh -c "cd /workspace/back && go mod download"
        print_success "Modules Go téléchargés"
    else
        print_error "Dossier back/ non trouvé"
    fi

    print_success "Configuration terminée !"
}

test_frontend() {
    print_header "Tests Frontend"
    docker-compose -f docker-compose.dev.yml exec devtools sh -c "cd /workspace/front && npm run test:unit"
}

test_backend() {
    print_header "Tests Backend"
    docker-compose -f docker-compose.dev.yml exec devtools sh -c "cd /workspace/back && go test ./..."
}

lint_frontend() {
    print_header "Lint Frontend"
    docker-compose -f docker-compose.dev.yml exec devtools sh -c "cd /workspace/front && npm run lint"
}

lint_frontend_fix() {
    print_header "Lint Frontend (Fix)"
    docker-compose -f docker-compose.dev.yml exec devtools sh -c "cd /workspace/front && npm run lint:oxlint:fix && npm run lint:prettier:fix"
}

format_backend() {
    print_header "Format Backend"
    docker-compose -f docker-compose.dev.yml exec devtools sh -c "cd /workspace/back && goimports -w ."
}

shell() {
    print_header "Shell dans le Dev Container"
    docker-compose -f docker-compose.dev.yml exec devtools bash
}

db_shell() {
    print_header "Shell PostgreSQL"
    docker-compose -f docker-compose.dev.yml exec postgres psql -U slotfinder -d slotfinder
}

db_backup() {
    print_header "Backup Base de Données"
    BACKUP_FILE="backup_$(date +%Y%m%d_%H%M%S).sql"
    docker-compose -f docker-compose.dev.yml exec -T postgres pg_dump -U slotfinder slotfinder > "$BACKUP_FILE"
    print_success "Backup créé : $BACKUP_FILE"
}

db_restore() {
    if [ -z "$1" ]; then
        print_error "Usage: $0 db:restore <backup_file>"
        exit 1
    fi
    print_header "Restauration Base de Données"
    docker-compose -f docker-compose.dev.yml exec -T postgres psql -U slotfinder slotfinder < "$1"
    print_success "Base de données restaurée depuis $1"
}

help() {
    cat << EOF
${BLUE}Dev Container Helper pour SlotFinder${NC}

${YELLOW}Usage:${NC}
  $0 <command> [options]

${YELLOW}Commandes disponibles:${NC}

  ${GREEN}Gestion des services:${NC}
    status              Affiche le statut des services
    logs [service]      Affiche les logs (optionnel: pour un service spécifique)
    restart             Redémarre tous les services
    rebuild             Reconstruit le dev container
    clean               Arrête et nettoie les services

  ${GREEN}Configuration:${NC}
    setup               Configuration initiale (première installation)

  ${GREEN}Tests:${NC}
    test:front          Lance les tests frontend
    test:back           Lance les tests backend

  ${GREEN}Linting/Formatting:${NC}
    lint:front          Lint le frontend
    lint:front:fix      Lint et fix le frontend
    format:back         Formate le code backend

  ${GREEN}Accès Shell:${NC}
    shell               Ouvre un shell dans le dev container
    db:shell            Ouvre un shell PostgreSQL

  ${GREEN}Base de données:${NC}
    db:backup           Crée un backup de la base de données
    db:restore <file>   Restaure un backup de la base de données

  ${GREEN}Aide:${NC}
    help                Affiche cette aide

${YELLOW}Exemples:${NC}
  $0 setup              # Première installation
  $0 status             # Voir l'état des services
  $0 logs frontend      # Logs du frontend uniquement
  $0 test:front         # Lancer les tests frontend
  $0 shell              # Accéder au shell du container
  $0 db:backup          # Faire un backup de la DB

EOF
}

# Main
case "$1" in
    status)
        status
        ;;
    logs)
        logs "$2"
        ;;
    rebuild)
        rebuild
        ;;
    restart)
        restart
        ;;
    clean)
        clean
        ;;
    setup)
        setup
        ;;
    test:front)
        test_frontend
        ;;
    test:back)
        test_backend
        ;;
    lint:front)
        lint_frontend
        ;;
    lint:front:fix)
        lint_frontend_fix
        ;;
    format:back)
        format_backend
        ;;
    shell)
        shell
        ;;
    db:shell)
        db_shell
        ;;
    db:backup)
        db_backup
        ;;
    db:restore)
        db_restore "$2"
        ;;
    help|--help|-h|"")
        help
        ;;
    *)
        print_error "Commande inconnue: $1"
        echo ""
        help
        exit 1
        ;;
esac
