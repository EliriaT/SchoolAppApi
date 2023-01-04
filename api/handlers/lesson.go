package api

import (
	"github.com/EliriaT/SchoolAppApi/api/token"
	"github.com/EliriaT/SchoolAppApi/service"
	"github.com/EliriaT/SchoolAppApi/service/dto"
	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
	"net/http"
)

func (server *Server) createLesson(ctx *gin.Context) {
	var req dto.CreateLessonRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	response, err := server.service.CreateLesson(ctx, authPayload, req)

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

func (server *Server) getLessons(ctx *gin.Context) {

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	response, err := server.service.GetLessons(ctx, authPayload)
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

func (server *Server) changeLesson(ctx *gin.Context) {
	var req dto.UpdateLessonParams
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	response, err := server.service.ChangeLesson(ctx, authPayload, req)
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

// lessons of a course /lesson/course/:id
func (server *Server) getCourseLessonsByCourseID(ctx *gin.Context) {
	var req dto.GetCourseLessonsRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	response, err := server.service.GetCourseLessons(ctx, authPayload, req)
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
