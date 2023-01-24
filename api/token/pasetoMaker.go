package token

import (
	"fmt"
	"time"

	"github.com/vk-rv/pvx"
	"golang.org/x/crypto/chacha20poly1305"
)

// PasetoMaker is a PASETO token maker which implements the TokenMaker interface
type PasetoMaker struct {
	paseto       *pvx.ProtoV4Local
	symmetricKey *pvx.SymKey
}

// CreateToken creates a new token for a specific user with unique email,
func (p *PasetoMaker) CreateToken(email string, role []int64, SchoolID int64, ClassID int64, UserID int64, duration time.Duration) (string, *Payload, error) {
	payload, err := NewPayload(email, role, SchoolID, ClassID, UserID, duration)
	if err != nil {
		return "", payload, err
	}

	tokenStr, err := p.paseto.Encrypt(p.symmetricKey, payload, pvx.WithFooter(""))
	return tokenStr, payload, err

}

// VerifyToken checks if the tocken is valid, or not and returns the decrypted payload
func (p *PasetoMaker) VerifyToken(token string) (*Payload, error) {
	payload := &Payload{}

	err := p.paseto.Decrypt(token, p.symmetricKey, pvx.WithFooter("")).ScanClaims(payload)
	if err != nil {
		return nil, err
	}

	err = payload.Valid()
	if err != nil {
		return nil, err
	}
	return payload, nil
}

// AuthenticateToken marks authentitcated field in the token payload as true, after 2fa is succesful,
func (p *PasetoMaker) AuthenticateToken(payload Payload) (string, error) {

	payload.Authenticated = true

	return p.paseto.Encrypt(p.symmetricKey, &payload, pvx.WithFooter(""))
}

func (p *PasetoMaker) CreatePasswordRecoveryToken(email string, duration time.Duration) (string, error) {
	passwordRecoveryPayload, err := NewPasswordRecoveryPayload(email, duration)
	if err != nil {

		return "", err
	}

	return p.paseto.Encrypt(p.symmetricKey, &passwordRecoveryPayload, pvx.WithFooter(""))
}

// VerifyToken checks if the tocken is valid, or not and returns the decrypted payload
func (p *PasetoMaker) VerifyPasswordToken(token string) (PasswordRecoveryPayload, error) {
	payload := &PasswordRecoveryPayload{}

	err := p.paseto.Decrypt(token, p.symmetricKey, pvx.WithFooter("")).ScanClaims(payload)
	if err != nil {
		return PasswordRecoveryPayload{}, ErrInvalidToken
	}

	err = payload.Valid()
	if err != nil {
		return PasswordRecoveryPayload{}, err
	}
	return *payload, nil
}

// NewPasetoMaker creates a new PasetoMaker
func NewPasetoMaker(symmetricKey string) (TokenMaker, error) {
	if len(symmetricKey) != chacha20poly1305.KeySize {
		return nil, fmt.Errorf("invalid key size: must be at %d characthers length", chacha20poly1305.KeySize)
	}

	tokenMaker := &PasetoMaker{
		paseto:       pvx.NewPV4Local(),
		symmetricKey: pvx.NewSymmetricKey([]byte(symmetricKey), pvx.Version4),
	}
	return tokenMaker, nil
}
