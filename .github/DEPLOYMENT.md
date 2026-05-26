# Déploiement CI/CD - SlotFinder

## Vue d'ensemble

Le projet SlotFinder utilise un pipeline CI/CD complet avec GitHub Actions pour le déploiement automatisé en staging et production.

## Architecture

```
Frontend (Node/React) + Backend (Go)
         ↓
    Docker Build
         ↓
   Traefik Router
         ↓
   Staging/Production
```

## Déploiement Staging

**Déclencheur:** Push sur la branche `main`

1. **Build**: Les images Docker du frontend et backend sont construites
2. **Deploy**: Les services sont déployés via `docker-compose-stg.yml`
3. **Domaine**: `stg.slotfinder.fr`
4. **Routing**: 
   - Frontend: `https://stg.slotfinder.fr/`
   - Backend API: `https://stg.slotfinder.fr/api/*`

## Déploiement Production

**Déclencheur:** Création d'un tag `v*` (ex: `v1.0.0`, `v2.1.3`)

1. **Build**: Les images Docker du frontend et backend sont construites
2. **Deploy**: Les services sont déployés via `docker-compose-prd.yml`
3. **Domaine**: `slotfinder.fr` (sans sous-domaine)
4. **Routing**:
   - Frontend: `https://slotfinder.fr/`
   - Backend API: `https://slotfinder.fr/api/*`

## Flux de déploiement

### Staging

```
git push main
    ↓
setup job
    ↓
build-staging (frontend + backend images)
    ↓
deploy-staging (docker-compose-stg.yml up -d)
    ↓
✅ Application disponible sur stg.slotfinder.fr
```

### Production

```
git tag v1.0.0
git push origin v1.0.0
    ↓
setup job
    ↓
build-production (frontend + backend images with tag)
    ↓
deploy-production (docker-compose-prd.yml up -d)
    ↓
✅ Application disponible sur slotfinder.fr
```

## Configuration des environnements

### Variables d'environnement du workflow

Le workflow utilise les variables d'environnement suivantes:

```yaml
IMAGE_NAME: slotfinder
STG_DOMAIN: stg.slotfinder.fr
PRD_DOMAIN: slotfinder.fr
```

### Configuration Traefik

Chaque environnement a sa propre configuration Traefik:

- **Staging**: `docker/traefik-dynamic.stg.yml`
- **Production**: `docker/traefik-dynamic.prd.yml`

Les routes sont configurées pour:
- Diriger les requêtes sans `/api` vers le frontend
- Diriger les requêtes avec `/api` vers le backend
- Retirer automatiquement le préfixe `/api` lors du routage vers le backend

## Images Docker

Les images Docker sont tagées de la façon suivante:

### Staging
- Frontend: `stg-slotfinder-front:<commit-sha>`
- Backend: `stg-slotfinder-back:<commit-sha>`

### Production
- Frontend: `prd-slotfinder-front:<tag-version>`
- Backend: `prd-slotfinder-back:<tag-version>`

## Fichiers Docker

- **Frontend**: `front/Dockerfile.prod` (utilisé pour le build)
- **Backend**: `back/Dockerfile`

## Architecture docker-compose

### Services

1. **Traefik**: Reverse proxy et routeur
   - Port HTTP: 80
   - Port HTTPS: 443
   - Gère les certificats SSL avec Let's Encrypt

2. **Frontend**: Application React
   - Serveur Nginx
   - Port 80 (interne)
   - Accessible via Traefik

3. **Backend**: API Go
   - Port 3000 (interne)
   - Accessible via Traefik avec préfixe `/api`

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
