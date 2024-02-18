package team

import (
	"context"
	"encoding/json"

	"github.com/MorhafAlshibly/coanda/api"
	"github.com/MorhafAlshibly/coanda/internal/team/model"
	"github.com/MorhafAlshibly/coanda/pkg/conversion"
)

type UpdateTeamCommand struct {
	service *Service
	In      *api.UpdateTeamRequest
	Out     *api.UpdateTeamResponse
}

func NewUpdateTeamCommand(service *Service, in *api.UpdateTeamRequest) *UpdateTeamCommand {
	return &UpdateTeamCommand{
		service: service,
		In:      in,
	}
}

func (c *UpdateTeamCommand) Execute(ctx context.Context) error {
	tErr := c.service.checkForTeamRequestError(c.In.Team)
	// Check if error is found
	if tErr != nil {
		c.Out = &api.UpdateTeamResponse{
			Success: false,
			Error:   conversion.Enum(*tErr, api.UpdateTeamResponse_Error_value, api.UpdateTeamResponse_NOT_FOUND),
		}
		return nil
	}
	// Check if no update is specified
	if c.In.Score == nil && c.In.Data == nil {
		c.Out = &api.UpdateTeamResponse{
			Success: false,
			Error:   api.UpdateTeamResponse_NO_UPDATE_SPECIFIED,
		}
		return nil
	}
	// Check if score is specified without whether to increment it
	if c.In.Score != nil && c.In.IncrementScore == nil {
		c.Out = &api.UpdateTeamResponse{
			Success: false,
			Error:   api.UpdateTeamResponse_INCREMENT_SCORE_NOT_SPECIFIED,
		}
		return nil
	}
	// Prepare data
	data := json.RawMessage(nil)
	if c.In.Data != nil {
		var err error
		data, err = conversion.ProtobufStructToRawJson(c.In.Data)
		if err != nil {
			return err
		}
	}
	result, err := c.service.database.UpdateTeam(ctx, model.UpdateTeamParams{
		Team: model.GetTeamParams{
			Name:   conversion.StringToSqlNullString(c.In.Team.Name),
			Owner:  conversion.Uint64ToSqlNullInt64(c.In.Team.Owner),
			Member: conversion.Uint64ToSqlNullInt64(c.In.Team.Member),
		},
		Data:           data,
		Score:          conversion.Int64ToSqlNullInt64(c.In.Score),
		IncrementScore: conversion.PointerBoolToValue(c.In.IncrementScore),
	})
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		c.Out = &api.UpdateTeamResponse{
			Success: false,
			Error:   api.UpdateTeamResponse_NOT_FOUND,
		}
		return nil
	}
	c.Out = &api.UpdateTeamResponse{
		Success: true,
		Error:   api.UpdateTeamResponse_NONE,
	}
	return nil

}
