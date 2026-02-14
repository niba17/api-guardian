# Stage 1: Kompilasi (Gunakan versi yang lebih baru dari 1.25.7)
FROM golang:1.26-alpine AS builder

WORKDIR /app
RUN apk add --no-cache git ca-certificates

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Build aplikasi
RUN CGO_ENABLED=0 GOOS=linux go build -o guardian-app ./cmd/server

# Stage 2: Image Berjalan
FROM alpine:latest
WORKDIR /app
RUN apk --no-cache add ca-certificates tzdata

COPY --from=builder /app/guardian-app .
COPY --from=builder /app/configs ./configs
RUN mkdir logs

EXPOSE 8080
CMD ["./guardian-app"]