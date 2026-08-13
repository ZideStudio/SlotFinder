# docker

Docker Compose files and Traefik configuration used for local development and CI/CD builds/deployments.

## Compose files

| File | Purpose |
|---|---|
| `docker-compose.dev.yml` | Local dev dependencies: Traefik (HTTPS reverse proxy) + Postgres + `devtools` service (VS Code devcontainer). Used by `make docker-deps`. |
| `docker-compose.traefik-host.yml` | Override for `docker-compose.dev.yml` that points Traefik at the front/back running natively on the host machine (`host.docker.internal`) instead of the `devtools` container. Used by `make docker-deps` together with the file above. |
| `docker-compose-build-prd.yml` | Builds the front/back Docker images for production (CI). |
| `docker-compose-build-stg.yml` | Builds the front/back Docker images for staging (CI). |
| `docker-compose-prd.yml` | Deploys the front/back services in production, with Traefik labels (CI, self-hosted environment). |
| `docker-compose-stg.yml` | Deploys the front/back services in staging, with Traefik labels (CI, self-hosted environment). |

## Other files

| File | Purpose |
|---|---|
| `Dockerfile.traefik` | Traefik image for local dev, with a self-signed TLS certificate generated at build time (enables HTTPS on `localhost`). |
| `traefik-dynamic.yml` | Traefik dynamic config for local dev: routes to the `devtools` container (front `:3000`, back `:3001`). |
| `traefik-dynamic.host.yml` | Traefik dynamic config for host mode: routes to `host.docker.internal` instead of the `devtools` container. |

## Notes

Relative paths in these files (`context`, volumes, `env_file`) are resolved relative to this `docker/` folder, not the repo root — this is Docker Compose's default behavior (the directory of the first `-f` file is used as the project directory). All commands (`Makefile`, CI workflows, `devcontainer.json`) therefore reference these files as `docker/docker-compose-*.yml` from the repo root.
