package token

import "time"

type TokenMaker interface {
	// CreateToken creates a new token for a specific user with unique email,
	CreateToken(email string, role []int64, SchoolID int64, ClassID int64, UserID int64, duration time.Duration) (string, error)

	CreatePasswordRecoveryToken(email string, duration time.Duration) (string, error)

	//VerifyToken checks if the tocken is valid, or not
	VerifyToken(token string) (*Payload, error)

	VerifyPasswordToken(token string) (PasswordRecoveryPayload, error)

	// AuthentificateToken marks authentitcated field in the token payload as true, after 2fa is succesful,
	AuthenticateToken(payload Payload) (string, error)
}
