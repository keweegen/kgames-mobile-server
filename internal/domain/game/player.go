package game

import (
	"github.com/shopspring/decimal"
	"sort"
)

type Player struct {
	ID        string
	GameID    string
	StateCode PlayerStateCode
	Position  int
	Profit    decimal.Decimal
	Fee       decimal.Decimal
	GaveUp    bool
	Draw      bool
}

type Players []Player

func (p Players) First() Player {
	if len(p) == 0 {
		return Player{}
	}
	return p[0]
}

func (p Players) OneByPosition(position int) Player {
	for _, player := range p {
		if player.Position == position {
			return player
		}
	}
	return Player{}
}

func (p Players) Map() MapPlayers {
	m := make(map[string]Player, len(p))
	for _, player := range p {
		m[player.ID] = player
	}
	return m
}

func (p Players) LenInt64() int64 {
	return int64(len(p))
}

func (p Players) LenDecimal() decimal.Decimal {
	return decimal.NewFromInt(p.LenInt64())
}

func (p Players) IDs() []string {
	unique := make(map[string]struct{}, len(p))
	for _, player := range p {
		unique[player.ID] = struct{}{}
	}

	ids := make([]string, 0, len(unique))
	for id := range unique {
		ids = append(ids, id)
	}

	return ids
}

func (p Players) Exists(id string) bool {
	for _, player := range p {
		if player.ID == id {
			return true
		}
	}
	return false
}

func (p Players) SortByPositions() Players {
	c := make(Players, len(p))
	copy(c, p)

	sort.Slice(c, func(i, j int) bool {
		return c[i].Position < c[j].Position
	})

	return c
}

type MapPlayers map[string]Player

func (p MapPlayers) Slice() Players {
	players := make(Players, 0, len(p))
	for _, player := range p {
		players = append(players, player)
	}

	sort.Slice(players, func(i, j int) bool {
		return players[i].Position < players[j].Position
	})

	return players
}
