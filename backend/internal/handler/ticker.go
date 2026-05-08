package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"ticker-backend/internal/cache"
	"time"
)

const (
	upstreamURL = "https://65q02hjc6c.execute-api.eu-central-1.amazonaws.com/prod"
	tickerName  = "wpftest"
)

// TickerHandler verwaltet alle Ticker-Anfragen
type TickerHandler struct {
	cache      *cache.Cache
	httpClient *http.Client
}

// New erstellt einen neuen TickerHandler
func New(c *cache.Cache) *TickerHandler {
	return &TickerHandler{
		cache: c,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// GetEntries gibt alle Ticker-Einträge zurück (mit Caching)
func (h *TickerHandler) GetEntries(w http.ResponseWriter, r *http.Request) {
	cacheKey := fmt.Sprintf("ticker:%s", tickerName)

	// Cache prüfen
	if cached, found := h.cache.Get(cacheKey); found {
		log.Println("[CACHE HIT] GET /ticker")
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("X-Cache", "HIT")
		w.Write(cached)
		return
	}

	// Upstream API aufrufen
	log.Println("[CACHE MISS] GET /ticker — rufe upstream API")
	url := fmt.Sprintf("%s/ticker/?tickernames=%s", upstreamURL, tickerName)
	resp, err := h.httpClient.Get(url)
	if err != nil {
		http.Error(w, `{"error": "Upstream nicht erreichbar"}`, http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, `{"error": "Fehler beim Lesen der Antwort"}`, http.StatusInternalServerError)
		return
	}

	// Im Cache speichern
	h.cache.Set(cacheKey, body)

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("X-Cache", "MISS")
	w.Write(body)
}

// CreateEntry erstellt einen neuen Eintrag mit Validierung
func (h *TickerHandler) CreateEntry(w http.ResponseWriter, r *http.Request) {
	var payload map[string]interface{}

	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, `{"error": "Ungültiges JSON"}`, http.StatusBadRequest)
		return
	}

	// Validierung
	if err := validateCreatePayload(payload); err != nil {
		http.Error(w, fmt.Sprintf(`{"error": "%s"}`, err.Error()), http.StatusBadRequest)
		return
	}

	// ticker_name immer setzen
	payload["ticker_name"] = tickerName

	// An upstream weiterleiten
	body, _ := json.Marshal(payload)
	resp, err := h.httpClient.Post(
		fmt.Sprintf("%s/ticker/", upstreamURL),
		"application/json",
		bytes.NewBuffer(body),
	)
	if err != nil {
		http.Error(w, `{"error": "Upstream nicht erreichbar"}`, http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	// Cache invalidieren — neue Daten vorhanden
	h.cache.Invalidate(fmt.Sprintf("ticker:%s", tickerName))

	respBody, _ := io.ReadAll(resp.Body)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(resp.StatusCode)
	w.Write(respBody)
}

// UpdateEntry bearbeitet einen Eintrag
func (h *TickerHandler) UpdateEntry(w http.ResponseWriter, r *http.Request) {
	var payload map[string]interface{}

	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, `{"error": "Ungültiges JSON"}`, http.StatusBadRequest)
		return
	}

	// ticker_id muss vorhanden sein
	if _, ok := payload["ticker_id"]; !ok {
		http.Error(w, `{"error": "ticker_id fehlt"}`, http.StatusBadRequest)
		return
	}

	body, _ := json.Marshal(payload)
	req, _ := http.NewRequest(http.MethodPatch,
		fmt.Sprintf("%s/ticker/", upstreamURL),
		bytes.NewBuffer(body),
	)
	req.Header.Set("Content-Type", "application/json")

	resp, err := h.httpClient.Do(req)
	if err != nil {
		http.Error(w, `{"error": "Upstream nicht erreichbar"}`, http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	// Cache invalidieren
	h.cache.Invalidate(fmt.Sprintf("ticker:%s", tickerName))

	respBody, _ := io.ReadAll(resp.Body)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(resp.StatusCode)
	w.Write(respBody)
}

// DeleteEntry löscht einen Eintrag
func (h *TickerHandler) DeleteEntry(w http.ResponseWriter, r *http.Request) {
	tickerID := r.PathValue("id")
	if tickerID == "" {
		http.Error(w, `{"error": "ID fehlt"}`, http.StatusBadRequest)
		return
	}

	req, _ := http.NewRequest(http.MethodDelete,
		fmt.Sprintf("%s/ticker/%s/", upstreamURL, tickerID),
		nil,
	)

	resp, err := h.httpClient.Do(req)
	if err != nil {
		http.Error(w, `{"error": "Upstream nicht erreichbar"}`, http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	// Cache invalidieren
	h.cache.Invalidate(fmt.Sprintf("ticker:%s", tickerName))

	w.WriteHeader(resp.StatusCode)
}

// validateCreatePayload prüft Pflichtfelder
func validateCreatePayload(payload map[string]interface{}) error {
	required := []string{"title", "message", "creator"}
	for _, field := range required {
		val, ok := payload[field]
		if !ok {
			return fmt.Errorf("Feld '%s' fehlt", field)
		}
		str, ok := val.(string)
		if !ok || str == "" {
			return fmt.Errorf("Feld '%s' darf nicht leer sein", field)
		}
	}
	return nil
}