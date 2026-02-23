#!/bin/bash
set -e

if [ -S /var/run/docker.sock ]; then
    sudo chmod 666 /var/run/docker.sock
fi

exec "$@"
