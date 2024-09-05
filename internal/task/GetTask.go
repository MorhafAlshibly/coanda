package task

import (
	"context"
	"database/sql"
	"errors"

	"github.com/MorhafAlshibly/coanda/api"
	"github.com/MorhafAlshibly/coanda/internal/task/model"
	"github.com/MorhafAlshibly/coanda/pkg/conversion"
)

type GetTaskCommand struct {
	service *Service
	In      *api.TaskRequest
	Out     *api.GetTaskResponse
}

func NewGetTaskCommand(service *Service, in *api.TaskRequest) *GetTaskCommand {
	return &GetTaskCommand{
		service: service,
		In:      in,
	}
}

func (c *GetTaskCommand) Execute(ctx context.Context) error {
	iErr := c.service.checkForTaskRequestError(c.In)
	if iErr != nil {
		c.Out = &api.GetTaskResponse{
			Success: false,
			Error:   conversion.Enum(*iErr, api.GetTaskResponse_Error_value, api.GetTaskResponse_NOT_FOUND),
		}
		return nil
	}
	result, err := c.service.database.GetTask(ctx, model.GetTaskParams{
		ID:   c.In.Id,
		Type: c.In.Type,
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			c.Out = &api.GetTaskResponse{
				Success: false,
				Task:    nil,
				Error:   api.GetTaskResponse_NOT_FOUND,
			}
			return nil
		}
		return err
	}
	task, err := unmarshalTask(&result)
	if err != nil {
		return err
	}
	c.Out = &api.GetTaskResponse{
		Success: true,
		Task:    task,
		Error:   api.GetTaskResponse_NONE,
	}
	return nil
}
