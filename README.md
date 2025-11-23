# SlotFinder

No more endless discussions to find a date: let us suggest the best times that suit everyone.

## Development Setup

### Prerequisites

Clone the repository:

```bash
git clone https://github.com/ZideStudio/SlotFinder
cd SlotFinder
```

Start the development environment with Docker:

```bash
make dev
```

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
make dev        # Start development environment
make dev-build  # Build and start development environment
make dev-down   # Stop development environment
make clean      # Clean development environment (remove containers and volumes)
```

## Technology Stack

- **Frontend**: React with Rsbuild, TypeScript, Tailwind CSS
- **Backend**: Go with Gin framework, PostgreSQL
- **Development**: Docker Compose with Traefik reverse proxy
- **Hot Reload**: Automatic updates for both frontend and backend

## License

This project is licensed under the MIT License, see the [LICENSE](LICENSE) file for details.
