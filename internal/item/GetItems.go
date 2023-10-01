package item

import (
	"context"

	"github.com/MorhafAlshibly/coanda/api"
	"github.com/MorhafAlshibly/coanda/pkg/storage"
)

type GetItemsCommand struct {
	service *Service
	In      *api.GetItemsRequest
	Out     *api.GetItemsResponse
}

func NewGetItemsCommand(service *Service, in *api.GetItemsRequest) *GetItemsCommand {
	return &GetItemsCommand{
		service: service,
		In:      in,
	}
}

func (c *GetItemsCommand) Execute(ctx context.Context) error {
	var items []*storage.Object
	var outs []*api.Item
	// If the type is not nil, set the filter to the type
	filter := ""
	if c.In.Type != nil {
		filter = "PartitionKey eq '" + *c.In.Type + "'"
	}
	if c.In.Max == nil {
		c.In.Max = new(uint64)
		*c.In.Max = c.service.defaultMaxPageLength
	}
	if c.In.Page == nil {
		c.In.Page = new(uint64)
		*c.In.Page = 1
	}
	items, err := c.service.store.Query(ctx, filter, int32(*c.In.Max), int(*c.In.Page))
	if err != nil {
		return err
	}
	for _, item := range items {
		out, err := objectToItem(item)
		if err != nil {
			return err
		}
		outs = append(outs, out)
	}
	c.Out = &api.GetItemsResponse{
		Success: true,
		Items:   outs,
	}
	return nil
}
