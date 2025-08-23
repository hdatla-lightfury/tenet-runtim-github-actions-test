package test_events

import (
	"context"
	"database/sql"

	"github.com/heroiclabs/nakama-common/runtime"
)

// ONE InitModule per domain - handles ALL user stuff
func InitModule(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, initializer runtime.Initializer) error {
	logger.Info("Initializing emit event domain...")
	if err := initializer.RegisterRpc("test_emit_event", handleEmitEvent); err != nil {
		return err
	}

	logger.Info("Emit event initialized")
	return nil
}
