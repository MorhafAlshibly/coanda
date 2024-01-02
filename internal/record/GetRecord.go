package record

import (
	"context"
	"errors"

	"github.com/MorhafAlshibly/coanda/api"
	"github.com/MorhafAlshibly/coanda/pkg/database/dynamoTable"
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
	if len(c.In.Name) < int(c.service.minRecordNameLength) {
		c.Out = &api.GetRecordResponse{
			Success: false,
			Record:  nil,
			Error:   api.GetRecordResponse_NAME_TOO_SHORT,
		}
		return nil
	}
	if len(c.In.Name) > int(c.service.maxRecordNameLength) {
		c.Out = &api.GetRecordResponse{
			Success: false,
			Record:  nil,
			Error:   api.GetRecordResponse_NAME_TOO_LONG,
		}
		return nil
	}
	if c.In.UserId == 0 {
		c.Out = &api.GetRecordResponse{
			Success: false,
			Record:  nil,
			Error:   api.GetRecordResponse_USER_ID_REQUIRED,
		}
		return nil
	}
	object, err := c.service.db.GetItem(ctx, &dynamoTable.GetItemInput{
		Key: map[string]any{
			"name":   c.In.Name,
			"userId": c.In.UserId,
		},
	})
	if err != nil {
		if errors.Is(err, &dynamoTable.ItemNotFoundError{}) {
			c.Out = &api.GetRecordResponse{
				Success: false,
				Record:  nil,
				Error:   api.GetRecordResponse_NOT_FOUND,
			}
			return nil
		}
		return err
	}
	// Get rank by getting count of all records with a faster time than the current record
	rank, err := c.service.db.Query(ctx, &dynamoTable.QueryInput{
		KeyConditionExpression: "#name = :name",
		FilterExpression:       "#record < :record",
		ExpressionAttributeNames: map[string]string{
			"#name":   "name",
			"#record": "record",
		},
		ExpressionAttributeValues: map[string]any{
			":name":   c.In.Name,
			":record": object["record"],
		},
	})
	if err != nil {
		return err
	}
	record, err := UnmarshalRecord(object)
	if err != nil {
		return err
	}
	c.Out = &api.GetRecordResponse{
		Success: true,
		Record:  record,
		Error:   api.GetRecordResponse_NONE,
	}
	return nil
}
