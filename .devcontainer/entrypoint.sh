#!/bin/bash
set -e

if [ -S /var/run/docker.sock ]; then
    sudo chmod 666 /var/run/docker.sock
fi

Xvfb :99 -screen 0 1280x1024x24

exec "$@"
