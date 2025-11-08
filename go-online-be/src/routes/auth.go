package routes

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/YoDobchev/Go-Online/src/database"
	"github.com/go-chi/chi/v5"
	"golang.org/x/crypto/bcrypt"
)

type RegisterReq struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func AuthRoutes() *chi.Mux {
	db := database.DB
	r := chi.NewRouter()

	r.Post("/login", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte("heloo"))

	})

	r.Post("/register", func(w http.ResponseWriter, r *http.Request) {

		fmt.Println("bruuu")
		var req RegisterReq

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid JSON", http.StatusBadRequest)
			return
		}

		var count int64
		if err := db.Model(&database.User{}).
			Where("username = ?", req.Username).
			Count(&count).Error; err != nil {
			http.Error(w, "db error", http.StatusInternalServerError)
			return
		}

		if count > 0 {
			http.Error(w, "already exists", http.StatusConflict)
			return
		}

		hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			http.Error(w, "hash err", http.StatusInternalServerError)
			return
		}

		newUser := database.User{
			Username: req.Username,
			Password: string(hash),
		}

		if err := db.Create(&newUser).Error; err != nil {
			http.Error(w, "db create err", http.StatusInternalServerError)
			return
		}

		response := map[string]string{"message": "register successful"}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	})

	return r
}
