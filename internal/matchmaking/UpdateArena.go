package matchmaking

import (
	"context"
	"database/sql"
	"encoding/json"

	"github.com/MorhafAlshibly/coanda/api"
	"github.com/MorhafAlshibly/coanda/internal/matchmaking/model"
	"github.com/MorhafAlshibly/coanda/pkg/conversion"
)

type UpdateArenaCommand struct {
	service *Service
	In      *api.UpdateArenaRequest
	Out     *api.UpdateArenaResponse
}

func NewUpdateArenaCommand(service *Service, in *api.UpdateArenaRequest) *UpdateArenaCommand {
	return &UpdateArenaCommand{
		service: service,
		In:      in,
	}
}

func (c *UpdateArenaCommand) Execute(ctx context.Context) error {
	aErr := c.service.checkForArenaRequestError(c.In.Arena)
	// Check if error is found
	if aErr != nil {
		c.Out = &api.UpdateArenaResponse{
			Success: false,
			Error:   conversion.Enum(*aErr, api.UpdateArenaResponse_Error_value, api.UpdateArenaResponse_NOT_FOUND),
		}
		return nil
	}
	// Check if no update is specified
	if c.In.MinPlayers == nil && c.In.MaxPlayers == nil && c.In.Data == nil {
		c.Out = &api.UpdateArenaResponse{
			Success: false,
			Error:   api.UpdateArenaResponse_NO_UPDATE_SPECIFIED,
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
	if c.In.MinPlayers != nil || c.In.MaxPlayers != nil {
		if c.In.MinPlayers == nil || c.In.MaxPlayers == nil {
			c.Out = &api.UpdateArenaResponse{
				Success: false,
				Error:   api.UpdateArenaResponse_BOTH_MIN_AND_MAX_PLAYERS_REQUIRED,
			}
			return nil
		}
		if *c.In.MinPlayers > *c.In.MaxPlayers {
			c.Out = &api.UpdateArenaResponse{
				Success: false,
				Error:   api.UpdateArenaResponse_MAX_PLAYERS_MUST_BE_GREATER_THAN_MIN_PLAYERS,
			}
			return nil
		}
	}
	result, err := c.service.database.UpdateArena(ctx, model.UpdateArenaParams{
		Arena: model.GetArenaParams{
			ID:   conversion.Uint64ToSqlNullInt64(c.In.Arena.Id),
			Name: conversion.StringToSqlNullString(c.In.Arena.Name),
		},
		Data:       data,
		MinPlayers: conversion.Uint32ToSqlNullInt32(c.In.MinPlayers),
		MaxPlayers: conversion.Uint32ToSqlNullInt32(c.In.MaxPlayers),
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
		_, err = c.service.database.GetArena(ctx, model.GetArenaParams{
			ID:   conversion.Uint64ToSqlNullInt64(c.In.Arena.Id),
			Name: conversion.StringToSqlNullString(c.In.Arena.Name),
		})
		if err != nil {
			if err == sql.ErrNoRows {
				// If we didn't find a row
				c.Out = &api.UpdateArenaResponse{
					Success: false,
					Error:   api.UpdateArenaResponse_NOT_FOUND,
				}
				return nil
			}
			return err
		}
	}
	c.Out = &api.UpdateArenaResponse{
		Success: true,
		Error:   api.UpdateArenaResponse_NONE,
	}
	return nil
}
