package model

import (
	"context"
	"database/sql"
	"encoding/json"
	"time"

	"github.com/MorhafAlshibly/coanda/pkg/conversion"
	"github.com/MorhafAlshibly/coanda/pkg/goquOptions"
	"github.com/doug-martin/goqu/v9"
	_ "github.com/doug-martin/goqu/v9/dialect/mysql"
)

var gq = goqu.Dialect("mysql")

type ArenaParams struct {
	ID   sql.NullInt64  `db:"id"`
	Name sql.NullString `db:"name"`
}

func filterGetArenaParams(arg ArenaParams) goqu.Expression {
	expressions := goqu.Ex{}
	if arg.ID.Valid {
		expressions["id"] = arg.ID
	}
	if arg.Name.Valid {
		expressions["name"] = arg.Name
	}
	return expressions
}

func (q *Queries) GetArena(ctx context.Context, arg ArenaParams, opts *goquOptions.SelectDataset) (MatchmakingArena, error) {
	arena := gq.From("matchmaking_arena").Prepared(true)
	query, args, err := opts.Apply(arena.Where(filterGetArenaParams(arg)).Limit(1)).ToSQL()
	if err != nil {
		return MatchmakingArena{}, err
	}
	var i MatchmakingArena
	err = q.db.QueryRowContext(ctx, query, args...).Scan(
		&i.ID,
		&i.Name,
		&i.MinPlayers,
		&i.MaxPlayersPerTicket,
		&i.MaxPlayers,
		&i.Data,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

type UpdateArenaParams struct {
	Arena               ArenaParams
	Data                json.RawMessage `db:"data"`
	MinPlayers          sql.NullInt32   `db:"min_players"`
	MaxPlayersPerTicket sql.NullInt32   `db:"max_players_per_ticket"`
	MaxPlayers          sql.NullInt32   `db:"max_players"`
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
	if arg.MaxPlayersPerTicket.Valid {
		updates["max_players_per_ticket"] = arg.MaxPlayersPerTicket
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

type MatchmakingUserParams struct {
	ID           sql.NullInt64 `db:"id"`
	ClientUserID sql.NullInt64 `db:"user_id"`
}

func filterMatchmakingUserParams(arg MatchmakingUserParams) goqu.Expression {
	expressions := goqu.Ex{}
	if arg.ID.Valid {
		expressions["id"] = arg.ID
	}
	if arg.ClientUserID.Valid {
		expressions["client_user_id"] = arg.ClientUserID
	}
	return expressions
}

func (q *Queries) GetMatchmakingUser(ctx context.Context, arg MatchmakingUserParams, opts *goquOptions.SelectDataset) (MatchmakingUser, error) {
	matchmakingUser := gq.From("matchmaking_user").Prepared(true)
	query, args, err := opts.Apply(matchmakingUser.Where(filterMatchmakingUserParams(arg)).Limit(1)).ToSQL()
	if err != nil {
		return MatchmakingUser{}, err
	}
	var i MatchmakingUser
	err = q.db.QueryRowContext(ctx, query, args...).Scan(
		&i.ID,
		&i.MatchmakingTicketID,
		&i.ClientUserID,
		&i.Elo,
		&i.Data,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

type UpdateMatchmakingUserParams struct {
	MatchmakingUser MatchmakingUserParams
	Data            json.RawMessage `db:"data"`
	Elo             sql.NullInt32   `db:"elo"`
}

func (q *Queries) UpdateMatchmakingUser(ctx context.Context, arg UpdateMatchmakingUserParams) (sql.Result, error) {
	matchmakingUser := gq.Update("matchmaking_user").Prepared(true)
	updates := goqu.Record{}
	if arg.Data != nil {
		updates["data"] = []byte(arg.Data)
	}
	if arg.Elo.Valid {
		updates["elo"] = arg.Elo
	}
	matchmakingUser = matchmakingUser.Set(updates)
	query, args, err := matchmakingUser.Where(filterMatchmakingUserParams(arg.MatchmakingUser)).Limit(1).ToSQL()
	if err != nil {
		return nil, err
	}
	return q.db.ExecContext(ctx, query, args...)
}

func (q *Queries) DeleteMatchmakingUser(ctx context.Context, arg MatchmakingUserParams) (sql.Result, error) {
	matchmakingUser := gq.Delete("matchmaking_user").Prepared(true)
	query, args, err := matchmakingUser.Where(
		goqu.And(
			filterMatchmakingUserParams(arg),
			goqu.C("matchmaking_ticket_id").IsNull(),
		),
	).Limit(1).ToSQL()
	if err != nil {
		return nil, err
	}
	return q.db.ExecContext(ctx, query, args...)
}

type MatchmakingTicketParams struct {
	MatchmakingUser MatchmakingUserParams
	ID              sql.NullInt64 `db:"id"`
}

func filterMatchmakingTicketParams(arg MatchmakingTicketParams, idColumnName *string) goqu.Expression {
	if idColumnName == nil {
		idColumnName = conversion.ValueToPointer("id")
	}
	expressions := []goqu.Expression{}
	if arg.ID.Valid {
		expressions = append(expressions, goqu.C(*idColumnName).Eq(arg.ID))
	}
	if arg.MatchmakingUser.ID.Valid {
		expressions = append(expressions, goqu.C(*idColumnName).Eq(gq.From(gq.From("matchmaking_ticket_with_user").Where(goqu.Ex{"matchmaking_user_id": arg.MatchmakingUser.ID}).Select("ticket_id").Limit(1))))
	}
	if arg.MatchmakingUser.ClientUserID.Valid {
		expressions = append(expressions, goqu.C(*idColumnName).Eq(gq.From(gq.From("matchmaking_ticket_with_user").Where(goqu.Ex{"client_user_id": arg.MatchmakingUser.ClientUserID}).Select("ticket_id").Limit(1))))
	}
	return goqu.And(expressions...)
}

type GetMatchmakingTicketParams struct {
	MatchmakingTicket MatchmakingTicketParams
	UserLimit         uint64
	UserOffset        uint64
	ArenaLimit        uint64
	ArenaOffset       uint64
}

func filterGetMatchmakingTicketParams(arg GetMatchmakingTicketParams) goqu.Expression {
	return goqu.And(
		filterMatchmakingTicketParams(arg.MatchmakingTicket, conversion.ValueToPointer("ticket_id")),
		goqu.C("user_number").Gt(arg.UserOffset),
		goqu.C("user_number").Lte(arg.UserOffset+arg.UserLimit),
		goqu.C("arena_number").Gt(arg.ArenaOffset),
		goqu.C("arena_number").Lte(arg.ArenaOffset+arg.ArenaLimit),
	)
}

func (q *Queries) GetMatchmakingTicket(ctx context.Context, arg GetMatchmakingTicketParams) ([]MatchmakingTicketWithUserAndArena, error) {
	matchmakingTicket := gq.From("matchmaking_ticket_with_user_and_arena").Prepared(true)
	query, args, err := matchmakingTicket.Where(filterGetMatchmakingTicketParams(arg)).ToSQL()
	if err != nil {
		return nil, err
	}
	rows, err := q.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []MatchmakingTicketWithUserAndArena
	for rows.Next() {
		var i MatchmakingTicketWithUserAndArena
		if err := rows.Scan(
			&i.TicketID,
			&i.MatchmakingMatchID,
			&i.Status,
			&i.UserCount,
			&i.TicketData,
			&i.TicketCreatedAt,
			&i.TicketUpdatedAt,
			&i.MatchmakingUserID,
			&i.ClientUserID,
			&i.Elo,
			&i.UserNumber,
			&i.UserData,
			&i.UserCreatedAt,
			&i.UserUpdatedAt,
			&i.ArenaID,
			&i.ArenaName,
			&i.ArenaMinPlayers,
			&i.ArenaMaxPlayersPerTicket,
			&i.ArenaMaxPlayers,
			&i.ArenaNumber,
			&i.ArenaData,
			&i.ArenaCreatedAt,
			&i.ArenaUpdatedAt,
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

type GetMatchmakingTicketsParams struct {
	MatchmakingUser    MatchmakingUserParams
	MatchmakingMatchID sql.NullInt64 `db:"matchmaking_match_id"`
	Arena              ArenaParams
	Statuses           []string
	Limit              uint64
	Offset             uint64
	UserLimit          uint64
	UserOffset         uint64
	ArenaLimit         uint64
	ArenaOffset        uint64
}

func filterGetMatchmakingTicketsParams(arg GetMatchmakingTicketsParams) goqu.Expression {
	expressions := []goqu.Expression{}
	if arg.MatchmakingMatchID.Valid {
		expressions = append(expressions, goqu.C("matchmaking_match_id").Eq(arg.MatchmakingMatchID))
	}
	if arg.MatchmakingUser.ID.Valid {
		expressions = append(expressions, goqu.C("matchmaking_user_id").Eq(arg.MatchmakingUser.ID))
	}
	if arg.MatchmakingUser.ClientUserID.Valid {
		expressions = append(expressions, goqu.C("client_user_id").Eq(arg.MatchmakingUser.ClientUserID))
	}
	if arg.Arena.ID.Valid {
		expressions = append(expressions, goqu.C("arena_id").Eq(arg.Arena.ID))
	}
	if arg.Arena.Name.Valid {
		expressions = append(expressions, goqu.C("arena_name").Eq(arg.Arena.Name))
	}
	if len(arg.Statuses) > 0 {
		expressions = append(expressions, goqu.C("status").In(arg.Statuses))
	}
	finalExpression := goqu.And(
		goqu.C("ticket_id").In(gq.From(gq.From("matchmaking_ticket_with_user_and_arena").Select("ticket_id").Where(goqu.And(expressions...)).GroupBy("ticket_id").Order(goqu.C("ticket_id").Asc()).Limit(uint(arg.Limit)).Offset(uint(arg.Offset)))),
		goqu.C("user_number").Gt(arg.UserOffset),
		goqu.C("user_number").Lte(arg.UserOffset+arg.UserLimit),
		goqu.C("arena_number").Gt(arg.ArenaOffset),
		goqu.C("arena_number").Lte(arg.ArenaOffset+arg.ArenaLimit),
	)
	return finalExpression
}

func (q *Queries) GetMatchmakingTickets(ctx context.Context, arg GetMatchmakingTicketsParams) ([]MatchmakingTicketWithUserAndArena, error) {
	matchmakingTicket := gq.From("matchmaking_ticket_with_user_and_arena").Prepared(true)
	query, args, err := matchmakingTicket.Where(filterGetMatchmakingTicketsParams(arg)).ToSQL()
	if err != nil {
		return nil, err
	}
	rows, err := q.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []MatchmakingTicketWithUserAndArena
	for rows.Next() {
		var i MatchmakingTicketWithUserAndArena
		if err = rows.Scan(
			&i.TicketID,
			&i.MatchmakingMatchID,
			&i.Status,
			&i.UserCount,
			&i.TicketData,
			&i.TicketCreatedAt,
			&i.TicketUpdatedAt,
			&i.MatchmakingUserID,
			&i.ClientUserID,
			&i.Elo,
			&i.UserNumber,
			&i.UserData,
			&i.UserCreatedAt,
			&i.UserUpdatedAt,
			&i.ArenaID,
			&i.ArenaName,
			&i.ArenaMinPlayers,
			&i.ArenaMaxPlayersPerTicket,
			&i.ArenaMaxPlayers,
			&i.ArenaNumber,
			&i.ArenaData,
			&i.ArenaCreatedAt,
			&i.ArenaUpdatedAt,
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

type UpdateMatchmakingTicketParams struct {
	MatchmakingTicket MatchmakingTicketParams
	Data              json.RawMessage `db:"data"`
}

func (q *Queries) UpdateMatchmakingTicket(ctx context.Context, arg UpdateMatchmakingTicketParams) (sql.Result, error) {
	matchmakingTicket := gq.Update("matchmaking_ticket").Prepared(true)
	updates := goqu.Record{}
	if arg.Data != nil {
		updates["data"] = []byte(arg.Data)
	}
	matchmakingTicket = matchmakingTicket.Set(updates)
	query, args, err := matchmakingTicket.Where(filterMatchmakingTicketParams(arg.MatchmakingTicket, nil)).Limit(1).ToSQL()
	if err != nil {
		return nil, err
	}
	return q.db.ExecContext(ctx, query, args...)
}

func (q *Queries) DeleteMatchmakingTicket(ctx context.Context, arg MatchmakingTicketParams) (sql.Result, error) {
	matchmakingTicket := gq.Delete("matchmaking_ticket").Prepared(true)
	query, args, err := matchmakingTicket.Where(
		goqu.And(
			filterMatchmakingTicketParams(arg, nil),
			goqu.C("matchmaking_match_id").IsNull(),
		),
	).Limit(1).ToSQL()
	if err != nil {
		return nil, err
	}
	return q.db.ExecContext(ctx, query, args...)
}

type MatchParams struct {
	MatchmakingTicket MatchmakingTicketParams
	ID                sql.NullInt64 `db:"id"`
}

func filterMatchParams(arg MatchParams, idColumnName *string) goqu.Expression {
	if idColumnName == nil {
		idColumnName = conversion.ValueToPointer("id")
	}
	expressions := []goqu.Expression{}
	if arg.ID.Valid {
		expressions = append(expressions, goqu.C(*idColumnName).Eq(arg.ID))
	}
	if arg.MatchmakingTicket.ID.Valid {
		expressions = append(expressions, goqu.C(*idColumnName).Eq(gq.From(gq.From("matchmaking_ticket").Where(goqu.Ex{"id": arg.MatchmakingTicket.ID}).Select("matchmaking_match_id").Limit(1))))
	}
	if arg.MatchmakingTicket.MatchmakingUser.ID.Valid {
		expressions = append(expressions, goqu.C(*idColumnName).Eq(gq.From(gq.From("matchmaking_ticket_with_user").Where(goqu.Ex{"matchmaking_user_id": arg.MatchmakingTicket.MatchmakingUser.ID}).Select("matchmaking_match_id").Limit(1))))
	}
	if arg.MatchmakingTicket.MatchmakingUser.ClientUserID.Valid {
		expressions = append(expressions, goqu.C(*idColumnName).Eq(gq.From(gq.From("matchmaking_ticket_with_user").Where(goqu.Ex{"client_user_id": arg.MatchmakingTicket.MatchmakingUser.ClientUserID}).Select("matchmaking_match_id").Limit(1))))
	}
	return goqu.And(expressions...)
}

type GetMatchParams struct {
	Match        MatchParams
	TicketLimit  uint64
	TicketOffset uint64
	UserLimit    uint64
	UserOffset   uint64
	ArenaLimit   uint64
	ArenaOffset  uint64
}

func filterGetMatchParams(arg GetMatchParams) goqu.Expression {
	return goqu.And(
		filterMatchParams(arg.Match, conversion.ValueToPointer("match_id")),
		goqu.C("ticket_number").Gt(arg.TicketOffset),
		goqu.C("ticket_number").Lte(arg.TicketOffset+arg.TicketLimit),
		goqu.C("user_number").Gt(arg.UserOffset),
		goqu.C("user_number").Lte(arg.UserOffset+arg.UserLimit),
		goqu.C("arena_number").Gt(arg.ArenaOffset),
		goqu.C("arena_number").Lte(arg.ArenaOffset+arg.ArenaLimit),
	)
}

func (q *Queries) GetMatch(ctx context.Context, arg GetMatchParams) ([]MatchmakingMatchWithArenaAndTicket, error) {
	matchmakingMatch := gq.From("matchmaking_match_with_arena_and_ticket").Prepared(true)
	query, args, err := matchmakingMatch.Where(filterGetMatchParams(arg)).ToSQL()
	if err != nil {
		return nil, err
	}
	rows, err := q.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []MatchmakingMatchWithArenaAndTicket
	for rows.Next() {
		var i MatchmakingMatchWithArenaAndTicket
		if err = rows.Scan(
			&i.MatchID,
			&i.PrivateServerID,
			&i.MatchStatus,
			&i.TicketCount,
			&i.UserCount,
			&i.MatchData,
			&i.LockedAt,
			&i.StartedAt,
			&i.EndedAt,
			&i.MatchCreatedAt,
			&i.MatchUpdatedAt,
			&i.ArenaID,
			&i.ArenaName,
			&i.ArenaMinPlayers,
			&i.ArenaMaxPlayersPerTicket,
			&i.ArenaMaxPlayers,
			&i.ArenaData,
			&i.ArenaCreatedAt,
			&i.ArenaUpdatedAt,
			&i.TicketID,
			&i.MatchmakingUserID,
			&i.TicketStatus,
			&i.TicketUserCount,
			&i.TicketNumber,
			&i.TicketData,
			&i.TicketCreatedAt,
			&i.TicketUpdatedAt,
			&i.ClientUserID,
			&i.Elo,
			&i.UserNumber,
			&i.UserData,
			&i.UserCreatedAt,
			&i.UserUpdatedAt,
			&i.TicketArenaID,
			&i.TicketArenaName,
			&i.TicketArenaMinPlayers,
			&i.TicketArenaMaxPlayersPerTicket,
			&i.TicketArenaMaxPlayers,
			&i.ArenaNumber,
			&i.TicketArenaData,
			&i.TicketArenaCreatedAt,
			&i.TicketArenaUpdatedAt,
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

type GetMatchesParams struct {
	Arena           ArenaParams
	MatchmakingUser MatchmakingUserParams
	Statuses        []string
	Limit           uint64
	Offset          uint64
	TicketLimit     uint64
	TicketOffset    uint64
	UserLimit       uint64
	UserOffset      uint64
	ArenaLimit      uint64
	ArenaOffset     uint64
}

func filterGetMatchesParams(arg GetMatchesParams) goqu.Expression {
	expressions := []goqu.Expression{}
	if arg.Arena.ID.Valid {
		expressions = append(expressions, goqu.C("arena_id").Eq(arg.Arena.ID))
	}
	if arg.Arena.Name.Valid {
		expressions = append(expressions, goqu.C("arena_name").Eq(arg.Arena.Name))
	}
	if arg.MatchmakingUser.ID.Valid {
		expressions = append(expressions, goqu.C("matchmaking_user_id").Eq(arg.MatchmakingUser.ID))
	}
	if arg.MatchmakingUser.ClientUserID.Valid {
		expressions = append(expressions, goqu.C("client_user_id").Eq(arg.MatchmakingUser.ClientUserID))
	}
	if len(arg.Statuses) > 0 {
		expressions = append(expressions, goqu.C("match_status").In(arg.Statuses))
	}
	finalExpression := goqu.And(
		goqu.C("match_id").In(gq.From(gq.From("matchmaking_match_with_arena_and_ticket").Select("match_id").Where(goqu.And(expressions...)).GroupBy("match_id").Order(goqu.C("match_id").Asc()).Limit(uint(arg.Limit)).Offset(uint(arg.Offset)))),
		goqu.C("ticket_number").Gt(arg.TicketOffset),
		goqu.C("ticket_number").Lte(arg.TicketOffset+arg.TicketLimit),
		goqu.C("user_number").Gt(arg.UserOffset),
		goqu.C("user_number").Lte(arg.UserOffset+arg.UserLimit),
		goqu.C("arena_number").Gt(arg.ArenaOffset),
		goqu.C("arena_number").Lte(arg.ArenaOffset+arg.ArenaLimit),
	)
	return finalExpression
}

func (q *Queries) GetMatches(ctx context.Context, arg GetMatchesParams) ([]MatchmakingMatchWithArenaAndTicket, error) {
	matchmakingMatch := gq.From("matchmaking_match_with_arena_and_ticket").Prepared(true)
	query, args, err := matchmakingMatch.Where(filterGetMatchesParams(arg)).ToSQL()
	if err != nil {
		return nil, err
	}
	rows, err := q.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []MatchmakingMatchWithArenaAndTicket
	for rows.Next() {
		var i MatchmakingMatchWithArenaAndTicket
		if err = rows.Scan(
			&i.MatchID,
			&i.PrivateServerID,
			&i.MatchStatus,
			&i.TicketCount,
			&i.UserCount,
			&i.MatchData,
			&i.LockedAt,
			&i.StartedAt,
			&i.EndedAt,
			&i.MatchCreatedAt,
			&i.MatchUpdatedAt,
			&i.ArenaID,
			&i.ArenaName,
			&i.ArenaMinPlayers,
			&i.ArenaMaxPlayersPerTicket,
			&i.ArenaMaxPlayers,
			&i.ArenaData,
			&i.ArenaCreatedAt,
			&i.ArenaUpdatedAt,
			&i.TicketID,
			&i.MatchmakingUserID,
			&i.TicketStatus,
			&i.TicketUserCount,
			&i.TicketNumber,
			&i.TicketData,
			&i.TicketCreatedAt,
			&i.TicketUpdatedAt,
			&i.ClientUserID,
			&i.Elo,
			&i.UserNumber,
			&i.UserData,
			&i.UserCreatedAt,
			&i.UserUpdatedAt,
			&i.TicketArenaID,
			&i.TicketArenaName,
			&i.TicketArenaMinPlayers,
			&i.TicketArenaMaxPlayersPerTicket,
			&i.TicketArenaMaxPlayers,
			&i.ArenaNumber,
			&i.TicketArenaData,
			&i.TicketArenaCreatedAt,
			&i.TicketArenaUpdatedAt,
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

type StartMatchParams struct {
	Match     MatchParams
	LockTime  time.Time
	StartTime time.Time
}

func (q *Queries) StartMatch(ctx context.Context, arg StartMatchParams) (sql.Result, error) {
	matchmakingMatch := gq.Update("matchmaking_match").Prepared(true)
	updates := goqu.Record{"locked_at": arg.LockTime, "started_at": arg.StartTime}
	matchmakingMatch = matchmakingMatch.Set(updates)
	query, args, err := matchmakingMatch.Where(
		goqu.And(
			filterMatchParams(arg.Match, nil),
			goqu.C("started_at").IsNull(),
		),
	).Limit(1).ToSQL()
	if err != nil {
		return nil, err
	}
	return q.db.ExecContext(ctx, query, args...)
}

type EndMatchParams struct {
	Match   MatchParams
	EndTime time.Time
}

func (q *Queries) EndMatch(ctx context.Context, arg EndMatchParams) (sql.Result, error) {
	matchmakingMatch := gq.Update("matchmaking_match").Prepared(true)
	updates := goqu.Record{"ended_at": arg.EndTime}
	matchmakingMatch = matchmakingMatch.Set(updates)
	query, args, err := matchmakingMatch.Where(
		goqu.And(
			filterMatchParams(arg.Match, nil),
			goqu.C("ended_at").IsNull(),
			goqu.C("started_at").IsNotNull(),
			goqu.C("started_at").Lt(arg.EndTime),
		),
	).Limit(1).ToSQL()
	if err != nil {
		return nil, err
	}
	return q.db.ExecContext(ctx, query, args...)
}

type UpdateMatchParams struct {
	Match MatchParams
	Data  json.RawMessage `db:"data"`
}

func (q *Queries) UpdateMatch(ctx context.Context, arg UpdateMatchParams) (sql.Result, error) {
	matchmakingMatch := gq.Update("matchmaking_match").Prepared(true)
	updates := goqu.Record{}
	if arg.Data != nil {
		updates["data"] = []byte(arg.Data)
	}
	matchmakingMatch = matchmakingMatch.Set(updates)
	query, args, err := matchmakingMatch.Where(filterMatchParams(arg.Match, nil)).Limit(1).ToSQL()
	if err != nil {
		return nil, err
	}
	return q.db.ExecContext(ctx, query, args...)
}

type SetMatchPrivateServerParams struct {
	Match           MatchParams
	PrivateServerID string `db:"private_server_id"`
}

func (q *Queries) SetMatchPrivateServer(ctx context.Context, arg SetMatchPrivateServerParams) (sql.Result, error) {
	matchmakingMatch := gq.Update("matchmaking_match").Prepared(true).Set(goqu.Record{"private_server_id": arg.PrivateServerID})
	query, args, err := matchmakingMatch.Where(
		filterMatchParams(arg.Match, nil),
		goqu.C("private_server_id").IsNull(),
	).Limit(1).ToSQL()
	if err != nil {
		return nil, err
	}
	return q.db.ExecContext(ctx, query, args...)
}

func (q *Queries) DeleteMatch(ctx context.Context, arg MatchParams) (sql.Result, error) {
	matchmakingMatch := gq.Delete("matchmaking_match").Prepared(true)
	query, args, err := matchmakingMatch.Where(filterMatchParams(arg, nil)).Limit(1).ToSQL()
	if err != nil {
		return nil, err
	}
	return q.db.ExecContext(ctx, query, args...)
}
