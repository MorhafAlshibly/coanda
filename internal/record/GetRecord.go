package record

import (
	"context"
	"database/sql"
	"errors"

	"github.com/MorhafAlshibly/coanda/api"
	"github.com/MorhafAlshibly/coanda/internal/record/model"
	"github.com/MorhafAlshibly/coanda/pkg/conversion"
)

type GetRecordCommand struct {
	service *Service
	In      *api.RecordRequest
	Out     *api.GetRecordResponse
}

func NewGetRecordCommand(service *Service, in *api.RecordRequest) *GetRecordCommand {
	return &GetRecordCommand{
		service: service,
		In:      in,
	}
}

func (c *GetRecordCommand) Execute(ctx context.Context) error {
	rErr := c.service.checkForRecordRequestError(c.In)
	if rErr != nil {
		c.Out = &api.GetRecordResponse{
			Success: false,
			Error:   conversion.Enum(*rErr, api.GetRecordResponse_Error_value, api.GetRecordResponse_NOT_FOUND),
		}
		return nil
	}
	result, err := c.service.database.GetRecord(ctx, model.GetRecordParams{
		Name:   c.In.Name,
		UserID: c.In.UserId,
	})
	// Check if record is found
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
	// Unmarshal the record
	record, err := unmarshalRecord(&result)
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
