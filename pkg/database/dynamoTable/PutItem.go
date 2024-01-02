package dynamoTable

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

type PutItemInput struct {
	Item                      map[string]any
	ConditionExpression       string
	ExpressionAttributeNames  map[string]string
	ExpressionAttributeValues map[string]any
}

type PutItemCommand struct {
	table *DynamoTable
	In    *PutItemInput
}

func (c *PutItemCommand) Execute(ctx context.Context) error {
	item, err := attributevalue.MarshalMap(c.In.Item)
	if err != nil {
		return err
	}
	expressionAttributeValues, err := attributevalue.MarshalMap(c.In.ExpressionAttributeValues)
	if err != nil {
		return err
	}
	_, err = c.table.client.PutItem(ctx, &dynamodb.PutItemInput{
		Item:                      item,
		TableName:                 aws.String(c.table.tableName),
		ConditionExpression:       NullifyEmptyString(aws.String(c.In.ConditionExpression)),
		ExpressionAttributeNames:  c.In.ExpressionAttributeNames,
		ExpressionAttributeValues: NullifyEmptyMapAttributeValue(c.In.ExpressionAttributeValues, expressionAttributeValues),
	})
	return err
}
