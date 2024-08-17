// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0
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
    created_at,
    updated_at,
    expires_at
FROM item
WHERE id = ?
    AND type = ?
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
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.ExpiresAt,
	)
	return i, err
}

const GetItems = `-- name: GetItems :many
SELECT id,
    type,
    data,
    created_at,
    updated_at,
    expires_at
FROM item
WHERE type = CASE
        WHEN ? IS NOT NULL THEN ?
        ELSE type
    END
ORDER BY id ASC
LIMIT ? OFFSET ?
`

type GetItemsParams struct {
	Type   sql.NullString `db:"type"`
	Limit  int32          `db:"limit"`
	Offset int32          `db:"offset"`
}

func (q *Queries) GetItems(ctx context.Context, arg GetItemsParams) ([]Item, error) {
	rows, err := q.db.QueryContext(ctx, GetItems,
		arg.Type,
		arg.Type,
		arg.Limit,
		arg.Offset,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Item
	for rows.Next() {
		var i Item
		if err := rows.Scan(
			&i.ID,
			&i.Type,
			&i.Data,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.ExpiresAt,
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

const UpdateItem = `-- name: UpdateItem :execresult
UPDATE item
SET data = CASE
        WHEN CAST(? as unsigned) != 0 THEN ?
        ELSE data
    END,
    expires_at = CASE
        WHEN ? IS NOT NULL THEN ?
        ELSE expires_at
    END
WHERE id = ?
    AND type = ?
LIMIT 1
`

type UpdateItemParams struct {
	DataExists int64           `db:"data_exists"`
	Data       json.RawMessage `db:"data"`
	ExpiresAt  sql.NullTime    `db:"expires_at"`
	ID         string          `db:"id"`
	Type       string          `db:"type"`
}

func (q *Queries) UpdateItem(ctx context.Context, arg UpdateItemParams) (sql.Result, error) {
	return q.db.ExecContext(ctx, UpdateItem,
		arg.DataExists,
		arg.Data,
		arg.ExpiresAt,
		arg.ExpiresAt,
		arg.ID,
		arg.Type,
	)
}
