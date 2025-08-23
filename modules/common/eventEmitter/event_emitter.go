package eventemitter

import (
	"context"

	"github.com/heroiclabs/nakama-common/api"
	"github.com/heroiclabs/nakama-common/runtime"
)

func EmitEvent(ctx context.Context, nk runtime.NakamaModule, eventName string, properties map[string]string) error {
	evt := &api.Event{
		Name:       eventName,
		Properties: properties,
		External:   false,
	}
	return nk.Event(ctx, evt)
}
