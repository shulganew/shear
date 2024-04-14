package app

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/shulganew/shear.git/internal/config"
	pb "github.com/shulganew/shear.git/internal/handler/grpc/proto"
	"github.com/shulganew/shear.git/internal/handler/grpc/servgrpc/ghandlers"
	"github.com/shulganew/shear.git/internal/handler/grpc/servgrpc/interceptors"
	"github.com/shulganew/shear.git/internal/service"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

const timeoutServerShutdown = time.Second * 5

// Manage web server.
func StartREST(ctx context.Context, conf *config.Config, componentsErrs chan error, r *chi.Mux) (restDone chan struct{}) {
	// Start web server.
	var srv = http.Server{Addr: conf.GetAddrREST(), Handler: r}
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
	restDone = make(chan struct{})
	go func() {
		defer zap.S().Infoln("Server web has been graceful shutdown.")
		defer close(restDone)
		<-ctx.Done()
		shutdownTimeoutCtx, cancelShutdownTimeoutCtx := context.WithTimeout(context.Background(), timeoutServerShutdown)
		defer cancelShutdownTimeoutCtx()
		if err := srv.Shutdown(shutdownTimeoutCtx); err != nil {
			zap.S().Infoln("an error occurred during server shutdown: %v", err)
		}
	}()
	return
}

// Manage gRPC server.
func StartGRPC(ctx context.Context, serviceURL *service.Shorten, conf *config.Config, db *sql.DB, sd *service.Delete, componentsErrs chan error) (grpcDone chan struct{}) {
	// Add pass value to interceptors
	initCtx := func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		ctx = context.WithValue(ctx, config.CtxIP{}, conf.GetIP())
		ctx = context.WithValue(ctx, config.CtxPassKey{}, conf.GetPass())
		return handler(ctx, req)
	}

	s := grpc.NewServer(grpc.ChainUnaryInterceptor(initCtx, interceptors.AuthInterceptor))
	us := ghandlers.NewUsersServer(serviceURL, conf, db, sd)

	pb.RegisterUsersServer(s, us)
	go func() {
		// Start gRPC server.
		listen, err := net.Listen("tcp", conf.GetAddrGRPC())
		if err != nil {
			componentsErrs <- fmt.Errorf("listen gRPC failed: %w", err)
		}
		if err := s.Serve(listen); err != nil {
			componentsErrs <- fmt.Errorf("serve gRPC failed: %w", err)
		}
	}()

	// Graceful shutdown.
	grpcDone = make(chan struct{})
	go func() {
		defer zap.S().Infoln("Server gRPC has been graceful shutdown.")
		defer close(grpcDone)
		<-ctx.Done()
		s.GracefulStop()
	}()
	return
}
