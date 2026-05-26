# Checklist - SlotFinder CD Deployment

Utilisez cette checklist pour vérifier que tout est prêt pour le déploiement.

## ✅ Prérequis

- [ ] Docker et Docker Compose installés sur le serveur
- [ ] Runner self-hosted GitHub Actions configuré et disponible
- [ ] Réseau Docker `zide-traefik-n` existe: `docker network ls | grep zide-traefik-n`
- [ ] Traefik `zide` tourne et écoute sur 80/443: `docker ps | grep traefik`
- [ ] Base de données PostgreSQL accessible (ou prévoir les détails de connexion)

## ✅ Configuration des Secrets GitHub

### Clés JWT (à générer)

```bash
# Générer les clés
mkdir -p ~/slotfinder-keys && cd ~/slotfinder-keys
openssl genrsa -out private-stg.pem 2048 && openssl rsa -in private-stg.pem -pubout -out public-stg.pem
openssl genrsa -out private-prd.pem 2048 && openssl rsa -in private-prd.pem -pubout -out public-prd.pem
```

- [ ] Secret `JWT_PRIVATE_KEY_STG` créé (contenu du fichier `private-stg.pem`)
- [ ] Secret `JWT_PUBLIC_KEY_STG` créé (contenu du fichier `public-stg.pem`)
- [ ] Secret `JWT_PRIVATE_KEY_PRD` créé (contenu du fichier `private-prd.pem`)
- [ ] Secret `JWT_PUBLIC_KEY_PRD` créé (contenu du fichier `public-prd.pem`)

### Environnement Frontend

- [ ] Secret `FRONT_ENV_STG` créé avec:
  ```
  FRONT_API_URL=https://stg.slotfinder.fr/api
  ```
- [ ] Secret `FRONT_ENV_PRD` créé avec:
  ```
  FRONT_API_URL=https://slotfinder.fr/api
  ```

### Environnement Backend

- [ ] Secret `BACK_ENV_STG` créé avec au minimum:
  ```
  DB_HOST=postgres-stg
  DB_PORT=5432
  DB_USER=slotfinder
  DB_PASSWORD=your_password
  DB_NAME=slotfinder_stg
  APP_HOST=0.0.0.0
  APP_PORT=3000
  ENV=staging
  ORIGIN=https://stg.slotfinder.fr
  DOMAIN=stg.slotfinder.fr
  DB_TIMEZONE=UTC
  ```

- [ ] Secret `BACK_ENV_PRD` créé avec au minimum:
  ```
  DB_HOST=postgres-prd
  DB_PORT=5432
  DB_USER=slotfinder
  DB_PASSWORD=your_password
  DB_NAME=slotfinder_prd
  APP_HOST=0.0.0.0
  APP_PORT=3000
  ENV=production
  ORIGIN=https://slotfinder.fr
  DOMAIN=slotfinder.fr
  DB_TIMEZONE=UTC
  ```

### Vérification des Secrets

```bash
# Consultez https://github.com/Jules-Zide/SlotFinder/settings/secrets/actions
# Vous devriez voir 8 secrets au total
```

- [ ] 8 secrets visibles dans GitHub Settings → Secrets and variables → Actions

## ✅ Vérification du Code

- [ ] Fichier `back/commons/guard/jwt_guard.go` utilise base64 encoding
- [ ] Dockerfile `back/Dockerfile` ne copie plus les fichiers JWT
- [ ] Workflow `.github/workflows/build-and-deploy.yml` encode les clés en base64

```bash
# Vérifier le code
grep -n "AUTH_PRIVATE_KEY_B64" SlotFinder/back/commons/guard/jwt_guard.go
grep -n "base64" SlotFinder/.github/workflows/build-and-deploy.yml
```

## ✅ Configuration Docker Compose

- [ ] Fichier `docker-compose-stg.yml` utilise `env_file: .env` pour le backend
- [ ] Fichier `docker-compose-prd.yml` utilise `env_file: .env` pour le backend
- [ ] Les labels Traefik sont configurés correctement

```bash
# Vérifier les configurations
grep -A 2 "env_file" SlotFinder/docker-compose-stg.yml
grep -A 2 "env_file" SlotFinder/docker-compose-prd.yml
```

## ✅ Test Staging

### 1. Créer et merger une branche sur `ci/cd`

```bash
# Sur votre machine locale
git checkout -b test/deployment
git push origin test/deployment

# Créer une PR et la merger sur ci/cd
# OU directement sur ci/cd
git checkout ci/cd
git pull origin
git merge test/deployment
git push origin ci/cd
```

- [ ] Workflow `build-staging` exécuté avec succès
- [ ] Workflow `deploy-staging` exécuté avec succès
- [ ] Vérifier les logs: https://github.com/Jules-Zide/SlotFinder/actions

### 2. Tester l'accès

```bash
# Frontend
curl -v https://stg.slotfinder.fr/ 2>&1 | head -20

# Backend
curl -v https://stg.slotfinder.fr/api/v1/health 2>&1 | head -20
```

- [ ] Frontend accessible sur https://stg.slotfinder.fr/
- [ ] Backend API accessible sur https://stg.slotfinder.fr/api/v1/health
- [ ] Pas d'erreur "permission denied" dans les logs du backend

### 3. Vérifier les logs du backend

```bash
docker compose -f docker-compose-stg.yml logs slotfinder-stg-backend | tail -50
```

- [ ] Pas d'erreur "permission denied on /config/jwt/private.pem"
- [ ] Backend démarre correctement
- [ ] JWT keys sont lues depuis les variables d'env base64

## ✅ Test Production

### 1. Créer un tag de version

```bash
# Sur votre machine locale
git tag v1.0.0-test
git push origin v1.0.0-test
```

- [ ] Workflow `build-production` exécuté avec succès
- [ ] Workflow `deploy-production` exécuté avec succès
- [ ] Vérifier les logs: https://github.com/Jules-Zide/SlotFinder/actions

### 2. Tester l'accès

```bash
# Frontend
curl -v https://slotfinder.fr/ 2>&1 | head -20

# Backend
curl -v https://slotfinder.fr/api/v1/health 2>&1 | head -20
```

- [ ] Frontend accessible sur https://slotfinder.fr/
- [ ] Backend API accessible sur https://slotfinder.fr/api/v1/health
- [ ] Pas d'erreur "permission denied" dans les logs du backend

### 3. Vérifier les logs du backend

```bash
docker compose -f docker-compose-prd.yml logs slotfinder-prd-backend | tail -50
```

- [ ] Pas d'erreur "permission denied on /config/jwt/private.pem"
- [ ] Backend démarre correctement

## ✅ Endpoint Test Fonctionnels

### Test d'authentification

```bash
# Staging
curl https://stg.slotfinder.fr/api/v1/auth/status -v

# Production
curl https://slotfinder.fr/api/v1/auth/status -v
```

- [ ] Réponse 401 (non authentifié) = bon signe
- [ ] Pas d'erreur "permission denied" ou de socket timeout

### Test OAuth (Discord)

- [ ] Accéder à https://stg.slotfinder.fr/
- [ ] Cliquer sur "Login with Discord"
- [ ] Vérifier que le callback fonctionne (devrait être un 302 suivi d'une redirection)

```bash
# Vérifier les logs pour "permission denied"
docker compose -f docker-compose-stg.yml logs slotfinder-stg-backend | grep -i "permission\|error"
```

- [ ] Pas d'erreur "permission denied on /config/jwt/private.pem"

## ✅ Mise à jour du `.gitignore`

Assurez-vous que les fichiers sensibles ne sont pas commitées:

```bash
# Vérifier que .gitignore est à jour
cat SlotFinder/.gitignore | grep -E "\.env|\.pem|secret|key"
```

- [ ] `.env` est dans `.gitignore`
- [ ] `*.pem` est dans `.gitignore`
- [ ] `config/jwt/` est dans `.gitignore`

## ✅ Documenter les Changements

- [ ] Lire `.github/DEPLOYMENT.md` pour comprendre le flux complet
- [ ] Lire `.github/SETUP_SECRETS.md` pour la configuration des secrets
- [ ] Cette checklist est complétée et signée

## 🎉 Déploiement Finalisé

Une fois tout coché:

1. **Supprimez les clés locales**:
   ```bash
   rm -rf ~/slotfinder-keys
   ```

2. **Assurez-vous que les tags de test sont supprimés**:
   ```bash
   git tag -d v1.0.0-test
   git push origin :refs/tags/v1.0.0-test
   ```

3. **Vous êtes prêt pour les vrais déploiements !**

---

## 📝 Notes et Problèmes Rencontrés

### Erreur: "permission denied on /config/jwt/private.pem"

**Solution appliquée**: Les clés JWT sont maintenant lues depuis les variables d'environnement en base64 (`AUTH_PRIVATE_KEY_B64` et `AUTH_PUBLIC_KEY_B64`), au lieu de fichiers.

**Modification**:
- `back/commons/guard/jwt_guard.go`: Support du base64 avec fallback vers fichiers
- `back/Dockerfile`: Suppression de la copie des fichiers JWT
- `.github/workflows/build-and-deploy.yml`: Encode les clés en base64 avant de les passer au container

**Vérification**:
```bash
# Le backend doit démarrer sans erreur
docker compose -f docker-compose-stg.yml logs | grep -i "permission\|error\|started"
```

---

**Date de completion**: _______________

**Signé par**: _______________
