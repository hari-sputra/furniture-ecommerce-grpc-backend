# Go gRPC E-Commerce Backend

This project is the gRPC backend for a furniture e-commerce platform. It is written in Go and provides services for user authentication and other e-commerce functionalities.

## Technologies Used

- **Go**: The programming language used for the backend.
- **gRPC**: The framework used for communication between the backend and client applications.
- **Protocol Buffers**: The language-neutral, platform-neutral, extensible mechanism for serializing structured data.
- **PostgreSQL**: The relational database used for data storage.
- **Docker**: The containerization platform used for building and running the application.
- **gRPC Web Proxy**: To make the gRPC backend accessible from web browsers.

## Getting Started

### Prerequisites

- Docker
- Docker Compose
- Go
- `migrate` tool
- `protoc` compiler

### Installation

1. **Clone the repository:**

   ```bash
   git clone <repository-url>
   ```

2. **Create a `.env` file** in the root directory with the following content:

   ```bash
   DB_USER=your_db_user
   DB_PASSWORD=your_db_password
   DB_HOST=your_db_host
   DB_PORT=your_db_port
   DB_NAME=furniture_db
   ```

3. **Run the application using Docker Compose:**

   ```bash
   docker-compose up -d
   ```

## Usage

### Compiling Protocol Buffers

To generate the Go and TypeScript code from the `.proto` files, use the following commands:

**Go:**

```bash
protoc --go_out=./pb \
 --go-grpc_out=./pb \
 --proto_path=./proto --go_opt=paths=source_relative \
 --go-grpc_opt=paths=source_relative \
 service/service.proto
```

**TypeScript:**

```bash
npx protoc --ts_out=./pb \
 --proto_path=./proto auth/auth.proto
```

### Running the gRPC Web Proxy

To make the backend accessible from a web browser, run the gRPC web proxy with the following command:

```bash
grpcwebproxy --backend_addr=localhost:50051 \
  --server_bind_address=0.0.0.0 \
  --server_http_debug_port=8080 \
  --run_tls_server=false \
  --backend_max_call_recv_msg_size=577659248 \
  --allow_all_origins
```

### Running Database Migrations

To apply the database schema changes, use the `migrate` tool. Make sure your database is running and accessible.

```bash
migrate \
  -path pkg/database/migrations \
  -database "postgres://postgres:password@localhost:5432/furniture_db?sslmode=disable" \
  up
```
