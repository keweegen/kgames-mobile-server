package game

import (
	"context"
	"time"
)

type Move struct {
	ID        string
	GameID    string
	StrikerID string
	BatterID  string
	UnitCode  UnitCode
	Position  Position
	CreatedAt time.Time
}

type Position struct {
	X int `json:"x,omitempty"`
	Y int `json:"y,omitempty"`
}

type MoveRepository interface {
	One(ctx context.Context, id string) (Move, error)
	Add(ctx context.Context, move Move) (Move, error)
}

type Mover interface {
	Do(ctx context.Context, game Game, newMove Move) (move Move, finished bool, params FinishParams, err error)
}
