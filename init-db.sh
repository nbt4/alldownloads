#!/bin/bash
set -e

# Wait for PostgreSQL to be ready
until pg_isready -h localhost -p 5432 -U alldl; do
  echo "Waiting for database to be ready..."
  sleep 2
done

echo "Database is ready. Running migrations..."

# Install golang-migrate if not present
if ! command -v migrate &> /dev/null; then
    echo "Installing golang-migrate..."
    apk add --no-cache curl
    curl -L https://github.com/golang-migrate/migrate/releases/latest/download/migrate.linux-amd64.tar.gz | tar xvz
    mv migrate /usr/local/bin/
fi

# Run migrations if they exist
if [ -d "/docker-entrypoint-initdb.d/migrations" ]; then
    echo "Running database migrations..."
    migrate -path /docker-entrypoint-initdb.d/migrations -database "postgres://alldl:alldl@localhost:5432/alldownloads?sslmode=disable" up
    echo "Migrations completed successfully."
else
    echo "No migrations directory found, skipping migrations."
fi