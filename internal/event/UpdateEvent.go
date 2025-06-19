package event

import (
	"context"
	"database/sql"

	"github.com/MorhafAlshibly/coanda/api"
	"github.com/MorhafAlshibly/coanda/internal/event/model"
	"github.com/MorhafAlshibly/coanda/pkg/conversion"
)

type UpdateEventCommand struct {
	service *Service
	In      *api.UpdateEventRequest
	Out     *api.UpdateEventResponse
}

func NewUpdateEventCommand(service *Service, in *api.UpdateEventRequest) *UpdateEventCommand {
	return &UpdateEventCommand{
		service: service,
		In:      in,
	}
}

func (c *UpdateEventCommand) Execute(ctx context.Context) error {
	eErr := c.service.checkForEventRequestError(c.In.Event)
	// Check if error is found
	if eErr != nil {
		c.Out = &api.UpdateEventResponse{
			Success: false,
			Error:   conversion.Enum(*eErr, api.UpdateEventResponse_Error_value, api.UpdateEventResponse_ID_OR_NAME_REQUIRED),
		}
		return nil
	}
	// Create empty event request if not provided
	if c.In.Event == nil {
		c.In.Event = &api.EventRequest{}
	}
	// Check if no update is specified
	if c.In.Data == nil {
		c.Out = &api.UpdateEventResponse{
			Success: false,
			Error:   api.UpdateEventResponse_DATA_REQUIRED,
		}
		return nil
	}
	// Prepare data
	data, err := conversion.ProtobufStructToRawJson(c.In.Data)
	if err != nil {
		return err
	}
	tx, err := c.service.sql.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	qtx := c.service.database.WithTx(tx)
	// Update event
	result, err := qtx.UpdateEvent(ctx, model.UpdateEventParams{
		Event: model.GetEventParams{
			ID:   conversion.Uint64ToSqlNullInt64(c.In.Event.Id),
			Name: conversion.StringToSqlNullString(c.In.Event.Name),
		},
		Data: data,
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
		_, err = qtx.GetEvent(ctx, model.GetEventParams{
			ID:   conversion.Uint64ToSqlNullInt64(c.In.Event.Id),
			Name: conversion.StringToSqlNullString(c.In.Event.Name),
		}, nil)
		if err != nil {
			if err == sql.ErrNoRows {
				c.Out = &api.UpdateEventResponse{
					Success: false,
					Error:   api.UpdateEventResponse_NOT_FOUND,
				}
				return nil
			}
			return err
		}
	}
	c.Out = &api.UpdateEventResponse{
		Success: true,
		Error:   api.UpdateEventResponse_NONE,
	}
	return nil
}
