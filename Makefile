.PHONY: dev dev-build dev-down clean logs

# Default target
all: dev-build

# Start development environment
dev:
	docker compose -f docker-compose.dev.yml up -d
	@echo "\n🚀 Development environment started!"
	@echo "📱 Front: https://localhost"
	@echo "🔧 API: https://localhost/api"
	@echo "🔧 API Doc: https://localhost/api/swagger/index.html"
	@echo "📊 Traefik Dashboard: http://localhost:9000"
	@echo "🗄️ Database: localhost:5432 (user: slotfinder, password: slotfinder, db: slotfinder)"
	docker compose -f docker-compose.dev.yml logs -f backend frontend

# Build and start development environment
dev-build:
	docker compose -f docker-compose.dev.yml up -d --build
	@echo "\n🚀 Development environment built and started!"
	@echo "📱 Front: https://localhost"
	@echo "🔧 API: https://localhost/api"
	@echo "🔧 API Doc: https://localhost/api/swagger/index.html"
	@echo "📊 Traefik Dashboard: http://localhost:9000"
	@echo "🗄️ Database: localhost:5432 (user: slotfinder, password: slotfinder, db: slotfinder)"
	docker compose -f docker-compose.dev.yml logs -f backend frontend


# Stop development environment
dev-down:
	docker compose -f docker-compose.dev.yml down

# Clean development environment (remove containers and volumes)
clean:
	docker compose -f docker-compose.dev.yml down -v --remove-orphans
