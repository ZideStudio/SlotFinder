.PHONY: dev dev-build dev-down clean logs

# Default target
all: start

# Extract Volta versions from package.json
get-volta-versions:
	@command -v jq >/dev/null 2>&1 || { echo "Error: jq is required but not installed. Please install jq to continue." >&2; exit 1; }
	@NODE_VERSION=$$(jq -r '.volta.node // empty' front/package.json) && \
	NPM_VERSION=$$(jq -r '.volta.npm // empty' front/package.json) && \
	if [ -z "$$NODE_VERSION" ] || [ -z "$$NPM_VERSION" ]; then \
		echo "Error: Volta configuration not found in front/package.json" >&2; \
		exit 1; \
	fi && \
	echo "NODE_VERSION=$$NODE_VERSION" > .env.volta && \
	echo "NPM_VERSION=$$NPM_VERSION" >> .env.volta

# Start development environment
start: get-volta-versions
	docker compose -f docker-compose.dev.yml up -d
	@echo "\n🚀 Development environment started!"
	@echo "📱 Front: https://localhost"
	@echo "🔧 API: https://localhost/api"
	@echo "🔧 API Doc: https://localhost/api/swagger/index.html"
	@echo "📊 Traefik Dashboard: http://localhost:9000"
	@echo "🗄️ Database: localhost:5432 (user: slotfinder, password: slotfinder, db: slotfinder)"

# Build and start development environment
build-start: get-volta-versions
	docker compose -f docker-compose.dev.yml up -d --build
	@echo "\n🚀 Development environment built and started!"
	@echo "📱 Front: https://localhost"
	@echo "🔧 API: https://localhost/api"
	@echo "🔧 API Doc: https://localhost/api/swagger/index.html"
	@echo "📊 Traefik Dashboard: http://localhost:9000"
	@echo "🗄️ Database: localhost:5432 (user: slotfinder, password: slotfinder, db: slotfinder)"

# Stop development environment
stop:
	docker compose -f docker-compose.dev.yml stop

# Stop and remove development environment containers
down:
	docker compose -f docker-compose.dev.yml down

# Clean development environment (remove containers and volumes)
clean:
	docker compose -f docker-compose.dev.yml down -v --remove-orphans
	@rm -f .env.volta

# Tail logs of backend and frontend services
logs:
	docker compose -f docker-compose.dev.yml logs -f backend frontend
