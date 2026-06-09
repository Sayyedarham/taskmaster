#!/bin/bash
set -e

echo "Stopping migrate container..."
docker-compose stop migrate

echo "Resetting migration state..."
docker-compose exec -T postgres psql -U taskmaster -d taskmaster -c "UPDATE schema_migrations SET version = 1, dirty = false;"

echo "Re-running migrations..."
docker-compose up migrate

echo "✅ Done!"
