package convert

import (
	"github.com/keweegen/tic-toe/internal/db/model"
	domain "github.com/keweegen/tic-toe/internal/domain/game"
)

var Player player

type player struct{}

func (player) Domain(m *model.GamePlayer) domain.Player {
	return domain.Player{
		ID:     m.PlayerID,
		GameID: m.GameID,
		Profit: m.Profit,
		Fee:    m.Fee,
		GaveUp: m.GaveUp,
		Draw:   m.Draw,
	}
}

func (player) DomainSlice(ms model.GamePlayerSlice) []domain.Player {
	players := make([]domain.Player, len(ms))
	for i, m := range ms {
		players[i] = Player.Domain(m)
	}
	return players
}

func (player) Model(d domain.Player) *model.GamePlayer {
	return &model.GamePlayer{
		PlayerID: d.ID,
		GameID:   d.GameID,
		Profit:   d.Profit,
		Fee:      d.Fee,
		GaveUp:   d.GaveUp,
		Draw:     d.Draw,
	}
}

func (player) ModelSlice(ds []domain.Player) model.GamePlayerSlice {
	players := make(model.GamePlayerSlice, len(ds))
	for i, d := range ds {
		players[i] = Player.Model(d)
	}
	return players
}
