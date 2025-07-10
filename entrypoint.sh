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

# Start the Go gRPC application in the background
/app/main &

# Wait for the Go gRPC application to be ready
echo "Waiting for Go gRPC application..."
while ! nc -z localhost 50051; do
  sleep 1
done
echo "Go gRPC application is ready."

# Start grpcwebproxy in the foreground
# It will listen on 8080 and forward to the Go app on 50051
exec /app/grpcwebproxy \
  --backend_addr=localhost:50051 \
  --backend_tls=false \
  --run_tls_server=false \
  --server_http_debug_port=8080