package api

import (
	"github.com/EliriaT/SchoolAppApi/api/token"
	"github.com/EliriaT/SchoolAppApi/service"
	"github.com/EliriaT/SchoolAppApi/service/dto"
	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
	"net/http"
)

func (server *Server) createMark(ctx *gin.Context) {
	var req dto.CreateMarkRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	response, err := server.service.CreateMark(ctx, authPayload, req)

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

func (server *Server) changeMark(ctx *gin.Context) {
	var req dto.UpdateMarkRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	response, err := server.service.ChangeMark(ctx, authPayload, req)
	if err != nil {
		if err == service.ErrUnAuthorized {
			ctx.JSON(http.StatusUnauthorized, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, response)
}

func (server *Server) deleteMark(ctx *gin.Context) {
	var req dto.DeleteMarkRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	err := server.service.DeleteMark(ctx, authPayload, req)
	if err != nil {
		if err == service.ErrUnAuthorized {
			ctx.JSON(http.StatusUnauthorized, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusNoContent, dto.MarkResponse{})
}
