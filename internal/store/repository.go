package store

import (
	"database/sql"
	gamedomain "github.com/keweegen/tic-toe/internal/domain/game"
	gamerepo "github.com/keweegen/tic-toe/internal/domain/game/repository/postgres"
	userdomain "github.com/keweegen/tic-toe/internal/domain/user"
	userrepo "github.com/keweegen/tic-toe/internal/domain/user/repository/postgres"
)

type Repository struct {
	User     userdomain.Repository
	Game     gamedomain.Repository
	GameMove gamedomain.MoveRepository
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{
		User:     userrepo.NewUserRepository(db),
		Game:     gamerepo.NewGameRepository(db),
		GameMove: gamerepo.NewMoveRepository(db),
	}
}
