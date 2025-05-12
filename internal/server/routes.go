package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"tanggalan-api/internal/scraper"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

func (s *Server) RegisterRoutes() http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	r.Get("/", s.HelloWorldHandler)

	// r.Get("/health", s.healthHandler)

	r.Get("/api", func(w http.ResponseWriter, r *http.Request) {
		monthParam := r.URL.Query().Get("month")
		yearParam := r.URL.Query().Get("year")

		if monthParam == "" || yearParam == "" {
			http.Error(w, "month and year query params required", http.StatusBadRequest)
			return
		}

		var month, year int
		_, err := fmt.Sscanf(monthParam, "%d", &month)
		_, err2 := fmt.Sscanf(yearParam, "%d", &year)
		if err != nil || err2 != nil {
			http.Error(w, "invalid query params", http.StatusBadRequest)
			return
		}

		monthName, err := scraper.GetMonthName(month)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		events, err := scraper.ScrapeMonthlyEvents(monthName, year)
		if err != nil {
			http.Error(w, fmt.Sprintf("failed to scrape: %v", err), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(events)
	})

	return r
}

func (s *Server) HelloWorldHandler(w http.ResponseWriter, r *http.Request) {
	resp := make(map[string]string)
	resp["message"] = "Hello World"

	jsonResp, err := json.Marshal(resp)
	if err != nil {
		log.Fatalf("error handling JSON marshal. Err: %v", err)
	}

	_, _ = w.Write(jsonResp)
}

// func (s *Server) healthHandler(w http.ResponseWriter, r *http.Request) {
// 	jsonResp, _ := json.Marshal(s.db.Health())
// 	_, _ = w.Write(jsonResp)
// }
