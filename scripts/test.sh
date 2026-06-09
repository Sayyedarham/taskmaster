#!/bin/bash
set -e

echo "🧪 Running tests..."

cd backend

echo "Running unit tests..."
go test ./tests/unit/... -v -race -cover

echo ""
echo "Running integration tests..."
go test ./tests/integration/... -v

echo ""
echo "✅ All tests passed"
