# Dockerfile untuk Go

# --- Tahap Build ---
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Install migrate tool
RUN apk --no-cache add curl && \
    curl -L https://github.com/golang-migrate/migrate/releases/download/v4.18.3/migrate.linux-amd64.tar.gz | tar xvz && \
    mv migrate /usr/local/bin/migrate

# Install grpcwebproxy tool
RUN curl -L -o /usr/local/bin/grpcwebproxy https://github.com/grpc/grpcweb/releases/download/1.5.0/grpcwebproxy-v1.5.0-linux-x86_64 && \
    chmod +x /usr/local/bin/grpcwebproxy

# Salin file dependensi dan unduh
COPY go.mod go.sum ./
RUN go mod download

# Salin sisa kode sumber
COPY . .

# Build aplikasi Go
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./main.go

# --- Tahap Produksi ---
FROM alpine:latest

WORKDIR /app

RUN apk --no-cache add ca-certificates netcat-openbsd

# Salin hasil build dari tahap sebelumnya
COPY --from=builder /app/main .
COPY --from=builder /usr/local/bin/migrate .
COPY --from=builder /usr/local/bin/grpcwebproxy .
COPY --from=builder /app/pkg ./pkg

COPY entrypoint.sh .
RUN chmod +x /app/entrypoint.sh

# Copy the production environment file and rename it to .env
COPY prod.env .env

# Port yang akan diekspos oleh aplikasi Go Anda
EXPOSE 8080

# Perintah untuk menjalankan aplikasi saat container dimulai
# Variabel lingkungan akan disuntikkan oleh Docker Compose
CMD ["./entrypoint.sh"]
