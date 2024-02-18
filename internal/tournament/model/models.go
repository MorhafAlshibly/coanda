// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.24.0

package model

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"
)

type TournamentTournamentInterval string

const (
	TournamentTournamentIntervalDaily     TournamentTournamentInterval = "daily"
	TournamentTournamentIntervalWeekly    TournamentTournamentInterval = "weekly"
	TournamentTournamentIntervalMonthly   TournamentTournamentInterval = "monthly"
	TournamentTournamentIntervalUnlimited TournamentTournamentInterval = "unlimited"
)

func (e *TournamentTournamentInterval) Scan(src interface{}) error {
	switch s := src.(type) {
	case []byte:
		*e = TournamentTournamentInterval(s)
	case string:
		*e = TournamentTournamentInterval(s)
	default:
		return fmt.Errorf("unsupported scan type for TournamentTournamentInterval: %T", src)
	}
	return nil
}

type NullTournamentTournamentInterval struct {
	TournamentTournamentInterval TournamentTournamentInterval
	Valid                        bool // Valid is true if TournamentTournamentInterval is not NULL
}

// Scan implements the Scanner interface.
func (ns *NullTournamentTournamentInterval) Scan(value interface{}) error {
	if value == nil {
		ns.TournamentTournamentInterval, ns.Valid = "", false
		return nil
	}
	ns.Valid = true
	return ns.TournamentTournamentInterval.Scan(value)
}

// Value implements the driver Valuer interface.
func (ns NullTournamentTournamentInterval) Value() (driver.Value, error) {
	if !ns.Valid {
		return nil, nil
	}
	return string(ns.TournamentTournamentInterval), nil
}

type RankedTournament struct {
	ID                  uint64                       `db:"id"`
	Name                string                       `db:"name"`
	TournamentInterval  TournamentTournamentInterval `db:"tournament_interval"`
	UserID              uint64                       `db:"user_id"`
	Score               int64                        `db:"score"`
	Ranking             uint64                       `db:"ranking"`
	Data                json.RawMessage              `db:"data"`
	TournamentStartedAt time.Time                    `db:"tournament_started_at"`
	CreatedAt           time.Time                    `db:"created_at"`
	UpdatedAt           time.Time                    `db:"updated_at"`
}

type Tournament struct {
	ID                  uint64                       `db:"id"`
	Name                string                       `db:"name"`
	TournamentInterval  TournamentTournamentInterval `db:"tournament_interval"`
	UserID              uint64                       `db:"user_id"`
	Score               int64                        `db:"score"`
	Data                json.RawMessage              `db:"data"`
	TournamentStartedAt time.Time                    `db:"tournament_started_at"`
	CreatedAt           time.Time                    `db:"created_at"`
	UpdatedAt           time.Time                    `db:"updated_at"`
}