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

type TickerHandler struct {
	cache      *cache.Cache
	httpClient *http.Client
}

func New(c *cache.Cache) *TickerHandler {
	return &TickerHandler{
		cache: c,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

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

	// Nur aktive Einträge zurückgeben
	var allEntries []map[string]interface{}
	if err := json.Unmarshal(body, &allEntries); err == nil {
		activeEntries := []map[string]interface{}{}
		for _, entry := range allEntries {
			if active, ok := entry["active"].(bool); ok && active {
				activeEntries = append(activeEntries, entry)
			}
		}
		filteredBody, _ := json.Marshal(activeEntries)
		h.cache.Set(cacheKey, filteredBody)

		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("X-Cache", "MISS")
		w.Write(filteredBody)
		return
	}

	// Fallback: ungefiltert zurückgeben
	h.cache.Set(cacheKey, body)
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("X-Cache", "MISS")
	w.Write(body)
}

func (h *TickerHandler) CreateEntry(w http.ResponseWriter, r *http.Request) {
	var payload map[string]interface{}

	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, `{"error": "Ungültiges JSON"}`, http.StatusBadRequest)
		return
	}

	if err := validateCreatePayload(payload); err != nil {
		http.Error(w, fmt.Sprintf(`{"error": "%s"}`, err.Error()), http.StatusBadRequest)
		return
	}

	// ticker_name immer setzen
	payload["ticker_name"] = tickerName

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

	h.cache.Invalidate(fmt.Sprintf("ticker:%s", tickerName))

	respBody, _ := io.ReadAll(resp.Body)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(resp.StatusCode)
	w.Write(respBody)
}

func (h *TickerHandler) UpdateEntry(w http.ResponseWriter, r *http.Request) {
	var payload map[string]interface{}

	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, `{"error": "Ungültiges JSON"}`, http.StatusBadRequest)
		return
	}

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

	h.cache.Invalidate(fmt.Sprintf("ticker:%s", tickerName))

	respBody, _ := io.ReadAll(resp.Body)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(resp.StatusCode)
	w.Write(respBody)
}

// DeleteEntry — simuliert Loeschen via PATCH mit active: false
// Die externe API unterstuetzt kein DELETE
func (h *TickerHandler) DeleteEntry(w http.ResponseWriter, r *http.Request) {
	tickerID := r.PathValue("id")
	if tickerID == "" {
		http.Error(w, `{"error": "ID fehlt"}`, http.StatusBadRequest)
		return
	}

	// Eintrag aus Cache oder API laden
	cacheKey := fmt.Sprintf("ticker:%s", tickerName)
	var entries []map[string]interface{}

	if cached, found := h.cache.Get(cacheKey); found {
		json.Unmarshal(cached, &entries)
	} else {
		url := fmt.Sprintf("%s/ticker/?tickernames=%s", upstreamURL, tickerName)
		resp, err := h.httpClient.Get(url)
		if err != nil {
			http.Error(w, `{"error": "Upstream nicht erreichbar"}`, http.StatusBadGateway)
			return
		}
		defer resp.Body.Close()
		body, _ := io.ReadAll(resp.Body)
		json.Unmarshal(body, &entries)
	}

	// Eintrag anhand der ID finden
	var target map[string]interface{}
	for _, entry := range entries {
		id, _ := entry["ticker_id"].(float64)
		if fmt.Sprintf("%d", int(id)) == tickerID {
			target = entry
			break
		}
	}

	if target == nil {
		http.Error(w, `{"error": "Eintrag nicht gefunden"}`, http.StatusNotFound)
		return
	}

	// active: false setzen — alle Pflichtfelder mitschicken
	target["active"] = false

	body, _ := json.Marshal(target)
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
	h.cache.Invalidate(cacheKey)

	log.Printf("[DELETE] Eintrag %s auf active:false gesetzt", tickerID)
	w.WriteHeader(http.StatusNoContent)
}

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