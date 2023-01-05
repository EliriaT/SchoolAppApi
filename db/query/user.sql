-- name: CreateUser :one
INSERT INTO "User"(
    email, password, totp_secret ,last_name, first_name, gender, phone_number, domicile, birth_date
)VALUES (
            $1,$2,$3,$4,$5,$6,$7,$8, $9
        ) RETURNING *;

-- name: GetUserbyId :one
SELECT * FROM "User"
WHERE id = $1 LIMIT 1;


-- name: GetUserbyEmail :one
SELECT * FROM "User"
WHERE email = $1 LIMIT 1;

-- name: UpdateUserPassword :one
UPDATE  "User"
SET  password = $2, updated_at = now()
where id = $1
RETURNING *;

-- name: UpdateUserEmail :one
UPDATE  "User"
SET  email = $2, updated_at = now()
where id = $1
RETURNING *;

-- name: DeleteUser :exec
DELETE FROM "User"
WHERE id = $1;

-- name: GetTeachers :many
Select * from "User"
INNER JOIN "UserRoles"
ON  "User".id = "UserRoles".user_id and "UserRoles".school_id = $1
INNER JOIN "Role"
ON  "UserRoles".role_id = "Role".id and ("Role".name = 'Teacher' or "Role".name = 'Director' or "Role".name = 'School_Manager')
LEFT JOIN "UserRoleClass"
on "UserRoleClass".user_role_id = "UserRoles".id
WHERE "UserRoleClass".user_role_id is Null;