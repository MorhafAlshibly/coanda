package event

import (
	"context"
	"errors"
	"time"

	"github.com/MorhafAlshibly/coanda/api"
	"github.com/MorhafAlshibly/coanda/internal/event/model"
	"github.com/MorhafAlshibly/coanda/pkg/conversion"
	"github.com/MorhafAlshibly/coanda/pkg/errorcodes"
	"github.com/go-sql-driver/mysql"
)

type CreateEventCommand struct {
	service *Service
	In      *api.CreateEventRequest
	Out     *api.CreateEventResponse
}

func NewCreateEventCommand(service *Service, in *api.CreateEventRequest) *CreateEventCommand {
	return &CreateEventCommand{
		service: service,
		In:      in,
	}
}

func (c *CreateEventCommand) Execute(ctx context.Context) error {
	// Check if event name is large enough
	if len(c.In.Name) < int(c.service.minEventNameLength) {
		c.Out = &api.CreateEventResponse{
			Success: false,
			Error:   api.CreateEventResponse_NAME_TOO_SHORT,
		}
		return nil
	}
	// Check if event name is small enough
	if len(c.In.Name) > int(c.service.maxEventNameLength) {
		c.Out = &api.CreateEventResponse{
			Success: false,
			Error:   api.CreateEventResponse_NAME_TOO_LONG,
		}
		return nil
	}
	// Check if data is provided
	if c.In.Data == nil {
		c.Out = &api.CreateEventResponse{
			Success: false,
			Error:   api.CreateEventResponse_DATA_REQUIRED,
		}
		return nil
	}
	// Check if event started at is provided
	if c.In.StartedAt == nil {
		c.Out = &api.CreateEventResponse{
			Success: false,
			Error:   api.CreateEventResponse_STARTED_AT_REQUIRED,
		}
		return nil
	}
	// Check if event started at is in the future
	if c.In.StartedAt.AsTime().Before(time.Now()) {
		c.Out = &api.CreateEventResponse{
			Success: false,
			Error:   api.CreateEventResponse_STARTED_AT_IN_THE_PAST,
		}
		return nil
	}
	// Check if we have any rounds
	if len(c.In.Rounds) == 0 {
		c.Out = &api.CreateEventResponse{
			Success: false,
			Error:   api.CreateEventResponse_ROUNDS_REQUIRED,
		}
		return nil
	}
	// Check if we have too many rounds
	if len(c.In.Rounds) > int(c.service.maxNumberOfRounds) {
		c.Out = &api.CreateEventResponse{
			Success: false,
			Error:   api.CreateEventResponse_TOO_MANY_ROUNDS,
		}
		return nil
	}
	for _, round := range c.In.Rounds {
		// Check if round name is large enough
		if len(round.Name) < int(c.service.minRoundNameLength) {
			c.Out = &api.CreateEventResponse{
				Success: false,
				Error:   api.CreateEventResponse_ROUND_NAME_TOO_SHORT,
			}
			return nil
		}
		// Check if round name is small enough
		if len(round.Name) > int(c.service.maxRoundNameLength) {
			c.Out = &api.CreateEventResponse{
				Success: false,
				Error:   api.CreateEventResponse_ROUND_NAME_TOO_LONG,
			}
			return nil
		}
		// Check if round data is provided
		if round.Data == nil {
			c.Out = &api.CreateEventResponse{
				Success: false,
				Error:   api.CreateEventResponse_ROUND_DATA_REQUIRED,
			}
			return nil
		}
		// Check if round ended at is provided
		if round.EndedAt == nil {
			c.Out = &api.CreateEventResponse{
				Success: false,
				Error:   api.CreateEventResponse_ROUND_ENDED_AT_REQUIRED,
			}
			return nil
		}
		// Check if round ended at is after event started at
		if round.EndedAt.AsTime().Before(c.In.StartedAt.AsTime()) {
			c.Out = &api.CreateEventResponse{
				Success: false,
				Error:   api.CreateEventResponse_ROUND_ENDED_AT_BEFORE_STARTED_AT,
			}
			return nil
		}
		// Check if round scoring is provided
		if len(round.Scoring) == 0 {
			c.Out = &api.CreateEventResponse{
				Success: false,
				Error:   api.CreateEventResponse_ROUND_SCORING_REQUIRED,
			}
			return nil
		}
	}
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
	// Create the event
	eventResult, err := qtx.CreateEvent(ctx, model.CreateEventParams{
		Name:      c.In.Name,
		Data:      data,
		StartedAt: c.In.StartedAt.AsTime(),
	})
	// If the event already exists, return an error
	if err != nil {
		var mysqlErr *mysql.MySQLError
		if errors.As(err, &mysqlErr) {
			if errorcodes.IsDuplicateEntry(mysqlErr, "event", "name") {
				c.Out = &api.CreateEventResponse{
					Success: false,
					Error:   api.CreateEventResponse_ALREADY_EXISTS,
				}
				return nil
			}
		}
		return err
	}
	// Get the event id
	eventID, err := eventResult.LastInsertId()
	if err != nil {
		return err
	}
	// Create the rounds
	for _, round := range c.In.Rounds {
		roundData, err := conversion.ProtobufStructToRawJson(round.Data)
		if err != nil {
			return err
		}
		// Convert scoring to json
		roundScoring, err := conversion.MapToRawJson(map[string]interface{}{"scoring": round.Scoring})
		if err != nil {
			return err
		}
		_, err = qtx.CreateEventRound(ctx, model.CreateEventRoundParams{
			EventID: uint64(eventID),
			Name:    round.Name,
			Data:    roundData,
			Scoring: roundScoring,
			EndedAt: round.EndedAt.AsTime(),
		})
		// If the round already exists, return an error
		if err != nil {
			var mysqlErr *mysql.MySQLError
			if errors.As(err, &mysqlErr) {
				if errorcodes.IsDuplicateEntry(mysqlErr, "event_round", "event_round_name_event_id_idx") {
					c.Out = &api.CreateEventResponse{
						Success: false,
						Error:   api.CreateEventResponse_DUPLICATE_ROUND_NAME,
					}
					return nil
				}
				if errorcodes.IsDuplicateEntry(mysqlErr, "event_round", "event_round_ended_at_event_id_idx") {
					c.Out = &api.CreateEventResponse{
						Success: false,
						Error:   api.CreateEventResponse_DUPLICATE_ROUND_ENDED_AT,
					}
					return nil
				}
			}
			return err
		}
	}
	// Commit the transaction
	err = tx.Commit()
	if err != nil {
		return err
	}
	uint64EventID := uint64(eventID)
	c.Out = &api.CreateEventResponse{
		Success: true,
		Id:      &uint64EventID,
		Error:   api.CreateEventResponse_NONE,
	}
	return nil
}
