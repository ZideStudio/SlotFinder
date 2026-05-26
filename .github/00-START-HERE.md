# 🚀 START HERE - SlotFinder JWT Fix

Bienvenue ! Ce fichier vous guide dans les étapes suivantes après la correction JWT.

## ⚡ Résumé Rapide (30 secondes)

**Problème**: `error=open /config/jwt/private.pem: permission denied`

**Solution**: Les clés JWT passent maintenant via des variables d'env en base64

**Quoi faire**: 
1. Générer les clés
2. Créer les secrets GitHub
3. Tester le déploiement

## 📋 Checklist Rapide

```bash
# 1️⃣  Générer les clés (auto)
bash .github/QUICK_START.sh

# 2️⃣  Créer les secrets GitHub (manual)
# Allez à: https://github.com/Jules-Zide/SlotFinder/settings/secrets/actions
# Créez 8 secrets (voir SETUP_SECRETS.md)

# 3️⃣  Tester le staging
git push origin ci/cd

# 4️⃣  Vérifier les logs
# https://github.com/Jules-Zide/SlotFinder/actions
# Doit afficher "Server started" SANS erreur "permission denied"

# 5️⃣  Tester production
git tag v1.0.0
git push origin v1.0.0
```

## 📚 Documentation Complète

| Document | Objectif | Temps |
|----------|----------|-------|
| **README.md** | Vue d'ensemble et navigation | 5 min |
| **SETUP_SECRETS.md** | Générer clés et créer secrets | 15 min |
| **DEPLOYMENT_CHECKLIST.md** | Checklist complète avant déploiement | 10 min |
| **JWT_FIX_SUMMARY.md** | Explications techniques | 10 min |
| **VALIDATION.md** | Valider l'implémentation | 10 min |
| **DEPLOYMENT.md** | Documentation complète du CD | 20 min |

## 🎯 Next Steps

### 1. Étape 1: Générer les clés (5 minutes)

```bash
bash .github/QUICK_START.sh
```

Cela créera 4 fichiers `.pem` et vous montrera le contenu à copier.

### 2. Étape 2: Créer les secrets (10 minutes)

Consultez **SETUP_SECRETS.md** pour les instructions détaillées.

Vous devez créer 8 secrets GitHub:
- 4 secrets JWT (clés PEM)
- 4 secrets d'environnement (config backend/frontend)

### 3. Étape 3: Tester (10 minutes)

```bash
git push origin ci/cd
```

Puis vérifiez:
- https://github.com/Jules-Zide/SlotFinder/actions
- Logs du backend: pas d'erreur "permission denied"
- https://stg.slotfinder.fr/ doit être accessible

## ❓ Questions Fréquentes

**Q: Quel fichier dois-je lire en premier ?**
A: README.md pour la vue d'ensemble, puis SETUP_SECRETS.md pour commencer.

**Q: Je n'ai pas de clés JWT, que faire ?**
A: Exécutez `bash .github/QUICK_START.sh` pour les générer automatiquement.

**Q: Où je peux voir les logs du déploiement ?**
A: https://github.com/Jules-Zide/SlotFinder/actions

**Q: Quelle est la différence avant/après ?**
A: Consultez JWT_FIX_SUMMARY.md pour la comparaison détaillée.

## 🔒 Sécurité

⚠️ **Important:**
- Ne JAMAIS committer les fichiers `.pem`
- Supprimez les clés locales après création des secrets
- Les secrets GitHub sont chiffrés

## 📞 Support

| Besoin | Fichier |
|--------|---------|
| Général | README.md |
| Setup | SETUP_SECRETS.md |
| Déploiement | DEPLOYMENT_CHECKLIST.md |
| Technique | JWT_FIX_SUMMARY.md |
| Validation | VALIDATION.md |

## ✅ Vérification Finale

Avant de déployer en production, vérifiez:

```bash
# 1. Les fichiers ont été modifiés
grep "AUTH_PRIVATE_KEY_B64" back/commons/guard/jwt_guard.go
# ✅ Doit afficher la variable

# 2. Dockerfile ne copie plus les fichiers JWT
grep "jwt/" back/Dockerfile | wc -l
# ✅ Doit afficher 0

# 3. Workflow encode les clés en base64
grep "base64" .github/workflows/build-and-deploy.yml | wc -l
# ✅ Doit afficher un nombre > 0
```

## 🚀 Prêt ?

Une fois le checklist complété:
1. Les clés JWT sont générées ✅
2. Les secrets GitHub sont créés ✅
3. Le staging déploie sans erreur ✅
4. Les logs ne montrent aucune erreur "permission denied" ✅

Vous pouvez déployer en production !

---

**Commencez par:** `.github/SETUP_SECRETS.md`

**Questions ?** Consultez `.github/README.md`

**Prêt à déployer ?** Consultez `.github/DEPLOYMENT_CHECKLIST.md`
