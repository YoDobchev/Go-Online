package routes

import (
	"encoding/json"
	"net/http"

	gogame "github.com/YoDobchev/Go-Online/src/game/go"
	"github.com/YoDobchev/Go-Online/src/middleware"
	"github.com/go-chi/chi/v5"
)

func GamesRoutes() *chi.Mux {
	r := chi.NewRouter()

	r.Get("/{id}", getGameHandler)
	r.Group(func(r chi.Router) {
		r.Use(middleware.IsLoggedIn)

		r.Get("/", getIfInGameHandler)
		r.Post("/", postNewGameHandler)
	})

	return r
}

func getGameHandler(w http.ResponseWriter, r *http.Request) {
	gameID := chi.URLParam(r, "id")

	game, exists := gogame.GameInstances[gameID]
	if !exists {
		http.Error(w, "game not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(struct {
		Players [2]string     `json:"players"`
		Turn    uint8         `json:"turn"`
		Board   *gogame.Board `json:"board"`
	}{
		Players: game.Players,
		Turn:    game.CurrectTurn,
		Board:   game.Board,
	})
}

func getIfInGameHandler(w http.ResponseWriter, r *http.Request) {
	user, ok := middleware.GetUserFromCtx(r)
	if !ok {
		http.Error(w, "user missing", http.StatusUnauthorized)
		return
	}

	game, inGame := gogame.PlayerToGame[user.Username]
	w.Header().Set("Content-Type", "application/json")
	if inGame {
		json.NewEncoder(w).Encode(map[string]string{
			"id": game.ID,
		})
	} else {
		json.NewEncoder(w).Encode(map[string]string{
			"id": "",
		})
	}
}

func postNewGameHandler(w http.ResponseWriter, r *http.Request) {
	user, ok := middleware.GetUserFromCtx(r)
	if !ok {
		http.Error(w, "user missing", http.StatusUnauthorized)
		return
	}

	newGame, err := gogame.NewGame(19, user.Username)
	if err != nil {
		http.Error(w, "could not create game", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"id": newGame.ID,
	})
}
