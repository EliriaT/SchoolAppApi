package api

import (
	"database/sql"
	"errors"
	"github.com/google/uuid"
	"log"
	"net/http"
	"time"

	"github.com/EliriaT/SchoolAppApi/api/token"
	"github.com/EliriaT/SchoolAppApi/service"
	"github.com/EliriaT/SchoolAppApi/service/dto"
	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

var incorrectCredentialsError = errors.New("Incorrect email or password")

func (server *Server) createUser(ctx *gin.Context) {
	var req dto.CreateUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	log.Println(ctx.ClientIP())

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	tokenPassReset, err := server.tokenMaker.CreatePasswordRecoveryToken(req.Email, server.config.EmailRecoveryTokenDuration)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	response, err := server.service.Register(ctx, authPayload, req, tokenPassReset)
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
			ctx.JSON(http.StatusUnauthorized, errorResponse(incorrectCredentialsError))
			return
		} else if err == bcrypt.ErrMismatchedHashAndPassword {
			ctx.JSON(http.StatusUnauthorized, errorResponse(incorrectCredentialsError))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	accessToken, _, err := server.tokenMaker.CreateToken(response.User.Email, roles, schoolID, classID, response.User.ID, server.config.AccessTokenDuration)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, accessToken)
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

	accessToken, err := server.tokenMaker.AuthenticateToken(*authPayload)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	refreshToken, refreshPayload, err := server.tokenMaker.CreateToken(response.User.Email, authPayload.Role, authPayload.SchoolID, authPayload.ClassID, response.User.ID, server.config.RefreshTokenDuration)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	session, err := server.service.CreateSession(ctx, refreshToken, refreshPayload, ctx.Request.UserAgent(), ctx.ClientIP())
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	response.SessionID = session.ID
	response.AccessToken = accessToken
	response.AccessTokenExpiresAt = authPayload.ExpiredAt
	response.RefreshToken = refreshToken
	response.RefreshTokenExpiresAt = refreshPayload.ExpiredAt
	viorelResponse := struct {
		SessionID             uuid.UUID `json:"session_id"`
		AccessToken           string    `json:"access_token"`
		AccessTokenExpiresAt  time.Time `json:"access_token_expires_at"`
		RefreshToken          string    `json:"refresh_token"`
		RefreshTokenExpiresAt time.Time `json:"refresh_token_expires_at"`
	}{
		SessionID:             session.ID,
		AccessToken:           accessToken,
		AccessTokenExpiresAt:  authPayload.ExpiredAt,
		RefreshToken:          refreshToken,
		RefreshTokenExpiresAt: refreshPayload.ExpiredAt,
	}

	ctx.JSON(http.StatusOK, viorelResponse)
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

func (server *Server) recoverAndChangePassword(ctx *gin.Context) {
	var req dto.PasswordChangeURIRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	passwordPayload, err := server.tokenMaker.VerifyPasswordToken(req.Token)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
		return
	}
	if passwordPayload.Email != req.Email {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	var jsonReq dto.PasswordChangeRequest
	if err := ctx.ShouldBindJSON(&jsonReq); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	err = server.service.ChangePassword(ctx, req.Email, jsonReq.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusUnauthorized, errorResponse(err))
			return
		} else if err == service.ErrEasyPassword {
			ctx.JSON(http.StatusBadRequest, errorResponse(service.ErrEasyPassword))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.String(http.StatusOK, "Your password was reset!")
}

func (server *Server) accountRecoveryRequest(ctx *gin.Context) {
	var req dto.AccountRecoveryRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	resetToken, err := server.tokenMaker.CreatePasswordRecoveryToken(req.Email, server.config.EmailRecoveryTokenDuration)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	link := server.service.CreatePasswordRecoveryLink(resetToken, req.Email)
	err = server.service.SendChangePasswordEmail(link, req.Email)
	if err != nil {
		if err == service.ErrInvalidEmail {
			ctx.JSON(http.StatusUnauthorized, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.String(http.StatusOK, "Please check your email, for a reset password link")

}
