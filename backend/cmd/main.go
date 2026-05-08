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
	c := cache.New(30 * time.Second)
	tickerHandler := handler.New(c)
	limiter := middleware.NewRateLimiter(60, time.Minute)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/ticker", tickerHandler.GetEntries)
	mux.HandleFunc("POST /api/ticker", tickerHandler.CreateEntry)
	mux.HandleFunc("PATCH /api/ticker", tickerHandler.UpdateEntry)
	mux.HandleFunc("DELETE /api/ticker/{id}", tickerHandler.DeleteEntry)

	c2 := cors.New(cors.Options{
		AllowedOrigins: []string{"http://localhost:5173"},
		AllowedMethods: []string{"GET", "POST", "PATCH", "DELETE"},
		AllowedHeaders: []string{"Content-Type"},
	})

	handler := limiter.Middleware(c2.Handler(mux))

	log.Println("Backend läuft auf http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", handler))
}