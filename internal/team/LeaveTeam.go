package team

import (
	"context"

	"github.com/MorhafAlshibly/coanda/api"
	"github.com/MorhafAlshibly/coanda/internal/team/model"
	"github.com/MorhafAlshibly/coanda/pkg/conversion"
)

type LeaveTeamCommand struct {
	service *Service
	In      *api.TeamMemberRequest
	Out     *api.LeaveTeamResponse
}

func NewLeaveTeamCommand(service *Service, in *api.TeamMemberRequest) *LeaveTeamCommand {
	return &LeaveTeamCommand{
		service: service,
		In:      in,
	}
}

func (c *LeaveTeamCommand) Execute(ctx context.Context) error {
	tErr := c.service.checkForTeamMemberRequestError(c.In)
	// Check if error is found
	if tErr != nil {
		c.Out = &api.LeaveTeamResponse{
			Success: false,
			Error:   conversion.Enum(*tErr, api.LeaveTeamResponse_Error_value, api.LeaveTeamResponse_NO_FIELD_SPECIFIED),
		}
		return nil
	}
	tx, err := c.service.sql.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	qtx := c.service.database.WithTx(tx)
	result, err := qtx.DeleteTeamMember(ctx, model.GetTeamMemberParams{
		ID:     conversion.Uint64ToSqlNullInt64(c.In.Id),
		UserID: conversion.Uint64ToSqlNullInt64(c.In.UserId),
	})
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		c.Out = &api.LeaveTeamResponse{
			Success: false,
			Error:   api.LeaveTeamResponse_NOT_FOUND,
		}
		return nil
	}
	// Check if that was the last member of the team
	teamMembers, err := qtx.GetTeamMembers(ctx, model.GetTeamMembersParams{
		Team: model.GetTeamParams{
			Member: model.GetTeamMemberParams{
				ID:     conversion.Uint64ToSqlNullInt64(c.In.Id),
				UserID: conversion.Uint64ToSqlNullInt64(c.In.UserId),
			},
		},
		Limit:  1,
		Offset: 0,
	})
	if err != nil {
		return err
	}
	if len(teamMembers) == 0 {
		_, err = qtx.DeleteTeam(ctx, model.GetTeamParams{
			ID: conversion.Uint64ToSqlNullInt64(c.In.Id),
		})
		if err != nil {
			return err
		}
	}
	err = tx.Commit()
	if err != nil {
		return err
	}
	c.Out = &api.LeaveTeamResponse{
		Success: true,
		Error:   api.LeaveTeamResponse_NONE,
	}
	return nil
}
