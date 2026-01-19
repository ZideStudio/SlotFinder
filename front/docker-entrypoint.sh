#!/bin/sh
set -e

# Generate i18next resources for TypeScript
echo "Generating i18next resources for TypeScript..."
npm run i18next-resources-for-ts

# Execute the command passed to the container
exec "$@"
