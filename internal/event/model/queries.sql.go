// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.24.0
// source: queries.sql

package model

import (
	"context"
	"database/sql"
	"encoding/json"
	"time"
)

const CreateEvent = `-- name: CreateEvent :execresult
INSERT INTO event (name, data, started_at)
VALUES (?, ?, ?)
`

type CreateEventParams struct {
	Name      string          `db:"name"`
	Data      json.RawMessage `db:"data"`
	StartedAt time.Time       `db:"started_at"`
}

func (q *Queries) CreateEvent(ctx context.Context, arg CreateEventParams) (sql.Result, error) {
	return q.db.ExecContext(ctx, CreateEvent, arg.Name, arg.Data, arg.StartedAt)
}

const CreateEventRound = `-- name: CreateEventRound :execresult
INSERT INTO event_round (event_id, name, data, scoring, ended_at)
VALUES (?, ?, ?, ?, ?)
`

type CreateEventRoundParams struct {
	EventID uint64          `db:"event_id"`
	Name    string          `db:"name"`
	Data    json.RawMessage `db:"data"`
	Scoring json.RawMessage `db:"scoring"`
	EndedAt time.Time       `db:"ended_at"`
}

func (q *Queries) CreateEventRound(ctx context.Context, arg CreateEventRoundParams) (sql.Result, error) {
	return q.db.ExecContext(ctx, CreateEventRound,
		arg.EventID,
		arg.Name,
		arg.Data,
		arg.Scoring,
		arg.EndedAt,
	)
}

const CreateEventRoundUser = `-- name: CreateEventRoundUser :execresult
INSERT INTO event_round_user (event_user_id, event_round_id, result, data)
VALUES (?, ?, ?, ?)
`

type CreateEventRoundUserParams struct {
	EventUserID  uint64          `db:"event_user_id"`
	EventRoundID uint64          `db:"event_round_id"`
	Result       uint64          `db:"result"`
	Data         json.RawMessage `db:"data"`
}

func (q *Queries) CreateEventRoundUser(ctx context.Context, arg CreateEventRoundUserParams) (sql.Result, error) {
	return q.db.ExecContext(ctx, CreateEventRoundUser,
		arg.EventUserID,
		arg.EventRoundID,
		arg.Result,
		arg.Data,
	)
}

const CreateOrUpdateEventUser = `-- name: CreateOrUpdateEventUser :execresult
INSERT INTO event_user (event_id, user_id, data)
VALUES (?, ?, ?) ON DUPLICATE KEY
UPDATE id = LAST_INSERT_ID(id),
    data = ?
`

type CreateOrUpdateEventUserParams struct {
	EventID uint64          `db:"event_id"`
	UserID  uint64          `db:"user_id"`
	Data    json.RawMessage `db:"data"`
}

func (q *Queries) CreateOrUpdateEventUser(ctx context.Context, arg CreateOrUpdateEventUserParams) (sql.Result, error) {
	return q.db.ExecContext(ctx, CreateOrUpdateEventUser,
		arg.EventID,
		arg.UserID,
		arg.Data,
		arg.Data,
	)
}

const DeleteEventRoundUser = `-- name: DeleteEventRoundUser :execresult
DELETE FROM event_round_user
WHERE id = ?
LIMIT 1
`

func (q *Queries) DeleteEventRoundUser(ctx context.Context, id uint64) (sql.Result, error) {
	return q.db.ExecContext(ctx, DeleteEventRoundUser, id)
}

const GetEventRoundUserByEventUserId = `-- name: GetEventRoundUserByEventUserId :one
SELECT eru.event_user_id,
    eru.event_round_id,
    eru.result,
    eru.data,
    eru.created_at,
    eru.updated_at
FROM event_round_user eru
    JOIN event_round er ON eru.event_round_id = er.id
WHERE eru.event_user_id = ?
    AND er.ended_at > NOW()
ORDER BY er.ended_at ASC
LIMIT 1
`

type GetEventRoundUserByEventUserIdRow struct {
	EventUserID  uint64          `db:"event_user_id"`
	EventRoundID uint64          `db:"event_round_id"`
	Result       uint64          `db:"result"`
	Data         json.RawMessage `db:"data"`
	CreatedAt    time.Time       `db:"created_at"`
	UpdatedAt    time.Time       `db:"updated_at"`
}

func (q *Queries) GetEventRoundUserByEventUserId(ctx context.Context, eventUserID uint64) (GetEventRoundUserByEventUserIdRow, error) {
	row := q.db.QueryRowContext(ctx, GetEventRoundUserByEventUserId, eventUserID)
	var i GetEventRoundUserByEventUserIdRow
	err := row.Scan(
		&i.EventUserID,
		&i.EventRoundID,
		&i.Result,
		&i.Data,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const UpdateEventRoundUserResult = `-- name: UpdateEventRoundUserResult :execresult
UPDATE event_round_user eru
SET eru.result = ?,
    eru.data = ?
WHERE eru.event_user_id = ?
    AND eru.event_round_id = ?
LIMIT 1
`

type UpdateEventRoundUserResultParams struct {
	Result       uint64          `db:"result"`
	Data         json.RawMessage `db:"data"`
	EventUserID  uint64          `db:"event_user_id"`
	EventRoundID uint64          `db:"event_round_id"`
}

func (q *Queries) UpdateEventRoundUserResult(ctx context.Context, arg UpdateEventRoundUserResultParams) (sql.Result, error) {
	return q.db.ExecContext(ctx, UpdateEventRoundUserResult,
		arg.Result,
		arg.Data,
		arg.EventUserID,
		arg.EventRoundID,
	)
}
