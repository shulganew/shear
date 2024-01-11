package router

import (
	"context"
	"database/sql"
	"net/http"
	"sync"

	"github.com/go-chi/chi/v5"
	"github.com/shulganew/shear.git/internal/config"
	"github.com/shulganew/shear.git/internal/middlewares"
	"github.com/shulganew/shear.git/internal/service"
	"github.com/shulganew/shear.git/internal/web/handlers"
)

// Chi Router for application
func RouteShear(conf *config.Config, stor *service.StorageURL, db *sql.DB, delete *service.Delete, finalCh chan service.DelBatch, waitDel *sync.WaitGroup) (r *chi.Mux) {

	webHand := handlers.NewHandlerWeb(conf, stor)
	r = chi.NewRouter()

	//send password for enctription to middlewares
	r.Use(func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			ctx := context.WithValue(r.Context(), config.CtxPassKey{}, conf.Pass)
			h.ServeHTTP(w, r.WithContext(ctx))
		})
	})

	r.Route("/", func(r chi.Router) {

		r.Use(middlewares.MidlewLog)
		r.Use(middlewares.MidlewZip)
		r.Use(middlewares.Auth)
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

		handCookieID := handlers.NewHandlerAuthUser(conf, stor)
		r.Get("/api/user/urls", http.HandlerFunc(handCookieID.GetUserURLs))

		delID := handlers.NewHandlerDelShorts(delete)
		r.Delete("/api/user/urls", http.HandlerFunc(delID.DelUserURLs))
	})

	return
}
