package token

import (
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"time"
)

const minSecretKeySize = 32

// JWTMaker is a JSON Web Token maker
type JWTMaker struct {
	secretKey string
}

// CreateToken creates a new token for a specific user with unique email,
func (j *JWTMaker) CreateToken(email string, role []int64, SchoolID int64, ClassID int64, UserID int64, duration time.Duration) (string, error) {
	payload, err := NewPayload(email, role, SchoolID, ClassID, UserID, duration)
	if err != nil {
		return "", err
	}

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)

	return jwtToken.SignedString([]byte(j.secretKey))

}

// VerifyToken checks if the tocken is valid, or not
func (j *JWTMaker) VerifyToken(token string) (*Payload, error) {
	keyFunc := func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, ErrInvalidToken
		}
		return []byte(j.secretKey), nil
	}

	jwtToken, err := jwt.ParseWithClaims(token, &Payload{}, keyFunc)
	if err != nil {
		verr, ok := err.(*jwt.ValidationError)
		if ok && errors.Is(verr.Inner, ErrExpiredToken) {
			return nil, ErrExpiredToken
		}
		return nil, ErrInvalidToken
	}

	payload, ok := jwtToken.Claims.(*Payload)
	if !ok {
		return nil, ErrInvalidToken
	}

	return payload, nil
}

// AuthenticateToken marks authentitcated field in the token payload as true, after 2fa is succesful,
func (j *JWTMaker) AuthenticateToken(payload Payload) (string, error) {

	payload.Authenticated = true

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, &payload)

	return jwtToken.SignedString([]byte(j.secretKey))

}

// Atention, not implemented
func (p *JWTMaker) CreatePasswordRecoveryToken(email string, duration time.Duration) (string, error) {
	return "", nil
}

// Atention, not implemented
func (p *JWTMaker) VerifyPasswordToken(token string) (PasswordRecoveryPayload, error) {
	return PasswordRecoveryPayload{}, nil
}

// NewJWTMaker creates a new JWTMaker that implements the TokenMaker interface
func NewJWTMaker(secretKey string) (TokenMaker, error) {
	if len(secretKey) < minSecretKeySize {
		return nil, fmt.Errorf("invalid key size: must be at least %d characthers length", minSecretKeySize)
	}
	return &JWTMaker{secretKey: secretKey}, nil
}
