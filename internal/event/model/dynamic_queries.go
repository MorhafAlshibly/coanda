package model

import (
	"context"
	"database/sql"
	"encoding/json"

	"github.com/doug-martin/goqu/v9"
	_ "github.com/doug-martin/goqu/v9/dialect/mysql"
	"github.com/doug-martin/goqu/v9/exp"
)

var gq = goqu.Dialect("mysql")

type GetEventParams struct {
	ID   sql.NullInt64  `db:"id"`
	Name sql.NullString `db:"name"`
}

// filterGetEventParams filters the GetEventParams to exp.Expression
func filterGetEventParams(arg GetEventParams) exp.Expression {
	expressions := goqu.Ex{}
	if arg.ID.Valid {
		expressions["id"] = arg.ID
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

func (q *Queries) DeleteEvent(ctx context.Context, arg GetEventParams) (sql.Result, error) {
	event := gq.Delete("event").Prepared(true)
	query, args, err := event.Where(filterGetEventParams(arg)).ToSQL()
	if err != nil {
		return nil, err
	}
	return q.db.ExecContext(ctx, query, args...)
}

func (q *Queries) GetEventWithRound(ctx context.Context, arg GetEventParams) ([]EventWithRound, error) {
	event := gq.From("event_with_round").Prepared(true).Select("id", "name", "current_round_id", "current_round_name", "data", "round_id", "round_name", "round_scoring", "round_data", "round_ended_at", "round_created_at", "round_updated_at", "started_at", "created_at", "updated_at")
	query, args, err := event.Where(filterGetEventParams(arg)).Limit(10).ToSQL()
	if err != nil {
		return nil, err
	}
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
	if arg.Event.ID.Valid {
		expressions["event_id"] = arg.Event.ID
	}
	if arg.Event.Name.Valid {
		expressions["event_id"] = gq.From(gq.From("event").Select("id").Where(goqu.Ex{"name": arg.Event.Name}).Limit(1))
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

type UpdateEventParams struct {
	Event GetEventParams
	Data  json.RawMessage `db:"data"`
}

func (q *Queries) UpdateEvent(ctx context.Context, arg UpdateEventParams) (sql.Result, error) {
	event := gq.Update("event").Prepared(true).Set(
		goqu.Record{
			"data": arg.Data,
		},
	)
	query, args, err := event.Where(filterGetEventParams(arg.Event)).ToSQL()
	if err != nil {
		return nil, err
	}
	return q.db.ExecContext(ctx, query, args...)
}

type GetEventRoundParams struct {
	Event GetEventParams
	ID    sql.NullInt64  `db:"id"`
	Name  sql.NullString `db:"name"`
}

// filterGetEventRoundParams filters the GetEventRoundParams to exp.Expression
func filterGetEventRoundParams(arg GetEventRoundParams) exp.Expression {
	expressions := goqu.Ex{}
	if arg.Event.ID.Valid {
		expressions["event_id"] = arg.Event.ID
	}
	if arg.Event.Name.Valid {
		expressions["event_id"] = gq.From(gq.From("event").Select("id").Where(goqu.Ex{"name": arg.Event.Name}).Limit(1))
	}
	if arg.ID.Valid {
		expressions["id"] = arg.ID
	}
	if arg.Name.Valid {
		expressions["name"] = arg.Name
	}
	// If ID and Name are not provided, we want to use the current round
	if !arg.ID.Valid && !arg.Name.Valid {
		// This means the first round greater than the current time
		expressions["ended_at"] = goqu.Op{"gt": goqu.Func("NOW")}
	}
	return expressions
}

func (q *Queries) GetEventRound(ctx context.Context, arg GetEventRoundParams) (EventRound, error) {
	eventRound := gq.From("event_round").Prepared(true)
	// Must order by ended_at to get the current round
	query, args, err := eventRound.Where(filterGetEventRoundParams(arg)).Order(goqu.I("ended_at").Asc()).Limit(1).ToSQL()
	if err != nil {
		return EventRound{}, err
	}
	var i EventRound
	err = q.db.QueryRowContext(ctx, query, args...).Scan(
		&i.ID,
		&i.EventID,
		&i.Name,
		&i.Scoring,
		&i.Data,
		&i.EndedAt,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

type GetEventRoundLeaderboardParams struct {
	EventRound GetEventRoundParams
	Limit      uint64
	Offset     uint64
}

// filterGetEventRoundLeaderboardParams filters the GetEventRoundLeaderboardParams to exp.Expression
func filterGetEventRoundLeaderboardParams(arg GetEventRoundLeaderboardParams) exp.Expression {
	expressions := goqu.Ex{}
	roundExpressions := goqu.Ex{}
	if arg.EventRound.Event.ID.Valid {
		expressions["event_id"] = arg.EventRound.Event.ID
		roundExpressions["event_id"] = arg.EventRound.Event.ID
	}
	if arg.EventRound.Event.Name.Valid {
		expressions["event_id"] = gq.From(gq.From("event").Select("id").Where(goqu.Ex{"name": arg.EventRound.Event.Name}).Limit(1))
		roundExpressions["event_id"] = gq.From(gq.From("event").Select("id").Where(goqu.Ex{"name": arg.EventRound.Event.Name}).Limit(1))
	}
	if arg.EventRound.ID.Valid {
		expressions["event_round_id"] = arg.EventRound.ID
	}
	if arg.EventRound.Name.Valid {
		expressions["round_name"] = arg.EventRound.Name
	}
	// If ID and Name are not provided, we want to use the current round
	if !arg.EventRound.ID.Valid && !arg.EventRound.Name.Valid {
		// This means the first round greater than the current time
		roundExpressions["ended_at"] = goqu.Op{"gt": goqu.Func("NOW")}
		expressions["event_round_id"] = gq.From("event_round").Select("id").Where(roundExpressions).Order(goqu.I("ended_at").Asc()).Limit(1)
	}
	return expressions
}

func (q *Queries) GetEventRoundLeaderboard(ctx context.Context, arg GetEventRoundLeaderboardParams) ([]EventRoundLeaderboard, error) {
	leaderboard := gq.From("event_round_leaderboard").Prepared(true).Select("id", "event_id", "round_name", "event_user_id", "event_round_id", "result", "score", "ranking", "created_at", "updated_at")
	query, args, err := leaderboard.Where(filterGetEventRoundLeaderboardParams(arg)).Limit(uint(arg.Limit)).Offset(uint(arg.Offset)).ToSQL()
	if err != nil {
		return nil, err
	}
	rows, err := q.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []EventRoundLeaderboard
	for rows.Next() {
		var i EventRoundLeaderboard
		if err := rows.Scan(
			&i.ID,
			&i.EventID,
			&i.RoundName,
			&i.EventUserID,
			&i.EventRoundID,
			&i.Result,
			&i.Score,
			&i.Ranking,
			&i.CreatedAt,
			&i.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	return items, nil
}

type UpdateEventRoundParams struct {
	EventRound GetEventRoundParams
	Data       json.RawMessage `db:"data"`
	Scoring    json.RawMessage `db:"scoring"`
}

func (q *Queries) UpdateEventRound(ctx context.Context, arg UpdateEventRoundParams) (sql.Result, error) {
	updateRecord := goqu.Record{}
	if arg.Data != nil {
		updateRecord["data"] = arg.Data
	}
	if arg.Scoring != nil {
		updateRecord["scoring"] = arg.Scoring
	}
	eventRound := gq.Update("event_round").Prepared(true).Set(updateRecord)
	query, args, err := eventRound.Where(filterGetEventRoundParams(arg.EventRound)).ToSQL()
	if err != nil {
		return nil, err
	}
	return q.db.ExecContext(ctx, query, args...)
}

type GetEventUserParams struct {
	Event  GetEventParams
	ID     sql.NullInt64 `db:"id"`
	UserID sql.NullInt64 `db:"user_id"`
}

// filterGetEventUserParams filters the GetEventUserParams to exp.Expression
func filterGetEventUserParams(arg GetEventUserParams) exp.Expression {
	expressions := goqu.Ex{}
	if arg.Event.ID.Valid {
		expressions["event_id"] = arg.Event.ID
	}
	if arg.Event.Name.Valid {
		expressions["event_id"] = gq.From(gq.From("event").Select("id").Where(goqu.Ex{"name": arg.Event.Name}).Limit(1))
	}
	if arg.ID.Valid {
		expressions["id"] = arg.ID
	}
	if arg.UserID.Valid {
		expressions["user_id"] = arg.UserID
	}
	return expressions
}

func (q *Queries) GetEventUser(ctx context.Context, arg GetEventUserParams) (EventLeaderboard, error) {
	eventUser := gq.From("event_leaderboard").Prepared(true)
	query, args, err := eventUser.Where(filterGetEventUserParams(arg)).Limit(1).ToSQL()
	if err != nil {
		return EventLeaderboard{}, err
	}
	var i EventLeaderboard
	err = q.db.QueryRowContext(ctx, query, args...).Scan(
		&i.ID,
		&i.EventID,
		&i.UserID,
		&i.Score,
		&i.Ranking,
		&i.Data,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

type GetEventRoundUsersParams struct {
	EventUser GetEventUserParams
	Limit     uint64
	Offset    uint64
}

// filterGetEventRoundUsersParams filters the GetEventRoundUsersParams to exp.Expression
func filterGetEventRoundUsersParams(arg GetEventRoundUsersParams) exp.Expression {
	expressions := goqu.Ex{}
	if arg.EventUser.Event.ID.Valid {
		expressions["event_id"] = arg.EventUser.Event.ID
	}
	if arg.EventUser.Event.Name.Valid {
		expressions["event_id"] = gq.From(gq.From("event").Select("id").Where(goqu.Ex{"name": arg.EventUser.Event.Name}).Limit(1))
	}
	if arg.EventUser.ID.Valid {
		expressions["event_user_id"] = arg.EventUser.ID
		// If EventUser ID is not provided, then the user ID must be provided
	} else if arg.EventUser.UserID.Valid {
		// We can get the event_user_id from the event_user table, we know that expression["event_id"] is already set
		expressions["event_user_id"] = gq.From(gq.From("event_user").Select("id").Where(goqu.Ex{
			"event_id": expressions["event_id"],
			"user_id":  arg.EventUser.UserID,
		}).Limit(1))
	}
	return expressions
}

func (q *Queries) GetEventRoundUsers(ctx context.Context, arg GetEventRoundUsersParams) ([]EventRoundLeaderboard, error) {
	eventRoundUsers := gq.From("event_round_leaderboard").Prepared(true).Select("id", "event_id", "round_name", "event_user_id", "event_round_id", "result", "score", "ranking", "data", "created_at", "updated_at")
	query, args, err := eventRoundUsers.Where(filterGetEventRoundUsersParams(arg)).Limit(uint(arg.Limit)).Offset(uint(arg.Offset)).ToSQL()
	if err != nil {
		return nil, err
	}
	rows, err := q.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []EventRoundLeaderboard
	for rows.Next() {
		var i EventRoundLeaderboard
		if err := rows.Scan(
			&i.ID,
			&i.EventID,
			&i.RoundName,
			&i.EventUserID,
			&i.EventRoundID,
			&i.Result,
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

type UpdateEventUserParams struct {
	User GetEventUserParams
	Data json.RawMessage `db:"data"`
}

func (q *Queries) UpdateEventUser(ctx context.Context, arg UpdateEventUserParams) (sql.Result, error) {
	eventUser := gq.Update("event_leaderboard").Prepared(true).Set(
		goqu.Record{
			"data": arg.Data,
		},
	)
	query, args, err := eventUser.Where(filterGetEventUserParams(arg.User)).ToSQL()
	if err != nil {
		return nil, err
	}
	return q.db.ExecContext(ctx, query, args...)
}

func (q *Queries) DeleteEventUser(ctx context.Context, arg GetEventUserParams) (sql.Result, error) {
	eventUser := gq.Delete("event_user").Prepared(true)
	query, args, err := eventUser.Where(filterGetEventUserParams(arg)).ToSQL()
	if err != nil {
		return nil, err
	}
	return q.db.ExecContext(ctx, query, args...)
}

type GetEventRoundUserParams struct {
	EventUser GetEventUserParams
	ID        sql.NullInt64  `db:"id"`
	Round     sql.NullString `db:"round"`
}

// filterGetEventRoundUserParams filters the GetEventRoundUserParams to exp.Expression
func filterGetEventRoundUserParams(arg GetEventRoundUserParams) exp.Expression {
	expressions := goqu.Ex{}
	roundExpressions := goqu.Ex{}
	if arg.EventUser.UserID.Valid {
		if arg.EventUser.Event.ID.Valid {
			expressions["event_user_id"] = gq.From("event_user").Select("id").Where(goqu.Ex{
				"event_id": arg.EventUser.Event.ID,
				"user_id":  arg.EventUser.UserID,
			}).Limit(1)
			roundExpressions["event_id"] = arg.EventUser.Event.ID
		}
		if arg.EventUser.Event.Name.Valid {
			expressions["event_user_id"] = gq.From("event_user").Select("id").Where(goqu.Ex{
				"event_id": gq.From("event").Select("id").Where(goqu.Ex{"name": arg.EventUser.Event.Name}).Limit(1),
				"user_id":  arg.EventUser.UserID,
			}).Limit(1)
			roundExpressions["event_id"] = gq.From("event").Select("id").Where(goqu.Ex{"name": arg.EventUser.Event.Name}).Limit(1)
		}
	}
	if arg.ID.Valid {
		expressions["id"] = arg.ID
	}
	if arg.Round.Valid {
		expressions["round"] = arg.Round
	}
	// If ID and Round are not provided, we want to use the current round
	if !arg.ID.Valid && !arg.Round.Valid {
		// This means the first round greater than the current time
		roundExpressions["ended_at"] = goqu.Op{"gt": goqu.Func("NOW")}
		expressions["event_round_id"] = gq.From("event_round").Select("id").Where(roundExpressions).Order(goqu.I("ended_at").Asc()).Limit(1)
	}
	return expressions
}

func (q *Queries) DeleteEventRoundUser(ctx context.Context, arg GetEventRoundUserParams) (sql.Result, error) {
	eventRoundUser := gq.Delete("event_round_user").Prepared(true)
	query, args, err := eventRoundUser.Where(filterGetEventRoundUserParams(arg)).ToSQL()
	if err != nil {
		return nil, err
	}
	return q.db.ExecContext(ctx, query, args...)
}
