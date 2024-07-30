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
	tErr := c.service.checkForTeamRequestError(c.In)
	// Check if error is found
	if tErr != nil {
		c.Out = &api.GetTeamResponse{
			Success: false,
			Error:   conversion.Enum(*tErr, api.GetTeamResponse_Error_value, api.GetTeamResponse_NOT_FOUND),
		}
		return nil
	}
	team, err := c.service.database.GetTeam(ctx, model.GetTeamParams{
		Name:   conversion.StringToSqlNullString(c.In.Name),
		Owner:  conversion.Uint64ToSqlNullInt64(c.In.Owner),
		Member: conversion.Uint64ToSqlNullInt64(c.In.Member),
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
	out, err := unmarshalTeam(team)
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
