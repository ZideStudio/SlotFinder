#!/bin/sh
set -e

# Run npm install to ensure dependencies are up to date
echo "Running npm install..."
npm install

# Generate i18next resources for TypeScript
echo "Generating i18next resources for TypeScript..."
npm run i18next-resources-for-ts

# Execute the command passed to the container
exec "$@"
