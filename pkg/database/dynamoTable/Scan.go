package dynamoTable

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

type ScanInput struct {
	ExclusiveStartKey         map[string]string
	FilterExpression          string
	ProjectionExpression      string
	IndexName                 string
	Max                       uint8
	ExpressionAttributeNames  map[string]string
	ExpressionAttributeValues map[string]any
}

type ScanCommand struct {
	table *DynamoTable
	In    *ScanInput
	Out   []map[string]any
}

func (c *ScanCommand) Execute(ctx context.Context) error {
	exclusiveStartKey, err := attributevalue.MarshalMap(c.In.ExclusiveStartKey)
	if err != nil {
		return err
	}
	expressionAttributeValues, err := attributevalue.MarshalMap(c.In.ExpressionAttributeValues)
	if err != nil {
		return err
	}
	result, err := c.table.client.Scan(ctx, &dynamodb.ScanInput{
		TableName:                 aws.String(c.table.readTableName),
		ExclusiveStartKey:         NullifyEmptyStringMapAttributeValue(c.In.ExclusiveStartKey, exclusiveStartKey),
		FilterExpression:          NullifyEmptyString(aws.String(c.In.FilterExpression)),
		ProjectionExpression:      NullifyEmptyString(aws.String(c.In.ProjectionExpression)),
		IndexName:                 NullifyEmptyString(aws.String(c.In.IndexName)),
		Limit:                     aws.Int32(int32(c.In.Max)),
		ExpressionAttributeNames:  c.In.ExpressionAttributeNames,
		ExpressionAttributeValues: NullifyEmptyMapAttributeValue(c.In.ExpressionAttributeValues, expressionAttributeValues),
	})
	if err != nil {
		return err
	}
	return attributevalue.UnmarshalListOfMaps(result.Items, &c.Out)
}
