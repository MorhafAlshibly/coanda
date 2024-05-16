package record

import (
	"context"
	"errors"

	"github.com/MorhafAlshibly/coanda/api"
	"github.com/MorhafAlshibly/coanda/internal/record/model"
	"github.com/MorhafAlshibly/coanda/pkg/conversion"
	"github.com/MorhafAlshibly/coanda/pkg/errorcodes"
	"github.com/go-sql-driver/mysql"
)

type CreateRecordCommand struct {
	service *Service
	In      *api.CreateRecordRequest
	Out     *api.CreateRecordResponse
}

func NewCreateRecordCommand(service *Service, in *api.CreateRecordRequest) *CreateRecordCommand {
	return &CreateRecordCommand{
		service: service,
		In:      in,
	}
}

func (c *CreateRecordCommand) Execute(ctx context.Context) error {
	// Check if record name is large enough
	if len(c.In.Name) < int(c.service.minRecordNameLength) {
		c.Out = &api.CreateRecordResponse{
			Success: false,
			Error:   api.CreateRecordResponse_NAME_TOO_SHORT,
		}
		return nil
	}
	// Check if record name is small enough
	if len(c.In.Name) > int(c.service.maxRecordNameLength) {
		c.Out = &api.CreateRecordResponse{
			Success: false,
			Error:   api.CreateRecordResponse_NAME_TOO_LONG,
		}
		return nil
	}
	// Check if user id is valid
	if c.In.UserId == 0 {
		c.Out = &api.CreateRecordResponse{
			Success: false,
			Error:   api.CreateRecordResponse_USER_ID_REQUIRED,
		}
		return nil
	}
	// Check if record is valid
	if c.In.Record == 0 {
		c.Out = &api.CreateRecordResponse{
			Success: false,
			Error:   api.CreateRecordResponse_RECORD_REQUIRED,
		}
		return nil
	}
	// Check if data is provided
	if c.In.Data == nil {
		c.Out = &api.CreateRecordResponse{
			Success: false,
			Error:   api.CreateRecordResponse_DATA_REQUIRED,
		}
		return nil
	}
	data, err := conversion.ProtobufStructToRawJson(c.In.Data)
	if err != nil {
		return err
	}
	// Insert the record into the database
	result, err := c.service.database.CreateRecord(ctx, model.CreateRecordParams{
		Name:   c.In.Name,
		UserID: c.In.UserId,
		Record: c.In.Record,
		Data:   data,
	})
	// Check if the record already exists
	if err != nil {
		var mysqlError *mysql.MySQLError
		if errors.As(err, &mysqlError) && mysqlError.Number == errorcodes.MySQLErrorCodeDuplicateEntry {
			c.Out = &api.CreateRecordResponse{
				Success: false,
				Error:   api.CreateRecordResponse_RECORD_EXISTS,
			}
			return nil
		}
		return err
	}
	// Get the id of the record
	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	c.Out = &api.CreateRecordResponse{
		Success: true,
		Error:   api.CreateRecordResponse_NONE,
		Id:      uint64(id),
	}
	return nil
}
