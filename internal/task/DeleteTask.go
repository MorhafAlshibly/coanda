package task

import (
	"context"

	"github.com/MorhafAlshibly/coanda/api"
	"github.com/MorhafAlshibly/coanda/internal/task/model"
	"github.com/MorhafAlshibly/coanda/pkg/conversion"
)

type DeleteTaskCommand struct {
	service *Service
	In      *api.TaskRequest
	Out     *api.TaskResponse
}

func NewDeleteTaskCommand(service *Service, in *api.TaskRequest) *DeleteTaskCommand {
	return &DeleteTaskCommand{
		service: service,
		In:      in,
	}
}

func (c *DeleteTaskCommand) Execute(ctx context.Context) error {
	iErr := c.service.checkForTaskRequestError(c.In)
	if iErr != nil {
		c.Out = &api.TaskResponse{
			Success: false,
			Error:   conversion.Enum(*iErr, api.TaskResponse_Error_value, api.TaskResponse_ID_REQUIRED),
		}
		return nil
	}
	result, err := c.service.database.DeleteTask(ctx, model.DeleteTaskParams{
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
		c.Out = &api.TaskResponse{
			Success: false,
			Error:   api.TaskResponse_NOT_FOUND,
		}
		return nil
	}
	c.Out = &api.TaskResponse{
		Success: true,
		Error:   api.TaskResponse_NONE,
	}
	return nil
}
