# 📚 Documentation CI/CD - SlotFinder

Bienvenue dans la documentation du pipeline CI/CD de SlotFinder. Consultez les fichiers appropriés selon vos besoins.

## 🚀 Pour Commencer

**Je dois configurer le déploiement ?**
→ Lire **[SETUP_SECRETS.md](./SETUP_SECRETS.md)** en premier

**Je dois comprendre comment fonctionne le déploiement ?**
→ Lire **[DEPLOYMENT.md](./DEPLOYMENT.md)**

**Je vais faire un déploiement ?**
→ Utiliser **[DEPLOYMENT_CHECKLIST.md](./DEPLOYMENT_CHECKLIST.md)**

**Je dois comprendre la correction JWT ?**
→ Lire **[JWT_FIX_SUMMARY.md](./JWT_FIX_SUMMARY.md)**

## 📋 Documents Disponibles

| Document | Description | Público |
|----------|-------------|---------|
| **[SETUP_SECRETS.md](./SETUP_SECRETS.md)** | 🔑 Guide pour générer les clés JWT et créer les secrets GitHub | **À faire en premier** |
| **[DEPLOYMENT.md](./DEPLOYMENT.md)** | 📖 Documentation complète du pipeline CI/CD | Tous |
| **[DEPLOYMENT_CHECKLIST.md](./DEPLOYMENT_CHECKLIST.md)** | ✅ Checklist avant/pendant/après le déploiement | Ops |
| **[JWT_FIX_SUMMARY.md](./JWT_FIX_SUMMARY.md)** | 🐛 Explique la correction du problème de permissions JWT | Techs |

## 🎯 Flux Rapide

### 1️⃣ Setup Initial (Une seule fois)

```bash
# Générer les clés
mkdir -p ~/slotfinder-keys && cd ~/slotfinder-keys
openssl genrsa -out private-stg.pem 2048 && openssl rsa -in private-stg.pem -pubout -out public-stg.pem
openssl genrsa -out private-prd.pem 2048 && openssl rsa -in private-prd.pem -pubout -out public-prd.pem

# Créer les secrets GitHub
# → Consultez SETUP_SECRETS.md étape par étape
```

### 2️⃣ Déploiement Staging

```bash
# Faire des changements
git checkout -b feature/something
git commit -am "my change"
git push origin feature/something

# Merger sur ci/cd
git checkout ci/cd
git pull origin
git merge feature/something
git push origin ci/cd

# Le workflow s'exécute automatiquement !
# Vérifier: https://github.com/Jules-Zide/SlotFinder/actions
```

### 3️⃣ Déploiement Production

```bash
# Créer un tag
git tag v1.0.0
git push origin v1.0.0

# Le workflow s'exécute automatiquement !
# Vérifier: https://github.com/Jules-Zide/SlotFinder/actions
```

## 🔍 Vérification Rapide

```bash
# Vérifier les secrets créés
# https://github.com/Jules-Zide/SlotFinder/settings/secrets/actions

# Vérifier les workflows
# https://github.com/Jules-Zide/SlotFinder/actions

# Vérifier le code
grep "AUTH_PRIVATE_KEY_B64" back/commons/guard/jwt_guard.go
grep "base64" .github/workflows/build-and-deploy.yml

# Vérifier l'application
curl https://stg.slotfinder.fr/
curl https://slotfinder.fr/
```

## 📊 Architecture

```
┌─────────────────────────────────────┐
│     Developer (Local Machine)       │
│  git push origin ci/cd              │
│  (or git tag v1.0.0)                │
└────────────────┬────────────────────┘
                 │
                 ▼
┌─────────────────────────────────────┐
│    GitHub Actions Workflow           │
│  - Generate secrets from env vars    │
│  - Encode JWT keys to base64         │
│  - Build Docker images               │
│  - Deploy via docker-compose         │
└────────────────┬────────────────────┘
                 │
                 ▼
┌─────────────────────────────────────┐
│    Self-Hosted Docker Runner        │
│  - Frontend (Nginx)                 │
│  - Backend (Go API)                 │
│  - PostgreSQL (optional)             │
└────────────────┬────────────────────┘
                 │
                 ▼
┌─────────────────────────────────────┐
│      External Traefik Router        │
│  - SSL/TLS (Let's Encrypt)          │
│  - Domain routing                    │
│  - Load balancing                    │
└────────────────┬────────────────────┘
                 │
                 ▼
┌─────────────────────────────────────┐
│          Public Internet             │
│  https://stg.slotfinder.fr/         │
│  https://slotfinder.fr/             │
└─────────────────────────────────────┘
```

## 🔒 Sécurité

✅ Bonnes pratiques appliquées:
- Secrets stockés dans GitHub (chiffrés)
- Clés JWT jamais en plain text dans les logs
- Variables d'environnement chiffées via secrets
- Pas de fichiers sensibles dans le repository
- Fallback automatique pour compatibilité

## 🐛 Problèmes Fréquents

### "permission denied on /config/jwt/private.pem"

**Résolu !** Les clés JWT sont maintenant passées via variables d'env en base64, pas via fichiers.

Pour vérifier:
```bash
docker compose -f docker-compose-stg.yml logs | grep -i "permission"
# Aucun résultat = tout va bien ✅
```

### "JWT secret not found"

Les secrets n'ont pas été créés. Consultez **[SETUP_SECRETS.md](./SETUP_SECRETS.md)**.

### "Workflow won't start"

1. Vérifiez que le runner self-hosted est en ligne
2. Vérifiez les logs: Settings → Actions → Runners
3. Vérifiez la branche push ou le tag push

## 📞 Support

| Question | Réponse |
|----------|--------|
| Comment générer les clés JWT ? | [SETUP_SECRETS.md](./SETUP_SECRETS.md#étape-1-générer-les-clés-jwt) |
| Comment créer les secrets ? | [SETUP_SECRETS.md](./SETUP_SECRETS.md#étape-2-créer-les-secrets-dans-github) |
| Comment déployer en staging ? | [DEPLOYMENT.md](./DEPLOYMENT.md#déploiement-staging) |
| Comment déployer en prod ? | [DEPLOYMENT.md](./DEPLOYMENT.md#déploiement-production) |
| Comment tester avant déploiement ? | [DEPLOYMENT_CHECKLIST.md](./DEPLOYMENT_CHECKLIST.md) |
| Pourquoi base64 pour les clés ? | [JWT_FIX_SUMMARY.md](./JWT_FIX_SUMMARY.md) |

## 🎓 Ressources Externes

- [GitHub Actions Documentation](https://docs.github.com/en/actions)
- [Docker Compose Reference](https://docs.docker.com/compose/compose-file/)
- [Traefik Documentation](https://doc.traefik.io/)
- [JWT Best Practices](https://tools.ietf.org/html/rfc8949)

---

**Dernière mise à jour**: 2024
**Version**: 1.0
**Status**: ✅ Production Ready
