package mover

import (
	"context"
	"errors"
	domain "github.com/keweegen/tic-toe/internal/domain/game"
	domainerrs "github.com/keweegen/tic-toe/internal/domain/game/errors"
	"github.com/shopspring/decimal"
)

type mover interface {
	do(moves []domain.Move, newMove *domain.Move, players domain.Players) (finished bool, err error)
	results(game domain.Game) (domain.FinishParams, error)
}

type Mover struct {
	repo   domain.MoveRepository
	movers map[domain.UnitCode]mover
}

var _ domain.Mover = (*Mover)(nil)

func New(repo domain.MoveRepository, profitFeePercent decimal.Decimal) *Mover {
	return &Mover{
		repo: repo,
		movers: map[domain.UnitCode]mover{
			domain.UnitCodeTicTacToeX: newTicTacToeMover(domain.UnitCodeTicTacToeX, profitFeePercent),
			domain.UnitCodeTicTacToeO: newTicTacToeMover(domain.UnitCodeTicTacToeO, profitFeePercent),
		},
	}
}

func (m *Mover) Do(
	ctx context.Context,
	game domain.Game,
	newMove domain.Move,
) (move domain.Move, finished bool, params domain.FinishParams, err error) {
	unitMover, ok := m.movers[newMove.UnitCode]
	if !ok {
		return domain.Move{}, false, domain.FinishParams{}, domainerrs.ErrMoveNotFound
	}

	finished, err = unitMover.do(game.Moves, &newMove, game.Players)
	if err != nil {
		if errors.Is(err, domainerrs.ErrGameInternal) {
			return domain.Move{}, true, m.resultsForInternalError(game), nil
		}

		return domain.Move{}, false, domain.FinishParams{}, err
	}

	if finished {
		params, err = unitMover.results(game)
		if err != nil {
			return domain.Move{}, false, domain.FinishParams{}, err
		}
	}

	move, err = m.repo.Add(ctx, newMove)
	if err != nil {
		return domain.Move{}, false, domain.FinishParams{}, err
	}

	return move, finished, params, nil
}

func (m *Mover) resultsForInternalError(game domain.Game) domain.FinishParams {
	players := make([]domain.PlayerFinishParams, len(game.Players))
	for i, player := range game.Players {
		players[i] = domain.PlayerFinishParams{
			PlayerID: player.ID,
			Profit:   game.Bid,
			Fee:      decimal.Zero,
		}
	}

	return domain.FinishParams{
		ReasonCode: domain.FinishReasonCodeInvalid,
		Players:    players,
	}
}
