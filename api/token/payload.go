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
	Email     string    `json:"Email"`
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

func NewPayload(email string, duration time.Duration) (*Payload, error) {
	tokenId, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}

	payload := &Payload{
		ID:        tokenId,
		Email:     email,
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
