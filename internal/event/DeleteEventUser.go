package event

import (
	"context"

	"github.com/MorhafAlshibly/coanda/api"
	"github.com/MorhafAlshibly/coanda/internal/event/model"
	"github.com/MorhafAlshibly/coanda/pkg/conversion"
)

type DeleteEventUserCommand struct {
	service *Service
	In      *api.EventUserRequest
	Out     *api.EventUserResponse
}

func NewDeleteEventUserCommand(service *Service, in *api.EventUserRequest) *DeleteEventUserCommand {
	return &DeleteEventUserCommand{
		service: service,
		In:      in,
	}
}

func (c *DeleteEventUserCommand) Execute(ctx context.Context) error {
	eErr := c.service.checkForEventUserRequestError(c.In)
	// Check if error is found
	if eErr != nil {
		c.Out = &api.EventUserResponse{
			Success: false,
			Error:   conversion.Enum(*eErr, api.EventUserResponse_Error_value, api.EventUserResponse_NOT_FOUND),
		}
		return nil
	}
	// Create empty event request if not provided
	if c.In.Event == nil {
		c.In.Event = &api.EventRequest{}
	}
	result, err := c.service.database.DeleteEventUser(ctx, model.GetEventUserParams{
		Event: model.GetEventParams{
			ID:   conversion.Uint64ToSqlNullInt64(c.In.Event.Id),
			Name: conversion.StringToSqlNullString(c.In.Event.Name),
		},
		ID:     conversion.Uint64ToSqlNullInt64(c.In.Id),
		UserID: conversion.Uint64ToSqlNullInt64(c.In.UserId),
	})
	// If an error occurs, it is an internal server error
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	// If no rows are affected, the event user is not found
	if rowsAffected == 0 {
		c.Out = &api.EventUserResponse{
			Success: false,
			Error:   api.EventUserResponse_NOT_FOUND,
		}
		return nil
	}
	c.Out = &api.EventUserResponse{
		Success: true,
		Error:   api.EventUserResponse_NONE,
	}
	return nil
}
