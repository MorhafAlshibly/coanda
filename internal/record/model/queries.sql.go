// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.24.0
// source: queries.sql

package model

import (
	"context"
	"database/sql"
	"encoding/json"
)

const createRecord = `-- name: CreateRecord :execresult
INSERT INTO record (name, user_id, record, data)
VALUES (?, ?, ?, ?)
`

type CreateRecordParams struct {
	Name   string
	UserID uint64
	Record uint64
	Data   json.RawMessage
}

func (q *Queries) CreateRecord(ctx context.Context, arg CreateRecordParams) (sql.Result, error) {
	return q.db.ExecContext(ctx, createRecord,
		arg.Name,
		arg.UserID,
		arg.Record,
		arg.Data,
	)
}

const deleteRecord = `-- name: DeleteRecord :execresult
DELETE FROM record
WHERE name = ?
  AND user_id = ?
LIMIT 1
`

type DeleteRecordParams struct {
	Name   string
	UserID uint64
}

func (q *Queries) DeleteRecord(ctx context.Context, arg DeleteRecordParams) (sql.Result, error) {
	return q.db.ExecContext(ctx, deleteRecord, arg.Name, arg.UserID)
}

const getRecord = `-- name: GetRecord :one
SELECT name,
  user_id,
  record,
  ranking,
  data,
  created_at,
  updated_at
FROM ranked_record
WHERE name = ?
  AND user_id = ?
LIMIT 1
`

type GetRecordParams struct {
	Name   string
	UserID uint64
}

func (q *Queries) GetRecord(ctx context.Context, arg GetRecordParams) (RankedRecord, error) {
	row := q.db.QueryRowContext(ctx, getRecord, arg.Name, arg.UserID)
	var i RankedRecord
	err := row.Scan(
		&i.Name,
		&i.UserID,
		&i.Record,
		&i.Ranking,
		&i.Data,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getRecords = `-- name: GetRecords :many
SELECT name,
  user_id,
  record,
  ranking,
  data,
  created_at,
  updated_at
FROM ranked_record
WHERE name = ?
  OR user_id = ?
ORDER BY record ASC
LIMIT ? OFFSET ?
`

type GetRecordsParams struct {
	Name   sql.NullString
	UserID sql.NullInt64
	Limit  int32
	Offset int32
}

func (q *Queries) GetRecords(ctx context.Context, arg GetRecordsParams) ([]RankedRecord, error) {
	rows, err := q.db.QueryContext(ctx, getRecords,
		arg.Name,
		arg.UserID,
		arg.Limit,
		arg.Offset,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []RankedRecord
	for rows.Next() {
		var i RankedRecord
		if err := rows.Scan(
			&i.Name,
			&i.UserID,
			&i.Record,
			&i.Ranking,
			&i.Data,
			&i.CreatedAt,
			&i.UpdatedAt,
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

const updateRecordData = `-- name: UpdateRecordData :execresult
UPDATE record
SET data = ?
WHERE name = ?
  AND user_id = ?
LIMIT 1
`

type UpdateRecordDataParams struct {
	Data   json.RawMessage
	Name   string
	UserID uint64
}

func (q *Queries) UpdateRecordData(ctx context.Context, arg UpdateRecordDataParams) (sql.Result, error) {
	return q.db.ExecContext(ctx, updateRecordData, arg.Data, arg.Name, arg.UserID)
}

const updateRecordRecord = `-- name: UpdateRecordRecord :execresult
UPDATE record
SET record = ?
WHERE name = ?
  AND user_id = ?
LIMIT 1
`

type UpdateRecordRecordParams struct {
	Record uint64
	Name   string
	UserID uint64
}

func (q *Queries) UpdateRecordRecord(ctx context.Context, arg UpdateRecordRecordParams) (sql.Result, error) {
	return q.db.ExecContext(ctx, updateRecordRecord, arg.Record, arg.Name, arg.UserID)
}
