package routes

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/YoDobchev/Go-Online/src/database"
	"github.com/YoDobchev/Go-Online/src/elo"
	gogame "github.com/YoDobchev/Go-Online/src/game/go"
	"github.com/YoDobchev/Go-Online/src/middleware"
)

func UsersRoutes() *chi.Mux {
	r := chi.NewRouter()

	r.Get("/{username}/elo", getEloForUserHandler)

	r.Group(func(r chi.Router) {
		r.Use(middleware.IsModerator)
		r.Get("/", getAllUsersHandler)
		r.Get("/{username}", getUserHandler)
		r.Get("/", getAllUsersHandler)
		r.Delete("/{username}", deleteUserHandler)
	})

	return r
}

func getUserHandler(w http.ResponseWriter, r *http.Request) {
	username := chi.URLParam(r, "username")
	if username == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "missing username"})
		return
	}

	var user database.User
	if err := database.DB.Where("username = ?", username).First(&user).Error; err != nil {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "user not found"})
		return
	}

	rank, _ := elo.GetRank(user.Elo)

	game, inGame := gogame.PlayerToGame[user.Username]

	var gameID string
	if inGame && game != nil {
		gameID = game.ID
	} else {
		gameID = ""
	}

	resp := map[string]any{
		"email":          user.Email,
		"username":       user.Username,
		"isInGameWithID": gameID,
		"elo":            user.Elo,
		"role":           user.Role,
		"rank":           rank,
	}

	writeJSON(w, http.StatusOK, resp)
}

func getEloForUserHandler(w http.ResponseWriter, r *http.Request) {
	username := chi.URLParam(r, "username")
	if username == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "missing username"})
		return
	}

	var user database.User
	if err := database.DB.Where("username = ?", username).First(&user).Error; err != nil {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "user not found"})
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"elo": user.Elo,
	})
}

func deleteUserHandler(w http.ResponseWriter, r *http.Request) {
	username := chi.URLParam(r, "username")
	if username == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "missing username"})
		return
	}

	var user database.User
	if err := database.DB.Where("username = ?", username).First(&user).Error; err != nil {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "user not found"})
		return
	}

	if err := database.DB.Delete(&user).Error; err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "could not delete user"})
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"message": "user deleted"})
}

func getAllUsersHandler(w http.ResponseWriter, r *http.Request) {
	var users []database.User

	if err := database.DB.Find(&users).Error; err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{
			"error": "could not load users",
		})
		return
	}

	resp := make([]map[string]any, 0, len(users))

	for _, user := range users {
		rank, _ := elo.GetRank(user.Elo)

		resp = append(resp, map[string]any{
			"id":       user.ID,
			"email":    user.Email,
			"username": user.Username,
			"elo":      user.Elo,
			"role":     user.Role,
			"rank":     rank,
		})
	}

	writeJSON(w, http.StatusOK, resp)
}
