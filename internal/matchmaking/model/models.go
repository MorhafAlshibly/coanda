// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.24.0

package model

import (
	"database/sql"
	"encoding/json"
	"time"
)

type MatchmakingArena struct {
	ID                  uint64          `db:"id"`
	Name                string          `db:"name"`
	MinPlayers          uint32          `db:"min_players"`
	MaxPlayersPerTicket uint32          `db:"max_players_per_ticket"`
	MaxPlayers          uint32          `db:"max_players"`
	Data                json.RawMessage `db:"data"`
	CreatedAt           time.Time       `db:"created_at"`
	UpdatedAt           time.Time       `db:"updated_at"`
}

type MatchmakingMatch struct {
	ID                 uint64          `db:"id"`
	MatchmakingArenaID uint64          `db:"matchmaking_arena_id"`
	Data               json.RawMessage `db:"data"`
	LockedAt           sql.NullTime    `db:"locked_at"`
	StartedAt          sql.NullTime    `db:"started_at"`
	EndedAt            sql.NullTime    `db:"ended_at"`
	CreatedAt          time.Time       `db:"created_at"`
	UpdatedAt          time.Time       `db:"updated_at"`
}

type MatchmakingTicket struct {
	ID                 uint64          `db:"id"`
	MatchmakingMatchID sql.NullInt64   `db:"matchmaking_match_id"`
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
