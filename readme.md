# Project Documentation

## Project Overview

### Build file proto

Example:

```bash

protoc --go_out=./pb \
 --go-grpc_out=./pb \
 --proto_path=./proto --go_opt=paths=source_relative \
 --go-grpc_opt=paths=source_relative \
 service/service.proto

```

### migrate SQL

Example:

```bash

migrate \
  -path pkg/database/migrations \
  -database "postgres://postgres:password@localhost:5432/furniture_db?sslmode=disable" \
  up


```
