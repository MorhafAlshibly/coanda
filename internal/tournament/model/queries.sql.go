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

const CreateTournament = `-- name: CreateTournament :execresult
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
	Name                string                       `db:"name"`
	TournamentInterval  TournamentTournamentInterval `db:"tournament_interval"`
	UserID              uint64                       `db:"user_id"`
	Score               int64                        `db:"score"`
	Data                json.RawMessage              `db:"data"`
	TournamentStartedAt time.Time                    `db:"tournament_started_at"`
}

func (q *Queries) CreateTournament(ctx context.Context, arg CreateTournamentParams) (sql.Result, error) {
	return q.db.ExecContext(ctx, CreateTournament,
		arg.Name,
		arg.TournamentInterval,
		arg.UserID,
		arg.Score,
		arg.Data,
		arg.TournamentStartedAt,
	)
}

const GetTournamentsBeforeWipe = `-- name: GetTournamentsBeforeWipe :many
SELECT id, name, tournament_interval, user_id, score, ranking, data, tournament_started_at, created_at, updated_at
FROM ranked_tournament
WHERE tournament_started_at < ?
    AND tournament_interval = ?
LIMIT ? OFFSET ?
`

type GetTournamentsBeforeWipeParams struct {
	TournamentStartedAt time.Time                    `db:"tournament_started_at"`
	TournamentInterval  TournamentTournamentInterval `db:"tournament_interval"`
	Limit               int32                        `db:"limit"`
	Offset              int32                        `db:"offset"`
}

func (q *Queries) GetTournamentsBeforeWipe(ctx context.Context, arg GetTournamentsBeforeWipeParams) ([]RankedTournament, error) {
	rows, err := q.db.QueryContext(ctx, GetTournamentsBeforeWipe,
		arg.TournamentStartedAt,
		arg.TournamentInterval,
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
			&i.ID,
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

const WipeTournaments = `-- name: WipeTournaments :execresult
DELETE FROM tournament
WHERE tournament_started_at < ?
    AND tournament_interval = ?
`

type WipeTournamentsParams struct {
	TournamentStartedAt time.Time                    `db:"tournament_started_at"`
	TournamentInterval  TournamentTournamentInterval `db:"tournament_interval"`
}

func (q *Queries) WipeTournaments(ctx context.Context, arg WipeTournamentsParams) (sql.Result, error) {
	return q.db.ExecContext(ctx, WipeTournaments, arg.TournamentStartedAt, arg.TournamentInterval)
}
