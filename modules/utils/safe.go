package utils

import (
	"context"
	"time"

	"github.com/heroiclabs/nakama-common/runtime"
)

func SpawnSafe(parent context.Context, logger runtime.Logger, fn func(ctx context.Context)) {
	ctx, cancel := context.WithTimeout(parent, 5*time.Second)
	go func() {
		defer cancel()
		defer func() {
			if r := recover(); r != nil {
				logger.Warn("Recovered in goroutine: %v", r)
			}
		}()
		fn(ctx)
	}()
}
