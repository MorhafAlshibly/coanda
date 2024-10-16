// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: queries.sql

package model

import (
	"context"
	"database/sql"
	"encoding/json"
	"time"
)

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
	MinPlayers          uint8           `db:"min_players"`
	MaxPlayersPerTicket uint8           `db:"max_players_per_ticket"`
	MaxPlayers          uint8           `db:"max_players"`
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
INSERT INTO matchmaking_ticket (data, elo_window, expires_at)
SELECT ?,
    0,
    ?
FROM DUAL
WHERE NOT EXISTS (
        SELECT 1
        FROM matchmaking_ticket_user mtu
            JOIN matchmaking_ticket mt ON mtu.matchmaking_ticket_id = mt.id
            LEFT JOIN matchmaking_match mm ON mt.matchmaking_match_id = mm.id
        WHERE FIND_IN_SET(
                mtu.matchmaking_user_id,
                ?
            )
            AND (
                (
                    mt.matchmaking_match_id IS NULL
                    AND mt.expires_at > NOW()
                )
                OR (
                    mt.matchmaking_match_id IS NOT NULL
                    AND mm.ended_at > NOW()
                )
            )
    )
`

type CreateMatchmakingTicketParams struct {
	Data                json.RawMessage `db:"data"`
	ExpiresAt           time.Time       `db:"expires_at"`
	IdsSeparatedByComma string          `db:"ids_separated_by_comma"`
}

func (q *Queries) CreateMatchmakingTicket(ctx context.Context, arg CreateMatchmakingTicketParams) (sql.Result, error) {
	return q.db.ExecContext(ctx, CreateMatchmakingTicket, arg.Data, arg.ExpiresAt, arg.IdsSeparatedByComma)
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

const CreateMatchmakingTicketUser = `-- name: CreateMatchmakingTicketUser :execresult
INSERT INTO matchmaking_ticket_user (matchmaking_ticket_id, matchmaking_user_id)
VALUES (?, ?)
`

type CreateMatchmakingTicketUserParams struct {
	MatchmakingTicketID uint64 `db:"matchmaking_ticket_id"`
	MatchmakingUserID   uint64 `db:"matchmaking_user_id"`
}

func (q *Queries) CreateMatchmakingTicketUser(ctx context.Context, arg CreateMatchmakingTicketUserParams) (sql.Result, error) {
	return q.db.ExecContext(ctx, CreateMatchmakingTicketUser, arg.MatchmakingTicketID, arg.MatchmakingUserID)
}

const CreateMatchmakingUser = `-- name: CreateMatchmakingUser :execresult
INSERT INTO matchmaking_user (client_user_id, data)
VALUES (?, ?)
`

type CreateMatchmakingUserParams struct {
	ClientUserID uint64          `db:"client_user_id"`
	Data         json.RawMessage `db:"data"`
}

func (q *Queries) CreateMatchmakingUser(ctx context.Context, arg CreateMatchmakingUserParams) (sql.Result, error) {
	return q.db.ExecContext(ctx, CreateMatchmakingUser, arg.ClientUserID, arg.Data)
}

const CreateMatchmakingUserElo = `-- name: CreateMatchmakingUserElo :execresult
INSERT INTO matchmaking_user_elo (matchmaking_user_id, matchmaking_arena_id, elo)
VALUES (?, ?, ?)
`

type CreateMatchmakingUserEloParams struct {
	MatchmakingUserID  uint64 `db:"matchmaking_user_id"`
	MatchmakingArenaID uint64 `db:"matchmaking_arena_id"`
	Elo                int32  `db:"elo"`
}

func (q *Queries) CreateMatchmakingUserElo(ctx context.Context, arg CreateMatchmakingUserEloParams) (sql.Result, error) {
	return q.db.ExecContext(ctx, CreateMatchmakingUserElo, arg.MatchmakingUserID, arg.MatchmakingArenaID, arg.Elo)
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
    client_user_id,
    elos,
    data,
    created_at,
    updated_at
FROM matchmaking_user_with_elo
ORDER BY client_user_id ASC
LIMIT ? OFFSET ?
`

type GetMatchmakingUsersParams struct {
	Limit  int32 `db:"limit"`
	Offset int32 `db:"offset"`
}

func (q *Queries) GetMatchmakingUsers(ctx context.Context, arg GetMatchmakingUsersParams) ([]MatchmakingUserWithElo, error) {
	rows, err := q.db.QueryContext(ctx, GetMatchmakingUsers, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []MatchmakingUserWithElo
	for rows.Next() {
		var i MatchmakingUserWithElo
		if err := rows.Scan(
			&i.ID,
			&i.ClientUserID,
			&i.Elos,
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

const SetAllMatchmakingUserElos = `-- name: SetAllMatchmakingUserElos :execresult
INSERT INTO matchmaking_user_elo (matchmaking_user_id, matchmaking_arena_id, elo)
SELECT mu.id,
    ma.id,
    ?
FROM matchmaking_user mu
    JOIN matchmaking_arena ma
WHERE mu.id = ?
`

type SetAllMatchmakingUserElosParams struct {
	Elo int32  `db:"elo"`
	ID  uint64 `db:"id"`
}

func (q *Queries) SetAllMatchmakingUserElos(ctx context.Context, arg SetAllMatchmakingUserElosParams) (sql.Result, error) {
	return q.db.ExecContext(ctx, SetAllMatchmakingUserElos, arg.Elo, arg.ID)
}
