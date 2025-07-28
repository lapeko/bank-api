#!/bin/sh
set -e

echo "🚀 Running migrations..."
migrate -source file://app/internal/db/migration -database "$POSTGRES_URL" -verbose up

echo "✅ Starting API..."
exec bank-api