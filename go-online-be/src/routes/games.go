package routes

import (
	"net/http"
	"strconv"

	gogame "github.com/YoDobchev/Go-Online/src/game/go"
	"github.com/YoDobchev/Go-Online/src/middleware"
	ws "github.com/YoDobchev/Go-Online/src/websocket"
	"github.com/go-chi/chi/v5"
)

func GamesRoutes() *chi.Mux {
	r := chi.NewRouter()

	r.Get("/{id}/ws", ws.WsGameHandler)
	r.Get("/{id}/state", getGameStateHandler)

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
	if inGame {
		writeJSON(w, http.StatusOK, map[string]string{
			"id": game.ID,
		})
	} else {
		writeJSON(w, http.StatusOK, map[string]string{
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

	writeJSON(w, http.StatusOK, map[string]string{
		"id": newGame.ID,
	})
}

func getGameStateHandler(w http.ResponseWriter, r *http.Request) {
	gameID := chi.URLParam(r, "id")
	game, exists := gogame.GameInstances[gameID]
	if !exists {
		http.Error(w, "game not found", http.StatusNotFound)
		return
	}
	moveNum := r.URL.Query().Get("moveNum")
	if moveNum == "" {
		http.Error(w, "enter a valid move number", http.StatusInternalServerError)
		return
	}

	moveNumi, err := strconv.Atoi(moveNum)
	if err != nil {
		http.Error(w, "enter a valid move number", http.StatusInternalServerError)
		return
	}

	board, err := gogame.GetBoardStateOnMoveNoFromDB(game.ID, (moveNumi))
	if err != nil {
		http.Error(w, "could not get board state", http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, board)
}
