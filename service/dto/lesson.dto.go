package dto

import (
	db "github.com/EliriaT/SchoolAppApi/db/sqlc"
	"time"
)

type CreateLessonRequest struct {
	Name      string    `json:"name" binding:"required"`
	CourseID  int64     `json:"course_id" binding:"required"`
	StartHour time.Time `json:"start_hour" binding:"required ltefield=EndHour" time_format:"15:04"`
	EndHour   time.Time `json:"end_hour" binding:"required" time_format:"15:04"`
	WeekDay   string    `json:"week_day" binding:"required,oneof=Monday Tuesday Wednesday Thursday Friday Saturday Sunday"`
	Classroom string    `json:"classroom"`
}
type LessonResponse struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	CourseID  int64     `json:"course_id"`
	StartHour time.Time `json:"start_hour"`
	EndHour   time.Time `json:"end_hour"`
	WeekDay   string    `json:"week_day"`
	Classroom string    `json:"classroom"`
	CreatedBy int64     `json:"createdBy"`
	UpdatedBy int64     `json:"updatedBy"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	TeacherID int64     `json:"teacher_id,omitempty"`
	ClassID   int64     `json:"class_id,omitempty"`
}

type UpdateLessonParams struct {
	LessonID  int64     `json:"lesson_id"`
	Name      string    `json:"name"`
	CourseID  int64     `json:"course_id"`
	StartHour time.Time `json:"start_hour" binding:"required ltefield=EndHour" time_format:"15:04"`
	EndHour   time.Time `json:"end_hour" binding:"required" time_format:"15:04"`
	WeekDay   string    `json:"week_day"`
	Classroom string    `json:"classroom"`
}

type GetCourseLessonsRequest struct {
	CourseID int64 `uri:"id" binding:"required"`
}

func NewLessonsResponse(lesson db.Lesson) LessonResponse {
	return LessonResponse{
		ID:        lesson.ID,
		Name:      lesson.Name,
		CourseID:  lesson.CourseID,
		StartHour: lesson.StartHour,
		EndHour:   lesson.EndHour,
		WeekDay:   lesson.WeekDay,
		Classroom: lesson.Classroom.String,
		CreatedBy: lesson.CreatedBy.Int64,
		UpdatedBy: lesson.UpdatedBy.Int64,
		CreatedAt: lesson.CreatedAt.Time,
		UpdatedAt: lesson.UpdatedAt.Time,
	}

}
