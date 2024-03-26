package shutdown

import (
	"context"
	"time"

	"go.uber.org/zap"
)

const timeoutShutdown = time.Second * 20

// Block gorutine until context done or get errer from componentsErrs channel.
func Graceful(ctx context.Context, cancel context.CancelFunc, componentsErrs chan error) {
	// Timer hardreset shutdown.
	context.AfterFunc(ctx, func() {
		ctx, cancelCtx := context.WithTimeout(context.Background(), timeoutShutdown)
		defer cancelCtx()
		<-ctx.Done()
		zap.S().Fatalln("failed to gracefully shutdown the service")
	})

	// Graceful shutdown.
	select {
	// Exit on root context done.
	case <-ctx.Done():
	// Exit on errors.
	case err := <-componentsErrs:
		zap.S().Errorln("Get server error: ", err)
		cancel()
	}
}
