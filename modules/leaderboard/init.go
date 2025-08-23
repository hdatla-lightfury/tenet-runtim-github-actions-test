package leaderboard

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/heroiclabs/nakama-common/runtime"
)

// ---- types for meta config ----
type LBConfig struct {
	Version     int       `json:"version"`
	Node        LBTypeCfg `json:"node"`
	Daily       LBTypeCfg `json:"daily"`
	Season      LBTypeCfg `json:"season"`
	Constraints struct {
		MinNodes int `json:"min_nodes"`
		MaxNodes int `json:"max_nodes"`
	} `json:"constraints"`
}
type LBTypeCfg struct {
	IDTemplate string                 `json:"id_template"`
	Sort       string                 `json:"sort"`
	Operator   string                 `json:"operator"`
	Reset      *string                `json:"reset"`
	Metadata   map[string]interface{} `json:"metadata"`
}

// ---- event data type ----
type Event struct {
	ID          string
	NodeCount   int
	SeasonEndTs int64
}

// ---- init ----
func InitModule(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, initializer runtime.Initializer) error {
	// 1. load meta config
	meta, err := loadLBConfig("modules/leaderboard/leaderboard_meta.json")
	logger.Info("meta_data processing happened")
	if err != nil {
		logger.Error(fmt.Sprintf("Error occured: %s", err))
		return err
	}

	// 2. get events from your event system (mocked here)
	events := []Event{
		{ID: "ipl_2025", NodeCount: 3, SeasonEndTs: 1748563200},
	}

	logger.Info("starting creation of leaderboards")
	// 3. create leaderboards per event
	for _, ev := range events {
		if err := createLeaderboardsForEvent(ctx, logger, nk, meta, ev); err != nil {
			logger.Error(fmt.Sprintf("failed to create leaderboards for event %s: %v", ev.ID, err))
		}
	}
	logger.Info("creation of leaderboards completed")

	return nil
}

// ---- helpers ----

func loadLBConfig(path string) (*LBConfig, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var cfg LBConfig
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

func createLeaderboardsForEvent(ctx context.Context, logger runtime.Logger, nk runtime.NakamaModule, meta *LBConfig, ev Event) error {
	// constraints
	if ev.NodeCount < meta.Constraints.MinNodes || ev.NodeCount > meta.Constraints.MaxNodes {
		return fmt.Errorf("node count out of allowed range: %d", ev.NodeCount)
	}

	// common template vars
	vars := map[string]string{
		"eventId":   ev.ID,
		"nodeIndex": "",
	}

	// --- 1) SEASON FIRST ---
	seasonID := expandTemplate(meta.Season.IDTemplate, vars)
	if err := nk.LeaderboardCreate(
		ctx,
		seasonID,
		true,
		meta.Season.Sort,
		meta.Season.Operator,
		stringOrEmpty(meta.Season.Reset),
		merge(meta.Season.Metadata, map[string]interface{}{
			"event_id":      ev.ID,
			"season_end_ts": ev.SeasonEndTs,
		}),
		true,
	); err != nil {
		logger.Debug(fmt.Sprintf("season leaderboard exists for event %s: %v", ev.ID, err))
	} else {
		logger.Info(fmt.Sprintf("created season leaderboard for event %s", ev.ID))
	}

	// --- 2) DAILY, with season_leaderboard_id in metadata ---
	dailyID := expandTemplate(meta.Daily.IDTemplate, vars)
	if err := nk.LeaderboardCreate(
		ctx,
		dailyID,
		true,
		meta.Daily.Sort,
		meta.Daily.Operator,
		stringOrEmpty(meta.Daily.Reset),
		merge(meta.Daily.Metadata, map[string]interface{}{
			"event_id":              ev.ID,
			"season_leaderboard_id": seasonID,
		}),
		true,
	); err != nil {
		logger.Debug(fmt.Sprintf("daily leaderboard exists for event %s: %v", ev.ID, err))
	} else {
		logger.Info(fmt.Sprintf("created daily leaderboard for event %s", ev.ID))
	}

	// --- 3) NODES, each with daily_leaderboard_id in metadata ---
	for i := 1; i <= ev.NodeCount; i++ {
		vars["nodeIndex"] = strconv.Itoa(i)
		nodeID := expandTemplate(meta.Node.IDTemplate, vars)

		if err := nk.LeaderboardCreate(
			ctx,
			nodeID,
			true,
			meta.Node.Sort,
			meta.Node.Operator,
			stringOrEmpty(meta.Node.Reset),
			merge(meta.Node.Metadata, map[string]interface{}{
				"event_id":             ev.ID,
				"node_index":           i,
				"daily_leaderboard_id": dailyID,
			}),
			true,
		); err != nil {
			logger.Debug(fmt.Sprintf("node leaderboard exists for event %s, node %d: %v", ev.ID, i, err))
		} else {
			logger.Info(fmt.Sprintf("created node leaderboard for event %s, node %d", ev.ID, i))
		}
	}

	return nil
}

func expandTemplate(tmpl string, vars map[string]string) string {
	out := tmpl
	for k, v := range vars {
		out = strings.ReplaceAll(out, "${"+k+"}", v)
	}
	return out
}

func stringOrEmpty(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

func merge(a, b map[string]interface{}) map[string]interface{} {
	out := make(map[string]interface{}, len(a)+len(b))
	for k, v := range a {
		out[k] = v
	}
	for k, v := range b {
		out[k] = v
	}
	return out
}
