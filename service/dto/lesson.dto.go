package dto

import (
	"encoding/json"
	db "github.com/EliriaT/SchoolAppApi/db/sqlc"
	"time"
)

type CreateLessonRequest struct {
	Name      string    `json:"name" binding:"required"`
	CourseID  int64     `json:"course_id" binding:"required"`
	StartHour time.Time `json:"start_hour" binding:"required,ltefield=EndHour" time_format:"15:04"`
	EndHour   time.Time `json:"end_hour" binding:"required" time_format:"15:04"`
	WeekDay   string    `json:"week_day" binding:"required,oneof=Monday Tuesday Wednesday Thursday Friday Saturday Sunday"`
	Classroom string    `json:"classroom"`
}

func (st *CreateLessonRequest) UnmarshalJSON(data []byte) error {
	type parseType struct {
		Name      string `json:"name" binding:"required"`
		CourseID  int64  `json:"course_id" binding:"required"`
		StartHour string `json:"start_hour" binding:"required`
		EndHour   string `json:"end_hour" binding:"required"`
		WeekDay   string `json:"week_day" binding:"required,oneof=Monday Tuesday Wednesday Thursday Friday Saturday Sunday"`
		Classroom string `json:"classroom"`
	}
	var res parseType
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	parsedStartHour, err := time.Parse("15:04", res.StartHour)
	if err != nil {
		return err
	}
	parsedEndHour, err := time.Parse("15:04", res.EndHour)
	if err != nil {
		return err
	}

	now := time.Now()
	st.Name = res.Name
	st.CourseID = res.CourseID
	st.Classroom = res.Classroom
	st.WeekDay = res.WeekDay
	st.StartHour = time.Date(now.Year(), now.Month(), now.Day(), parsedStartHour.Hour(), parsedStartHour.Minute(), parsedStartHour.Second(), 0, now.Location())
	st.EndHour = time.Date(now.Year(), now.Month(), now.Day(), parsedEndHour.Hour(), parsedEndHour.Minute(), parsedEndHour.Second(), 0, now.Location())
	return nil
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
