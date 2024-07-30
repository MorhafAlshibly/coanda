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

func (q *Queries) GetMatchmakingUser(ctx context.Context, arg GetMatchmakingUserParams) (MatchmakingUserWithElo, error) {
	matchmakingUser := gq.From("matchmaking_user_with_elo").Prepared(true)
	query, args, err := matchmakingUser.Where(filterGetMatchmakingUserParams(arg)).Limit(1).ToSQL()
	if err != nil {
		return MatchmakingUserWithElo{}, err
	}
	var i MatchmakingUserWithElo
	err = q.db.QueryRowContext(ctx, query, args...).Scan(
		&i.ID,
		&i.ClientUserID,
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

type MatchmakingTicketParams struct {
	MatchmakingUser GetMatchmakingUserParams
	ID              sql.NullInt64 `db:"id"`
}

type GetMatchmakingTicketParams struct {
	MatchmakingTicket MatchmakingTicketParams
	Limit             uint64
	Offset            uint64
}

func filterMatchmakingTicketParams(arg MatchmakingTicketParams) goqu.Expression {
	expressions := goqu.Ex{}
	if arg.ID.Valid {
		expressions["id"] = arg.ID
	}
	if arg.MatchmakingUser.ID.Valid {
		expressions["id"] = gq.From(gq.From("matchmaking_ticket_with_user").Where(goqu.And(goqu.Ex{"matchmaking_user_id": arg.MatchmakingUser.ID}, goqu.Or(goqu.Ex{"status": "PENDING"}, goqu.Ex{"status": "MATCHED"}))).Select("id").Limit(1))
	}
	if arg.MatchmakingUser.ClientUserID.Valid {
		expressions["id"] = gq.From(gq.From("matchmaking_ticket_with_user").Where(goqu.And(goqu.Ex{"client_user_id": arg.MatchmakingUser.ClientUserID}, goqu.Or(goqu.Ex{"status": "PENDING"}, goqu.Ex{"status": "MATCHED"}))).Select("id").Limit(1))
	}
	return expressions
}

func (q *Queries) GetMatchmakingTicket(ctx context.Context, arg GetMatchmakingTicketParams) ([]MatchmakingTicketWithUserAndArena, error) {
	matchmakingTicket := gq.From("matchmaking_ticket_with_user_and_arena").Prepared(true)
	query, args, err := matchmakingTicket.Where(filterMatchmakingTicketParams(arg.MatchmakingTicket)).Limit(uint(arg.Limit)).Offset(uint(arg.Offset)).ToSQL()
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

type PollMatchmakingTicketParams struct {
	MatchmakingTicket MatchmakingTicketParams
	ExpiryTimeWindow  time.Duration
}

func filterPollMatchmakingTicketParams(arg PollMatchmakingTicketParams) goqu.Expression {
	expressions := goqu.Ex{}
	if arg.MatchmakingTicket.ID.Valid {
		expressions["id"] = arg.MatchmakingTicket.ID
	}
	if arg.MatchmakingTicket.MatchmakingUser.ID.Valid {
		expressions["id"] = gq.From(gq.From("matchmaking_ticket_with_user").Where(goqu.And(goqu.Ex{"matchmaking_user_id": arg.MatchmakingTicket.MatchmakingUser.ID}, goqu.Ex{"status": "PENDING"})).Select("id").Limit(1))
	}
	if arg.MatchmakingTicket.MatchmakingUser.ClientUserID.Valid {
		expressions["id"] = gq.From(gq.From("matchmaking_ticket_with_user").Where(goqu.And(goqu.Ex{"client_user_id": arg.MatchmakingTicket.MatchmakingUser.ClientUserID}, goqu.Ex{"status": "PENDING"})).Select("id").Limit(1))
	}
	return expressions
}

func (q *Queries) PollMatchmakingTicket(ctx context.Context, arg PollMatchmakingTicketParams) (sql.Result, error) {
	matchmakingTicket := gq.Update("matchmaking_ticket").Prepared(true)
	updates := goqu.Record{"expires_at": time.Now().Add(arg.ExpiryTimeWindow)}
	matchmakingTicket = matchmakingTicket.Set(updates)
	query, args, err := matchmakingTicket.Where(filterPollMatchmakingTicketParams(arg)).Limit(1).ToSQL()
	if err != nil {
		return nil, err
	}
	return q.db.ExecContext(ctx, query, args...)
}

type GetMatchmakingTicketsParams struct {
	MatchmakingUser    GetMatchmakingUserParams
	MatchmakingMatchID sql.NullInt64  `db:"matchmaking_match_id"`
	Status             sql.NullString `db:"status"`
	Limit              uint64
	Offset             uint64
	UserLimit          uint64
	UserOffset         uint64
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
	if arg.Status.Valid {
		expressions["status"] = arg.Status
	}
	return expressions
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
	query, args, err := matchmakingTicket.Where(filterMatchmakingTicketParams(arg.MatchmakingTicket)).Limit(1).ToSQL()
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
		filterMatchmakingTicketParams(arg),
		goqu.Ex{"expires_at": goqu.Op{">": time.Now()}},
	)).Limit(1).ToSQL()
	if err != nil {
		return nil, err
	}
	return q.db.ExecContext(ctx, query, args...)
}
