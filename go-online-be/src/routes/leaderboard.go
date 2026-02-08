package routes

import (
	"net/http"

	gogame "github.com/YoDobchev/Go-Online/src/game/go"
	"github.com/go-chi/chi/v5"
)

func LeaderboardRoutes() *chi.Mux {
	r := chi.NewRouter()
	r.Get("/", getLeaderBoard)

	return r
}

func getLeaderBoard(w http.ResponseWriter, r *http.Request) {
	leaderboard, err := gogame.GetLeaderBoard()
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, leaderboard)
}
