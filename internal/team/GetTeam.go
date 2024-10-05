package team

import (
	"context"
	"database/sql"
	"errors"

	"github.com/MorhafAlshibly/coanda/api"
	"github.com/MorhafAlshibly/coanda/internal/team/model"
	"github.com/MorhafAlshibly/coanda/pkg/conversion"
)

type GetTeamCommand struct {
	service *Service
	In      *api.GetTeamRequest
	Out     *api.GetTeamResponse
}

func NewGetTeamCommand(service *Service, in *api.GetTeamRequest) *GetTeamCommand {
	return &GetTeamCommand{
		service: service,
		In:      in,
	}
}

func (c *GetTeamCommand) Execute(ctx context.Context) error {
	limit, offset := conversion.PaginationToLimitOffset(c.In.Pagination, c.service.defaultMaxPageLength, c.service.maxMaxPageLength)
	tErr := c.service.checkForTeamRequestError(c.In.Team)
	// Check if error is found
	if tErr != nil {
		c.Out = &api.GetTeamResponse{
			Success: false,
			Error:   conversion.Enum(*tErr, api.GetTeamResponse_Error_value, api.GetTeamResponse_NO_FIELD_SPECIFIED),
		}
		return nil
	}
	// Check if team member is initialised
	if c.In.Team.Member == nil {
		c.In.Team.Member = &api.TeamMemberRequest{}
	}
	team, err := c.service.database.GetTeam(ctx, model.GetTeamParams{
		Team: model.TeamParams{
			ID:   conversion.Uint64ToSqlNullInt64(c.In.Team.Id),
			Name: conversion.StringToSqlNullString(c.In.Team.Name),
			Member: model.GetTeamMemberParams{
				ID:     conversion.Uint64ToSqlNullInt64(c.In.Team.Member.Id),
				UserID: conversion.Uint64ToSqlNullInt64(c.In.Team.Member.UserId),
			},
		},
		Limit:  limit,
		Offset: offset,
	})
	// Check if team is found
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			c.Out = &api.GetTeamResponse{
				Success: false,
				Error:   api.GetTeamResponse_NOT_FOUND,
			}
			return nil
		}
		return err
	}
	out, err := unmarshalTeamWithMembers(team)
	if err != nil {
		return err
	}
	c.Out = &api.GetTeamResponse{
		Success: true,
		Team:    out,
		Error:   api.GetTeamResponse_NONE,
	}
	return nil
}
