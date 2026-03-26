.PHONY: start build-start stop down clean logs front back storybook

# Default target
all: start

# Start frontend and backend with combined logs (auto-installs node_modules if missing)
start:
	@cd front && [ -d node_modules ] || npm install
	@echo "\n🚀 Starting frontend and backend..."
	@echo "📱 Front: https://localhost"
	@echo "🔧 API: https://localhost/api"
	@echo "🔧 API Doc: https://localhost/api/swagger/index.html\n"
	@(cd front && npm run start) & (cd back && air); wait

# Start infrastructure (build docker images)
build-start:
	docker compose -p slotfinder -f docker-compose.dev.yml up -d --build
	@echo "\n🚀 Infrastructure built and started!"
	@echo "📊 Traefik Dashboard: http://localhost:9000"
	@echo "🗄️ Database: localhost:5432 (user: slotfinder, password: slotfinder, db: slotfinder)"

# Stop infrastructure
stop:
	docker compose -p slotfinder -f docker-compose.dev.yml stop

# Stop and remove infrastructure containers
down:
	docker compose -p slotfinder -f docker-compose.dev.yml down

# Clean infrastructure (remove containers and volumes)
clean:
	docker compose -p slotfinder -f docker-compose.dev.yml down -v --remove-orphans

# Tail logs of infrastructure services
logs:
	docker compose -p slotfinder -f docker-compose.dev.yml logs -f

# Start frontend only (auto-installs node_modules if missing)
front:
	@cd front && [ -d node_modules ] || npm install
	@echo "\n📱 Starting frontend on https://localhost"
	cd front && npm run start

# Start backend only with hot reload
back:
	@echo "\n🔧 Starting backend on https://localhost/api"
	cd back && air

# Start storybook (auto-installs node_modules if missing)
storybook:
	@cd front && [ -d node_modules ] || npm install
	@echo "\n📒 Starting Storybook on http://localhost:3002"
	cd front && npm run start:storybook
