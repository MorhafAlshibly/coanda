package record

import (
	"context"

	"github.com/MorhafAlshibly/coanda/api"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
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
	if c.In.NameUserId != nil {
		if len(c.In.NameUserId.Name) < int(c.service.minRecordNameLength) {
			c.Out = &api.GetRecordResponse{
				Success: false,
				Record:  nil,
				Error:   api.GetRecordResponse_NAME_TOO_SHORT,
			}
			return nil
		}
		if len(c.In.NameUserId.Name) > int(c.service.maxRecordNameLength) {
			c.Out = &api.GetRecordResponse{
				Success: false,
				Record:  nil,
				Error:   api.GetRecordResponse_NAME_TOO_LONG,
			}
			return nil
		}
	}
	// Get the item from the store
	pipelineWithMatch := mongo.Pipeline{
		bson.D{
			{Key: "$match", Value: filter},
		},
	}
	for _, stage := range pipeline {
		pipelineWithMatch = append(pipelineWithMatch, stage)
	}
	cursor, err := c.service.db.Aggregate(ctx, pipelineWithMatch)
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
