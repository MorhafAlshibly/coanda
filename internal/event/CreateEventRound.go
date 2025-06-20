package event

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/MorhafAlshibly/coanda/api"
	"github.com/MorhafAlshibly/coanda/internal/event/model"
	"github.com/MorhafAlshibly/coanda/pkg/conversion"
	"github.com/MorhafAlshibly/coanda/pkg/errorcode"
	"github.com/MorhafAlshibly/coanda/pkg/goquOptions"
	"github.com/go-sql-driver/mysql"
)

type CreateEventRoundCommand struct {
	service *Service
	In      *api.CreateEventRoundRequest
	Out     *api.CreateEventRoundResponse
}

func NewCreateEventRoundCommand(service *Service, in *api.CreateEventRoundRequest) *CreateEventRoundCommand {
	return &CreateEventRoundCommand{
		service: service,
		In:      in,
	}
}

func (c *CreateEventRoundCommand) Execute(ctx context.Context) error {
	eErr := c.service.checkForEventRequestError(c.In.Event)
	// Check if error is found
	if eErr != nil {
		c.Out = &api.CreateEventRoundResponse{
			Success: false,
			Error:   conversion.Enum(*eErr, api.CreateEventRoundResponse_Error_value, api.CreateEventRoundResponse_ID_OR_NAME_REQUIRED),
		}
		return nil
	}
	// Check if no round is specified
	if c.In.Round == nil {
		c.Out = &api.CreateEventRoundResponse{
			Success: false,
			Error:   api.CreateEventRoundResponse_ROUND_REQUIRED,
		}
		return nil
	}
	// Check if round name is too short
	if len(c.In.Round.Name) < int(c.service.minRoundNameLength) {
		c.Out = &api.CreateEventRoundResponse{
			Success: false,
			Error:   api.CreateEventRoundResponse_ROUND_NAME_TOO_SHORT,
		}
		return nil
	}
	// Check if round name is too long
	if len(c.In.Round.Name) > int(c.service.maxRoundNameLength) {
		c.Out = &api.CreateEventRoundResponse{
			Success: false,
			Error:   api.CreateEventRoundResponse_ROUND_NAME_TOO_LONG,
		}
		return nil
	}
	// Check if round data is specified
	if c.In.Round.Data == nil {
		c.Out = &api.CreateEventRoundResponse{
			Success: false,
			Error:   api.CreateEventRoundResponse_ROUND_DATA_REQUIRED,
		}
		return nil
	}
	// Check if round ended at is specified
	if c.In.Round.EndedAt == nil {
		c.Out = &api.CreateEventRoundResponse{
			Success: false,
			Error:   api.CreateEventRoundResponse_ROUND_ENDED_AT_REQUIRED,
		}
		return nil
	}
	// Check if round ended at is in the past
	if c.In.Round.EndedAt.AsTime().Before(time.Now()) {
		c.Out = &api.CreateEventRoundResponse{
			Success: false,
			Error:   api.CreateEventRoundResponse_ROUND_ENDED_AT_IN_THE_PAST,
		}
		return nil
	}
	tx, err := c.service.sql.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	qtx := c.service.database.WithTx(tx)
	// Get the event
	event, err := qtx.GetEvent(ctx, model.GetEventParams{
		ID:   conversion.Uint64ToSqlNullInt64(c.In.Event.Id),
		Name: conversion.StringToSqlNullString(c.In.Event.Name),
	}, &goquOptions.SelectDataset{Locked: true})
	if err != nil {
		// If the event does not exist, return an error
		if err == sql.ErrNoRows {
			c.Out = &api.CreateEventRoundResponse{
				Success: false,
				Error:   api.CreateEventRoundResponse_NOT_FOUND,
			}
			return nil
		}
	}
	// Check if round ended at is after event ended at
	if c.In.Round.EndedAt.AsTime().Before(event.StartedAt) {
		c.Out = &api.CreateEventRoundResponse{
			Success: false,
			Error:   api.CreateEventRoundResponse_ROUND_ENDED_AT_BEFORE_EVENT_STARTED_AT,
		}
		return nil
	}
	// Check if round scoring is specified
	if len(c.In.Round.Scoring) == 0 {
		c.Out = &api.CreateEventRoundResponse{
			Success: false,
			Error:   api.CreateEventRoundResponse_ROUND_SCORING_REQUIRED,
		}
		return nil
	}
	data, err := conversion.ProtobufStructToRawJson(c.In.Round.Data)
	if err != nil {
		return err
	}
	scoring, err := conversion.MapToRawJson(map[string]interface{}{"scoring": c.In.Round.Scoring})
	if err != nil {
		return err
	}
	// Create event round
	result, err := qtx.CreateEventRound(ctx, model.CreateEventRoundParams{
		EventID: event.ID,
		Name:    c.In.Round.Name,
		Data:    data,
		Scoring: scoring,
		EndedAt: c.In.Round.EndedAt.AsTime(),
	})
	// If the round already exists, return an error
	if err != nil {
		var mysqlErr *mysql.MySQLError
		if errors.As(err, &mysqlErr) {
			if errorcode.IsDuplicateEntry(mysqlErr, "event_round", "event_round_name_event_id_idx") {
				c.Out = &api.CreateEventRoundResponse{
					Success: false,
					Error:   api.CreateEventRoundResponse_DUPLICATE_ROUND_NAME,
				}
				return nil
			}
			if errorcode.IsDuplicateEntry(mysqlErr, "event_round", "event_round_event_id_ended_at_idx") {
				c.Out = &api.CreateEventRoundResponse{
					Success: false,
					Error:   api.CreateEventRoundResponse_DUPLICATE_ROUND_ENDED_AT,
				}
				return nil
			}
		}
		return err
	}
	// Get inserted ID
	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	uint64Id := uint64(id)
	err = tx.Commit()
	if err != nil {
		return err
	}
	c.Out = &api.CreateEventRoundResponse{
		Success: true,
		Id:      &uint64Id,
		Error:   api.CreateEventRoundResponse_NONE,
	}
	return nil
}
