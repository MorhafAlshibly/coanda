package model

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/doug-martin/goqu/v9"
	"github.com/doug-martin/goqu/v9/exp"
)

var gq = goqu.Dialect("mysql")

type GetEventParams struct {
	Id   sql.NullInt64  `db:"id"`
	Name sql.NullString `db:"name"`
}

// filterGetEventParams filters the GetEventParams to exp.Expression
func filterGetEventParams(arg GetEventParams) exp.Expression {
	expressions := goqu.Ex{}
	if arg.Id.Valid {
		expressions["id"] = arg.Id
	}
	if arg.Name.Valid {
		expressions["name"] = arg.Name
	}
	return expressions
}

func (q *Queries) GetEvent(ctx context.Context, arg GetEventParams) (Event, error) {
	event := gq.From("event").Prepared(true)
	query, args, err := event.Where(filterGetEventParams(arg)).Limit(1).ToSQL()
	if err != nil {
		return Event{}, err
	}
	var i Event
	err = q.db.QueryRowContext(ctx, query, args...).Scan(
		&i.ID,
		&i.Name,
		&i.Data,
		&i.StartedAt,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

func (q *Queries) GetEventWithRound(ctx context.Context, arg GetEventParams) ([]EventWithRound, error) {
	event := gq.From("event_with_round").Prepared(true).Select("id", "name", "current_round_id", "current_round_name", "data", "round_id", "round_name", "round_scoring", "round_data", "round_ended_at", "round_created_at", "round_updated_at", "started_at", "created_at", "updated_at")
	query, args, err := event.Where(filterGetEventParams(arg)).ToSQL()
	if err != nil {
		return nil, err
	}
	fmt.Println(query)
	rows, err := q.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []EventWithRound
	for rows.Next() {
		var i EventWithRound
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.CurrentRoundID,
			&i.CurrentRoundName,
			&i.Data,
			&i.RoundID,
			&i.RoundName,
			&i.RoundScoring,
			&i.RoundData,
			&i.RoundEndedAt,
			&i.RoundCreatedAt,
			&i.RoundUpdatedAt,
			&i.StartedAt,
			&i.CreatedAt,
			&i.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	return items, nil
}

type GetEventLeaderboardParams struct {
	Event  GetEventParams
	Limit  uint64
	Offset uint64
}

// filterGetEventLeaderboardParams filters the GetEventLeaderboardParams to exp.Expression
func filterGetEventLeaderboardParams(arg GetEventLeaderboardParams) exp.Expression {
	expressions := goqu.Ex{}
	if arg.Event.Id.Valid {
		expressions["event"] = arg.Event.Id
	}
	if arg.Event.Name.Valid {
		expressions["event"] = arg.Event.Name
	}
	return expressions
}

func (q *Queries) GetEventLeaderboard(ctx context.Context, arg GetEventLeaderboardParams) ([]EventLeaderboard, error) {
	leaderboard := gq.From("event_leaderboard").Prepared(true).Select("id", "event_id", "user_id", "score", "ranking", "data", "created_at", "updated_at")
	query, args, err := leaderboard.Where(filterGetEventLeaderboardParams(arg)).Limit(uint(arg.Limit)).Offset(uint(arg.Offset)).ToSQL()
	if err != nil {
		return nil, err
	}
	rows, err := q.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []EventLeaderboard
	for rows.Next() {
		var i EventLeaderboard
		if err := rows.Scan(
			&i.ID,
			&i.EventID,
			&i.UserID,
			&i.Score,
			&i.Ranking,
			&i.Data,
			&i.CreatedAt,
			&i.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	return items, nil
}
