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
	r.Post("/", handlers.MidlewLog(http.HandlerFunc(webHand.SetURL), conf.Applog))
	r.Get("/{id}", handlers.MidlewLog(http.HandlerFunc(webHand.GetURL), conf.Applog))

	apiHand := api.NewHandler(conf)
	//api
	r.Post("/api/shorten", handlers.MidlewLog(http.HandlerFunc(apiHand.SetAPI), conf.Applog))

	return
}
