# SlotFinder

No more endless discussions to find a date: let us suggest the best times that suit everyone.

## Development Setup

You have two options to set up your development environment:

1. **Dev Container (Recommended)** - Complete development environment in Docker with IDE integration
2. **Local Setup** - Traditional setup on your host machine

### Option 1: Dev Container (Recommended)

#### Quick Start

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

   # OR open in Zed
   zed .
   # Then: Cmd/Ctrl+Shift+P → "projects: open dev container"
   ```

3. **Wait for Setup**
   - The container builds automatically
   - Dependencies install automatically
   - All services start automatically

4. **Start Coding!**

### Option 2: Local Setup

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

## Technology Stack

- **Frontend**: React with Rsbuild, TypeScript, Sass
- **Backend**: Go with Gin framework, PostgreSQL
- **Development**: Docker Compose with Traefik reverse proxy
- **Hot Reload**: Automatic updates for both frontend and backend

## License

This project is licensed under the MIT License, see the [LICENSE](LICENSE) file for details.
