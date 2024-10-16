// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0

package model

import (
	"database/sql"
	"encoding/json"
	"time"
)

type MatchmakingArena struct {
	ID                  uint64          `db:"id"`
	Name                string          `db:"name"`
	MinPlayers          uint8           `db:"min_players"`
	MaxPlayersPerTicket uint8           `db:"max_players_per_ticket"`
	MaxPlayers          uint8           `db:"max_players"`
	Data                json.RawMessage `db:"data"`
	CreatedAt           time.Time       `db:"created_at"`
	UpdatedAt           time.Time       `db:"updated_at"`
}

type MatchmakingMatch struct {
	ID                 uint64          `db:"id"`
	MatchmakingArenaID uint64          `db:"matchmaking_arena_id"`
	PrivateServerID    sql.NullString  `db:"private_server_id"`
	Data               json.RawMessage `db:"data"`
	LockedAt           sql.NullTime    `db:"locked_at"`
	StartedAt          sql.NullTime    `db:"started_at"`
	EndedAt            sql.NullTime    `db:"ended_at"`
	CreatedAt          time.Time       `db:"created_at"`
	UpdatedAt          time.Time       `db:"updated_at"`
}

type MatchmakingMatchWithTicket struct {
	ID                       uint64          `db:"id"`
	ArenaID                  sql.NullInt64   `db:"arena_id"`
	ArenaName                sql.NullString  `db:"arena_name"`
	ArenaMinPlayers          sql.NullInt16   `db:"arena_min_players"`
	ArenaMaxPlayersPerTicket sql.NullInt16   `db:"arena_max_players_per_ticket"`
	ArenaMaxPlayers          sql.NullInt16   `db:"arena_max_players"`
	ArenaData                json.RawMessage `db:"arena_data"`
	ArenaCreatedAt           sql.NullTime    `db:"arena_created_at"`
	ArenaUpdatedAt           sql.NullTime    `db:"arena_updated_at"`
	PrivateServerID          sql.NullString  `db:"private_server_id"`
	MatchStatus              string          `db:"match_status"`
	MatchData                json.RawMessage `db:"match_data"`
	LockedAt                 sql.NullTime    `db:"locked_at"`
	StartedAt                sql.NullTime    `db:"started_at"`
	EndedAt                  sql.NullTime    `db:"ended_at"`
	MatchCreatedAt           time.Time       `db:"match_created_at"`
	MatchUpdatedAt           time.Time       `db:"match_updated_at"`
	MatchmakingTicketID      sql.NullInt64   `db:"matchmaking_ticket_id"`
	MatchmakingUserID        sql.NullInt64   `db:"matchmaking_user_id"`
	ClientUserID             sql.NullInt64   `db:"client_user_id"`
	Elos                     json.RawMessage `db:"elos"`
	UserData                 json.RawMessage `db:"user_data"`
	UserCreatedAt            sql.NullTime    `db:"user_created_at"`
	UserUpdatedAt            sql.NullTime    `db:"user_updated_at"`
	Arenas                   json.RawMessage `db:"arenas"`
	MatchmakingMatchID       sql.NullInt64   `db:"matchmaking_match_id"`
	TicketStatus             sql.NullString  `db:"ticket_status"`
	TicketData               json.RawMessage `db:"ticket_data"`
	ExpiresAt                sql.NullTime    `db:"expires_at"`
	TicketCreatedAt          sql.NullTime    `db:"ticket_created_at"`
	TicketUpdatedAt          sql.NullTime    `db:"ticket_updated_at"`
}

type MatchmakingTicket struct {
	ID                 uint64          `db:"id"`
	MatchmakingMatchID sql.NullInt64   `db:"matchmaking_match_id"`
	EloWindow          uint32          `db:"elo_window"`
	Data               json.RawMessage `db:"data"`
	ExpiresAt          time.Time       `db:"expires_at"`
	CreatedAt          time.Time       `db:"created_at"`
	UpdatedAt          time.Time       `db:"updated_at"`
}

type MatchmakingTicketArena struct {
	MatchmakingTicketID uint64 `db:"matchmaking_ticket_id"`
	MatchmakingArenaID  uint64 `db:"matchmaking_arena_id"`
}

type MatchmakingTicketUser struct {
	MatchmakingTicketID uint64 `db:"matchmaking_ticket_id"`
	MatchmakingUserID   uint64 `db:"matchmaking_user_id"`
}

type MatchmakingTicketWithUser struct {
	ID                 uint64          `db:"id"`
	MatchmakingUserID  uint64          `db:"matchmaking_user_id"`
	ClientUserID       uint64          `db:"client_user_id"`
	UserData           json.RawMessage `db:"user_data"`
	UserCreatedAt      time.Time       `db:"user_created_at"`
	UserUpdatedAt      time.Time       `db:"user_updated_at"`
	MatchmakingMatchID sql.NullInt64   `db:"matchmaking_match_id"`
	Status             string          `db:"status"`
	TicketData         json.RawMessage `db:"ticket_data"`
	ExpiresAt          time.Time       `db:"expires_at"`
	TicketCreatedAt    time.Time       `db:"ticket_created_at"`
	TicketUpdatedAt    time.Time       `db:"ticket_updated_at"`
}

type MatchmakingTicketWithUserAndArena struct {
	ID                 uint64          `db:"id"`
	MatchmakingUserID  uint64          `db:"matchmaking_user_id"`
	ClientUserID       uint64          `db:"client_user_id"`
	Elos               json.RawMessage `db:"elos"`
	UserData           json.RawMessage `db:"user_data"`
	UserCreatedAt      time.Time       `db:"user_created_at"`
	UserUpdatedAt      time.Time       `db:"user_updated_at"`
	Arenas             json.RawMessage `db:"arenas"`
	MatchmakingMatchID sql.NullInt64   `db:"matchmaking_match_id"`
	Status             string          `db:"status"`
	TicketData         json.RawMessage `db:"ticket_data"`
	ExpiresAt          time.Time       `db:"expires_at"`
	TicketCreatedAt    time.Time       `db:"ticket_created_at"`
	TicketUpdatedAt    time.Time       `db:"ticket_updated_at"`
}

type MatchmakingUser struct {
	ID           uint64          `db:"id"`
	ClientUserID uint64          `db:"client_user_id"`
	Data         json.RawMessage `db:"data"`
	CreatedAt    time.Time       `db:"created_at"`
	UpdatedAt    time.Time       `db:"updated_at"`
}

type MatchmakingUserElo struct {
	ID                 uint64 `db:"id"`
	Elo                int32  `db:"elo"`
	MatchmakingUserID  uint64 `db:"matchmaking_user_id"`
	MatchmakingArenaID uint64 `db:"matchmaking_arena_id"`
}

type MatchmakingUserWithElo struct {
	ID           uint64          `db:"id"`
	ClientUserID uint64          `db:"client_user_id"`
	Elos         json.RawMessage `db:"elos"`
	Data         json.RawMessage `db:"data"`
	CreatedAt    time.Time       `db:"created_at"`
	UpdatedAt    time.Time       `db:"updated_at"`
}
