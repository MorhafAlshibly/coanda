package team

import (
	"context"
	"database/sql"
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
	var data json.RawMessage
	dataExists := int64(0)
	scoreOffset := int64(0)
	incrementScore := false
	var err error
	if c.In.Data != nil {
		data, err = conversion.ProtobufStructToRawJson(c.In.Data)
		if err != nil {
			return err
		}
		dataExists = 1
	}
	if c.In.Score != nil {
		if c.In.IncrementScore == nil {
			c.Out = &api.UpdateTeamResponse{
				Success: false,
				Error:   api.UpdateTeamResponse_INCREMENT_SCORE_NOT_SPECIFIED,
			}
			return nil
		}
		incrementScore = *c.In.IncrementScore
		scoreOffset = int64(*c.In.Score)
	}
	if c.In.Score == nil && c.In.Data == nil {
		c.Out = &api.UpdateTeamResponse{
			Success: false,
			Error:   api.UpdateTeamResponse_NO_UPDATE_SPECIFIED,
		}
		return nil
	}
	field := c.service.GetTeamField(c.In.Team)
	var result sql.Result
	// Check if name or owner is provided
	if field == NAME || field == OWNER {
		result, err = c.service.database.UpdateTeam(ctx, model.UpdateTeamParams{
			Name: sql.NullString{
				String: conversion.PointerToValue(c.In.Team.Name, ""),
				Valid:  field == NAME,
			},
			Owner: sql.NullInt64{
				Int64: int64(conversion.PointerToValue(c.In.Team.Owner, 0)),
				Valid: field == OWNER,
			},
			DataExists: dataExists,
			Data:       data,
			ScoreOffset: sql.NullInt64{
				Int64: scoreOffset,
				Valid: c.In.Score != nil,
			},
			IncrementScore: incrementScore,
		})
		// Check if member is provided
	} else if field == MEMBER {
		result, err = c.service.database.UpdateTeamByMember(
			ctx,
			model.UpdateTeamByMemberParams{
				UserID:     *c.In.Team.Member,
				DataExists: dataExists,
				Data:       data,
				ScoreOffset: sql.NullInt64{
					Int64: scoreOffset,
					Valid: c.In.Score != nil,
				},
				IncrementScore: incrementScore,
			})
	} else {
		c.Out = &api.UpdateTeamResponse{
			Success: false,
			Error:   conversion.Enum(field, api.UpdateTeamResponse_Error_value, api.UpdateTeamResponse_NOT_FOUND),
		}
		return nil
	}
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
