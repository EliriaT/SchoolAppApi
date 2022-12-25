-- name: CreateClass :one
INSERT INTO "Class"(
    name, head_teacher
)VALUES (
            $1,$2
        ) RETURNING *;

-- name: GetClassById :one
SELECT * FROM "Class"
WHERE id = $1 LIMIT 1;

-- name: ListAllClasses :many
SELECT * FROM "Class"
ORDER BY name
LIMIT $1
OFFSET $2;

-- name: UpdateClassHeadTeacher :one
UPDATE  "Class"
SET  head_teacher = $2, updated_at = now()
where id = $1
RETURNING *;