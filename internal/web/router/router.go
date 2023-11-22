package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/shulganew/shear.git/internal/api"
	"github.com/shulganew/shear.git/internal/config"
	"github.com/shulganew/shear.git/internal/web/handlers"
)

// Chi Router for application
func RouteShear(conf *config.Shear) (r *chi.Mux) {

	webHand := handlers.NewHandler(conf)
	r = chi.NewRouter()
	r.Use(handlers.MidlewLog)
	r.Post("/", http.HandlerFunc(webHand.SetURL))
	r.Get("/{id}", http.HandlerFunc(webHand.GetURL))

	//api
	apiHand := api.NewHandler(conf)
	r.Route("/api/shorten", func(r chi.Router) {
		r.Use(handlers.MidlewLog)
		r.Use(handlers.MidlewZip)
		r.Post("/", http.HandlerFunc(apiHand.SetAPI))

	})

	return
}
