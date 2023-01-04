package api

import (
	"database/sql"
	"github.com/EliriaT/SchoolAppApi/api/token"
	"github.com/EliriaT/SchoolAppApi/service"
	"github.com/EliriaT/SchoolAppApi/service/dto"
	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
	"net/http"
)

func (server *Server) createSemester(ctx *gin.Context) {
	var req dto.CreateSemesterRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	response, err := server.service.CreateSemester(ctx, authPayload, req)

	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "foreign_key_violation", "unique_violation":
				ctx.JSON(http.StatusForbidden, errorResponse(err))
				return

			}
		}
		if err.Error() == service.ErrUnAuthorized.Error() {
			ctx.JSON(http.StatusUnauthorized, errorResponse(err))
			return
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusCreated, response)
}

// only a school manager can get its school, should not be from Id, but from userid that is in the token payload
func (server *Server) getSemesters(ctx *gin.Context) {

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	response, err := server.service.ListSemesters(ctx, authPayload)
	if err != nil {
		if err.Error() == sql.ErrNoRows.Error() {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		} else if err == service.ErrUnAuthorized {
			ctx.JSON(http.StatusUnauthorized, errorResponse(err))
			return
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, response)
}

// only the admin can list schools
func (server *Server) getCurrentSemester(ctx *gin.Context) {

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	response, err := server.service.GetCurrentSemester(ctx, authPayload)
	if err != nil {
		if err.Error() == sql.ErrNoRows.Error() {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		} else if err == service.ErrUnAuthorized {
			ctx.JSON(http.StatusUnauthorized, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, response)
}
