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

if ! command -v rtk >/dev/null 2>&1; then
  warn "RTK is not installed; continuing without automatic Copilot hook setup."
else
  info "Initializing RTK..."
  rtk init --copilot --auto-patch >/tmp/rtk-init.log 2>&1
  status=$?

  if [ "$status" -ne 0 ]; then
    warn "RTK initialization failed; continuing without automatic Copilot hook setup."
    if [ -s /tmp/rtk-init.log ]; then
      cat /tmp/rtk-init.log
    fi
  else
    success "RTK initialized successfully."
  fi
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
