-- name: CreateRoleForUser :one
INSERT INTO "UserRoles"(
user_id,role_id,school_id
)VALUES (
$1,$2,$3
) RETURNING *;

-- name: GetUserRoleByUserId :one
SELECT * FROM "UserRoles"
WHERE user_id = $1 LIMIT 1;

-- name: AddUserToClass :one
INSERT INTO "UserRoleClass"(
    user_role_id,class_id
)VALUES (
            $1,$2
        ) RETURNING *;

-- name: GetUserClassByUserRoleId :one
SELECT * FROM "UserRoleClass"
WHERE user_role_id = $1 LIMIT 1;