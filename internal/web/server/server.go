package server

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/shulganew/shear.git/internal/config"
	"github.com/shulganew/shear.git/internal/service"
	"github.com/shulganew/shear.git/internal/web/router"
	"go.uber.org/zap"
)

const timeoutServerShutdown = time.Second * 5

// Manage web server.
func ShortenerServer(ctx context.Context, wgroot *sync.WaitGroup, wgdel *sync.WaitGroup, db *sql.DB, conf *config.Config, short *service.Shorten, del *service.Delete, componentsErrs chan error) {
	defer zap.S().Infoln("Server shutdown done.")
	// Start web server.
	var srv = http.Server{Addr: conf.GetAddress(), Handler: router.RouteShear(conf, short, db, del, wgroot)}
	go func() {
		// Public sertificate: server.crt
		//
		// Private key: server.pem
		if conf.IsSequre() {
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
	wgroot.Add(1)
	go func() {
		defer zap.S().Infoln("Server web has been graceful shutdown.")
		defer wgroot.Done()
		<-ctx.Done()
		// Wait until all del async short will be saved.
		wgdel.Wait()
		shutdownTimeoutCtx, cancelShutdownTimeoutCtx := context.WithTimeout(context.Background(), timeoutServerShutdown)
		defer cancelShutdownTimeoutCtx()
		if err := srv.Shutdown(shutdownTimeoutCtx); err != nil {
			zap.S().Infoln("an error occurred during server shutdown: %v", err)
		}
	}()
}
