package leaderboard

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/heroiclabs/nakama-common/api"
	"github.com/heroiclabs/nakama-common/runtime"
)

// TO DO : Integrate a good rule base validator, refer the below article
// For data structures that are being imported from external libraries like here
// redefine the data structure with struct validate tags
// https://dev.to/kittipat1413/a-guide-to-input-validation-in-go-with-validator-v10-56bp

func validateGenericLeaderboardEventInputs(evt *api.Event) error {
	var errs []string

	if evt.Name == "" {
		errs = append(errs, "id is required")
	}
	if evt.Properties == nil || len(evt.Properties) == 0 {
		errs = append(errs, "Properties map cannot be empty")
	}

	props := evt.GetProperties()
	leaderboardType, ok := props["leaderboard_type"]
	if !ok {
		errs = append(errs, "leaderboard_type key must be present")
	}
	if ok {
		switch leaderboardType {
		case "node", "daily", "season":
			// valid
		default:
			errs = append(errs, fmt.Sprintf("invalid leaderboard_type: %q", leaderboardType))
		}
	}

	nodeLbId := props["node_leaderboard_id"]
	scoreStr := props["score"]
	userId := props["user_id"]
	userName := props["user_name"]
	dailyLbId := props["daily_leaderboard_id"]
	seasonLbId := props["season_leaderboard_id"]

	if nodeLbId == "" {
		errs = append(errs, "missing node_leaderboard_id")
	}
	if scoreStr == "" {
		errs = append(errs, "missing score")
	}
	if userId == "" {
		errs = append(errs, "missing user_id")
	}
	if userName == "" {
		errs = append(errs, "missing user_name")
	}
	if dailyLbId == "" {
		errs = append(errs, "missing daily_leaderboard_id")
	}
	if seasonLbId == "" {
		errs = append(errs, "missing season_leaderboard_id")
	}

	if scoreStr != "" {
		if _, err := parseScore(scoreStr); err != nil {
			errs = append(errs, fmt.Sprintf("invalid score value %q: %v", scoreStr, err))
		}
	}

	if len(errs) > 0 {
		return errors.New(strings.Join(errs, "; "))
	}
	return nil
}

func validateNodeLeaderboardEventInputs(evt *api.Event) error {
	// all the validation are done in GenericLeaderboardEventInputs
	// any custom ones to be added here in future
	return nil
}

func validateDailyLeaderboardEventInputs(evt *api.Event) error {
	var errs []string

	props := evt.GetProperties()
	deltaStr := props["delta"]

	_, err := parseScore(deltaStr)
	if err != nil {
		errs = append(errs, "Invalid delta value for daily leaderboard event")
	}

	if len(errs) > 0 {
		return errors.New(strings.Join(errs, "; "))
	}
	return nil
}

func validateSeasonLeaderboardEventInputs(evt *api.Event) error {
	// all the validation are done in GenericLeaderboardEventInputs
	// any custom ones to be added here in future
	return nil
}

func validateDailyLeaderboardResetInputs(
	ctx context.Context,
	db *sql.DB,
	nk runtime.NakamaModule,
	lb *api.Leaderboard,
) error {

	dailyLbId := lb.GetId()
	metaData, err := parseLeaderboardMetadata(lb.GetMetadata())
	if err != nil {
		errMsg := fmt.Sprintf("failed to parse leaderboard metadata for id : %s", dailyLbId)
		return errors.New(errMsg)
	}

	associatedSeasonLbId, ok := metaData["season_leaderboard_id"]
	if !ok || associatedSeasonLbId == "" {
		errMsg := "an associated season leaderboard id for the daily leaderboard id" +
			"must be present in leaderboard metadata"
		return errors.New(errMsg)
	}

	lbExists, err := checkLeaderboardExists(ctx, nk, lb)
	if !lbExists || err != nil {
		errMsg := fmt.Sprintf("leaderboard with id %s doesn't exist", dailyLbId) + err.Error()
		return errors.New(errMsg)
	}
	return nil
}
