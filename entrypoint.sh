#!/bin/sh

set -e

# Run database migrations
/app/migrate -path /app/pkg/database/migrations -database "postgres://${DB_USER}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=disable" up

# Start the application
/app/main
