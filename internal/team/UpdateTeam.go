package team

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"

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
			Error:   conversion.Enum(*tErr, api.UpdateTeamResponse_Error_value, api.UpdateTeamResponse_NO_FIELD_SPECIFIED),
		}
		return nil
	}
	// Check if member is initialized
	if c.In.Team.Member == nil {
		c.In.Team.Member = &api.TeamMemberRequest{}
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
		Team: model.TeamParams{
			ID:   conversion.Uint64ToSqlNullInt64(c.In.Team.Id),
			Name: conversion.StringToSqlNullString(c.In.Team.Name),
			Member: model.GetTeamMemberParams{
				ID:     conversion.Uint64ToSqlNullInt64(c.In.Team.Member.Id),
				UserID: conversion.Uint64ToSqlNullInt64(c.In.Team.Member.UserId),
			},
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
		// Check if we didn't find a row
		_, err = c.service.database.GetTeam(ctx, model.GetTeamParams{
			Team: model.TeamParams{
				ID:   conversion.Uint64ToSqlNullInt64(c.In.Team.Id),
				Name: conversion.StringToSqlNullString(c.In.Team.Name),
				Member: model.GetTeamMemberParams{
					ID:     conversion.Uint64ToSqlNullInt64(c.In.Team.Member.Id),
					UserID: conversion.Uint64ToSqlNullInt64(c.In.Team.Member.UserId),
				},
			},
			Limit:  1,
			Offset: 0,
		})
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				// If we didn't find a row
				c.Out = &api.UpdateTeamResponse{
					Success: false,
					Error:   api.UpdateTeamResponse_NOT_FOUND,
				}
				return nil
			}
			return err
		}
	}
	c.Out = &api.UpdateTeamResponse{
		Success: true,
		Error:   api.UpdateTeamResponse_NONE,
	}
	return nil
}
