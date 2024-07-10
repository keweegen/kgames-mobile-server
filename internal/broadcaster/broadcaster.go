package broadcaster

import (
	"context"
	gamepb "github.com/keweegen/tic-toe/api/grpc/game"
	gamedomain "github.com/keweegen/tic-toe/internal/domain/game"
)

type Broadcaster interface {
	Subscribe(ctx context.Context, gameID, playerID string, stream gamepb.Service_GameStreamingServer)
	Unsubscribe(gameID, playerID string)
	UnsubscribeAll()
	Channel() <-chan Message
	Notify(ctx context.Context, gameID string, message Message, excludePlayerIDs ...string)
	Shutdown() error
}

type Message struct {
	Error       error
	Description string
	GameID      string
	UserID      string
	Action      gamedomain.StreamAction
	Data        any
}
