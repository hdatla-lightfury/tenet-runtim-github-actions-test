package account

import (
	"context"

	"github.com/heroiclabs/nakama-common/api"
	"github.com/heroiclabs/nakama-common/runtime"
	eventEmitter "github.com/titan/titan-runtime/modules/common/eventEmitter"
	"github.com/titan/titan-runtime/modules/common/models"
	"github.com/titan/titan-runtime/modules/common/notifier"
)

func HandleAccountUpdatedEvent(ctx context.Context, logger runtime.Logger, evt *api.Event) {
	logger.Debug("account_updated event received")

}

func UpdateAccount(ctx context.Context, nk runtime.NakamaModule, logger runtime.Logger, userID string, req *models.UpdateProfileRequest) (*models.UpdateProfileResponse, error) {
	logger.Info("Updating account: %+v", req)

	if err := eventEmitter.EmitEvent(ctx, nk, "account_updated", map[string]string{"profile": req.DisplayName}); err != nil {
		logger.Error("Failed to emit account_updated event: %v", err)
	}
	content := map[string]interface{}{
		"display_name": req.DisplayName,
	}
	if err := notifier.SendNotifications(ctx, nk, logger, userID, userID, content, false, int(UserProfileUpdated)); err != nil {
		logger.Error("Failed to send notifications: %v", err)
	}
	return &models.UpdateProfileResponse{
		Success: true,
		Message: "Profile updated successfully",
	}, nil
}
