package api

import (
	"database/sql"
	"github.com/EliriaT/SchoolAppApi/api/token"
	"github.com/EliriaT/SchoolAppApi/service"
	"github.com/EliriaT/SchoolAppApi/service/dto"
	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
)

func (server *Server) createUser(ctx *gin.Context) {
	var req dto.CreateUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	log.Println(ctx.ClientIP())

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	response, err := server.service.Register(ctx, authPayload, req)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "foreign_key_violation", "unique_violation":
				ctx.JSON(http.StatusForbidden, errorResponse(err))
				return
			}
		}
		switch err.Error() {
		case service.ErrBadRequest.Error():
			ctx.JSON(http.StatusBadRequest, errorResponse(err))
			return

		case sql.ErrNoRows.Error():
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return

		case service.ErrUnAuthorized.Error():
			ctx.JSON(http.StatusUnauthorized, errorResponse(err))
			return
		default:
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}

	}

	ctx.JSON(http.StatusCreated, response)
}

func (server *Server) loginUser(ctx *gin.Context) {
	var req dto.LoginUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	response, roles, schoolID, classID, err := server.service.Login(ctx, req)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		} else if err == bcrypt.ErrMismatchedHashAndPassword {
			ctx.JSON(http.StatusUnauthorized, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	accessToken, err := server.tokenMaker.CreateToken(response.User.Email, roles, schoolID, classID, response.User.ID, server.config.AccessTokenDuration)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	response.AccessToken = accessToken

	ctx.JSON(http.StatusOK, response)
}

func (server *Server) twoFactorLoginUser(ctx *gin.Context) {
	var req dto.CheckTOTPRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	response, err := server.service.CheckTOTP(ctx, authPayload, req)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	authToken, err := server.tokenMaker.AuthenticateToken(*authPayload)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	response.AccessToken = authToken
	ctx.JSON(http.StatusOK, response)
}

func (server *Server) getRoles(ctx *gin.Context) {
	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	response, err := server.service.GetRoles(ctx, authPayload)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, response)
}

func (server *Server) getTeacher(ctx *gin.Context) {
	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	response, err := server.service.GetTeachers(ctx, authPayload)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, response)
}
