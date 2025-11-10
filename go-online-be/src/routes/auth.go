package routes

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/YoDobchev/Go-Online/src/database"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type RegisterReq struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginReq struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func AuthRoutes() *chi.Mux {
	db := database.DB
	r := chi.NewRouter()

	r.Post("/login", func(w http.ResponseWriter, r *http.Request) {
		var req LoginReq
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid JSON", http.StatusBadRequest)
			return
		}

		var user database.User
		if err := db.Where("username = ?", req.Username).First(&user).Error; err != nil {
			http.Error(w, "invalid credentials", http.StatusUnauthorized)
			return
		}

		if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
			http.Error(w, "invalid credentials", http.StatusUnauthorized)
			return
		}

		var secure bool
		if env := os.Getenv("ENV"); env == "prod" {
			secure = true
		}

		sessionToken := uuid.NewString()
		expiry := time.Now().Add(24 * time.Hour)

		session := database.Session{
			UserID:    user.Id,
			Token:     sessionToken,
			ExpiresAt: expiry,
		}

		if err := db.Create(&session).Error; err != nil {
			http.Error(w, "could not create session", http.StatusInternalServerError)
			return
		}

		http.SetCookie(w, &http.Cookie{
			Name:     "session",
			Value:    sessionToken,
			Path:     "/",
			Expires:  expiry,
			HttpOnly: true,
			SameSite: http.SameSiteLaxMode,
			Secure:   secure,
		})

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"message": "login successful",
		})
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
			Email:    req.Email,
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
