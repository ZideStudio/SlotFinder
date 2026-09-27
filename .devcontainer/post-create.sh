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

if command -v go >/dev/null 2>&1; then
  info "Installing Go developer tools..."
  if (
    cd /workspace &&
      go install golang.org/x/tools/gopls@v0.21.1 &&
      go install github.com/go-delve/delve/cmd/dlv@v1.26.0 &&
      go install golang.org/x/tools/cmd/goimports@v0.44.0
  ); then
    success "Go developer tools installed."
  else
    warn "Failed to install Go developer tools."
  fi
else
  warn "The go command is not available; Go developer tools installation skipped."
fi

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
