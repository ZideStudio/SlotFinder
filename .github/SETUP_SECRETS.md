# Configuration des Secrets GitHub - SlotFinder

Ce guide vous aide à créer tous les secrets nécessaires pour le déploiement CI/CD.

## Prérequis

- Accès administrateur au repository GitHub SlotFinder
- OpenSSL installé (pour générer les clés JWT)
- Terminal/Shell

## Étape 1: Générer les clés JWT

### Pour le Staging

```bash
# Créer un répertoire temporaire pour les clés
mkdir -p ~/slotfinder-keys
cd ~/slotfinder-keys

# Générer la clé privée (2048 bits, standard pour JWT/RS256)
openssl genrsa -out private-stg.pem 2048

# Extraire la clé publique
openssl rsa -in private-stg.pem -pubout -out public-stg.pem

# Vérifier les clés (optionnel)
openssl rsa -in private-stg.pem -check
```

### Pour la Production

```bash
# Générer la clé privée
openssl genrsa -out private-prd.pem 2048

# Extraire la clé publique
openssl rsa -in private-prd.pem -pubout -out public-prd.pem

# Vérifier les clés (optionnel)
openssl rsa -in private-prd.pem -check
```

Les fichiers générés ressembleront à:

```
-----BEGIN RSA PRIVATE KEY-----
MIIEpAIBAAKCAQEA1234567890...
... (many more lines)
-----END RSA PRIVATE KEY-----
```

```
-----BEGIN PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA1234567890...
... (many more lines)
-----END PUBLIC KEY-----
```

## Étape 2: Créer les secrets dans GitHub

### Accéder aux secrets

1. Allez à: https://github.com/Jules-Zide/SlotFinder/settings/secrets/actions
2. Cliquez sur "New repository secret"

### Ajouter les secrets JWT

Créez les 4 secrets JWT suivants:

#### 1. `JWT_PRIVATE_KEY_STG`

```bash
# Afficher le contenu à copier
cat ~/slotfinder-keys/private-stg.pem
```

- **Name**: `JWT_PRIVATE_KEY_STG`
- **Value**: Copiez le contenu **entier** du fichier `private-stg.pem` (BEGIN et END lines incluses)

#### 2. `JWT_PUBLIC_KEY_STG`

```bash
# Afficher le contenu à copier
cat ~/slotfinder-keys/public-stg.pem
```

- **Name**: `JWT_PUBLIC_KEY_STG`
- **Value**: Copiez le contenu **entier** du fichier `public-stg.pem`

#### 3. `JWT_PRIVATE_KEY_PRD`

```bash
# Afficher le contenu à copier
cat ~/slotfinder-keys/private-prd.pem
```

- **Name**: `JWT_PRIVATE_KEY_PRD`
- **Value**: Copiez le contenu **entier** du fichier `private-prd.pem`

#### 4. `JWT_PUBLIC_KEY_PRD`

```bash
# Afficher le contenu à copier
cat ~/slotfinder-keys/public-prd.pem
```

- **Name**: `JWT_PUBLIC_KEY_PRD`
- **Value**: Copiez le contenu **entier** du fichier `public-prd.pem`

### Ajouter les secrets d'environnement

#### 5. `FRONT_ENV_STG`

- **Name**: `FRONT_ENV_STG`
- **Value**:
```
FRONT_API_URL=https://stg.slotfinder.fr/api
FRONT_DEBUG=true
```

Ajustez selon vos besoins (autres variables FRONT_*).

#### 6. `BACK_ENV_STG`

- **Name**: `BACK_ENV_STG`
- **Value**:
```
DB_HOST=postgres-stg
DB_PORT=5432
DB_USER=slotfinder
DB_PASSWORD=your_secret_password
DB_NAME=slotfinder_stg
APP_HOST=0.0.0.0
APP_PORT=3000
ENV=staging
IMGBB_API_KEY=your_imgbb_key
ORIGIN=https://stg.slotfinder.fr
DOMAIN=stg.slotfinder.fr
DB_TIMEZONE=UTC
```

#### 7. `FRONT_ENV_PRD`

- **Name**: `FRONT_ENV_PRD`
- **Value**:
```
FRONT_API_URL=https://slotfinder.fr/api
FRONT_DEBUG=false
```

#### 8. `BACK_ENV_PRD`

- **Name**: `BACK_ENV_PRD`
- **Value**:
```
DB_HOST=postgres-prd
DB_PORT=5432
DB_USER=slotfinder
DB_PASSWORD=your_secret_password
DB_NAME=slotfinder_prd
APP_HOST=0.0.0.0
APP_PORT=3000
ENV=production
IMGBB_API_KEY=your_imgbb_key
ORIGIN=https://slotfinder.fr
DOMAIN=slotfinder.fr
DB_TIMEZONE=UTC
```

## Étape 3: Vérifier les secrets

```bash
# Lister les secrets créés (à faire dans l'interface GitHub)
# Settings → Secrets and variables → Actions

# Vous devriez voir:
# ✅ JWT_PRIVATE_KEY_STG
# ✅ JWT_PUBLIC_KEY_STG
# ✅ JWT_PRIVATE_KEY_PRD
# ✅ JWT_PUBLIC_KEY_PRD
# ✅ FRONT_ENV_STG
# ✅ BACK_ENV_STG
# ✅ FRONT_ENV_PRD
# ✅ BACK_ENV_PRD
```

## Étape 4: Nettoyer les clés locales

Une fois les secrets créés dans GitHub, supprimez les fichiers locaux:

```bash
# Supprimer le répertoire temporaire
rm -rf ~/slotfinder-keys

# Vérifier qu'il n'y a pas de trace
ls ~/slotfinder-keys  # Ne doit pas exister
```

## Étape 5: Tester le déploiement

### Test Staging

```bash
# Créer une branche de test et faire un push
git checkout -b test/deployment
echo "test" >> README.md
git add README.md
git commit -m "test deployment"
git push origin test/deployment

# Créer une PR et merger sur la branche ci/cd
# Le workflow devrait s'exécuter automatiquement
```

Vérifiez les logs du workflow:
- Allez à: https://github.com/Jules-Zide/SlotFinder/actions
- Cliquez sur le workflow en cours
- Vérifiez l'étape "Create .env for backend staging"

### Test Production

```bash
# Créer un tag et le pousser
git tag v1.0.0-test
git push origin v1.0.0-test

# Vérifier les logs du workflow
# https://github.com/Jules-Zide/SlotFinder/actions
```

## Dépannage

### Erreur: "permission denied on /config/jwt/private.pem"

Cette erreur a été **corrigée** - les clés sont maintenant lues depuis les variables d'environnement en base64, pas depuis des fichiers.

Si vous la voyez encore:
1. Vérifiez que `JWT_PRIVATE_KEY_STG` et `JWT_PUBLIC_KEY_STG` sont dans les secrets
2. Vérifiez que les clés sont correctement formatées (BEGIN/END lines incluses)
3. Relancez le déploiement

### Erreur: "failed to decode base64"

Cela signifie que les clés n'ont pas bien été encodées. Vérifiez:
1. Le contenu des secrets JWT - ils doivent contenir des clés PEM complètes
2. Pas de caractères invisibles ou espaces supplémentaires

### Erreur: "AUTH_PRIVATE_KEY_B64 not found"

C'est normal - le code a un fallback vers les fichiers. Si vous continuez à avoir "permission denied", cela signifie que:
1. La variable n'est pas passée au container
2. Le fichier `/config/jwt/private.pem` n'existe pas

Assurez-vous que les variables sont bien dans le `.env` du deployment.

## Rotation des clés

Pour changer les clés (rotation de sécurité):

1. Générez de nouvelles clés:
   ```bash
   openssl genrsa -out private-stg-new.pem 2048
   openssl rsa -in private-stg-new.pem -pubout -out public-stg-new.pem
   ```

2. Mettez à jour les secrets GitHub avec les nouvelles clés

3. Relancez le déploiement:
   ```bash
   git push origin ci/cd  # Force un redéploiement
   ```

Les tokens existants générés avec les anciennes clés resteront valides jusqu'à expiration. Les nouveaux tokens utiliseront les nouvelles clés.

## Sécurité

⚠️ **Points importants:**

- 🔒 Les clés privées ne doivent **JAMAIS** être commitées dans le repository
- 🔒 Les clés privées ne doivent **JAMAIS** être en plain text dans les logs
- 🔒 Utilisez GitHub Secrets (HTTPS, chiffré) pour stocker les clés
- 🔒 Accès limité aux secrets - seulement les administrateurs
- 🔒 Rotation régulière des clés (au moins une fois par an)
- 🔒 Supprimez les copies locales des clés après création

## Support

Si vous avez des problèmes:

1. Consultez `.github/DEPLOYMENT.md` pour plus de détails
2. Vérifiez les logs du workflow: https://github.com/Jules-Zide/SlotFinder/actions
3. Vérifiez les logs du container: `docker compose logs backend`
