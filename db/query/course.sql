-- name: CreateCourse :one
INSERT INTO "Course"(
    name,teacher_id,semester_id,class_id
)VALUES (
            $1,$2,$3,$4
        ) RETURNING *;

-- name: UpdateCourseDates :one
UPDATE  "Course"
SET  dates = $2, updated_at = now()
where id = $1
RETURNING *;

-- name: ListCoursesOfClass :many
SELECT * FROM "Course"
WHERE class_id = $1
ORDER BY name
LIMIT $2
OFFSET $3;

-- name: ListCoursesOfTeacher :many
SELECT * FROM "Course"
WHERE teacher_id = $1
ORDER BY name
LIMIT $2
OFFSET $3;

-- name: UpdateCourseTeacher :one
UPDATE  "Course"
SET  teacher_id = $2, updated_at = now()
where id = $1
RETURNING *;