package routes

import (
	"net/http"
	"strconv"

	gogame "github.com/YoDobchev/Go-Online/src/game/go"
	"github.com/go-chi/chi/v5"
)

func BlogsRoutes() *chi.Mux {
	r := chi.NewRouter()
	r.Get("/", getBlogsHandler)
	r.Get("/{id}", getBlogByIDHandler)

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
