package convert

import (
	"github.com/keweegen/tic-toe/internal/db/model"
	domain "github.com/keweegen/tic-toe/internal/domain/user"
	"github.com/keweegen/tic-toe/internal/helper"
)

var User user

type user struct{}

func (user) Domain(m *model.User) domain.User {
	return domain.User{
		ID:        m.ID,
		Name:      m.Name,
		Email:     m.Email,
		Locale:    helper.Locale(m.Locale),
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	}
}

func (user) Model(d domain.User) *model.User {
	return &model.User{
		ID:        d.ID,
		Name:      d.Name,
		Email:     d.Email,
		Locale:    d.Locale.String(),
		CreatedAt: d.CreatedAt,
		UpdatedAt: d.UpdatedAt,
	}
}
