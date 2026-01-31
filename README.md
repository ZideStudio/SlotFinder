# SlotFinder

No more endless discussions to find a date: let us suggest the best times that suit everyone.

## Development Setup

You have two options to set up your development environment:

1. **Dev Container (Recommended)** - Complete development environment in Docker with IDE integration
2. **Local Setup** - Traditional setup on your host machine

### Option 1: Dev Container (Recommended)

The dev container provides a fully configured development environment with:
- ✅ Go 1.25.5 with gopls (LSP), goimports, delve, and air
- ✅ Node.js 24.13.0 with TypeScript Language Server
- ✅ Prettier and oxlint from project dependencies
- ✅ PostgreSQL with automatic setup
- ✅ All services (frontend, backend, database, traefik) running automatically
- ✅ Format on save, autocompletion, and hot reload enabled
- ✅ Works with VS Code and Zed

#### Quick Start with Dev Container

1. **Prerequisites**
   - Docker and Docker Compose installed
   - VS Code with [Dev Containers extension](https://marketplace.visualstudio.com/items?itemName=ms-vscode-remote.remote-containers) OR Zed (recent version)

2. **Open in Container**
   ```bash
   # Clone the repository
   git clone https://github.com/ZideStudio/SlotFinder
   cd SlotFinder
   
   # Open in VS Code
   code .
   # Then: Cmd/Ctrl+Shift+P → "Dev Containers: Reopen in Container"
   
   # OR open in Zed (will auto-detect dev container)
   zed .
   ```

3. **Wait for Setup**
   - The container builds automatically (first time only)
   - Dependencies install automatically
   - All services start automatically

4. **Start Coding!**
   - Open any `.go` or `.tsx` file
   - Enjoy autocompletion, format on save, and hot reload

#### Dev Container Documentation

- **Quick Start Guide**: [.devcontainer/QUICKSTART.md](.devcontainer/QUICKSTART.md)
- **Complete Documentation**: [.devcontainer/README.md](.devcontainer/README.md)
- **Verification Checklist**: [.devcontainer/CHECKLIST.md](.devcontainer/CHECKLIST.md)
- **Helper Script**: Run `./devcontainer/dev.sh help` for available commands

### Option 2: Local Setup

#### Prerequisites

#### Clone the repository

```bash
git clone https://github.com/ZideStudio/SlotFinder
cd SlotFinder
```

#### Install frontend dependencies

Install the required packages for the frontend:

```bash
cd front
npm install
cd ..
```

#### Set up environment variables

Clone the env `backend/.env.model` file to `backend/.env` and modify the variables as needed.

Note that the default values prefixed with `DB_` are already set and work with the dockerized development environment. You can change them if you want to connect to an external database.

#### Start the development environment

Start the development environment with Docker:

```bash
make
```

Access the application:

- **Frontend**: https://localhost
- **Storybook**: http://localhost:3002
- **Backend API**: https://localhost/api
- **Traefik Dashboard**: http://localhost:9000

**Note**: The development environment uses self-signed certificates. Your browser will show a security warning - this is normal for local development. You can safely proceed by clicking "Advanced" > "Proceed to localhost".

### Database Access

Connect to the development database from your host machine:

```
Host: localhost
Port: 5432
Username: slotfinder
Password: slotfinder
Database: slotfinder
```

Example connection:

```bash
psql -h localhost -p 5432 -U slotfinder -d slotfinder
```

### Available Commands

```bash
make                # Start development environment
make start          # Start development environment
make start-build    # Build and start development environment
make stop           # Stop development environment
make down           # Tear down development environment
make clean          # Clean development environment (remove containers and volumes)
make logs           # View logs of frontend and backend services
```

## Development Tools

### IDE Configuration

Both Zed and VS Code are fully configured for this project:

- **Format on Save**: Automatic code formatting for Go, TypeScript, JavaScript, CSS, SCSS, JSON, YAML, Markdown
- **LSP Integration**: Full autocompletion and IntelliSense for Go and TypeScript
- **Hot Reload**: Changes are automatically reflected without manual restarts
- **Linting**: oxlint for JavaScript/TypeScript, integrated Go linting
- **Debugging**: Debug configurations included for both frontend and backend

Configuration files:
- `.zed/settings.json` - Zed IDE configuration
- `.vscode/settings.json` - VS Code configuration
- `.vscode/tasks.json` - Predefined tasks
- `.vscode/launch.json` - Debug configurations

### Testing

```bash
# Frontend tests
cd front
npm run test:unit           # Unit tests
npm run test:unit:watch     # Unit tests in watch mode
npm run test:browser        # Browser tests
npm run lint:tsc            # TypeScript type checking

# Backend tests
cd back
go test ./...               # Run all tests
go test -v ./...            # Verbose output
```

### Linting and Formatting

```bash
# Frontend
cd front
npm run lint                # Lint everything
npm run lint:oxlint:fix     # Fix oxlint issues
npm run lint:prettier:fix   # Fix prettier formatting

# Backend
cd back
goimports -w .              # Format Go code
go vet ./...                # Go static analysis
```

## Technology Stack

- **Frontend**: React with Rsbuild, TypeScript, Sass
- **Backend**: Go with Gin framework, PostgreSQL
- **Development**: Docker Compose with Traefik reverse proxy
- **Hot Reload**: Automatic updates for both frontend and backend

## License

This project is licensed under the MIT License, see the [LICENSE](LICENSE) file for details.
