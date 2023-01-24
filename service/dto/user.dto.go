package dto

import (
	"database/sql"
	db "github.com/EliriaT/SchoolAppApi/db/sqlc"
	"github.com/google/uuid"
	"time"
)

// When the User is created, if it is admin, it will indicate the school of the director/manager. Otherwise, the school is taken from the token.
// The class is indicated only for students. Teachers and Head Teachers are assigned separately.
type CreateUserRequest struct {
	Email string `json:"email" binding:"required,email"`
	//Password    string    `json:"password" binding:"required,min=10"`
	LastName    string    `json:"lastName" binding:"required,alpha"`
	FirstName   string    `json:"firstName"  binding:"required,alpha"`
	Gender      string    `json:"gender" binding:"required,oneof=F M"`
	PhoneNumber string    `json:"phoneNumber" binding:"required,e164"`
	Domicile    string    `json:"domicile"`
	BirthDate   time.Time `json:"birthDate" binding:"required" time_format:"2006-01-02"`
	SchoolID    int64     `json:"school"`
	RoleID      int64     `json:"role_id" binding:"required"`
	ClassID     int64     `json:"class_id"`
}

type UserResponse struct {
	ID                int64          `json:"id,omitempty"`
	Email             string         `json:"email,omitempty"`
	TOTPSecret        string         `json:"authentificator_secret,omitempty"`
	Qrcode            string         `json:"qrcode,omitempty"`
	LastName          string         `json:"lastName"`
	FirstName         string         `json:"firstName"`
	Gender            string         `json:"gender,omitempty"`
	PhoneNumber       sql.NullString `json:"phoneNumber,omitempty"`
	Domicile          sql.NullString `json:"domicile,omitempty"`
	BirthDate         sql.NullTime   `json:"birthDate,omitempty"`
	PasswordChangedAt time.Time      `json:"passwordChangedAt,omitempty"`
	CreatedAt         time.Time      `json:"createdAt,omitempty"`
	UserSchool        int64          `json:"user_school,omitempty"`
	UserRole          int64          `json:"user_role,omitempty"`
	UserClass         int64          `json:"user_class,omitempty"`
}

func NewUserResponse(user db.User) UserResponse {

	return UserResponse{
		ID:                user.ID,
		Email:             user.Email,
		LastName:          user.LastName,
		FirstName:         user.FirstName,
		Gender:            user.Gender,
		PhoneNumber:       user.PhoneNumber,
		Domicile:          user.Domicile,
		BirthDate:         user.BirthDate,
		PasswordChangedAt: user.PasswordChangedAt,
		CreatedAt:         user.CreatedAt,
	}
}

type LoginUserRequest struct {
	Email    string `json:"email" form:"email" binding:"required,email"`
	Password string `json:"password" form:"password" binding:"required,min=6"`
}

type LoginUserResponse struct {
	SessionID             uuid.UUID `json:"session_id"` //id of refresh token
	AccessToken           string    `json:"access_token"`
	AccessTokenExpiresAt  time.Time `json:"access_token_expires_at"`
	RefreshToken          string    `json:"refresh_token"`
	RefreshTokenExpiresAt time.Time `json:"refresh_token_expires_at"`
	User                  UserResponse
}

type CheckTOTPRequest struct {
	//Email     string `json:"email" form:"email" binding:"required,email"`
	TotpToken string `json:"totp_token" binding:"required"`
}

type CheckTOTPResponse struct {
	AccessToken string `json:"access_token"`
	User        UserResponse
}

type AccountRecoveryRequest struct {
	Email string `json:"email"`
}

type PasswordChangeURIRequest struct {
	Email string `uri:"email" binding:"required,email"`
	Token string `uri:"token" binding:"required"`
}

type PasswordChangeRequest struct {
	Password string `json:"password" binding:"required,min=10"`
}

type TeacherResponse struct {
	UserID     int64  `json:"user_id"`
	UserName   string `json:"user_name"`
	UserRoleID int64  `json:"user_role_id"`
	RoleName   string `json:"role_name"`
}
