package dynamoTable

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

type DeleteItemInput struct {
	Key                       map[string]any
	ConditionExpression       string
	ExpressionAttributeNames  map[string]string
	ExpressionAttributeValues map[string]any
}

type DeleteItemCommand struct {
	table *DynamoTable
	In    *DeleteItemInput
}

func (c *DeleteItemCommand) Execute(ctx context.Context) error {
	key, err := attributevalue.MarshalMap(c.In.Key)
	if err != nil {
		return err
	}
	expressionAttributeValues, err := attributevalue.MarshalMap(c.In.ExpressionAttributeValues)
	if err != nil {
		return err
	}
	_, err = c.table.client.DeleteItem(ctx, &dynamodb.DeleteItemInput{
		Key:                       key,
		TableName:                 aws.String(c.table.tableName),
		ConditionExpression:       NullifyEmptyString(aws.String(c.In.ConditionExpression)),
		ExpressionAttributeNames:  c.In.ExpressionAttributeNames,
		ExpressionAttributeValues: NullifyEmptyMapAttributeValue(c.In.ExpressionAttributeValues, expressionAttributeValues),
	})
	return err
}
