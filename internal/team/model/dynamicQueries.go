package model

import (
	"context"
	"database/sql"
	"encoding/json"

	"github.com/doug-martin/goqu/v9"
	_ "github.com/doug-martin/goqu/v9/dialect/mysql"
	"github.com/doug-martin/goqu/v9/exp"
)

var gq = goqu.Dialect("mysql")

type GetTeamMemberParams struct {
	ID     sql.NullInt64 `db:"id"`
	UserID sql.NullInt64 `db:"user_id"`
}

// filterGetTeamMemberParams filters the GetTeamMemberParams to exp.Expression
func filterGetTeamMemberParams(arg GetTeamMemberParams) exp.Expression {
	expressions := goqu.Ex{}
	if arg.ID.Valid {
		expressions["id"] = arg.ID
	}
	if arg.UserID.Valid {
		expressions["user_id"] = arg.UserID
	}
	return expressions
}

func (q *Queries) GetTeamMember(ctx context.Context, arg GetTeamMemberParams) (TeamMember, error) {
	teamMember := gq.From("team_member").Prepared(true).Select("id", "user_id", "team_id", "member_number", "data", "joined_at", "updated_at")
	query, args, err := teamMember.Where(filterGetTeamMemberParams(arg)).Limit(1).ToSQL()
	if err != nil {
		return TeamMember{}, err
	}
	var i TeamMember
	err = q.db.QueryRowContext(ctx, query, args...).Scan(
		&i.ID,
		&i.UserID,
		&i.TeamID,
		&i.MemberNumber,
		&i.Data,
		&i.JoinedAt,
		&i.UpdatedAt,
	)
	return i, err
}

func (q *Queries) DeleteTeamMember(ctx context.Context, arg GetTeamMemberParams) (sql.Result, error) {
	teamMember := gq.Delete("team_member").Prepared(true)
	query, args, err := teamMember.Where(filterGetTeamMemberParams(arg)).Limit(1).ToSQL()
	if err != nil {
		return nil, err
	}
	return q.db.ExecContext(ctx, query, args...)
}

type UpdateTeamMemberParams struct {
	TeamMember GetTeamMemberParams
	Data       json.RawMessage `db:"data"`
}

func (q *Queries) UpdateTeamMember(ctx context.Context, arg UpdateTeamMemberParams) (sql.Result, error) {
	teamMember := gq.Update("team_member").Prepared(true)
	updates := goqu.Record{}
	if arg.Data != nil {
		updates["data"] = []byte(arg.Data)
	}
	teamMember = teamMember.Set(updates)
	query, args, err := teamMember.Where(filterGetTeamMemberParams(arg.TeamMember)).Limit(1).ToSQL()
	if err != nil {
		return nil, err
	}
	return q.db.ExecContext(ctx, query, args...)
}

type TeamParams struct {
	ID     sql.NullInt64  `db:"id"`
	Name   sql.NullString `db:"name"`
	Member GetTeamMemberParams
}

type GetTeamParams struct {
	Team   TeamParams
	Limit  uint64
	Offset uint64
}

// filterGetTeamParams filters the GetTeamParams to exp.Expression and exp.LiteralExpression
func filterGetTeamParams(arg GetTeamParams) (exp.Expression, exp.LiteralExpression) {
	expressions := goqu.Ex{}
	if arg.Team.ID.Valid {
		expressions["id"] = arg.Team.ID
	}
	if arg.Team.Name.Valid {
		expressions["name"] = arg.Team.Name
	}
	if arg.Team.Member.ID.Valid {
		expressions["id"] = gq.From(gq.From("team_member").Select("team_id").Where(goqu.Ex{"id": arg.Team.Member.ID}).Limit(1))
	}
	if arg.Team.Member.UserID.Valid {
		expressions["id"] = gq.From(gq.From("team_member").Select("team_id").Where(goqu.Ex{"user_id": arg.Team.Member.UserID}).Limit(1))
	}
	limitMembers := goqu.L("member_number_without_gaps < ? AND member_number_without_gaps >= ?", arg.Limit+arg.Offset, arg.Offset)
	return expressions, limitMembers
}

func (q *Queries) GetTeam(ctx context.Context, arg GetTeamParams) ([]RankedTeamWithMember, error) {
	team := gq.From("ranked_team_with_member").Prepared(true).Select("id", "name", "score", "ranking", "data", "created_at", "updated_at", "member_id", "member_user_id", "member_number", "member_data", "joined_at", "member_updated_at")
	query, args, err := team.Where(filterGetTeamParams(arg)).ToSQL()
	if err != nil {
		return nil, err
	}
	rows, err := q.db.QueryContext(ctx, query, args...)
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
	return items, nil
}

type GetTeamMembersParams struct {
	Team   TeamParams
	Limit  uint64
	Offset uint64
}

// filterGetTeamMembersParams filters the GetTeamMembersParams to exp.Expression
func filterGetTeamMembersParams(arg GetTeamMembersParams) exp.Expression {
	expressions := goqu.Ex{}
	if arg.Team.ID.Valid {
		expressions["team_id"] = arg.Team.ID
	}
	if arg.Team.Name.Valid {
		expressions["team_id"] = gq.From(gq.From("team").Select("id").Where(goqu.Ex{"name": arg.Team.Name}).Limit(1))
	}
	if arg.Team.Member.ID.Valid {
		expressions["team_id"] = gq.From(gq.From("team_member").Select("team_id").Where(goqu.Ex{"id": arg.Team.Member.ID}).Limit(1))
	}
	if arg.Team.Member.UserID.Valid {
		expressions["team_id"] = gq.From(gq.From("team_member").Select("team_id").Where(goqu.Ex{"user_id": arg.Team.Member.UserID}).Limit(1))
	}
	return expressions
}

func (q *Queries) GetTeamMembers(ctx context.Context, arg GetTeamMembersParams) ([]TeamMember, error) {
	teamMember := gq.From("team_member").Prepared(true).Select("id", "user_id", "team_id", "member_number", "data", "joined_at", "updated_at")
	query, args, err := teamMember.Where(filterGetTeamMembersParams(arg)).Order(goqu.C("member_number").Asc()).Limit(uint(arg.Limit)).Offset(uint(arg.Offset)).ToSQL()
	if err != nil {
		return nil, err
	}
	rows, err := q.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []TeamMember
	for rows.Next() {
		var i TeamMember
		if err := rows.Scan(
			&i.ID,
			&i.UserID,
			&i.TeamID,
			&i.MemberNumber,
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

func filterTeamParams(arg TeamParams) exp.Expression {
	expressions := goqu.Ex{}
	if arg.ID.Valid {
		expressions["id"] = arg.ID
	}
	if arg.Name.Valid {
		expressions["name"] = arg.Name
	}
	if arg.Member.ID.Valid {
		expressions["id"] = gq.From(gq.From("team_member").Select("team_id").Where(goqu.Ex{"id": arg.Member.ID}).Limit(1))
	}
	if arg.Member.UserID.Valid {
		expressions["id"] = gq.From(gq.From("team_member").Select("team_id").Where(goqu.Ex{"user_id": arg.Member.UserID}).Limit(1))
	}
	return expressions
}

func (q *Queries) DeleteTeam(ctx context.Context, arg TeamParams) (sql.Result, error) {
	team := gq.Delete("team").Prepared(true)
	query, args, err := team.Where(filterTeamParams(arg)).Limit(1).ToSQL()
	if err != nil {
		return nil, err
	}
	return q.db.ExecContext(ctx, query, args...)
}

type UpdateTeamParams struct {
	Team           TeamParams
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
	query, args, err := team.Where(filterTeamParams(arg.Team)).Limit(1).ToSQL()
	if err != nil {
		return nil, err
	}
	return q.db.ExecContext(ctx, query, args...)
}
