# SlotFinder

No more endless discussions to find a date: let us suggest the best times that suit everyone.

## Development Setup

### Prerequisites

- **Docker** and **Docker Compose** installed on your system
- **jq** (command-line JSON processor) - required for extracting Volta versions
  - macOS: `brew install jq`
  - Ubuntu/Debian: `sudo apt-get install jq`
  - Windows: Download from [jqlang.github.io/jq](https://jqlang.github.io/jq/)

#### Clone the repository

```bash
git clone https://github.com/ZideStudio/SlotFinder
cd SlotFinder
```

#### Set up environment variables

Clone the env `backend/.env.model` file to `backend/.env` and modify the variables as needed.

Note that the default values prefixed with `DB_` are already set and work with the dockerized development environment. You can change them if you want to connect to an external database.

#### Start the development environment

Start the development environment with Docker:

```bash
make
```

The Makefile will automatically extract Node.js and npm versions from the Volta configuration in `front/package.json` and use them for the Docker build.

Access the application:

- **Frontend**: https://localhost
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
