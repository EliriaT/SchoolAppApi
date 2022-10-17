-- name: CreateSchool :one
INSERT INTO "School"(
    name
)VALUES (
    $1
) RETURNING *;

-- name: GetSchoolbyId :one
SELECT * FROM "School"
WHERE id = $1 LIMIT 1;

-- name: GetSchoolbyName :one
SELECT * FROM "School"
WHERE name = $1 LIMIT 1;

-- name: ListSchools :many
SELECT * FROM "School"
ORDER BY name
LIMIT $1
OFFSET $2;

-- name: UpdateSchool :one
UPDATE  "School"
SET  name = $2
where id = $1
RETURNING *;

-- name: DeleteSchool :exec
DELETE FROM "School"
WHERE id = $1;