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

-- name: GetRolebyId :one
SELECT * FROM "Role"
WHERE id = $1 LIMIT 1;

-- name: GetRoles :many
SELECT * FROM "Role";
