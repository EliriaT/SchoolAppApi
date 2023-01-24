package service

import (
	"context"
	db "github.com/EliriaT/SchoolAppApi/db/sqlc"
	"github.com/EliriaT/SchoolAppApi/service/dto"
	"github.com/google/uuid"
)

type SessionService interface {
	GetSession(ctx context.Context, tokenID uuid.UUID) (dto.SessionResponse, error)
}

type sessionService struct {
	db db.Store
	//RolesService
	//roles map[string]db.Role
}

func (s *sessionService) GetSession(ctx context.Context, tokenID uuid.UUID) (dto.SessionResponse, error) {
	session, err := s.db.GetSession(ctx, tokenID)
	if err != nil {
		return dto.SessionResponse{}, err
	}
	return dto.NewSessionResponse(session), nil
}

func NewSessionService(database db.Store) SessionService {
	return &sessionService{db: database}
}
