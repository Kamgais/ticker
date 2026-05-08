# ─── Stage 1: Build ───────────────────────────────────
FROM golang:1.23-alpine AS builder

WORKDIR /app

# Erst nur go.mod und go.sum kopieren → besseres Layer-Caching
COPY go.mod go.sum ./
RUN go mod download

# Dann Source Code kopieren und kompilieren
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o ticker-backend ./cmd/main.go

# ─── Stage 2: Production ──────────────────────────────
# Distroless = minimales Image, nur das Binary
FROM gcr.io/distroless/static-debian12 AS production

WORKDIR /app

# Nur das kompilierte Binary kopieren
COPY --from=builder /app/ticker-backend .

EXPOSE 8080

CMD ["/app/ticker-backend"]