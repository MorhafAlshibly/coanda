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
		event, err := c.service.database.GetEventByName(ctx, *c.In.Event.Name)
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
	// Try to create event user
	eventUserResult, err := c.service.database.CreateEventUser(ctx, model.CreateEventUserParams{
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
	eventRoundUserResult, err := c.service.database.CreateEventRoundUser(ctx, model.CreateEventRoundUserParams{
		EventUserID: eventUserId,
		Result:      c.In.Result,
		Data:        roundUserData,
	})
	if err != nil {
		var mysqlErr *mysql.MySQLError
		if errors.As(err, &mysqlErr) {
			if mysqlErr.Number != errorcodes.MySQLErrorCodeDuplicateEntry {
				c.Out = &api.AddEventResultResponse{
					Success: false,
					Error:   api.AddEventResultResponse_ALREADY_EXISTS,
				}
				return nil
			}
		}
		return err
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
	c.Out = &api.AddEventResultResponse{
		Success: true,
		Error:   api.AddEventResultResponse_NONE,
	}
	return nil
}
