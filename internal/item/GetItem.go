package item

import (
	"context"
	"errors"

	"github.com/MorhafAlshibly/coanda/api"
	"github.com/MorhafAlshibly/coanda/pkg/storage"
)

type GetItemCommand struct {
	service *Service
	In      *api.GetItemRequest
	Out     *api.GetItemResponse
}

func NewGetItemCommand(service *Service, in *api.GetItemRequest) *GetItemCommand {
	return &GetItemCommand{
		service: service,
		In:      in,
	}
}

func (c *GetItemCommand) Execute(ctx context.Context) error {
	if c.In.Id == "" {
		c.Out = &api.GetItemResponse{
			Success: false,
			Item:    nil,
			Error:   api.GetItemResponse_ID_NOT_SET,
		}
		return nil
	}
	if len(c.In.Type) < int(c.service.minTypeLength) {
		c.Out = &api.GetItemResponse{
			Success: false,
			Item:    nil,
			Error:   api.GetItemResponse_TYPE_TOO_SHORT,
		}
		return nil
	}
	object, err := c.service.store.Get(ctx, c.In.Id, c.In.Type)
	if err != nil {
		if errors.Is(err, &storage.ObjectNotFoundError{}) {
			c.Out = &api.GetItemResponse{
				Success: false,
				Item:    nil,
				Error:   api.GetItemResponse_NOT_FOUND,
			}
			return nil
		}
		return err
	}
	item, err := objectToItem(object)
	if err != nil {
		return err
	}
	c.Out = &api.GetItemResponse{
		Success: true,
		Item:    item,
		Error:   api.GetItemResponse_NONE,
	}
	return nil
}
