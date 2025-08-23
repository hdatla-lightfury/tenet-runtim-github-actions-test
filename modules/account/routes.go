package account

import (
	"context"
	"database/sql"

	"github.com/heroiclabs/nakama-common/runtime"
)

// ONE InitModule per domain - handles ALL user stuff
func InitModule(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, initializer runtime.Initializer) error {
	logger.Info("Initializing User domain...")
	if err := initializer.RegisterRpc("update_account", UpdateAccountHandler); err != nil {
		return err
	}
	if err := initializer.RegisterBeforeAuthenticateDevice(BeforeAuthenticateDevice); err != nil {
		return err
	}
	if err := initializer.RegisterAfterAuthenticateDevice(AfterAuthenticateDevice); err != nil {
		return err
	}

	if err := initializer.RegisterAfterUpdateAccount(AfterUpdateAccount); err != nil {
		return err
	}

	logger.Info("User domain initialized")
	return nil
}
