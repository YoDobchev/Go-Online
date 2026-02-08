package routes

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func BlogsRoutes() *chi.Mux {
	r := chi.NewRouter()
	r.Get("/", getBlogsHandler)

	return r
}

func getBlogsHandler(w http.ResponseWriter, r *http.Request) {

}
