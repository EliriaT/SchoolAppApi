-- name: CreateSemester :one
INSERT INTO "Semester"(
    name,start_date,end_date
)VALUES (
            $1,$2,$3
        ) RETURNING *;

-- name: GetSemesterbyId :one
SELECT * FROM "Semester"
WHERE id = $1 LIMIT 1;

-- name: GetSemesters :one
SELECT * FROM "Semester"
ORDER BY start_date DESC
LIMIT $1
OFFSET $2;

-- name: GetCurrentSemester :one
SELECT * FROM "Semester"
WHERE NOW() BETWEEN start_date AND end_date
LIMIT 1;