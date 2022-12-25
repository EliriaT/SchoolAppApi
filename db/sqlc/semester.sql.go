// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.15.0
// source: semester.sql

package db

import (
	"context"
	"database/sql"
)

const createSemester = `-- name: CreateSemester :one
INSERT INTO "Semester"(
    name,start_date,end_date
)VALUES (
            $1,$2,$3
        ) RETURNING id, name, start_date, end_date, created_by, updated_by, created_at, updated_at
`

type CreateSemesterParams struct {
	Name      sql.NullString `json:"name"`
	StartDate sql.NullTime   `json:"startDate"`
	EndDate   sql.NullTime   `json:"endDate"`
}

func (q *Queries) CreateSemester(ctx context.Context, arg CreateSemesterParams) (Semester, error) {
	row := q.db.QueryRowContext(ctx, createSemester, arg.Name, arg.StartDate, arg.EndDate)
	var i Semester
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.StartDate,
		&i.EndDate,
		&i.CreatedBy,
		&i.UpdatedBy,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getCurrentSemester = `-- name: GetCurrentSemester :one
SELECT id, name, start_date, end_date, created_by, updated_by, created_at, updated_at FROM "Semester"
WHERE NOW() BETWEEN start_date AND end_date
LIMIT 1
`

func (q *Queries) GetCurrentSemester(ctx context.Context) (Semester, error) {
	row := q.db.QueryRowContext(ctx, getCurrentSemester)
	var i Semester
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.StartDate,
		&i.EndDate,
		&i.CreatedBy,
		&i.UpdatedBy,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getSemesterbyId = `-- name: GetSemesterbyId :one
SELECT id, name, start_date, end_date, created_by, updated_by, created_at, updated_at FROM "Semester"
WHERE id = $1 LIMIT 1
`

func (q *Queries) GetSemesterbyId(ctx context.Context, id int64) (Semester, error) {
	row := q.db.QueryRowContext(ctx, getSemesterbyId, id)
	var i Semester
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.StartDate,
		&i.EndDate,
		&i.CreatedBy,
		&i.UpdatedBy,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getSemesters = `-- name: GetSemesters :one
SELECT id, name, start_date, end_date, created_by, updated_by, created_at, updated_at FROM "Semester"
ORDER BY start_date DESC
LIMIT $1
OFFSET $2
`

type GetSemestersParams struct {
	Limit  int32 `json:"limit"`
	Offset int32 `json:"offset"`
}

func (q *Queries) GetSemesters(ctx context.Context, arg GetSemestersParams) (Semester, error) {
	row := q.db.QueryRowContext(ctx, getSemesters, arg.Limit, arg.Offset)
	var i Semester
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.StartDate,
		&i.EndDate,
		&i.CreatedBy,
		&i.UpdatedBy,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}
