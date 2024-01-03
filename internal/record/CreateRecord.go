package record

import (
	"context"

	"github.com/MorhafAlshibly/coanda/api"
	"github.com/MorhafAlshibly/coanda/pkg/conversion"
	"github.com/MorhafAlshibly/coanda/pkg/database/sqlc"
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
	if c.In.Record == 0 {
		c.Out = &api.CreateRecordResponse{
			Success: false,
			Error:   api.CreateRecordResponse_RECORD_REQUIRED,
		}
		return nil
	}
	data, err := conversion.ProtobufStructToRawJson(c.In.Data)
	if err != nil {
		return err
	}
	// Insert the record into the database
	_, err = c.service.database.CreateRecord(ctx, sqlc.CreateRecordParams{
		Name:   c.In.Name,
		UserID: c.In.UserId,
		Record: c.In.Record,
		Data:   data,
	})
	// Check if the record already exists
	if err.(*mysql.MySQLError).Number == sqlc.ERR_DUP_ENTRY {
		c.Out = &api.CreateRecordResponse{
			Success: false,
			Error:   api.CreateRecordResponse_RECORD_EXISTS,
		}
		return nil
	}
	if err != nil {
		return err
	}
	c.Out = &api.CreateRecordResponse{
		Success: true,
		Error:   api.CreateRecordResponse_NONE,
	}
	return nil
}
