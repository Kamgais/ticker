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
- [Frontend — Technische Entscheidungen](#-frontend--technische-entscheidungen)
- [Backend — Technische Entscheidungen](#-backend--technische-entscheidungen)
- [DevOps & Containerisierung](#-devops--containerisierung)
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
- **Alle 8 Zusatzfeatures** aus der Aufgabenstellung vollständig implementiert
- **Optionaler Backend-Teil** in Go vollständig umgesetzt

### Mein Denkprozess — bevor ich Code geschrieben habe

Bevor ich mit der Implementierung begann, habe ich die Aufgabe systematisch analysiert
und folgende Designentscheidungen bewusst getroffen:

**1. API zuerst verstehen:**
Ich habe die externe API direkt getestet und die echte JSON-Antwort analysiert.
Erst danach habe ich die TypeScript-Interfaces definiert — nicht geraten, sondern
aus echten Daten abgeleitet. Das verhinderte spätere Typfehler und Refactoring-Aufwand.

**2. Schichtenarchitektur von Anfang an:**
Ich habe bewusst eine klare Trennung eingeführt:
`Types → API-Service → Pinia Store → Komponenten`
Jede Schicht hat genau eine Verantwortung. Der Store weiß nichts von Axios,
die Komponenten wissen nichts von der API-URL.

**3. Warum ein Go-Backend statt direkter API-Anbindung:**
Eine direkte Anbindung im Frontend wäre schneller gewesen — aber ich habe
das Go-Backend bewusst hinzugefügt, um drei Probleme zu lösen:
Caching (weniger API-Last), Kapselung (externe URL nur im Backend bekannt)
und Schutz (Rate Limiting gegen Missbrauch).

**4. DevOps von Anfang an mitgedacht:**
Dockerfiles wurden mit Multi-Stage Builds designed — nicht nachträglich
hinzugefügt. Das zeigt sich in der Image-Größe: Backend ~10MB statt ~800MB.

---

## 🚀 Demo & Login

```bash
git clone https://github.com/Kamgais/ticker.git
cd ticker
docker compose up --build
```

App öffnen: **http://localhost**

| Benutzer | Passwort  | Rolle  |
|----------|-----------|--------|
| admin    | admin123  | Admin  |
| cyril    | cyril123  | User   |

---

## ⚙️ Setup & Start

### Voraussetzungen

| Tool | Version | Zweck |
|------|---------|-------|
| Docker Desktop | Latest | Empfohlene Ausführung (ein Befehl) |
| Node.js | 22+ | Lokale Frontend-Entwicklung |
| Go | 1.23+ | Lokale Backend-Entwicklung |

### Option 1: Docker (empfohlen)

Ein Befehl startet Frontend + Backend + Nginx:

```bash
docker compose up --build
```

| Service | URL |
|---------|-----|
| Frontend (Nginx) | http://localhost |
| Backend API | http://localhost:8080/api/ticker |

Stoppen:
```bash
docker compose down
```

### Option 2: Lokal ohne Docker

**Terminal 1 — Go Backend starten:**
```bash
cd backend
go run cmd/main.go
# Backend läuft auf http://localhost:8080
```

**Terminal 2 — Vue Frontend starten:**
```bash
cd ticker-app
npm install
npm run dev
# Frontend läuft auf http://localhost:5173
```

> **Hinweis:** Bei lokaler Entwicklung zeigt das Backend `[CACHE MISS]` und
> `[CACHE HIT]` Logs im Terminal — so ist der Cache-Status jederzeit sichtbar.

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
├── ticker-app/                    # Vue 3 Frontend
│   ├── src/
│   │   ├── types/
│   │   │   └── ticker.ts          # TypeScript Interfaces (aus echter API abgeleitet)
│   │   ├── services/
│   │   │   └── tickerApi.ts       # Axios HTTP-Client mit Retry-Logik
│   │   ├── stores/
│   │   │   └── ticker.ts          # Pinia Store (State, Getters, Actions)
│   │   ├── components/
│   │   │   ├── TickerCard.vue     # Einzelne Karte mit Edit/Delete
│   │   │   ├── TickerList.vue     # Liste aller Karten
│   │   │   ├── TickerModal.vue    # Erstellen/Bearbeiten Modal
│   │   │   ├── SearchBar.vue      # Suche + Sortierung
│   │   │   ├── Pagination.vue     # Seitennavigation
│   │   │   ├── DarkModeToggle.vue # Dark/Light Mode Toggle
│   │   │   ├── LoginScreen.vue    # Auth-Simulation
│   │   │   ├── ErrorBanner.vue    # Fehleranzeige
│   │   │   └── LoadingSpinner.vue # Ladezustand
│   │   ├── views/
│   │   │   └── HomeView.vue       # Hauptseite mit Auto-Refresh
│   │   ├── App.vue                # Root mit Auth-Guard
│   │   └── main.ts                # Einstiegspunkt
│   ├── Dockerfile                 # Multi-Stage: Node Builder -> Nginx Production
│   └── nginx.conf                 # SPA-Routing + API-Proxy
├── backend/                       # Go Backend
│   ├── cmd/
│   │   └── main.go                # HTTP-Server, Router, Middleware-Chain
│   ├── internal/
│   │   ├── cache/
│   │   │   └── cache.go           # Thread-sicherer In-Memory Cache (RWMutex)
│   │   ├── handler/
│   │   │   └── ticker.go          # HTTP-Handler + Validierung + Upstream-Proxy
│   │   └── middleware/
│   │       └── ratelimit.go       # IP-basierter Rate Limiter
│   ├── go.mod
│   └── Dockerfile                 # Multi-Stage: Go Builder -> Distroless Production
├── docker-compose.yml             # Orchestrierung beider Services
└── README.md
```

### Datenfluss — GET /api/ticker

```
Browser
  |
  |  GET http://localhost/api/ticker
  v
Nginx (Port 80)
  |
  |  proxy_pass http://backend:8080
  v
Go Backend (Port 8080)
  |
  +-- Rate Limiter pruefen --> 429 wenn ueberschritten
  |
  +-- Cache pruefen
  |     +-- HIT  --> sofort antworten (< 1ms) + Header: X-Cache: HIT
  |     +-- MISS --> weiter zur upstream API
  |
  |  GET https://65q02hjc6c.execute-api.eu-central-1.amazonaws.com/prod/ticker/?tickernames=wpftest
  v
Wettstar AWS API
  |
  |  JSON Response
  v
Go Backend
  |
  +-- Im Cache speichern (TTL: 30s)
  +-- Header setzen: X-Cache: MISS
  |
  v
Nginx --> Browser
```

---

## 🖥️ Frontend — Technische Entscheidungen

### TypeScript-First Ansatz

Bevor ich irgendeine Komponente gebaut habe, habe ich zuerst die echte API-Antwort
analysiert und daraus präzise TypeScript-Interfaces abgeleitet:

```typescript
// Abgeleitet aus der echten API-Antwort — nicht geraten
export interface TickerEntry {
  ticker_id: number      // nicht "id"
  ticker_name: string
  created: string        // nicht "date"
  modified: string
  title: string
  message: string        // nicht "text"
  type: number
  creator: string        // nicht "author"
  highlight: boolean
  automatic: boolean
  active: boolean
}
```

### Pinia Store — Pipeline-Architektur

Der Store implementiert eine bewusste Daten-Pipeline:

```
entries (Rohdaten aus API)
    |
    v
filteredEntries (Suchfilter angewendet)
    |
    v
sortedEntries (Sortierung angewendet)
    |
    v
paginatedEntries (Aktuelle Seite ausgeschnitten)
```

Jede Schicht ist als `computed()` definiert — ändert sich `searchQuery`,
aktualisieren sich alle drei Schichten automatisch ohne manuelle Neuberechnung.

### API-Service — Retry mit exponentiellem Backoff

Bei Netzwerkproblemen versucht der Service automatisch bis zu 3 Mal:

```
Versuch 1 --> Fehler --> warte 1000ms
Versuch 2 --> Fehler --> warte 2000ms
Versuch 3 --> Fehler --> Fehlermeldung an UI
```

### Auto-Refresh mit Countdown-Anzeige

Der Auto-Refresh ist bewusst transparent gestaltet:
- Grüner Punkt: Daten aktuell
- Rotierendes blaues Icon: Wird gerade aktualisiert
- Countdown zeigt, wann die nächste Aktualisierung kommt

### Dark Mode — Kein Flash beim Laden

Dark Mode wird in `localStorage` gespeichert und in `App.vue` sofort
beim Start geladen — bevor Vue rendert. Das verhindert das "Flash of
Light" beim Seitenaufruf im Dark Mode.

---

## ⚙️ Backend — Technische Entscheidungen

### Warum Go?

| Kriterium | Go | Node.js |
|-----------|----|---------| 
| Kompiliertes Binary | Ja, ~6MB | Nein, braucht Runtime |
| Docker Image Größe | ~10MB | ~150MB |
| Concurrency | Goroutines (leichtgewichtig) | Event Loop |
| Typsicherheit | Statisch typisiert | Optional |
| Startup-Zeit | < 10ms | ~200ms |

### Thread-sicherer Cache mit sync.RWMutex

Der Cache verwendet `sync.RWMutex` für sichere parallele Zugriffe:

```go
// Mehrere Goroutines koennen gleichzeitig lesen
c.mu.RLock()
defer c.mu.RUnlock()

// Aber nur eine kann schreiben
c.mu.Lock()
defer c.mu.Unlock()
```

Ohne Mutex koennte ein gleichzeitiger Schreib- und Lesevorgang aus
verschiedenen HTTP-Handler-Goroutines den Cache korrumpieren.

### Cache-Invalidierung bei Schreiboperationen

Bei POST, PATCH und DELETE wird der Cache sofort invalidiert:

```go
h.cache.Invalidate(fmt.Sprintf("ticker:%s", tickerName))
```

Der naechste GET-Request bekommt garantiert frische Daten.

### Validierung im Backend

```go
// Pflichtfelder prüfen bevor upstream-Call
func validateCreatePayload(payload map[string]interface{}) error {
    required := []string{"title", "message", "creator"}
    // ...
}
```

Unnoetige Requests zur externen API mit ungueltigen Daten werden verhindert.

---

## 🐳 DevOps & Containerisierung

### Multi-Stage Dockerfile — Frontend

```dockerfile
# Stage 1: Build (~800MB Node.js Image)
FROM node:22-alpine AS builder
RUN npm ci && npm run build

# Stage 2: Production (~25MB Nginx)
FROM nginx:alpine AS production
COPY --from=builder /app/dist /usr/share/nginx/html
COPY nginx.conf /etc/nginx/conf.d/default.conf
```

**Resultat:** Production Image ~25MB statt ~800MB

### Multi-Stage Dockerfile — Backend

```dockerfile
# Stage 1: Build (~800MB Go Image)
FROM golang:1.23-alpine AS builder
RUN CGO_ENABLED=0 GOOS=linux go build -o ticker-backend ./cmd/main.go

# Stage 2: Production (~10MB Distroless)
FROM gcr.io/distroless/static-debian12 AS production
COPY --from=builder /app/ticker-backend .
```

**Resultat:** Production Image ~10MB statt ~800MB  
**Sicherheit:** Distroless enthaelt keine Shell — minimale Angriffsfläche.

### Nginx als Reverse Proxy

```nginx
# SPA-Routing: Alle Routen auf index.html
location / {
    try_files $uri $uri/ /index.html;
}

# API-Proxy: /api an Go Backend weiterleiten
location /api {
    proxy_pass http://backend:8080;
}

# Asset-Caching: 1 Jahr fuer statische Dateien
location ~* \.(js|css|png|svg|ico)$ {
    expires 1y;
    add_header Cache-Control "public, immutable";
}
```

---

## ✅ Implementierte Features

### Pflichtfunktionen

| Feature | Status | Details |
|---------|--------|---------|
| Einträge laden & anzeigen | ✅ | GET /api/ticker mit Caching |
| Neuen Eintrag erstellen | ✅ | POST mit Formularvalidierung |
| Eintrag bearbeiten | ✅ | PATCH mit vorausgefülltem Modal |
| Eintrag löschen | ✅ | DELETE mit Bestätigungsdialog |
| Responsive UI | ✅ | Tailwind CSS Breakpoints |
| Ladezustände | ✅ | Spinner + Refresh-Indikator mit Countdown |
| API-Fehlerbehandlung | ✅ | Deutsche Fehlermeldungen pro HTTP-Status-Code |
| Benutzerführung | ✅ | Lösch-Bestätigung, Validierung, Highlight-Badge |

### Zusatzfeatures (alle 8 implementiert)

| Feature | Status | Details |
|---------|--------|---------|
| Auto-Refresh | ✅ | Alle 10s mit sichtbarem Countdown im Header |
| Suche / Filter | ✅ | Echtzeit-Suche in Titel, Nachricht und Ersteller |
| Sortierung | ✅ | Nach Datum, Titel oder ID — auf-/absteigend |
| Pagination | ✅ | 5 Einträge pro Seite mit vollständiger Navigation |
| Offline/Retry-Logik | ✅ | 3 Versuche mit exponentiellem Backoff |
| Dark Mode | ✅ | Persistiert in localStorage, kein Flash beim Laden |
| Auth-Simulation | ✅ | Login-Screen mit Session-Persistenz |
| Docker-Setup | ✅ | Multi-Stage Builds + Docker Compose |

### Optionaler Backend-Teil (Go)

| Feature | Status | Details |
|---------|--------|---------|
| API-Proxy | ✅ | Kapselt externe Wettstar API vollständig |
| In-Memory Caching | ✅ | TTL: 30s, Thread-sicher mit RWMutex + Goroutine-Cleanup |
| Request-Validierung | ✅ | Pflichtfelder geprüft vor upstream-Call |
| Rate Limiting | ✅ | 60 req/min/IP mit automatischer Bereinigung |

---

## 📡 API-Dokumentation

### Externe Wettstar API (upstream)

| Method | Endpoint | Beschreibung |
|--------|----------|--------------|
| GET | `/ticker/?tickernames=wpftest` | Alle Einträge laden |
| POST | `/ticker/` | Neuen Eintrag erstellen |
| PATCH | `/ticker/` | Eintrag bearbeiten |
| DELETE | `/ticker/{id}/` | Eintrag löschen |

### Go Backend API (intern)

| Method | Endpoint | Beschreibung |
|--------|----------|--------------|
| GET | `/api/ticker` | Einträge laden (mit Cache) |
| POST | `/api/ticker` | Eintrag erstellen (mit Validierung) |
| PATCH | `/api/ticker` | Eintrag bearbeiten (Cache invalidieren) |
| DELETE | `/api/ticker/{id}` | Eintrag löschen (Cache invalidieren) |

### Response Headers

| Header | Wert | Bedeutung |
|--------|------|-----------|
| X-Cache | HIT | Antwort aus In-Memory Cache (< 1ms) |
| X-Cache | MISS | Antwort direkt von upstream API |

---

## 🤖 KI-Unterstützung

### Genutzte Tools

- **Claude (Anthropic)** — als Mentor, Lernbegleiter und Code-Assistent

### Wofür KI genutzt wurde

**Architektur & Konzepte:**
- Erklärung von Vue 3 Composition API Konzepten (ref, reactive, computed, watch)
- Erklärung von Pinia Store-Architektur — wann Getters vs. Actions sinnvoll sind
- Erklärung von Go-Konzepten (Goroutines, sync.RWMutex, HTTP-Handler)
- Diskussion der Architekturentscheidungen (warum Go-Backend, warum Caching)

**UI & Design:**
- Tailwind CSS Klassen für responsives Layout und Dark Mode
- Komponenten-Struktur (Modal mit Teleport, Pagination-Buttons)
- Farbschema und visuelle Hierarchie der Ticker-Karten

**Code-Generierung:**
- Grundstruktur der Go-Dateien (cache.go, ratelimit.go, handler.go)
- Grundstruktur der Vue-Komponenten
- Docker- und Nginx-Konfigurationen

### Was ich selbst erarbeitet und angepasst habe

**Analyse & Entscheidungen:**
- Die externe API selbst getestet und die echte JSON-Struktur analysiert
- TypeScript-Interfaces aus der echten API-Antwort abgeleitet — nicht geraten
- Entschieden welche Zusatzfeatures sinnvoll sind und wie sie zusammenspielen
- Architekturentscheidung: Go-Backend mit Caching statt direkter API-Anbindung

**Debugging & Problemlösung:**
- Docker Build-Fehler (Go-Version 1.26.2 vs. 1.22 im Image) eigenständig analysiert
- Nginx-Proxy-Konfiguration debuggt (BASE_URL auf relativen Pfad geändert)
- Healthcheck-Problem mit distroless Image identifiziert und gelöst
- GitHub-Authentifizierungsproblem eigenständig behoben

**Inhaltliches Verständnis:**
- Jeden Schritt verstanden bevor weitergemacht wurde
- Bei unklaren Konzepten (Vue Routing, Pinia, Go Concurrency) gezielt nachgefragt
  und erst nach dem Verständnis weitergemacht
- Code nicht blind übernommen — jeden Block gelesen, verstanden und bei Bedarf
  angepasst

### Fazit zur KI-Nutzung

KI wurde als intelligentes Werkzeug und Mentor eingesetzt — nicht als Ersatz
für eigenes Denken. Die Architekturentscheidungen, das Debugging, das Verständnis
der Konzepte und die Problemlösung lagen bei mir. KI hat die Umsetzungsgeschwindigkeit
erhöht, nicht die Denkarbeit ersetzt — genau wie es die Aufgabenstellung beschreibt.

---

## 🔮 Verbesserungsmöglichkeiten

### Kurzfristig

1. **Tests** — Unit-Tests für den Pinia Store (Vitest) und Go-Handler (testing package)
2. **Optimistic Updates** — UI sofort aktualisieren, bei API-Fehler zurückrollen
3. **Umgebungsvariablen** — API-URL und Cache-TTL über `.env` konfigurierbar

### Mittelfristig

4. **WebSocket** — Echtzeit-Updates statt Polling alle 10 Sekunden
5. **CI/CD Pipeline** — GitHub Actions für Tests, Docker-Build und Deployment
6. **HTTPS** — SSL-Zertifikat mit Let's Encrypt

### Langfristig

7. **Backend Persistenz** — Eigene Datenbank im Go-Backend für Offline-Modus
   und bessere Businesslogik
8. **Kubernetes Deployment** — Helm Chart für skalierbares Deployment
   auf einem Kubernetes Cluster (EKS)

---

## 👨‍💻 Über den Entwickler

**Cyril Kamgais Totso**  
M.Sc. Informatik (Cloud Computing & AWS) — TH Brandenburg  
3+ Jahre Erfahrung in Full-Stack-Entwicklung und Cloud-Architektur

| Kontakt | Link |
|---------|------|
| Email | cyrilkamgais1203@gmail.com |
| Telefon | +49 176 68484380 |
| LinkedIn | [cyril-kamgais-totso-a86491207](https://linkedin.com/in/cyril-kamgais-totso-a86491207) |
| GitHub | [Kamgais](https://github.com/Kamgais) |
| Portfolio | [personal-portfolio-v2-delta.vercel.app](https://personal-portfolio-v2-delta.vercel.app) |