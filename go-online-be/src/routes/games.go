package routes

import (
	"encoding/json"
	"net/http"

	gogame "github.com/YoDobchev/Go-Online/src/game/go"
	"github.com/YoDobchev/Go-Online/src/middleware"
	ws "github.com/YoDobchev/Go-Online/src/websocket"
	"github.com/go-chi/chi/v5"
)

func GamesRoutes() *chi.Mux {
	r := chi.NewRouter()

	r.Get("/{id}/ws", ws.WsGameHandler)

	r.Group(func(r chi.Router) {
		r.Use(middleware.IsLoggedIn)

		r.Get("/", getIfInGameHandler)
		r.Post("/", postNewGameHandler)
	})

	return r
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
