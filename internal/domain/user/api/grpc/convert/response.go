package convert

import (
	userpb "github.com/keweegen/tic-toe/api/grpc/user"
	domain "github.com/keweegen/tic-toe/internal/domain/user"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var Response response

type response struct{}

func (response) User(u domain.User) *userpb.Response {
	return &userpb.Response{
		Id:        u.ID,
		Name:      u.Name,
		Email:     u.Email,
		Locale:    string(u.Locale),
		CreatedAt: timestamppb.New(u.CreatedAt),
		UpdatedAt: timestamppb.New(u.UpdatedAt),
	}
}
