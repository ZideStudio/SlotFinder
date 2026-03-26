#!/bin/bash
set -e

if [ -S /var/run/docker.sock ]; then
    sudo chmod 666 /var/run/docker.sock
fi

# Recover from stale Xvfb lock after an unclean previous shutdown.
LOCK_FILE="/tmp/.X99-lock"
SOCKET_FILE="/tmp/.X11-unix/X99"
if [ -f "$LOCK_FILE" ]; then
    LOCK_PID="$(cat "$LOCK_FILE" 2>/dev/null || true)"
    if [ -z "$LOCK_PID" ] || ! kill -0 "$LOCK_PID" 2>/dev/null; then
        rm -f "$LOCK_FILE" "$SOCKET_FILE"
    fi
fi

# Start virtual display for browser tooling. If it is already running, continue.
Xvfb :99 -screen 0 1280x1024x24 -nolisten tcp >/tmp/xvfb.log 2>&1 || true

exec "$@"
