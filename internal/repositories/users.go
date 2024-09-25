package repositories

import (
	"bookstore_api/models"
	"context"
	"errors"
	"fmt"
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
	Update(ctx context.Context, user *models.User) (*models.User, error)
}

func (repo *UserRepository) Get(ctx context.Context, email string) (*models.User, error) {
	user := &models.User{}
	err := repo.Db.GetContext(ctx, user, "SELECT * FROM users WHERE email = $1", email)
	if err != nil {
		return nil, fmt.Errorf("error getting user: %s", err)
	}

	return user, nil
}

func (repo *UserRepository) Register(ctx context.Context, user *models.UserRegister) (*models.User, error) {
	query := `		
		WITH email_conflict AS (
			SELECT id FROM users WHERE email = :email
		)
		INSERT INTO users (name, email, password)
		SELECT :name, :email, :password
		WHERE NOT EXISTS (SELECT 1 FROM email_conflict)
		RETURNING *;
	`

	// Prepare the named query
	stmt, err := repo.Db.PrepareNamedContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("error preparing query: %v", err)
	}

	registeredUser := &models.User{}
	err = stmt.GetContext(ctx, registeredUser, user)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			return nil, fmt.Errorf("error creating user: email already exists")
		}
		return nil, fmt.Errorf("error creating user: %v", err)
	}

	return registeredUser, nil
}

func (repo *UserRepository) Update(ctx context.Context, user *models.User) (*models.User, error) {
	query := `
		WITH email_conflict AS (
			SELECT id FROM users WHERE email = :email AND NOT id = :id
		)
		UPDATE users
		SET name=:name, email=:email, password=:password, updated_at=:updated_at 
		WHERE id=:id AND NOT EXISTS (SELECT 1 FROM email_conflict)
		RETURNING *
	`

	stmt, err := repo.Db.PrepareNamedContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("error preparing query: %v", err)
	}

	updatedUser := &models.User{}
	err = stmt.GetContext(ctx, updatedUser, user)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			return nil, errors.New("error updating user: email already exists")
		}
		return nil, fmt.Errorf("error updating user: %v", err)
	}

	return updatedUser, nil
}
