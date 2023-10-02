package record

import (
	"context"

	"github.com/MorhafAlshibly/coanda/api"
	"go.mongodb.org/mongo-driver/bson"
)

type GetRecordCommand struct {
	service *Service
	In      *api.GetRecordRequest
	Out     *api.GetRecordResponse
}

func NewGetRecordCommand(service *Service, in *api.GetRecordRequest) *GetRecordCommand {
	return &GetRecordCommand{
		service: service,
		In:      in,
	}
}

func (c *GetRecordCommand) Execute(ctx context.Context) error {
	filter, err := getFilter(c.In)
	if err != nil {
		c.Out = &api.GetRecordResponse{
			Success: false,
			Record:  nil,
			Error:   api.GetRecordResponse_INVALID,
		}
		return nil
	}
	// Get the item from the store
	matchStage := bson.D{
		{Key: "$match", Value: filter},
	}
	cursor, err := c.service.db.Aggregate(ctx, append(pipeline, matchStage))
	if err != nil {
		return err
	}
	defer cursor.Close(ctx)
	cursor.Next(ctx)
	record, err := toRecord(cursor)
	if err != nil {
		if err.Error() == "EOF" {
			c.Out = &api.GetRecordResponse{
				Success: false,
				Record:  nil,
				Error:   api.GetRecordResponse_NOT_FOUND,
			}
			return nil
		}
		return err
	}
	c.Out = &api.GetRecordResponse{
		Success: true,
		Record:  record,
		Error:   api.GetRecordResponse_NONE,
	}
	return nil
}
