package user

import (
	"context"
	"github.com/keweegen/tic-toe/internal/helper"
	"time"
)

type User struct {
	ID        string
	Name      string
	Email     string
	Locale    helper.Locale
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Repository interface {
	One(ctx context.Context, id string) (User, error)
	Add(ctx context.Context, user User) (User, error)
	Update(ctx context.Context, user User) (User, error)
	Delete(ctx context.Context, id string) error
}

type CreateParams struct {
	Name   string
	Email  string
	Locale helper.Locale
}

type UpdateParams struct {
	Name   string
	Locale helper.Locale
}

type Service interface {
	One(ctx context.Context, id string) (User, error)
	Create(ctx context.Context, params CreateParams) (User, error)
	Update(ctx context.Context, id string, params UpdateParams) (User, error)
	Delete(ctx context.Context, id string) error
}
