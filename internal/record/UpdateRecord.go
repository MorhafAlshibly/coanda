package record

import (
	"context"

	"github.com/MorhafAlshibly/coanda/api"
	"github.com/MorhafAlshibly/coanda/internal/record/model"
	"github.com/MorhafAlshibly/coanda/pkg/conversion"
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
	rErr := c.service.CheckForRecordRequestError(c.In.Request)
	if rErr != nil {
		c.Out = &api.UpdateRecordResponse{
			Success: false,
			Error:   conversion.Enum(*rErr, api.UpdateRecordResponse_Error_value, api.UpdateRecordResponse_NOT_FOUND),
		}
		return nil
	}
	if c.In.Record == nil && c.In.Data == nil {
		c.Out = &api.UpdateRecordResponse{
			Success: false,
			Error:   api.UpdateRecordResponse_NO_UPDATE_SPECIFIED,
		}
		return nil
	}
	// Update the record in the store
	tx, err := c.service.sql.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	if c.In.Record != nil {
		result, err := c.service.database.UpdateRecordRecord(ctx, model.UpdateRecordRecordParams{
			Name:   c.In.Request.Name,
			UserID: c.In.Request.UserId,
			Record: *c.In.Record,
		})
		if err != nil {
			return err
		}
		rowsAffected, err := result.RowsAffected()
		if err != nil {
			return err
		}
		if rowsAffected == 0 {
			c.Out = &api.UpdateRecordResponse{
				Success: false,
				Error:   api.UpdateRecordResponse_NOT_FOUND,
			}
			return nil
		}
	}
	if c.In.Data != nil {
		data, err := conversion.ProtobufStructToRawJson(c.In.Data)
		if err != nil {
			return err
		}
		result, err := c.service.database.UpdateRecordData(ctx, model.UpdateRecordDataParams{
			Name:   c.In.Request.Name,
			UserID: c.In.Request.UserId,
			Data:   data,
		})
		if err != nil {
			return err
		}
		rowsAffected, err := result.RowsAffected()
		if err != nil {
			return err
		}
		if rowsAffected == 0 {
			c.Out = &api.UpdateRecordResponse{
				Success: false,
				Error:   api.UpdateRecordResponse_NOT_FOUND,
			}
			return nil
		}
	}
	err = tx.Commit()
	if err != nil {
		return err
	}
	c.Out = &api.UpdateRecordResponse{
		Success: true,
		Error:   api.UpdateRecordResponse_NONE,
	}
	return nil
}
