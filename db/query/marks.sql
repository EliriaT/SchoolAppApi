-- name: CreateMark :one
INSERT INTO "Marks"(
    course_id,mark_date,is_absent, mark, student_id
)VALUES (
            $1,$2,$3,$4,$5
        ) RETURNING *;

-- name: GetCourseMarks :one
SELECT * FROM "Marks"
WHERE course_id = $1;


-- name: UpdateCourseMarksbyId :one
UPDATE  "Marks"
SET  mark = $2, updated_at = now()
where id = $1
RETURNING *;

