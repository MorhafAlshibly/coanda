package team

import (
	"context"

	"github.com/MorhafAlshibly/coanda/api"
	"github.com/MorhafAlshibly/coanda/internal/team/model"
	"github.com/MorhafAlshibly/coanda/pkg/conversion"
)

type DeleteTeamCommand struct {
	service *Service
	In      *api.TeamRequest
	Out     *api.TeamResponse
}

func NewDeleteTeamCommand(service *Service, in *api.TeamRequest) *DeleteTeamCommand {
	return &DeleteTeamCommand{
		service: service,
		In:      in,
	}
}

func (c *DeleteTeamCommand) Execute(ctx context.Context) error {
	tErr := c.service.checkForTeamRequestError(c.In)
	// Check if error is found
	if tErr != nil {
		c.Out = &api.TeamResponse{
			Success: false,
			Error:   conversion.Enum(*tErr, api.TeamResponse_Error_value, api.TeamResponse_NOT_FOUND),
		}
		return nil
	}
	result, err := c.service.database.DeleteTeam(ctx, model.GetTeamParams{
		Name:   conversion.StringToSqlNullString(c.In.Name),
		Owner:  conversion.Uint64ToSqlNullInt64(c.In.Owner),
		Member: conversion.Uint64ToSqlNullInt64(c.In.Member),
	})
	// If an error occurs, it is an internal server error
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	// If no rows are affected, the team is not found
	if rowsAffected == 0 {
		c.Out = &api.TeamResponse{
			Success: false,
			Error:   api.TeamResponse_NOT_FOUND,
		}
		return nil
	}
	c.Out = &api.TeamResponse{
		Success: true,
		Error:   api.TeamResponse_NONE,
	}
	return nil
}
