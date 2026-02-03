package middleware

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/YoDobchev/Go-Online/src/database"
)

type ctxKey string

const userCtxKey ctxKey = "sessionUser"

func IsLoggedIn(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, err := GetUserInfo(r)
		if err != nil {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), userCtxKey, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func GetUserFromCtx(r *http.Request) (*database.User, bool) {
	user, ok := r.Context().Value(userCtxKey).(*database.User)
	return user, ok
}

func GetUserInfo(r *http.Request) (*database.User, error) {
	cookie, err := r.Cookie("session")
	if err != nil {
		return nil, err
	}

	sessionToken := cookie.Value

	var session database.Session
	if err := database.DB.
		Preload("User").
		Where("token = ?", sessionToken).
		First(&session).Error; err != nil {
		return nil, err
	}

	if time.Now().After(session.ExpiresAt) {
		return nil, fmt.Errorf("session expired")
	}

	return &session.User, nil
}
