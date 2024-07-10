package convert

import (
	"github.com/keweegen/tic-toe/internal/db/model"
	domain "github.com/keweegen/tic-toe/internal/domain/game"
	"github.com/volatiletech/null/v8"
)

var Game game

type game struct{}

func (game) Domain(m *model.Game) (domain.Game, error) {
	var moves []domain.Move
	if m.R != nil && len(m.R.GameMoves) > 0 {
		var err error
		moves, err = Move.DomainSlice(m.R.GameMoves)
		if err != nil {
			return domain.Game{}, err
		}
	}

	var players []domain.Player
	if m.R != nil && len(m.R.GamePlayers) > 0 {
		players = Player.DomainSlice(m.R.GamePlayers)
	}

	return domain.Game{
		ID:               m.ID,
		TypeCode:         domain.TypeCode(m.TypeCode),
		StateCode:        domain.StateCode(m.StateCode),
		Bid:              m.Bid,
		FinishReasonCode: domain.FinishReasonCode(m.FinishReasonCode.String),
		Players:          players,
		MaxPlayers:       int(m.MaxPlayers),
		Moves:            moves,
		CreatedAt:        m.CreatedAt,
		CreatedBy:        m.CreatedBy,
		StartedAt:        m.StartedAt,
		FinishedAt:       m.FinishedAt,
	}, nil
}

func (game) Model(d domain.Game) *model.Game {
	return &model.Game{
		ID:               d.ID,
		TypeCode:         d.TypeCode.String(),
		StateCode:        d.StateCode.String(),
		Bid:              d.Bid,
		MaxPlayers:       int16(d.MaxPlayers),
		FinishReasonCode: null.NewString(d.FinishReasonCode.String(), !d.FinishReasonCode.Unknown()),
		CreatedAt:        d.CreatedAt,
		CreatedBy:        d.CreatedBy,
		StartedAt:        d.StartedAt,
		FinishedAt:       d.FinishedAt,
	}
}
