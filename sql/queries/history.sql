-- name: CreateHistory :one
INSERT INTO history (id, user_id, created_at, prompt, reply) 
VALUES (
    $1,
    $2,
    $3,
    $4,
    $5
)
RETURNING *;

-- name: GetLastHistoryByUserID :one
SELECT * FROM history WHERE user_id = $1
ORDER BY created_at DESC
LIMIT 1;

-- name: GetAllHistoryByUserID :many
SELECT * FROM history WHERE user_id = $1
ORDER BY created_at DESC;

-- name: DeleteAllHistoryByUserID :many
DELETE FROM history WHERE user_id = $1
RETURNING *;
