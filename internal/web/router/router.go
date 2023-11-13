package router

import (
	"github.com/go-chi/chi/v5"
	webhandl "github.com/shulganew/shear.git/internal/web/handlers"
)

// Chi Router for application
func RouteShear(hadler webhandl.URLHandler) *chi.Mux {
	r := chi.NewRouter()
	r.Get("/{id}", hadler.GetURL)
	r.Post("/", hadler.SetUrl)

	return r
}
