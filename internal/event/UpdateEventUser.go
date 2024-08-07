package event

import (
	"context"

	"github.com/MorhafAlshibly/coanda/api"
	"github.com/MorhafAlshibly/coanda/internal/event/model"
	"github.com/MorhafAlshibly/coanda/pkg/conversion"
)

type UpdateEventUserCommand struct {
	service *Service
	In      *api.UpdateEventUserRequest
	Out     *api.UpdateEventUserResponse
}

func NewUpdateEventUserCommand(service *Service, in *api.UpdateEventUserRequest) *UpdateEventUserCommand {
	return &UpdateEventUserCommand{
		service: service,
		In:      in,
	}
}

func (c *UpdateEventUserCommand) Execute(ctx context.Context) error {
	eErr := c.service.checkForEventUserRequestError(c.In.User)
	// Check if error is found
	if eErr != nil {
		c.Out = &api.UpdateEventUserResponse{
			Success: false,
			Error:   conversion.Enum(*eErr, api.UpdateEventUserResponse_Error_value, api.UpdateEventUserResponse_NOT_FOUND),
		}
		return nil
	}
	// Create blank event request if not provided
	if c.In.User.Event == nil {
		c.In.User.Event = &api.EventRequest{}
	}
	// Check if no update is specified
	if c.In.Data == nil {
		c.Out = &api.UpdateEventUserResponse{
			Success: false,
			Error:   api.UpdateEventUserResponse_DATA_REQUIRED,
		}
		return nil
	}
	// Prepare data
	data, err := conversion.ProtobufStructToRawJson(c.In.Data)
	if err != nil {
		return err
	}
	// Update event user
	result, err := c.service.database.UpdateEventUser(ctx, model.UpdateEventUserParams{
		User: model.GetEventUserParams{
			Event: model.GetEventParams{
				ID:   conversion.Uint64ToSqlNullInt64(c.In.User.Event.Id),
				Name: conversion.StringToSqlNullString(c.In.User.Event.Name),
			},
			ID:     conversion.Uint64ToSqlNullInt64(c.In.User.Id),
			UserID: conversion.Uint64ToSqlNullInt64(c.In.User.Id),
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
	// If no rows are affected, the event user is not found
	if rowsAffected == 0 {
		// Check if no rows were affected
		_, err = c.service.database.GetEventUser(ctx, model.GetEventUserParams{
			Event: model.GetEventParams{
				ID:   conversion.Uint64ToSqlNullInt64(c.In.User.Event.Id),
				Name: conversion.StringToSqlNullString(c.In.User.Event.Name),
			},
			ID:     conversion.Uint64ToSqlNullInt64(c.In.User.Id),
			UserID: conversion.Uint64ToSqlNullInt64(c.In.User.Id),
		})
		if err != nil {
			c.Out = &api.UpdateEventUserResponse{
				Success: false,
				Error:   api.UpdateEventUserResponse_NOT_FOUND,
			}
			return nil
		}
	}
	c.Out = &api.UpdateEventUserResponse{
		Success: true,
		Error:   api.UpdateEventUserResponse_NONE,
	}
	return nil
}
