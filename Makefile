.PHONY: start front back storybook docker-deps docker-deps-down

# Default target
all: start

# Start frontend and backend with combined logs (installs frontend dependencies)
start:
	@cd front && mise install && npm install
	@cd back && mise install
	@echo "\n🚀 Starting frontend and backend..."
	@echo "📱 Front: https://localhost"
	@echo "🔧 API: https://localhost/api"
	@echo "🔧 API Doc: https://localhost/api/swagger/index.html\n"
	@(cd front && npm run start) & (cd back && go tool air); wait

# Start frontend only (installs frontend dependencies)
front:
	@cd front && mise install && npm install
	@echo "\n📱 Starting frontend on https://localhost"
	cd front && npm run start

# Start backend only with hot reload
back:
	@cd back && mise install
	@echo "\n🔧 Starting backend on https://localhost/api"
	cd back && go tool air

# Start storybook (installs frontend dependencies)
storybook:
	@cd front && mise install && npm install
	@echo "\n📒 Starting Storybook on http://localhost:3002"
	@cd front && mise exec -- sh -c 'npm install && npm run start:storybook'
