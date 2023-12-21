package router

import (
	"database/sql"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/shulganew/shear.git/internal/config"
	"github.com/shulganew/shear.git/internal/middlewares"
	"github.com/shulganew/shear.git/internal/service"
	"github.com/shulganew/shear.git/internal/web/handlers"
)

// Chi Router for application
func RouteShear(conf *config.Config, stor *service.StorageURL, db *sql.DB) (r *chi.Mux) {

	webHand := handlers.NewHandlerWeb(conf, stor)
	r = chi.NewRouter()
	r.Use(middlewares.MidlewLog)
	r.Use(middlewares.MidlewZip)
	r.Use(middlewares.Cookie)
	r.Post("/", http.HandlerFunc(webHand.SetURL))
	r.Get("/{id}", http.HandlerFunc(webHand.GetURL))

	//api
	apiHand := handlers.NewHandlerAPI(conf, stor)
	r.Post("/api/shorten", http.HandlerFunc(apiHand.GetBrief))

	//DB Postgres Ping
	dbHand := handlers.NewDB(db)
	r.Get("/ping", http.HandlerFunc(dbHand.Ping))

	//DB Postgres Batch request
	batchHand := handlers.NewHandlerBatch(conf, stor)
	r.Post("/api/shorten/batch", http.HandlerFunc(batchHand.BatchSet))

	return
}
