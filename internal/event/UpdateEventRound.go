package event

import (
	"context"
	"encoding/json"

	"github.com/MorhafAlshibly/coanda/api"
	"github.com/MorhafAlshibly/coanda/internal/event/model"
	"github.com/MorhafAlshibly/coanda/pkg/conversion"
)

type UpdateEventRoundCommand struct {
	service *Service
	In      *api.UpdateEventRoundRequest
	Out     *api.UpdateEventRoundResponse
}

func NewUpdateEventRoundCommand(service *Service, in *api.UpdateEventRoundRequest) *UpdateEventRoundCommand {
	return &UpdateEventRoundCommand{
		service: service,
		In:      in,
	}
}

func (c *UpdateEventRoundCommand) Execute(ctx context.Context) error {
	erErr := c.service.checkForEventRoundRequestError(c.In.Round)
	if erErr != nil {
		c.Out = &api.UpdateEventRoundResponse{
			Success: false,
			Error:   conversion.Enum(*erErr, api.UpdateEventRoundResponse_Error_value, api.UpdateEventRoundResponse_ID_OR_NAME_REQUIRED),
		}
		return nil
	}
	// Create empty event request if not provided
	if c.In.Round.Event == nil {
		c.In.Round.Event = &api.EventRequest{}
	}
	var data json.RawMessage
	var scoring json.RawMessage
	var err error
	if c.In.Data != nil {
		data, err = conversion.ProtobufStructToRawJson(c.In.Data)
		if err != nil {
			return err
		}
	}
	if len(c.In.Scoring) > 0 {
		scoring, err = conversion.MapToRawJson(map[string]interface{}{"scoring": c.In.Scoring})
		if err != nil {
			return err
		}
	}
	if c.In.Data == nil && len(c.In.Scoring) == 0 {
		c.Out = &api.UpdateEventRoundResponse{
			Success: false,
			Error:   api.UpdateEventRoundResponse_NO_UPDATE_SPECIFIED,
		}
		return nil
	}
	result, err := c.service.database.UpdateEventRound(ctx, model.UpdateEventRoundParams{
		EventRound: model.GetEventRoundParams{
			Event: model.GetEventParams{
				ID:   conversion.Uint64ToSqlNullInt64(c.In.Round.Event.Id),
				Name: conversion.StringToSqlNullString(c.In.Round.Event.Name),
			},
			ID:   conversion.Uint64ToSqlNullInt64(c.In.Round.Id),
			Name: conversion.StringToSqlNullString(c.In.Round.RoundName),
		},
		Data:    data,
		Scoring: scoring,
	})
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		// Check if no rows were affected
		_, err = c.service.database.GetEventRound(ctx, model.GetEventRoundParams{
			Event: model.GetEventParams{
				ID:   conversion.Uint64ToSqlNullInt64(c.In.Round.Event.Id),
				Name: conversion.StringToSqlNullString(c.In.Round.Event.Name),
			},
			ID:   conversion.Uint64ToSqlNullInt64(c.In.Round.Id),
			Name: conversion.StringToSqlNullString(c.In.Round.RoundName),
		})
		if err != nil {
			c.Out = &api.UpdateEventRoundResponse{
				Success: false,
				Error:   api.UpdateEventRoundResponse_NOT_FOUND,
			}
			return nil
		}
	}
	c.Out = &api.UpdateEventRoundResponse{
		Success: true,
		Error:   api.UpdateEventRoundResponse_NONE,
	}
	return nil
}
