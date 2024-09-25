package repositories

import (
	"bookstore_api/models"
	"context"
	"fmt"
	"log"
)

type SessionRepository struct {
	*Repository
}

func NewSessionRepository(repository *Repository) *SessionRepository {
	return &SessionRepository{
		repository,
	}
}

type ISessionRepository interface {
	Create(ctx context.Context, session *models.Sessions) (*models.Sessions, error)
	Get(ctx context.Context, id string) (*models.Sessions, error)
	Revoke(ctx context.Context, id string) error
	Delete(ctx context.Context, id string) error
}

func (repo *SessionRepository) Create(ctx context.Context, session *models.Sessions) (*models.Sessions, error) {
	query := `
		INSERT INTO Sessions (id, user_email, refresh_token, expires_at) 
		VALUES (:id, :user_email, :refresh_token, :expires_at) 
		RETURNING *
    `

	_, err := repo.Db.NamedExecContext(ctx, query, session)
	if err != nil {
		log.Printf("Here is the sql statement error: %v", err)
		return nil, err
	}

	return session, nil
}

func (repo *SessionRepository) Get(ctx context.Context, id string) (*models.Sessions, error) {
	session := &models.Sessions{}
	err := repo.Db.GetContext(ctx, session, "SELECT * FROM sessions WHERE id = $1", id)
	if err != nil {
		return nil, fmt.Errorf("error getting session: %s", err)
	}

	return session, nil
}

func (repo *SessionRepository) Revoke(ctx context.Context, id string) error {
	_, err := repo.Db.NamedExecContext(ctx, "UPDATE sessions SET is_revoked=TRUE WHERE id = :id", map[string]interface{}{"id": id})
	if err != nil {
		return fmt.Errorf("error revoking session: %s", err)
	}

	return nil
}

func (repo *SessionRepository) Delete(ctx context.Context, id string) error {
	_, err := repo.Db.ExecContext(ctx, "DELETE FROM sessions WHERE id = $1", id)
	if err != nil {
		return fmt.Errorf("error deleting session: %s", err)
	}

	return nil
}
