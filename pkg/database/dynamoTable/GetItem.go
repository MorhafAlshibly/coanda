package dynamoTable

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

type GetItemInput struct {
	Key                      map[string]any
	ProjectionExpression     string
	ExpressionAttributeNames map[string]string
}

type GetItemCommand struct {
	table *DynamoTable
	In    *GetItemInput
	Out   map[string]any
}

type ItemNotFoundError struct{}

func (e *ItemNotFoundError) Error() string {
	return "Item not found"
}

func (c *GetItemCommand) Execute(ctx context.Context) error {
	key, err := attributevalue.MarshalMap(c.In.Key)
	if err != nil {
		return err
	}
	result, err := c.table.client.GetItem(ctx, &dynamodb.GetItemInput{
		TableName:                aws.String(c.table.readTableName),
		Key:                      key,
		ProjectionExpression:     NullifyEmptyString(aws.String(c.In.ProjectionExpression)),
		ExpressionAttributeNames: c.In.ExpressionAttributeNames,
	})
	if err != nil {
		return err
	}
	if result.Item == nil {
		return &ItemNotFoundError{}
	}
	return attributevalue.UnmarshalMap(result.Item, &c.Out)
}
