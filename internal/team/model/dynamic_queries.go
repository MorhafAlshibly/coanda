package model

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/doug-martin/goqu/v9"
	_ "github.com/doug-martin/goqu/v9/dialect/mysql"
	"github.com/doug-martin/goqu/v9/exp"
)

var gq = goqu.Dialect("mysql")

type GetTeamParams struct {
	Name   sql.NullString `db:"name"`
	Owner  sql.NullInt64  `db:"owner"`
	Member sql.NullInt64  `db:"member"`
}

// filterGetTeamParams filters the GetTeamParams to exp.Expression
func filterGetTeamParams(arg GetTeamParams) exp.Expression {
	expressions := goqu.Ex{}
	if arg.Name.Valid {
		expressions["name"] = arg.Name
	}
	if arg.Owner.Valid {
		expressions["owner"] = arg.Owner
	}
	if arg.Member.Valid {
		expressions["name"] = gq.From("team_member").Select("team").Where(goqu.Ex{"user_id": arg.Member}).Limit(1)
	}
	return expressions
}

func (q *Queries) GetTeam(ctx context.Context, arg GetTeamParams) (RankedTeam, error) {
	team := gq.From("ranked_team").Prepared(true)
	query, args, err := team.Where(filterGetTeamParams(arg)).Limit(1).ToSQL()
	if err != nil {
		return RankedTeam{}, err
	}
	var i RankedTeam
	err = q.db.QueryRowContext(ctx, query, args...).Scan(
		&i.Name,
		&i.Owner,
		&i.Score,
		&i.Ranking,
		&i.Data,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

type GetTeamMembersParams struct {
	Team   GetTeamParams
	Limit  uint64
	Offset uint64
}

// filterGetTeamMembersParams filters the GetTeamMembersParams to exp.Expression
func filterGetTeamMembersParams(arg GetTeamMembersParams) exp.Expression {
	expressions := goqu.Ex{}
	if arg.Team.Name.Valid {
		expressions["team"] = arg.Team.Name
	}
	if arg.Team.Owner.Valid {
		expressions["team"] = gq.From("team_owner").Select("team").Where(goqu.Ex{"user_id": arg.Team.Owner}).Limit(1)
	}
	if arg.Team.Member.Valid {
		expressions["team"] = gq.From("team_member").Select("team").Where(goqu.Ex{"user_id": arg.Team.Member}).Limit(1)
	}
	return expressions
}

func (q *Queries) GetTeamMembers(ctx context.Context, arg GetTeamMembersParams) ([]TeamMember, error) {
	teamMember := gq.From("team_member").Prepared(true).Select("team", "user_id", "data", "joined_at", "updated_at")
	query, args, err := teamMember.Where(filterGetTeamMembersParams(arg)).Limit(uint(arg.Limit)).Offset(uint(arg.Offset)).ToSQL()
	if err != nil {
		return nil, err
	}
	fmt.Println(query, args)
	rows, err := q.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []TeamMember
	for rows.Next() {
		var i TeamMember
		if err := rows.Scan(
			&i.Team,
			&i.UserID,
			&i.Data,
			&i.JoinedAt,
			&i.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	return items, nil
}

func (q *Queries) DeleteTeam(ctx context.Context, arg GetTeamParams) (sql.Result, error) {
	team := gq.Delete("team").Prepared(true)
	query, args, err := team.Where(filterGetTeamParams(arg)).Limit(1).ToSQL()
	if err != nil {
		return nil, err
	}
	return q.db.ExecContext(ctx, query, args...)
}

func (q *Queries) DeleteTeamOwner(ctx context.Context, arg GetTeamParams) (sql.Result, error) {
	team := gq.Delete("team_owner").Prepared(true)
	query, args, err := team.Where(filterGetTeamParams(arg)).Limit(1).ToSQL()
	if err != nil {
		return nil, err
	}
	return q.db.ExecContext(ctx, query, args...)
}

type UpdateTeamParams struct {
	Team           GetTeamParams
	Data           json.RawMessage `db:"data"`
	Score          sql.NullInt64   `db:"score"`
	IncrementScore bool
}

func (q *Queries) UpdateTeam(ctx context.Context, arg UpdateTeamParams) (sql.Result, error) {
	team := gq.Update("team").Prepared(true)
	updates := goqu.Record{}
	if arg.Data != nil {
		updates["data"] = []byte(arg.Data)
	}
	if arg.Score.Valid {
		if arg.IncrementScore {
			updates["score"] = goqu.L("score + ?", arg.Score)
		} else {
			updates["score"] = arg.Score
		}
	}
	team = team.Set(updates)
	query, args, err := team.Where(filterGetTeamParams(arg.Team)).Limit(1).ToSQL()
	if err != nil {
		return nil, err
	}
	fmt.Println(query, args)
	return q.db.ExecContext(ctx, query, args...)
}

type UpdateTeamMembersParams struct {
	Team GetTeamParams
	Data json.RawMessage `db:"data"`
}

func (q *Queries) UpdateTeamMembers(ctx context.Context, arg UpdateTeamMembersParams) (sql.Result, error) {
	teamMember := gq.Update("team_member").Prepared(true)
	if arg.Data != nil {
		teamMember = teamMember.Set(goqu.Record{"data": arg.Data})
	}
	query, args, err := teamMember.Where(filterGetTeamParams(arg.Team)).ToSQL()
	if err != nil {
		return nil, err
	}
	return q.db.ExecContext(ctx, query, args...)
}
