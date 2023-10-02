package record

import (
	"context"

	"github.com/MorhafAlshibly/coanda/api"
)

type DeleteRecordCommand struct {
	service *Service
	In      *api.DeleteRecordRequest
	Out     *api.DeleteRecordResponse
}

func NewDeleteRecordCommand(service *Service, in *api.DeleteRecordRequest) *DeleteRecordCommand {
	return &DeleteRecordCommand{
		service: service,
		In:      in,
	}
}

func (c *DeleteRecordCommand) Execute(ctx context.Context) error {
	filter, err := getFilter(c.In.Record)
	if err != nil {
		c.Out = &api.DeleteRecordResponse{
			Success: false,
			Error:   api.DeleteRecordResponse_INVALID,
		}
		return nil
	}
	result, writeErr := c.service.db.DeleteOne(ctx, filter)
	if writeErr != nil {
		return writeErr
	}
	if result.DeletedCount == 0 {
		c.Out = &api.DeleteRecordResponse{
			Success: false,
			Error:   api.DeleteRecordResponse_NOT_FOUND,
		}
		return nil
	}
	c.Out = &api.DeleteRecordResponse{
		Success: true,
		Error:   api.DeleteRecordResponse_NONE,
	}
	return nil
}
