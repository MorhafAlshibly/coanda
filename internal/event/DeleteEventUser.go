package event

import (
	"context"
	"database/sql"

	"github.com/MorhafAlshibly/coanda/api"
	"github.com/MorhafAlshibly/coanda/internal/event/model"
	"github.com/MorhafAlshibly/coanda/pkg/conversion"
	"github.com/MorhafAlshibly/coanda/pkg/goquOptions"
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
			Error:   conversion.Enum(*eErr, api.EventUserResponse_Error_value, api.EventUserResponse_ID_OR_NAME_REQUIRED),
		}
		return nil
	}
	// Create empty event request if not provided
	if c.In.Event == nil {
		c.In.Event = &api.EventRequest{}
	}
	tx, err := c.service.sql.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	qtx := c.service.database.WithTx(tx)
	// Check if event name is provided
	if c.In.Event.Name != nil && c.In.Event.Id == nil {
		event, err := qtx.GetEvent(ctx, model.GetEventParams{
			Name: conversion.StringToSqlNullString(c.In.Event.Name),
		}, &goquOptions.SelectDataset{Locked: true})
		if err != nil {
			if err == sql.ErrNoRows {
				c.Out = &api.EventUserResponse{
					Success: false,
					Error:   api.EventUserResponse_NOT_FOUND,
				}
				return nil
			}
			return err
		}
		c.In.Event.Id = &event.ID
	}
	result, err := qtx.DeleteEventUser(ctx, model.GetEventUserWithoutWriteLockingParams{
		EventID:      conversion.Uint64ToSqlNullInt64(c.In.Event.Id),
		ID:           conversion.Uint64ToSqlNullInt64(c.In.Id),
		ClientUserID: conversion.Uint64ToSqlNullInt64(c.In.ClientUserId),
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
	err = tx.Commit()
	if err != nil {
		return err
	}
	c.Out = &api.EventUserResponse{
		Success: true,
		Error:   api.EventUserResponse_NONE,
	}
	return nil
}
