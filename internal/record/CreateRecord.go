package record

import (
	"context"

	"github.com/MorhafAlshibly/coanda/api"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type CreateRecordCommand struct {
	service *Service
	In      *api.CreateRecordRequest
	Out     *api.CreateRecordResponse
}

func NewCreateRecordCommand(service *Service, in *api.CreateRecordRequest) *CreateRecordCommand {
	return &CreateRecordCommand{
		service: service,
		In:      in,
	}
}

func (c *CreateRecordCommand) Execute(ctx context.Context) error {
	// Check if record name is large enough
	if len(c.In.Name) < int(c.service.minRecordNameLength) {
		c.Out = &api.CreateRecordResponse{
			Success: false,
			Error:   api.CreateRecordResponse_NAME_TOO_SHORT,
		}
		return nil
	}
	// Check if user id is valid
	if c.In.UserId == 0 {
		c.Out = &api.CreateRecordResponse{
			Success: false,
			Error:   api.CreateRecordResponse_USER_ID_REQUIRED,
		}
		return nil
	}
	if c.In.Record == 0 {
		c.Out = &api.CreateRecordResponse{
			Success: false,
			Error:   api.CreateRecordResponse_RECORD_REQUIRED,
		}
		return nil
	}
	// Insert the record into the database
	id, writeErr := c.service.db.InsertOne(ctx, bson.D{
		{Key: "name", Value: c.In.Name},
		{Key: "userId", Value: c.In.UserId},
		{Key: "record", Value: c.In.Record},
		{Key: "data", Value: c.In.Data},
	})
	if writeErr != nil {
		if mongo.IsDuplicateKeyError(writeErr) {
			c.Out = &api.CreateRecordResponse{
				Success: false,
				Error:   api.CreateRecordResponse_RECORD_EXISTS,
			}
			return nil
		}
		return writeErr
	}
	c.Out = &api.CreateRecordResponse{
		Success: true,
		Id:      id.Hex(),
		Error:   api.CreateRecordResponse_NONE,
	}
	return nil
}
