package services

import (
	"context"

	"github.com/heroiclabs/nakama-common/runtime"
	"github.com/titan/titan-runtime/modules/common/models"
)

func AccountUpdateId(ctx context.Context, nk runtime.NakamaModule, logger runtime.Logger, userID string, userName string, avatarURL string, langTag string, metadata map[string]interface{}, displayName string, timezone string, location string) error {
	return nk.AccountUpdateId(ctx, userID, userName, metadata, avatarURL, timezone, langTag, displayName, location)
}

func GetAccountId(ctx context.Context, nk runtime.NakamaModule, logger runtime.Logger, userID string) (*models.Account, error) {
	resp, err := nk.AccountGetId(ctx, userID)
	if err != nil {
		logger.Error("Error getting account id: %v", err)
		return nil, err
	}
	return &models.Account{
		UserID:      resp.User.Id,
		Username:    resp.User.Username,
		DisplayName: resp.User.DisplayName,
		AvatarURL:   resp.User.AvatarUrl,
		LangTag:     resp.User.LangTag,
		Location:    resp.User.Location,
		Timezone:    resp.User.Timezone,
	}, nil
}

func WalletUpdate(ctx context.Context, nk runtime.NakamaModule, logger runtime.Logger, userID string, changeSet map[string]int64, metadata map[string]interface{}, persistent bool) (models.Wallet, models.Wallet, error) {
	wallet1, wallet2, err := nk.WalletUpdate(ctx, userID, changeSet, metadata, persistent)
	if err != nil {
		logger.Error("Error updating wallet: %v", err)
		return models.Wallet{}, models.Wallet{}, err
	}
	return models.Wallet{
			Coins:    wallet1["coins"],
			Diamonds: wallet1["diamonds"],
		}, models.Wallet{
			Coins:    wallet2["coins"],
			Diamonds: wallet2["diamonds"],
		}, nil
}
