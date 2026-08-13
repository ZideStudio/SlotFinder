# docker

Fichiers Docker Compose et configuration Traefik utilisés en développement local et pour les builds/déploiements CI/CD.

## Compose files

| Fichier | But |
|---|---|
| `docker-compose.dev.yml` | Dépendances de dev local : Traefik (reverse proxy HTTPS) + Postgres + service `devtools` (devcontainer VS Code). Utilisé par `make docker-deps`. |
| `docker-compose.traefik-host.yml` | Override de `docker-compose.dev.yml` pour faire pointer Traefik vers le front/back lancés en natif sur la machine hôte (`host.docker.internal`) au lieu du container `devtools`. Utilisé par `make docker-deps` en combinaison avec le fichier ci-dessus. |
| `docker-compose-build-prd.yml` | Build des images Docker front/back pour la production (CI). |
| `docker-compose-build-stg.yml` | Build des images Docker front/back pour la staging (CI). |
| `docker-compose-prd.yml` | Déploiement des services front/back en production, avec labels Traefik (CI, environnement self-hosted). |
| `docker-compose-stg.yml` | Déploiement des services front/back en staging, avec labels Traefik (CI, environnement self-hosted). |

## Autres fichiers

| Fichier | But |
|---|---|
| `Dockerfile.traefik` | Image Traefik pour le dev local, avec certificat TLS auto-signé généré au build (permet HTTPS sur `localhost`). |
| `traefik-dynamic.yml` | Config dynamique Traefik pour le dev local : route vers le container `devtools` (front `:3000`, back `:3001`). |
| `traefik-dynamic.host.yml` | Config dynamique Traefik pour le mode "host" : route vers `host.docker.internal` au lieu du container `devtools`. |

## Notes

Les chemins relatifs dans ces fichiers (`context`, volumes, `env_file`) sont résolus par rapport à ce dossier `docker/`, pas à la racine du repo — c'est le comportement par défaut de Docker Compose (le dossier du premier fichier `-f` sert de project directory). Toutes les commandes (`Makefile`, workflows CI, `devcontainer.json`) référencent donc ces fichiers via `docker/docker-compose-*.yml` depuis la racine du repo.
