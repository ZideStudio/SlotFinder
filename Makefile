.PHONY: start front back storybook docker-deps docker-deps-down

# Default target
all: start

# Start frontend and backend with combined logs (installs frontend dependencies)
start:
	@mise install
	@cd front && npm install
	@echo "\n🚀 Starting frontend and backend..."
	@echo "📱 Front: https://localhost"
	@echo "🔧 API: https://localhost/api"
	@echo "🔧 API Doc: https://localhost/api/swagger/index.html\n"
	@(cd front && npm run start) & (cd back && go tool air); wait

# Start frontend only (installs frontend dependencies)
front:
	@mise install
	@cd front && npm install
	@echo "\n📱 Starting frontend on https://localhost"
	cd front && npm run start

# Start backend only with hot reload
back:
	@mise install
	@echo "\n🔧 Starting backend on https://localhost/api"
	cd back && go tool air

# Start storybook (installs frontend dependencies)
storybook:
	@mise install
	@cd front && npm install
	@echo "\n📒 Starting Storybook on http://localhost:3002"
	@cd front && npm run start:storybook
