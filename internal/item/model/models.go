// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0

package model

import (
	"database/sql"
	"encoding/json"
	"time"
)

type Item struct {
	ID        string          `db:"id"`
	Type      string          `db:"type"`
	Data      json.RawMessage `db:"data"`
	ExpiresAt sql.NullTime    `db:"expires_at"`
	CreatedAt time.Time       `db:"created_at"`
	UpdatedAt time.Time       `db:"updated_at"`
}
