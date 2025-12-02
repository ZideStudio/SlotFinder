.PHONY: dev dev-build dev-down clean logs

# Default target
all: start

# Extract Volta versions from package.json
get-volta-versions:
	@command -v jq >/dev/null 2>&1 || { echo "Error: jq is required but not installed. Please install jq to continue." >&2; exit 1; }
	@VOLTA_CONFIG=$$(jq -r '.volta // empty' front/package.json) && \
	if [ -z "$$VOLTA_CONFIG" ] || [ "$$VOLTA_CONFIG" = "null" ]; then \
		echo "Error: Volta configuration not found in front/package.json" >&2; \
		exit 1; \
	fi && \
	NODE_VERSION=$$(echo "$$VOLTA_CONFIG" | jq -r '.node // empty') && \
	NPM_VERSION=$$(echo "$$VOLTA_CONFIG" | jq -r '.npm // empty') && \
	if [ -z "$$NODE_VERSION" ] || [ -z "$$NPM_VERSION" ]; then \
		echo "Error: Volta node or npm version not found in front/package.json" >&2; \
		exit 1; \
	fi && \
	printf "NODE_VERSION=%s\nNPM_VERSION=%s\n" "$$NODE_VERSION" "$$NPM_VERSION" > .env.volta.tmp && \
	mv .env.volta.tmp .env.volta

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
	@rm -f .env.volta .env.volta.tmp

# Tail logs of backend and frontend services
logs:
	docker compose -f docker-compose.dev.yml logs -f backend frontend
