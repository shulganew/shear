package app

import (
	"context"

	"github.com/shulganew/shear.git/internal/service"
)

// Deleted short URL from common channel.
func DeleteShort(ctx context.Context, short *service.Shorten, delCh chan service.DelBatch) (done chan struct{}) {
	done = make(chan struct{})
	go func() {
		defer close(done)
		for {
			select {
			case <-ctx.Done():
				return
			case delBatch := <-delCh:
				short.DeleteBatch(ctx, delBatch)
			}
		}
	}()
	return
}
