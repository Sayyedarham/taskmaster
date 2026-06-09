#!/bin/bash
set -e

echo "🚀 TaskMaster Setup Script"
echo "==========================="

# Check prerequisites
echo "Checking prerequisites..."
command -v docker >/dev/null 2>&1 || { echo "❌ Docker is required but not installed. Aborting." >&2; exit 1; }
command -v docker-compose >/dev/null 2>&1 || command -v docker compose >/dev/null 2>&1 || { echo "❌ Docker Compose is required but not installed. Aborting." >&2; exit 1; }

echo "✅ Docker and Docker Compose found"

# Create .env from example if it doesn't exist
if [ ! -f .env ]; then
    echo "Creating .env from example..."
    cp .env.example .env
    echo "✅ .env created. Edit it if you need custom configuration."
fi

# Determine docker compose command
if command -v docker-compose >/dev/null 2>&1; then
    COMPOSE="docker-compose"
else
    COMPOSE="docker compose"
fi

# Stop any existing containers to avoid conflicts
echo "Stopping any existing containers..."
$COMPOSE down -v 2>/dev/null || true

# Start infrastructure
echo "Starting PostgreSQL and Redis..."
$COMPOSE up -d postgres redis

# Wait for PostgreSQL to be ready
echo "Waiting for PostgreSQL to be ready..."
RETRIES=30
until $COMPOSE exec -T postgres pg_isready -U taskmaster > /dev/null 2>&1; do
    RETRIES=$((RETRIES - 1))
    if [ $RETRIES -eq 0 ]; then
        echo "❌ PostgreSQL failed to start after 30 attempts"
        exit 1
    fi
    echo "  PostgreSQL not ready yet, retrying... ($RETRIES left)"
    sleep 2
done
echo "✅ PostgreSQL is ready"

# Run migrations
echo "Running database migrations..."
$COMPOSE up migrate

# Check if migration succeeded
if [ $? -ne 0 ]; then
    echo "⚠️ Migration may have issues, but continuing..."
fi

# Seed data
echo "Seeding demo data..."
if [ -f scripts/seed.sh ]; then
    bash scripts/seed.sh
else
    echo "⚠️ Seed script not found, skipping seed data"
fi

# Start API and Frontend
echo "Starting API and Frontend..."
$COMPOSE up -d api frontend

echo ""
echo "✅ Setup complete!"
echo ""
echo "📍 Access points:"
echo "   API:      http://localhost:8080"
echo "   Frontend: http://localhost:3000"
echo "   Health:   http://localhost:8080/health"
echo ""
echo "🔑 Demo credentials:"
echo "   Admin:   admin@taskmaster.com / password123"
echo "   Manager: manager@taskmaster.com / password123"
echo "   Member:  member@taskmaster.com / password123"
echo ""
echo "📖 Next steps:"
echo "   - View logs:     $COMPOSE logs -f api"
echo "   - Run tests:     make test (requires Go installed)"
echo "   - Stop stack:    $COMPOSE down"
echo ""
