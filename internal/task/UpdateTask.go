package task

import (
	"context"
	"database/sql"

	"github.com/MorhafAlshibly/coanda/api"
	"github.com/MorhafAlshibly/coanda/internal/task/model"
	"github.com/MorhafAlshibly/coanda/pkg/conversion"
)

type UpdateTaskCommand struct {
	service *Service
	In      *api.UpdateTaskRequest
	Out     *api.UpdateTaskResponse
}

func NewUpdateTaskCommand(service *Service, in *api.UpdateTaskRequest) *UpdateTaskCommand {
	return &UpdateTaskCommand{
		service: service,
		In:      in,
	}
}

func (c *UpdateTaskCommand) Execute(ctx context.Context) error {
	iErr := c.service.checkForTaskRequestError(c.In.Task)
	if iErr != nil {
		c.Out = &api.UpdateTaskResponse{
			Success: false,
			Error:   conversion.Enum(*iErr, api.UpdateTaskResponse_Error_value, api.UpdateTaskResponse_ID_REQUIRED),
		}
		return nil
	}
	if c.In.Data == nil {
		c.Out = &api.UpdateTaskResponse{
			Success: false,
			Error:   api.UpdateTaskResponse_DATA_REQUIRED,
		}
		return nil
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
	result, err := qtx.UpdateTask(ctx, model.UpdateTaskParams{
		ID:   c.In.Task.Id,
		Type: c.In.Task.Type,
		Data: data,
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
		_, err := qtx.GetTask(ctx, model.GetTaskParams{
			ID:   c.In.Task.Id,
			Type: c.In.Task.Type,
		})
		if err != nil {
			if err == sql.ErrNoRows {
				c.Out = &api.UpdateTaskResponse{
					Success: false,
					Error:   api.UpdateTaskResponse_NOT_FOUND,
				}
				return nil
			}
			return err
		}
	}
	c.Out = &api.UpdateTaskResponse{
		Success: true,
		Error:   api.UpdateTaskResponse_NONE,
	}
	return nil
}
