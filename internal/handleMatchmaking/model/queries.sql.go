// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0
// source: queries.sql

package model

import (
	"context"
	"database/sql"
	"encoding/json"
	"time"
)

const AddMatchIDToTicket = `-- name: AddMatchIDToTicket :execresult
UPDATE matchmaking_ticket
SET matchmaking_match_id = ?
WHERE id = ?
    AND matchmaking_match_id IS NULL
`

type AddMatchIDToTicketParams struct {
	MatchmakingMatchID sql.NullInt64 `db:"matchmaking_match_id"`
	ID                 uint64        `db:"id"`
}

func (q *Queries) AddMatchIDToTicket(ctx context.Context, arg AddMatchIDToTicketParams) (sql.Result, error) {
	return q.db.ExecContext(ctx, AddMatchIDToTicket, arg.MatchmakingMatchID, arg.ID)
}

const CreateMatch = `-- name: CreateMatch :execresult
INSERT INTO matchmaking_match (matchmaking_arena_id, data)
VALUES (?, "{}")
`

func (q *Queries) CreateMatch(ctx context.Context, matchmakingArenaID uint64) (sql.Result, error) {
	return q.db.ExecContext(ctx, CreateMatch, matchmakingArenaID)
}

const GetAgedMatchmakingTickets = `-- name: GetAgedMatchmakingTickets :many
SELECT id,
    matchmaking_match_id,
    elo_window,
    data,
    expires_at,
    created_at,
    updated_at
FROM matchmaking_ticket
WHERE expires_at < NOW()
    AND matchmaking_match_id IS NULL
    AND elo_window >= ?
LIMIT ? OFFSET ?
`

type GetAgedMatchmakingTicketsParams struct {
	EloWindowMax uint32 `db:"elo_window_max"`
	Limit        int32  `db:"limit"`
	Offset       int32  `db:"offset"`
}

func (q *Queries) GetAgedMatchmakingTickets(ctx context.Context, arg GetAgedMatchmakingTicketsParams) ([]MatchmakingTicket, error) {
	rows, err := q.db.QueryContext(ctx, GetAgedMatchmakingTickets, arg.EloWindowMax, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []MatchmakingTicket
	for rows.Next() {
		var i MatchmakingTicket
		if err := rows.Scan(
			&i.ID,
			&i.MatchmakingMatchID,
			&i.EloWindow,
			&i.Data,
			&i.ExpiresAt,
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

const GetClosestMatch = `-- name: GetClosestMatch :one
WITH ticket_info AS (
    SELECT mtu.matchmaking_ticket_id,
        COUNT(*) AS user_count,
        AVG(mu.elo) AS avg_elo
    FROM matchmaking_ticket_user mtu
        JOIN matchmaking_user mu ON mtu.matchmaking_user_id = mu.id
    WHERE mtu.matchmaking_ticket_id = ?
    GROUP BY mtu.matchmaking_ticket_id
),
ticket_arenas AS (
    SELECT mta.matchmaking_arena_id
    FROM matchmaking_ticket_arena mta
    WHERE mta.matchmaking_ticket_id = ?
)
SELECT mm.id AS match_id,
    ma.name AS arena_name,
    ma.max_players,
    current_players,
    (ma.max_players - current_players) AS remaining_capacity,
    ti.user_count AS ticket_user_count,
    match_avg_elo,
    ti.avg_elo AS ticket_avg_elo,
    ABS(match_avg_elo - ti.avg_elo) AS elo_difference,
    mm.locked_at -- Added for visibility
FROM matchmaking_match mm
    JOIN matchmaking_arena ma ON mm.matchmaking_arena_id = ma.id
    JOIN ticket_arenas ta ON mm.matchmaking_arena_id = ta.matchmaking_arena_id
    JOIN (
        SELECT mt.matchmaking_match_id,
            COUNT(*) AS current_players,
            AVG(mu.elo) AS match_avg_elo
        FROM matchmaking_ticket mt
            JOIN matchmaking_ticket_user mtu ON mt.id = mtu.matchmaking_ticket_id
            JOIN matchmaking_user mu ON mtu.matchmaking_user_id = mu.id
        WHERE mt.matchmaking_match_id IS NOT NULL
        GROUP BY mt.matchmaking_match_id
    ) match_stats ON mm.id = match_stats.matchmaking_match_id
    JOIN ticket_info ti
WHERE ti.user_count <= (ma.max_players - current_players)
    AND (
        mm.locked_at IS NULL
        OR mm.locked_at > NOW()
    )
ORDER BY elo_difference ASC
LIMIT 1
`

type GetClosestMatchParams struct {
	TicketID uint64 `db:"ticket_id"`
}

type GetClosestMatchRow struct {
	MatchID           uint64       `db:"match_id"`
	ArenaName         string       `db:"arena_name"`
	MaxPlayers        uint32       `db:"max_players"`
	CurrentPlayers    int64        `db:"current_players"`
	RemainingCapacity int32        `db:"remaining_capacity"`
	TicketUserCount   int64        `db:"ticket_user_count"`
	MatchAvgElo       interface{}  `db:"match_avg_elo"`
	TicketAvgElo      interface{}  `db:"ticket_avg_elo"`
	EloDifference     int32        `db:"elo_difference"`
	LockedAt          sql.NullTime `db:"locked_at"`
}

func (q *Queries) GetClosestMatch(ctx context.Context, arg GetClosestMatchParams) (GetClosestMatchRow, error) {
	row := q.db.QueryRowContext(ctx, GetClosestMatch, arg.TicketID, arg.TicketID)
	var i GetClosestMatchRow
	err := row.Scan(
		&i.MatchID,
		&i.ArenaName,
		&i.MaxPlayers,
		&i.CurrentPlayers,
		&i.RemainingCapacity,
		&i.TicketUserCount,
		&i.MatchAvgElo,
		&i.TicketAvgElo,
		&i.EloDifference,
		&i.LockedAt,
	)
	return i, err
}

const GetMostPopularArenaOnTicket = `-- name: GetMostPopularArenaOnTicket :one
SELECT ma.id,
    ma.name,
    ma.min_players,
    ma.max_players_per_ticket,
    ma.max_players,
    ma.data,
    ma.created_at,
    ma.updated_at,
    COUNT(mt.id) AS ticket_count
FROM matchmaking_ticket mt
    JOIN matchmaking_ticket_arena mta ON mt.id = mta.matchmaking_ticket_id
    JOIN matchmaking_arena ma ON mta.matchmaking_arena_id = ma.id
WHERE ma.id IN (
        SELECT mta.matchmaking_arena_id
        FROM matchmaking_ticket_arena mta
        WHERE mta.matchmaking_ticket_id = ?
    )
GROUP BY ma.id
ORDER BY ticket_count DESC
LIMIT 1
`

type GetMostPopularArenaOnTicketRow struct {
	ID                  uint64          `db:"id"`
	Name                string          `db:"name"`
	MinPlayers          uint32          `db:"min_players"`
	MaxPlayersPerTicket uint32          `db:"max_players_per_ticket"`
	MaxPlayers          uint32          `db:"max_players"`
	Data                json.RawMessage `db:"data"`
	CreatedAt           time.Time       `db:"created_at"`
	UpdatedAt           time.Time       `db:"updated_at"`
	TicketCount         int64           `db:"ticket_count"`
}

func (q *Queries) GetMostPopularArenaOnTicket(ctx context.Context, matchmakingTicketID uint64) (GetMostPopularArenaOnTicketRow, error) {
	row := q.db.QueryRowContext(ctx, GetMostPopularArenaOnTicket, matchmakingTicketID)
	var i GetMostPopularArenaOnTicketRow
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.MinPlayers,
		&i.MaxPlayersPerTicket,
		&i.MaxPlayers,
		&i.Data,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.TicketCount,
	)
	return i, err
}

const GetNonAgedMatchmakingTickets = `-- name: GetNonAgedMatchmakingTickets :many
SELECT id,
    matchmaking_match_id,
    elo_window,
    data,
    expires_at,
    created_at,
    updated_at
FROM matchmaking_ticket
WHERE expires_at < NOW()
    AND matchmaking_match_id IS NULL
    AND elo_window < ?
LIMIT ? OFFSET ?
`

type GetNonAgedMatchmakingTicketsParams struct {
	EloWindowMax uint32 `db:"elo_window_max"`
	Limit        int32  `db:"limit"`
	Offset       int32  `db:"offset"`
}

func (q *Queries) GetNonAgedMatchmakingTickets(ctx context.Context, arg GetNonAgedMatchmakingTicketsParams) ([]MatchmakingTicket, error) {
	rows, err := q.db.QueryContext(ctx, GetNonAgedMatchmakingTickets, arg.EloWindowMax, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []MatchmakingTicket
	for rows.Next() {
		var i MatchmakingTicket
		if err := rows.Scan(
			&i.ID,
			&i.MatchmakingMatchID,
			&i.EloWindow,
			&i.Data,
			&i.ExpiresAt,
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

const IncrementEloWindow = `-- name: IncrementEloWindow :execresult
UPDATE matchmaking_ticket
SET elo_window = elo_window + ?
WHERE expires_at < NOW()
    AND matchmaking_match_id IS NULL
    AND elo_window < ?
`

type IncrementEloWindowParams struct {
	EloWindowIncrement uint32 `db:"elo_window_increment"`
	EloWindowMax       uint32 `db:"elo_window_max"`
}

func (q *Queries) IncrementEloWindow(ctx context.Context, arg IncrementEloWindowParams) (sql.Result, error) {
	return q.db.ExecContext(ctx, IncrementEloWindow, arg.EloWindowIncrement, arg.EloWindowMax)
}
