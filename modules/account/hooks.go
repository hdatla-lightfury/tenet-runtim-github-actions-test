package account

import (
	"context"
	"database/sql"

	"github.com/heroiclabs/nakama-common/api"
	"github.com/heroiclabs/nakama-common/runtime"
	eventEmitter "github.com/titan/titan-runtime/modules/common/eventEmitter"
	"github.com/titan/titan-runtime/modules/common/notifier"
)

func BeforeAuthenticateDevice(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, in *api.AuthenticateDeviceRequest) (*api.AuthenticateDeviceRequest, error) {
	logger.Info("BeforeAuthenticateDevice:----------------- %+v", in.Username)
	return in, nil
}

func AfterAuthenticateDevice(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, out *api.Session, in *api.AuthenticateDeviceRequest) error {
	content := map[string]interface{}{
		"display_name": in.Username,
		"message":      "logged in successfully",
	}
	if err := notifier.SendNotifications(ctx, nk, logger, in.Account.Id, in.Account.Id, content, false, int(UserProfileUpdated)); err != nil {
		logger.Error("Failed to send notifications: %v", err)
	}
	if err := eventEmitter.EmitEvent(ctx, nk, "account_updated", map[string]string{"profile": in.Username}); err != nil {
		logger.Error("Failed to emit account_logged_in event: %v", err)
	}
	logger.Info("AfterAuthenticateDevice:----------------- %+v", out.Token)
	return nil
}

func AfterUpdateAccount(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, in *api.UpdateAccountRequest) error {
	logger.Info("AfterUpdateAccount:----------------- %+v", in.Username)
	if err := eventEmitter.EmitEvent(ctx, nk, "account_updated", map[string]string{"profile": in.DisplayName.Value}); err != nil {
		logger.Error("Failed to emit account_updated event: %v", err)
	}
	content := map[string]interface{}{
		"display_name": in.Username,
		"message":      "profile updated successfully",
	}
	if err := notifier.SendNotifications(ctx, nk, logger, in.Username.Value, in.Username.Value, content, false, int(UserProfileUpdated)); err != nil {
		logger.Error("Failed to send notifications: %v", err)
	}
	return nil
}
