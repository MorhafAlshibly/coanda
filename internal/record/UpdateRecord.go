package record

import (
	"context"

	"github.com/MorhafAlshibly/coanda/api"
	"go.mongodb.org/mongo-driver/bson"
)

type UpdateRecordCommand struct {
	service *Service
	In      *api.UpdateRecordRequest
	Out     *api.UpdateRecordResponse
}

func NewUpdateRecordCommand(service *Service, in *api.UpdateRecordRequest) *UpdateRecordCommand {
	return &UpdateRecordCommand{
		service: service,
		In:      in,
	}
}

func (c *UpdateRecordCommand) Execute(ctx context.Context) error {
	filter, err := getFilter(c.In.Request)
	if err != nil {
		c.Out = &api.UpdateRecordResponse{
			Success: false,
			Error:   api.UpdateRecordResponse_INVALID,
		}
		return nil
	}
	if c.In.Request.NameUserId != nil {
		if len(c.In.Request.NameUserId.Name) < int(c.service.minRecordNameLength) {
			c.Out = &api.UpdateRecordResponse{
				Success: false,
				Error:   api.UpdateRecordResponse_NAME_TOO_SHORT,
			}
			return nil
		}
		if len(c.In.Request.NameUserId.Name) > int(c.service.maxRecordNameLength) {
			c.Out = &api.UpdateRecordResponse{
				Success: false,
				Error:   api.UpdateRecordResponse_NAME_TOO_LONG,
			}
			return nil
		}
	}
	update := bson.D{
		{Key: "$set", Value: bson.D{
			{Key: "data", Value: c.In.Data},
		}},
	}
	// Check if record is given
	if c.In.Record != nil {
		update = bson.D{
			{Key: "$set", Value: bson.D{
				{Key: "record", Value: *c.In.Record},
				{Key: "data", Value: c.In.Data},
			}},
		}
	}
	// Update the item in the store
	_, writeErr := c.service.db.UpdateOne(ctx, filter, update)
	if writeErr != nil {
		if writeErr.Error() == "EOF" {
			c.Out = &api.UpdateRecordResponse{
				Success: false,
				Error:   api.UpdateRecordResponse_NOT_FOUND,
			}
			return nil
		}
		return writeErr
	}
	c.Out = &api.UpdateRecordResponse{
		Success: true,
		Error:   api.UpdateRecordResponse_NONE,
	}
	return nil
}
