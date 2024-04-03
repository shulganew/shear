package router

import (
	"context"
	"database/sql"
	"net/http"
	"strings"

	"net/http/pprof"

	"github.com/go-chi/chi/v5"
	"github.com/shulganew/shear.git/internal/config"
	"github.com/shulganew/shear.git/internal/service"
	"github.com/shulganew/shear.git/internal/web/handlers"
	"github.com/shulganew/shear.git/internal/web/middlewares"
	"github.com/shulganew/shear.git/internal/web/validators"
	httpSwagger "github.com/swaggo/http-swagger"
)

// Chi Router for application.
func RouteShear(conf *config.Config, short *service.Shorten, db *sql.DB, delete *service.Delete) (r *chi.Mux) {
	r = chi.NewRouter()

	// Send password and ip/mask trusted network to middlewares.
	r.Use(func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := context.WithValue(r.Context(), config.CtxPassKey{}, conf.GetPass())
			ctx = context.WithValue(ctx, config.CtxIP{}, conf.GetIP())
			h.ServeHTTP(w, r.WithContext(ctx))
		})
	})

	r.Route("/", func(r chi.Router) {
		r.Use(middlewares.MiddlwLog)
		r.Use(middlewares.MiddlwZip)
		r.Use(middlewares.Auth)

		// Set short from URL.
		webHand := handlers.NewHandlerGetURL(conf, short)
		r.Post("/", http.HandlerFunc(webHand.SetURL))
		// Get URL by short.
		r.Get("/{id}", http.HandlerFunc(webHand.GetURL))

		// JSON API for shortener.
		apiHand := handlers.NewHandlerAPI(conf, short)
		r.Post("/api/shorten", http.HandlerFunc(apiHand.GetBrief))

		// Database test - ping.
		dbHand := handlers.NewDB(db)
		r.Get("/ping", http.HandlerFunc(dbHand.Ping))

		// DB Postgres Batch request (multiple JSON)
		batchHand := handlers.NewHandlerBatch(conf, short)
		r.Post("/api/shorten/batch", http.HandlerFunc(batchHand.BatchSet))

		// Get all users URLs.
		handCookieID := handlers.NewHandlerAuthUser(conf, short)
		r.Get("/api/user/urls", http.HandlerFunc(handCookieID.GetUserURLs))

		// Batch delete shorts from handlers (bulk postgres delete).
		delID := handlers.NewHandlerDelShorts(delete)
		r.Delete("/api/user/urls", http.HandlerFunc(delID.DelUserURLs))

		// Server statistic.
		stat := handlers.NewHandlerStat(conf, short)
		r.With(middlewares.NetAccess).Get("/api/internal/stats", http.HandlerFunc(stat.GetStat))

		if conf.IsPprof() {
			// Adding pprof.
			r.Get("/debug/pprof/*", http.HandlerFunc(pprof.Index))
			r.Get("/debug/pprof/cmdline", http.HandlerFunc(pprof.Cmdline))
			r.Get("/debug/pprof/profile", http.HandlerFunc(pprof.Profile))
			r.Get("/debug/pprof/symbol", http.HandlerFunc(pprof.Symbol))
			r.Get("/debug/pprof/trace", http.HandlerFunc(pprof.Trace))
		}
		// Add swagger page.
		// Check and parse URL.
		_, startport := validators.CheckURL(conf.GetAddress(), conf.IsSecure())
		r.Get("/swagger/*", httpSwagger.Handler(
			httpSwagger.URL(strings.Join([]string{conf.GetProtocol(), "://", "localhost:", startport, "/swagger/doc.json"}, "")),
		))
	})
	return
}
