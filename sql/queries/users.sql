-- name: CreateUser :one
INSERT INTO users (id, created_at, updated_at, email, password_hash)
VALUES (
    gen_random_uuid(),
    NOW(),
    NOW(),
    $1,
    $2
)
RETURNING *;


-- name: ResetUsers :exec
TRUNCATE TABLE users RESTART IDENTITY CASCADE;


-- name: GetUserByEmail :one
SELECT * FROM users WHERE email = $1;


-- name: UpdateUser :one
UPDATE users SET email = $2, password_hash = $3
WHERE id = $1
RETURNING *;


-- name: GetUserByID :one
SELECT * FROM users WHERE id = $1;

-- name: UpgradeToChirpyRed :one
UPDATE users SET is_chirpy_red = true, updated_at = NOW()
WHERE id = $1
RETURNING *;