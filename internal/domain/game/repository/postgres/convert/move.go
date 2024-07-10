package convert

import (
	"encoding/json"
	"github.com/keweegen/tic-toe/internal/db/model"
	domain "github.com/keweegen/tic-toe/internal/domain/game"
)

var Move move

type move struct{}

func (move) Domain(m *model.GameMove) (domain.Move, error) {
	var position domain.Position
	if err := json.Unmarshal(m.Position, &position); err != nil {
		return domain.Move{}, err
	}

	return domain.Move{
		ID:        m.ID,
		GameID:    m.GameID,
		StrikerID: m.StrikerID,
		BatterID:  m.BatterID,
		UnitCode:  domain.UnitCode(m.UnitCode),
		Position:  position,
		CreatedAt: m.CreatedAt,
	}, nil
}

func (move) DomainSlice(ms model.GameMoveSlice) ([]domain.Move, error) {
	moves := make([]domain.Move, len(ms))
	for i, m := range ms {
		move, err := Move.Domain(m)
		if err != nil {
			return nil, err
		}
		moves[i] = move
	}
	return moves, nil
}

func (move) Model(d domain.Move) (*model.GameMove, error) {
	position, err := json.Marshal(d.Position)
	if err != nil {
		return nil, err
	}

	return &model.GameMove{
		ID:        d.ID,
		GameID:    d.GameID,
		StrikerID: d.StrikerID,
		BatterID:  d.BatterID,
		UnitCode:  d.UnitCode.String(),
		Position:  position,
		CreatedAt: d.CreatedAt,
	}, nil
}
