package game

import (
	"context"
	"github.com/shopspring/decimal"
	"github.com/volatiletech/null/v8"
	"time"
)

type Game struct {
	ID               string
	TypeCode         TypeCode
	StateCode        StateCode
	Bid              decimal.Decimal
	FinishReasonCode FinishReasonCode
	Players          Players
	MaxPlayers       int
	Moves            []Move
	CreatedAt        time.Time
	CreatedBy        string
	StartedAt        null.Time
	FinishedAt       null.Time
}

func (g Game) Finished() bool {
	return !g.FinishedAt.IsZero() && g.StateCode == StateCodeFinished
}

type Repository interface {
	Add(ctx context.Context, game Game) (Game, error)
	One(ctx context.Context, id string) (Game, error)
	Update(ctx context.Context, game Game) error
}

type CreateParams struct {
	InitiatorID string
	TypeCode    TypeCode
	Bid         decimal.Decimal
	MaxPlayers  int
}

type MoveParams struct {
	StrikerID string
	UnitCode  UnitCode
}

type FinishParams struct {
	ReasonCode  FinishReasonCode
	InitiatorID string
	Players     []PlayerFinishParams
}

type PlayerFinishParams struct {
	PlayerID string
	Profit   decimal.Decimal
	Fee      decimal.Decimal
	Draw     bool
	GaveUp   bool
}

type Service interface {
	One(ctx context.Context, id string) (Game, error)
	Create(ctx context.Context, params CreateParams) (Game, error)
	Update(ctx context.Context, game Game) error
	Start(ctx context.Context, id string) (Game, error)
	Move(ctx context.Context, id string, params MoveParams) (Move, error)
	GaveUp(ctx context.Context, id string, playerID string) error
	Finish(ctx context.Context, id string, params FinishParams) error
	Notify(ctx context.Context, gameID, initiatorID string, action StreamAction, data any)
}

type StreamMessage struct {
	GameID string
	UserID string
	Action StreamAction
	Data   any
}

type StreamService interface {
	Handle(ctx context.Context, msg StreamMessage) error
}
