package dto

import (
	"database/sql"
	db "github.com/EliriaT/SchoolAppApi/db/sqlc"
	"time"
)

type CreateMarkRequest struct {
	CourseID  int64     `json:"course_id" binding:"required"`
	MarkDate  time.Time `json:"mark_date" binding:"required"`
	IsAbsent  bool      `json:"is_absent"`
	Mark      int32     `json:"mark"`
	StudentID int64     `json:"student_id" binding:"required"`
}

type MarkResponse struct {
	MarkID    int64         `json:"mark_id"`
	CourseID  int64         `json:"course_id"`
	MarkDate  time.Time     `json:"mark_date"`
	IsAbsent  sql.NullBool  `json:"is_absent"`
	Mark      sql.NullInt32 `json:"mark"`
	StudentID int64         `json:"student_id"`
	CreatedBy sql.NullInt64 `json:"createdBy"`
	UpdatedBy sql.NullInt64 `json:"updatedBy"`
	CreatedAt sql.NullTime  `json:"createdAt"`
	UpdatedAt sql.NullTime  `json:"updatedAt"`
}

type UpdateMarkRequest struct {
	MarkID    int64     `json:"mark_id" binding:"required"`
	CourseID  int64     `json:"course_id" binding:"required"`
	MarkDate  time.Time `json:"mark_date" binding:"required"`
	IsAbsent  bool      `json:"is_absent"`
	Mark      int32     `json:"mark" binding:"gte=1,lte=10"`
	StudentID int64     `json:"student_id" binding:"required"`
}

type DeleteMarkRequest struct {
	MarkID int64 `uri:"id" binding:"required"`
}

func NewMarkResponse(mark db.Mark) MarkResponse {
	return MarkResponse{
		mark.ID,
		mark.CourseID,
		mark.MarkDate,
		mark.IsAbsent,
		mark.Mark,
		mark.StudentID,
		mark.CreatedBy,
		mark.UpdatedBy,
		mark.CreatedAt,
		mark.UpdatedAt,
	}

}
