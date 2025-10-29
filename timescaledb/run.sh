#!/usr/bin/env bash
set -e

# Environment defaults
POSTGRES_USER=${POSTGRES_USER:-homeassistant}
POSTGRES_PASSWORD=${POSTGRES_PASSWORD:-secret}
POSTGRES_DB=${POSTGRES_DB:-homeassistant}

echo "Starting TimescaleDB..."
echo "User: $POSTGRES_USER, DB: $POSTGRES_DB"

# Initialize DB if first run
if [ ! -d "/var/lib/postgresql/data" ]; then
  mkdir -p /var/lib/postgresql/data
  chown -R postgres:postgres /var/lib/postgresql
fi

# Start Postgres
exec docker-entrypoint.sh postgres
