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

type GameRepository struct {
	db *sql.DB
}

var _ domain.Repository = (*GameRepository)(nil)

func NewGameRepository(db *sql.DB) GameRepository {
	return GameRepository{db: db}
}

func (r GameRepository) One(ctx context.Context, id string) (domain.Game, error) {
	mods := []qm.QueryMod{
		model.GameWhere.ID.EQ(id),
	}

	m, err := model.Games(mods...).One(ctx, r.db)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.Game{}, domainerrs.ErrGameNotFound
		}
		return domain.Game{}, err
	}

	return convert.Game.Domain(m)
}

func (r GameRepository) Add(ctx context.Context, game domain.Game) (domain.Game, error) {
	m := convert.Game.Model(game)

	if err := m.Insert(ctx, r.db, boil.Infer()); err != nil {
		return domain.Game{}, err
	}

	if len(game.Players) > 0 {
		players := convert.Player.ModelSlice(game.Players)
		if err := m.AddGamePlayers(ctx, r.db, true, players...); err != nil {
			return domain.Game{}, err
		}
	}

	return convert.Game.Domain(m)
}

func (r GameRepository) Update(ctx context.Context, game domain.Game) error {
	updateColumns := boil.Whitelist(
		model.GameColumns.StateCode,
		model.GameColumns.FinishReasonCode,
		model.GameColumns.StartedAt,
		model.GameColumns.FinishedAt,
	)

	m := convert.Game.Model(game)
	if _, err := m.Update(ctx, r.db, updateColumns); err != nil {
		return err
	}
	if err := r.updatePlayers(ctx, game.Players); err != nil {
		return err
	}

	return nil
}

func (r GameRepository) updatePlayers(ctx context.Context, players domain.Players) error {
	if len(players) == 0 {
		return nil
	}

	mPlayers := convert.Player.ModelSlice(players)
	updateColumns := boil.Whitelist(
		model.GamePlayerColumns.Profit,
		model.GamePlayerColumns.Fee,
		model.GamePlayerColumns.GaveUp,
		model.GamePlayerColumns.Draw,
	)

	for _, player := range mPlayers {
		if _, err := player.Update(ctx, r.db, updateColumns); err != nil {
			return err
		}
	}

	return nil
}
