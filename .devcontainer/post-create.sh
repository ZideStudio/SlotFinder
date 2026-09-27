#!/usr/bin/env bash
set -u

log() {
  echo "[$(date +'%H:%M:%S')] $1"
}

warn() {
  log "⚠️  $1"
}

info() {
  log "ℹ️  $1"
}

success() {
  log "✅ $1"
}

info "Starting post-create script..."

if [ -d "/workspace/front" ]; then
  info "Checking frontend dependencies..."
  if [ -f "/workspace/front/package.json" ]; then
    cd /workspace/front

    if [ ! -d "node_modules" ]; then
      info "Installing frontend npm dependencies..."
      npm ci
    else
      info "Frontend npm dependencies are already installed."
    fi

    info "Installing Playwright browsers..."
    npx playwright install --with-deps
    success "Playwright Chromium is ready for browser tests."
  else
    warn "The front folder does not contain a package.json; Playwright installation skipped."
  fi
else
  warn "The front folder is missing; Playwright installation skipped."
fi

success "Post-create script completed successfully."
exit 0
