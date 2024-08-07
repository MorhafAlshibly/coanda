// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.24.0
// source: queries.sql

package model

import (
	"context"
	"database/sql"
	"encoding/json"
)

const CreateTeam = `-- name: CreateTeam :execresult
INSERT INTO team (name, owner, score, data)
VALUES (?, ?, ?, ?)
`

type CreateTeamParams struct {
	Name  string          `db:"name"`
	Owner uint64          `db:"owner"`
	Score int64           `db:"score"`
	Data  json.RawMessage `db:"data"`
}

func (q *Queries) CreateTeam(ctx context.Context, arg CreateTeamParams) (sql.Result, error) {
	return q.db.ExecContext(ctx, CreateTeam,
		arg.Name,
		arg.Owner,
		arg.Score,
		arg.Data,
	)
}

const CreateTeamMember = `-- name: CreateTeamMember :execresult
INSERT INTO team_member (team, user_id, member_number, data)
VALUES (?, ?, ?, ?)
`

type CreateTeamMemberParams struct {
	Team         string          `db:"team"`
	UserID       uint64          `db:"user_id"`
	MemberNumber uint32          `db:"member_number"`
	Data         json.RawMessage `db:"data"`
}

func (q *Queries) CreateTeamMember(ctx context.Context, arg CreateTeamMemberParams) (sql.Result, error) {
	return q.db.ExecContext(ctx, CreateTeamMember,
		arg.Team,
		arg.UserID,
		arg.MemberNumber,
		arg.Data,
	)
}

const CreateTeamOwner = `-- name: CreateTeamOwner :execresult
INSERT INTO team_owner (team, user_id)
VALUES (?, ?)
`

type CreateTeamOwnerParams struct {
	Team   string `db:"team"`
	UserID uint64 `db:"user_id"`
}

func (q *Queries) CreateTeamOwner(ctx context.Context, arg CreateTeamOwnerParams) (sql.Result, error) {
	return q.db.ExecContext(ctx, CreateTeamOwner, arg.Team, arg.UserID)
}

const DeleteTeamMember = `-- name: DeleteTeamMember :execresult
DELETE FROM team_member
WHERE user_id = ?
LIMIT 1
`

func (q *Queries) DeleteTeamMember(ctx context.Context, userID uint64) (sql.Result, error) {
	return q.db.ExecContext(ctx, DeleteTeamMember, userID)
}

const GetHighestMemberNumber = `-- name: GetHighestMemberNumber :one
SELECT MAX(member_number) AS member_number
FROM team_member
WHERE team = ?
`

func (q *Queries) GetHighestMemberNumber(ctx context.Context, team string) (interface{}, error) {
	row := q.db.QueryRowContext(ctx, GetHighestMemberNumber, team)
	var member_number interface{}
	err := row.Scan(&member_number)
	return member_number, err
}

const GetTeamMember = `-- name: GetTeamMember :one
SELECT team,
  user_id,
  member_number,
  data,
  joined_at,
  updated_at
FROM team_member
WHERE user_id = ?
LIMIT 1
`

func (q *Queries) GetTeamMember(ctx context.Context, member uint64) (TeamMember, error) {
	row := q.db.QueryRowContext(ctx, GetTeamMember, member)
	var i TeamMember
	err := row.Scan(
		&i.Team,
		&i.UserID,
		&i.MemberNumber,
		&i.Data,
		&i.JoinedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const GetTeams = `-- name: GetTeams :many
SELECT name, owner, score, ranking, data, created_at, updated_at
FROM ranked_team
ORDER BY score DESC
LIMIT ? OFFSET ?
`

type GetTeamsParams struct {
	Limit  int32 `db:"limit"`
	Offset int32 `db:"offset"`
}

func (q *Queries) GetTeams(ctx context.Context, arg GetTeamsParams) ([]RankedTeam, error) {
	rows, err := q.db.QueryContext(ctx, GetTeams, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []RankedTeam
	for rows.Next() {
		var i RankedTeam
		if err := rows.Scan(
			&i.Name,
			&i.Owner,
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
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const SearchTeams = `-- name: SearchTeams :many
SELECT name, owner, score, ranking, data, created_at, updated_at
FROM ranked_team
WHERE name LIKE CONCAT('%', ?, '%')
ORDER BY score DESC
LIMIT ? OFFSET ?
`

type SearchTeamsParams struct {
	Query  interface{} `db:"query"`
	Limit  int32       `db:"limit"`
	Offset int32       `db:"offset"`
}

func (q *Queries) SearchTeams(ctx context.Context, arg SearchTeamsParams) ([]RankedTeam, error) {
	rows, err := q.db.QueryContext(ctx, SearchTeams, arg.Query, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []RankedTeam
	for rows.Next() {
		var i RankedTeam
		if err := rows.Scan(
			&i.Name,
			&i.Owner,
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
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const UpdateTeamMember = `-- name: UpdateTeamMember :execresult
UPDATE team_member
SET data = ?
WHERE user_id = ?
LIMIT 1
`

type UpdateTeamMemberParams struct {
	Data   json.RawMessage `db:"data"`
	UserID uint64          `db:"user_id"`
}

func (q *Queries) UpdateTeamMember(ctx context.Context, arg UpdateTeamMemberParams) (sql.Result, error) {
	return q.db.ExecContext(ctx, UpdateTeamMember, arg.Data, arg.UserID)
}
