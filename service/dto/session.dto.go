package dto

import (
	db "github.com/EliriaT/SchoolAppApi/db/sqlc"
	"github.com/google/uuid"
	"time"
)

type SessionResponse struct {
	ID           uuid.UUID `json:"id"`
	Email        string    `json:"email"`
	RefreshToken string    `json:"refreshToken"`
	UserAgent    string    `json:"userAgent"`
	ClientIp     string    `json:"clientIp"`
	IsBlocked    bool      `json:"isBlocked"`
	ExpiresAt    time.Time `json:"expiresAt"`
	CreatedAt    time.Time `json:"createdAt"`
}

type RenewAccessTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type RenewAccessTokenResponse struct {
	AccessToken          string    `json:"access_token"`
	AccessTokenExpiresAt time.Time `json:"access_token_expires_at"`
}

func NewSessionResponse(session db.Session) SessionResponse {
	return SessionResponse{
		ID:           session.ID,
		Email:        session.Email,
		RefreshToken: session.RefreshToken,
		UserAgent:    session.UserAgent,
		ClientIp:     session.ClientIp,
		IsBlocked:    session.IsBlocked,
		ExpiresAt:    session.ExpiresAt,
		CreatedAt:    session.CreatedAt,
	}
}
