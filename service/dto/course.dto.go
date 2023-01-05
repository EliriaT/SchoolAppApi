package dto

import (
	db "github.com/EliriaT/SchoolAppApi/db/sqlc"
	"time"
)

type CreateCourseRequest struct {
	Name       string `json:"name" binding:"required"`
	TeacherID  int64  `json:"teacher_id" binding:"required"`
	SemesterID int64  `json:"semester_id" binding:"required"`
	ClassID    int64  `json:"class_id"  binding:"required"`
}
type CourseResponse struct {
	ID         int64     `json:"id,omitempty"`
	Name       string    `json:"name"`
	TeacherID  int64     `json:"teacher_id"`
	SemesterID int64     `json:"semester_id"`
	ClassID    int64     `json:"class_id"`
	CreatedBy  int64     `json:"createdBy"`
	UpdatedBy  int64     `json:"updatedBy"`
	CreatedAt  time.Time `json:"createdAt"`
	UpdatedAt  time.Time `json:"updatedAt"`
}

type UpdateCourseParams struct {
	CourseID   int64  `json:"id" binding:"required"`
	Name       string `json:"name"`
	TeacherID  int64  `json:"teacher_id"`
	SemesterID int64  `json:"semester_id"`
	ClassID    int64  `json:"class_id"`
}

type GetCourseRequest struct {
	CourseID int64 `uri:"id" binding:"required"`
}

type GetCourseMarksResponse struct {
	CourseName string         `json:"course_name"`
	TeacherID  int64          `json:"teacher_id"`
	SemesterID int64          `json:"semester_id"`
	ClassID    int64          `json:"class_id"`
	Dates      []time.Time    `json:"dates"`
	Marks      []MarkResponse `json:"marks"`
}

func NewCourseResponse(course db.Course) CourseResponse {
	return CourseResponse{
		course.ID,
		course.Name,
		course.TeacherID,
		course.SemesterID,
		course.ClassID,
		course.CreatedBy.Int64,
		course.UpdatedBy.Int64,
		course.CreatedAt.Time,
		course.UpdatedAt.Time,
	}

}
