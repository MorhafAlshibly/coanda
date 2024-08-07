package record

import (
	"context"

	"github.com/MorhafAlshibly/coanda/api"
	"github.com/MorhafAlshibly/coanda/internal/record/model"
	"github.com/MorhafAlshibly/coanda/pkg/conversion"
)

type DeleteRecordCommand struct {
	service *Service
	In      *api.RecordRequest
	Out     *api.DeleteRecordResponse
}

func NewDeleteRecordCommand(service *Service, in *api.RecordRequest) *DeleteRecordCommand {
	return &DeleteRecordCommand{
		service: service,
		In:      in,
	}
}

func (c *DeleteRecordCommand) Execute(ctx context.Context) error {
	rErr := c.service.checkForRecordRequestError(c.In)
	// Check if error is found
	if rErr != nil {
		c.Out = &api.DeleteRecordResponse{
			Success: false,
			Error:   conversion.Enum(*rErr, api.DeleteRecordResponse_Error_value, api.DeleteRecordResponse_NOT_FOUND),
		}
		return nil
	}
	result, err := c.service.database.DeleteRecord(ctx, model.GetRecordParams{
		Id:         conversion.Uint64ToSqlNullInt64(c.In.Id),
		NameUserId: convertNameUserIdToNullNameUserId(c.In.NameUserId),
	})
	// If an error occurs, it is an internal server error
	if err != nil {
		return err
	}
	// If no rows are affected, the record is not found
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
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
