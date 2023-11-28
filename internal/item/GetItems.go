package item

import (
	"context"

	"github.com/MorhafAlshibly/coanda/api"
	"github.com/MorhafAlshibly/coanda/pkg"
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
	var filter map[string]any
	if c.In.Type != nil {
		if len(*c.In.Type) < int(c.service.minTypeLength) {
			c.Out = &api.GetItemsResponse{
				Success: false,
				Items:   nil,
				Error:   api.GetItemsResponse_TYPE_TOO_SHORT,
			}
			return nil
		}
		if len(*c.In.Type) > int(c.service.maxTypeLength) {
			c.Out = &api.GetItemsResponse{
				Success: false,
				Items:   nil,
				Error:   api.GetItemsResponse_TYPE_TOO_LONG,
			}
			return nil
		}
		filter = map[string]any{
			"PartitionKey": *c.In.Type,
		}
	}
	max, page := pkg.ParsePagination(c.In.Max, c.In.Page, c.service.defaultMaxPageLength, c.service.maxMaxPageLength)
	items, err := c.service.store.Query(ctx, filter, int32(max), int(page))
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
		Error:   api.GetItemsResponse_NONE,
	}
	return nil
}
