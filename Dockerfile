# Dockerfile untuk Go

# --- Tahap Build ---
FROM golang:1.24-alpine AS builder

WORKDIR /app

RUN go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

# Salin file dependensi dan unduh
COPY go.mod go.sum ./
RUN go mod download

# Salin sisa kode sumber
COPY . .

# Build aplikasi Go
RUN CGO_ENABLED=0 GOOS=linux go build -o /main .

# --- Tahap Produksi ---
FROM alpine:latest

WORKDIR /

# Salin hasil build dari tahap sebelumnya
COPY --from=builder /main /main

# Port yang akan diekspos oleh aplikasi Go Anda
EXPOSE 8003

# Perintah untuk menjalankan aplikasi saat container dimulai
# Variabel lingkungan akan disuntikkan oleh Docker Compose
CMD ["/main"]