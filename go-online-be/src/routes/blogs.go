package routes

import (
	"encoding/json"
	"net/http"
	"strconv"

	gogame "github.com/YoDobchev/Go-Online/src/game/go"
	"github.com/YoDobchev/Go-Online/src/middleware"
	"github.com/go-chi/chi/v5"
)

func BlogsRoutes() *chi.Mux {
	r := chi.NewRouter()
	r.Get("/", getBlogsHandler)
	r.Get("/{id}", getBlogByIDHandler)

	r.Group(func(r chi.Router) {
		r.Use(middleware.IsLoggedIn)
		r.Post("/", postNewBlogHandler)
	})

	return r
}

func getBlogsHandler(w http.ResponseWriter, r *http.Request) {
	blogs, err := gogame.LoadAllBlogsFromDB()
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, blogs)
}

func getBlogByIDHandler(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	blogID, err := strconv.Atoi(id)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid blog id"})
		return
	}

	blog, err := gogame.LoadBlogFromDB(blogID)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, blog)
}

type createBlogReq struct {
	Id       int    `json:"id,omitempty"`
	AuthorId int    `json:"authorId,omitempty"`
	Title    string `json:"title"`
	Content  string `json:"content"`
}

func postNewBlogHandler(w http.ResponseWriter, r *http.Request) {
	user, ok := middleware.GetUserFromCtx(r)
	if !ok {
		http.Error(w, "user missing", http.StatusUnauthorized)
		return
	}

	if user.Role != "admin" {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}

	var req createBlogReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}

	err := gogame.CreateBlog(req.Id, req.AuthorId, req.Title, req.Content)
	if err != nil {
		http.Error(w, "could not create blog", http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"id": strconv.Itoa(req.Id),
	})
}
