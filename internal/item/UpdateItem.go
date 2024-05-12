package item

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"

	"github.com/MorhafAlshibly/coanda/api"
	"github.com/MorhafAlshibly/coanda/internal/item/model"
	"github.com/MorhafAlshibly/coanda/pkg/conversion"
)

type UpdateItemCommand struct {
	service *Service
	In      *api.UpdateItemRequest
	Out     *api.UpdateItemResponse
}

func NewUpdateItemCommand(service *Service, in *api.UpdateItemRequest) *UpdateItemCommand {
	return &UpdateItemCommand{
		service: service,
		In:      in,
	}
}

func (c *UpdateItemCommand) Execute(ctx context.Context) error {
	iErr := c.service.checkForItemRequestError(c.In.Item)
	if iErr != nil {
		c.Out = &api.UpdateItemResponse{
			Success: false,
			Error:   conversion.Enum(*iErr, api.UpdateItemResponse_Error_value, api.UpdateItemResponse_NOT_FOUND),
		}
		return nil
	}
	if c.In.ExpiresAt == nil && c.In.Data == nil {
		c.Out = &api.UpdateItemResponse{
			Success: false,
			Error:   api.UpdateItemResponse_NO_UPDATE_SPECIFIED,
		}
		return nil
	}
	var data json.RawMessage
	dataExists := int64(0)
	if c.In.Data != nil {
		var err error
		data, err = conversion.ProtobufStructToRawJson(c.In.Data)
		if err != nil {
			return err
		}
		dataExists = 1
	}
	result, err := c.service.database.UpdateItem(ctx, model.UpdateItemParams{
		ID:         c.In.Item.Id,
		Type:       c.In.Item.Type,
		Data:       data,
		DataExists: dataExists,
		ExpiresAt:  conversion.TimestampToSqlNullTime(c.In.ExpiresAt),
	})
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		// Check if item is found
		_, err := c.service.database.GetItem(ctx, model.GetItemParams{
			ID:   c.In.Item.Id,
			Type: c.In.Item.Type,
		})
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {

				c.Out = &api.UpdateItemResponse{
					Success: false,
					Error:   api.UpdateItemResponse_NOT_FOUND,
				}
				return nil
			}
			return err
		}
	}
	c.Out = &api.UpdateItemResponse{
		Success: true,
		Error:   api.UpdateItemResponse_NONE,
	}
	return nil
}
