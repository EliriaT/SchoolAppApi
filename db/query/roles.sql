-- name: CreateRoles :many
INSERT INTO "Role"(
    name
)VALUES (
            $1
        ),
        (
            $2
        ),
        (
            $3
        ),
        (
            $4
        ),
        (
            $5
        )
 RETURNING *;

-- name: GetRolebyName :one
SELECT * FROM "Role"
WHERE name = $1 LIMIT 1;

-- name: GetRoles :many
SELECT * FROM "Role";
