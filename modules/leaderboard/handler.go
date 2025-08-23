package leaderboard

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/heroiclabs/nakama-common/api"
	"github.com/heroiclabs/nakama-common/runtime"
)

// no logic on node leaderboard reset for now, place holder function
func handleSeasonLeaderboardReset(
	ctx context.Context,
	logger runtime.Logger,
	db *sql.DB,
	nk runtime.NakamaModule,
	lb *api.Leaderboard,
	resetUnix int64,
) error {
	return nil
}

// daily leaderboardReset update the season leaderboard
// it takes the todays dailyleaderboard values and updates
// season leaderboard
func handleDailyLeaderboardReset(
	ctx context.Context,
	logger runtime.Logger,
	db *sql.DB,
	nk runtime.NakamaModule,
	lb *api.Leaderboard,
	reset int64,
) error {

	validationErr := validateDailyLeaderboardResetInputs(ctx, db, nk, lb)
	if validationErr != nil {
		logger.Error(validationErr.Error())
		return validationErr
	}

	dailyLbId := lb.GetId()
	metaData, _ := parseLeaderboardMetadata(lb.GetMetadata())

	associatedSeasonLbId, _ := metaData["season_leaderboard_id"]

	var cursor string
	totalProcessed := 0

	for {
		records, _, nextCursor, _, err := nk.LeaderboardRecordsList(
			ctx,
			lb.GetId(),
			nil,
			pageSize,
			cursor,
			reset,
		)
		if err != nil {
			return fmt.Errorf("error parsing through records (id=%s, cursor=%q): %w", dailyLbId, cursor, err)
		}

		for _, r := range records {
			if _, err := nk.LeaderboardRecordWrite(
				ctx,
				associatedSeasonLbId,
				r.GetOwnerId(),
				"",
				r.GetScore(),
				r.GetSubscore(),
				nil,
				nil,
			); err != nil {
				logger.Error(fmt.Sprintf("failed to write daily score to season score for user id %v: %s", r.GetOwnerId(), err))
			}

			totalProcessed++
		}

		if nextCursor == "" {
			break
		}
		cursor = nextCursor
	}

	logger.Info("season leaderboard id : %s got updated from daily leaderboard : %v reset", dailyLbId, associatedSeasonLbId)

	return nil
}

// no logic on node leaderboard reset for now, place holder function
func handleNodeLeaderboardReset(
	ctx context.Context,
	logger runtime.Logger,
	db *sql.DB,
	nk runtime.NakamaModule,
	lb *api.Leaderboard,
	resetUnix int64,
) error {
	return nil
}

func leaderBoardResetHandler(
	ctx context.Context,
	logger runtime.Logger,
	db *sql.DB,
	nk runtime.NakamaModule,
	lb *api.Leaderboard,
	reset int64,
) error {
	var meta = make(map[string]string)

	_ = json.Unmarshal([]byte(lb.GetMetadata()), &meta)
	lbType, ok := meta["leaderboard_type"]
	if !ok || lbType == "" {
		errMsg := "leaderboard_type key must be present in leaderboard metadata"
		logger.Error(errMsg)
		return errors.New(errMsg)
	}

	switch lbType {
	case "season":
		return handleSeasonLeaderboardReset(ctx, logger, db, nk, lb, reset)
	case "daily":
		return handleDailyLeaderboardReset(ctx, logger, db, nk, lb, reset)
	case "node":
		return handleNodeLeaderboardReset(ctx, logger, db, nk, lb, reset)
	default:
		errMsg := fmt.Sprintf("leaderboard_type %s doesn't exist or is not configured", lbType)
		logger.Error(errMsg)
	}

	return nil
}

// HandleUpdateLeaderBoardEvent routes incoming leaderboard update events to the appropriate handler
func HandleUpdateLeaderboardEvent(ctx context.Context, logger runtime.Logger, nk runtime.NakamaModule, evt *api.Event) {
	validationError := validateGenericLeaderboardEventInputs(evt)
	if validationError != nil {
		logger.Error(validationError.Error())
	}

	props := evt.GetProperties()
	leaderboardType, _ := props["leaderboard_type"]
	switch leaderboardType {
	case "node":
		processNodeLeaderboardEvent(ctx, logger, nk, evt)
	case "daily":
		processDailyLeaderboardEvent(ctx, logger, nk, evt)
	case "season":
		processSeasonLeaderboardEvent(ctx, logger, nk, evt)
	default:
		logger.Error("there is no default leaderboards, leaderboard_type must be passed")
	}
}
