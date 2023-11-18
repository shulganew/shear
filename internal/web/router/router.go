package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/shulganew/shear.git/internal/web/handlers"
	"go.uber.org/zap"
)

// Chi Router for application
func RouteShear(handler handlers.URLHandler, appLogger zap.SugaredLogger) (r *chi.Mux) {
	r = chi.NewRouter()
	r.Get("/{id}", handlers.MidlewLog(http.HandlerFunc(handler.GetURL), appLogger))
	r.Post("/", handlers.MidlewLog(http.HandlerFunc(handler.SetURL), appLogger))

	return
}
