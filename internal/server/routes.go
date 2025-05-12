package server

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	db "tanggalan-api/internal/database"
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

	r.Get("/api", s.getAnnualEventsByMonthEventsHandler(s.q))
	r.HandleFunc("/api/sync", s.syncEventsHandler(s.q, s.dbConn))

	return r
}

func (s *Server) getAnnualEventsByMonthEventsHandler(q *db.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		monthParam := r.URL.Query().Get("month")
		yearParam := r.URL.Query().Get("year")

		if monthParam == "" || yearParam == "" {
			http.Error(w, "month and year query params are required", http.StatusBadRequest)
			return
		}

		// Format jadi 2 digit (e.g. "04" bukan "4")
		monthInt, err := strconv.Atoi(monthParam)
		if err != nil || monthInt < 1 || monthInt > 12 {
			http.Error(w, "invalid month", http.StatusBadRequest)
			return
		}
		monthStr := fmt.Sprintf("%02d", monthInt)

		// Get events from DB
		events, err := q.GetEventsByMonthAndYear(r.Context(), db.GetEventsByMonthAndYearParams{
			Year:  yearParam,
			Month: monthStr,
		})
		if err != nil {
			http.Error(w, "failed to get events", http.StatusInternalServerError)
			return
		}

		// Convert to JSON
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(events)
	}
}

func (s *Server) syncEventsHandler(q *db.Queries, dbConn *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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

		status, err := scraper.SyncEventsFromTanggalan(r.Context(), q, monthName, year)
		if err != nil {
			http.Error(w, "failed to sync events", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(status))
	}
}
