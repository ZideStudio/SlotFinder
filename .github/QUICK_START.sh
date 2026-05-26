#!/bin/bash

# SlotFinder JWT Fix - Quick Start Script
# This script helps you set up the JWT keys for SlotFinder deployment

set -e

echo "╔════════════════════════════════════════════════════════════════╗"
echo "║       SlotFinder JWT Fix - Quick Start Setup                  ║"
echo "╚════════════════════════════════════════════════════════════════╝"
echo ""

# Step 1: Create JWT keys directory
echo "📁 STEP 1: Creating JWT keys directory..."
KEYS_DIR="${HOME}/slotfinder-keys"
mkdir -p "$KEYS_DIR"
cd "$KEYS_DIR"
echo "✅ Directory created: $KEYS_DIR"
echo ""

# Step 2: Generate Staging Keys
echo "🔑 STEP 2: Generating Staging keys..."
openssl genrsa -out private-stg.pem 2048 2>/dev/null
openssl rsa -in private-stg.pem -pubout -out public-stg.pem 2>/dev/null
echo "✅ Staging keys generated"
echo "   - private-stg.pem"
echo "   - public-stg.pem"
echo ""

# Step 3: Generate Production Keys
echo "🔑 STEP 3: Generating Production keys..."
openssl genrsa -out private-prd.pem 2048 2>/dev/null
openssl rsa -in private-prd.pem -pubout -out public-prd.pem 2>/dev/null
echo "✅ Production keys generated"
echo "   - private-prd.pem"
echo "   - public-prd.pem"
echo ""

# Step 4: Display keys for copying
echo "📋 STEP 4: Key Contents (for GitHub Secrets)"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""

echo "🔐 JWT_PRIVATE_KEY_STG:"
echo "────────────────────────"
cat private-stg.pem
echo ""
echo ""

echo "🔐 JWT_PUBLIC_KEY_STG:"
echo "────────────────────"
cat public-stg.pem
echo ""
echo ""

echo "🔐 JWT_PRIVATE_KEY_PRD:"
echo "────────────────────────"
cat private-prd.pem
echo ""
echo ""

echo "🔐 JWT_PUBLIC_KEY_PRD:"
echo "────────────────────"
cat public-prd.pem
echo ""
echo ""

# Step 5: Instructions
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "📝 NEXT STEPS:"
echo ""
echo "1️⃣  Copy each key content above"
echo ""
echo "2️⃣  Create 4 GitHub Secrets:"
echo "    Go to: https://github.com/Jules-Zide/SlotFinder/settings/secrets/actions"
echo ""
echo "    Create these secrets:"
echo "    • JWT_PRIVATE_KEY_STG (copy content above)"
echo "    • JWT_PUBLIC_KEY_STG  (copy content above)"
echo "    • JWT_PRIVATE_KEY_PRD (copy content above)"
echo "    • JWT_PUBLIC_KEY_PRD  (copy content above)"
echo ""
echo "3️⃣  Create 4 Environment Secrets:"
echo "    • FRONT_ENV_STG   (see .github/SETUP_SECRETS.md)"
echo "    • BACK_ENV_STG    (see .github/SETUP_SECRETS.md)"
echo "    • FRONT_ENV_PRD   (see .github/SETUP_SECRETS.md)"
echo "    • BACK_ENV_PRD    (see .github/SETUP_SECRETS.md)"
echo ""
echo "4️⃣  Clean up local keys:"
echo "    rm -rf $KEYS_DIR"
echo ""
echo "5️⃣  Test the deployment:"
echo "    git push origin ci/cd"
echo ""
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""
echo "📖 Full documentation: .github/SETUP_SECRETS.md"
echo "✅ Deployment checklist: .github/DEPLOYMENT_CHECKLIST.md"
echo ""
