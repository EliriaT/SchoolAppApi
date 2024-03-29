// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.15.0
// source: marks.sql

package db

import (
	"context"
	"database/sql"
	"time"
)

const createMark = `-- name: CreateMark :one
INSERT INTO "Marks"(
    course_id,mark_date,is_absent, mark, student_id
)VALUES (
            $1,$2,$3,$4,$5
        ) RETURNING id, course_id, mark_date, is_absent, mark, student_id, created_by, updated_by, created_at, updated_at
`

type CreateMarkParams struct {
	CourseID  int64         `json:"courseID"`
	MarkDate  time.Time     `json:"markDate"`
	IsAbsent  sql.NullBool  `json:"isAbsent"`
	Mark      sql.NullInt32 `json:"mark"`
	StudentID int64         `json:"studentID"`
}

func (q *Queries) CreateMark(ctx context.Context, arg CreateMarkParams) (Mark, error) {
	row := q.db.QueryRowContext(ctx, createMark,
		arg.CourseID,
		arg.MarkDate,
		arg.IsAbsent,
		arg.Mark,
		arg.StudentID,
	)
	var i Mark
	err := row.Scan(
		&i.ID,
		&i.CourseID,
		&i.MarkDate,
		&i.IsAbsent,
		&i.Mark,
		&i.StudentID,
		&i.CreatedBy,
		&i.UpdatedBy,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const deleteMark = `-- name: DeleteMark :exec
DELETE FROM "Marks"
WHERE id = $1
`

func (q *Queries) DeleteMark(ctx context.Context, id int64) error {
	_, err := q.db.ExecContext(ctx, deleteMark, id)
	return err
}

const getMarkByID = `-- name: GetMarkByID :one
SELECT id, course_id, mark_date, is_absent, mark, student_id, created_by, updated_by, created_at, updated_at FROM "Marks"
where id = $1
`

func (q *Queries) GetMarkByID(ctx context.Context, id int64) (Mark, error) {
	row := q.db.QueryRowContext(ctx, getMarkByID, id)
	var i Mark
	err := row.Scan(
		&i.ID,
		&i.CourseID,
		&i.MarkDate,
		&i.IsAbsent,
		&i.Mark,
		&i.StudentID,
		&i.CreatedBy,
		&i.UpdatedBy,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const updateCourseAbsencebyId = `-- name: UpdateCourseAbsencebyId :one
UPDATE  "Marks"
SET  mark = 0, is_absent = true,updated_at = now()
where id = $1
RETURNING id, course_id, mark_date, is_absent, mark, student_id, created_by, updated_by, created_at, updated_at
`

func (q *Queries) UpdateCourseAbsencebyId(ctx context.Context, id int64) (Mark, error) {
	row := q.db.QueryRowContext(ctx, updateCourseAbsencebyId, id)
	var i Mark
	err := row.Scan(
		&i.ID,
		&i.CourseID,
		&i.MarkDate,
		&i.IsAbsent,
		&i.Mark,
		&i.StudentID,
		&i.CreatedBy,
		&i.UpdatedBy,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const updateCourseMarksbyId = `-- name: UpdateCourseMarksbyId :one
UPDATE  "Marks"
SET  mark = $2, is_absent = false,updated_at = now()
where id = $1
RETURNING id, course_id, mark_date, is_absent, mark, student_id, created_by, updated_by, created_at, updated_at
`

type UpdateCourseMarksbyIdParams struct {
	ID   int64         `json:"id"`
	Mark sql.NullInt32 `json:"mark"`
}

func (q *Queries) UpdateCourseMarksbyId(ctx context.Context, arg UpdateCourseMarksbyIdParams) (Mark, error) {
	row := q.db.QueryRowContext(ctx, updateCourseMarksbyId, arg.ID, arg.Mark)
	var i Mark
	err := row.Scan(
		&i.ID,
		&i.CourseID,
		&i.MarkDate,
		&i.IsAbsent,
		&i.Mark,
		&i.StudentID,
		&i.CreatedBy,
		&i.UpdatedBy,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}
