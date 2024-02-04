package dynamoTable

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

type QueryInput struct {
	ExclusiveStartKey         map[string]string
	KeyConditionExpression    string
	FilterExpression          string
	ProjectionExpression      string
	IndexName                 string
	Order                     bool
	Max                       uint8
	ExpressionAttributeNames  map[string]string
	ExpressionAttributeValues map[string]any
}

type QueryCommand struct {
	table *DynamoTable
	In    *QueryInput
	Out   []map[string]any
}

func (c *QueryCommand) Execute(ctx context.Context) error {
	exclusiveStartKey, err := attributevalue.MarshalMap(c.In.ExclusiveStartKey)
	if err != nil {
		return err
	}
	expressionAttributeValues, err := attributevalue.MarshalMap(c.In.ExpressionAttributeValues)
	if err != nil {
		return err
	}
	result, err := c.table.client.Query(ctx, &dynamodb.QueryInput{
		TableName:                 aws.String(c.table.readTableName),
		ExclusiveStartKey:         NullifyEmptyStringMapAttributeValue(c.In.ExclusiveStartKey, exclusiveStartKey),
		KeyConditionExpression:    NullifyEmptyString(aws.String(c.In.KeyConditionExpression)),
		FilterExpression:          NullifyEmptyString(aws.String(c.In.FilterExpression)),
		ProjectionExpression:      NullifyEmptyString(aws.String(c.In.ProjectionExpression)),
		IndexName:                 NullifyEmptyString(aws.String(c.In.IndexName)),
		ScanIndexForward:          aws.Bool(c.In.Order),
		Limit:                     aws.Int32(int32(c.In.Max)),
		ExpressionAttributeNames:  c.In.ExpressionAttributeNames,
		ExpressionAttributeValues: NullifyEmptyMapAttributeValue(c.In.ExpressionAttributeValues, expressionAttributeValues),
	})
	if err != nil {
		return err
	}
	return attributevalue.UnmarshalListOfMaps(result.Items, &c.Out)
}
