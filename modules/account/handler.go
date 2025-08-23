package account

import (
	"context"
	"database/sql"
	"encoding/json"

	"github.com/heroiclabs/nakama-common/runtime"
	"github.com/titan/titan-runtime/modules/common/models"
)

func UpdateAccountHandler(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, payload string) (string, error) {
	userID, ok := ctx.Value(runtime.RUNTIME_CTX_USER_ID).(string)
	if !ok {
		return "", runtime.NewError("Authentication required", 16)
	}
	// parse request
	var req models.UpdateProfileRequest
	if err := json.Unmarshal([]byte(payload), &req); err != nil {
		return "", runtime.NewError("Invalid request", 14)
	}
	// Call service
	profile, err := UpdateAccount(ctx, nk, logger, userID, &req)
	if err != nil {
		logger.Error("Failed to get profile: %v", err)
		return "", runtime.NewError("Failed to get profile", 13)
	}

	// Return response
	responseJSON, _ := json.Marshal(profile)
	return string(responseJSON), nil
}
