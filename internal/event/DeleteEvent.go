package event

import (
	"context"

	"github.com/MorhafAlshibly/coanda/api"
	"github.com/MorhafAlshibly/coanda/internal/event/model"
	"github.com/MorhafAlshibly/coanda/pkg/conversion"
)

type DeleteEventCommand struct {
	service *Service
	In      *api.EventRequest
	Out     *api.EventResponse
}

func NewDeleteEventCommand(service *Service, in *api.EventRequest) *DeleteEventCommand {
	return &DeleteEventCommand{
		service: service,
		In:      in,
	}
}

func (c *DeleteEventCommand) Execute(ctx context.Context) error {
	eErr := c.service.checkForEventRequestError(c.In)
	// Check if error is found
	if eErr != nil {
		c.Out = &api.EventResponse{
			Success: false,
			Error:   conversion.Enum(*eErr, api.EventResponse_Error_value, api.EventResponse_ID_OR_NAME_REQUIRED),
		}
		return nil
	}
	result, err := c.service.database.DeleteEvent(ctx, model.GetEventParams{
		ID:   conversion.Uint64ToSqlNullInt64(c.In.Id),
		Name: conversion.StringToSqlNullString(c.In.Name),
	})
	// If an error occurs, it is an internal server error
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	// If no rows are affected, the event is not found
	if rowsAffected == 0 {
		c.Out = &api.EventResponse{
			Success: false,
			Error:   api.EventResponse_NOT_FOUND,
		}
		return nil
	}
	c.Out = &api.EventResponse{
		Success: true,
		Error:   api.EventResponse_NONE,
	}
	return nil
}
