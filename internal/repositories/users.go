package repositories

import (
	"bookstore_api/models"
	"context"
)

type UserRepository struct {
	*Repository
}

func NewUserRepository(repository *Repository) *UserRepository {
	return &UserRepository{
		repository,
	}
}

type IUserRepository interface {
	Get(ctx context.Context, email string) (*models.User, error)
	Register(ctx context.Context, user *models.UserRegister) (*models.User, error)
	Login(ctx context.Context, user *models.UserLogin) (*models.User, error)
	Update(ctx context.Context, user *models.User) (*models.User, error)
}

func (repo *UserRepository) Get(ctx context.Context, email string) (*models.User, error) {
	return nil, nil
}

func (repo *UserRepository) Register(ctx context.Context, user *models.UserRegister) (*models.User, error) {
	return nil, nil
}

func (repo *UserRepository) Login(ctx context.Context, user *models.UserLogin) (*models.User, error) {
	return nil, nil
}

func (repo *UserRepository) Update(ctx context.Context, user *models.User) (*models.User, error) {
	return nil, nil
}
