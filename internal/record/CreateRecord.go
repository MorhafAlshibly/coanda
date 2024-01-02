package record

import (
	"context"
	"fmt"
	"strings"

	"github.com/MorhafAlshibly/coanda/api"
	"github.com/MorhafAlshibly/coanda/pkg/database/dynamoTable"
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
	// Check if record name is small enough
	if len(c.In.Name) > int(c.service.maxRecordNameLength) {
		c.Out = &api.CreateRecordResponse{
			Success: false,
			Error:   api.CreateRecordResponse_NAME_TOO_LONG,
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
	err := c.service.db.PutItem(ctx, &dynamoTable.PutItemInput{
		Item: map[string]any{
			"name":   c.In.Name,
			"userId": c.In.UserId,
			"record": c.In.Record,
			"data":   c.In.Data,
		},
		ConditionExpression: "attribute_not_exists(#name)",
		ExpressionAttributeNames: map[string]string{
			"#name": "name",
		},
	})
	if err != nil {
		fmt.Println(err)
		if strings.Contains("ConditionalCheckFailedException", err.Error()) {
			c.Out = &api.CreateRecordResponse{
				Success: false,
				Error:   api.CreateRecordResponse_RECORD_EXISTS,
			}
			return nil
		}
		return err
	}
	c.Out = &api.CreateRecordResponse{
		Success: true,
		Error:   api.CreateRecordResponse_NONE,
	}
	return nil
}
