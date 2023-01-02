// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.15.0
// source: user_roles.sql

package db

import (
	"context"
	"database/sql"
)

const addUserToClass = `-- name: AddUserToClass :one
INSERT INTO "UserRoleClass"(
    user_role_id,class_id
)VALUES (
            $1,$2
        ) RETURNING id, user_role_id, class_id
`

type AddUserToClassParams struct {
	UserRoleID int64 `json:"userRoleID"`
	ClassID    int64 `json:"classID"`
}

func (q *Queries) AddUserToClass(ctx context.Context, arg AddUserToClassParams) (UserRoleClass, error) {
	row := q.db.QueryRowContext(ctx, addUserToClass, arg.UserRoleID, arg.ClassID)
	var i UserRoleClass
	err := row.Scan(&i.ID, &i.UserRoleID, &i.ClassID)
	return i, err
}

const createRoleForUser = `-- name: CreateRoleForUser :one
INSERT INTO "UserRoles"(
user_id,role_id,school_id
)VALUES (
$1,$2,$3
) RETURNING id, user_id, role_id, school_id
`

type CreateRoleForUserParams struct {
	UserID   int64         `json:"userID"`
	RoleID   int64         `json:"roleID"`
	SchoolID sql.NullInt64 `json:"schoolID"`
}

func (q *Queries) CreateRoleForUser(ctx context.Context, arg CreateRoleForUserParams) (UserRole, error) {
	row := q.db.QueryRowContext(ctx, createRoleForUser, arg.UserID, arg.RoleID, arg.SchoolID)
	var i UserRole
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.RoleID,
		&i.SchoolID,
	)
	return i, err
}

const getUserClassByUserRoleId = `-- name: GetUserClassByUserRoleId :one
SELECT id, user_role_id, class_id FROM "UserRoleClass"
WHERE user_role_id = $1 LIMIT 1
`

func (q *Queries) GetUserClassByUserRoleId(ctx context.Context, userRoleID int64) (UserRoleClass, error) {
	row := q.db.QueryRowContext(ctx, getUserClassByUserRoleId, userRoleID)
	var i UserRoleClass
	err := row.Scan(&i.ID, &i.UserRoleID, &i.ClassID)
	return i, err
}

const getUserRoleById = `-- name: GetUserRoleById :one
SELECT id, user_id, role_id, school_id FROM "UserRoles"
WHERE id = $1
`

func (q *Queries) GetUserRoleById(ctx context.Context, id int64) (UserRole, error) {
	row := q.db.QueryRowContext(ctx, getUserRoleById, id)
	var i UserRole
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.RoleID,
		&i.SchoolID,
	)
	return i, err
}

const getUserRoleByUserId = `-- name: GetUserRoleByUserId :many
SELECT id, user_id, role_id, school_id FROM "UserRoles"
WHERE user_id = $1
`

func (q *Queries) GetUserRoleByUserId(ctx context.Context, userID int64) ([]UserRole, error) {
	rows, err := q.db.QueryContext(ctx, getUserRoleByUserId, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []UserRole{}
	for rows.Next() {
		var i UserRole
		if err := rows.Scan(
			&i.ID,
			&i.UserID,
			&i.RoleID,
			&i.SchoolID,
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
