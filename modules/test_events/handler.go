package test_events

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/heroiclabs/nakama-common/api"
	"github.com/heroiclabs/nakama-common/runtime"
	eventEmitter "github.com/titan/titan-runtime/modules/common/eventEmitter"
)

func handleEmitEvent(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, payload string) (string, error) {
	// Unmarshal payload into an api.Event struct
	evt := &api.Event{}
	if err := json.Unmarshal([]byte(payload), evt); err != nil {
		logger.Error("Failed to parse event payload: %v", err)
		return "", runtime.NewError("invalid payload", 3)
	}

	eventName := evt.Name
	props := evt.Properties

	// Call your custom emitter with the event name and properties from the payload.
	if err := eventEmitter.EmitEvent(ctx, nk, eventName, props); err != nil {
		logger.Error("Failed to emit %s event: %v", eventName, err)
		return "", runtime.NewError("event emit failed", 13)
	}
	time.Sleep(50000 * time.Millisecond)

	logger.Info("Processed incoming event: %s", evt.Name)
	return fmt.Sprintf(`{"status":"ok","event":"%s"}`, evt.Name), nil
}
