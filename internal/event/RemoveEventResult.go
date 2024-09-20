package event

import (
	"context"

	"github.com/MorhafAlshibly/coanda/api"
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
	if c.In == nil {
		c.Out = &api.RemoveEventResultResponse{
			Success: false,
			Error:   api.RemoveEventResultResponse_ID_REQUIRED,
		}
		return nil
	}
	if c.In.Id == 0 {
		c.Out = &api.RemoveEventResultResponse{
			Success: false,
			Error:   api.RemoveEventResultResponse_ID_REQUIRED,
		}
		return nil
	}
	result, err := c.service.database.DeleteEventRoundUser(ctx, c.In.Id)
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
