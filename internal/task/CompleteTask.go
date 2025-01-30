package task

import (
	"context"
	"database/sql"
	"errors"

	"github.com/MorhafAlshibly/coanda/api"
	"github.com/MorhafAlshibly/coanda/internal/task/model"
	"github.com/MorhafAlshibly/coanda/pkg/conversion"
)

type CompleteTaskCommand struct {
	service *Service
	In      *api.TaskRequest
	Out     *api.CompleteTaskResponse
}

func NewCompleteTaskCommand(service *Service, in *api.TaskRequest) *CompleteTaskCommand {
	return &CompleteTaskCommand{
		service: service,
		In:      in,
	}
}

func (c *CompleteTaskCommand) Execute(ctx context.Context) error {
	iErr := c.service.checkForTaskRequestError(c.In)
	if iErr != nil {
		c.Out = &api.CompleteTaskResponse{
			Success: false,
			Error:   conversion.Enum(*iErr, api.CompleteTaskResponse_Error_value, api.CompleteTaskResponse_ID_REQUIRED),
		}
		return nil
	}
	result, err := c.service.database.CompleteTask(ctx, model.CompleteTaskParams{
		ID:   c.In.Id,
		Type: c.In.Type,
	})
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		// Check if task is found
		_, err := c.service.database.GetTask(ctx, model.GetTaskParams{
			ID:   c.In.Id,
			Type: c.In.Type,
		})
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				c.Out = &api.CompleteTaskResponse{
					Success: false,
					Error:   api.CompleteTaskResponse_NOT_FOUND,
				}
				return nil
			}
			return err
		}
		c.Out = &api.CompleteTaskResponse{
			Success: false,
			Error:   api.CompleteTaskResponse_ALREADY_COMPLETED,
		}
		return nil
	}
	c.Out = &api.CompleteTaskResponse{
		Success: true,
		Error:   api.CompleteTaskResponse_NONE,
	}
	return nil
}
