package model

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"
)

type GetTournamentParams struct {
	ID                          sql.NullInt64 `db:"id"`
	NameIntervalUserIDStartedAt NullNameIntervalUserIDStartedAt
}

type NullNameIntervalUserIDStartedAt struct {
	Name                string                       `db:"name"`
	TournamentInterval  TournamentTournamentInterval `db:"tournament_interval"`
	UserID              uint64                       `db:"user_id"`
	TournamentStartedAt time.Time                    `db:"tournament_started_at"`
	Valid               bool
}

func (q *Queries) GetTournament(ctx context.Context, arg GetTournamentParams) (RankedTournament, error) {
	tournaments := sq.Select("*").From("ranked_tournament")
	if arg.ID.Valid {
		tournaments = tournaments.Where(sq.Eq{"id": arg.ID})
	}
	if arg.NameIntervalUserIDStartedAt.Valid {
		tournaments = tournaments.Where(sq.Eq{"name": arg.NameIntervalUserIDStartedAt.Name})
		tournaments = tournaments.Where(sq.Eq{"tournament_interval": arg.NameIntervalUserIDStartedAt.TournamentInterval})
		tournaments = tournaments.Where(sq.Eq{"user_id": arg.NameIntervalUserIDStartedAt.UserID})
		tournaments = tournaments.Where(sq.Eq{"tournament_started_at": arg.NameIntervalUserIDStartedAt.TournamentStartedAt})
	}
	sql, args, err := tournaments.Limit(1).ToSql()
	if err != nil {
		return RankedTournament{}, err
	}
	row := q.db.QueryRowContext(ctx, sql, args...)
	var i RankedTournament
	err = row.Scan(
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
	)
	return i, err
}

type GetTournamentsParams struct {
	Name                sql.NullString               `db:"name"`
	UserID              sql.NullInt64                `db:"user_id"`
	TournamentInterval  TournamentTournamentInterval `db:"tournament_interval"`
	TournamentStartedAt time.Time                    `db:"tournament_started_at"`
	Limit               uint64                       `db:"limit"`
	Offset              uint64                       `db:"offset"`
}

func (q *Queries) GetTournaments(ctx context.Context, arg GetTournamentsParams) ([]RankedTournament, error) {
	tournaments := sq.Select("*").From("ranked_tournament")
	if arg.Name.Valid {
		tournaments = tournaments.Where(sq.Eq{"name": arg.Name})
	}
	if arg.UserID.Valid {
		tournaments = tournaments.Where(sq.Eq{"user_id": arg.UserID})
	}
	tournaments = tournaments.Where(sq.Eq{"tournament_interval": arg.TournamentInterval})
	tournaments = tournaments.Where(sq.Eq{"tournament_started_at": arg.TournamentStartedAt})
	sql, args, err := tournaments.Limit(arg.Limit).Offset(arg.Offset).ToSql()
	if err != nil {
		return nil, err
	}
	fmt.Println(sql)
	fmt.Printf("%+v\n", args)
	rows, err := q.db.QueryContext(ctx, sql, args...)
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

func (q *Queries) DeleteTournament(ctx context.Context, arg GetTournamentParams) (sql.Result, error) {
	tournaments := sq.Delete("tournament")
	if arg.ID.Valid {
		tournaments = tournaments.Where(sq.Eq{"id": arg.ID})
	} else if arg.NameIntervalUserIDStartedAt.Valid {
		tournaments = tournaments.Where(sq.Eq{"name": arg.NameIntervalUserIDStartedAt.Name})
		tournaments = tournaments.Where(sq.Eq{"tournament_interval": arg.NameIntervalUserIDStartedAt.TournamentInterval})
		tournaments = tournaments.Where(sq.Eq{"user_id": arg.NameIntervalUserIDStartedAt.UserID})
		tournaments = tournaments.Where(sq.Eq{"tournament_started_at": arg.NameIntervalUserIDStartedAt.TournamentStartedAt})
	} else {
		return nil, errors.New("no valid parameters")
	}
	sql, args, err := tournaments.Limit(1).ToSql()
	if err != nil {
		return nil, err
	}
	return q.db.ExecContext(ctx, sql, args...)
}

type UpdateTournamentParams struct {
	ID                          sql.NullInt64 `db:"id"`
	NameIntervalUserIDStartedAt NullNameIntervalUserIDStartedAt
	Data                        json.RawMessage `db:"data"`
	Score                       sql.NullInt64   `db:"score"`
	IncrementScore              bool
}

func (q *Queries) UpdateTournament(ctx context.Context, arg UpdateTournamentParams) (sql.Result, error) {
	tournaments := sq.Update("tournament")
	if arg.ID.Valid {
		tournaments = tournaments.Where(sq.Eq{"id": arg.ID})
	}
	if arg.NameIntervalUserIDStartedAt.Valid {
		tournaments = tournaments.Where(sq.Eq{"name": arg.NameIntervalUserIDStartedAt.Name})
		tournaments = tournaments.Where(sq.Eq{"tournament_interval": arg.NameIntervalUserIDStartedAt.TournamentInterval})
		tournaments = tournaments.Where(sq.Eq{"user_id": arg.NameIntervalUserIDStartedAt.UserID})
		tournaments = tournaments.Where(sq.Eq{"tournament_started_at": arg.NameIntervalUserIDStartedAt.TournamentStartedAt})
	}
	if arg.Data != nil {
		tournaments = tournaments.Set("data", arg.Data)
	}
	if arg.Score.Valid {
		if arg.IncrementScore {
			tournaments = tournaments.Set("score", sq.Expr("score + ?", arg.Score))
		} else {
			tournaments = tournaments.Set("score", arg.Score)
		}
	}
	sql, args, err := tournaments.Limit(1).ToSql()
	if err != nil {
		return nil, err
	}
	return q.db.ExecContext(ctx, sql, args...)
}
