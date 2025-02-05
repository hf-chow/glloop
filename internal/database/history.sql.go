// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: history.sql

package database

import (
	"context"
	"time"

	"github.com/google/uuid"
)

const createHistory = `-- name: CreateHistory :one
INSERT INTO history (id, user_id, created_at, prompt, reply) 
VALUES (
    $1,
    $2,
    $3,
    $4,
    $5
)
RETURNING id, user_id, created_at, prompt, reply
`

type CreateHistoryParams struct {
	ID        uuid.UUID
	UserID    uuid.UUID
	CreatedAt time.Time
	Prompt    string
	Reply     string
}

func (q *Queries) CreateHistory(ctx context.Context, arg CreateHistoryParams) (History, error) {
	row := q.db.QueryRowContext(ctx, createHistory,
		arg.ID,
		arg.UserID,
		arg.CreatedAt,
		arg.Prompt,
		arg.Reply,
	)
	var i History
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.CreatedAt,
		&i.Prompt,
		&i.Reply,
	)
	return i, err
}

const deleteAllHistoryByUserID = `-- name: DeleteAllHistoryByUserID :many
DELETE FROM history WHERE user_id = $1
RETURNING id, user_id, created_at, prompt, reply
`

func (q *Queries) DeleteAllHistoryByUserID(ctx context.Context, userID uuid.UUID) ([]History, error) {
	rows, err := q.db.QueryContext(ctx, deleteAllHistoryByUserID, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []History
	for rows.Next() {
		var i History
		if err := rows.Scan(
			&i.ID,
			&i.UserID,
			&i.CreatedAt,
			&i.Prompt,
			&i.Reply,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getAllHistoryByUserID = `-- name: GetAllHistoryByUserID :many
SELECT id, user_id, created_at, prompt, reply FROM history WHERE user_id = $1
ORDER BY created_at DESC
`

func (q *Queries) GetAllHistoryByUserID(ctx context.Context, userID uuid.UUID) ([]History, error) {
	rows, err := q.db.QueryContext(ctx, getAllHistoryByUserID, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []History
	for rows.Next() {
		var i History
		if err := rows.Scan(
			&i.ID,
			&i.UserID,
			&i.CreatedAt,
			&i.Prompt,
			&i.Reply,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getLastHistoryByUserID = `-- name: GetLastHistoryByUserID :one
SELECT id, user_id, created_at, prompt, reply FROM history WHERE user_id = $1
ORDER BY created_at DESC
LIMIT 1
`

func (q *Queries) GetLastHistoryByUserID(ctx context.Context, userID uuid.UUID) (History, error) {
	row := q.db.QueryRowContext(ctx, getLastHistoryByUserID, userID)
	var i History
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.CreatedAt,
		&i.Prompt,
		&i.Reply,
	)
	return i, err
}
