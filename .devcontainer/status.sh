#!/bin/bash

# Dev Container Status Dashboard
# Affiche un tableau de bord visuel de l'état du dev container

# Colors
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
CYAN='\033[0;36m'
NC='\033[0m'

echo ""
echo -e "${BLUE}╔════════════════════════════════════════════════════════════╗${NC}"
echo -e "${BLUE}║                 SlotFinder Dev Container                  ║${NC}"
echo -e "${BLUE}║                     Status Dashboard                       ║${NC}"
echo -e "${BLUE}╚════════════════════════════════════════════════════════════╝${NC}"
echo ""

# Check if inside container
if [ -f "/.dockerenv" ] || [ -n "$DEVCONTAINER" ]; then
    echo -e "${GREEN}✓ Running inside dev container${NC}"
    INSIDE_CONTAINER=true
else
    echo -e "${YELLOW}⚠ Running on host machine${NC}"
    INSIDE_CONTAINER=false
fi

echo ""
echo -e "${CYAN}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo -e "${CYAN}  Tools & Versions${NC}"
echo -e "${CYAN}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"

# Check Go
if command -v go &> /dev/null; then
    GO_VERSION=$(go version | awk '{print $3}')
    echo -e "${GREEN}✓${NC} Go:        $GO_VERSION"
else
    echo -e "${RED}✗${NC} Go:        Not installed"
fi

# Check Node
if command -v node &> /dev/null; then
    NODE_VERSION=$(node --version)
    echo -e "${GREEN}✓${NC} Node.js:   $NODE_VERSION"
else
    echo -e "${RED}✗${NC} Node.js:   Not installed"
fi

# Check gopls
if command -v gopls &> /dev/null; then
    echo -e "${GREEN}✓${NC} gopls:     Installed"
else
    echo -e "${RED}✗${NC} gopls:     Not installed"
fi

# Check goimports
if command -v goimports &> /dev/null; then
    echo -e "${GREEN}✓${NC} goimports: Installed"
else
    echo -e "${RED}✗${NC} goimports: Not installed"
fi

# Check TypeScript LSP
if command -v typescript-language-server &> /dev/null; then
    echo -e "${GREEN}✓${NC} TS LSP:    Installed"
else
    echo -e "${RED}✗${NC} TS LSP:    Not installed"
fi

# Check Prettier (local)
if [ -f "/workspace/front/node_modules/.bin/prettier" ]; then
    echo -e "${GREEN}✓${NC} Prettier:  Installed (local)"
elif command -v prettier &> /dev/null; then
    echo -e "${YELLOW}⚠${NC} Prettier:  Installed (global)"
else
    echo -e "${RED}✗${NC} Prettier:  Not installed"
fi

# Check oxlint (local)
if [ -f "/workspace/front/node_modules/.bin/oxlint" ]; then
    echo -e "${GREEN}✓${NC} oxlint:    Installed (local)"
else
    echo -e "${RED}✗${NC} oxlint:    Not installed"
fi

echo ""
echo -e "${CYAN}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo -e "${CYAN}  Services Status${NC}"
echo -e "${CYAN}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"

if ! $INSIDE_CONTAINER; then
    # On host, check Docker containers
    if command -v docker &> /dev/null; then
        FRONTEND=$(docker ps --filter "name=frontend" --format "{{.Status}}" 2>/dev/null | grep -q "Up" && echo "Running" || echo "Stopped")
        BACKEND=$(docker ps --filter "name=backend" --format "{{.Status}}" 2>/dev/null | grep -q "Up" && echo "Running" || echo "Stopped")
        POSTGRES=$(docker ps --filter "name=postgres" --format "{{.Status}}" 2>/dev/null | grep -q "Up" && echo "Running" || echo "Stopped")
        TRAEFIK=$(docker ps --filter "name=traefik" --format "{{.Status}}" 2>/dev/null | grep -q "Up" && echo "Running" || echo "Stopped")
        
        [ "$FRONTEND" = "Running" ] && echo -e "${GREEN}✓${NC} Frontend:  $FRONTEND" || echo -e "${RED}✗${NC} Frontend:  $FRONTEND"
        [ "$BACKEND" = "Running" ] && echo -e "${GREEN}✓${NC} Backend:   $BACKEND" || echo -e "${RED}✗${NC} Backend:   $BACKEND"
        [ "$POSTGRES" = "Running" ] && echo -e "${GREEN}✓${NC} Postgres:  $POSTGRES" || echo -e "${RED}✗${NC} Postgres:  $POSTGRES"
        [ "$TRAEFIK" = "Running" ] && echo -e "${GREEN}✓${NC} Traefik:   $TRAEFIK" || echo -e "${RED}✗${NC} Traefik:   $TRAEFIK"
    else
        echo -e "${RED}✗${NC} Docker not available"
    fi
else
    # Inside container, check connectivity
    if nc -z postgres 5432 2>/dev/null; then
        echo -e "${GREEN}✓${NC} PostgreSQL: Reachable"
    else
        echo -e "${RED}✗${NC} PostgreSQL: Unreachable"
    fi
fi

echo ""
echo -e "${CYAN}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo -e "${CYAN}  Access URLs${NC}"
echo -e "${CYAN}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"

echo -e "🌐 Frontend:         ${BLUE}https://localhost${NC}"
echo -e "🔌 Backend API:      ${BLUE}https://localhost/api${NC}"
echo -e "📚 Storybook:        ${BLUE}http://localhost:3002${NC}"
echo -e "🔄 Traefik:          ${BLUE}http://localhost:9000${NC}"
echo -e "🗄️  PostgreSQL:       ${BLUE}localhost:5432${NC}"

echo ""
echo -e "${CYAN}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo -e "${CYAN}  Quick Commands${NC}"
echo -e "${CYAN}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"

if $INSIDE_CONTAINER; then
    echo "cd /workspace/front && npm run test:unit    # Frontend tests"
    echo "cd /workspace/back && go test ./...         # Backend tests"
    echo "cd /workspace/front && npm run lint         # Lint frontend"
else
    echo "./devcontainer/dev.sh status                # Check status"
    echo "./devcontainer/dev.sh logs                  # View logs"
    echo "./devcontainer/dev.sh shell                 # Enter container"
    echo "./devcontainer/dev.sh test:front            # Run frontend tests"
    echo "./devcontainer/dev.sh test:back             # Run backend tests"
fi

echo ""
echo -e "${BLUE}════════════════════════════════════════════════════════════${NC}"
echo ""
