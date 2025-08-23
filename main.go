package main

import (
	"context"
	"database/sql"
	"time"

	"github.com/heroiclabs/nakama-common/runtime"
	"github.com/titan/titan-runtime/modules/account"
	eventProcessor "github.com/titan/titan-runtime/modules/common/eventProcessor"
	"github.com/titan/titan-runtime/modules/leaderboard"
	"github.com/titan/titan-runtime/modules/test_events"
)

func InitModule(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, initializer runtime.Initializer) error {
	initStart := time.Now()
	logger.Info("Initializing Titan Runtime")
	account.InitModule(ctx, logger, db, nk, initializer)
	eventProcessor.InitModule(ctx, logger, db, nk, initializer)
	leaderboard.InitModuleCallbacks(ctx, logger, db, nk, initializer)
	test_events.InitModule(ctx, logger, db, nk, initializer)
	logger.Info("Titan Runtime initialized in %s", time.Since(initStart))
	return nil
}
