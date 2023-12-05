package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/shulganew/shear.git/internal/config"
	"github.com/shulganew/shear.git/internal/middlewares"
	"github.com/shulganew/shear.git/internal/web/handlers"
)

// Chi Router for application
func RouteShear(conf *config.Shear) (r *chi.Mux) {

	webHand := handlers.NewHandlerWeb(conf)
	r = chi.NewRouter()
	r.Use(middlewares.MidlewLog)
	r.Use(middlewares.MidlewZip)
	r.Post("/", http.HandlerFunc(webHand.SetURL))
	r.Get("/{id}", http.HandlerFunc(webHand.GetURL))

	//api
	apiHand := handlers.NewHandlerAPI(conf)
	r.Post("/api/shorten", http.HandlerFunc(apiHand.GetShortURL))

	//DB Postgres
	dbHand := handlers.NewDB(conf)
	r.Get("/ping", http.HandlerFunc(dbHand.Ping))

	return
}
