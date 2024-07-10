package service

import (
	"context"
	"github.com/keweegen/tic-toe/internal/cache"
	"github.com/keweegen/tic-toe/internal/cache/redis"
	domain "github.com/keweegen/tic-toe/internal/domain/game"
	domainerrs "github.com/keweegen/tic-toe/internal/domain/game/errors"
	"log/slog"
	"time"
)

var _ domain.StreamService = (*GameStream)(nil)

type streamActionHandler func(ctx context.Context, msg domain.StreamMessage) error

type GameStream struct {
	gameTTL        time.Duration
	gameService    domain.Service
	actionHandlers map[domain.StreamAction]streamActionHandler

	cacheGamePlayer              cache.Cache[domain.Player]
	cacheGameCurrentMovePlayerID cache.Cache[string]
}

func NewGameStream(gameTTL time.Duration, cacheClient redis.CacheClient, gameService domain.Service) *GameStream {
	s := &GameStream{
		gameTTL:                      gameTTL,
		gameService:                  gameService,
		cacheGamePlayer:              redis.NewCache[domain.Player]("player", cacheClient),
		cacheGameCurrentMovePlayerID: redis.NewCache[string]("current_player_move", cacheClient),
	}

	s.actionHandlers = map[domain.StreamAction]streamActionHandler{
		domain.StreamActionPlayerJoin:       s.playerJoin,
		domain.StreamActionPlayerReady:      s.playerReady,
		domain.StreamActionPlayerMove:       s.playerMove,
		domain.StreamActionPlayerTimeout:    s.playerTimeout,
		domain.StreamActionPlayerDraw:       s.playerDraw,
		domain.StreamActionPlayerGaveUp:     s.playerGaveUp,
		domain.StreamActionPlayerDisconnect: s.playerDisconnect,
		domain.StreamActionGameStart:        s.gameStart,
		domain.StreamActionGameFinish:       s.gameFinish,
	}

	return s
}

func (s *GameStream) Handle(ctx context.Context, msg domain.StreamMessage) error {
	handler, ok := s.actionHandlers[msg.Action]
	if !ok || handler == nil {
		return nil
	}
	if err := handler(ctx, msg); err != nil {
		return err
	}

	return nil
}

func (s *GameStream) playerJoin(ctx context.Context, msg domain.StreamMessage) error {
	game, err := s.gameService.One(ctx, msg.GameID)
	if err != nil {
		return err
	}

	if len(game.Players) >= game.MaxPlayers {
		return domainerrs.ErrGameHasMaxPlayers
	}

	player := domain.Player{
		ID:        msg.UserID,
		GameID:    game.ID,
		StateCode: domain.PlayerStateCodeJoined,
		Position:  1, // TODO: set position
	}

	game.Players = append(game.Players, player)

	if err = s.cacheGamePlayer.Set(ctx, player.ID, player, s.gameTTL); err != nil {
		return err
	}

	s.gameService.Notify(ctx, game.ID, player.ID, domain.StreamActionPlayerJoin, player)

	return nil
}

func (s *GameStream) playerReady(ctx context.Context, msg domain.StreamMessage) error {
	player, err := s.cacheGamePlayer.One(ctx, msg.UserID)
	if err != nil {
		return err
	}

	if player.StateCode.Ready() {
		return nil
	}
	if !player.StateCode.Connected() {
		return domainerrs.ErrPlayerHasIncorrectState
	}

	player.StateCode = domain.PlayerStateCodeReady

	if err = s.cacheGamePlayer.Set(ctx, player.ID, player, s.gameTTL); err != nil {
		return err
	}

	s.gameService.Notify(ctx, msg.GameID, player.ID, domain.StreamActionPlayerJoin, player)

	return nil
}

func (s *GameStream) playerMove(ctx context.Context, msg domain.StreamMessage) error {
	input, ok := msg.Data.(domain.StreamDataPlayerMove)
	if !ok {
		return domainerrs.ErrInvalidStreamData
	}

	currentMovePlayerID, err := s.cacheGameCurrentMovePlayerID.One(ctx, msg.GameID)
	if err != nil {
		return err
	}
	if currentMovePlayerID != msg.UserID {
		return domainerrs.ErrNotYourTurn
	}

	moveParams := domain.MoveParams{
		StrikerID: input.StrikerID,
		UnitCode:  input.UnitCode,
	}

	move, err := s.gameService.Move(ctx, msg.GameID, moveParams)
	if err != nil {
		return err
	}

	if err = s.cacheGameCurrentMovePlayerID.Set(ctx, msg.GameID, move.BatterID, timeForMove); err != nil {
		return err
	}

	s.gameService.Notify(ctx, msg.GameID, msg.UserID, domain.StreamActionPlayerMove, move)

	return nil
}

func (s *GameStream) playerTimeout(ctx context.Context, msg domain.StreamMessage) error {
	slog.InfoContext(ctx, "player timeout", slog.String("game_id", msg.GameID), slog.String("user_id", msg.UserID))
	return nil
}

func (s *GameStream) playerDraw(ctx context.Context, msg domain.StreamMessage) error {
	slog.InfoContext(ctx, "player draw", slog.String("game_id", msg.GameID), slog.String("user_id", msg.UserID))
	return nil
}

func (s *GameStream) playerGaveUp(ctx context.Context, msg domain.StreamMessage) error {
	slog.InfoContext(ctx, "player gave up", slog.String("game_id", msg.GameID), slog.String("user_id", msg.UserID))
	return nil
}

func (s *GameStream) playerDisconnect(ctx context.Context, msg domain.StreamMessage) error {
	player, err := s.cacheGamePlayer.One(ctx, msg.UserID)
	if err != nil {
		return err
	}

	if player.StateCode.Disconnected() {
		return nil
	}

	player.StateCode = domain.PlayerStateCodeDisconnected

	if err = s.cacheGamePlayer.Set(ctx, player.ID, player, s.gameTTL); err != nil {
		return err
	}

	s.gameService.Notify(ctx, msg.GameID, msg.UserID, domain.StreamActionPlayerDisconnect, player)

	return nil
}

func (s *GameStream) playerLeave(ctx context.Context, msg domain.StreamMessage) error {
	game, err := s.gameService.One(ctx, msg.GameID)
	if err != nil {
		return err
	}

	for i, player := range game.Players {
		if player.ID == msg.UserID {
			game.Players = append(game.Players[:i], game.Players[i+1:]...)
			break
		}
	}

	if err = s.gameService.Update(ctx, game); err != nil {
		return err
	}
	if err = s.cacheGamePlayer.Delete(ctx, msg.UserID); err != nil {
		return err
	}

	return nil
}

func (s *GameStream) gameStart(ctx context.Context, msg domain.StreamMessage) error {
	game, err := s.gameService.Start(ctx, msg.GameID)
	if err != nil {
		return err
	}

	data := domain.StreamDataGameStart{
		StrikerID:   game.Players.First().ID,
		TimeForMove: timeForMove,
	}

	if err = s.cacheGameCurrentMovePlayerID.Set(ctx, game.ID, data.StrikerID, data.TimeForMove); err != nil {
		return err
	}

	// TODO: Add check time for move with offset +1 second

	s.gameService.Notify(ctx, game.ID, "", domain.StreamActionGameStart, data)

	return nil
}

func (s *GameStream) gameFinish(ctx context.Context, msg domain.StreamMessage) error {
	slog.InfoContext(ctx, "game finish", slog.String("game_id", msg.GameID), slog.String("user_id", msg.UserID))
	return nil
}
