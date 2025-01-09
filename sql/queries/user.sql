-- name: CreateUser :one
INSERT INTO users (id, created_at, updated_at, name)
VALUES (
    $1,
    $2,
    $3,
    $4
)
RETURNING *;

-- name: GetUser :one
SELECT * FROM users WHERE name = $1 LIMIT 1;

-- name: GetIDByUsername :one
SELECT id FROM users WHERE name = $1 LIMIT 1;

-- name: GetUsernameByID :one
SELECT name FROM users WHERE id = $1 LIMIT 1;

-- name: UsernameExists :one
SELECT EXISTS (
    SELECT 2
    FROM users
    WHERE name = $1
) AS exists;
