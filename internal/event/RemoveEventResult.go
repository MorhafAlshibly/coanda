package event

import (
	"context"

	"github.com/MorhafAlshibly/coanda/api"
	"github.com/MorhafAlshibly/coanda/internal/event/model"
	"github.com/MorhafAlshibly/coanda/pkg/conversion"
)

type RemoveEventResultCommand struct {
	service *Service
	In      *api.EventRoundUserRequest
	Out     *api.RemoveEventResultResponse
}

func NewRemoveEventResultCommand(service *Service, in *api.EventRoundUserRequest) *RemoveEventResultCommand {
	return &RemoveEventResultCommand{
		service: service,
		In:      in,
	}
}

func (c *RemoveEventResultCommand) Execute(ctx context.Context) error {
	eruErr := c.service.checkForEventRoundUserRequestError(c.In)
	// Check if error is found
	if eruErr != nil {
		c.Out = &api.RemoveEventResultResponse{
			Success: false,
			Error:   conversion.Enum(*eruErr, api.RemoveEventResultResponse_Error_value, api.RemoveEventResultResponse_NOT_FOUND),
		}
		return nil
	}
	// Create event user request if not provided
	if c.In.User == nil {
		c.In.User = &api.EventUserRequest{
			Event: &api.EventRequest{},
		}
	}
	if c.In.User.Event == nil {
		c.In.User.Event = &api.EventRequest{}
	}
	result, err := c.service.database.DeleteEventRoundUser(ctx, model.GetEventRoundUserParams{
		EventUser: model.GetEventUserParams{
			Event: model.GetEventParams{
				ID:   conversion.Uint64ToSqlNullInt64(c.In.User.Event.Id),
				Name: conversion.StringToSqlNullString(c.In.User.Event.Name),
			},
			ID:     conversion.Uint64ToSqlNullInt64(c.In.User.Id),
			UserID: conversion.Uint64ToSqlNullInt64(c.In.User.UserId),
		},
		ID:    conversion.Uint64ToSqlNullInt64(c.In.Id),
		Round: conversion.StringToSqlNullString(c.In.Round),
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
		c.Out = &api.RemoveEventResultResponse{
			Success: false,
			Error:   api.RemoveEventResultResponse_NOT_FOUND,
		}
		return nil
	}
	c.Out = &api.RemoveEventResultResponse{
		Success: true,
		Error:   api.RemoveEventResultResponse_NONE,
	}
	return nil
}
