package leaderboard

import (
	"context"
	"database/sql"

	"github.com/heroiclabs/nakama-common/runtime"
)

// ---- init ----
func InitModuleCallbacks(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, initializer runtime.Initializer) error {
	// ls1. load meta config
	if err := initializer.RegisterLeaderboardReset(leaderBoardResetHandler); err != nil {
		return err
	}

	logger.Info("Leaderboard CallBacks initialized")
	return nil
}
