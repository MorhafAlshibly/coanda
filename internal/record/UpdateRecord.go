package record

import (
	"context"
	"database/sql"
	"errors"

	"github.com/MorhafAlshibly/coanda/api"
	"github.com/MorhafAlshibly/coanda/pkg/conversion"
	"github.com/MorhafAlshibly/coanda/pkg/database/sqlc"
)

type UpdateRecordCommand struct {
	service *Service
	In      *api.UpdateRecordRequest
	Out     *api.UpdateRecordResponse
}

func NewUpdateRecordCommand(service *Service, in *api.UpdateRecordRequest) *UpdateRecordCommand {
	return &UpdateRecordCommand{
		service: service,
		In:      in,
	}
}

func (c *UpdateRecordCommand) Execute(ctx context.Context) error {
	if len(c.In.Request.Name) < int(c.service.minRecordNameLength) {
		c.Out = &api.UpdateRecordResponse{
			Success: false,
			Error:   api.UpdateRecordResponse_NAME_TOO_SHORT,
		}
		return nil
	}
	if len(c.In.Request.Name) > int(c.service.maxRecordNameLength) {
		c.Out = &api.UpdateRecordResponse{
			Success: false,
			Error:   api.UpdateRecordResponse_NAME_TOO_LONG,
		}
		return nil
	}
	data, err := conversion.ProtobufStructToRawJson(c.In.Data)
	if err != nil {
		return err
	}
	// Update the item in the store
	if c.In.Record != nil && c.In.Data != nil {
		err = c.service.database.UpdateRecord(ctx, sqlc.UpdateRecordParams{
			Name:   c.In.Request.Name,
			UserID: c.In.Request.UserId,
			Data:   data,
			Record: *c.In.Record,
		})
	} else if c.In.Record != nil {
		err = c.service.database.UpdateRecordRecord(ctx, sqlc.UpdateRecordRecordParams{
			Name:   c.In.Request.Name,
			UserID: c.In.Request.UserId,
			Record: *c.In.Record,
		})
	} else if c.In.Data != nil {
		err = c.service.database.UpdateRecordData(ctx, sqlc.UpdateRecordDataParams{
			Name:   c.In.Request.Name,
			UserID: c.In.Request.UserId,
			Data:   data,
		})
	} else {
		c.Out = &api.UpdateRecordResponse{
			Success: false,
			Error:   api.UpdateRecordResponse_INVALID,
		}
		return nil
	}
	// Record not found
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			c.Out = &api.UpdateRecordResponse{
				Success: false,
				Error:   api.UpdateRecordResponse_NOT_FOUND,
			}
			return nil
		}
		return err
	}
	c.Out = &api.UpdateRecordResponse{
		Success: true,
		Error:   api.UpdateRecordResponse_NONE,
	}
	return nil
}
