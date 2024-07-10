package service

import (
	"context"
	domain "github.com/keweegen/tic-toe/internal/domain/user"
)

var _ domain.Service = (*UserService)(nil)

type UserService struct {
	userRepo domain.Repository
}

func NewUserService(userRepo domain.Repository) *UserService {
	return &UserService{
		userRepo: userRepo,
	}
}

func (s *UserService) One(ctx context.Context, id string) (domain.User, error) {
	return s.userRepo.One(ctx, id)
}

func (s *UserService) Create(ctx context.Context, params domain.CreateParams) (domain.User, error) {
	user := domain.User{
		Name:   params.Name,
		Email:  params.Email,
		Locale: params.Locale,
	}
	return s.userRepo.Add(ctx, user)
}

func (s *UserService) Update(ctx context.Context, id string, params domain.UpdateParams) (domain.User, error) {
	user, err := s.One(ctx, id)
	if err != nil {
		return domain.User{}, err
	}

	user.Name = params.Name
	user.Locale = params.Locale

	return s.userRepo.Update(ctx, user)
}

func (s *UserService) Delete(ctx context.Context, id string) error {
	return s.userRepo.Delete(ctx, id)
}
