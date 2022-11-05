package api

import (
	"database/sql"
	"github.com/EliriaT/SchoolAppApi/db/service"
	db "github.com/EliriaT/SchoolAppApi/db/sqlc"
	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
	"log"
	"net/http"
	"time"
)

type createUserRequest struct {
	Email       string    `json:"email" form:"email" binding:"required,email"`
	Password    string    `json:"password" form:"password" binding:"required,min=6"`
	LastName    string    `json:"lastName" form:"lastName" binding:"required,alpha"`
	FirstName   string    `json:"firstName" form:"firstName" binding:"required,alpha"`
	Gender      string    `json:"gender" form:"gender" binding:"required,oneof=F M"`
	PhoneNumber string    `json:"phoneNumber" form:"phoneNumber" binding:"required,e164"`
	Domicile    string    `json:"domicile" form:"domicile"`
	BirthDate   time.Time `json:"birthDate" form:"birthDate" binding:"required" time_format:"2006-01-02"`
}

type userResponse struct {
	ID                int64          `json:"id"`
	Email             string         `json:"email"`
	LastName          string         `json:"lastName"`
	FirstName         string         `json:"firstName"`
	Gender            string         `json:"gender"`
	PhoneNumber       sql.NullString `json:"phoneNumber"`
	Domicile          sql.NullString `json:"domicile"`
	BirthDate         time.Time      `json:"birthDate"`
	PasswordChangedAt time.Time      `json:"passwordChangedAt"`
	CreatedAt         time.Time      `json:"createdAt"`
}

func newUserResponse(user db.User) userResponse {
	return userResponse{
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

func (server *Server) createUser(ctx *gin.Context) {
	var req createUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Println("Here")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	hashedPassword, err := service.HashPassword(req.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	arg := db.CreateUserParams{
		Email:     req.Email,
		Password:  hashedPassword,
		LastName:  req.LastName,
		FirstName: req.FirstName,
		Gender:    req.Gender,
		PhoneNumber: sql.NullString{
			String: req.PhoneNumber,
			Valid:  true,
		},
		Domicile: sql.NullString{
			String: req.Domicile,
			Valid:  true,
		},
		BirthDate: req.BirthDate,
	}

	user, err := server.store.CreateUser(ctx, arg)

	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "unique_violation":
				ctx.JSON(http.StatusForbidden, errorResponse(err))
				return
			}
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	response := newUserResponse(user)
	ctx.JSON(http.StatusOK, response)
}

type loginUserRequest struct {
	Email    string `json:"email" form:"email" binding:"required,email"`
	Password string `json:"password" form:"password" binding:"required,min=6"`
}

type loginUserResponse struct {
	AccessToken string `json:"access_token"`
	//Here i should set role, maybe user id, or should i set it in AccessToken ?
	User userResponse
}

func (server *Server) loginUser(ctx *gin.Context) {
	var req loginUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	user, err := server.store.GetUserbyEmail(ctx, req.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	err = service.CheckPassword(req.Password, user.Password)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	accessToken, err := server.tokenMaker.CreateToken(user.Email, server.config.AccessTokenDuration)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
	}

	response := loginUserResponse{
		AccessToken: accessToken,
		User:        newUserResponse(user),
	}
	ctx.JSON(http.StatusOK, response)
}
