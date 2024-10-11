// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: queries.sql

package model

import (
	"context"
	"database/sql"
	"encoding/json"
)

const CreateTeam = `-- name: CreateTeam :execresult
INSERT INTO team (name, score, data)
VALUES (?, ?, ?)
`

type CreateTeamParams struct {
	Name  string          `db:"name"`
	Score int64           `db:"score"`
	Data  json.RawMessage `db:"data"`
}

func (q *Queries) CreateTeam(ctx context.Context, arg CreateTeamParams) (sql.Result, error) {
	return q.db.ExecContext(ctx, CreateTeam, arg.Name, arg.Score, arg.Data)
}

const CreateTeamMember = `-- name: CreateTeamMember :execresult
INSERT INTO team_member (user_id, team_id, member_number, data)
VALUES (?, ?, ?, ?)
`

type CreateTeamMemberParams struct {
	UserID       uint64          `db:"user_id"`
	TeamID       uint64          `db:"team_id"`
	MemberNumber uint32          `db:"member_number"`
	Data         json.RawMessage `db:"data"`
}

func (q *Queries) CreateTeamMember(ctx context.Context, arg CreateTeamMemberParams) (sql.Result, error) {
	return q.db.ExecContext(ctx, CreateTeamMember,
		arg.UserID,
		arg.TeamID,
		arg.MemberNumber,
		arg.Data,
	)
}

const GetFirstOpenMemberNumber = `-- name: GetFirstOpenMemberNumber :one
SELECT first_open_member
FROM team_with_first_open_member
WHERE id = ?
LIMIT 1
`

func (q *Queries) GetFirstOpenMemberNumber(ctx context.Context, team uint64) (uint32, error) {
	row := q.db.QueryRowContext(ctx, GetFirstOpenMemberNumber, team)
	var first_open_member uint32
	err := row.Scan(&first_open_member)
	return first_open_member, err
}

const GetTeams = `-- name: GetTeams :many
SELECT id,
  name,
  score,
  ranking,
  data,
  created_at,
  updated_at,
  member_id,
  user_id,
  member_number,
  member_data,
  joined_at,
  member_updated_at,
  member_number_without_gaps
FROM ranked_team_with_member
WHERE member_number_without_gaps <= CAST(? AS UNSIGNED)
  AND member_number_without_gaps > CAST(? AS UNSIGNED)
  AND id IN (
    SELECT id
    FROM (
        SELECT id
        FROM team
        ORDER BY score DESC,
          id
        LIMIT ? OFFSET ?
      ) temp_team
  )
ORDER BY score DESC,
  id,
  member_number
`

type GetTeamsParams struct {
	MemberLimitPlusOffset int64 `db:"member_limit_plus_offset"`
	MemberOffset          int64 `db:"member_offset"`
	Limit                 int32 `db:"limit"`
	Offset                int32 `db:"offset"`
}

func (q *Queries) GetTeams(ctx context.Context, arg GetTeamsParams) ([]RankedTeamWithMember, error) {
	rows, err := q.db.QueryContext(ctx, GetTeams,
		arg.MemberLimitPlusOffset,
		arg.MemberOffset,
		arg.Limit,
		arg.Offset,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []RankedTeamWithMember
	for rows.Next() {
		var i RankedTeamWithMember
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.Score,
			&i.Ranking,
			&i.Data,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.MemberID,
			&i.UserID,
			&i.MemberNumber,
			&i.MemberData,
			&i.JoinedAt,
			&i.MemberUpdatedAt,
			&i.MemberNumberWithoutGaps,
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
SELECT id,
  name,
  score,
  ranking,
  data,
  created_at,
  updated_at,
  member_id,
  user_id,
  member_number,
  member_data,
  joined_at,
  member_updated_at,
  member_number_without_gaps
FROM ranked_team_with_member
WHERE name LIKE CONCAT('%', ?, '%')
  AND member_number_without_gaps <= CAST(? AS UNSIGNED)
  AND member_number_without_gaps > CAST(? AS UNSIGNED)
  AND id IN (
    SELECT id
    FROM (
        SELECT id
        FROM team
        WHERE name LIKE CONCAT('%', ?, '%')
        ORDER BY score DESC,
          id
        LIMIT ? OFFSET ?
      ) temp_team
  )
ORDER BY score DESC,
  id,
  member_number
`

type SearchTeamsParams struct {
	Query                 interface{} `db:"query"`
	MemberLimitPlusOffset int64       `db:"member_limit_plus_offset"`
	MemberOffset          int64       `db:"member_offset"`
	Limit                 int32       `db:"limit"`
	Offset                int32       `db:"offset"`
}

func (q *Queries) SearchTeams(ctx context.Context, arg SearchTeamsParams) ([]RankedTeamWithMember, error) {
	rows, err := q.db.QueryContext(ctx, SearchTeams,
		arg.Query,
		arg.MemberLimitPlusOffset,
		arg.MemberOffset,
		arg.Query,
		arg.Limit,
		arg.Offset,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []RankedTeamWithMember
	for rows.Next() {
		var i RankedTeamWithMember
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.Score,
			&i.Ranking,
			&i.Data,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.MemberID,
			&i.UserID,
			&i.MemberNumber,
			&i.MemberData,
			&i.JoinedAt,
			&i.MemberUpdatedAt,
			&i.MemberNumberWithoutGaps,
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
