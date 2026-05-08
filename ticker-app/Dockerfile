# ─── Stage 1: Build ───────────────────────────────────
FROM node:22-alpine AS builder

WORKDIR /app

# Erst nur package.json kopieren → besseres Layer-Caching
COPY package*.json ./
RUN npm ci

# Dann Source Code kopieren und bauen
COPY . .
RUN npm run build

# ─── Stage 2: Production ──────────────────────────────
FROM nginx:alpine AS production

# Gebautes Frontend in Nginx kopieren
COPY --from=builder /app/dist /usr/share/nginx/html

# Nginx Konfiguration für Vue Router (SPA)
COPY nginx.conf /etc/nginx/conf.d/default.conf

EXPOSE 80

CMD ["nginx", "-g", "daemon off;"]