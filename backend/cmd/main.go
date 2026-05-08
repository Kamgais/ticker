package main

import (
	"log"
	"net/http"
	"ticker-backend/internal/cache"
	"ticker-backend/internal/handler"
	"ticker-backend/internal/middleware"
	"time"

	"github.com/rs/cors"
)

func main() {
	// Cache mit 30 Sekunden TTL
	c := cache.New(30 * time.Second)

	// Handler
	tickerHandler := handler.New(c)

	// Rate Limiter: 60 Anfragen pro Minute pro IP
	limiter := middleware.NewRateLimiter(60, time.Minute)

	// Router
	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/ticker", tickerHandler.GetEntries)
	mux.HandleFunc("POST /api/ticker", tickerHandler.CreateEntry)
	mux.HandleFunc("PATCH /api/ticker", tickerHandler.UpdateEntry)
	mux.HandleFunc("DELETE /api/ticker/{id}", tickerHandler.DeleteEntry)

	// CORS — erlaubt Anfragen vom Vue Frontend
	c2 := cors.New(cors.Options{
		AllowedOrigins: []string{"http://localhost:5173"},
		AllowedMethods: []string{"GET", "POST", "PATCH", "DELETE"},
		AllowedHeaders: []string{"Content-Type"},
	})

	// Middleware chain
	handler := limiter.Middleware(c2.Handler(mux))

	log.Println("Backend läuft auf http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", handler))
}