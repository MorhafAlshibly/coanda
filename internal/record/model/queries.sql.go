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

const CreateRecord = `-- name: CreateRecord :execresult
INSERT INTO record (name, user_id, record, data)
VALUES (?, ?, ?, ?)
`

type CreateRecordParams struct {
	Name   string          `db:"name"`
	UserID uint64          `db:"user_id"`
	Record uint64          `db:"record"`
	Data   json.RawMessage `db:"data"`
}

// -- name: GetRecord :one
// SELECT id,
//   name,
//   user_id,
//   record,
//   ranking,
//   data,
//   created_at,
//   updated_at
// FROM ranked_record
// WHERE (
//     name = sqlc.narg(name)
//     AND user_id = sqlc.narg(user_id)
//   )
//   OR id = sqlc.narg(id)
// LIMIT 1;
// -- name: GetRecords :many
// SELECT id,
//   name,
//   user_id,
//   record,
//   ranking,
//   data,
//   created_at,
//   updated_at
// FROM ranked_record
// WHERE name = sqlc.narg(name)
//   OR user_id = sqlc.narg(user_id)
// ORDER BY record ASC
// LIMIT ? OFFSET ?;
func (q *Queries) CreateRecord(ctx context.Context, arg CreateRecordParams) (sql.Result, error) {
	return q.db.ExecContext(ctx, CreateRecord,
		arg.Name,
		arg.UserID,
		arg.Record,
		arg.Data,
	)
}
