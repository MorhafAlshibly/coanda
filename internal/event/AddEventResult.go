package event

import (
	"context"
	"database/sql"
	"errors"

	"github.com/MorhafAlshibly/coanda/api"
	"github.com/MorhafAlshibly/coanda/internal/event/model"
	"github.com/MorhafAlshibly/coanda/pkg/conversion"
	"github.com/MorhafAlshibly/coanda/pkg/goquOptions"
)

type AddEventResultCommand struct {
	service *Service
	In      *api.AddEventResultRequest
	Out     *api.AddEventResultResponse
}

func NewAddEventResultCommand(service *Service, in *api.AddEventResultRequest) *AddEventResultCommand {
	return &AddEventResultCommand{
		service: service,
		In:      in,
	}
}

func (c *AddEventResultCommand) Execute(ctx context.Context) error {
	eErr := c.service.checkForEventRequestError(c.In.Event)
	if eErr != nil {
		c.Out = &api.AddEventResultResponse{
			Success: false,
			Error:   conversion.Enum(*eErr, api.AddEventResultResponse_Error_value, api.AddEventResultResponse_ID_OR_NAME_REQUIRED),
		}
		return nil
	}
	// Check if user is provided
	if c.In.ClientUserId == 0 {
		c.Out = &api.AddEventResultResponse{
			Success: false,
			Error:   api.AddEventResultResponse_CLIENT_USER_ID_REQUIRED,
		}
		return nil
	}
	// Check if result is provided
	if c.In.Result == 0 {
		c.Out = &api.AddEventResultResponse{
			Success: false,
			Error:   api.AddEventResultResponse_RESULT_REQUIRED,
		}
		return nil
	}
	// Check if user data is provided
	if c.In.UserData == nil {
		c.Out = &api.AddEventResultResponse{
			Success: false,
			Error:   api.AddEventResultResponse_USER_DATA_REQUIRED,
		}
		return nil
	}
	// Check if round user data is provided
	if c.In.RoundUserData == nil {
		c.Out = &api.AddEventResultResponse{
			Success: false,
			Error:   api.AddEventResultResponse_ROUND_USER_DATA_REQUIRED,
		}
		return nil
	}
	userData, err := conversion.ProtobufStructToRawJson(c.In.UserData)
	if err != nil {
		return err
	}
	tx, err := c.service.sql.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	qtx := c.service.database.WithTx(tx)
	// Check if event exists and lock it for update
	event, err := qtx.GetEvent(ctx, model.GetEventParams{
		ID:   conversion.Uint64ToSqlNullInt64(c.In.Event.Id),
		Name: conversion.StringToSqlNullString(c.In.Event.Name),
	}, &goquOptions.SelectDataset{Locked: true})
	if err != nil {
		// If event does not exist, return not found error
		if err == sql.ErrNoRows {
			c.Out = &api.AddEventResultResponse{
				Success: false,
				Error:   api.AddEventResultResponse_NOT_FOUND,
			}
			return nil
		}
		return err
	}
	// Get current event round and lock it for update
	eventRound, err := qtx.GetEventRound(ctx, model.GetEventRoundParams{
		Event: model.GetEventParams{
			ID: conversion.Uint64ToSqlNullInt64(&event.ID),
		},
	}, &goquOptions.SelectDataset{Locked: true})
	if err != nil {
		// If event round does not exist then the event has ended
		if err == sql.ErrNoRows {
			c.Out = &api.AddEventResultResponse{
				Success: false,
				Error:   api.AddEventResultResponse_EVENT_ENDED,
			}
			return nil
		}
		return err
	}
	// Check if event user exists and lock it for update
	// If it does not exist, it will be created in the next step
	eventUser, err := qtx.GetEventUser(ctx, model.GetEventUserParams{
		Event: model.GetEventParams{
			ID: conversion.Uint64ToSqlNullInt64(&event.ID),
		},
		ClientUserID: conversion.Uint64ToSqlNullInt64(&c.In.ClientUserId),
	}, &goquOptions.SelectDataset{Locked: true})
	var eventUserId int64
	if err != nil {
		// If event user does not exist, we will create it in the next step
		if err != sql.ErrNoRows {
			return err
		}
		eventUserResult, err := qtx.CreateEventUser(ctx, model.CreateEventUserParams{
			EventID:      event.ID,
			ClientUserID: c.In.ClientUserId,
			Data:         userData,
		})
		if err != nil {
			return err
		}
		// Get the last inserted ID of the event user
		eventUserId, err = eventUserResult.LastInsertId()
		if err != nil {
			return err
		}
	} else {
		// If event user exists, we will use its ID
		eventUserId = int64(eventUser.ID)
		// Update the event user data
		eventUserResult, err := qtx.UpdateEventUser(ctx, model.UpdateEventUserParams{
			User: model.GetEventUserWithoutWriteLockingParams{
				ID: conversion.Int64ToSqlNullInt64(&eventUserId),
			},
		})
		if err != nil {
			return err
		}
		rowsAffected, err := eventUserResult.RowsAffected()
		if err != nil {
			return err
		}
		if rowsAffected == 0 {
			// If no rows were affected, the event user was not updated
			return errors.New("event user was not updated")
		}
	}
	roundUserData, err := conversion.ProtobufStructToRawJson(c.In.RoundUserData)
	if err != nil {
		return err
	}
	// Get event round user and lock it for update, if it doesnt exist, it will be created in the next step
	_, err = qtx.GetEventRoundUserByEventUserIdAndRoundId(ctx, model.GetEventRoundUserByEventUserIdAndRoundIdParams{
		EventUserID:  uint64(eventUserId),
		EventRoundID: eventRound.ID,
	}, &goquOptions.SelectDataset{Locked: true})
	if err != nil {
		// Event round user does not exist
		if err != sql.ErrNoRows {
			return err
		}
		// Try to create event user round
		_, err = qtx.CreateEventRoundUser(ctx, model.CreateEventRoundUserParams{
			EventUserID:  uint64(eventUserId),
			EventRoundID: eventRound.ID,
			Result:       c.In.Result,
			Data:         roundUserData,
		})
		if err != nil {
			return err
		}
	} else {
		// Update the event round user
		eventRoundUserResult, err := qtx.UpdateEventRoundUser(ctx, model.UpdateEventRoundUserParams{
			EventUserID:  uint64(eventUserId),
			EventRoundID: eventRound.ID,
			Result:       c.In.Result,
			Data:         roundUserData,
		})
		if err != nil {
			return err
		}
		rowsAffected, err := eventRoundUserResult.RowsAffected()
		if err != nil {
			return err
		}
		if rowsAffected == 0 {
			return errors.New("event round user was not updated")
		}
	}
	// Commit the transaction
	err = tx.Commit()
	if err != nil {
		return err
	}
	c.Out = &api.AddEventResultResponse{
		Success: true,
		Error:   api.AddEventResultResponse_NONE,
	}
	return nil
}
