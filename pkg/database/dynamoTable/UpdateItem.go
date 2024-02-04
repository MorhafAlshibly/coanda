package dynamoTable

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

type UpdateItemInput struct {
	Key                       map[string]string
	ConditionExpression       string
	UpdateExpression          string
	ExpressionAttributeNames  map[string]string
	ExpressionAttributeValues map[string]any
}

type UpdateItemCommand struct {
	table *DynamoTable
	In    *UpdateItemInput
}

func (c *UpdateItemCommand) Execute(ctx context.Context) error {
	key, err := attributevalue.MarshalMap(c.In.Key)
	if err != nil {
		return err
	}
	expressionAttributeValues, err := attributevalue.MarshalMap(c.In.ExpressionAttributeValues)
	if err != nil {
		return err
	}
	_, err = c.table.client.UpdateItem(ctx, &dynamodb.UpdateItemInput{
		Key:                       key,
		TableName:                 aws.String(c.table.tableName),
		ConditionExpression:       NullifyEmptyString(aws.String(c.In.ConditionExpression)),
		UpdateExpression:          NullifyEmptyString(aws.String(c.In.UpdateExpression)),
		ExpressionAttributeNames:  c.In.ExpressionAttributeNames,
		ExpressionAttributeValues: NullifyEmptyMapAttributeValue(c.In.ExpressionAttributeValues, expressionAttributeValues),
	})
	return err
}
