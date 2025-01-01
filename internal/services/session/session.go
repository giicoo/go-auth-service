package session

import (
	"context"
	"fmt"

	"github.com/giicoo/go-auth-service/internal/entity"
	"github.com/giicoo/go-auth-service/internal/repository"
	"github.com/giicoo/go-auth-service/internal/repository/redisRepo"
)

type SessionService struct {
	ctx  context.Context
	repo repository.SessionRepo
}

func NewSessionManager() *SessionService {
	return &SessionService{
		ctx:  context.Background(),
		repo: redisRepo.NewSessionRepo(),
	}
}

func (sm *SessionService) CreateSession(s *entity.Session) (*entity.Session, error) {
	id, err := sm.repo.CreateSession(sm.ctx, s)
	if err != nil {
		return nil, fmt.Errorf("session manager create: %w", err)
	}

	return id, nil
}

func (sm *SessionService) GetSession(id string, user_id int) (*entity.Session, error) {
	session, err := sm.repo.GetSession(sm.ctx, id, user_id)
	if err != nil {
		return nil, fmt.Errorf("session manager get '%s:%d': %w", id, user_id, err)
	}

	return session, nil
}

func (sm *SessionService) DeleteSession(id string, user_id int) error {
	if err := sm.repo.DeleteSession(sm.ctx, id, user_id); err != nil {
		return fmt.Errorf("session manager delete '%s:%d': %w", id, user_id, err)
	}
	return nil
}

func (sm *SessionService) GetListSession(user_id int) ([]*entity.Session, error) {
	sessions, err := sm.repo.GetListSession(sm.ctx, user_id)
	if err != nil {
		return nil, fmt.Errorf("session manager get list '%d': %w", user_id, err)
	}
	return sessions, err
}

func (sm *SessionService) DeleteListSession(user_id int) error {
	if err := sm.repo.DeleteListSession(sm.ctx, user_id); err != nil {
		return fmt.Errorf("session manager del list '%d': %w", user_id, err)
	}
	return nil
}
