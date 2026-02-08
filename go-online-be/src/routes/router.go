package routes

import (
	"encoding/json"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
)

func New() *chi.Mux {
	r := chi.NewRouter()

	VITE_URL := os.Getenv("VITE_URL")
	if VITE_URL == "" {
		VITE_URL = "http://localhost:5173"
	}

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{VITE_URL},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	r.Mount("/api/auth", AuthRoutes())
	r.Mount("/api/game", GamesRoutes())
	r.Mount("/api/search", SearchRoutes())
	r.Mount("/api/blogs", BlogsRoutes())
	r.Mount("/api/reports", ReportRoutes())
	r.Mount("/api/leaderboard", LeaderboardRoutes())

	return r
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}
