.PHONY: start build-start stop down clean logs front back storybook

# Default target
all: start

# Start infrastructure (traefik + postgres)
start:
	docker compose -p slotfinder -f docker-compose.dev.yml up -d
	@echo "\n🚀 Infrastructure started!"
	@echo "📊 Traefik Dashboard: http://localhost:9000"
	@echo "🗄️ Database: localhost:5432 (user: slotfinder, password: slotfinder, db: slotfinder)"
	@echo "\nRun 'make front' to start the frontend, 'make back' to start the backend."

# Build and start infrastructure
build-start:
	docker compose -p slotfinder -f docker-compose.dev.yml up -d --build
	@echo "\n🚀 Infrastructure built and started!"
	@echo "📊 Traefik Dashboard: http://localhost:9000"
	@echo "🗄️ Database: localhost:5432 (user: slotfinder, password: slotfinder, db: slotfinder)"
	@echo "\nRun 'make front' to start the frontend, 'make back' to start the backend."

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

# Start frontend (auto-installs node_modules if missing)
front:
	@cd front && [ -d node_modules ] || npm install
	@echo "\n📱 Starting frontend on https://localhost"
	cd front && npm run start

# Start backend with hot reload
back:
	@echo "\n🔧 Starting backend on https://localhost/api"
	cd back && air

# Start storybook (auto-installs node_modules if missing)
storybook:
	@cd front && [ -d node_modules ] || npm install
	@echo "\n📒 Starting Storybook on http://localhost:3002"
	cd front && npm run start:storybook
