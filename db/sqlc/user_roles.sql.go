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
	UserRoleID sql.NullInt64 `json:"userRoleID"`
	ClassID    sql.NullInt64 `json:"classID"`
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
	UserID   sql.NullInt64 `json:"userID"`
	RoleID   sql.NullInt64 `json:"roleID"`
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

func (q *Queries) GetUserClassByUserRoleId(ctx context.Context, userRoleID sql.NullInt64) (UserRoleClass, error) {
	row := q.db.QueryRowContext(ctx, getUserClassByUserRoleId, userRoleID)
	var i UserRoleClass
	err := row.Scan(&i.ID, &i.UserRoleID, &i.ClassID)
	return i, err
}

const getUserRoleByUserId = `-- name: GetUserRoleByUserId :one
SELECT id, user_id, role_id, school_id FROM "UserRoles"
WHERE user_id = $1 LIMIT 1
`

func (q *Queries) GetUserRoleByUserId(ctx context.Context, userID sql.NullInt64) (UserRole, error) {
	row := q.db.QueryRowContext(ctx, getUserRoleByUserId, userID)
	var i UserRole
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.RoleID,
		&i.SchoolID,
	)
	return i, err
}
