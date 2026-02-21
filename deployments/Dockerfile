# Stage 1: Kompilasi
FROM golang:1.26-alpine AS builder

# Tidak perlu git atau proxy karena kita pakai vendor
RUN apk add --no-cache tzdata ca-certificates git

WORKDIR /app

COPY go.mod go.sum ./

ENV GOPROXY=https://proxy.golang.org,direct

RUN go mod download

# Copy seluruh source code sisanya
COPY . .

# Build dengan flag -mod=vendor
RUN CGO_ENABLED=0 GOOS=linux go build -mod=vendor -o guardian-app ./cmd/server

# Stage 2: Image Berjalan
FROM alpine:latest
WORKDIR /app
RUN apk --no-cache add ca-certificates tzdata

COPY --from=builder /app/guardian-app .
COPY --from=builder /app/configs ./configs
RUN mkdir logs

EXPOSE 8080
CMD ["./guardian-app"]