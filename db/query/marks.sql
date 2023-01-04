-- name: CreateMark :one
INSERT INTO "Marks"(
    course_id,mark_date,is_absent, mark, student_id
)VALUES (
            $1,$2,$3,$4,$5
        ) RETURNING *;

-- name: GetCourseMarks :many
SELECT * FROM "Marks"
INNER JOIN "Course"
ON  "Course".id = "Marks".course_id AND "Course".id = $1;

-- name: GetMarkByID :one
SELECT * FROM "Marks"
where id = $1;

-- name: GetStudentCourseMarks :many
SELECT * FROM "Marks"
INNER JOIN "Course"
ON  "Course".id = "Marks".course_id AND "Course".id = $1 AND "Marks".student_id = $2;

-- name: UpdateCourseMarksbyId :one
UPDATE  "Marks"
SET  mark = $2, is_absent = false,updated_at = now()
where id = $1
RETURNING *;

-- name: UpdateCourseAbsencebyId :one
UPDATE  "Marks"
SET  mark = 0, is_absent = true,updated_at = now()
where id = $1
RETURNING *;

-- name: DeleteMark :exec
DELETE FROM "Marks"
WHERE id = $1;