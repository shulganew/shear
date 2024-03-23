package app

import (
	"context"
	"sync"

	"github.com/shulganew/shear.git/internal/config"
	"github.com/shulganew/shear.git/internal/service"
	"go.uber.org/zap"
)

// Graceful shutdown.
func Shutdown(ctx context.Context, wgroot *sync.WaitGroup, wgdel *sync.WaitGroup, conf *config.Config, short *service.Shorten, backup *service.Backup) {
	wgroot.Add(1)
	go func() {
		defer zap.S().Infoln("Graceful shutdown done.")
		defer wgroot.Done()
		<-ctx.Done()
		// Wait until all short will be async deleted.
		wgdel.Wait()
		if conf.IsBackup() {
			service.BackupShorts(short, *backup)
			return
		}

	}()
}

// Deleted short URL from common channel.
func DeleteShort(ctx context.Context, wgroot *sync.WaitGroup, short *service.Shorten, finalCh chan service.DelBatch) {
	wgroot.Add(1)
	go func() {
		defer zap.S().Infoln("Graceful shutdown done.")
		defer wgroot.Done()
		for {
			select {
			case <-ctx.Done():
				return
			case delBatch := <-finalCh:
				short.DeleteBatch(ctx, delBatch)
			}
		}
	}()
}
