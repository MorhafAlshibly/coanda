package item

import (
	"context"

	"github.com/MorhafAlshibly/coanda/api/gql"
)

type GetItemCommand struct {
	service *Service
	In      *gql.GetItem
	Out     *gql.Item
}

func NewGetItemCommand(service *Service, in *gql.GetItem) *GetItemCommand {
	return &GetItemCommand{
		service: service,
		In:      in,
	}
}

func (c *GetItemCommand) Execute(ctx context.Context) error {
	object, err := c.service.store.Get(ctx, c.In.ID, c.In.Type)
	if err != nil {
		return err
	}
	c.Out, err = objectToItem(object)
	if err != nil {
		return err
	}
	return nil
}
