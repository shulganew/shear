package service

import (
	"context"
)

// Deleted short URL from common channel.
func DeleteShort(ctx context.Context, short *Shorten, delCh chan DelBatch) (done chan struct{}) {
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
