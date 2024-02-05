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

const createTournament = `-- name: CreateTournament :execresult
INSERT INTO tournament (
        name,
        tournament_interval,
        user_id,
        score,
        data,
        tournament_started_at
    )
VALUES (?, ?, ?, ?, ?, ?)
`

type CreateTournamentParams struct {
	Name                string
	TournamentInterval  TournamentTournamentInterval
	UserID              uint64
	Score               int64
	Data                json.RawMessage
	TournamentStartedAt time.Time
}

func (q *Queries) CreateTournament(ctx context.Context, arg CreateTournamentParams) (sql.Result, error) {
	return q.db.ExecContext(ctx, createTournament,
		arg.Name,
		arg.TournamentInterval,
		arg.UserID,
		arg.Score,
		arg.Data,
		arg.TournamentStartedAt,
	)
}

const deleteTournament = `-- name: DeleteTournament :execresult
DELETE FROM tournament
WHERE name = ?
    AND tournament_interval = ?
    AND user_id = ?
LIMIT 1
`

type DeleteTournamentParams struct {
	Name               string
	TournamentInterval TournamentTournamentInterval
	UserID             uint64
}

func (q *Queries) DeleteTournament(ctx context.Context, arg DeleteTournamentParams) (sql.Result, error) {
	return q.db.ExecContext(ctx, deleteTournament, arg.Name, arg.TournamentInterval, arg.UserID)
}

const getTournament = `-- name: GetTournament :one
SELECT name,
    tournament_interval,
    user_id,
    score,
    ranking,
    data,
    tournament_started_at,
    created_at,
    updated_at
FROM ranked_tournament
WHERE name = ?
    AND tournament_interval = ?
    AND user_id = ?
LIMIT 1
`

type GetTournamentParams struct {
	Name               string
	TournamentInterval TournamentTournamentInterval
	UserID             uint64
}

func (q *Queries) GetTournament(ctx context.Context, arg GetTournamentParams) (RankedTournament, error) {
	row := q.db.QueryRowContext(ctx, getTournament, arg.Name, arg.TournamentInterval, arg.UserID)
	var i RankedTournament
	err := row.Scan(
		&i.Name,
		&i.TournamentInterval,
		&i.UserID,
		&i.Score,
		&i.Ranking,
		&i.Data,
		&i.TournamentStartedAt,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getTournaments = `-- name: GetTournaments :many
SELECT name,
    tournament_interval,
    user_id,
    score,
    ranking,
    data,
    tournament_started_at,
    created_at,
    updated_at
FROM ranked_tournament
WHERE name = ?
    OR tournament_interval = ?
    OR user_id = ?
ORDER BY name ASC,
    tournament_interval ASC,
    score DESC
LIMIT ? OFFSET ?
`

type GetTournamentsParams struct {
	Name               sql.NullString
	TournamentInterval NullTournamentTournamentInterval
	UserID             sql.NullInt64
	Limit              int32
	Offset             int32
}

func (q *Queries) GetTournaments(ctx context.Context, arg GetTournamentsParams) ([]RankedTournament, error) {
	rows, err := q.db.QueryContext(ctx, getTournaments,
		arg.Name,
		arg.TournamentInterval,
		arg.UserID,
		arg.Limit,
		arg.Offset,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []RankedTournament
	for rows.Next() {
		var i RankedTournament
		if err := rows.Scan(
			&i.Name,
			&i.TournamentInterval,
			&i.UserID,
			&i.Score,
			&i.Ranking,
			&i.Data,
			&i.TournamentStartedAt,
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

const updateTournamentData = `-- name: UpdateTournamentData :execresult
UPDATE tournament
SET data = ?
WHERE name = ?
    AND tournament_interval = ?
    AND user_id = ?
LIMIT 1
`

type UpdateTournamentDataParams struct {
	Data               json.RawMessage
	Name               string
	TournamentInterval TournamentTournamentInterval
	UserID             uint64
}

func (q *Queries) UpdateTournamentData(ctx context.Context, arg UpdateTournamentDataParams) (sql.Result, error) {
	return q.db.ExecContext(ctx, updateTournamentData,
		arg.Data,
		arg.Name,
		arg.TournamentInterval,
		arg.UserID,
	)
}

const updateTournamentScore = `-- name: UpdateTournamentScore :execresult
UPDATE tournament
SET score = CASE
        WHEN ? IS NOT NULL THEN ? + CASE
            WHEN CAST(? as unsigned) != 0 THEN score
            ELSE 0
        END
        ELSE score
    END
WHERE name = ?
    AND tournament_interval = ?
    AND user_id = ?
LIMIT 1
`

type UpdateTournamentScoreParams struct {
	Score              int64
	IncrementScore     int64
	Name               string
	TournamentInterval TournamentTournamentInterval
	UserID             uint64
}

func (q *Queries) UpdateTournamentScore(ctx context.Context, arg UpdateTournamentScoreParams) (sql.Result, error) {
	return q.db.ExecContext(ctx, updateTournamentScore,
		arg.Score,
		arg.Score,
		arg.IncrementScore,
		arg.Name,
		arg.TournamentInterval,
		arg.UserID,
	)
}
