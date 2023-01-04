// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.15.0
// source: lesson.sql

package db

import (
	"context"
	"database/sql"
	"time"

	"github.com/lib/pq"
)

const createLesson = `-- name: CreateLesson :one
INSERT INTO "Lesson"(
    name,course_id,start_hour,end_hour, week_day, classroom
)VALUES (
            $1,$2,$3,$4,$5,$6
        ) RETURNING id, name, course_id, teacher_id, start_hour, end_hour, week_day, classroom, created_by, updated_by, created_at, updated_at
`

type CreateLessonParams struct {
	Name      string         `json:"name"`
	CourseID  int64          `json:"courseID"`
	StartHour time.Time      `json:"startHour"`
	EndHour   time.Time      `json:"endHour"`
	WeekDay   string         `json:"weekDay"`
	Classroom sql.NullString `json:"classroom"`
}

func (q *Queries) CreateLesson(ctx context.Context, arg CreateLessonParams) (Lesson, error) {
	row := q.db.QueryRowContext(ctx, createLesson,
		arg.Name,
		arg.CourseID,
		arg.StartHour,
		arg.EndHour,
		arg.WeekDay,
		arg.Classroom,
	)
	var i Lesson
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.CourseID,
		&i.TeacherID,
		&i.StartHour,
		&i.EndHour,
		&i.WeekDay,
		&i.Classroom,
		&i.CreatedBy,
		&i.UpdatedBy,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getClassSchedule = `-- name: GetClassSchedule :many
SELECT "Lesson".id, "Lesson".name, course_id, "Lesson".teacher_id, start_hour, end_hour, week_day, classroom, "Lesson".created_by, "Lesson".updated_by, "Lesson".created_at, "Lesson".updated_at, "Course".id, "Course".name, "Course".teacher_id, semester_id, class_id, dates, "Course".created_by, "Course".updated_by, "Course".created_at, "Course".updated_at
FROM "Lesson"
INNER JOIN "Course"
ON  "Lesson".course_id = "Course".id AND "Course".class_id = $1
`

type GetClassScheduleRow struct {
	ID          int64          `json:"id"`
	Name        string         `json:"name"`
	CourseID    int64          `json:"courseID"`
	TeacherID   int64          `json:"teacherID"`
	StartHour   time.Time      `json:"startHour"`
	EndHour     time.Time      `json:"endHour"`
	WeekDay     string         `json:"weekDay"`
	Classroom   sql.NullString `json:"classroom"`
	CreatedBy   sql.NullInt64  `json:"createdBy"`
	UpdatedBy   sql.NullInt64  `json:"updatedBy"`
	CreatedAt   sql.NullTime   `json:"createdAt"`
	UpdatedAt   sql.NullTime   `json:"updatedAt"`
	ID_2        int64          `json:"id2"`
	Name_2      string         `json:"name2"`
	TeacherID_2 int64          `json:"teacherID2"`
	SemesterID  int64          `json:"semesterID"`
	ClassID     int64          `json:"classID"`
	Dates       []time.Time    `json:"dates"`
	CreatedBy_2 sql.NullInt64  `json:"createdBy2"`
	UpdatedBy_2 sql.NullInt64  `json:"updatedBy2"`
	CreatedAt_2 sql.NullTime   `json:"createdAt2"`
	UpdatedAt_2 sql.NullTime   `json:"updatedAt2"`
}

func (q *Queries) GetClassSchedule(ctx context.Context, classID int64) ([]GetClassScheduleRow, error) {
	rows, err := q.db.QueryContext(ctx, getClassSchedule, classID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []GetClassScheduleRow{}
	for rows.Next() {
		var i GetClassScheduleRow
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.CourseID,
			&i.TeacherID,
			&i.StartHour,
			&i.EndHour,
			&i.WeekDay,
			&i.Classroom,
			&i.CreatedBy,
			&i.UpdatedBy,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.ID_2,
			&i.Name_2,
			&i.TeacherID_2,
			&i.SemesterID,
			&i.ClassID,
			pq.Array(&i.Dates),
			&i.CreatedBy_2,
			&i.UpdatedBy_2,
			&i.CreatedAt_2,
			&i.UpdatedAt_2,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getLessonsOfCourse = `-- name: GetLessonsOfCourse :many
SELECT id, name, course_id, teacher_id, start_hour, end_hour, week_day, classroom, created_by, updated_by, created_at, updated_at FROM "Lesson"
WHERE course_id = $1
`

func (q *Queries) GetLessonsOfCourse(ctx context.Context, courseID int64) ([]Lesson, error) {
	rows, err := q.db.QueryContext(ctx, getLessonsOfCourse, courseID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Lesson{}
	for rows.Next() {
		var i Lesson
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.CourseID,
			&i.TeacherID,
			&i.StartHour,
			&i.EndHour,
			&i.WeekDay,
			&i.Classroom,
			&i.CreatedBy,
			&i.UpdatedBy,
			&i.CreatedAt,
			&i.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getTeacherSchedule = `-- name: GetTeacherSchedule :many
SELECT "Lesson".id, "Lesson".name, course_id, "Lesson".teacher_id, start_hour, end_hour, week_day, classroom, "Lesson".created_by, "Lesson".updated_by, "Lesson".created_at, "Lesson".updated_at, "Course".id, "Course".name, "Course".teacher_id, semester_id, class_id, dates, "Course".created_by, "Course".updated_by, "Course".created_at, "Course".updated_at
FROM "Lesson"
INNER JOIN "Course"
ON  "Lesson".course_id = "Course".id AND "Course".teacher_id = $1
`

type GetTeacherScheduleRow struct {
	ID          int64          `json:"id"`
	Name        string         `json:"name"`
	CourseID    int64          `json:"courseID"`
	TeacherID   int64          `json:"teacherID"`
	StartHour   time.Time      `json:"startHour"`
	EndHour     time.Time      `json:"endHour"`
	WeekDay     string         `json:"weekDay"`
	Classroom   sql.NullString `json:"classroom"`
	CreatedBy   sql.NullInt64  `json:"createdBy"`
	UpdatedBy   sql.NullInt64  `json:"updatedBy"`
	CreatedAt   sql.NullTime   `json:"createdAt"`
	UpdatedAt   sql.NullTime   `json:"updatedAt"`
	ID_2        int64          `json:"id2"`
	Name_2      string         `json:"name2"`
	TeacherID_2 int64          `json:"teacherID2"`
	SemesterID  int64          `json:"semesterID"`
	ClassID     int64          `json:"classID"`
	Dates       []time.Time    `json:"dates"`
	CreatedBy_2 sql.NullInt64  `json:"createdBy2"`
	UpdatedBy_2 sql.NullInt64  `json:"updatedBy2"`
	CreatedAt_2 sql.NullTime   `json:"createdAt2"`
	UpdatedAt_2 sql.NullTime   `json:"updatedAt2"`
}

func (q *Queries) GetTeacherSchedule(ctx context.Context, teacherID int64) ([]GetTeacherScheduleRow, error) {
	rows, err := q.db.QueryContext(ctx, getTeacherSchedule, teacherID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []GetTeacherScheduleRow{}
	for rows.Next() {
		var i GetTeacherScheduleRow
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.CourseID,
			&i.TeacherID,
			&i.StartHour,
			&i.EndHour,
			&i.WeekDay,
			&i.Classroom,
			&i.CreatedBy,
			&i.UpdatedBy,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.ID_2,
			&i.Name_2,
			&i.TeacherID_2,
			&i.SemesterID,
			&i.ClassID,
			pq.Array(&i.Dates),
			&i.CreatedBy_2,
			&i.UpdatedBy_2,
			&i.CreatedAt_2,
			&i.UpdatedAt_2,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updateLesson = `-- name: UpdateLesson :one
UPDATE  "Lesson"
SET  start_hour= $2,end_hour=$3, week_day=$4,classroom=$5,updated_at = now()
where id = $1
RETURNING id, name, course_id, teacher_id, start_hour, end_hour, week_day, classroom, created_by, updated_by, created_at, updated_at
`

type UpdateLessonParams struct {
	ID        int64          `json:"id"`
	StartHour time.Time      `json:"startHour"`
	EndHour   time.Time      `json:"endHour"`
	WeekDay   string         `json:"weekDay"`
	Classroom sql.NullString `json:"classroom"`
}

func (q *Queries) UpdateLesson(ctx context.Context, arg UpdateLessonParams) (Lesson, error) {
	row := q.db.QueryRowContext(ctx, updateLesson,
		arg.ID,
		arg.StartHour,
		arg.EndHour,
		arg.WeekDay,
		arg.Classroom,
	)
	var i Lesson
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.CourseID,
		&i.TeacherID,
		&i.StartHour,
		&i.EndHour,
		&i.WeekDay,
		&i.Classroom,
		&i.CreatedBy,
		&i.UpdatedBy,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}
