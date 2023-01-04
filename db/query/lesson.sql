-- name: CreateLesson :one
INSERT INTO "Lesson"(
    name,course_id,start_hour,end_hour, week_day, classroom
)VALUES (
            $1,$2,$3,$4,$5,$6
        ) RETURNING *;

-- name: GetLessonsOfCourse :many
SELECT * FROM "Lesson"
WHERE course_id = $1;

-- name: GetClassSchedule :many
SELECT *
FROM "Lesson"
INNER JOIN "Course"
ON  "Lesson".course_id = "Course".id AND "Course".class_id = $1;

-- name: GetTeacherSchedule :many
SELECT *
FROM "Lesson"
INNER JOIN "Course"
ON  "Lesson".course_id = "Course".id AND "Course".teacher_id = $1;

-- name: UpdateLesson :one
UPDATE  "Lesson"
SET  start_hour= $2,end_hour=$3, week_day=$4,classroom=$5,updated_at = now()
where id = $1
RETURNING *;
