package leaderboard

import (
	"context"
	"strconv"
	"time"

	"github.com/heroiclabs/nakama-common/api"
	"github.com/heroiclabs/nakama-common/runtime"
	eventemitter "github.com/titan/titan-runtime/modules/common/eventEmitter"
)

func processNodeLeaderboardEvent(ctx context.Context, logger runtime.Logger, nk runtime.NakamaModule, evt *api.Event) {
	validationErr := validateNodeLeaderboardEventInputs(evt)
	if validationErr != nil {
		logger.Error(validationErr.Error())
		return
	}
	// once the validations are done, you can safely ignore the validation checks
	// as the validation logs any missing inputs and returns
	props := evt.GetProperties()
	nodeLbId := props["node_leaderboard_id"]
	scoreStr := props["score"]
	userId := props["user_id"]
	userName := props["user_name"]

	newScore, _ := parseScore(scoreStr)
	oldBest := getCurrentScore(ctx, logger, nk, nodeLbId, userId)
	if newScore <= oldBest {
		logger.Debug("No improvement in score; skipping update")
		return
	}

	if _, err := nk.LeaderboardRecordWrite(
		ctx,
		nodeLbId,
		userId,
		userName,
		newScore,
		0,
		map[string]interface{}{
			"source_event": evt.GetName(),
		},
		nil,
	); err != nil {
		logger.WithFields(
			map[string]interface{}{
				"properties": props,
			},
		).Error("Failed to write new record to node leaderboard " + err.Error())
		return
	}

	delta := newScore - oldBest
	logger.Info("Updated node leaderboard with new best score")

	// update existing properties for daily leaderboard
	props["leaderboard_type"] = "daily"
	props["delta"] = strconv.FormatInt(delta, 10)

	if err := eventemitter.EmitEvent(ctx, nk, "update_leaderboard", props); err != nil {
		logger.Error("Failed to emit daily leaderboard update event")
	}
	return
}

func processDailyLeaderboardEvent(ctx context.Context, logger runtime.Logger, nk runtime.NakamaModule, evt *api.Event) {
	validationErr := validateDailyLeaderboardEventInputs(evt)
	if validationErr != nil {
		logger.Error(validationErr.Error())
		return
	}

	props := evt.GetProperties()
	dailyLbId := props["daily_leaderboard_id"]
	deltaStr := props["delta"]
	userId := props["user_id"]
	userName := props["user_name"]

	delta, _ := parseScore(deltaStr)
	if delta == 0 {
		logger.Debug("Delta is zero; skipping daily leaderboard update")
		return
	}

	if _, err := nk.LeaderboardRecordWrite(
		ctx,
		dailyLbId,
		userId,
		userName, delta, 0, map[string]interface{}{
			"source_event": evt.GetName(),
		}, nil); err != nil {
		logger.Error("Failed to increment daily leaderboard score")
		return
	}

	logger.Info("Updated daily leaderboard with delta")
}

func processSeasonLeaderboardEvent(ctx context.Context, logger runtime.Logger, nk runtime.NakamaModule, evt *api.Event) {
	validationErr := validateSeasonLeaderboardEventInputs(evt)
	if validationErr != nil {
		logger.Error(validationErr.Error())
		return
	}

	props := evt.GetProperties()
	seasonLbId := props["season_leaderboard_id"]
	userId := props["user_id"]
	scoreStr := props["score"]
	userName := props["user_name"]

	newScore, _ := parseScore(scoreStr)
	oldBest := getCurrentScore(ctx, logger, nk, seasonLbId, userId)

	if newScore <= oldBest {
		logger.Debug("No improvement in score; skipping season leaderboard update")
		return
	}

	ctx2, cancel := context.WithTimeout(context.Background(), 500*time.Microsecond)
	defer cancel()
	if _, err := nk.LeaderboardRecordWrite(ctx2, seasonLbId, userId, userName, newScore, 0, map[string]interface{}{
		"source_event": evt.GetName(),
		"from_daily":   props["source_daily_id"],
	}, nil); err != nil {
		logger.Error("Failed to write new record to season leaderboard")
		return
	}

	logger.Info("Updated season leaderboard with new best score")
}
