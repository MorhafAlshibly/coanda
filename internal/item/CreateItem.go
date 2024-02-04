package item

import (
	"context"
	"time"

	"github.com/MorhafAlshibly/coanda/api"
	"github.com/MorhafAlshibly/coanda/pkg/database/dynamoTable"
	"github.com/google/uuid"
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
	if len(c.In.Type) < int(c.service.minTypeLength) {
		c.Out = &api.CreateItemResponse{
			Success: false,
			Error:   api.CreateItemResponse_TYPE_TOO_SHORT,
		}
		return nil
	}
	if len(c.In.Type) > int(c.service.maxTypeLength) {
		c.Out = &api.CreateItemResponse{
			Success: false,
			Error:   api.CreateItemResponse_TYPE_TOO_LONG,
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
	if c.In.ExpireAt == nil {
		// If the item has no expiry, set it to the Unix epoch time (0)
		c.In.ExpireAt = timestamppb.New(time.Unix(0, 0))
	}
	id := uuid.New().String()
	item := &api.Item{
		Id:       id,
		Type:     c.In.Type,
		Data:     c.In.Data,
		ExpireAt: c.In.ExpireAt,
	}
	// Add the item to the store
	object, err := MarshalItem(item)
	if err != nil {
		return err
	}
	err = c.service.database.PutItem(ctx, &dynamoTable.PutItemInput{
		Item: object,
	})
	if err != nil {
		return err
	}
	// Allot the output
	c.Out = &api.CreateItemResponse{
		Success: true,
		Item:    item,
	}
	return nil
}
