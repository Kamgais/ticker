# 📰 Ticker App — Full-Stack Webanwendung

> Arbeitsprobe für die Position **Full-Stack Developer (m/w/d)**  
> German Tote Service- und Beteiligungs GmbH / Wettstar  
> Entwickelt von **Cyril Kamgais Totso**  
> GitHub: [Kamgais](https://github.com/Kamgais) · Portfolio: [personal-portfolio-v2-delta.vercel.app](https://personal-portfolio-v2-delta.vercel.app)

---

## 📋 Inhaltsverzeichnis

- [Über das Projekt](#-über-das-projekt)
- [Demo & Login](#-demo--login)
- [Setup & Start](#-setup--start)
- [Architektur & Systemdesign](#-architektur--systemdesign)
- [Technische Entscheidungen](#-technische-entscheidungen)
- [DevOps & CI/CD](#-devops--cicd)
- [Tests](#-tests)
- [Implementierte Features](#-implementierte-features)
- [API-Dokumentation](#-api-dokumentation)
- [KI-Unterstützung](#-ki-unterstützung)
- [Verbesserungsmöglichkeiten](#-verbesserungsmöglichkeiten)

---

## 🎯 Über das Projekt

Diese Anwendung wurde als Arbeitsprobe entwickelt und demonstriert meine Fähigkeit,
eine vollständige Full-Stack-Webanwendung von der Architektur bis zum Deployment
eigenständig zu konzipieren, zu entwickeln und zu containerisieren.

### Was wurde gebaut?

- **Vue 3 Frontend** mit TypeScript, Pinia State Management und Tailwind CSS
- **Go Backend** als intelligenter Proxy mit In-Memory Caching und Rate Limiting
- **Docker Setup** mit Multi-Stage Builds, Nginx als Reverse Proxy und Docker Compose
- **GitHub Actions CI/CD Pipeline** mit automatischen Tests und DockerHub Push
- **16 Unit Tests** für den Pinia Store mit Vitest
- **Alle 8 Zusatzfeatures** aus der Aufgabenstellung vollständig implementiert
- **Optionaler Backend-Teil** in Go vollständig umgesetzt

### Mein Denkprozess — bevor ich Code geschrieben habe

Bevor ich mit der Implementierung begann, habe ich die Aufgabe systematisch
analysiert und folgende Designentscheidungen bewusst getroffen:

**1. API zuerst verstehen:**
Ich habe die externe API direkt getestet und die echte JSON-Antwort analysiert.
Erst danach habe ich die TypeScript-Interfaces definiert — nicht geraten, sondern
aus echten Daten abgeleitet. Das verhinderte spätere Typfehler und Refactoring-Aufwand.

**2. Schichtenarchitektur von Anfang an:**
Ich habe bewusst eine klare Trennung eingeführt:
`Types → API-Service → Pinia Store → Komponenten`
Jede Schicht hat genau eine Verantwortung.

**3. Warum ein Go-Backend statt direkter API-Anbindung:**
Eine direkte Anbindung im Frontend wäre schneller gewesen — aber ich habe
das Go-Backend bewusst hinzugefügt um Caching, Kapselung und Rate Limiting
zu realisieren.

**4. DevOps von Anfang an mitgedacht:**
Dockerfiles wurden mit Multi-Stage Builds designed, eine CI/CD Pipeline mit
GitHub Actions aufgebaut und Images automatisch zu DockerHub gepusht.

---

## 🚀 Demo & Login

```bash
git clone https://github.com/Kamgais/ticker.git
cd ticker
docker compose up --build
```

App öffnen: **http://localhost**

| Benutzer | Passwort  |
|----------|-----------|
| admin    | admin123  |
| cyril    | cyril123  |

Oder direkt von DockerHub:

```bash
docker pull kamgais/ticker-frontend:latest
docker pull kamgais/ticker-backend:latest
```

---

## ⚙️ Setup & Start

### Voraussetzungen

| Tool | Version | Zweck |
|------|---------|-------|
| Docker Desktop | Latest | Empfohlene Ausführung |
| Node.js | 22+ | Lokale Frontend-Entwicklung |
| Go | 1.23+ | Lokale Backend-Entwicklung |

### Option 1: Docker (empfohlen)

```bash
docker compose up --build
```

| Service | URL |
|---------|-----|
| Frontend (Nginx) | http://localhost |
| Backend API | http://localhost:8080/api/ticker |
| Backend Health | http://localhost:8080/health |

Stoppen:
```bash
docker compose down
```

### Option 2: Production — direkt von DockerHub

Kein lokales Bauen nötig — Images werden direkt von DockerHub geholt:

```bash
docker compose -f docker-compose.prod.yml up
```

App öffnen: **http://localhost**

| | docker-compose.yml | docker-compose.prod.yml |
|--|-------------------|------------------------|
| Images | Lokal gebaut | Von DockerHub geholt |
| Geschwindigkeit | ~2-3 Minuten | ~30 Sekunden |
| Zweck | Entwicklung | Production / Demo |

### Option 3: Lokal ohne Docker

**Terminal 1 — Go Backend:**
```bash
cd backend
go run cmd/main.go
```

**Terminal 2 — Vue Frontend:**
```bash
cd ticker-app
npm install
npm run dev
```

### Option 4: Tests ausführen

```bash
cd ticker-app
npm run test
```

---

## 🏗️ Architektur & Systemdesign

### Systemübersicht

```
+------------------------------------------------------------------+
|                         Docker Network                            |
|                                                                  |
|  +-------------------------+      +---------------------------+  |
|  |   Vue 3 Frontend        |      |      Go Backend           |  |
|  |   Nginx :80             +----->+      :8080                |  |
|  |                         |      |                           |  |
|  |  [ Pinia Store        ] |      |  [ In-Memory Cache      ] |  |
|  |  [ filteredEntries    ] |      |  [ TTL: 30 Sekunden      ] |  |
|  |  [ sortedEntries      ] |      |  [ Rate Limiter          ] |  |
|  |  [ paginatedEntries   ] |      |  [ 60 req/min/IP         ] |  |
|  +-------------------------+      +-------------+-------------+  |
+--------------------------------------------------+---------------+
                                                   |
                                                   v
                                  +--------------------------------+
                                  |   Wettstar AWS API             |
                                  |   eu-central-1.amazonaws.com   |
                                  +--------------------------------+
```

### Projektstruktur

```
ticker/
├── .github/
│   └── workflows/
│       └── ci.yml                 # GitHub Actions CI/CD Pipeline
├── ticker-app/                    # Vue 3 Frontend
│   ├── src/
│   │   ├── __tests__/
│   │   │   └── ticker.store.test.ts  # 16 Unit Tests (Vitest)
│   │   ├── types/
│   │   │   └── ticker.ts          # TypeScript Interfaces
│   │   ├── services/
│   │   │   └── tickerApi.ts       # Axios HTTP-Client mit Retry-Logik
│   │   ├── stores/
│   │   │   └── ticker.ts          # Pinia Store (State, Getters, Actions)
│   │   ├── components/
│   │   │   ├── TickerCard.vue     # Karte mit Edit/Delete
│   │   │   ├── TickerList.vue     # Liste aller Karten
│   │   │   ├── TickerModal.vue    # Erstellen/Bearbeiten Modal
│   │   │   ├── SearchBar.vue      # Suche + Sortierung
│   │   │   ├── Pagination.vue     # Seitennavigation
│   │   │   ├── DarkModeToggle.vue # Dark/Light Mode
│   │   │   ├── LoginScreen.vue    # Auth-Simulation
│   │   │   ├── ErrorBanner.vue    # Fehleranzeige
│   │   │   └── LoadingSpinner.vue # Ladezustand
│   │   ├── views/
│   │   │   └── HomeView.vue       # Hauptseite mit Auto-Refresh
│   │   ├── App.vue                # Root mit Auth-Guard
│   │   └── main.ts                # Einstiegspunkt
│   ├── Dockerfile                 # Multi-Stage: Node -> Nginx
│   └── nginx.conf                 # SPA-Routing + API-Proxy
├── backend/                       # Go Backend
│   ├── cmd/
│   │   └── main.go                # HTTP-Server, Router, Middleware
│   ├── internal/
│   │   ├── cache/
│   │   │   └── cache.go           # Thread-sicherer In-Memory Cache
│   │   ├── handler/
│   │   │   └── ticker.go          # HTTP-Handler + Validierung + Proxy
│   │   └── middleware/
│   │       └── ratelimit.go       # IP-basierter Rate Limiter
│   └── Dockerfile                 # Multi-Stage: Go -> Distroless
└── docker-compose.yml             # Lokale Orchestrierung
```

---

## 💡 Technische Entscheidungen

### 1. Warum Go-Backend statt direkter API-Anbindung im Frontend

Die einfachste Lösung wäre gewesen, die externe Wettstar API direkt
im Vue Frontend aufzurufen. Ich habe mich bewusst dagegen entschieden
und ein Go-Backend dazwischengeschaltet — aus drei konkreten Gründen:

**Caching:** Ohne Backend macht jeder Auto-Refresh alle 10 Sekunden
einen neuen Request zur externen API. Mit dem Go-Backend cache ich
die Antwort für 30 Sekunden. Das bedeutet: bei 10 gleichzeitigen
Nutzern macht nicht jeder seinen eigenen API-Call — alle bekommen
die gecachte Antwort in unter 1ms. Der `X-Cache: HIT` Header im
Response zeigt jederzeit ob der Cache genutzt wurde.

**Kapselung:** Die externe API-URL ist nur im Go-Backend bekannt.
Das Frontend kennt nur `/api/ticker`. Wenn sich die externe
API-URL ändert, muss ich nur das Backend anpassen — kein
Frontend-Deployment nötig.

**Rate Limiting:** Mein Rate Limiter erlaubt maximal 60 Requests
pro Minute pro IP — darüber hinaus gibt es einen 429-Fehler,
bevor der Request die externe API überhaupt erreicht.

```
Ohne Backend:    Browser --> AWS API (jeder Request!)
Mit Go-Backend:  Browser --> Nginx --> Go --> Cache --> AWS API (selten)
```

---

### 2. Pinia Pipeline-Architektur

Ich habe eine Pipeline aus `computed()`-Werten gebaut:

```
entries           (Rohdaten aus API)
    |
    v
filteredEntries   (Suchfilter angewendet)
    |
    v
sortedEntries     (Sortierung angewendet)
    |
    v
paginatedEntries  (Aktuelle Seite — das zeigt die UI)
```

Jede Schicht ist ein `computed()` — ändert sich `searchQuery`,
aktualisieren sich alle drei Schichten automatisch und reaktiv.
Kein manuelles Triggern, keine Race Conditions. Außerdem setze
ich `currentPage` bei jeder Suche und Sortierung auf 1 zurück.

---

### 3. Multi-Stage Dockerfile — ~10MB statt ~800MB

```dockerfile
# Stage 1: Bauen (800MB)
FROM golang:1.23-alpine AS builder
RUN go build -o ticker-backend ./cmd/main.go

# Stage 2: Production (10MB)
FROM gcr.io/distroless/static-debian12 AS production
COPY --from=builder /app/ticker-backend .
```

Das finale Image enthält nur das kompilierte Binary — kein
Go-Compiler, keine Shell, kein Package Manager. Schnellere
Deployments und minimale Angriffsfläche.

---

### 4. Nginx als Reverse Proxy

```nginx
location /api {
    proxy_pass http://backend:8080;
}
```

Das Frontend macht relative Requests (`/api/ticker`) —
keine hardcodierte URL, kein CORS-Problem, funktioniert
in jeder Umgebung. Nginx cached statische Assets ein Jahr.

---

### 5. TypeScript-Types aus echter API abgeleitet

Ich habe die API direkt getestet bevor ich Types definiert habe.
Die Felder hießen anders als erwartet:

```typescript
ticker_id   // nicht "id"
message     // nicht "text"
creator     // nicht "author"
created     // nicht "date"
```

Durch Analyse zuerst war mein Code von Anfang an korrekt —
kein nachträgliches Refactoring nötig.

---

### 6. Delete — Soft Delete via PATCH active: false

Die externe Wettstar API stellt keinen DELETE-Endpunkt zur
Verfügung (`Access-Control-Allow-Methods: OPTIONS, POST`).

Ich habe eine pragmatische Lösung entwickelt: Ich setze
`active: false` via PATCH um einen Eintrag zu deaktivieren.
Da PATCH alle Pflichtfelder benötigt, lade ich den Eintrag
zuerst aus dem Cache, setze `active: false` und schicke
den vollständigen Payload:

```go
target["active"] = false
req, _ := http.NewRequest(http.MethodPatch, upstreamURL+"/ticker/", body)
```

Im Backend filtere ich beim GET alle inaktiven Einträge
heraus — die UI zeigt ausschließlich `active: true` Einträge.

```
DELETE Button --> PATCH active: false --> Cache invalidiert
             --> GET filtert active: false raus
             --> Eintrag dauerhaft verschwunden ✅
```

---

## 🐳 DevOps & CI/CD

### GitHub Actions Pipeline

Bei jedem Push auf `main` läuft automatisch:

```
Push auf main
    |
    +-- Job 1: Frontend
    |     ├── TypeScript prüfen (tsc --noEmit)
    |     ├── 16 Unit Tests (vitest run)
    |     └── Production Build (npm run build)
    |
    +-- Job 2: Backend
    |     ├── Go Build (go build ./...)
    |     └── Go Vet (go vet ./...)
    |
    +-- Job 3: Docker (nur wenn Job 1+2 grün)
          ├── Login zu DockerHub
          ├── Frontend Image bauen + pushen
          └── Backend Image bauen + pushen
```

**Zwei Image-Tags pro Push:**
- `latest` — immer der neueste Stand
- `git-sha` (z.B. `a3f2c1d`) — exakte Version für Rollbacks

```bash
# Rollback auf eine bestimmte Version
docker pull kamgais/ticker-frontend:a3f2c1d
```

### Docker Images auf DockerHub

```bash
docker pull kamgais/ticker-frontend:latest
docker pull kamgais/ticker-backend:latest
```

### Multi-Stage Builds

| Image | Build Stage | Production Stage | Größe |
|-------|-------------|-----------------|-------|
| Frontend | node:22-alpine | nginx:alpine | ~25MB |
| Backend | golang:1.23-alpine | distroless | ~10MB |

### Kubernetes (mögliche Erweiterung)

Als nächster Schritt könnte die App auf einem Kubernetes Cluster
deployt werden. Dank der bestehenden Docker Images auf DockerHub
wären folgende Manifeste schnell umsetzbar:

- `Deployment` — 2 Replicas mit Liveness/Readiness Probes
- `HPA` — Auto-Scaling bei CPU > 70% (bis 10 Pods)
- `Ingress` — Traffic-Routing über Nginx Ingress Controller

Ich habe praktische Erfahrung mit Kubernetes (Helm, HPA, CronJobs,
StatefulSets, Ingress) aus akademischen Projekten — die Umsetzung
wäre ein logischer nächster Schritt.

---

## 🧪 Tests

### 16 Unit Tests — Pinia Store (Vitest)

```bash
cd ticker-app
npm run test
```

```
✓ hat korrekten initialen State
✓ filtert Eintraege nach Titel
✓ filtert Eintraege nach Ersteller
✓ filtert Eintraege nach Nachricht
✓ gibt alle Eintraege zurueck wenn Suche leer ist
✓ setzt currentPage auf 1 bei neuer Suche
✓ sortiert Eintraege nach Titel aufsteigend
✓ sortiert Eintraege nach Titel absteigend
✓ sortiert Eintraege nach ID aufsteigend
✓ berechnet totalPages korrekt
✓ gibt korrekte Eintraege fuer Seite 1 zurueck
✓ gibt korrekte Eintraege fuer letzte Seite zurueck
✓ laedt Eintraege erfolgreich
✓ setzt status auf error bei API-Fehler
✓ fuegt neuen Eintrag am Anfang der Liste hinzu
✓ entfernt Eintrag nach deleteEntry

Tests: 16 passed (16)
```

### Test-Kategorien

| Kategorie | Tests | Was wird getestet |
|-----------|-------|-------------------|
| State | 1 | Initialer Store-Zustand |
| Suche | 5 | Filter nach Titel, Nachricht, Ersteller |
| Sortierung | 3 | Auf-/absteigende Sortierung |
| Pagination | 3 | Seitenberechnung und Navigation |
| API Actions | 4 | fetchEntries, createEntry, deleteEntry, Fehlerfall |

### Test-Strategie

Ich habe den API-Service komplett gemockt (`vi.mock`) — so teste
ich den Store isoliert ohne echte HTTP-Requests. Das macht Tests
schnell, deterministisch und unabhängig von der externen API.

---

## ✅ Implementierte Features

### Pflichtfunktionen

| Feature | Status | Details |
|---------|--------|---------|
| Einträge laden & anzeigen | ✅ | GET /api/ticker mit Caching |
| Neuen Eintrag erstellen | ✅ | POST mit Formularvalidierung |
| Eintrag bearbeiten | ✅ | PATCH mit vorausgefülltem Modal |
| Eintrag löschen | ✅ | Soft Delete via PATCH active: false |
| Responsive UI | ✅ | Tailwind CSS Breakpoints |
| Ladezustände | ✅ | Spinner + Refresh-Indikator mit Countdown |
| API-Fehlerbehandlung | ✅ | Deutsche Fehlermeldungen pro HTTP-Status |
| Benutzerführung | ✅ | Lösch-Bestätigung, Validierung, Highlight-Badge |

### Zusatzfeatures (alle 8 implementiert)

| Feature | Status | Details |
|---------|--------|---------|
| Auto-Refresh | ✅ | Alle 10s mit sichtbarem Countdown im Header |
| Suche / Filter | ✅ | Echtzeit-Suche in Titel, Nachricht, Ersteller |
| Sortierung | ✅ | Nach Datum, Titel, ID — auf-/absteigend |
| Pagination | ✅ | 3 Einträge pro Seite mit vollständiger Navigation |
| Offline/Retry-Logik | ✅ | 3 Versuche mit exponentiellem Backoff |
| Dark Mode | ✅ | Persistiert in localStorage |
| Auth-Simulation | ✅ | Login-Screen mit Session-Persistenz |
| Docker-Setup | ✅ | Multi-Stage Builds + Docker Compose |

### Optionaler Backend-Teil (Go)

| Feature | Status | Details |
|---------|--------|---------|
| API-Proxy | ✅ | Kapselt externe Wettstar API |
| In-Memory Caching | ✅ | TTL: 30s, Thread-sicher mit RWMutex |
| Request-Validierung | ✅ | Pflichtfelder geprüft vor upstream-Call |
| Rate Limiting | ✅ | 60 req/min/IP |
| Health Endpoint | ✅ | GET /health für Liveness Probes |

### Extras (über Aufgabenstellung hinaus)

| Feature | Details |
|---------|---------|
| GitHub Actions CI/CD | Tests + Build + DockerHub Push bei jedem Commit |
| DockerHub Images | Automatisch gepusht mit latest + git-sha Tags |
| 16 Unit Tests | Vitest — State, Suche, Sortierung, Pagination, Actions |
| /health Endpoint | Für Kubernetes Liveness/Readiness Probes |

---

## 📡 API-Dokumentation

### Externe Wettstar API (upstream)

| Method | Endpoint | Beschreibung |
|--------|----------|--------------|
| GET | `/ticker/?tickernames=wpftest` | Alle Einträge laden |
| POST | `/ticker/` | Neuen Eintrag erstellen |
| PATCH | `/ticker/` | Eintrag bearbeiten |
| DELETE | — | Nicht unterstützt |

### Go Backend API (intern)

| Method | Endpoint | Beschreibung |
|--------|----------|--------------|
| GET | `/api/ticker` | Einträge laden (mit Cache, nur active: true) |
| POST | `/api/ticker` | Eintrag erstellen (mit Validierung) |
| PATCH | `/api/ticker` | Eintrag bearbeiten (Cache invalidieren) |
| DELETE | `/api/ticker/{id}` | Soft Delete via PATCH active: false |
| GET | `/health` | Health Check für Kubernetes Probes |

### Response Headers

| Header | Wert | Bedeutung |
|--------|------|-----------|
| X-Cache | HIT | Antwort aus Cache (< 1ms) |
| X-Cache | MISS | Antwort von upstream API |

---

## 🤖 KI-Unterstützung

### Genutzte Tools

- **Claude (Anthropic)** — als Mentor, Lernbegleiter und Code-Assistent

### Wofür KI genutzt wurde

**Architektur & Konzepte:**
- Erklärung von Vue 3 Composition API (ref, reactive, computed, watch)
- Erklärung von Pinia Store-Architektur
- Erklärung von Go-Konzepten (Goroutines, sync.RWMutex, HTTP-Handler)
- Diskussion der Architekturentscheidungen

**UI & Design:**
- Tailwind CSS Klassen für responsives Layout und Dark Mode
- Komponenten-Struktur (Modal mit Teleport, Pagination)
- Farbschema und visuelle Hierarchie der Ticker-Karten

**Code-Generierung:**
- Grundstruktur der Go-Dateien
- Grundstruktur der Vue-Komponenten
- Docker-, Nginx- und GitHub Actions Konfigurationen

### Was ich selbst erarbeitet und angepasst habe

- Die externe API selbst getestet und JSON-Struktur analysiert
- TypeScript-Interfaces aus echter API-Antwort abgeleitet
- Delete-Problem selbst identifiziert und Lösung via PATCH entwickelt
- Docker Build-Fehler (Go-Version) debuggt und behoben
- Nginx-Proxy-Konfiguration debuggt (BASE_URL angepasst)
- GitHub Secrets konfiguriert und DockerHub Integration eingerichtet
- Alle Konzepte verstanden bevor weitergemacht wurde

### Fazit

KI wurde als intelligentes Werkzeug eingesetzt — nicht als Ersatz
für eigenes Denken. Architekturentscheidungen, Debugging und
Problemlösung lagen bei mir. KI hat die Geschwindigkeit erhöht,
nicht die Denkarbeit ersetzt.

---

## 🔮 Verbesserungsmöglichkeiten

### Kurzfristig

1. **E2E Tests** — Cypress für vollständige Browser-Tests
2. **Optimistic Updates** — UI sofort aktualisieren, bei Fehler zurückrollen
3. **Umgebungsvariablen** — API-URL über `.env` konfigurierbar

### Mittelfristig

4. **WebSocket** — Echtzeit-Updates statt Polling alle 10 Sekunden
5. **HTTPS** — SSL-Zertifikat mit Let's Encrypt
6. **Staging Environment** — Separate Pipeline für Staging vs. Production

### Langfristig

7. **Backend Persistenz** — Eigene Datenbank für Offline-Modus
8. **Kubernetes Deployment** — Dank der bestehenden DockerHub Images
   wäre ein Kubernetes Deployment mit HPA und Ingress der nächste
   logische Schritt für produktionsreife Skalierung.

---

## 👨‍💻 Über den Entwickler

**Cyril Kamgais Totso**  
M.Sc. Informatik (Cloud Computing & AWS) — TH Brandenburg  
3+ Jahre Erfahrung in Full-Stack-Entwicklung und Cloud-Architektur

| Kontakt | |
|---------|---|
| Email | cyrilkamgais1203@gmail.com |
| Telefon | +49 176 68484380 |
| LinkedIn | [cyril-kamgais-totso-a86491207](https://linkedin.com/in/cyril-kamgais-totso-a86491207) |
| GitHub | [Kamgais](https://github.com/Kamgais) |
| Portfolio | [personal-portfolio-v2-delta.vercel.app](https://personal-portfolio-v2-delta.vercel.app) |