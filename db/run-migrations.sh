#!/bin/bash
set -e

echo "Waiting for Postgres to be ready..."
sleep 5  # optionally increase if needed

for f in /migrations/*.sql; do
  echo "Running migration: $f"
  psql "postgresql://${DB_USER}:${DB_PASS}@${DB_HOST}:${DB_PORT}/${DB_NAME}" -f "$f"
done

echo "All migrations executed successfully."
