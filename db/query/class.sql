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
ORDER BY name;

-- name: UpdateClassHeadTeacher :one
UPDATE  "Class"
SET  head_teacher = $2, updated_at = now()
where id = $1
RETURNING *;

-- name: GetClassWithStudents :many
SELECT email, last_name, first_name, gender, phone_number,domicile,birth_date, role_id FROM "Class"
INNER JOIN "UserRoleClass"
ON  "Class".id = "UserRoleClass".class_id AND "Class".id = $1
INNER JOIN "UserRoles"
on "UserRoles".id = "UserRoleClass".user_role_id
INNER JOIN "User"
on "User".id = "UserRoles".user_id;