-- name: CreateUser :one
INSERT INTO "User"(
    email, password, last_name, first_name, gender, phone_number, domicile, birth_date
)VALUES (
            $1,$2,$3,$4,$5,$6,$7,$8
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