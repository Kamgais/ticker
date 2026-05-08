package middleware

import (
	"net/http"
	"sync"
	"time"
)

// RateLimiter begrenzt Anfragen pro IP
type RateLimiter struct {
	mu       sync.Mutex
	requests map[string][]time.Time
	limit    int
	window   time.Duration
}

// NewRateLimiter erstellt einen neuen RateLimiter
func NewRateLimiter(limit int, window time.Duration) *RateLimiter {
	return &RateLimiter{
		requests: make(map[string][]time.Time),
		limit:    limit,
		window:   window,
	}
}

// Allow prüft ob eine IP noch Anfragen machen darf
func (rl *RateLimiter) Allow(ip string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	windowStart := now.Add(-rl.window)

	// Alte Anfragen entfernen
	filtered := []time.Time{}
	for _, t := range rl.requests[ip] {
		if t.After(windowStart) {
			filtered = append(filtered, t)
		}
	}

	rl.requests[ip] = filtered

	// Limit prüfen
	if len(filtered) >= rl.limit {
		return false
	}

	rl.requests[ip] = append(rl.requests[ip], now)
	return true
}

// Middleware gibt einen HTTP-Handler zurück
func (rl *RateLimiter) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := r.RemoteAddr

		if !rl.Allow(ip) {
			http.Error(w, `{"error": "Too many requests"}`, http.StatusTooManyRequests)
			return
		}

		next.ServeHTTP(w, r)
	})
}