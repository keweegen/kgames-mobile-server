package service

import (
	"context"
	"github.com/keweegen/tic-toe/internal/broadcaster"
	domain "github.com/keweegen/tic-toe/internal/domain/game"
	domainerrs "github.com/keweegen/tic-toe/internal/domain/game/errors"
	"github.com/shopspring/decimal"
	"github.com/volatiletech/null/v8"
	"time"
)

const (
	timeForMove = time.Second * 10
)

type Game struct {
	profitFeePercent decimal.Decimal
	repo             domain.Repository
	moveRepo         domain.MoveRepository
	mover            domain.Mover
	broadcaster      broadcaster.Broadcaster
}

var _ domain.Service = (*Game)(nil)

func NewGame(
	profitFeePercent decimal.Decimal,
	repo domain.Repository,
	moveRepo domain.MoveRepository,
	mover domain.Mover,
	broadcaster broadcaster.Broadcaster,
) Game {
	return Game{
		profitFeePercent: profitFeePercent,
		repo:             repo,
		moveRepo:         moveRepo,
		mover:            mover,
		broadcaster:      broadcaster,
	}
}

func (s Game) One(ctx context.Context, id string) (domain.Game, error) {
	return s.repo.One(ctx, id)
}

func (s Game) Create(ctx context.Context, params domain.CreateParams) (domain.Game, error) {
	players := domain.Players{
		{
			ID:       params.InitiatorID,
			Position: 0,
		},
	}

	game := domain.Game{
		TypeCode:   params.TypeCode,
		StateCode:  domain.StateCodeCreated,
		Bid:        params.Bid,
		MaxPlayers: params.MaxPlayers,
		Players:    players,
		CreatedBy:  params.InitiatorID,
	}

	return s.repo.Add(ctx, game)
}

func (s Game) Update(ctx context.Context, game domain.Game) error {
	return s.repo.Update(ctx, game)
}

func (s Game) Start(ctx context.Context, id string) (domain.Game, error) {
	game, err := s.repo.One(ctx, id)
	if err != nil {
		return domain.Game{}, err
	}
	if game.StateCode != domain.StateCodeCreated {
		return game, nil
	}
	if game.MaxPlayers != len(game.Players) {
		return domain.Game{}, domainerrs.ErrGameHasNotEnoughPlayers
	}

	game.StateCode = domain.StateCodeActive
	game.StartedAt = null.TimeFrom(time.Now())

	if err = s.Update(ctx, game); err != nil {
		return domain.Game{}, err
	}

	return game, nil
}

func (s Game) Move(ctx context.Context, id string, params domain.MoveParams) (domain.Move, error) {
	game, err := s.repo.One(ctx, id)
	if err != nil {
		return domain.Move{}, err
	}
	if game.Finished() {
		return domain.Move{}, domainerrs.ErrGameIsFinished
	}
	if game.StateCode != domain.StateCodeActive {
		return domain.Move{}, domainerrs.ErrGameIsNotActive
	}

	move := domain.Move{
		GameID:    game.ID,
		StrikerID: params.StrikerID,
		UnitCode:  params.UnitCode,
	}

	move, finished, finishParams, err := s.mover.Do(ctx, game, move)
	if err != nil {
		return domain.Move{}, err
	}

	if finished {
		if err = s.Finish(ctx, game.ID, finishParams); err != nil {
			return domain.Move{}, err
		}
	}

	return move, nil
}

func (s Game) GaveUp(ctx context.Context, id, playerID string) error {
	game, err := s.repo.One(ctx, id)
	if err != nil {
		return err
	}
	if game.Finished() {
		return domainerrs.ErrGameIsFinished
	}
	if !game.Players.Exists(playerID) {
		return domainerrs.ErrPlayerNotFound
	}

	finishParams := domain.FinishParams{
		ReasonCode: domain.FinishReasonCodeGaveUp,
		Players:    make([]domain.PlayerFinishParams, len(game.Players)),
	}

	fee := game.Bid.Mul(s.profitFeePercent)
	profit := game.Bid.Sub(fee)

	for i, player := range game.Players {
		if player.ID == playerID {
			finishParams.Players[i] = domain.PlayerFinishParams{
				PlayerID: player.ID,
				Profit:   game.Bid.Neg(),
				Fee:      decimal.Zero,
				GaveUp:   true,
			}
			continue
		}

		finishParams.Players[i] = domain.PlayerFinishParams{
			PlayerID: player.ID,
			Profit:   profit,
			Fee:      fee,
		}
	}

	if err = s.finish(ctx, game, finishParams, true); err != nil {
		return err
	}

	return nil
}

func (s Game) Finish(ctx context.Context, id string, params domain.FinishParams) error {
	game, err := s.repo.One(ctx, id)
	if err != nil {
		return err
	}
	return s.finish(ctx, game, params, true)
}

func (s Game) finish(ctx context.Context, game domain.Game, params domain.FinishParams, notify bool) error {
	if game.Finished() {
		return nil
	}

	mapOfPlayers := game.Players.Map()
	for _, playerParams := range params.Players {
		player, ok := mapOfPlayers[playerParams.PlayerID]
		if !ok {
			return domainerrs.ErrPlayerNotFound
		}

		player.Profit = playerParams.Profit
		player.Fee = playerParams.Fee
		player.Draw = playerParams.Draw
		player.GaveUp = playerParams.GaveUp

		mapOfPlayers[player.ID] = player
	}

	game.StateCode = domain.StateCodeFinished
	game.FinishReasonCode = params.ReasonCode
	game.Players = mapOfPlayers.Slice()
	game.FinishedAt = null.TimeFrom(time.Now())

	if err := s.repo.Update(ctx, game); err != nil {
		return err
	}

	if notify {
		s.Notify(ctx, game.ID, params.InitiatorID, domain.StreamActionGameFinish, params)
	}

	return nil
}

func (s Game) Notify(ctx context.Context, gameID, initiatorID string, action domain.StreamAction, data any) {
	msg := broadcaster.Message{
		GameID: gameID,
		UserID: initiatorID,
		Action: action,
		Data:   data,
	}

	excludeUserIDs := make([]string, 0)
	if initiatorID != "" {
		excludeUserIDs = append(excludeUserIDs, initiatorID)
	}

	s.broadcaster.Notify(ctx, gameID, msg, excludeUserIDs...)
}
