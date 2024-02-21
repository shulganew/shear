package router

import (
	"context"
	"database/sql"
	"net/http"
	"sync"

	"github.com/go-chi/chi/v5"
	"github.com/shulganew/shear.git/internal/config"
	"github.com/shulganew/shear.git/internal/service"
	"github.com/shulganew/shear.git/internal/web/handlers"
	"github.com/shulganew/shear.git/internal/web/middlewares"
)

// Chi Router for application
func RouteShear(conf *config.Config, stor service.StorageURL, db *sql.DB, delete *service.Deleter, finalCh chan service.DelBatch, waitDel *sync.WaitGroup) (r *chi.Mux) {

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

		// Set short from URL.
		webHand := handlers.NewHandlerGetURL(conf, stor)
		r.Post("/", http.HandlerFunc(webHand.SetURL))
		// Get URL by short.
		r.Get("/{id}", http.HandlerFunc(webHand.GetURL))

		// JSON API for shortener.
		apiHand := handlers.NewHandlerAPI(conf, stor)
		r.Post("/api/shorten", http.HandlerFunc(apiHand.GetBrief))

		// Databes teset - ping.
		dbHand := handlers.NewDB(db)
		r.Get("/ping", http.HandlerFunc(dbHand.Ping))

		// DB Postgres Batch request (multiple JSON)
		batchHand := handlers.NewHandlerBatch(conf, stor)
		r.Post("/api/shorten/batch", http.HandlerFunc(batchHand.BatchSet))

		// Get all users URLs.
		handCookieID := handlers.NewHandlerAuthUser(conf, stor)
		r.Get("/api/user/urls", http.HandlerFunc(handCookieID.GetUserURLs))

		// Batch delete shorts from handlers (bulk postgers delete).
		delID := handlers.NewHandlerDelShorts(delete)
		r.Delete("/api/user/urls", http.HandlerFunc(delID.DelUserURLs))
	})

	return
}
