package routes

import (
	"encoding/json"
	"net/http"
	"os"
	"time"

	"github.com/YoDobchev/Go-Online/src/database"
	gogame "github.com/YoDobchev/Go-Online/src/game/go"
	"github.com/YoDobchev/Go-Online/src/middleware"
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
	Identifier string `json:"identifier"`
	Password   string `json:"password"`
}

func AuthRoutes() *chi.Mux {
	r := chi.NewRouter()

	r.Post("/login", loginHandler)
	r.Post("/register", registerHandler)

	r.Group(func(r chi.Router) {
		r.Use(middleware.IsLoggedIn)

		r.Get("/me", meHandler)
		r.Delete("/logout", logoutHandler)
	})

	return r
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	var req LoginReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}

	var user database.User
	if err := database.DB.Where("username = ? OR email = ?", req.Identifier, req.Identifier).First(&user).Error; err != nil {
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
		UserID:    user.ID,
		Token:     sessionToken,
		ExpiresAt: expiry,
	}

	if err := database.DB.Create(&session).Error; err != nil {
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

	writeJSON(w, http.StatusOK, map[string]string{
		"message": "login successful",
	})
}
func registerHandler(w http.ResponseWriter, r *http.Request) {
	var req RegisterReq

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}

	var count int64
	if err := database.DB.Model(&database.User{}).
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

	if err := database.DB.Create(&newUser).Error; err != nil {
		http.Error(w, "db create err", http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"message": "register successful"})
}

func meHandler(w http.ResponseWriter, r *http.Request) {
	user, ok := middleware.GetUserFromCtx(r)
	if !ok {
		http.Error(w, "user missing", http.StatusUnauthorized)
		return
	}

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
	}

	writeJSON(w, http.StatusOK, resp)
}

func logoutHandler(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session")
	if err != nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	sessionToken := cookie.Value

	if err := database.DB.Where("token = ?", sessionToken).Delete(&database.Session{}).Error; err != nil {
		http.Error(w, "could not logout", http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "session",
		Value:    "",
		Path:     "/",
		Expires:  time.Unix(0, 0),
		MaxAge:   -1,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Secure:   os.Getenv("ENV") == "prod",
	})

	writeJSON(w, http.StatusOK, map[string]string{
		"message": "logout successful",
	})
}
