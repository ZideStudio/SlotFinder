#!/bin/sh
set -e

# Run npm install to ensure dependencies are up to date
echo "Running npm install..."
npm install

# Execute the command passed to the container
exec "$@"
