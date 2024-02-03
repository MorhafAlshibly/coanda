package record

import (
	"context"
	"database/sql"
	"errors"

	"github.com/MorhafAlshibly/coanda/api"
	// "github.com/MorhafAlshibly/coanda/pkg/database/sqlc"
)

type DeleteRecordCommand struct {
	service *Service
	In      *api.GetRecordRequest
	Out     *api.DeleteRecordResponse
}

func NewDeleteRecordCommand(service *Service, in *api.GetRecordRequest) *DeleteRecordCommand {
	return &DeleteRecordCommand{
		service: service,
		In:      in,
	}
}

func (c *DeleteRecordCommand) Execute(ctx context.Context) error {
	if len(c.In.Name) < int(c.service.minRecordNameLength) {
		c.Out = &api.DeleteRecordResponse{
			Success: false,
			Error:   api.DeleteRecordResponse_NAME_TOO_SHORT,
		}
		return nil
	}
	if len(c.In.Name) > int(c.service.maxRecordNameLength) {
		c.Out = &api.DeleteRecordResponse{
			Success: false,
			Error:   api.DeleteRecordResponse_NAME_TOO_LONG,
		}
		return nil
	}
	err := c.service.database.DeleteRecord(ctx, sqlc.DeleteRecordParams{
		Name:   c.In.Name,
		UserID: c.In.UserId,
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			c.Out = &api.DeleteRecordResponse{
				Success: false,
				Error:   api.DeleteRecordResponse_NOT_FOUND,
			}
			return nil
		}
		return err
	}
	c.Out = &api.DeleteRecordResponse{
		Success: true,
		Error:   api.DeleteRecordResponse_NONE,
	}
	return nil
}
