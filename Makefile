.PHONY: start front back storybook

# Default target
all: start

# Start frontend and backend with combined logs (auto-installs node_modules if missing)
start:
	@cd front && npm install
	@echo "\n🚀 Starting frontend and backend..."
	@echo "📱 Front: https://localhost"
	@echo "🔧 API: https://localhost/api"
	@echo "🔧 API Doc: https://localhost/api/swagger/index.html\n"
	@(cd front && npm run start) & (cd back && air); wait

# Start frontend only (auto-installs node_modules if missing)
front:
	@cd front && npm install
	@echo "\n📱 Starting frontend on https://localhost"
	cd front && npm run start

# Start backend only with hot reload
back:
	@echo "\n🔧 Starting backend on https://localhost/api"
	cd back && air

# Start storybook (auto-installs node_modules if missing)
storybook:
	@cd front && npm install
	@echo "\n📒 Starting Storybook on http://localhost:3002"
	cd front && npm run start:storybook
