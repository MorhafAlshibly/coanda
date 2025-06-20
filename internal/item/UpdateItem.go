package item

import (
	"context"
	"database/sql"

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
			Error:   conversion.Enum(*iErr, api.UpdateItemResponse_Error_value, api.UpdateItemResponse_ID_REQUIRED),
		}
		return nil
	}
	if c.In.Data == nil {
		c.Out = &api.UpdateItemResponse{
			Success: false,
			Error:   api.UpdateItemResponse_DATA_REQUIRED,
		}
		return nil
	}
	data, err := conversion.ProtobufStructToRawJson(c.In.Data)
	if err != nil {
		return err
	}
	tx, err := c.service.sql.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	qtx := c.service.database.WithTx(tx)
	result, err := qtx.UpdateItem(ctx, model.UpdateItemParams{
		ID:   c.In.Item.Id,
		Type: c.In.Item.Type,
		Data: data,
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
		_, err := qtx.GetItem(ctx, model.GetItemParams{
			ID:   c.In.Item.Id,
			Type: c.In.Item.Type,
		})
		if err != nil {
			if err == sql.ErrNoRows {
				c.Out = &api.UpdateItemResponse{
					Success: false,
					Error:   api.UpdateItemResponse_NOT_FOUND,
				}
				return nil
			}
			return err
		}
	}
	err = tx.Commit()
	if err != nil {
		return err
	}
	c.Out = &api.UpdateItemResponse{
		Success: true,
		Error:   api.UpdateItemResponse_NONE,
	}
	return nil
}
