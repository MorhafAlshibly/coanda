package model

import (
	"context"
	"database/sql"
	"encoding/json"
	"time"

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
	Arena               GetArenaParams
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

type GetMatchmakingUserParams struct {
	ID           sql.NullInt64 `db:"id"`
	ClientUserID sql.NullInt64 `db:"user_id"`
}

func filterGetMatchmakingUserParams(arg GetMatchmakingUserParams) goqu.Expression {
	expressions := goqu.Ex{}
	if arg.ID.Valid {
		expressions["id"] = arg.ID
	}
	if arg.ClientUserID.Valid {
		expressions["user_id"] = arg.ClientUserID
	}
	return expressions
}

func (q *Queries) GetMatchmakingUser(ctx context.Context, arg GetMatchmakingUserParams) (MatchmakingUser, error) {
	matchmakingUser := gq.From("matchmaking_user").Prepared(true)
	query, args, err := matchmakingUser.Where(filterGetMatchmakingUserParams(arg)).Limit(1).ToSQL()
	if err != nil {
		return MatchmakingUser{}, err
	}
	var i MatchmakingUser
	err = q.db.QueryRowContext(ctx, query, args...).Scan(
		&i.ID,
		&i.ClientUserID,
		&i.Elo,
		&i.Data,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

type UpdateMatchmakingUserParams struct {
	MatchmakingUser GetMatchmakingUserParams
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
	query, args, err := matchmakingUser.Where(filterGetMatchmakingUserParams(arg.MatchmakingUser)).Limit(1).ToSQL()
	if err != nil {
		return nil, err
	}
	return q.db.ExecContext(ctx, query, args...)
}

type MatchmakingTicketParams struct {
	MatchmakingUser GetMatchmakingUserParams
	ID              sql.NullInt64 `db:"id"`
	Statuses        []string
}

type GetMatchmakingTicketParams struct {
	MatchmakingTicket MatchmakingTicketParams
	UserLimit         uint64
	UserOffset        uint64
	ArenaLimit        uint64
	ArenaOffset       uint64
}

func filterGetMatchmakingTicketParams(arg GetMatchmakingTicketParams) goqu.Expression {
	expressions := goqu.Ex{}
	if arg.MatchmakingTicket.ID.Valid {
		expressions["ticket_id"] = arg.MatchmakingTicket.ID
	}
	if arg.MatchmakingTicket.MatchmakingUser.ID.Valid {
		expressions["ticket_id"] = gq.From(gq.From("matchmaking_ticket_with_user").Where(goqu.And(goqu.Ex{"matchmaking_user_id": arg.MatchmakingTicket.MatchmakingUser.ID})).Select("ticket_id").Limit(1))
	}
	if arg.MatchmakingTicket.MatchmakingUser.ClientUserID.Valid {
		expressions["ticket_id"] = gq.From(gq.From("matchmaking_ticket_with_user").Where(goqu.And(goqu.Ex{"client_user_id": arg.MatchmakingTicket.MatchmakingUser.ClientUserID})).Select("ticket_id").Limit(1))
	}
	if len(arg.MatchmakingTicket.Statuses) > 0 {
		expressions["status"] = goqu.Op{"IN": arg.MatchmakingTicket.Statuses}
	}
	orExpression := goqu.And(
		goqu.C("user_number").Gt(arg.UserOffset),
		goqu.C("user_number").Lt(arg.UserOffset+arg.UserLimit),
		goqu.C("arena_number").Gt(arg.ArenaOffset),
		goqu.C("arena_number").Lt(arg.ArenaOffset+arg.ArenaLimit),
	)
	return goqu.And(expressions, orExpression)
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
		if err = q.db.QueryRowContext(ctx, query, args...).Scan(
			&i.TicketID,
			&i.MatchmakingMatchID,
			&i.Status,
			&i.TicketData,
			&i.ExpiresAt,
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
	return items, nil
}

type PollMatchmakingTicketParams struct {
	MatchmakingTicket MatchmakingTicketParams
	ExpiryTimeWindow  time.Duration
}

func filterMatchmakingTicketParams(arg MatchmakingTicketParams) goqu.Expression {
	expressions := goqu.Ex{}
	if arg.ID.Valid {
		expressions["id"] = arg.ID
	}
	if arg.MatchmakingUser.ID.Valid {
		expressions["id"] = gq.From(gq.From("matchmaking_ticket_with_user").Where(goqu.And(goqu.Ex{"matchmaking_user_id": arg.MatchmakingUser.ID})).Select("id").Limit(1))
	}
	if arg.MatchmakingUser.ClientUserID.Valid {
		expressions["id"] = gq.From(gq.From("matchmaking_ticket_with_user").Where(goqu.And(goqu.Ex{"client_user_id": arg.MatchmakingUser.ClientUserID})).Select("id").Limit(1))
	}
	return goqu.And(expressions, goqu.C("expires_at").Gt(time.Now()))
}

func (q *Queries) PollMatchmakingTicket(ctx context.Context, arg PollMatchmakingTicketParams) (sql.Result, error) {
	matchmakingTicket := gq.Update("matchmaking_ticket").Prepared(true)
	updates := goqu.Record{"expires_at": time.Now().Add(arg.ExpiryTimeWindow)}
	matchmakingTicket = matchmakingTicket.Set(updates)
	query, args, err := matchmakingTicket.Where(filterMatchmakingTicketParams(arg.MatchmakingTicket)).Limit(1).ToSQL()
	if err != nil {
		return nil, err
	}
	return q.db.ExecContext(ctx, query, args...)
}

type GetMatchmakingTicketsParams struct {
	MatchmakingUser    GetMatchmakingUserParams
	MatchmakingMatchID sql.NullInt64 `db:"matchmaking_match_id"`
	Statuses           []string
	Limit              uint64
	Offset             uint64
	UserLimit          uint64
	UserOffset         uint64
	ArenaLimit         uint64
	ArenaOffset        uint64
}

func filterGetMatchmakingTicketsParams(arg GetMatchmakingTicketsParams) goqu.Expression {
	expressions := goqu.Ex{}
	if arg.MatchmakingMatchID.Valid {
		expressions["matchmaking_match_id"] = arg.MatchmakingMatchID
	}
	if arg.MatchmakingUser.ID.Valid {
		expressions["id"] = goqu.Op{"IN": gq.From(gq.From("matchmaking_ticket_with_user").Where(goqu.Ex{"matchmaking_user_id": arg.MatchmakingUser.ID}).Select("id").Limit(1))}
	}
	if arg.MatchmakingUser.ClientUserID.Valid {
		expressions["id"] = goqu.Op{"IN": gq.From(gq.From("matchmaking_ticket_with_user").Where(goqu.Ex{"client_user_id": arg.MatchmakingUser.ClientUserID}).Select("id").Limit(1))}
	}
	if len(arg.Statuses) > 0 {
		expressions["status"] = goqu.Op{"IN": arg.Statuses}
	}
	orExpression := goqu.And(
		goqu.C("user_number").Gt(arg.UserOffset),
		goqu.C("user_number").Lt(arg.UserOffset+arg.UserLimit),
		goqu.C("arena_number").Gt(arg.ArenaOffset),
		goqu.C("arena_number").Lt(arg.ArenaOffset+arg.ArenaLimit),
		goqu.C("ticket_id").Gt(arg.Offset),
		goqu.C("ticket_id").Lt(arg.Offset+arg.Limit),
	)
	return goqu.And(expressions, orExpression)
}

func (q *Queries) GetMatchmakingTickets(ctx context.Context, arg GetMatchmakingTicketsParams) ([]MatchmakingTicketWithUserAndArena, error) {
	matchmakingTicket := gq.From("matchmaking_ticket_with_user_and_arena").Prepared(true)
	query, args, err := matchmakingTicket.Where(filterGetMatchmakingTicketsParams(arg)).Limit(uint(arg.Limit)).Offset(uint(arg.Offset)).ToSQL()
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
		if err = q.db.QueryRowContext(ctx, query, args...).Scan(
			&i.ID,
			&i.MatchmakingUserID,
			&i.ClientUserID,
			&i.Elos,
			&i.UserData,
			&i.UserCreatedAt,
			&i.UserUpdatedAt,
			&i.Arenas,
			&i.MatchmakingMatchID,
			&i.Status,
			&i.TicketData,
			&i.ExpiresAt,
			&i.TicketCreatedAt,
			&i.TicketUpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
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
	query, args, err := matchmakingTicket.Where(filterGetMatchmakingTicketParams(arg.MatchmakingTicket)).Limit(1).ToSQL()
	if err != nil {
		return nil, err
	}
	return q.db.ExecContext(ctx, query, args...)
}

func (q *Queries) ExpireMatchmakingTicket(ctx context.Context, arg MatchmakingTicketParams) (sql.Result, error) {
	matchmakingTicket := gq.Update("matchmaking_ticket").Prepared(true)
	updates := goqu.Record{"expires_at": time.Now()}
	matchmakingTicket = matchmakingTicket.Set(updates)
	// Only expire if the expires_at is in the future
	query, args, err := matchmakingTicket.Where(goqu.And(
		filterGetMatchmakingTicketParams(arg),
		goqu.Ex{"expires_at": goqu.Op{">": time.Now()}},
	)).Limit(1).ToSQL()
	if err != nil {
		return nil, err
	}
	return q.db.ExecContext(ctx, query, args...)
}

type MatchParams struct {
	MatchmakingTicket MatchmakingTicketParams
	ID                sql.NullInt64 `db:"id"`
}

type GetMatchParams struct {
	Match  MatchParams
	Limit  uint64
	Offset uint64
}

func filterMatchParams(arg MatchParams) goqu.Expression {
	expressions := goqu.Ex{}
	if arg.ID.Valid {
		expressions["id"] = arg.ID
	}
	if arg.MatchmakingTicket.ID.Valid {
		expressions["id"] = gq.From(gq.From("matchmaking_ticket").Where(goqu.Ex{"id": arg.MatchmakingTicket.ID}).Select("matchmaking_match_id").Limit(1))
	}
	if arg.MatchmakingTicket.MatchmakingUser.ID.Valid {
		expressions["id"] = gq.From(gq.From("matchmaking_ticket_with_user").Where(goqu.Ex{"matchmaking_user_id": arg.MatchmakingTicket.MatchmakingUser.ID}).Select("matchmaking_match_id").Limit(1))
	}
	if arg.MatchmakingTicket.MatchmakingUser.ClientUserID.Valid {
		expressions["id"] = gq.From(gq.From("matchmaking_ticket_with_user").Where(goqu.Ex{"client_user_id": arg.MatchmakingTicket.MatchmakingUser.ClientUserID}).Select("matchmaking_match_id").Limit(1))
	}
	return expressions
}

func (q *Queries) GetMatch(ctx context.Context, arg GetMatchParams) ([]MatchmakingMatchWithTicket, error) {
	matchmakingMatch := gq.From("matchmaking_match_with_ticket").Prepared(true)
	query, args, err := matchmakingMatch.Where(filterMatchParams(arg.Match)).Limit(uint(arg.Limit)).Offset(uint(arg.Offset)).ToSQL()
	if err != nil {
		return nil, err
	}
	rows, err := q.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []MatchmakingMatchWithTicket
	for rows.Next() {
		var i MatchmakingMatchWithTicket
		if err = q.db.QueryRowContext(ctx, query, args...).Scan(
			&i.ID,
			&i.ArenaID,
			&i.ArenaName,
			&i.ArenaMinPlayers,
			&i.ArenaMaxPlayersPerTicket,
			&i.ArenaMaxPlayers,
			&i.ArenaData,
			&i.ArenaCreatedAt,
			&i.ArenaUpdatedAt,
			&i.PrivateServerID,
			&i.MatchStatus,
			&i.MatchData,
			&i.LockedAt,
			&i.StartedAt,
			&i.EndedAt,
			&i.MatchCreatedAt,
			&i.MatchUpdatedAt,
			&i.MatchmakingTicketID,
			&i.MatchmakingUserID,
			&i.ClientUserID,
			&i.Elos,
			&i.UserData,
			&i.UserCreatedAt,
			&i.UserUpdatedAt,
			&i.Arenas,
			&i.MatchmakingMatchID,
			&i.TicketStatus,
			&i.TicketData,
			&i.ExpiresAt,
			&i.TicketCreatedAt,
			&i.TicketUpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	return items, nil
}

type GetMatchesParams struct {
	Arena           GetArenaParams
	MatchmakingUser GetMatchmakingUserParams
	Status          sql.NullString `db:"status"`
	Limit           uint64
	Offset          uint64
}

func filterGetMatchesParams(arg GetMatchesParams) goqu.Expression {
	expressions := goqu.Ex{}
	if arg.Arena.ID.Valid {
		expressions["arena_id"] = arg.Arena.ID
	}
	if arg.Arena.Name.Valid {
		expressions["arena_name"] = arg.Arena.Name
	}
	if arg.MatchmakingUser.ID.Valid {
		expressions["id"] = gq.From(gq.From("matchmaking_match_with_ticket").Where(goqu.Ex{"matchmaking_user_id": arg.MatchmakingUser.ID}).Select("id").Limit(1))
	}
	if arg.MatchmakingUser.ClientUserID.Valid {
		expressions["id"] = gq.From(gq.From("matchmaking_match_with_ticket").Where(goqu.Ex{"client_user_id": arg.MatchmakingUser.ClientUserID}).Select("id").Limit(1))
	}
	if arg.Status.Valid {
		expressions["status"] = arg.Status
	}
	return expressions
}

func (q *Queries) GetMatches(ctx context.Context, arg GetMatchesParams) ([]MatchmakingMatchWithTicket, error) {
	matchmakingMatch := gq.From("matchmaking_match_with_ticket").Prepared(true)
	query, args, err := matchmakingMatch.Where(filterGetMatchesParams(arg)).Limit(uint(arg.Limit)).Offset(uint(arg.Offset)).ToSQL()
	if err != nil {
		return nil, err
	}
	rows, err := q.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []MatchmakingMatchWithTicket
	for rows.Next() {
		var i MatchmakingMatchWithTicket
		if err = q.db.QueryRowContext(ctx, query, args...).Scan(
			&i.ID,
			&i.ArenaID,
			&i.ArenaName,
			&i.ArenaMinPlayers,
			&i.ArenaMaxPlayersPerTicket,
			&i.ArenaMaxPlayers,
			&i.ArenaData,
			&i.ArenaCreatedAt,
			&i.ArenaUpdatedAt,
			&i.PrivateServerID,
			&i.MatchStatus,
			&i.MatchData,
			&i.LockedAt,
			&i.StartedAt,
			&i.EndedAt,
			&i.MatchCreatedAt,
			&i.MatchUpdatedAt,
			&i.MatchmakingTicketID,
			&i.MatchmakingUserID,
			&i.ClientUserID,
			&i.Elos,
			&i.UserData,
			&i.UserCreatedAt,
			&i.UserUpdatedAt,
			&i.Arenas,
			&i.MatchmakingMatchID,
			&i.TicketStatus,
			&i.TicketData,
			&i.ExpiresAt,
			&i.TicketCreatedAt,
			&i.TicketUpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
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
			filterGetMatchmakingTicketParams(arg.Match.MatchmakingTicket),
			goqu.Ex{"started_at": goqu.Op{"IS": nil}},
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
			filterMatchParams(arg.Match),
			goqu.Ex{"ended_at": goqu.Op{"IS": nil}},
			goqu.Ex{"started_at": goqu.Op{"IS NOT": nil}},
			goqu.Ex{"started_at": goqu.Op{"<": arg.EndTime}},
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
	query, args, err := matchmakingMatch.Where(filterMatchParams(arg.Match)).Limit(1).ToSQL()
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
		filterMatchParams(arg.Match),
		goqu.Ex{"private_server_id": nil},
	).Limit(1).ToSQL()
	if err != nil {
		return nil, err
	}
	return q.db.ExecContext(ctx, query, args...)
}
