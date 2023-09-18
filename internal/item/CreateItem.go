package item

import (
	"context"
	"time"

	"github.com/MorhafAlshibly/coanda/api"
	"github.com/bytedance/sonic"
)

type CreateItemCommand struct {
	service *Service
	In      *api.CreateItemRequest
	Out     *api.Item
}

func NewCreateItemCommand(service *Service, in *api.CreateItemRequest) *CreateItemCommand {
	return &CreateItemCommand{
		service: service,
		In:      in,
	}
}

func (c *CreateItemCommand) Execute(ctx context.Context) error {
	marshalled, err := sonic.Marshal(c.In.Data)
	if err != nil {
		return err
	}
	mapData := map[string]string{
		"Type": c.In.Type,
		"Data": string(marshalled),
	}
	// If the item has an expiry, add it to the map
	if c.In.Expire != nil {
		mapData["Expire"] = c.In.Expire.AsTime().Format(time.RFC3339)
	}
	// Add the item to the store
	object, err := c.service.store.Add(ctx, c.In.Type, mapData)
	if err != nil {
		return err
	}
	// Allot the output
	c.Out = &api.Item{
		Id:     object.Key,
		Type:   c.In.Type,
		Data:   c.In.Data,
		Expire: c.In.Expire,
	}
	return nil
}
