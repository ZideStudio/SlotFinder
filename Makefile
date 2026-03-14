.PHONY: dev dev-build dev-down clean logs

# Default target
all: start

# Start development environment
start:
	docker compose -p slotfinder -f docker-compose.dev.yml up -d
	@echo "\n🚀 Development environment started!"
	@echo "📱 Front: https://localhost"
	@echo "📒 Storybook: http://localhost:3002"
	@echo "🔧 API: https://localhost/api"
	@echo "🔧 API Doc: https://localhost/api/swagger/index.html"
	@echo "📊 Traefik Dashboard: http://localhost:9000"
	@echo "🗄️ Database: localhost:5432 (user: slotfinder, password: slotfinder, db: slotfinder)"

# Build and start development environment
build-start:
	docker compose -p slotfinder -f docker-compose.dev.yml up -d --build
	@echo "\n🚀 Development environment built and started!"
	@echo "📱 Front: https://localhost"
	@echo "📒 Storybook: http://localhost:3002"
	@echo "🔧 API: https://localhost/api"
	@echo "🔧 API Doc: https://localhost/api/swagger/index.html"
	@echo "📊 Traefik Dashboard: http://localhost:9000"
	@echo "🗄️ Database: localhost:5432 (user: slotfinder, password: slotfinder, db: slotfinder)"

# Stop development environment
stop:
	docker compose -p slotfinder -f docker-compose.dev.yml stop

# Stop and remove development environment containers
down:
	docker compose -p slotfinder -f docker-compose.dev.yml down

# Clean development environment (remove containers and volumes)
clean:
	docker compose -p slotfinder -f docker-compose.dev.yml down -v --remove-orphans

# Tail logs of backend and frontend services
logs:
	docker compose -p slotfinder -f docker-compose.dev.yml logs -f backend frontend
