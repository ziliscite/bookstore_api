package services

import (
	"bookstore_api/internal/repositories"
	"bookstore_api/models"
	"context"
)

type UserService struct {
	*Service
	userRepo repositories.IUserRepository
}

func NewUserService(service *Service, userRepo repositories.IUserRepository) *UserService {
	return &UserService{
		Service:  service,
		userRepo: userRepo,
	}
}

// Validate things from here, not handler

func (r *UserService) RegisterUser(ctx context.Context, user *models.UserRegister) (*models.User, error) {
	return r.userRepo.Register(ctx, user)
}

func (r *UserService) LoginUser(ctx context.Context, user *models.UserLogin) (*models.User, error) {
	return r.userRepo.Login(ctx, user)
}

func (r *UserService) UpdateUser(ctx context.Context, user *models.User) (*models.User, error) {
	return r.userRepo.Update(ctx, user)
}
