package convert

import (
	userpb "github.com/keweegen/tic-toe/api/grpc/user"
	domain "github.com/keweegen/tic-toe/internal/domain/user"
	"github.com/keweegen/tic-toe/internal/helper"
)

var Request request

type request struct{}

func (request) CreateUser(req *userpb.CreateRequest) domain.CreateParams {
	return domain.CreateParams{
		Name:   req.Name,
		Email:  req.Email,
		Locale: helper.Locale(req.Locale),
	}
}

func (request) UpdateUser(req *userpb.UpdateRequest) domain.UpdateParams {
	return domain.UpdateParams{
		Name:   req.Name,
		Locale: helper.Locale(req.Locale),
	}
}
