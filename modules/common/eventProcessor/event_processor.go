package eventprocessor

import (
	"context"
	"database/sql"

	"github.com/heroiclabs/nakama-common/api"
	"github.com/heroiclabs/nakama-common/runtime"
	"github.com/titan/titan-runtime/modules/account"
	"github.com/titan/titan-runtime/modules/leaderboard"
)

func InitModule(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, initializer runtime.Initializer) error {
	logger.Info("Initializing EventProcessor domain...")
	if err := initializer.RegisterEvent(ProcessEvent(nk, db)); err != nil {
		return err
	}
	logger.Info("EventProcessor domain initialized")
	return nil
}

func ProcessEvent(nk runtime.NakamaModule, db *sql.DB) func(ctx context.Context, logger runtime.Logger, evt *api.Event) {
	return func(ctx context.Context, logger runtime.Logger, evt *api.Event) {
		switch evt.GetName() {
		case "account_updated":
			logger.Debug("[WORKER]account_updated event received")
			account.HandleAccountUpdatedEvent(ctx, logger, evt)
		case "profile_updated":
			logger.Debug("[WORKER]profile_updated event received")
			// profile.HandleProfileUpdatedEvent(ctx, logger, evt)
		case "update_leaderboard":
			logger.Debug("[WORKER]update_leaderboard event received")
			leaderboard.HandleUpdateLeaderboardEvent(ctx, logger, nk, evt)
		default:
			logger.Error("unrecognized event: %+v", evt)
		}
	}
}
