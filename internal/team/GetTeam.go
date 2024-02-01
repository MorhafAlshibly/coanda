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
	In      *api.TeamRequest
	Out     *api.GetTeamResponse
}

func NewGetTeamCommand(service *Service, in *api.TeamRequest) *GetTeamCommand {
	return &GetTeamCommand{
		service: service,
		In:      in,
	}
}

func (c *GetTeamCommand) Execute(ctx context.Context) error {
	field := c.service.GetTeamField(c.In)
	var team model.RankedTeam
	var err error
	// Check if name or owner is provided
	if field == NAME || field == OWNER {
		team, err = c.service.database.GetTeam(ctx, model.GetTeamParams{
			Name: sql.NullString{
				String: conversion.PointerToValue(c.In.Name, ""),
				Valid:  field == NAME,
			},
			Owner: sql.NullInt64{
				Int64: int64(conversion.PointerToValue(c.In.Owner, 0)),
				Valid: field == OWNER,
			}})
		// Check if member is provided
	} else if field == MEMBER {
		team, err = c.service.database.GetTeamByMember(
			ctx,
			*c.In.Member,
		)
	} else {
		c.Out = &api.GetTeamResponse{
			Success: false,
			Error:   conversion.Enum(field, api.GetTeamResponse_Error_value, api.GetTeamResponse_NOT_FOUND),
		}
		return nil
	}
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
	out, err := UnmarshalTeam(team)
	if err != nil {
		return err
	}
	c.Out = &api.GetTeamResponse{
		Success: true,
		Error:   api.GetTeamResponse_NONE,
		Team:    out,
	}
	return nil
}
