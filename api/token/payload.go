package token

import (
	"errors"
	"github.com/google/uuid"
	"time"
)

// Different types of error returned by the VerifyToken function
var (
	ErrInvalidToken = errors.New("token is invalid")
	ErrExpiredToken = errors.New("token has expired")
)

// Payload contains the payload data of the token
type Payload struct {
	ID        uuid.UUID `json:"id"`
	Email     string    `json:"email"`
	Role      []int64   `json:"role"`
	SchoolID  int64     `json:"school_id"`
	ClassID   int64     `json:"class_id"`
	UserID    int64     `json:"user_id"`
	IssuedAt  time.Time `json:"issued_at"`
	ExpiredAt time.Time `json:"expired_at"`
	//jwt.RegisteredClaims
}

// Valid checks if the token payload is valid or not
// TODO Validate all cases and with IP
func (payload *Payload) Valid() error {
	if time.Now().After(payload.ExpiredAt) {
		return ErrExpiredToken
	}
	return nil
}

func NewPayload(email string, role []int64, SchoolID int64, ClassID int64, UserID int64, duration time.Duration) (*Payload, error) {
	tokenId, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}

	// I should have here User ID which is a uuid, and user role
	payload := &Payload{
		ID:        tokenId,
		Email:     email,
		Role:      role,
		SchoolID:  SchoolID,
		ClassID:   ClassID,
		UserID:    UserID,
		IssuedAt:  time.Now(),
		ExpiredAt: time.Now().Add(duration),
		//RegisteredClaims: jwt.RegisteredClaims{
		//	ExpiresAt: jwt.NewNumericDate(time.Now().Add(duration)),
		//	IssuedAt:  jwt.NewNumericDate(time.Now()),
		//	NotBefore: jwt.NewNumericDate(time.Now()),
		//},
	}
	return payload, nil
}
