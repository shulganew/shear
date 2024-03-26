package server

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/shulganew/shear.git/internal/config"
	"go.uber.org/zap"
)

const timeoutServerShutdown = time.Second * 5

// Manage web server.
func ShortenerServer(ctx context.Context, conf *config.Config, componentsErrs chan error, r *chi.Mux) (webDone chan struct{}) {
	// Start web server.
	var srv = http.Server{Addr: conf.GetAddress(), Handler: r}
	go func() {
		// Public certificate: server.crt
		//
		// Private key: server.pem
		if conf.IsSecure() {
			if err := srv.ListenAndServeTLS("./cert/server.crt", "./cert/server.pem"); err != nil {
				if errors.Is(err, http.ErrServerClosed) {
					return
				}
				componentsErrs <- fmt.Errorf("listen and server has failed: %w", err)
			}
		} else {
			if err := srv.ListenAndServe(); err != nil {
				if errors.Is(err, http.ErrServerClosed) {
					return
				}
				componentsErrs <- fmt.Errorf("listen and server has failed: %w", err)
			}
		}
	}()

	// Graceful shutdown.
	webDone = make(chan struct{})
	go func() {
		defer zap.S().Infoln("Server web has been graceful shutdown.")
		defer close(webDone)
		<-ctx.Done()
		// Wait until all del async short will be saved.
		shutdownTimeoutCtx, cancelShutdownTimeoutCtx := context.WithTimeout(context.Background(), timeoutServerShutdown)
		defer cancelShutdownTimeoutCtx()
		if err := srv.Shutdown(shutdownTimeoutCtx); err != nil {
			zap.S().Infoln("an error occurred during server shutdown: %v", err)
		}
	}()
	return
}
