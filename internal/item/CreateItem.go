package item

import (
	"context"
	"errors"
	"time"

	"github.com/MorhafAlshibly/coanda/api"
	"github.com/MorhafAlshibly/coanda/internal/item/model"
	"github.com/MorhafAlshibly/coanda/pkg/conversion"
	"github.com/MorhafAlshibly/coanda/pkg/errorcodes"
	"github.com/MorhafAlshibly/coanda/pkg/validation"
	"github.com/go-sql-driver/mysql"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type CreateItemCommand struct {
	service *Service
	In      *api.CreateItemRequest
	Out     *api.CreateItemResponse
}

func NewCreateItemCommand(service *Service, in *api.CreateItemRequest) *CreateItemCommand {
	return &CreateItemCommand{
		service: service,
		In:      in,
	}
}

func (c *CreateItemCommand) Execute(ctx context.Context) error {
	if c.In.Id == "" {
		c.Out = &api.CreateItemResponse{
			Success: false,
			Error:   api.CreateItemResponse_ID_REQUIRED,
		}
		return nil
	}
	if c.In.Type == "" {
		c.Out = &api.CreateItemResponse{
			Success: false,
			Error:   api.CreateItemResponse_TYPE_REQUIRED,
		}
		return nil
	}
	if c.In.Data == nil {
		c.Out = &api.CreateItemResponse{
			Success: false,
			Error:   api.CreateItemResponse_DATA_REQUIRED,
		}
		return nil
	}
	if c.In.ExpiresAt == nil {
		// If the item has no expiry, set it to the Unix epoch time (0)
		c.In.ExpiresAt = timestamppb.New(time.Unix(0, 0))
	}
	data, err := conversion.ProtobufStructToRawJson(c.In.Data)
	if err != nil {
		return err
	}
	_, err = c.service.database.CreateItem(ctx, model.CreateItemParams{
		ID:        c.In.Id,
		Type:      c.In.Type,
		Data:      data,
		ExpiresAt: validation.ValidateATimestampToSqlNullTime(c.In.ExpiresAt),
	})
	// Check if the item already exists
	if err != nil {
		var mysqlErr *mysql.MySQLError
		if errors.As(err, &mysqlErr) && mysqlErr.Number == errorcodes.MySQLErrorCodeDuplicateEntry {
			c.Out = &api.CreateItemResponse{
				Success: false,
				Error:   api.CreateItemResponse_ALREADY_EXISTS,
			}
			return nil
		}
		return err
	}
	c.Out = &api.CreateItemResponse{
		Success: true,
		Error:   api.CreateItemResponse_NONE,
	}
	return nil
}
