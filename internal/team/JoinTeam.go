package team

import (
	"context"
	"database/sql"
	"errors"

	"github.com/MorhafAlshibly/coanda/api"
	"github.com/MorhafAlshibly/coanda/internal/team/model"
	"github.com/MorhafAlshibly/coanda/pkg/conversion"
	errorcode "github.com/MorhafAlshibly/coanda/pkg/errorcode"
	"github.com/go-sql-driver/mysql"
)

type JoinTeamCommand struct {
	service *Service
	In      *api.JoinTeamRequest
	Out     *api.JoinTeamResponse
}

func NewJoinTeamCommand(service *Service, in *api.JoinTeamRequest) *JoinTeamCommand {
	return &JoinTeamCommand{
		service: service,
		In:      in,
	}
}

func (c *JoinTeamCommand) Execute(ctx context.Context) error {
	// Check if user id is provided
	if c.In.UserId == 0 {
		c.Out = &api.JoinTeamResponse{
			Success: false,
			Error:   api.JoinTeamResponse_USER_ID_REQUIRED,
		}
		return nil
	}
	// Check if data is provided
	if c.In.Data == nil {
		c.Out = &api.JoinTeamResponse{
			Success: false,
			Error:   api.JoinTeamResponse_DATA_REQUIRED,
		}
		return nil
	}
	data, err := conversion.ProtobufStructToRawJson(c.In.Data)
	if err != nil {
		return err
	}
	tErr := c.service.checkForTeamRequestError(c.In.Team)
	// Check if field is provided
	if tErr != nil {
		c.Out = &api.JoinTeamResponse{
			Success: false,
			Error:   conversion.Enum(*tErr, api.JoinTeamResponse_Error_value, api.JoinTeamResponse_NO_FIELD_SPECIFIED),
		}
		return nil
	}
	// Check if team member is initialized
	if c.In.Team.Member == nil {
		c.In.Team.Member = &api.TeamMemberRequest{}
	}
	tx, err := c.service.sql.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	qtx := c.service.database.WithTx(tx)
	team, err := qtx.GetTeam(ctx, model.GetTeamParams{
		ID:   conversion.Uint64ToSqlNullInt64(c.In.Team.Id),
		Name: conversion.StringToSqlNullString(c.In.Team.Name),
		Member: model.GetTeamMemberParams{
			ID:     conversion.Uint64ToSqlNullInt64(c.In.Team.Member.Id),
			UserID: conversion.Uint64ToSqlNullInt64(c.In.Team.Member.UserId),
		},
	})
	// If the team is not found, return appropriate error
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			c.Out = &api.JoinTeamResponse{
				Success: false,
				Error:   api.JoinTeamResponse_NOT_FOUND,
			}
			return nil
		}
		return err
	}
	// Get highest member number
	firstOpenMemberNumber, err := qtx.GetFirstOpenMemberNumber(ctx, team.ID)
	if err != nil {
		return err
	}
	if firstOpenMemberNumber > uint32(c.service.maxMembers) {
		c.Out = &api.JoinTeamResponse{
			Success: false,
			Error:   api.JoinTeamResponse_TEAM_FULL,
		}
		return nil
	}
	// Add the member to the team
	result, err := qtx.CreateTeamMember(ctx, model.CreateTeamMemberParams{
		UserID:       c.In.UserId,
		TeamID:       team.ID,
		Data:         data,
		MemberNumber: firstOpenMemberNumber,
	})
	// If we have a duplicate entry, either user is already in a team or the team is full
	if err != nil {
		var mysqlErr *mysql.MySQLError
		if errors.As(err, &mysqlErr) {
			if errorcode.IsDuplicateEntry(mysqlErr, "team_member", "team_member_user_id_idx") {
				c.Out = &api.JoinTeamResponse{
					Success: false,
					Error:   api.JoinTeamResponse_ALREADY_IN_A_TEAM,
				}
				return nil
			}
		}
		return err
	}
	// Get the rows affected
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	// If the team is full, return appropriate error
	if rowsAffected == 0 {
		c.Out = &api.JoinTeamResponse{
			Success: false,
			Error:   api.JoinTeamResponse_TEAM_FULL,
		}
		return nil
	}
	c.Out = &api.JoinTeamResponse{
		Success: true,
		Error:   api.JoinTeamResponse_NONE,
	}
	return nil
}
