package event

import (
	"context"
	"errors"

	"github.com/MorhafAlshibly/coanda/api"
	"github.com/MorhafAlshibly/coanda/internal/event/model"
	"github.com/MorhafAlshibly/coanda/pkg/conversion"
	"github.com/MorhafAlshibly/coanda/pkg/errorcodes"
	"github.com/go-sql-driver/mysql"
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
	if c.In.UserId == 0 {
		c.Out = &api.AddEventResultResponse{
			Success: false,
			Error:   api.AddEventResultResponse_USER_ID_REQUIRED,
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
	var eventId uint64
	// If ID is null, try to get event by name
	if c.In.Event.Id == nil {
		event, err := c.service.database.GetEvent(ctx, model.GetEventParams{
			Name: conversion.StringToSqlNullString(c.In.Event.Name),
		})
		if err != nil {
			c.Out = &api.AddEventResultResponse{
				Success: false,
				Error:   api.AddEventResultResponse_NOT_FOUND,
			}
			return nil
		}
		eventId = event.ID
	} else {
		eventId = *c.In.Event.Id
	}
	tx, err := c.service.sql.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	qtx := c.service.database.WithTx(tx)
	// Try to create event user
	eventUserResult, err := qtx.CreateEventUser(ctx, model.CreateEventUserParams{
		EventID: eventId,
		UserID:  c.In.UserId,
		Data:    userData,
	})
	var eventUserId uint64
	if err != nil {
		// If the event user already exists, we can ignore the error
		var mysqlErr *mysql.MySQLError
		if errors.As(err, &mysqlErr) {
			if mysqlErr.Number != errorcodes.MySQLErrorCodeDuplicateEntry {
				return err
			}
		} else {
			return err
		}
		// Get the existing event user
		eventUser, err := c.service.database.GetEventUserByEventIdAndUserId(ctx, model.GetEventUserByEventIdAndUserIdParams{
			EventID: eventId,
			UserID:  c.In.UserId,
		})
		if err != nil {
			return err
		}
		eventUserId = eventUser.ID
		// Update the existing event user
		_, err = qtx.UpdateEventUser(ctx, model.UpdateEventUserParams{
			User: model.GetEventUserParams{
				ID: conversion.Uint64ToSqlNullInt64(&eventUserId),
			},
			Data: userData,
		})
		if err != nil {
			return err
		}

	} else {
		lastInsertedEventId, err := eventUserResult.LastInsertId()
		if err != nil {
			return err
		}
		eventUserId = uint64(lastInsertedEventId)
	}
	roundUserData, err := conversion.ProtobufStructToRawJson(c.In.RoundUserData)
	if err != nil {
		return err
	}
	// Try to create event user round
	eventRoundUserResult, err := qtx.CreateEventRoundUser(ctx, model.CreateEventRoundUserParams{
		EventUserID: eventUserId,
		Result:      c.In.Result,
		Data:        roundUserData,
	})
	if err != nil {
		var mysqlErr *mysql.MySQLError
		if errors.As(err, &mysqlErr) {
			if mysqlErr.Number != errorcodes.MySQLErrorCodeDuplicateEntry {
				return err
			}
		} else {
			return err
		}
		// If the event round user already exists, we can ignore the error and update the existing one
		updateEventRoundUserResultResult, err := qtx.UpdateEventRoundUserResult(ctx, model.UpdateEventRoundUserResultParams{
			EventUserID: eventUserId,
			Result:      c.In.Result,
			Data:        roundUserData,
		})
		if err != nil {
			return err
		}
		rowsAffected, err := updateEventRoundUserResultResult.RowsAffected()
		if err != nil {
			return err
		}
		if rowsAffected == 0 {
			// If no rows were affected, the event has either ended or the result is the same
			// Check if the result is the same
			eventRoundUser, err := c.service.database.GetEventRoundUserByEventUserId(ctx, eventUserId)
			if err != nil {
				return err
			}
			if eventRoundUser.Result == c.In.Result {
				c.Out = &api.AddEventResultResponse{
					Success: true,
					Error:   api.AddEventResultResponse_NONE,
				}
				return nil
			}
			// If the result is different, the round has ended or the event has been deleted
			return errors.New("event round user not found, unexpected error occurred")
		}
		// If rows were affected, the result was updated
		c.Out = &api.AddEventResultResponse{
			Success: true,
			Error:   api.AddEventResultResponse_NONE,
		}
		return nil

	}
	// If no round user was not created, the event has already ended
	eventRoundUserRowsAffected, err := eventRoundUserResult.RowsAffected()
	if err != nil {
		return err
	}
	if eventRoundUserRowsAffected == 0 {
		c.Out = &api.AddEventResultResponse{
			Success: false,
			Error:   api.AddEventResultResponse_EVENT_ENDED,
		}
		return nil
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
