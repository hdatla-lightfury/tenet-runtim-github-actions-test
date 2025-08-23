package leaderboard

import (
	"context"
	"encoding/json"
	"errors"
	"strconv"
	"strings"

	"github.com/heroiclabs/nakama-common/api"
	"github.com/heroiclabs/nakama-common/runtime"
)

func parseScore(s string) (int64, error) {
	return strconv.ParseInt(s, 10, 64)
}

func checkLeaderboardExists(
	ctx context.Context,
	nk runtime.NakamaModule,
	lb *api.Leaderboard,
) (bool, error) {
	_, _, _, _, err := nk.LeaderboardRecordsList(
		ctx,
		lb.GetId(),
		nil,
		1,
		"",
		0,
	)
	if err != nil {
		// Distinguish "not found" from other errors.
		// (There isn't a dedicated helper for this; string check is the common pattern.)
		// https://github.com/heroiclabs/nakama/blob/bab2d8f6f49d98ce5fcb774267b4cc0cf894c175/server/core_leaderboard.go#L41C2-L41C24
		if strings.Contains(err.Error(), "leaderboard not found") {
			return false, errors.New("leaderboard not found")
		}
		return false, err
	}

	return true, nil
}

func getCurrentScore(
	ctx context.Context,
	logger runtime.Logger,
	nk runtime.NakamaModule,
	leaderboardId string,
	userId string,
) int64 {
	_, userRecords, _, _, err := nk.LeaderboardRecordsList(
		ctx,
		leaderboardId,
		[]string{userId},
		1,
		"",
		0,
	)

	if err != nil {
		logger.Error("Failed to fetch current leaderboard record")
		return 0
	}
	if userRecords == nil || len(userRecords) == 0 {
		return 0
	}
	return userRecords[0].Score
}

func parseLeaderboardMetadata(meta string) (map[string]string, error) {
	if strings.TrimSpace(meta) == "" {
		return nil, nil
	}
	var m map[string]string
	if err := json.Unmarshal([]byte(meta), &m); err != nil {
		return nil, err
	}
	return m, nil
}
