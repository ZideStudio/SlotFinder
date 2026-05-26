# Déploiement CI/CD - SlotFinder

## Vue d'ensemble

Le projet SlotFinder utilise un pipeline CI/CD complet avec GitHub Actions pour le déploiement automatisé en staging et production.

## Architecture

```
Frontend (Node/React) + Backend (Go)
         ↓
    Docker Build
         ↓
    External Traefik
    (centralized router)
         ↓
   Staging/Production
```

**Note**: L'application utilise le réseau Traefik externe existant (`zide-traefik-n`) partagé. Assurez-vous que ce réseau existe et est configuré correctement.

## Déploiement Staging

**Déclencheur:** Push sur la branche `ci/cd` (test) ou `main` (production)

1. **Build**: Les images Docker du frontend et backend sont construites avec le tag du commit SHA
2. **Deploy**: Les services sont déployés via `docker-compose-stg.yml`
3. **Domaine**: `stg.slotfinder.fr`
4. **Routing**: Via labels Docker Traefik, sans ports exposés
   - Frontend: `https://stg.slotfinder.fr/`
   - Backend API: `https://stg.slotfinder.fr/api/*`

## Déploiement Production

**Déclencheur:** Création d'un tag `v*` (ex: `v1.0.0`, `v2.1.3`)

1. **Build**: Les images Docker du frontend et backend sont construites avec le tag de version
2. **Deploy**: Les services sont déployés via `docker-compose-prd.yml`
3. **Domaine**: `slotfinder.fr` (sans sous-domaine)
4. **Routing**: Via labels Docker Traefik, sans ports exposés
   - Frontend: `https://slotfinder.fr/`
   - Backend API: `https://slotfinder.fr/api/*`

## Flux de déploiement

### Staging

```
git push origin ci/cd (ou main)
    ↓
setup job
    ↓
build-staging (images avec tag SHA)
    ↓
deploy-staging (docker-compose-stg.yml up -d)
    ↓
✅ Application disponible sur stg.slotfinder.fr via Traefik
```

### Production

```
git tag v1.0.0
git push origin v1.0.0
    ↓
setup job
    ↓
build-production (images avec tag version)
    ↓
deploy-production (docker-compose-prd.yml up -d)
    ↓
✅ Application disponible sur slotfinder.fr via Traefik
```

## Configuration des environnements

### Variables d'environnement du workflow

Le workflow utilise les variables d'environnement suivantes:

```yaml
IMAGE_NAME: slotfinder
STG_DOMAIN: stg.slotfinder.fr
PRD_DOMAIN: slotfinder.fr
```

## Configuration Traefik

L'application utilise le **réseau Traefik externe partagé** (`zide-traefik-n`). Les routes et middlewares sont configurés via les **labels Docker** dans les fichiers `docker-compose`.

### Prérequis

Le réseau Traefik doit exister et être accessible:

```bash
docker network create zide-traefik-n
```

### Configuration des routes

Chaque container expose ses services via des labels Traefik:
- **Frontend**: Port 80, route sans `/api`
- **Backend**: Port 3000, route avec `/api` + middleware stripprefix

Les images Docker sont tagées de la façon suivante:

### Staging
- Frontend: `stg-slotfinder-front:<commit-sha>`
- Backend: `stg-slotfinder-back:<commit-sha>`

### Production
- Frontend: `prd-slotfinder-front:<tag-version>`
- Backend: `prd-slotfinder-back:<tag-version>`

Le workflow utilise `docker compose build` pour construire les images, identique au CI existant (`docker-build.yml`).

## Architecture docker-compose

### Services

Chaque environnement (`docker-compose-stg.yml` et `docker-compose-prd.yml`) contient:

1. **Frontend**: Application React servie par Nginx
   - Port interne: 80
   - Accessible via: `https://<domain>/`
   - Labels Traefik pour routage HTTP/HTTPS

2. **Backend**: API Go
   - Port interne: 3000
   - Accessible via: `https://<domain>/api/*`
   - Labels Traefik avec middleware stripprefix pour `/api`

### Réseau

Les deux services sont connectés au réseau Traefik externe:
- **Staging**: `zide-traefik-n`
- **Production**: `zide-traefik-n`

**Aucun port n'est exposé directement** - tout le trafic passe par Traefik.

## Secrets requis

**Note**: Le workflow utilise `self-hosted` runners. Assurez-vous que les secrets nécessaires sont configurés sur votre serveur.

## Considérations de sécurité

- Les certificats SSL sont gérés automatiquement par Let's Encrypt
- Le backend expose uniquement via Traefik (pas de port direct)
- Les images sont tagées par commit (staging) ou version (production)

## Dépannage

### Le déploiement échoue

1. Vérifiez les logs du workflow dans GitHub Actions
2. Vérifiez que le runner self-hosted est disponible
3. Vérifiez que Docker et Docker Compose sont installés sur le serveur

### Les certificats SSL ne se génèrent pas

1. Vérifiez que les ports 80 et 443 sont accessibles
2. Vérifiez les logs de Traefik: `docker logs <container-traefik>`
3. Vérifiez que le domaine pointe vers le serveur correct

### Le backend ne répond pas à `/api`

1. Vérifiez que le middleware stripprefix est configuré
2. Vérifiez que le backend écoute sur le port 3000
3. Vérifiez les logs du backend

## Mise à jour du déploiement

Pour mettre à jour le déploiement:

1. **Staging**: Créez un commit et pushez sur `main`
2. **Production**: Créez un tag git et pushez-le

```bash
# Release en production
git tag v1.0.0
git push origin v1.0.0
```

## Rollback

Pour revenir à une version précédente:

1. Pushez le tag précédent vers la production
2. Ou modifiez le déploiement manuellement via docker-compose

## Monitoring

Vous pouvez monitorer les déploiements:

1. Via les logs GitHub Actions
2. Via les commandes docker sur le serveur:
   ```bash
   docker compose -f docker-compose-stg.yml logs
   docker compose -f docker-compose-prd.yml logs
   ```
