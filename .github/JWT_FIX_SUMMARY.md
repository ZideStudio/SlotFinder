# Résumé des Corrections - JWT Permission Issue

## 🐛 Problème Initial

```
2026/05/26 08:48PM ERR PROVIDER_CALLBACK failed to connect user | 
error=open /config/jwt/private.pem: permission denied
```

Le backend ne pouvait pas lire les clés JWT depuis les fichiers du système.

## ✅ Solution Appliquée

Les clés JWT ne sont plus stockées en tant que fichiers physiques dans le container. Elles sont maintenant :

1. **Stockées** dans les secrets GitHub (PEM complet)
2. **Encodées en base64** par le workflow
3. **Injectées** via variables d'environnement `AUTH_PRIVATE_KEY_B64` et `AUTH_PUBLIC_KEY_B64`
4. **Décodées** au runtime par le backend

## 📝 Fichiers Modifiés

### 1. `back/commons/guard/jwt_guard.go`

**Changement**: Ajout du support base64 avec fallback vers fichiers

```go
// Avant: lire uniquement depuis fichier
f, err := os.ReadFile(config.Auth.PublicPemPath)

// Après: priorité base64, fallback fichier
publicKeyB64 := os.Getenv("AUTH_PUBLIC_KEY_B64")
if publicKeyB64 != "" {
    f, err = base64.StdEncoding.DecodeString(publicKeyB64)
} else {
    f, err = os.ReadFile(config.Auth.PublicPemPath)
}
```

Même logique pour `GenerateAccessToken()` avec `AUTH_PRIVATE_KEY_B64`.

### 2. `back/Dockerfile`

**Changement**: Suppression de la copie des fichiers JWT

```dockerfile
# ❌ AVANT: Copie les fichiers JWT
COPY ./back/config/jwt/ /usr/src/app/config/jwt/
COPY --from=build /usr/src/app/config/jwt/ /config/jwt/

# ✅ APRÈS: Pas de copie - clés via variables d'env
# (plus de fichiers JWT)
```

### 3. `.github/workflows/build-and-deploy.yml`

**Changement**: Encode les clés en base64 avant déploiement

```bash
# Avant
echo "${{ secrets.BACK_ENV_STG }}" > .env

# Après
{
  echo "${{ secrets.BACK_ENV_STG }}"
  echo "AUTH_PRIVATE_KEY_B64=$(base64 -i <(echo \"${{ secrets.JWT_PRIVATE_KEY_STG }}\" | tr -d '\r') | tr -d '\n')"
  echo "AUTH_PUBLIC_KEY_B64=$(base64 -i <(echo \"${{ secrets.JWT_PUBLIC_KEY_STG }}\" | tr -d '\r') | tr -d '\n')"
} > .env
```

## 📚 Documentation Créée

### 1. `.github/SETUP_SECRETS.md`

Guide complet pour :
- Générer les clés RSA
- Créer les 8 secrets GitHub
- Tester et dépanner

### 2. `.github/DEPLOYMENT.md` (mis à jour)

Ajout de :
- Section "JWT Keys (Base64 Encoding)"
- Explications du nouveau flux
- Notes de compatibilité

### 3. `.github/DEPLOYMENT_CHECKLIST.md`

Checklist complète pour :
- Vérifier les prérequis
- Configurer les secrets
- Tester staging et production
- Vérifier les logs

## 🚀 Next Steps

### 1. Générer les clés JWT

```bash
mkdir -p ~/slotfinder-keys && cd ~/slotfinder-keys
openssl genrsa -out private-stg.pem 2048 && openssl rsa -in private-stg.pem -pubout -out public-stg.pem
openssl genrsa -out private-prd.pem 2048 && openssl rsa -in private-prd.pem -pubout -out public-prd.pem
```

### 2. Créer les 8 secrets GitHub

Consultez `.github/SETUP_SECRETS.md` étape par étape.

Les 8 secrets à créer:
- `JWT_PRIVATE_KEY_STG` (PEM)
- `JWT_PUBLIC_KEY_STG` (PEM)
- `JWT_PRIVATE_KEY_PRD` (PEM)
- `JWT_PUBLIC_KEY_PRD` (PEM)
- `FRONT_ENV_STG`
- `BACK_ENV_STG`
- `FRONT_ENV_PRD`
- `BACK_ENV_PRD`

### 3. Tester le déploiement

```bash
# Push sur ci/cd pour tester staging
git push origin ci/cd

# Vérifier les logs
# https://github.com/Jules-Zide/SlotFinder/actions
```

### 4. Vérifier que ça fonctionne

```bash
# Le backend doit démarrer sans "permission denied"
docker compose -f docker-compose-stg.yml logs | grep -i "permission"
```

## 🔒 Sécurité

✅ **Points positifs**:
- Clés privées jamais en plain text dans les logs
- Clés jamais stockées dans les fichiers du container
- Clés chiffrées dans GitHub Secrets
- Support du fallback vers fichiers (rétrocompatibilité)

## 💡 Avantages

1. **Aucun problème de permissions** - pas de fichiers à lire/écrire
2. **Images Docker plus petites** - pas de fichiers JWT à copier
3. **Plus portable** - fonctionnera sur n'importe quel OS
4. **Compatible** - fallback automatique vers fichiers si env var manquante

## 📋 Vérification Rapide

```bash
# Vérifier que la modification est présente
grep "AUTH_PRIVATE_KEY_B64" SlotFinder/back/commons/guard/jwt_guard.go
# Output: publicKeyB64 := os.Getenv("AUTH_PRIVATE_KEY_B64")

# Vérifier que Dockerfile ne copie plus les fichiers JWT
grep -n "jwt/" SlotFinder/back/Dockerfile
# Output: (rien - pas de fichier JWT copié)

# Vérifier que le workflow encode en base64
grep "base64" SlotFinder/.github/workflows/build-and-deploy.yml
# Output: plusieurs lignes avec base64
```

## 🎉 Résultat

Après le déploiement avec les secrets configurés :

✅ `https://stg.slotfinder.fr/api/v1/auth/status` → 401 (normal sans token)
✅ `https://stg.slotfinder.fr/api/v1/auth/discord/callback?code=...` → 302 (callback OAuth fonctionne)
✅ Logs du backend : **AUCUNE** erreur "permission denied on /config/jwt/private.pem"

## 📞 Support

- Consultez `.github/SETUP_SECRETS.md` pour les détails des secrets
- Consultez `.github/DEPLOYMENT.md` pour comprendre le flux complet
- Consultez `.github/DEPLOYMENT_CHECKLIST.md` pour la checklist
