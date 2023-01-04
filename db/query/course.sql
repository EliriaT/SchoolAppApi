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
ORDER BY name;

-- name: GetCourseByID :one
SELECT * FROM "Course"
WHERE id = $1;

-- name: ListCoursesOfTeacher :many
SELECT * FROM "Course"
WHERE teacher_id = $1
ORDER BY name;

-- name: GetCoursesOfSchool :many
SELECT * FROM "Course"
INNER JOIN "UserRoles"
ON  "Course".teacher_id = "UserRoles".id AND "UserRoles".school_id = $1
ORDER BY name;

-- name: UpdateCourse :one
UPDATE  "Course"
SET  teacher_id = $2, name = $1,semester_id=$3,class_id = $4,updated_at = now()
where id = $1
RETURNING *;