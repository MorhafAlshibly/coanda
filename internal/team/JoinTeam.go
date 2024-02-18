package team

import (
	"context"
	"database/sql"
	"errors"

	"github.com/MorhafAlshibly/coanda/api"
	"github.com/MorhafAlshibly/coanda/internal/team/model"
	"github.com/MorhafAlshibly/coanda/pkg/conversion"
	errorcodes "github.com/MorhafAlshibly/coanda/pkg/errorcodes"
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
			Error:   conversion.Enum(*tErr, api.JoinTeamResponse_Error_value, api.JoinTeamResponse_NOT_FOUND),
		}
		return nil
	}
	team, err := c.service.database.GetTeam(ctx, model.GetTeamParams{
		Name:   conversion.StringToSqlNullString(c.In.Team.Name),
		Owner:  conversion.Uint64ToSqlNullInt64(c.In.Team.Owner),
		Member: conversion.Uint64ToSqlNullInt64(c.In.Team.Member),
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
	// Add the member to the team
	result, err := c.service.database.CreateTeamMember(ctx, model.CreateTeamMemberParams{
		Team:       team.Name,
		UserID:     c.In.UserId,
		Data:       data,
		MaxMembers: int64(c.service.maxMembers),
	})
	// If the user is already in the team, return appropriate error
	if err != nil {
		var mysqlErr *mysql.MySQLError
		if errors.As(err, &mysqlErr) && mysqlErr.Number == errorcodes.MySQLErrorCodeDuplicateEntry {
			c.Out = &api.JoinTeamResponse{
				Success: false,
				Error:   api.JoinTeamResponse_ALREADY_IN_A_TEAM,
			}
			return nil
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
