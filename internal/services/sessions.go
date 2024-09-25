package services

import (
	"bookstore_api/internal/repositories"
	"bookstore_api/models"
	"context"
)

type SessionService struct {
	*Service
	sessionRepo repositories.ISessionRepository
}

func NewSessionService(service *Service, sessionRepo repositories.ISessionRepository) *SessionService {
	return &SessionService{
		Service:     service,
		sessionRepo: sessionRepo,
	}
}

func (s *SessionService) CreateSession(ctx context.Context, session *models.Sessions) (*models.Sessions, error) {
	return s.sessionRepo.Create(ctx, session)
}

func (s *SessionService) GetSession(ctx context.Context, sessionID string) (*models.Sessions, error) {
	return s.sessionRepo.Get(ctx, sessionID)
}

func (s *SessionService) RevokeSession(ctx context.Context, sessionID string) error {
	return s.sessionRepo.Revoke(ctx, sessionID)
}

func (s *SessionService) DeleteSession(ctx context.Context, sessionID string) error {
	return s.sessionRepo.Delete(ctx, sessionID)
}
