.PHONY: start front back storybook docker-deps docker-deps-down docker-deps-logs docker-deps-host docker-deps-host-down docker-deps-host-logs

# Default target
all: start

# Start frontend and backend with combined logs (installs frontend dependencies)
start:
	@cd front && npm install
	@echo "\n🚀 Starting frontend and backend..."
	@echo "📱 Front: https://localhost"
	@echo "🔧 API: https://localhost/api"
	@echo "🔧 API Doc: https://localhost/api/swagger/index.html\n"
	@(cd front && npm run start) & (cd back && go tool air); wait

# Start frontend only (installs frontend dependencies)
front:
	@cd front && npm install
	@echo "\n📱 Starting frontend on https://localhost"
	cd front && npm run start

# Start backend only with hot reload
back:
	@echo "\n🔧 Starting backend on https://localhost/api"
	cd back && go tool air

# Start storybook (installs frontend dependencies)
storybook:
	@cd front && npm install
	@echo "\n📒 Starting Storybook on http://localhost:3002"
	cd front && npm run start:storybook

# Start dockerized dependencies (host mode)
docker-deps:
	@echo "\n🐳 Starting dockerized dependencies..."
	@docker compose -f docker-compose.dev.yml -f docker-compose.traefik-host.yml up -d traefik postgres
	@echo "✅ Traefik is starting on https://localhost (dashboard: http://localhost:9000)"
	@echo "✅ Postgres is starting on localhost:5432"
	@echo "📝 Make sure your frontend is running on http://localhost:3000"
	@echo "📝 Make sure your backend is running on http://localhost:3001"

# Stop dockerized dependencies (host mode)
docker-deps-down:
	@echo "\n🐳 Stopping dockerized dependencies..."
	@docker compose -f docker-compose.dev.yml -f docker-compose.traefik-host.yml down
	@echo "🛑 Dependencies stopped"
