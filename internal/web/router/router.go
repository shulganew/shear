package router

import (
	"github.com/go-chi/chi/v5"
	webhandl "github.com/shulganew/shear.git/internal/web/handlers"
)

// Chi Router for application
func RouteShear() *chi.Mux {
	r := chi.NewRouter()
	r.Get("/{id}", webhandl.GetURL)
	r.Post("/", webhandl.SetUrl)

	return r
}