// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: queries.sql

package model

import (
	"context"
	"database/sql"
	"encoding/json"
)

const CreateItem = `-- name: CreateItem :execresult
INSERT INTO item (
        id,
        type,
        data,
        expires_at
    )
VALUES (?, ?, ?, ?)
`

type CreateItemParams struct {
	ID        string          `db:"id"`
	Type      string          `db:"type"`
	Data      json.RawMessage `db:"data"`
	ExpiresAt sql.NullTime    `db:"expires_at"`
}

func (q *Queries) CreateItem(ctx context.Context, arg CreateItemParams) (sql.Result, error) {
	return q.db.ExecContext(ctx, CreateItem,
		arg.ID,
		arg.Type,
		arg.Data,
		arg.ExpiresAt,
	)
}

const DeleteItem = `-- name: DeleteItem :execresult
DELETE FROM item
WHERE id = ?
    AND type = ?
    AND (
        expires_at IS NULL
        OR expires_at > NOW()
    )
LIMIT 1
`

type DeleteItemParams struct {
	ID   string `db:"id"`
	Type string `db:"type"`
}

func (q *Queries) DeleteItem(ctx context.Context, arg DeleteItemParams) (sql.Result, error) {
	return q.db.ExecContext(ctx, DeleteItem, arg.ID, arg.Type)
}

const GetItem = `-- name: GetItem :one
SELECT id,
    type,
    data,
    expires_at,
    created_at,
    updated_at
FROM item
WHERE id = ?
    AND type = ?
    AND (
        expires_at IS NULL
        OR expires_at > NOW()
    )
LIMIT 1
`

type GetItemParams struct {
	ID   string `db:"id"`
	Type string `db:"type"`
}

func (q *Queries) GetItem(ctx context.Context, arg GetItemParams) (Item, error) {
	row := q.db.QueryRowContext(ctx, GetItem, arg.ID, arg.Type)
	var i Item
	err := row.Scan(
		&i.ID,
		&i.Type,
		&i.Data,
		&i.ExpiresAt,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const UpdateItem = `-- name: UpdateItem :execresult
UPDATE item
SET data = ?
WHERE id = ?
    AND type = ?
    AND (
        expires_at IS NULL
        OR expires_at > NOW()
    )
LIMIT 1
`

type UpdateItemParams struct {
	Data json.RawMessage `db:"data"`
	ID   string          `db:"id"`
	Type string          `db:"type"`
}

func (q *Queries) UpdateItem(ctx context.Context, arg UpdateItemParams) (sql.Result, error) {
	return q.db.ExecContext(ctx, UpdateItem, arg.Data, arg.ID, arg.Type)
}
