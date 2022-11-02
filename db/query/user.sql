-- name: CreateUser :one
INSERT INTO "User"(
    email, password, last_name, first_name, gender, phone_number, domicile, birth_date
)VALUES (
            $1,$2,$3,$4,$5,$6,$7,$8
        ) RETURNING *;

-- name: GetUserbyId :one
SELECT * FROM "User"
WHERE id = $1 LIMIT 1;
