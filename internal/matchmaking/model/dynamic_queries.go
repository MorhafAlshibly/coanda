package model

import (
	"context"
	"database/sql"
	"encoding/json"

	"github.com/doug-martin/goqu/v9"
	_ "github.com/doug-martin/goqu/v9/dialect/mysql"
)

var gq = goqu.Dialect("mysql")

type GetArenaParams struct {
	ID   sql.NullInt64  `db:"id"`
	Name sql.NullString `db:"name"`
}

func filterGetArenaParams(arg GetArenaParams) goqu.Expression {
	expressions := goqu.Ex{}
	if arg.ID.Valid {
		expressions["id"] = arg.ID
	}
	if arg.Name.Valid {
		expressions["name"] = arg.Name
	}
	return expressions
}

func (q *Queries) GetArena(ctx context.Context, arg GetArenaParams) (MatchmakingArena, error) {
	arena := gq.From("matchmaking_arena").Prepared(true)
	query, args, err := arena.Where(filterGetArenaParams(arg)).Limit(1).ToSQL()
	if err != nil {
		return MatchmakingArena{}, err
	}
	var i MatchmakingArena
	err = q.db.QueryRowContext(ctx, query, args...).Scan(
		&i.ID,
		&i.Name,
		&i.MinPlayers,
		&i.MaxPlayers,
		&i.Data,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

type UpdateArenaParams struct {
	Arena      GetArenaParams
	Data       json.RawMessage `db:"data"`
	MinPlayers sql.NullInt32   `db:"min_players"`
	MaxPlayers sql.NullInt32   `db:"max_players"`
}

func (q *Queries) UpdateArena(ctx context.Context, arg UpdateArenaParams) (sql.Result, error) {
	arena := gq.Update("matchmaking_arena").Prepared(true)
	updates := goqu.Record{}
	if arg.Data != nil {
		updates["data"] = []byte(arg.Data)
	}
	if arg.MinPlayers.Valid {
		updates["min_players"] = arg.MinPlayers
	}
	if arg.MaxPlayers.Valid {
		updates["max_players"] = arg.MaxPlayers
	}
	arena = arena.Set(updates)
	query, args, err := arena.Where(filterGetArenaParams(arg.Arena)).Limit(1).ToSQL()
	if err != nil {
		return nil, err
	}
	return q.db.ExecContext(ctx, query, args...)
}

type GetMatchmakingUserParams struct {
	ID     sql.NullInt64 `db:"id"`
	UserID sql.NullInt64 `db:"user_id"`
}

func filterGetMatchmakingUserParams(arg GetMatchmakingUserParams) goqu.Expression {
	expressions := goqu.Ex{}
	if arg.ID.Valid {
		expressions["id"] = arg.ID
	}
	if arg.UserID.Valid {
		expressions["user_id"] = arg.UserID
	}
	return expressions
}

func (q *Queries) GetMatchmakingUser(ctx context.Context, arg GetMatchmakingUserParams) (MatchmakingUserWithElo, error) {
	matchmakingUser := gq.From("matchmaking_user_with_elo").Prepared(true)
	query, args, err := matchmakingUser.Where(filterGetMatchmakingUserParams(arg)).Limit(1).ToSQL()
	if err != nil {
		return MatchmakingUserWithElo{}, err
	}
	var i MatchmakingUserWithElo
	err = q.db.QueryRowContext(ctx, query, args...).Scan(
		&i.ID,
		&i.UserID,
		&i.Elos,
		&i.Data,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

type UpdateMatchmakingUserParams struct {
	MatchmakingUser GetMatchmakingUserParams
	Data            json.RawMessage `db:"data"`
}

func (q *Queries) UpdateMatchmakingUser(ctx context.Context, arg UpdateMatchmakingUserParams) (sql.Result, error) {
	matchmakingUser := gq.Update("matchmaking_user").Prepared(true)
	updates := goqu.Record{}
	if arg.Data != nil {
		updates["data"] = []byte(arg.Data)
	}
	matchmakingUser = matchmakingUser.Set(updates)
	query, args, err := matchmakingUser.Where(filterGetMatchmakingUserParams(arg.MatchmakingUser)).Limit(1).ToSQL()
	if err != nil {
		return nil, err
	}
	return q.db.ExecContext(ctx, query, args...)
}

type SetMatchmakingUserEloParams struct {
	MatchmakingUser GetMatchmakingUserParams
	Elo             sql.NullInt64 `db:"elo"`
	IncrementElo    bool
}

func (q *Queries) SetMatchmakingUserElo(ctx context.Context, arg SetMatchmakingUserEloParams) (sql.Result, error) {
	matchmakingUser := gq.Update("matchmaking_user").Prepared(true)
	updates := goqu.Record{}
	if arg.Elo.Valid {
		if arg.IncrementElo {
			updates["elo"] = goqu.L("elo + ?", arg.Elo)
		} else {
			updates["elo"] = arg.Elo
		}
	}
	matchmakingUser = matchmakingUser.Set(updates)
	query, args, err := matchmakingUser.Where(filterGetMatchmakingUserParams(arg.MatchmakingUser)).Limit(1).ToSQL()
	if err != nil {
		return nil, err
	}
	return q.db.ExecContext(ctx, query, args...)
}
