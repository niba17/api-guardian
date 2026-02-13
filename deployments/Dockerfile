# Stage 1: Kompilasi (Koki Memasak)
FROM golang:alpine AS builder

WORKDIR /app

# Install git & sertifikat keamanan
RUN apk add --no-cache git ca-certificates

# Copy daftar belanjaan dulu (Layer Caching)
COPY go.mod go.sum ./
RUN go mod download

# Copy seluruh kode sumber dari Root Project
COPY . .

# Build aplikasi
RUN CGO_ENABLED=0 GOOS=linux go build -o guardian-app ./cmd/server

# Stage 2: Image Berjalan (Penyajian)
FROM alpine:latest

# 👇 UPDATE: Kita pakai /app biar konsisten dengan docker-compose
WORKDIR /app

# Install CA Certs (Penting buat HTTPS ke Google)
RUN apk --no-cache add ca-certificates

# Copy binary aplikasi dari Stage 1
COPY --from=builder /app/guardian-app .

# 👇 PERBAIKAN UTAMA DI SINI 👇
# Kita copy satu folder 'configs' (berisi geoip/GeoLite2-City.mmdb)
# ke dalam container di folder /app/configs
COPY configs/ ./configs/

# Buat folder logs kosong
RUN mkdir logs

# Buka Port
EXPOSE 8080

# Jalankan
CMD ["./guardian-app"]