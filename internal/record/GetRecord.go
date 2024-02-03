package record

import (
	"context"
	"database/sql"
	"errors"

	"github.com/MorhafAlshibly/coanda/api"
	// "github.com/MorhafAlshibly/coanda/pkg/database/sqlc"
)

type GetRecordCommand struct {
	service *Service
	In      *api.GetRecordRequest
	Out     *api.GetRecordResponse
}

func NewGetRecordCommand(service *Service, in *api.GetRecordRequest) *GetRecordCommand {
	return &GetRecordCommand{
		service: service,
		In:      in,
	}
}

func (c *GetRecordCommand) Execute(ctx context.Context) error {
	if len(c.In.Name) < int(c.service.minRecordNameLength) {
		c.Out = &api.GetRecordResponse{
			Success: false,
			Record:  nil,
			Error:   api.GetRecordResponse_NAME_TOO_SHORT,
		}
		return nil
	}
	if len(c.In.Name) > int(c.service.maxRecordNameLength) {
		c.Out = &api.GetRecordResponse{
			Success: false,
			Record:  nil,
			Error:   api.GetRecordResponse_NAME_TOO_LONG,
		}
		return nil
	}
	if c.In.UserId == 0 {
		c.Out = &api.GetRecordResponse{
			Success: false,
			Record:  nil,
			Error:   api.GetRecordResponse_USER_ID_REQUIRED,
		}
		return nil
	}
	result, err := c.service.database.GetRecord(ctx, sqlc.GetRecordParams{
		Name:   c.In.Name,
		UserID: c.In.UserId,
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			c.Out = &api.GetRecordResponse{
				Success: false,
				Record:  nil,
				Error:   api.GetRecordResponse_NOT_FOUND,
			}
			return nil
		}
		return err
	}
	record, err := UnmarshalRecord(&result)
	if err != nil {
		return err
	}
	c.Out = &api.GetRecordResponse{
		Success: true,
		Record:  record,
		Error:   api.GetRecordResponse_NONE,
	}
	return nil
}
