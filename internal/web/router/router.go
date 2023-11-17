package router

import (
	"github.com/go-chi/chi/v5"
	"github.com/shulganew/shear.git/internal/web/handlers"
)

// Chi Router for application
func RouteShear(hadler handlers.URLHandler) (r *chi.Mux) {
	r = chi.NewRouter()
	r.Get("/{id}", hadler.GetURL)
	r.Post("/", hadler.SetURL)
	return
}
