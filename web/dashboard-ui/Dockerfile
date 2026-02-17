# Gunakan image slim (Riset 3)
FROM node:20-slim

WORKDIR /app

# Copy package files
COPY package*.json ./

# KONFIGURASI JARINGAN "ANTI-BADAI":
# 1. Gunakan Mirror NPM China/Asia (Lebih dekat & stabil dari Indonesia)
# 2. Matikan Strict SSL (Bypass transparent proxy ISP)
# 3. Set Timeout & Retries tinggi via CONFIG (bukan flag)
RUN npm config set registry https://registry.npmmirror.com/ && \
    npm config set strict-ssl false && \
    npm config set fetch-retries 10 && \
    npm config set fetch-retry-mintimeout 60000 && \
    npm config set fetch-retry-maxtimeout 600000 && \
    npm config set fetch-timeout 900000 && \
    npm install --no-audit --no-fund --prefer-offline

COPY . .

EXPOSE 5173

CMD ["npm", "run", "dev", "--", "--host"]