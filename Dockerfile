# Dockerfile untuk Go

# --- Tahap Build ---
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Salin file dependensi dan unduh
COPY go.mod go.sum ./
RUN go mod download

# Salin sisa kode sumber
COPY . .

# Build aplikasi Go
RUN CGO_ENABLED=0 GOOS=linux go build -o /main .

# --- Tahap Produksi ---
FROM alpine:latest

WORKDIR /app

# Salin hasil build dari tahap sebelumnya
COPY --from=builder /app/main .

# Port yang akan diekspos oleh aplikasi Go Anda
EXPOSE 8003

# Perintah untuk menjalankan aplikasi saat container dimulai
# Variabel lingkungan akan disuntikkan oleh Docker Compose
CMD ["/main"]