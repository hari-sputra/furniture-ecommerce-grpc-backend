#!/bin/sh

set -e

# Wait for the database to be ready
echo "Waiting for database..."
while ! nc -z ${DB_HOST} ${DB_PORT}; do
  sleep 1
done
echo "Database is ready."

# Run database migrations
/app/migrate -path /app/pkg/database/migrations -database "postgres://${DB_USER}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=disable" up

# Start the application
/app/main