package grpc

import (
	"context"
	"errors"
	"fmt"
	gamepb "github.com/keweegen/tic-toe/api/grpc/game"
	"github.com/keweegen/tic-toe/internal/broadcaster"
	domain "github.com/keweegen/tic-toe/internal/domain/game"
	"github.com/keweegen/tic-toe/internal/domain/game/api/grpc/convert"
	"io"
)

var _ gamepb.ServiceServer = (*GameHandler)(nil)

type GameHandler struct {
	service       domain.Service
	streamService domain.StreamService
	broadcaster   broadcaster.Broadcaster
}

func NewGameHandler(
	service domain.Service,
	streamService domain.StreamService,
	bc broadcaster.Broadcaster,
) *GameHandler {
	return &GameHandler{
		service:       service,
		streamService: streamService,
		broadcaster:   bc,
	}
}

func (h *GameHandler) CreateGame(ctx context.Context, req *gamepb.CreateGameRequest) (*gamepb.GameResponse, error) {
	g, err := h.service.Create(ctx, convert.Request.CreateGame(req))
	if err != nil {
		return nil, err
	}
	return convert.Response.Game(g), nil
}

func (h *GameHandler) GameStreaming(stream gamepb.Service_GameStreamingServer) error {
	ctx := stream.Context()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			req, err := stream.Recv()
			if errors.Is(err, io.EOF) {
				return nil
			}
			if err != nil {
				return err
			}

			fmt.Println("Receive message from gRPC", req.GameId, req.UserId)

			h.broadcaster.Subscribe(
				ctx,
				req.GameId,
				req.UserId,
				stream,
			)

			if err := h.streamService.Handle(stream.Context(), convert.Request.Stream(req)); err != nil {
				if sendErr := stream.SendMsg(convert.Response.Error(err)); err != nil {
					return sendErr
				}
			}
		}
	}
}
