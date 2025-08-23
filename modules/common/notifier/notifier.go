package notifier

import (
	"context"

	"github.com/heroiclabs/nakama-common/runtime"
)

func SendNotifications(ctx context.Context, nk runtime.NakamaModule, logger runtime.Logger, user1, user2 string, content map[string]interface{}, persistent bool, code int) error {
	notifications := []*runtime.NotificationSend{
		{
			UserID:     user2,
			Subject:    "New Friend",
			Content:    content,
			Code:       code,
			Persistent: persistent,
		},
		{
			UserID:     user1,
			Subject:    "Friend Added",
			Content:    content,
			Code:       code,
			Persistent: persistent,
		},
	}

	if err := nk.NotificationsSend(ctx, notifications); err != nil {
		logger.Warn("Failed to send friend notifications: %v", err)
		return err
	}
	return nil
}
