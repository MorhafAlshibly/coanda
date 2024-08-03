package event

import (
	"context"
	"database/sql"

	"github.com/MorhafAlshibly/coanda/api"
	"github.com/MorhafAlshibly/coanda/internal/event/model"
	"github.com/MorhafAlshibly/coanda/pkg/conversion"
)

type GetEventUserCommand struct {
	service *Service
	In      *api.GetEventUserRequest
	Out     *api.GetEventUserResponse
}

func NewGetEventUserCommand(service *Service, in *api.GetEventUserRequest) *GetEventUserCommand {
	return &GetEventUserCommand{
		service: service,
		In:      in,
	}
}

func (c *GetEventUserCommand) Execute(ctx context.Context) error {
	euErr := c.service.checkForEventUserRequestError(c.In.User)
	// Check if error is found
	if euErr != nil {
		c.Out = &api.GetEventUserResponse{
			Success: false,
			Error:   conversion.Enum(*euErr, api.GetEventUserResponse_Error_value, api.GetEventUserResponse_NOT_FOUND),
		}
		return nil
	}
	// Create blank event request if not provided
	if c.In.User.Event == nil {
		c.In.User.Event = &api.EventRequest{}
	}
	limit, offset := conversion.PaginationToLimitOffset(c.In.Pagination, c.service.defaultMaxPageLength, c.service.maxMaxPageLength)
	eventUser, err := c.service.database.GetEventUser(ctx, model.GetEventUserParams{
		Event: model.GetEventParams{
			ID:   conversion.Uint64ToSqlNullInt64(c.In.User.Event.Id),
			Name: conversion.StringToSqlNullString(c.In.User.Event.Name),
		},
		ID:     conversion.Uint64ToSqlNullInt64(c.In.User.Id),
		UserID: conversion.Uint64ToSqlNullInt64(c.In.User.UserId),
	})
	if err != nil {
		if err == sql.ErrNoRows {
			c.Out = &api.GetEventUserResponse{
				Success: false,
				Error:   api.GetEventUserResponse_NOT_FOUND,
			}
			return nil
		}
		return err
	}
	apiEventUser, err := unmarshalEventUser(eventUser)
	if err != nil {
		return err
	}
	eventRoundUsers, err := c.service.database.GetEventRoundUsers(ctx, model.GetEventRoundUsersParams{
		EventUser: model.GetEventUserParams{
			Event: model.GetEventParams{
				ID:   conversion.Uint64ToSqlNullInt64(c.In.User.Event.Id),
				Name: conversion.StringToSqlNullString(c.In.User.Event.Name),
			},
			ID:     conversion.Uint64ToSqlNullInt64(c.In.User.Id),
			UserID: conversion.Uint64ToSqlNullInt64(c.In.User.Id),
		},
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		if err == sql.ErrNoRows {
			eventRoundUsers = []model.EventRoundLeaderboard{}
		} else {
			return err
		}
	}
	apiEventRoundUsers, err := unmarshalEventRoundLeaderboard(eventRoundUsers)
	if err != nil {
		return err
	}
	c.Out = &api.GetEventUserResponse{
		Success: true,
		User:    apiEventUser,
		Results: apiEventRoundUsers,
		Error:   api.GetEventUserResponse_NONE,
	}
	return nil
}
