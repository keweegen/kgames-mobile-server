package grpc

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v8"
	gamepb "github.com/keweegen/tic-toe/api/grpc/game"
	"github.com/keweegen/tic-toe/internal/broadcaster"
	"github.com/keweegen/tic-toe/internal/broadcaster/grpc/convert"
	"github.com/keweegen/tic-toe/internal/logger"
	"log/slog"
	"sync"
)

var _ broadcaster.Broadcaster = (*Broadcaster)(nil)

type Broadcaster struct {
	keyPrefix   string
	clients     map[string]*client
	message     chan broadcaster.Message
	redisClient *redis.Client
	mu          sync.RWMutex
}

type client struct {
	ctx     context.Context
	cancel  context.CancelFunc
	pubSub  *redis.PubSub
	players map[string]gamepb.Service_GameStreamingServer
}

func NewBroadcaster(keyPrefix string, redisClient *redis.Client) *Broadcaster {
	return &Broadcaster{
		keyPrefix:   keyPrefix,
		clients:     make(map[string]*client),
		message:     make(chan broadcaster.Message),
		redisClient: redisClient,
	}
}

func (b *Broadcaster) key(gameID string) string {
	return fmt.Sprintf("%s:%s", b.keyPrefix, gameID)
}

func (b *Broadcaster) Subscribe(
	ctx context.Context,
	gameID, playerID string,
	stream gamepb.Service_GameStreamingServer,
) {
	b.mu.Lock()
	defer b.mu.Unlock()

	if _, ok := b.clients[gameID]; !ok {
		clientCtx, cancel := context.WithCancel(ctx)

		c := &client{
			ctx:     clientCtx,
			cancel:  cancel,
			pubSub:  b.redisClient.Subscribe(ctx, b.key(gameID)),
			players: make(map[string]gamepb.Service_GameStreamingServer),
		}

		b.clients[gameID] = c
		go func() { b.read(gameID, c) }()
	}

	b.clients[gameID].players[playerID] = stream

	slog.DebugContext(ctx, "subscribe to game",
		slog.String(logger.FieldGameID, gameID),
		slog.String(logger.FieldUserID, playerID))

	fmt.Printf("CLIENTS: %+v\n", b.clients)
}

func (b *Broadcaster) read(gameID string, c *client) {
	ch := c.pubSub.Channel()

	for {
		select {
		case <-c.ctx.Done():
			return
		case payload := <-ch:
			var msg broadcaster.Message
			if err := json.Unmarshal([]byte(payload.Payload), &msg); err != nil {
				slog.ErrorContext(c.ctx, "unmarshal notification payload",
					slog.String(logger.FieldGameID, gameID),
					slog.String(logger.FieldNotificationPayload, payload.Payload),
					slog.Any(logger.FieldError, err),
				)
				continue
			}

			fmt.Println("Receive message from redis", msg)

			b.sendAll(c.ctx, msg)
		}
	}
}

func (b *Broadcaster) sendAll(ctx context.Context, msg broadcaster.Message) {
	excludeUserIDs := make([]string, 0)
	//if msg.UserID != "" {
	//	excludeUserIDs = append(excludeUserIDs, msg.UserID)
	//}

	mapExcludePlayer := make(map[string]struct{}, len(excludeUserIDs))
	for _, playerID := range excludeUserIDs {
		mapExcludePlayer[playerID] = struct{}{}
	}

	b.mu.RLock()
	defer b.mu.RUnlock()

	fmt.Println("Send message to all players", msg.GameID, len(b.clients[msg.GameID].players), b.clients[msg.GameID].players)

	for playerID, stream := range b.clients[msg.GameID].players {
		fmt.Println()
		fmt.Println("Range clients", playerID)

		if _, ok := mapExcludePlayer[playerID]; ok {
			continue
		}

		if err := stream.Send(convert.Broadcaster.Response(msg)); err != nil {
			slog.ErrorContext(ctx, "send message to player",
				slog.String(logger.FieldGameID, msg.GameID),
				slog.String(logger.FieldUserID, playerID),
				slog.String(logger.FieldNotificationAction, msg.Action.String()),
				slog.Any(logger.FieldNotificationData, msg.Data),
				slog.Any(logger.FieldError, err),
			)
			continue
		}

		fmt.Println("Send message to player", playerID, msg)
	}
}

func (b *Broadcaster) Unsubscribe(gameID, playerID string) {
	b.mu.Lock()
	defer b.mu.Unlock()

	if _, ok := b.clients[gameID]; ok {
		delete(b.clients[gameID].players, playerID)
	}
}

func (b *Broadcaster) UnsubscribeByGame(gameID string) {
	b.mu.Lock()
	defer b.mu.Unlock()

	if gameClients, ok := b.clients[gameID]; ok {
		if err := gameClients.pubSub.Close(); err != nil {
			slog.ErrorContext(gameClients.ctx, "close redis pub-sub",
				slog.String(logger.FieldGameID, gameID),
			)
		}

		gameClients.cancel()
		delete(b.clients, gameID)
	}
}

func (b *Broadcaster) UnsubscribeAll() {
	for gameID := range b.clients {
		b.UnsubscribeByGame(gameID)
	}
}

func (b *Broadcaster) Channel() <-chan broadcaster.Message {
	return b.message
}

func (b *Broadcaster) Notify(ctx context.Context, gameID string, msg broadcaster.Message, _ ...string) {
	fmt.Println("Receive task for notify", msg)

	binaryMessage, err := json.Marshal(msg)
	if err != nil {
		slog.ErrorContext(ctx, "marshal notification message",
			slog.String(logger.FieldGameID, gameID),
			slog.String(logger.FieldNotificationAction, msg.Action.String()),
			slog.Any(logger.FieldNotificationData, msg.Data),
			slog.Any(logger.FieldError, err),
		)
		return
	}

	fmt.Println("publish message in redis", msg)
	if err := b.redisClient.Publish(ctx, b.key(gameID), binaryMessage).Err(); err != nil {
		slog.ErrorContext(ctx, "publish message in redis",
			slog.String(logger.FieldGameID, gameID),
			slog.String(logger.FieldNotificationAction, msg.Action.String()),
			slog.Any(logger.FieldNotificationData, msg.Data),
			slog.Any(logger.FieldError, err),
		)
	}
}

func (b *Broadcaster) Shutdown() error {
	b.UnsubscribeAll()
	close(b.message)
	return nil
}
