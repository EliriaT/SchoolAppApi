// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.15.0
// source: class.sql

package db

import (
	"context"
	"database/sql"
)

const createClass = `-- name: CreateClass :one
INSERT INTO "Class"(
    name, head_teacher
)VALUES (
            $1,$2
        ) RETURNING id, name, head_teacher, created_by, updated_by, created_at, updated_at
`

type CreateClassParams struct {
	Name        string `json:"name"`
	HeadTeacher int64  `json:"headTeacher"`
}

func (q *Queries) CreateClass(ctx context.Context, arg CreateClassParams) (Class, error) {
	row := q.db.QueryRowContext(ctx, createClass, arg.Name, arg.HeadTeacher)
	var i Class
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.HeadTeacher,
		&i.CreatedBy,
		&i.UpdatedBy,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getClassById = `-- name: GetClassById :one
SELECT id, name, head_teacher, created_by, updated_by, created_at, updated_at FROM "Class"
WHERE id = $1 LIMIT 1
`

func (q *Queries) GetClassById(ctx context.Context, id int64) (Class, error) {
	row := q.db.QueryRowContext(ctx, getClassById, id)
	var i Class
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.HeadTeacher,
		&i.CreatedBy,
		&i.UpdatedBy,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getClassWithStudents = `-- name: GetClassWithStudents :many
SELECT email, last_name, first_name, gender, phone_number,domicile,birth_date, role_id FROM "Class"
INNER JOIN "UserRoleClass"
ON  "Class".id = "UserRoleClass".class_id AND "Class".id = $1
INNER JOIN "UserRoles"
on "UserRoles".id = "UserRoleClass".user_role_id
INNER JOIN "User"
on "User".id = "UserRoles".user_id
`

type GetClassWithStudentsRow struct {
	Email       string         `json:"email"`
	LastName    string         `json:"lastName"`
	FirstName   string         `json:"firstName"`
	Gender      string         `json:"gender"`
	PhoneNumber sql.NullString `json:"phoneNumber"`
	Domicile    sql.NullString `json:"domicile"`
	BirthDate   sql.NullTime   `json:"birthDate"`
	RoleID      int64          `json:"roleID"`
}

func (q *Queries) GetClassWithStudents(ctx context.Context, id int64) ([]GetClassWithStudentsRow, error) {
	rows, err := q.db.QueryContext(ctx, getClassWithStudents, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []GetClassWithStudentsRow{}
	for rows.Next() {
		var i GetClassWithStudentsRow
		if err := rows.Scan(
			&i.Email,
			&i.LastName,
			&i.FirstName,
			&i.Gender,
			&i.PhoneNumber,
			&i.Domicile,
			&i.BirthDate,
			&i.RoleID,
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

const listAllClasses = `-- name: ListAllClasses :many
SELECT id, name, head_teacher, created_by, updated_by, created_at, updated_at FROM "Class"
ORDER BY name
`

func (q *Queries) ListAllClasses(ctx context.Context) ([]Class, error) {
	rows, err := q.db.QueryContext(ctx, listAllClasses)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Class{}
	for rows.Next() {
		var i Class
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.HeadTeacher,
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

const updateClassHeadTeacher = `-- name: UpdateClassHeadTeacher :one
UPDATE  "Class"
SET  head_teacher = $2, updated_at = now()
where id = $1
RETURNING id, name, head_teacher, created_by, updated_by, created_at, updated_at
`

type UpdateClassHeadTeacherParams struct {
	ID          int64 `json:"id"`
	HeadTeacher int64 `json:"headTeacher"`
}

func (q *Queries) UpdateClassHeadTeacher(ctx context.Context, arg UpdateClassHeadTeacherParams) (Class, error) {
	row := q.db.QueryRowContext(ctx, updateClassHeadTeacher, arg.ID, arg.HeadTeacher)
	var i Class
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.HeadTeacher,
		&i.CreatedBy,
		&i.UpdatedBy,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}
