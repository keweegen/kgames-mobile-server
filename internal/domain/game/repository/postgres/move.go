package postgres

import (
	"context"
	"database/sql"
	"errors"
	"github.com/keweegen/tic-toe/internal/db/model"
	domain "github.com/keweegen/tic-toe/internal/domain/game"
	domainerrs "github.com/keweegen/tic-toe/internal/domain/game/errors"
	"github.com/keweegen/tic-toe/internal/domain/game/repository/postgres/convert"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

type MoveRepository struct {
	db *sql.DB
}

var _ domain.MoveRepository = (*MoveRepository)(nil)

func NewMoveRepository(db *sql.DB) MoveRepository {
	return MoveRepository{db: db}
}

func (r MoveRepository) One(ctx context.Context, id string) (domain.Move, error) {
	mods := []qm.QueryMod{
		model.GameMoveWhere.ID.EQ(id),
	}

	m, err := model.GameMoves(mods...).One(ctx, r.db)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.Move{}, domainerrs.ErrMoveNotFound
		}
		return domain.Move{}, err
	}

	return convert.Move.Domain(m)
}

func (r MoveRepository) Add(ctx context.Context, move domain.Move) (domain.Move, error) {
	m, err := convert.Move.Model(move)
	if err != nil {
		return domain.Move{}, err
	}

	if err = m.Insert(ctx, r.db, boil.Infer()); err != nil {
		return domain.Move{}, err
	}

	return convert.Move.Domain(m)
}
