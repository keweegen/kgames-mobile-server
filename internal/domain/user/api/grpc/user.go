package grpc

import (
	"context"
	userpb "github.com/keweegen/tic-toe/api/grpc/user"
	domain "github.com/keweegen/tic-toe/internal/domain/user"
	"github.com/keweegen/tic-toe/internal/domain/user/api/grpc/convert"
	"google.golang.org/protobuf/types/known/emptypb"
)

var _ userpb.ServiceServer = (*UserHandler)(nil)

type UserHandler struct {
	service domain.Service
}

func NewUserHandler(service domain.Service) *UserHandler {
	return &UserHandler{
		service: service,
	}
}

func (h *UserHandler) Create(ctx context.Context, req *userpb.CreateRequest) (*userpb.Response, error) {
	u, err := h.service.Create(ctx, convert.Request.CreateUser(req))
	if err != nil {
		return nil, err
	}
	return convert.Response.User(u), nil
}

func (h *UserHandler) Get(ctx context.Context, req *userpb.GetRequest) (*userpb.Response, error) {
	u, err := h.service.One(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return convert.Response.User(u), nil
}

func (h *UserHandler) Update(ctx context.Context, req *userpb.UpdateRequest) (*userpb.Response, error) {
	u, err := h.service.Update(ctx, req.Id, convert.Request.UpdateUser(req))
	if err != nil {
		return nil, err
	}
	return convert.Response.User(u), nil
}

func (h *UserHandler) Delete(ctx context.Context, req *userpb.DeleteRequest) (*emptypb.Empty, error) {
	if err := h.service.Delete(ctx, req.Id); err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}
