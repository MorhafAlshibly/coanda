package model

import (
	"context"
	"database/sql"
	"encoding/json"
	"time"

	"github.com/doug-martin/goqu/v9"
	_ "github.com/doug-martin/goqu/v9/dialect/mysql"
	"github.com/doug-martin/goqu/v9/exp"
)

var gq = goqu.Dialect("mysql")

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

// filterGetTournamentParams filters the GetTournamentParams to exp.Expression
func filterGetTournamentParams(arg GetTournamentParams) exp.Expression {
	expressions := goqu.Ex{}
	if arg.ID.Valid {
		expressions["id"] = arg.ID
	}
	if arg.NameIntervalUserIDStartedAt.Valid {
		expressions["name"] = arg.NameIntervalUserIDStartedAt.Name
		expressions["tournament_interval"] = arg.NameIntervalUserIDStartedAt.TournamentInterval
		expressions["user_id"] = arg.NameIntervalUserIDStartedAt.UserID
		expressions["tournament_started_at"] = arg.NameIntervalUserIDStartedAt.TournamentStartedAt
	}
	return expressions
}

func (q *Queries) GetTournament(ctx context.Context, arg GetTournamentParams) (RankedTournament, error) {
	tournament := gq.From("ranked_tournament").Prepared(true).Select("id", "name", "tournament_interval", "user_id", "score", "ranking", "data", "tournament_started_at", "sent_to_third_party_at", "created_at", "updated_at")
	query, args, err := tournament.Where(filterGetTournamentParams(arg)).Limit(1).ToSQL()
	if err != nil {
		return RankedTournament{}, err
	}
	var i RankedTournament
	err = q.db.QueryRowContext(ctx, query, args...).Scan(
		&i.ID,
		&i.Name,
		&i.TournamentInterval,
		&i.UserID,
		&i.Score,
		&i.Ranking,
		&i.Data,
		&i.TournamentStartedAt,
		&i.SentToThirdPartyAt,
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

// filterGetTournamentsParams filters the GetTournamentsParams to exp.Expression
func filterGetTournamentsParams(arg GetTournamentsParams) exp.Expression {
	expressions := goqu.Ex{}
	if arg.Name.Valid {
		expressions["name"] = arg.Name
	}
	if arg.UserID.Valid {
		expressions["user_id"] = arg.UserID
	}
	expressions["tournament_interval"] = arg.TournamentInterval
	expressions["tournament_started_at"] = arg.TournamentStartedAt
	return expressions
}

func (q *Queries) GetTournaments(ctx context.Context, arg GetTournamentsParams) ([]RankedTournament, error) {
	tournament := gq.From("ranked_tournament").Prepared(true).Select("id", "name", "tournament_interval", "user_id", "score", "ranking", "data", "tournament_started_at", "sent_to_third_party_at", "created_at", "updated_at")
	query, args, err := tournament.Where(filterGetTournamentsParams(arg)).Limit(uint(arg.Limit)).Offset(uint(arg.Offset)).ToSQL()
	if err != nil {
		return nil, err
	}
	rows, err := q.db.QueryContext(ctx, query, args...)
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
			&i.SentToThirdPartyAt,
			&i.CreatedAt,
			&i.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	return items, nil
}

func (q *Queries) DeleteTournament(ctx context.Context, arg GetTournamentParams) (sql.Result, error) {
	tournament := gq.Delete("tournament").Prepared(true)
	query, args, err := tournament.Where(filterGetTournamentParams(arg)).Limit(1).ToSQL()
	if err != nil {
		return nil, err
	}
	return q.db.ExecContext(ctx, query, args...)
}

type UpdateTournamentParams struct {
	Tournament     GetTournamentParams
	Data           json.RawMessage `db:"data"`
	Score          sql.NullInt64   `db:"score"`
	IncrementScore bool
}

func (q *Queries) UpdateTournament(ctx context.Context, arg UpdateTournamentParams) (sql.Result, error) {
	tournament := gq.Update("tournament").Prepared(true)
	updates := goqu.Record{}
	if arg.Data != nil {
		updates["data"] = []byte(arg.Data)
	}
	if arg.Score.Valid {
		if arg.IncrementScore {
			updates["score"] = goqu.L("score + ?", arg.Score)
		} else {
			updates["score"] = arg.Score
		}
	}
	query, args, err := tournament.Where(filterGetTournamentParams(arg.Tournament)).Set(updates).Limit(1).ToSQL()
	if err != nil {
		return nil, err
	}
	return q.db.ExecContext(ctx, query, args...)
}
