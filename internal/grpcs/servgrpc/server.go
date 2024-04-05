package servgrpc

import (
	"context"
	"database/sql"
	"fmt"
	"net"

	"github.com/shulganew/shear.git/internal/config"
	"github.com/shulganew/shear.git/internal/service"
	"go.uber.org/zap"
	"google.golang.org/grpc"

	pb "github.com/shulganew/shear.git/internal/grpcs/proto"
	"github.com/shulganew/shear.git/internal/grpcs/servgrpc/ghandlers"
	"github.com/shulganew/shear.git/internal/grpcs/servgrpc/interceptors"
)

// TODO GRPC tls
// https://github.com/grpc/grpc-go/blob/master/examples/features/encryption/TLS/server/main.go
// Manage gRPC server.
func Shortener(ctx context.Context, serviceURL *service.Shorten, conf *config.Config, db *sql.DB, sd *service.Delete, componentsErrs chan error) (rpcDone chan struct{}) {
	// Add pass value to interceptors
	initCtx := func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		newCtx := context.WithValue(ctx, config.CtxPassKey{}, conf.GetPass())
		return handler(newCtx, req)
	}

	s := grpc.NewServer(grpc.ChainUnaryInterceptor(initCtx, interceptors.LogInterceptor))
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
	rpcDone = make(chan struct{})
	go func() {
		defer zap.S().Infoln("Server gRPC has been graceful shutdown.")
		defer close(rpcDone)
		<-ctx.Done()
		s.GracefulStop()
	}()
	return
}
