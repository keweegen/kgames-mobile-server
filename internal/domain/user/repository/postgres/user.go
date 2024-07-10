package postgres

import (
	"context"
	"database/sql"
	"errors"
	"github.com/keweegen/tic-toe/internal/db/model"
	domain "github.com/keweegen/tic-toe/internal/domain/user"
	domainerrs "github.com/keweegen/tic-toe/internal/domain/user/errors"
	"github.com/keweegen/tic-toe/internal/domain/user/repository/postgres/convert"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

type UserRepository struct {
	db *sql.DB
}

var _ domain.Repository = (*UserRepository)(nil)

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (r UserRepository) One(ctx context.Context, id string) (domain.User, error) {
	mods := []qm.QueryMod{
		model.UserWhere.ID.EQ(id),
	}

	m, err := model.Users(mods...).One(ctx, r.db)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.User{}, domainerrs.ErrUserNotFound
		}
		return domain.User{}, err
	}

	return convert.User.Domain(m), nil
}

func (r UserRepository) Add(ctx context.Context, user domain.User) (domain.User, error) {
	m := convert.User.Model(user)
	if err := m.Insert(ctx, r.db, boil.Infer()); err != nil {
		return domain.User{}, err
	}
	return convert.User.Domain(m), nil
}

func (r UserRepository) Update(ctx context.Context, player domain.User) (domain.User, error) {
	updateColumns := boil.Whitelist(
		model.UserColumns.Name,
		model.UserColumns.Locale,
		model.UserColumns.UpdatedAt,
	)

	m := convert.User.Model(player)
	if _, err := m.Update(ctx, r.db, updateColumns); err != nil {
		return domain.User{}, err
	}

	return convert.User.Domain(m), nil
}

func (r UserRepository) Delete(ctx context.Context, id string) error {
	mods := []qm.QueryMod{
		model.UserWhere.ID.EQ(id),
	}

	if _, err := model.Users(mods...).DeleteAll(ctx, r.db); err != nil {
		return err
	}

	return nil
}
