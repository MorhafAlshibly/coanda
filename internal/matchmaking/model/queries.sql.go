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

const AddTicketIDToUser = `-- name: AddTicketIDToUser :execresult
UPDATE matchmaking_user
SET matchmaking_ticket_id = ?
WHERE id = ?
    AND matchmaking_ticket_id IS NULL
`

type AddTicketIDToUserParams struct {
	MatchmakingTicketID sql.NullInt64 `db:"matchmaking_ticket_id"`
	ID                  uint64        `db:"id"`
}

func (q *Queries) AddTicketIDToUser(ctx context.Context, arg AddTicketIDToUserParams) (sql.Result, error) {
	return q.db.ExecContext(ctx, AddTicketIDToUser, arg.MatchmakingTicketID, arg.ID)
}

const CreateArena = `-- name: CreateArena :execresult
INSERT INTO matchmaking_arena (
        name,
        min_players,
        max_players_per_ticket,
        max_players,
        data
    )
VALUES (?, ?, ?, ?, ?)
`

type CreateArenaParams struct {
	Name                string          `db:"name"`
	MinPlayers          uint32          `db:"min_players"`
	MaxPlayersPerTicket uint32          `db:"max_players_per_ticket"`
	MaxPlayers          uint32          `db:"max_players"`
	Data                json.RawMessage `db:"data"`
}

func (q *Queries) CreateArena(ctx context.Context, arg CreateArenaParams) (sql.Result, error) {
	return q.db.ExecContext(ctx, CreateArena,
		arg.Name,
		arg.MinPlayers,
		arg.MaxPlayersPerTicket,
		arg.MaxPlayers,
		arg.Data,
	)
}

const CreateMatchmakingTicket = `-- name: CreateMatchmakingTicket :execresult
INSERT INTO matchmaking_ticket (data)
VALUES (?)
`

func (q *Queries) CreateMatchmakingTicket(ctx context.Context, data json.RawMessage) (sql.Result, error) {
	return q.db.ExecContext(ctx, CreateMatchmakingTicket, data)
}

const CreateMatchmakingTicketArena = `-- name: CreateMatchmakingTicketArena :execresult
INSERT INTO matchmaking_ticket_arena (matchmaking_ticket_id, matchmaking_arena_id)
VALUES (?, ?)
`

type CreateMatchmakingTicketArenaParams struct {
	MatchmakingTicketID uint64 `db:"matchmaking_ticket_id"`
	MatchmakingArenaID  uint64 `db:"matchmaking_arena_id"`
}

func (q *Queries) CreateMatchmakingTicketArena(ctx context.Context, arg CreateMatchmakingTicketArenaParams) (sql.Result, error) {
	return q.db.ExecContext(ctx, CreateMatchmakingTicketArena, arg.MatchmakingTicketID, arg.MatchmakingArenaID)
}

const CreateMatchmakingUser = `-- name: CreateMatchmakingUser :execresult
INSERT INTO matchmaking_user (client_user_id, elo, data)
VALUES (?, ?, ?)
`

type CreateMatchmakingUserParams struct {
	ClientUserID uint64          `db:"client_user_id"`
	Elo          int64           `db:"elo"`
	Data         json.RawMessage `db:"data"`
}

func (q *Queries) CreateMatchmakingUser(ctx context.Context, arg CreateMatchmakingUserParams) (sql.Result, error) {
	return q.db.ExecContext(ctx, CreateMatchmakingUser, arg.ClientUserID, arg.Elo, arg.Data)
}

const GetArenas = `-- name: GetArenas :many
SELECT id,
    name,
    min_players,
    max_players_per_ticket,
    max_players,
    data,
    created_at,
    updated_at
FROM matchmaking_arena
ORDER BY created_at DESC
LIMIT ? OFFSET ?
`

type GetArenasParams struct {
	Limit  int32 `db:"limit"`
	Offset int32 `db:"offset"`
}

func (q *Queries) GetArenas(ctx context.Context, arg GetArenasParams) ([]MatchmakingArena, error) {
	rows, err := q.db.QueryContext(ctx, GetArenas, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []MatchmakingArena
	for rows.Next() {
		var i MatchmakingArena
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.MinPlayers,
			&i.MaxPlayersPerTicket,
			&i.MaxPlayers,
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

const GetMatchmakingUsers = `-- name: GetMatchmakingUsers :many
SELECT id,
    matchmaking_ticket_id,
    client_user_id,
    elo,
    data,
    created_at,
    updated_at
FROM matchmaking_user
ORDER BY client_user_id ASC
LIMIT ? OFFSET ?
`

type GetMatchmakingUsersParams struct {
	Limit  int32 `db:"limit"`
	Offset int32 `db:"offset"`
}

func (q *Queries) GetMatchmakingUsers(ctx context.Context, arg GetMatchmakingUsersParams) ([]MatchmakingUser, error) {
	rows, err := q.db.QueryContext(ctx, GetMatchmakingUsers, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []MatchmakingUser
	for rows.Next() {
		var i MatchmakingUser
		if err := rows.Scan(
			&i.ID,
			&i.MatchmakingTicketID,
			&i.ClientUserID,
			&i.Elo,
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
